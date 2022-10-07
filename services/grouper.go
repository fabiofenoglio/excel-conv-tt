package services

import (
	"sort"
	"strings"
)

func GroupByStartDay(rows []ParsedRow) []GroupedRows {

	// group by start date, ordering each group by start time ASC, end time ASC

	grouped := make([]*GroupedRows, 0)
	index := make(map[string]*GroupedRows)
	for _, activity := range rows {
		if activity.StartAt.IsZero() {
			continue
		}
		key := activity.StartAt.Format(layoutDateOrderable)

		group, ok := index[key]
		if !ok {
			group = &GroupedRows{
				Key: key,
			}
			index[key] = group
			grouped = append(grouped, group)
		}
		group.Rows = append(group.Rows, activity)
	}

	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Key < grouped[j].Key
	})

	for _, group := range grouped {
		sort.Slice(group.Rows, func(i, j int) bool {
			if group.Rows[i].StartAt.UnixMilli() < group.Rows[j].StartAt.UnixMilli() {
				return true
			}
			if group.Rows[i].StartAt.UnixMilli() > group.Rows[j].StartAt.UnixMilli() {
				return false
			}
			return group.Rows[i].EndAt.UnixMilli() < group.Rows[j].EndAt.UnixMilli()
		})
	}

	out := make([]GroupedRows, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func GroupByRoom(allData Parsed, rows []ParsedRow) []GroupedRows {

	// group by room, ordering each group by start time ASC, end time ASC

	grouped := make([]*GroupedRows, 0)
	index := make(map[string]*GroupedRows)
	for _, activity := range rows {
		key := activity.Room.Code

		group, ok := index[key]
		if !ok {
			group = &GroupedRows{
				Key: key,
			}
			index[key] = group
			grouped = append(grouped, group)
		}
		group.Rows = append(group.Rows, activity)
	}

	for _, knownRoom := range allData.Rooms {
		if _, ok := index[knownRoom.Code]; !ok {
			group := &GroupedRows{
				Key: knownRoom.Code,
			}
			index[knownRoom.Code] = group
			grouped = append(grouped, group)
		}
	}

	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Key < grouped[j].Key
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

	out := make([]GroupedRows, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}
