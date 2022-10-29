package reader

import (
	"strings"
	"unicode"
)

func stringToCode(raw string) string {
	cleaned := strings.ToLower(removeWhiteSpaces(raw))
	if len(cleaned) > 0 {
		return cleaned
	}
	return raw
}

func removeWhiteSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
