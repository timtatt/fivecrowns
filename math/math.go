package math

func Min[T int | float64](a, b T) T {
	if a > b {
		return b
	}
	return a
}
