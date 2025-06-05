package bots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardEncoding(t *testing.T) {

	cases := []struct {
		Arg      Card
		Expected string
	}{
		{
			Arg: Card{
				Joker: true,
			},
			Expected: "*",
		},
		{
			Arg: Card{
				Joker:  false,
				Number: 10,
				Suite:  'R',
			},
			Expected: "10-R",
		},
	}

	for _, tc := range cases {
		t.Run("should encode joker: "+tc.Expected, func(t *testing.T) {

			res := tc.Arg.Encode()
			assert.Equal(t, tc.Expected, res)
		})
	}
}

func TestSequenceEncode(t *testing.T) {

	seq := Sequence{
		Card{
			Joker: true,
		},
		Card{
			Number: 10,
			Suite:  'R',
		},
	}

	assert.Equal(t, "*:10-R", seq.Encode())

}

func TestSequenceEncodeCards(t *testing.T) {
	seq := Sequence{
		Card{
			Joker: true,
		},
		Card{
			Number: 10,
			Suite:  'R',
		},
	}

	assert.Equal(t, []string{"*", "10-R"}, seq.EncodeCards())
}

func TestSequenceDecode(t *testing.T) {

	encSeq := "*:10-R:7-G"

	seq, err := DecodeSequence(encSeq)

	require.NoError(t, err)

	assert.Equal(t, seq, Sequence{
		{
			Joker: true,
		},
		{
			Number: 10,
			Suite:  'R',
		},
		{
			Number: 7,
			Suite:  'G',
		},
	})
}
