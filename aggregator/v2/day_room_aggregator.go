package aggregator

import (
	"sort"
	"strings"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

func AggregateByRooomInCompetenceDay(_ config.WorkflowContext, days []scheduleForSingleDay, anagraphicsRef *parser.OutputAnagraphics) []ScheduleForSingleDayWithRooms {

	grouped := make([]*ScheduleForSingleDayWithRooms, 0)

	for _, daySchedule := range days {
		mapped := &ScheduleForSingleDayWithRooms{
			Day:            daySchedule.Day,
			VisitingGroups: daySchedule.VisitingGroups,
			RoomsSchedule:  make([]ScheduleForSingleDayAndRoom, 0, len(anagraphicsRef.Rooms)),
			StartAt:        daySchedule.StartAt,
			EndAt:          daySchedule.EndAt,
		}

		for _, room := range anagraphicsRef.Rooms {
			roomSchedule := ScheduleForSingleDayAndRoom{
				RoomCode: room.Code,
				Rows:     nil,
			}

			for _, row := range daySchedule.Rows {
				if row.InputRow.RoomCode == room.Code {
					roomSchedule.Rows = append(roomSchedule.Rows, ToOutputRow(row))
				}
			}

			mapped.RoomsSchedule = append(mapped.RoomsSchedule, roomSchedule)
		}

		grouped = append(grouped, mapped)
	}

	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Day.Before(grouped[j].Day)
	})

	for _, group := range grouped {
		visitingGroupSequentialCodeMap := make(map[string]string)
		for _, visitingGroupThisDay := range group.VisitingGroups {
			visitingGroupSequentialCodeMap[visitingGroupThisDay.VisitingGroupCode] = visitingGroupThisDay.SequentialCode
		}

		// sort rooms schedule by preferred order, room name, room code
		sort.Slice(group.RoomsSchedule, func(i, j int) bool {
			roomRefI := anagraphicsRef.Rooms[group.RoomsSchedule[i].RoomCode]
			roomRefJ := anagraphicsRef.Rooms[group.RoomsSchedule[j].RoomCode]
			diff := roomRefI.PreferredOrder - roomRefJ.PreferredOrder
			if diff != 0 {
				return diff < 0
			}

			diff = strings.Compare(roomRefI.Name, roomRefJ.Name)
			if diff != 0 {
				return diff < 0
			}

			return strings.Compare(roomRefI.Code, roomRefJ.Code) < 0
		})

		// order each rooms' activities by StartTime, VisitingGroup.SequentialCode, EndTime, BookingCode, ID
		for i, roomSchedule := range group.RoomsSchedule {
			rows := roomSchedule.Rows

			sort.Slice(rows, func(i, j int) bool {
				if rows[i].StartTime.UnixMilli() < rows[j].StartTime.UnixMilli() {
					return true
				}
				if rows[i].StartTime.UnixMilli() > rows[j].StartTime.UnixMilli() {
					return false
				}
				groupCodeI := visitingGroupSequentialCodeMap[rows[i].VisitingGroupCode]
				groupCodeJ := visitingGroupSequentialCodeMap[rows[j].VisitingGroupCode]
				diff := strings.Compare(groupCodeI, groupCodeJ)
				if diff != 0 {
					return diff < 0
				}
				if rows[i].EndTime.UnixMilli() < rows[j].EndTime.UnixMilli() {
					return true
				}
				if rows[i].EndTime.UnixMilli() > rows[j].EndTime.UnixMilli() {
					return false
				}
				diff = strings.Compare(rows[i].BookingCode, rows[j].BookingCode)
				if diff != 0 {
					return diff < 0
				}
				return rows[i].ID < rows[j].ID
			})

			group.RoomsSchedule[i].Rows = rows
		}
	}

	out := make([]ScheduleForSingleDayWithRooms, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}
