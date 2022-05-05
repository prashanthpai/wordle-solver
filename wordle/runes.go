package wordle

func containsRune(haystack []rune, needle rune) bool {
	for _, r := range haystack {
		if r == needle {
			return true
		}
	}

	return false
}
