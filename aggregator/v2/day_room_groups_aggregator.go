package aggregator

import (
	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

const (
	layoutDateHourMinsOrderable = "2006-01-02T15:04"
)

func AggregateByRooomGroupsOnSameActivity(
	_ config.WorkflowContext,
	days []ScheduleForSingleDayWithRooms,
	anagraphicsRef *parser.OutputAnagraphics,
) []ScheduleForSingleDayWithRoomsAndGroupedActivities {

	grouped := make([]*ScheduleForSingleDayWithRoomsAndGroupedActivities, 0)

	for _, daySchedule := range days {
		mapped := &ScheduleForSingleDayWithRoomsAndGroupedActivities{
			Day:            daySchedule.Day,
			VisitingGroups: daySchedule.VisitingGroups,
			RoomsSchedule:  make([]ScheduleForSingleDayAndRoomWithGroupedActivities, 0, len(daySchedule.RoomsSchedule)),
			StartAt:        daySchedule.StartAt,
			EndAt:          daySchedule.EndAt,
		}

		for _, roomSchedule := range daySchedule.RoomsSchedule {
			mapped.RoomsSchedule = append(mapped.RoomsSchedule,
				aggregateScheduleForSingleDayAndRoomWithGroupedActivities(roomSchedule, anagraphicsRef))
		}

		grouped = append(grouped, mapped)
	}

	cnt := 1

	for i1, group := range grouped {
		for i2, g := range group.RoomsSchedule {
			for i3 := range g.GroupedActivities {
				grouped[i1].RoomsSchedule[i2].GroupedActivities[i3].ID = cnt
				cnt++
			}
		}
	}

	out := make([]ScheduleForSingleDayWithRoomsAndGroupedActivities, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func aggregateScheduleForSingleDayAndRoomWithGroupedActivities(
	input ScheduleForSingleDayAndRoom,
	anagraphicsRef *parser.OutputAnagraphics,
) ScheduleForSingleDayAndRoomWithGroupedActivities {
	mapped := ScheduleForSingleDayAndRoomWithGroupedActivities{
		RoomCode:          input.RoomCode,
		GroupedActivities: make([]GroupedActivity, 0, len(input.Rows)),
	}

	indexMap := make(map[string]int)

	for _, row := range input.Rows {
		room := anagraphicsRef.Rooms[input.RoomCode]

		// we group if they have the same start time, end time (and activity code ???)
		if !room.GroupActivities || row.StartTime.IsZero() || row.EndTime.IsZero() || row.ActivityCode == "" {
			// cannot be grouped
			mapped.GroupedActivities = append(mapped.GroupedActivities, GroupedActivity{
				StartTime: row.StartTime,
				EndTime:   row.EndTime,
				Rows:      []OutputRow{row},
			})
			continue
		}

		hash := row.StartTime.Format(layoutDateHourMinsOrderable) + "/" +
			row.EndTime.Format(layoutDateHourMinsOrderable)

		mappedAt, already := indexMap[hash]

		if !already {
			// nothing mapped under this key yet
			mapped.GroupedActivities = append(mapped.GroupedActivities, GroupedActivity{
				StartTime: row.StartTime,
				EndTime:   row.EndTime,
				Rows:      []OutputRow{row},
			})
			indexMap[hash] = len(mapped.GroupedActivities) - 1
		} else {
			mapped.GroupedActivities[mappedAt].Rows = append(mapped.GroupedActivities[mappedAt].Rows, row)
		}
	}

	return mapped
}
