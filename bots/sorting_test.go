package bots

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareCard(t *testing.T) {

	cases := []struct {
		Args     []Card
		Expected int
	}{
		{
			Args: []Card{
				{
					Joker: true,
				},
				{
					Number: 10,
					Suite:  'R',
				},
			},
			Expected: 1,
		},
		{
			Args: []Card{
				{
					Number: 10,
					Suite:  'R',
				},
				{
					Joker: true,
				},
			},
			Expected: -1,
		},
		{
			Args: []Card{
				{
					Number: 10,
					Suite:  'R',
				},
				{
					Number: 10,
					Suite:  'B',
				},
			},
			Expected: 16,
		},
		{
			Args: []Card{
				{
					Number: 11,
					Suite:  'B',
				},
				{
					Number: 10,
					Suite:  'B',
				},
			},
			Expected: 1,
		},
	}

	for _, tc := range cases {

		t.Run("should compare two cards", func(t *testing.T) {

			res := CompareCard(tc.Args[0], tc.Args[1])

			assert.Equal(t, tc.Expected, res)
		})

	}
}

func TestSortSequence(t *testing.T) {

	seq := []Card{
		{
			Number: 9,
			Suite:  'Y',
		},
		{
			Joker: true,
		},
		{
			Number: 8,
			Suite:  'R',
		},
		{
			Number: 9,
			Suite:  'R',
		},
	}

	slices.SortFunc(seq, CompareCard)

	assert.Equal(t, "8-R:9-R:9-Y:*", EncodeSequence(seq))

}
