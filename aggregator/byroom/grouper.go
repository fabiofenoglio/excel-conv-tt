package byroom

import (
	"sort"
	"strings"
	"time"

	"github.com/fabiofenoglio/excelconv/model"
)

func GroupByRoom(allData model.ParsedData, rows []model.ParsedRow) []GroupedByRoom {

	activitiesByTime := sortActivitiesByStartTime(rows)

	// group by room, ordering each group by start time ASC, end time ASC

	grouped := make([]*GroupedByRoom, 0)
	index := make(map[string]*GroupedByRoom)
	for _, activity := range activitiesByTime {
		key := activity.Room.Code

		group, ok := index[key]
		if !ok {
			numSlotsToCreate := uint(1)
			if activity.Room.Slots > 0 {
				numSlotsToCreate = activity.Room.Slots
			}
			group = &GroupedByRoom{
				Room:            activity.Room,
				TotalInAllSlots: 0,
				Slots:           make([]GroupedByRoomSlot, numSlotsToCreate),
			}
			index[key] = group
			grouped = append(grouped, group)
		}

		// find which slot to fill
		fitsAt, fits := getIndexOfAvailableSlot(activity, group.Slots)

		if fits {
			group.Slots[fitsAt].Rows = append(group.Slots[fitsAt].Rows, activity)
		} else {
			// have to create a new slot
			newSlot := GroupedByRoomSlot{
				SlotIndex: uint(len(group.Slots)),
				Rows:      nil,
			}
			newSlot.Rows = append(newSlot.Rows, activity)
			group.Slots = append(group.Slots, newSlot)
		}

		group.TotalInAllSlots++
	}

	for _, knownRoom := range allData.Rooms {
		if _, ok := index[knownRoom.Code]; !ok {
			group := &GroupedByRoom{
				Room: knownRoom,
			}
			index[knownRoom.Code] = group
			grouped = append(grouped, group)
		}
	}

	sort.Slice(grouped, func(i, j int) bool {
		preferredSortRes := grouped[i].Room.PreferredOrder - grouped[j].Room.PreferredOrder
		if preferredSortRes != 0 {
			return preferredSortRes < 0
		}
		return grouped[i].Room.Code < grouped[j].Room.Code
	})

	for _, group := range grouped {
		sort.Slice(group.Slots, func(i, j int) bool {
			return group.Slots[i].SlotIndex < group.Slots[j].SlotIndex
		})

		for _, slot := range group.Slots {
			sort.Slice(slot.Rows, func(i, j int) bool {
				if slot.Rows[i].StartAt.UnixMilli() < slot.Rows[j].StartAt.UnixMilli() {
					return true
				}
				if slot.Rows[i].StartAt.UnixMilli() > slot.Rows[j].StartAt.UnixMilli() {
					return false
				}
				if slot.Rows[i].EndAt.UnixMilli() < slot.Rows[j].EndAt.UnixMilli() {
					return true
				}
				if slot.Rows[i].EndAt.UnixMilli() > slot.Rows[j].EndAt.UnixMilli() {
					return false
				}

				c := strings.Compare(slot.Rows[i].Operator.Code, slot.Rows[j].Operator.Code)
				return c < 0
			})
		}

	}

	out := make([]GroupedByRoom, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func getIndexOfAvailableSlot(act model.ParsedRow, slots []GroupedByRoomSlot) (int, bool) {
	if len(slots) == 0 {
		return -1, false
	}

	scoreMap := make(map[int]int)
	numSlots := len(slots)

	type scoresTable struct {
		pointsForFittingWithSmallClearance        int
		pointsForFittingWithLotClearance          int
		penaltyForActivityImmediatelyToTheRight   int
		penaltyForActivity2ndToTheRight           int
		penaltyForActivityImmediatelyToTheLeft    int
		pointsForOperatorHasOtherActivitiesInSlot int
		penaltyForEachOtherActivityInSlot         int
	}
	scoreSettings := scoresTable{
		pointsForFittingWithSmallClearance:        50,
		pointsForFittingWithLotClearance:          10,
		penaltyForActivityImmediatelyToTheRight:   30,
		penaltyForActivity2ndToTheRight:           5,
		penaltyForActivityImmediatelyToTheLeft:    0,
		pointsForOperatorHasOtherActivitiesInSlot: 30,
		penaltyForEachOtherActivityInSlot:         5,
	}

	for slotIndex, slot := range slots {
		if !activityFitsInTime(act, slot.Rows, 0) {
			// no way the activity fits, not considering this slot.
			continue
		}

		score := 0

		// prefer slots where it fits with some clearance before and after
		if activityFitsInTime(act, slot.Rows, time.Minute*1) {
			score += scoreSettings.pointsForFittingWithSmallClearance
		}
		if activityFitsInTime(act, slot.Rows, time.Minute*30) {
			score += scoreSettings.pointsForFittingWithLotClearance
		}

		// prefer slots without activities on the right
		if slotIndex < numSlots-1 {
			if !activityFitsInTime(act, slots[slotIndex+1].Rows, 0) {
				score -= scoreSettings.penaltyForActivityImmediatelyToTheRight
			}
		}
		if slotIndex < numSlots-2 {
			if !activityFitsInTime(act, slots[slotIndex+2].Rows, 0) {
				score -= scoreSettings.penaltyForActivity2ndToTheRight
			}
		}

		// prefer slots without activities on the left
		if slotIndex > 0 {
			if !activityFitsInTime(act, slots[slotIndex-1].Rows, 0) {
				score -= scoreSettings.penaltyForActivityImmediatelyToTheLeft
			}
		}

		// prefer slots where the same operator was
		if act.Operator.Code != "" && sameOperatorHasOtherActivitiesInSlot(act.Operator.Code, slot) {
			score += scoreSettings.pointsForOperatorHasOtherActivitiesInSlot
		}

		// prefer the more empty slots
		score -= scoreSettings.penaltyForEachOtherActivityInSlot * len(slot.Rows)

		// prefer slots on the left when everything else is the same
		score -= slotIndex

		scoreMap[slotIndex] = score
	}

	if len(scoreMap) == 0 {
		// no valid slots
		return -1, false
	}

	// pick the highest score
	highestIndex, highestScore := -1, 0
	for index, score := range scoreMap {
		if highestIndex < 0 || score > highestScore {
			highestScore = score
			highestIndex = index
		}
	}

	return highestIndex, true
}

func sameOperatorHasOtherActivitiesInSlot(operatorCode string, slot GroupedByRoomSlot) bool {
	for _, otherAct := range slot.Rows {
		if otherAct.Operator.Code == operatorCode {
			return true
		}
	}
	return false
}

func activityFitsInTime(act model.ParsedRow, occupied []model.ParsedRow, clearance time.Duration) bool {
	if act.StartAt.IsZero() || act.EndAt.IsZero() {
		return true
	}
	actStartWithClearance := act.StartAt.Add(-clearance)
	actEndWithClearance := act.EndAt.Add(clearance)
	for _, otherAct := range occupied {
		if actStartWithClearance.Before(otherAct.EndAt) && otherAct.StartAt.Before(actEndWithClearance) {
			return false
		}
	}
	return true
}

func sortActivitiesByStartTime(rows []model.ParsedRow) []model.ParsedRow {

	// PRE-SORT all activities by startTime, endTime
	activitiesByTime := make([]model.ParsedRow, 0, len(rows))
	for _, r := range rows {
		activitiesByTime = append(activitiesByTime, r)
	}

	sort.Slice(activitiesByTime, func(i, j int) bool {
		res := activitiesByTime[i].StartAt.UnixMilli() - activitiesByTime[j].StartAt.UnixMilli()
		if res != 0 {
			return res < 0
		}
		res = activitiesByTime[i].EndAt.UnixMilli() - activitiesByTime[j].EndAt.UnixMilli()
		if res != 0 {
			return res < 0
		}
		return activitiesByTime[i].Code < activitiesByTime[j].Code
	})

	return activitiesByTime
}
