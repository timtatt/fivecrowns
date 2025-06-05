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
)

type BotRequest struct {
	Action      Action   `json:"action"`
	PlayerCount int      `json:"playerCount"`
	Round       int      `json:"round"`
	LastTurn    bool     `json:"lastTurn"`
	NewestCard  string   `json:"newestCard"`
	Hand        []string `json:"hand"`
	Discard     []string `json:"discard"`
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
	Action    Action     `json:"action"`
	Card      string     `json:"card"`
	Flop      bool       `json:"flop"`
	Sequences [][]string `json:"sequences"`
}

type ScoreResponse struct {
	Action    Action     `json:"action"`
	Flop      bool       `json:"flop"`
	Sequences [][]string `json:"sequences"`
}
