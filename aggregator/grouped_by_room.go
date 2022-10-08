package aggregator

import "github.com/fabiofenoglio/excelconv/model"

type GroupedByRoom struct {
	Room model.Room
	Rows []model.ParsedRow
}
