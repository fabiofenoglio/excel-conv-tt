package aggregator

import "github.com/fabiofenoglio/excelconv/model"

type GroupedByRoom struct {
	Room            model.Room
	Slots           []GroupedByRoomSlot
	TotalInAllSlots uint
}

type GroupedByRoomSlot struct {
	SlotIndex uint
	Rows      []model.ParsedRow
}
