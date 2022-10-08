package aggregator

import (
	"sort"
	"strings"

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
				Slots: []GroupedByRoomSlot{
					{
						SlotIndex: 0,
						Rows:      nil,
					},
				},
			}
			index[key] = group
			grouped = append(grouped, group)
		}

		// find which slot to fill
		fits := false
		fitsAt := 0
		for slotIndex, slot := range group.Slots {
			fits = activityFitsInTime(activity, slot.Rows)
			if fits {
				fitsAt = slotIndex
				break
			}
		}

		if fits {
			group.Slots[fitsAt].Rows = append(group.Slots[fitsAt].Rows, activity)
		} else {
			// have to create a new slot
			newSlot := GroupedByRoomSlot{
				SlotIndex: uint(len(group.Slots)),
				Rows:      nil,
			}
			group.Slots = append(group.Slots, newSlot)
			newSlot.Rows = append(newSlot.Rows, activity)
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

func activityFitsInTime(act model.ParsedRow, occupied []model.ParsedRow) bool {
	if act.StartAt.IsZero() || act.EndAt.IsZero() {
		return true
	}
	for _, otherAct := range occupied {
		if act.StartAt.Before(otherAct.EndAt) && otherAct.StartAt.Before(act.EndAt) {
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
