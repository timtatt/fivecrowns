package bigbrainbot

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/timtatt/fivecrowns/bots"
	"github.com/timtatt/fivecrowns/bots/grugbot"
	"github.com/timtatt/fivecrowns/game"
)

type bigBrainBot struct{}

func NewBigBrainBot() bots.Bot {
	return &bigBrainBot{}
}

func (b *bigBrainBot) Score(req bots.BotRequest) (bots.ScoreResponse, error) {

	calculation, err := grugbot.Calculate(req)

	if err != nil {
		return bots.ScoreResponse{}, fmt.Errorf("cannot calculate response: %w", err)
	}

	return bots.ScoreResponse{
		Action:    req.Action,
		Flop:      game.CanFlop(calculation.Sequences),
		Sequences: game.EncodeSequences(calculation.Sequences),
	}, nil
}

func (b *bigBrainBot) Draw(req bots.BotRequest) (bots.DrawResponse, error) {

	// add the discard to the hand and determine if it gets added to a valid sequence

	topCard, err := game.DecodeCard(req.Discard[0])

	if err != nil {
		return bots.DrawResponse{}, fmt.Errorf("unable to decode discard card: %w", err)
	}

	hypothetical, err := Calculate(bots.BotRequest{
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

func (b *bigBrainBot) Discard(req bots.BotRequest) (bots.DiscardResponse, error) {

	calculation, err := Calculate(req)

	if err != nil {
		return bots.DiscardResponse{}, fmt.Errorf("unable to calculate score: %w", err)
	}

	// determine which card is the highest one that is not in a valid sequence
	worstCard := grugbot.WorstCard(req.Round, calculation.Sequences, req.LastTurn)

	slog.Info("worst card detected", "card", worstCard)

	// update the discard response to omit the discarded card
	seq := calculation.Sequences[worstCard.SequenceIdx]

	if len(seq) == 1 {
		// the sequence only has one card, delete the whole thing
		calculation.Sequences = slices.Delete(calculation.Sequences, worstCard.SequenceIdx, worstCard.SequenceIdx+1)
	} else {
		// remove the card from the its sequence
		calculation.Sequences[worstCard.SequenceIdx] = slices.Delete(calculation.Sequences[worstCard.SequenceIdx], worstCard.CardIdx, worstCard.CardIdx+1)
	}

	return bots.DiscardResponse{
		Flop:      game.CanFlop(calculation.Sequences),
		Sequences: game.EncodeSequences(calculation.Sequences),
		Action:    bots.ActionDiscard,
		Card:      worstCard.Card.Encode(),
	}, nil
}

type Calculation struct {
	Sequences [][]game.Card
	Flop      bool
	Hand      []game.Card
}

func Calculate(req bots.BotRequest) (Calculation, error) {
	hand, err := game.DecodeCards(req.Hand)

	if err != nil {
		return Calculation{}, fmt.Errorf("unable to decode cards: %w", err)
	}

	// sort the cards by suite and number
	slices.SortFunc(hand, game.CompareCard)

	// find all possible sequences in the hand
	seqs := grugbot.FindSequences(req.Round, hand)

	slog.Info("calculated possible sequences", "seqs", game.EncodeSequences(seqs))

	// filter out sequences if they have cards that have been used twiced
	// use the wilds to build more sequences
	seqs = grugbot.FilterSequences(req.Round, hand, seqs)

	return Calculation{
		Flop:      game.CanFlop(seqs),
		Sequences: seqs,
	}, nil
}
