package grugbot

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/timtatt/fivecrowns/bots"
	"github.com/timtatt/fivecrowns/game"
)

// Grugbot tries to put cards in the biggest possible sequences
// if a card is used in multiple sequences, the one with the biggest score is preferred

type grugBot struct{}

func NewGrugBot() bots.Bot {
	return &grugBot{}
}

func (b *grugBot) Score(req bots.BotRequest) (bots.ScoreResponse, error) {
	hand, err := game.DecodeCards(req.Hand)

	if err != nil {
		return bots.ScoreResponse{}, fmt.Errorf("unable to decode cards: %w", err)
	}

	// sort the cards by suite and number
	slices.SortFunc(hand, game.CompareCard)

	// find all possible sequences in the hand
	seqs := FindSequences(req.Round, hand)

	// filter out sequences if they have cards that have been used twiced
	// use the wilds to build more sequences
	seqs = FilterSequences(req.Round, hand, seqs)

	return bots.ScoreResponse{
		Action:    req.Action,
		Flop:      false,
		Sequences: game.EncodeSequences(seqs),
	}, nil
}

func (b *grugBot) Draw(req bots.BotRequest) (bots.DrawResponse, error) {
	return bots.DrawResponse{}, errors.New("not implemented")
}

func (b *grugBot) Discard(req bots.BotRequest) (bots.DiscardResponse, error) {
	return bots.DiscardResponse{}, errors.New("not implemented")
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

		if c.Suite == pc.Suite && c.Number == pc.Number+1 {
			// cards are in sequence, add this one to the sequence
			curSeq = append(curSeq, c)
		} else {
			// cards are not in a run

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
	for number := 3; number < 13; number++ {

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

// takes a hand and a list of sequences
// determines if any card is being used more than it should be
// if so, will chose the sequence with the highest score
// output will include any cards which dont fit within a sequence as a single-carded sequence
func FilterSequences(round int, hand []game.Card, seqs [][]game.Card) [][]game.Card {

	// score all of the sequences

	slices.SortFunc(seqs, func(a, b []game.Card) int {

		if len(a) < 3 && len(b) >= 3 {
			return 1
		} else if len(a) >= 3 && len(b) < 3 {
			return -1
		} else {
			return game.ScoreSequence(b) - game.ScoreSequence(a)
		}
	})

	slog.Info("sorting by sequence scores", "seqs", game.EncodeSequences(seqs))

	filteredSeqs := make([][]game.Card, 0)

	cardCounts := CardCounts(hand)

	for i, seq := range seqs {

		valid := true

		// TODO: if a sequence cannot be used, its subset should be added to the seqs

		// don't bother checking for card usage for first sequence
		if i > 0 {

			// check all the cards are available
			for card, reqCount := range CardCounts(seq) {
				if cardCounts[card] < reqCount {
					valid = false
					break
				}

			}
		}

		if !valid {
			continue
		}

		// add jokers to the sequence if the sequence is < 3
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
