package parser

import (
	"strings"
	"unicode"
)

func nameToCode(raw string) string {
	return strings.ToLower(removeWhiteSpaces(raw))
}

func removeWhiteSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) && !unicode.IsControl(ch) && unicode.IsPrint(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func cleanStringForVisualization(raw string) string {
	return strings.TrimSpace(raw)
}
