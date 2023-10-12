package loggingdrain

import (
	"strings"
	"unicode"
)

func getStringTokens(message string) []string {
	content := strings.TrimSpace(message)
	return strings.Fields(content)
}

func stringHasNumber(message string) bool {
	for _, char := range message {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
