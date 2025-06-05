package stupidbot

import (
	"math/rand"

	"github.com/timtatt/fivecrowns/bots"
)

type stupidBot struct{}

func NewStupidBot() bots.Bot {
	return &stupidBot{}
}

// StupidBot randomly picks a stack to draw from
func (*stupidBot) Draw(req bots.BotRequest) (bots.DrawResponse, error) {

	r := rand.Intn(2)
	stack := bots.StackDeck
	if r == 1 {
		stack = bots.StackDiscard
	}

	return bots.DrawResponse{
		Action: bots.ActionDraw,
		Stack:  stack,
	}, nil

}

// StupidBot is unable to make a sequence and will discard a card at random
func (s *stupidBot) Discard(req bots.BotRequest) (bots.DiscardResponse, error) {
	discardIdx := rand.Intn(len(req.Hand))

	sequences := make([][]string, 0, len(req.Hand))

	for i, card := range req.Hand {
		if i != discardIdx {
			sequences = append(sequences, []string{card})
		}
	}

	return bots.DiscardResponse{
		Action:    bots.ActionDiscard,
		Card:      req.Hand[discardIdx],
		Sequences: sequences,
	}, nil
}

func (s *stupidBot) Score(req bots.BotRequest) (bots.ScoreResponse, error) {

	sequences := make([][]string, 0, len(req.Hand))

	for _, card := range req.Hand {
		sequences = append(sequences, []string{card})
	}

	return bots.ScoreResponse{
		Action:    bots.ActionDiscard,
		Flop:      false,
		Sequences: sequences,
	}, nil

}
