package aggregator

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
			group = &GroupedByRoom{
				Room:            activity.Room,
				TotalInAllSlots: 0,
				Slots:           nil,
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

	// as preferred method of placement we want a new slot
	// if we can add a new one.
	if act.Room.Slots > 0 && len(slots) < int(act.Room.Slots) {
		// we have more available slots before reaching the limit.
		// first things
		// return saying that we didn't find a suitable spot and another slot should be added
		return -1, false
	}

	// TODO: not just the first slot available but the one available
	// will less activities

	// first attempt: look for available slot in slots that are there already,
	// with time clearance before and after to keep the grid as clean as possible.
	fits := false
	fitsAt := 0
	for slotIndex, slot := range slots {
		fits = activityFitsInTime(act, slot.Rows, time.Minute*1)
		if fits {
			fitsAt = slotIndex
			break
		}
	}

	// found some available place without adding more slots,
	// with clearance, returning immediately as that's the optimal outcome
	if fits {
		return fitsAt, true
	}

	if act.Room.Slots > 0 && len(slots) < int(act.Room.Slots) {
		// we have more available slots before reaching the limit.
		// return saying that we didn't find a suitable spot and another slot should be added
		return -1, false
	}

	// we didn't find an optimal place in an existing spot,
	// we don't have more room for available slots
	// so we try to fit without time clearance
	// this way we avoid adding slots if not absolutely necessary.
	fits = false
	fitsAt = 0
	for slotIndex, slot := range slots {
		fits = activityFitsInTime(act, slot.Rows, 0)
		if fits {
			fitsAt = slotIndex
			break
		}
	}
	if fits {
		return fitsAt, true
	}
	return -1, false
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
