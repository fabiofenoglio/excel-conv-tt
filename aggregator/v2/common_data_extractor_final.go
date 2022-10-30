package aggregator

import (
	"github.com/fabiofenoglio/excelconv/config"
)

func ExtractCommonDataFinal(
	_ config.WorkflowContext,
	currentCommonData CommonData,
	daysWithRoomsAndGroupingSlots []ScheduleForSingleDayWithRoomsAndGroupSlots,
) CommonData {
	output := currentCommonData

	maxVisitingGroups := 0

	for _, day := range daysWithRoomsAndGroupingSlots {
		if len(day.VisitingGroups) > maxVisitingGroups {
			maxVisitingGroups = len(day.VisitingGroups)
		}
	}

	output.MaxVisitingGroupsPerDay = maxVisitingGroups
	return output
}
