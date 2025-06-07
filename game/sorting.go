package game

// compares two cards
func CompareCard(c1, c2 Card) int {
	if c1 == c2 {
		return 0
	} else if c1.Joker {
		return 1
	} else if c2.Joker {
		return -1
	} else if c1.Suite == c2.Suite {
		return c1.Number - c2.Number
	} else {
		return int(c1.Suite) - int(c2.Suite)
	}
}
