package game

type Suite rune

var (
	SuiteBlue   rune = 'B'
	SuiteRed    rune = 'R'
	SuiteYellow rune = 'Y'
	SuiteGreen  rune = 'G'
	SuiteBlack  rune = 'X'
	Suites           = []rune{SuiteBlue, SuiteGreen, SuiteBlack, SuiteRed, SuiteYellow}
)

type Card struct {
	Joker  bool
	Number int
	Suite  rune
}

var (
	CardJoker = Card{Joker: true}
)

func (c Card) IsWild(round int) bool {
	return c.Joker || c.Number == round
}

type SequenceType string

var (
	SequenceTypeRun    SequenceType = "run"
	SequenceTypeSet    SequenceType = "set"
	SequenceTypeEither SequenceType = "either"
)

// determine if the sequence is a set
// ignores wild cards
func GetSequenceType(seq []Card, round int) SequenceType {
	c1 := -1

	for i := 0; i < len(seq); i++ {

		if seq[i].IsWild(round) {
			continue
		}

		if c1 == -1 {
			c1 = i
		} else if seq[c1].Number == seq[i].Number {
			return SequenceTypeSet
		} else {
			return SequenceTypeRun
		}
	}

	return SequenceTypeEither
}

func CanFlop(seqs [][]Card) bool {

	canFlop := true

	for _, seq := range seqs {
		if len(seq) < 3 {
			canFlop = false
		}
	}

	return canFlop
}
