package game

func ScoreSequence(cards []Card) int {
	score := 0
	for _, card := range cards {
		score += ScoreCard(card)
	}

	return score
}

func ScoreCard(card Card) int {
	if card.Joker {
		return 25
	} else {
		return card.Number
	}
}
