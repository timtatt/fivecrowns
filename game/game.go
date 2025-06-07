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
