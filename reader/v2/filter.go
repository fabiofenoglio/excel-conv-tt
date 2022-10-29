package reader

import (
	"regexp"
)

var (
	containsNAttivitaRegex = regexp.MustCompile(`(?i)\s*[n0-9]+\s*[°um]*\s*attivit[aà]\'?\s*`)
)

func Filter(rows []Row) ([]Row, error) {
	out := make([]Row, 0, len(rows))

	for _, row := range rows {
		if filterRow(row) {
			out = append(out, row)
		}
	}

	return out, nil
}

func filterRow(r Row) bool {
	if r.Activity != "" && containsNAttivitaRegex.MatchString(r.Activity) {
		return false
	}

	return true
}
