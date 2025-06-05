package grugbot

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"slices"

	"github.com/timtatt/fivecrowns/bots"
)

// Grugbot tries to put cards in the biggest possible sequences
// if a card is used in multiple sequences, the one with the biggest score is preferred

type grugBot struct{}

func NewGrugBot() bots.Bot {
	return &grugBot{}
}

func (b *grugBot) Score(req bots.BotRequest) (bots.ScoreResponse, error) {
	hand, err := bots.DecodeCards(req.Hand)

	if err != nil {
		return bots.ScoreResponse{}, fmt.Errorf("unable to decode cards: %w", err)
	}

	FindSequences(hand)

	return bots.ScoreResponse{}, errors.New("not implemented")
}

func (b *grugBot) Draw(req bots.BotRequest) (bots.DrawResponse, error) {
	return bots.DrawResponse{}, errors.New("not implemented")
}

func (b *grugBot) Discard(req bots.BotRequest) (bots.DiscardResponse, error) {
	return bots.DiscardResponse{}, errors.New("not implemented")
}

// takes a list of cards
func FindSequences(hand []bots.Card) {

	// sort the cards by suite and number
	slices.SortFunc(hand, bots.CompareCard)

	// put all the cards into a map
	cardCounts := make(map[bots.Card]int)

	for _, card := range hand {
		cardCounts[card] += 1
	}

	slog.Info("card counts in hand", "hand", cardCounts)

	seqs := make([][]bots.Card, 0)

	// go through the sorted cards and find runs of numbers in the same suite

	curSeq := []bots.Card{hand[0]}

	for i := 1; i < len(hand); i++ {

		// check if the current number is in sequence with previous

		c := hand[i]
		pc := curSeq[len(curSeq)-1]

		// if we are dealing with jokers, we can ignore them
		// if the current card and previous card in sequence is the same, we can ignore it
		if c.Joker || c == pc {
			continue
		}

		if c.Suite == pc.Suite && c.Number == pc.Number+1 {
			// cards are in sequence, add this one to the sequence
			curSeq = append(curSeq, c)
		} else {
			// cards are not in a run

			// save the sequence if it qualifies
			if len(curSeq) >= 3 {
				log.Println(bots.EncodeSequence(curSeq))
				seqs = append(seqs, curSeq)
			}

			// reset the current sequence
			curSeq = []bots.Card{c}
		}
	}

	// ensure we don't leave a curSeq hanging
	if len(curSeq) >= 3 {
		seqs = append(seqs, curSeq)
	}

	// check the cards for sets
	for number := 3; number < 13; number++ {

		curSeq := make([]bots.Card, 0)

		for _, suite := range bots.Suites {

			c := bots.Card{
				Joker:  false,
				Number: number,
				Suite:  suite,
			}

			for range cardCounts[c] {
				curSeq = append(curSeq, c)
			}

		}

		if len(curSeq) >= 3 {
			seqs = append(seqs, curSeq)
		}
	}

}
