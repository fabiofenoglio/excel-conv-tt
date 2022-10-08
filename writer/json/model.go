package json

import "github.com/fabiofenoglio/excelconv/model"

type SerializableOutput struct {
	Activities []model.ParsedRow `json:"activities"`
	Rooms      []model.Room      `json:"rooms"`
	Operators  []model.Operator  `json:"operators"`
}
