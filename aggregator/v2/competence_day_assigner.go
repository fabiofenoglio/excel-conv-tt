package aggregator

import (
	"time"

	"github.com/fabiofenoglio/excelconv/config"
)

func AssignCompetenceDay(_ config.WorkflowContext, rows []InputRow) []Row {
	out := make([]Row, 0, len(rows))

	for _, row := range rows {
		out = append(out, Row{
			InputRow:       row,
			CompetenceDate: extractCompetenceDay(row),
		})
	}

	return out
}

func extractCompetenceDay(row InputRow) time.Time {
	start := row.StartTime

	if start.IsZero() {
		return row.Date
	}

	if start.Hour() >= 6 {
		return time.Date(start.Year(), start.Month(), start.Day(), 12, 0, 0, 0, start.Location())
	}

	return time.Date(start.Year(), start.Month(), start.Day()-1, 12, 0, 0, 0, start.Location())
}
