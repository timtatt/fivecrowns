package bots

func ScoreSequence(cards []Card) int {
	score := 0
	for _, card := range cards {
		if card.Joker {
			score += 25
		} else {
			score += card.Number
		}
	}

	return score
}
