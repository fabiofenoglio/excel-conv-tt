package parser

import (
	"regexp"
	_ "time/tzdata"
)

var (
	containsNAttivitaRegex = regexp.MustCompile(`(?i)\s*[n0-9]+\s*[°um]*\s*attivit[aà]\'?\s*`)
)

func filter(r ExcelRow) bool {
	if r.Attivita != "" && containsNAttivitaRegex.MatchString(r.Attivita) {
		return false
	}

	return true
}
