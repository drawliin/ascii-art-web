package helpers

// Splits banner file's content into a slice of ascii representations
// Returns a slice of ascii represented symbols (each of them as a slice of 8 lines/strings)
func Split2D(s string) [][]string {
	arr := [][]string{}
	symbol := []string{}
	wordStart := 0
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			if s[wordStart:i] == "" {
				if count == 8 {
					arr = append(arr, symbol)
					count = 0
					symbol = []string{}
				}
			} else {
				symbol = append(symbol, s[wordStart:i])
				count++
			}
			wordStart = i + 1
		}
	}
	if len(symbol) > 0 {
		arr = append(arr, symbol)
	}
	return arr
}

func ContainOnlyNewLines(s string) bool {
	for i := 0; i < len(s); i += 2 {
		if i+2 > len(s) || s[i:i+2] != "\\n" {
			return false
		}
	}
	return true
}