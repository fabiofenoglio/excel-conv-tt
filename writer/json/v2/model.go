package json

import (
	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
)

type SerializableOutput struct {
	CommonData     aggregator2.CommonData                                   `json:"CommonData"`
	Days           []aggregator2.ScheduleForSingleDayWithRoomsAndGroupSlots `json:"Days"`
	AnagraphicsRef *parser2.OutputAnagraphics                               `json:"Anagraphics"`
}
