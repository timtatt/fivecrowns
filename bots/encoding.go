package bots

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Card struct {
	Joker  bool
	Number int
	Suite  rune
}

// encode the card into a string eg. "*" or "10-R"
func (c Card) Encode() string {
	if c.Joker {
		return "*"
	}

	return fmt.Sprintf("%d-%c", c.Number, c.Suite)
}

// decode the card into a struct
func DecodeCard(c string) (Card, error) {
	if c == "*" {
		return Card{
			Joker: true,
		}, nil
	}

	parts := strings.SplitN(string(c), "-", 2)
	number, err := strconv.Atoi(parts[0])

	if err != nil {
		return Card{}, fmt.Errorf("unable to decode number: %s", c)
	} else if number < 3 || number > 13 {
		return Card{}, fmt.Errorf("invalid number in card encoding: %s", c)
	}

	suite, _ := utf8.DecodeLastRuneInString(parts[1])

	switch suite {
	case SuiteBlue, SuiteRed, SuiteBlack, SuiteGreen, SuiteYellow:
		break
	default:
		return Card{}, fmt.Errorf("invalid suite in card encoding: %s", c)
	}

	return Card{
		Joker:  false,
		Number: number,
		Suite:  suite,
	}, nil
}

type Sequence []Card

// encode the sequence into a string eg. 10-R:8-Y:*:10-X
func (s Sequence) Encode() string {

	var out strings.Builder

	for i, card := range s {
		if i != 0 {
			out.WriteString(":")
		}

		out.WriteString(string(card.Encode()))
	}

	return out.String()
}

// convert the sequence into a list of CardCodes
// eg. ["10-R", "8-Y", "*"]
func (s Sequence) EncodeCards() []string {

	sequenceCode := make([]string, len(s))

	for i, card := range s {
		sequenceCode[i] = card.Encode()
	}

	return sequenceCode
}

// takes encoded sequence eg. 10-R:*:8-Y and decodes into a list of Cards
func DecodeSequence(s string) (Sequence, error) {

	cardCodes := strings.Split(s, ":")

	seq := make([]Card, len(cardCodes))

	var errs error
	for i, code := range cardCodes {
		card, err := DecodeCard(code)

		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		seq[i] = card
	}

	return seq, errs
}
