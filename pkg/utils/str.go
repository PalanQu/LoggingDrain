package utils

import (
	"strings"
	"unicode"
)

func GetStringTokens(message string) []string {
	content := strings.TrimSpace(message)
	return strings.Fields(content)
}

func StringHasNumber(message string) bool {
	for _, char := range message {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
