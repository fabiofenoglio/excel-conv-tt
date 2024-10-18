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
		if row.Activity != "" && containsNAttivitaRegex.MatchString(row.Activity) {
			row.IsPlaceholderNumeroAttivita = true
		}
		out = append(out, row)
	}

	return out, nil
}
