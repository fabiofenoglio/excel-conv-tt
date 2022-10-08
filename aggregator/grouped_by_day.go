package aggregator

import "github.com/fabiofenoglio/excelconv/model"

type GroupedByDay struct {
	Key  string
	Rows []model.ParsedRow
}
