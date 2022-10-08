package json

import "github.com/fabiofenoglio/excelconv/model"

type TemplateData struct {
	PageTitle           string
	GenerationTimeStamp string

	Activities []model.ParsedRow
	Rooms      []model.Room
	Operators  []model.Operator
}
