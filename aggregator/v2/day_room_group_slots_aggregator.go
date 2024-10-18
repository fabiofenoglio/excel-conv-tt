package aggregator

import (
	"fmt"
	"time"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/database"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

func AggregateByRooomGroupSlotInRoom(
	_ config.WorkflowContext,
	days []ScheduleForSingleDayWithRoomsAndGroupedActivities,
	anagraphicsRef *parser.OutputAnagraphics,
) []ScheduleForSingleDayWithRoomsAndGroupSlots {

	grouped := make([]*ScheduleForSingleDayWithRoomsAndGroupSlots, 0)

	for _, daySchedule := range days {
		mapped := &ScheduleForSingleDayWithRoomsAndGroupSlots{
			Day:                                   daySchedule.Day,
			VisitingGroups:                        daySchedule.VisitingGroups,
			RoomsSchedule:                         make([]ScheduleForSingleDayAndRoomWithGroupSlots, 0, len(daySchedule.RoomsSchedule)),
			StartAt:                               daySchedule.StartAt,
			EndAt:                                 daySchedule.EndAt,
			NumeroAttivitaMarkers:                 daySchedule.NumeroAttivitaMarkers,
			NumeroGruppiAttivitaConfermateMarkers: daySchedule.NumeroGruppiAttivitaConfermateMarkers,
		}

		for _, roomSchedule := range daySchedule.RoomsSchedule {

			numSlotsToCreate := 1
			room := anagraphicsRef.Rooms[roomSchedule.RoomCode]
			if room.Slots > 0 {
				numSlotsToCreate = int(room.Slots)
			}

			mappedForThisRoom := ScheduleForSingleDayAndRoomWithGroupSlots{
				RoomCode:           roomSchedule.RoomCode,
				Slots:              nil,
				NumTotalInAllSlots: 0,
			}

			for numSlotsToCreate > 0 {
				numSlotsToCreate--
				mappedForThisRoom.Slots = append(mappedForThisRoom.Slots, ScheduleForSingleDayAndRoomGroupSlot{
					SlotIndex: len(mappedForThisRoom.Slots),
				})
			}

			for _, activity := range roomSchedule.GroupedActivities {

				var numSlotsToSpan int
				var fitsAt int
				var fitComputationLog []string

				if room.GroupActivities && len(mappedForThisRoom.Slots) > 1 && noOtherActivitiesInSlots(activity, roomSchedule) {
					// take all the slots as they are all available
					fitsAt = 0
					numSlotsToSpan = len(mappedForThisRoom.Slots)

				} else {
					var fits bool
					// find which slot to fill
					fitsAt, _, fits, fitComputationLog = findBestSlotForGroupedActivity(activity, mappedForThisRoom.Slots, roomSchedule)
					if !fits {
						// will force to create a new slot
						fitsAt = len(mappedForThisRoom.Slots)
					}
					numSlotsToSpan = len(activity.Rows)
				}

				activity.StartingSlotIndex = fitsAt
				activity.FitComputationLog = fitComputationLog
				activity.NumOccupiedSlots = numSlotsToSpan

				for cnt := 0; cnt < numSlotsToSpan; cnt++ {
					effectiveIndex := fitsAt + cnt

					if len(mappedForThisRoom.Slots) <= effectiveIndex {
						newSlot := ScheduleForSingleDayAndRoomGroupSlot{
							SlotIndex:         len(mappedForThisRoom.Slots),
							GroupedActivities: nil,
						}
						mappedForThisRoom.Slots = append(mappedForThisRoom.Slots, newSlot)
					}

					mappedForThisRoom.Slots[effectiveIndex].GroupedActivities =
						append(mappedForThisRoom.Slots[effectiveIndex].GroupedActivities, activity)
				}

				mappedForThisRoom.NumTotalInAllSlots++
			}

			mapped.RoomsSchedule = append(mapped.RoomsSchedule, mappedForThisRoom)
		}

		grouped = append(grouped, mapped)
	}

	out := make([]ScheduleForSingleDayWithRoomsAndGroupSlots, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func findBestSlotForGroupedActivity(
	act GroupedActivity,
	slots []ScheduleForSingleDayAndRoomGroupSlot,
	parent ScheduleForSingleDayAndRoomWithGroupedActivities,
) (int, int, bool, []string) {
	if len(slots) == 0 {
		return -1, 0, false, nil
	}

	scoreMap := make(map[int]int)
	logMap := make([]string, 0)
	numSlots := len(slots)
	numSlotsForThisAct := len(act.Rows)

	scoreSettings := database.GetEffectiveSlotPlacementPreferencesForRoom(parent.RoomCode)

	doLog := false // should come from input

	for slotIndex, slot := range slots {
		if !groupedActivityFitsInTime(act, slotIndex, slots, 0) {
			// no way the activity fits, not considering this slot.
			continue
		}

		score := 0

		apply := func(target *int, amount int, reason string) {
			if amount == 0 {
				return
			}
			newValue := *target + amount
			if doLog {
				logMap = append(logMap,
					fmt.Sprintf("[%d] %v (%v -> %v) - %s", slotIndex, amount, *target, newValue, reason))
			}
			*target = newValue
		}

		// prefer slots where it fits with some clearance before and after
		if groupedActivityFitsInTime(act, slotIndex, slots, time.Minute*1) {
			apply(&score, scoreSettings.PointsForFittingWithSmallClearance, "has small clearance")
		}
		if groupedActivityFitsInTime(act, slotIndex, slots, time.Minute*30) {
			apply(&score, scoreSettings.PointsForFittingWithLotClearance, "has lot of clearance")
		}

		// prefer slots without activities on the right
		if slotIndex < numSlots-numSlotsForThisAct {
			if !groupedActivityFitsInTime(act, slotIndex+1, slots, 0) {
				apply(&score, -scoreSettings.PenaltyForActivityImmediatelyToTheRight, "has activity on the right")
			}
		}
		if slotIndex < numSlots-numSlotsForThisAct-1 {
			if !groupedActivityFitsInTime(act, slotIndex+2, slots, 0) {
				apply(&score, -scoreSettings.PenaltyForActivity2ndToTheRight, "has activity on the right (2nd)")
			}
		}

		// prefer slots without activities on the left
		if slotIndex > 0 {
			if !groupedActivityFitsInTime(act, slotIndex-1, slots, 0) {
				apply(&score, -scoreSettings.PenaltyForActivityImmediatelyToTheLeft, "has activity on the left")
			}
		}

		// prefer slots where the same operator was
		if len(act.OperatorCodes()) > 0 && sameOperatorHasOtherGroupedActivitiesInSlot(act.OperatorCodes(), slot) {
			apply(&score, scoreSettings.PointsForOperatorHasOtherActivitiesInSlot, "operator has other activities in slot")
		}

		// prefer slots where the same group was
		if len(act.BookingCodes()) > 0 && sameGroupHasOtherGroupedActivitiesInSlot(act.BookingCodes(), slot) {
			apply(&score, scoreSettings.PointsForGroupHasOtherActivitiesInSlot, "group has other activities in slot")
		}

		// prefer the more empty slots
		diff := -scoreSettings.PenaltyForEachOtherActivityInSlot * len(slot.GroupedActivities)
		apply(&score, diff, "other activities in the slot")

		// prefer slots on the left when everything else is the same
		diff = -slotIndex
		apply(&score, diff, "prefer leftmost positions")

		scoreMap[slotIndex] = score
	}

	if len(scoreMap) == 0 {
		// no valid slots
		return -1, 0, false, nil
	}

	// pick the highest score
	highestIndex, highestScore := -1, 0
	for index, score := range scoreMap {
		if highestIndex < 0 || score > highestScore {
			highestScore = score
			highestIndex = index
		}
	}

	return highestIndex, highestScore, true, logMap
}

func noOtherActivitiesInSlots(
	act GroupedActivity,
	parent ScheduleForSingleDayAndRoomWithGroupedActivities) bool {

	if act.StartTime.IsZero() || act.EndTime.IsZero() {
		return false
	}

	for _, groupedActivity := range parent.GroupedActivities {
		if groupedActivity.ID == act.ID {
			continue
		}

		if act.StartTime.Before(groupedActivity.EndTime) && groupedActivity.StartTime.Before(act.EndTime) {
			return false
		}
	}

	return true
}

func sameOperatorHasOtherGroupedActivitiesInSlot(operatorCodes []string, slot ScheduleForSingleDayAndRoomGroupSlot) bool {
	for _, operatorCode := range operatorCodes {
		for _, otherAct := range slot.GroupedActivities {
			for _, otherActOperatorCode := range otherAct.OperatorCodes() {
				if otherActOperatorCode == operatorCode {
					return true
				}
			}
		}
	}

	return false
}

func sameGroupHasOtherGroupedActivitiesInSlot(codes []string, slot ScheduleForSingleDayAndRoomGroupSlot) bool {
	for _, operatorCode := range codes {
		for _, otherAct := range slot.GroupedActivities {
			for _, otherActOperatorCode := range otherAct.BookingCodes() {
				if otherActOperatorCode == operatorCode {
					return true
				}
			}
		}
	}

	return false
}

func groupedActivityFitsInTime(act GroupedActivity, slotIndex int, slots []ScheduleForSingleDayAndRoomGroupSlot, clearance time.Duration) bool {
	if act.StartTime.IsZero() || act.EndTime.IsZero() {
		return true
	}
	actStartWithClearance := act.StartTime.Add(-clearance)
	actEndWithClearance := act.EndTime.Add(clearance)

	numToCheck := len(act.Rows)
	i := slotIndex

	for numToCheck > 0 {
		if i >= len(slots) {
			// not enough slots to fit
			return false
		}

		for _, otherAct := range slots[i].GroupedActivities {
			if actStartWithClearance.Before(otherAct.EndTime) && otherAct.StartTime.Before(actEndWithClearance) {
				return false
			}
		}

		i++
		numToCheck--
	}

	return true
}
