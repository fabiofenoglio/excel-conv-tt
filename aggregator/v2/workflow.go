package aggregator

import (
	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

func Execute(ctx config.WorkflowContext, rawInput parser.Output) (Output, error) {
	input := ToInput(rawInput)

	rowsWithCompetenceDate := AssignCompetenceDay(ctx, input.Rows)

	commonData := ExtractCommonData(ctx, rowsWithCompetenceDate)

	days := AggregateByCompetenceDay(ctx, rowsWithCompetenceDate, rawInput.Anagraphics)

	daysWithRooms := AggregateByRooomInCompetenceDay(ctx, days, rawInput.Anagraphics)

	// TODO remove and clear aggregator
	// daysWithRoomsAndSlots := AggregateByRooomSlotInRoom(ctx, daysWithRooms, rawInput.Anagraphics)

	daysWithRoomsAndGrouping := AggregateByRooomGroupsOnSameActivity(ctx, daysWithRooms, rawInput.Anagraphics)

	daysWithRoomsAndGroupingSlots := AggregateByRooomGroupSlotInRoom(ctx, daysWithRoomsAndGrouping, rawInput.Anagraphics)

	return Output{
		CommonData: commonData,
		Days:       daysWithRoomsAndGroupingSlots,
	}, nil
}
