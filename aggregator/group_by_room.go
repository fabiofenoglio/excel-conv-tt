package aggregator

import (
	"sort"
	"strings"

	"github.com/fabiofenoglio/excelconv/model"
)

func GroupByRoom(allData model.ParsedData, rows []model.ParsedRow) []GroupedByRoom {

	// group by room, ordering each group by start time ASC, end time ASC

	grouped := make([]*GroupedByRoom, 0)
	index := make(map[string]*GroupedByRoom)
	for _, activity := range rows {
		key := activity.Room.Code

		group, ok := index[key]
		if !ok {
			group = &GroupedByRoom{
				Room: activity.Room,
			}
			index[key] = group
			grouped = append(grouped, group)
		}
		group.Rows = append(group.Rows, activity)
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
		sort.Slice(group.Rows, func(i, j int) bool {
			if group.Rows[i].StartAt.UnixMilli() < group.Rows[j].StartAt.UnixMilli() {
				return true
			}
			if group.Rows[i].StartAt.UnixMilli() > group.Rows[j].StartAt.UnixMilli() {
				return false
			}
			if group.Rows[i].EndAt.UnixMilli() < group.Rows[j].EndAt.UnixMilli() {
				return true
			}
			if group.Rows[i].EndAt.UnixMilli() > group.Rows[j].EndAt.UnixMilli() {
				return false
			}

			c := strings.Compare(group.Rows[i].Operator.Code, group.Rows[j].Operator.Code)
			return c < 0
		})
	}

	out := make([]GroupedByRoom, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}
