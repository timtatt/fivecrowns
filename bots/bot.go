package bots

type Bot interface {
	Draw(req BotRequest) (DrawResponse, error)
	Discard(req BotRequest) (DiscardResponse, error)
	Score(req BotRequest) (ScoreResponse, error)
}

type Suite rune

var (
	SuiteBlue   rune = 'B'
	SuiteRed    rune = 'R'
	SuiteYellow rune = 'Y'
	SuiteGreen  rune = 'G'
	SuiteBlack  rune = 'X'
	Suites           = []rune{SuiteBlue, SuiteGreen, SuiteBlack, SuiteRed, SuiteYellow}
)

type BotRequest struct {
	Discard     []string `json:"discard"`
	Hand        []string `json:"hand"`
	Action      Action   `json:"action"`
	NewestCard  string   `json:"newestCard"`
	PlayerCount int      `json:"playerCount"`
	Round       int      `json:"round"`
	LastTurn    bool     `json:"lastTurn"`
}

type Action string

const (
	ActionDraw    Action = "draw"
	ActionDiscard Action = "discard"
	ActionScore   Action = "score"
)

type Stack string

const (
	StackDiscard Stack = "discard"
	StackDeck    Stack = "deck"
)

type DrawResponse struct {
	Action Action `json:"action"`
	Stack  Stack  `json:"stack"`
}

type DiscardResponse struct {
	Sequences [][]string `json:"sequences"`
	Action    Action     `json:"action"`
	Card      string     `json:"card"`
	Flop      bool       `json:"flop"`
}

type ScoreResponse struct {
	Sequences [][]string `json:"sequences"`
	Action    Action     `json:"action"`
	Flop      bool       `json:"flop"`
}
