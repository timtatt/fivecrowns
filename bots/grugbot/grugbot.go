package grugbot

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/timtatt/fivecrowns/bots"
	"github.com/timtatt/fivecrowns/game"
	"github.com/timtatt/fivecrowns/math"
)

// Grugbot tries to put cards in the biggest possible sequences
// if a card is used in multiple sequences, the one with the biggest score is preferred

type grugBot struct{}

func NewGrugBot() bots.Bot {
	return &grugBot{}
}

type Calculation struct {
	Sequences [][]game.Card
	Flop      bool
	Hand      []game.Card
}

func (b *grugBot) Score(req bots.BotRequest) (bots.ScoreResponse, error) {

	calculation, err := b.Calculate(req)

	if err != nil {
		return bots.ScoreResponse{}, fmt.Errorf("cannot calculate response: %w", err)
	}

	return bots.ScoreResponse{
		Action:    req.Action,
		Flop:      game.CanFlop(calculation.Sequences),
		Sequences: game.EncodeSequences(calculation.Sequences),
	}, nil
}

func (b *grugBot) Draw(req bots.BotRequest) (bots.DrawResponse, error) {

	// add the discard to the hand and determine if it gets added to a valid sequence

	topCard, err := game.DecodeCard(req.Discard[0])

	if err != nil {
		return bots.DrawResponse{}, fmt.Errorf("unable to decode discard card: %w", err)
	}

	hypothetical, err := b.Calculate(bots.BotRequest{
		Action:  req.Action,
		Hand:    append(req.Hand, req.Discard[0]),
		Round:   req.Round,
		Discard: req.Discard,
	})

	if err != nil {
		return bots.DrawResponse{}, fmt.Errorf("unable to hypothesise discard score: %w", err)
	}

	// determine if the topCard has been used in a sequence
	// goes in reverse and checks if there is an invalid sequence with only the topCard
	for i := len(hypothetical.Sequences) - 1; i >= 0; i-- {
		// available optimisation: if sequences are now valid, the top card will be there somewhere. we don't need to know where
		if len(hypothetical.Sequences[i]) >= 3 {
			return bots.DrawResponse{
				Action: req.Action,
				Stack:  bots.StackDiscard,
			}, nil
		}

		// these sequences are invalid
		if slices.Contains(hypothetical.Sequences[i], topCard) {
			return bots.DrawResponse{
				Action: req.Action,
				Stack:  bots.StackDeck,
			}, nil
		}
	}

	return bots.DrawResponse{}, errors.New("did not find the top card in the sequences")
}

func (b *grugBot) Discard(req bots.BotRequest) (bots.DiscardResponse, error) {

	calculation, err := b.Calculate(req)

	if err != nil {
		return bots.DiscardResponse{}, fmt.Errorf("unable to calculate score: %w", err)
	}

	// determine which card is the highest one that is not in a valid sequence

	// save the card and its location to easily remove it in the future
	type cardAndLocation struct {
		card        game.Card
		sequenceIdx int
		cardIdx     int
	}

	var worstCard cardAndLocation

	for i := len(calculation.Sequences) - 1; i >= 0; i-- {
		seq := calculation.Sequences[i]

		// if it is the last turn, we just want to throw out our highest card, even if it is in a partial sequence
		threshold := 2
		if req.LastTurn {
			threshold = 3
		}

		// skip seq if above keeping threshold
		if len(seq) >= threshold {
			continue
		}

		// get the highest card in the current sequence
		for j, card := range seq {
			if game.ScoreCard(card) > game.ScoreCard(worstCard.card) {
				worstCard = cardAndLocation{
					card:        card,
					sequenceIdx: i,
					cardIdx:     j,
				}
			}
		}
	}

	// if the hand can flop without discarding, a random card will need to be chosen to be omitted
	// TODO: edge case for round 5,8,11 which may break a sequence
	// TODO: for round 5,8,11 if a sequence is broken, ensure it is the least troublesome
	if worstCard.card.IsNil() {
		for i, seq := range calculation.Sequences {
			if len(seq) > 3 {
				// remove first a non-wild card
				for j, card := range seq {
					if !card.IsWild(req.Round) {
						worstCard = cardAndLocation{
							card:        card,
							sequenceIdx: i,
							cardIdx:     j,
						}

						break
					}
				}
			}
		}
	}

	slog.Info("worst card detected", "card", worstCard)

	// update the discard response to omit the discarded card

	seq := calculation.Sequences[worstCard.sequenceIdx]

	if len(seq) == 1 {
		// the sequence only has one card, delete the whole thing
		calculation.Sequences = slices.Delete(calculation.Sequences, worstCard.sequenceIdx, worstCard.sequenceIdx+1)
	} else {
		// remove the card from the its sequence
		calculation.Sequences[worstCard.sequenceIdx] = slices.Delete(calculation.Sequences[worstCard.sequenceIdx], worstCard.cardIdx, worstCard.cardIdx+1)
	}

	return bots.DiscardResponse{
		Flop:      game.CanFlop(calculation.Sequences),
		Sequences: game.EncodeSequences(calculation.Sequences),
		Action:    bots.ActionDiscard,
		Card:      worstCard.card.Encode(),
	}, nil
}

// calculate best possible sequences
func (b *grugBot) Calculate(req bots.BotRequest) (Calculation, error) {
	hand, err := game.DecodeCards(req.Hand)

	if err != nil {
		return Calculation{}, fmt.Errorf("unable to decode cards: %w", err)
	}

	// sort the cards by suite and number
	slices.SortFunc(hand, game.CompareCard)

	// find all possible sequences in the hand
	seqs := FindSequences(req.Round, hand)

	slog.Info("calculated possible sequences", "seqs", game.EncodeSequences(seqs))

	// filter out sequences if they have cards that have been used twiced
	// use the wilds to build more sequences
	seqs = FilterSequences(req.Round, hand, seqs)

	return Calculation{
		Flop:      game.CanFlop(seqs),
		Sequences: seqs,
	}, nil
}

// takes a list of cards and returns a map with the counts of each card
func CardCounts(hand []game.Card) map[game.Card]int {
	cardCounts := make(map[game.Card]int)

	for _, card := range hand {
		cardCounts[card] += 1
	}

	return cardCounts
}

// takes a sorted list of cards and returns a list of possible sequences
// note: a single card may be used multiple times
func FindSequences(round int, hand []game.Card) [][]game.Card {

	seqs := make([][]game.Card, 0)

	// go through the sorted cards and find runs of numbers in the same suite

	curSeq := []game.Card{hand[0]}

	for i := 1; i < len(hand); i++ {

		// check if the current number is in sequence with previous

		c := hand[i]
		pc := curSeq[len(curSeq)-1]

		// if we are dealing with jokers, we can ignore them
		// if the current card and previous card in sequence is the same, we can ignore it
		if c.IsWild(round) || c == pc {
			continue
		}

		// check if cards are sequential or have a single gap
		if c.Suite == pc.Suite && c.Number == pc.Number+1 {
			// cards are in sequence, add this one to the sequence
			curSeq = append(curSeq, c)
		} else {
			// cards are not in a complete run

			if c.Suite == pc.Suite && c.Number == pc.Number+2 && len(curSeq) == 1 {
				// if the sequence has a gap, save if length is < 2
				curSeq = append(curSeq, c)
			}

			// save the sequence if it qualifies
			if len(curSeq) >= 2 {
				slog.Info("adding sequence", "seq", game.EncodeSequence(curSeq))
				seqs = append(seqs, curSeq)
			}

			// reset the current sequence
			curSeq = []game.Card{c}
		}
	}

	// ensure we don't leave a curSeq hanging
	if len(curSeq) >= 2 {
		seqs = append(seqs, curSeq)
		slog.Info("adding sequence", "seq", game.EncodeSequence(curSeq))
	}

	cardCounts := CardCounts(hand)
	slog.Info("card counts in hand", "hand", cardCounts)

	// check the cards for sets
	for number := 3; number <= 13; number++ {

		// do not get a collection set of wilds
		if number == round {
			continue
		}

		curSeq := make([]game.Card, 0)

		for _, suite := range game.Suites {

			c := game.Card{
				Joker:  false,
				Number: number,
				Suite:  suite,
			}

			for range cardCounts[c] {
				curSeq = append(curSeq, c)
			}

		}

		if len(curSeq) >= 2 {
			// add all single cards as a sequence too
			seqs = append(seqs, curSeq)
			slog.Info("adding sequence", "seq", game.EncodeSequence(curSeq))
		}
	}

	// add individual cards as seqs
	for _, card := range hand {
		// do not get seqeunces of the wilds
		if !card.IsWild(round) {
			seqs = append(seqs, []game.Card{card})
		}

	}

	return seqs
}

func CompareSequence(round int) func(a, b []game.Card) int {
	return func(a, b []game.Card) int {
		// sorts in reverse to optimise the sequence processing

		// 1 = prefer a
		// -1 = prefer b

		if len(a) < 3 && len(b) >= 3 {
			return -1
		} else if len(a) >= 3 && len(b) < 3 {
			return 1
		}

		scoreDiff := game.ScoreSequence(a) - game.ScoreSequence(b)

		if scoreDiff != 0 {
			return scoreDiff
		}

		// prefer sets to runs, when everything else is equal
		aType := game.GetSequenceType(a, round)
		bType := game.GetSequenceType(b, round)

		if aType == bType {
			return 0
		} else if aType == game.SequenceTypeRun {
			return -1
		} else if bType == game.SequenceTypeRun {
			return 1
		}

		return 0
	}

}

// takes a hand and a list of sequences
// determines if any card is being used more than it should be
// if so, will chose the sequence with the highest score
// output will include any cards which dont fit within a sequence as a single-carded sequence
func FilterSequences(round int, hand []game.Card, seqs [][]game.Card) [][]game.Card {

	// score all of the sequences

	slices.SortFunc(seqs, CompareSequence(round))

	slog.Info("sorting by sequence scores", "seqs", game.EncodeSequences(seqs))

	cardCounts := CardCounts(hand)

	// keeps track of which sequences we can still try out
	remainingSeqs := slices.Clone(seqs)

	filteredSeqs := make([][]game.Card, 0)

	z := 0

	for len(remainingSeqs) > 0 {
		lastIdx := len(remainingSeqs) - 1

		// remainingSeqs will always be sorted by score
		// always use the top sequence (which is at the end)
		seq := remainingSeqs[lastIdx]

		z += 1

		if z == 100 {
			slog.Info("infinite loop", "filteredSeqs", game.EncodeSequences(filteredSeqs), "remainingSeqs", game.EncodeSequences(remainingSeqs))
			panic("infinite loop detected")
		}

		// available optimisation: don't bother checking for card usage for first sequence
		if len(filteredSeqs) != 0 {

			// check all the cards are available
			availableCards := make([]game.Card, 0)
			for card, reqCount := range CardCounts(seq) {
				for range math.Min(reqCount, cardCounts[card]) {
					availableCards = append(availableCards, card)
				}
			}

			if len(availableCards) == 0 {
				// if there are no cards available, we need to remove this as a possible sequence
				remainingSeqs = slices.Delete(remainingSeqs, lastIdx, lastIdx+1)

				continue
			} else if len(availableCards) != len(seq) {
				slog.Info("missing cards to build sequence", "needs", game.EncodeCards(seq), "available", game.EncodeCards(availableCards))
				// there are some missing cards, this seq cannot be used as is

				// add this remaining seq back into the remainingSeqs
				slices.SortFunc(availableCards, game.CompareCard)
				remainingSeqs[lastIdx] = availableCards

				// resort the remainingSeqs
				slices.SortFunc(remainingSeqs, CompareSequence(round))

				continue
			}
		}

		// add jokers to the sequence if the sequence is < 3 cards in length
		gap := 3 - len(seq)
		if gap > 0 && wildCount(cardCounts, round) >= gap {
			for range gap {

				// fetch an available wild card
				// decrement the wildcard after usage

				// ignoring the err, we know there will be a wild available here
				wc, err := getWild(cardCounts, round)
				slog.Info("add a wild to seq", "seq", game.EncodeSequence(seq), "wild", wc)

				if err != nil {
					// this should not occur
					panic(err)
				}

				seq = append(seq, wc)
				cardCounts[wc] -= 1
			}
		}

		filteredSeqs = append(filteredSeqs, seq)

		// decrement remaining cards which have been used
		for _, card := range seq {
			// do not decrement wilds, they have been decremented as above
			if !card.IsWild(round) {
				cardCounts[card] -= 1
			}
		}

		// remove the sequence from remaining seqs
		remainingSeqs = slices.Delete(remainingSeqs, lastIdx, lastIdx+1)
	}

	slog.Info("remaining card counts", "counts", cardCounts)

	for {
		wc, err := getWild(cardCounts, round)
		if err != nil {
			break
		}

		// place the remaining jokers in a sequence
		if len(filteredSeqs) > 0 {
			filteredSeqs[0] = append(filteredSeqs[0], wc)
		} else {
			// handle case when only wilds are in the hand
			filteredSeqs = append(filteredSeqs, []game.Card{wc})
		}

		// decrement the wild
		cardCounts[wc] -= 1
	}

	return filteredSeqs
}

func getWild(cardCounts map[game.Card]int, round int) (game.Card, error) {

	if cardCounts[game.CardJoker] > 0 {
		return game.CardJoker, nil
	}

	for _, s := range game.Suites {
		c := game.Card{
			Number: round,
			Suite:  s,
		}

		if cardCounts[c] > 0 {
			return c, nil
		}
	}

	return game.Card{}, errors.New("no wilds found")
}

func wildCount(cardCounts map[game.Card]int, round int) int {

	count := 0

	count += cardCounts[game.CardJoker]

	for _, s := range game.Suites {
		c := game.Card{
			Number: round,
			Suite:  s,
		}

		count += cardCounts[c]
	}

	return count
}
