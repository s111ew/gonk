package editor

func truncateString(s string, x int) string {
	runes := []rune(s)
	if len(runes) > x {
		return string(runes[:x])
	}
	return s
}
