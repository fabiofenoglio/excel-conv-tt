package aggregator

import (
	"time"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/database"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

func AggregateByRooomSlotInRoom(
	_ config.WorkflowContext,
	days []ScheduleForSingleDayWithRooms,
	anagraphicsRef *parser.OutputAnagraphics,
) []ScheduleForSingleDayWithRoomsAndSlots {

	grouped := make([]*ScheduleForSingleDayWithRoomsAndSlots, 0)

	for _, daySchedule := range days {
		mapped := &ScheduleForSingleDayWithRoomsAndSlots{
			Day:            daySchedule.Day,
			VisitingGroups: daySchedule.VisitingGroups,
			RoomsSchedule:  make([]ScheduleForSingleDayAndRoomWithSlots, 0, len(daySchedule.RoomsSchedule)),
			StartAt:        daySchedule.StartAt,
			EndAt:          daySchedule.EndAt,
		}

		for _, roomSchedule := range daySchedule.RoomsSchedule {
			numSlotsToCreate := 1
			room := anagraphicsRef.Rooms[roomSchedule.RoomCode]
			if room.Slots > 0 {
				numSlotsToCreate = int(room.Slots)
			}

			mappedForThisRoom := ScheduleForSingleDayAndRoomWithSlots{
				RoomCode:           roomSchedule.RoomCode,
				Slots:              nil,
				NumTotalInAllSlots: 0,
			}

			for numSlotsToCreate > 0 {
				numSlotsToCreate--
				mappedForThisRoom.Slots = append(mappedForThisRoom.Slots, ScheduleForSingleDayAndRoomSlot{
					SlotIndex: len(mappedForThisRoom.Slots),
				})
			}

			for _, activity := range roomSchedule.Rows {

				// find which slot to fill
				fitsAt, _, fits := findBestSlotForActivity(activity, mappedForThisRoom.Slots)

				if fits {
					mappedForThisRoom.Slots[fitsAt].Rows = append(mappedForThisRoom.Slots[fitsAt].Rows, activity)
				} else {
					// have to create a new slot
					newSlot := ScheduleForSingleDayAndRoomSlot{
						SlotIndex: len(mappedForThisRoom.Slots),
						Rows:      nil,
					}
					newSlot.Rows = append(newSlot.Rows, activity)
					mappedForThisRoom.Slots = append(mappedForThisRoom.Slots, newSlot)
				}
				mappedForThisRoom.NumTotalInAllSlots++
			}

			mapped.RoomsSchedule = append(mapped.RoomsSchedule, mappedForThisRoom)
		}

		grouped = append(grouped, mapped)
	}

	out := make([]ScheduleForSingleDayWithRoomsAndSlots, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func findBestSlotForActivity(act OutputRow, slots []ScheduleForSingleDayAndRoomSlot) (int, int, bool) {
	if len(slots) == 0 {
		return -1, 0, false
	}

	scoreMap := make(map[int]int)
	numSlots := len(slots)

	scoreSettings := database.GetEffectiveSlotPlacementPreferencesForRoom(act.RoomCode)

	for slotIndex, slot := range slots {
		if !activityFitsInTime(act, slot.Rows, 0) {
			// no way the activity fits, not considering this slot.
			continue
		}

		score := 0

		// prefer slots where it fits with some clearance before and after
		if activityFitsInTime(act, slot.Rows, time.Minute*1) {
			score += scoreSettings.PointsForFittingWithSmallClearance
		}
		if activityFitsInTime(act, slot.Rows, time.Minute*30) {
			score += scoreSettings.PointsForFittingWithLotClearance
		}

		// prefer slots without activities on the right
		if slotIndex < numSlots-1 {
			if !activityFitsInTime(act, slots[slotIndex+1].Rows, 0) {
				score -= scoreSettings.PenaltyForActivityImmediatelyToTheRight
			}
		}
		if slotIndex < numSlots-2 {
			if !activityFitsInTime(act, slots[slotIndex+2].Rows, 0) {
				score -= scoreSettings.PenaltyForActivity2ndToTheRight
			}
		}

		// prefer slots without activities on the left
		if slotIndex > 0 {
			if !activityFitsInTime(act, slots[slotIndex-1].Rows, 0) {
				score -= scoreSettings.PenaltyForActivityImmediatelyToTheLeft
			}
		}

		// prefer slots where the same operator was
		if act.OperatorCode != "" && sameOperatorHasOtherActivitiesInSlot(act.OperatorCode, slot) {
			score += scoreSettings.PointsForOperatorHasOtherActivitiesInSlot
		}

		// prefer slots where the same group was
		if act.BookingCode != "" && sameGroupHasOtherActivitiesInSlot(act.BookingCode, slot) {
			score += scoreSettings.PointsForGroupHasOtherActivitiesInSlot
		}

		// prefer the more empty slots
		score -= scoreSettings.PenaltyForEachOtherActivityInSlot * len(slot.Rows)

		// prefer slots on the left when everything else is the same
		score -= slotIndex

		scoreMap[slotIndex] = score
	}

	if len(scoreMap) == 0 {
		// no valid slots
		return -1, 0, false
	}

	// pick the highest score
	highestIndex, highestScore := -1, 0
	for index, score := range scoreMap {
		if highestIndex < 0 || score > highestScore {
			highestScore = score
			highestIndex = index
		}
	}

	return highestIndex, highestScore, true
}

func sameOperatorHasOtherActivitiesInSlot(operatorCode string, slot ScheduleForSingleDayAndRoomSlot) bool {
	for _, otherAct := range slot.Rows {
		if otherAct.OperatorCode == operatorCode {
			return true
		}
	}
	return false
}

func sameGroupHasOtherActivitiesInSlot(code string, slot ScheduleForSingleDayAndRoomSlot) bool {
	for _, otherAct := range slot.Rows {
		if otherAct.BookingCode == code {
			return true
		}
	}
	return false
}

func activityFitsInTime(act OutputRow, occupied []OutputRow, clearance time.Duration) bool {
	if act.StartTime.IsZero() || act.EndTime.IsZero() {
		return true
	}
	actStartWithClearance := act.StartTime.Add(-clearance)
	actEndWithClearance := act.EndTime.Add(clearance)
	for _, otherAct := range occupied {
		if actStartWithClearance.Before(otherAct.EndTime) && otherAct.StartTime.Before(actEndWithClearance) {
			return false
		}
	}
	return true
}
