package grugbot

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timtatt/fivecrowns/bots"
	"github.com/timtatt/fivecrowns/game"
)

func TestGrugbot(t *testing.T) {

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
				"3-B",
				"3-X",
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
