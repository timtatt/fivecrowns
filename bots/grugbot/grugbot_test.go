package grugbot

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timtatt/fivecrowns/bots"
	"github.com/timtatt/fivecrowns/game"
)

func TestGrugbotScore(t *testing.T) {

	cases := []struct {
		Hand     string
		Round    int
		Expected []string
	}{
		{
			Hand:  "5-B:*:5-R:4-B:6-B",
			Round: 10,
			Expected: []string{
				"4-B:5-B:6-B:*",
				"5-R",
			},
		},
		{
			Hand:  "5-B:5-R:4-B:6-B:6-R:7-X:*",
			Round: 7,
			Expected: []string{
				"4-B:5-B:6-B:7-X",
				"5-R:6-R:*",
			},
		},
		{
			Hand:  "5-B:5-R:4-B:6-B:7-X:*:3-Y",
			Round: 7,
			Expected: []string{
				"4-B:5-B:6-B",
				"5-R:*:7-X",
				"3-Y",
			},
		},
		{
			Hand:  "3-X:3-Y:3-B",
			Round: 3,
			Expected: []string{
				"3-B:3-X:3-Y",
			},
		},
		{
			Hand:  "3-X:4-Y:4-B:4-R",
			Round: 4,
			Expected: []string{
				"3-X:4-B:4-R:4-Y",
			},
		},
		{
			Hand:  "3-X:4-X:5-X:3-X:3-B",
			Round: 6,
			Expected: []string{
				"3-X:4-X:5-X",
				"3-B:3-X",
			},
		},
		{
			Hand:  "5-B:5-R:4-B:6-B:7-X:5-Y:6-G",
			Round: 7,
			Expected: []string{
				"5-B:5-R:5-Y",
				"6-B:6-G:7-X",
				"4-B",
			},
		},
		{
			Hand:  "6-Y:13-B:4-X:13-Y:3-X:7-B:8-X",
			Round: 7,
			Expected: []string{
				"13-B:13-Y:7-B",
				"8-X",
				"3-X:4-X",
				"6-Y",
			},
		},
	}

	for _, tc := range cases {

		gb := NewGrugBot()

		t.Run("test best sequence from hand: "+tc.Hand, func(t *testing.T) {
			hand := strings.Split(tc.Hand, ":")

			res, err := gb.Score(bots.BotRequest{
				Action: bots.ActionScore,
				Hand:   hand,
				Round:  tc.Round,
			})

			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, game.FlattenSequences(res.Sequences))

		})

	}

}

func TestGrugbotDraw(t *testing.T) {

	cases := []struct {
		Hand     string
		Round    int
		Discard  string
		Expected bots.Stack
	}{
		{
			Hand:     "9-R:10-R:5-X:8-R:6-B:8-B:11-R:11-Y:4-Y",
			Round:    9,
			Discard:  "11-R",
			Expected: bots.StackDiscard,
		},
		{
			Hand:     "9-R:10-R:5-X:8-R:6-B:8-B:11-R:11-Y:4-Y",
			Round:    9,
			Discard:  "4-G",
			Expected: bots.StackDeck,
		},
	}

	for _, tc := range cases {

		gb := NewGrugBot()

		t.Run("test best stack to draw from hand: "+tc.Hand, func(t *testing.T) {
			hand := strings.Split(tc.Hand, ":")

			res, err := gb.Draw(bots.BotRequest{
				Action:  bots.ActionScore,
				Hand:    hand,
				Round:   tc.Round,
				Discard: []string{tc.Discard},
			})

			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, res.Stack)

		})

	}

}

func TestGrugbotDiscard(t *testing.T) {

	cases := []struct {
		Hand     string
		Round    int
		Expected string
	}{
		{
			Hand:     "*:12-R:3-B:11-B:5-B:3-R:11-R:4-Y",
			Round:    7,
			Expected: "11-B",
		},
		{
			Hand:     "5-X:3-B:5-R:9-Y:13-B:11-X:6-Y",
			Round:    7,
			Expected: "13-B",
		},
		{
			Hand:     "4-G:7-R:7-X:8-R:7-B:4-R:*",
			Round:    7,
			Expected: "8-R",
		},
	}

	for _, tc := range cases {

		gb := NewGrugBot()

		t.Run("test best stack to draw from hand: "+tc.Hand, func(t *testing.T) {
			hand := strings.Split(tc.Hand, ":")

			res, err := gb.Discard(bots.BotRequest{
				Action: bots.ActionScore,
				Hand:   hand,
				Round:  tc.Round,
			})

			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, res.Card)

		})

	}

}
