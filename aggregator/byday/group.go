package byday

import "github.com/fabiofenoglio/excelconv/model"

type GroupedByDay struct {
	Key  string
	Rows []model.ParsedRow
}
