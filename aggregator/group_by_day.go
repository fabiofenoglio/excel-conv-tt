package aggregator

import (
	"sort"

	"github.com/fabiofenoglio/excelconv/model"
)

func GroupByStartDay(rows []model.ParsedRow) []GroupedByDay {

	// group by start date, ordering each group by start time ASC, end time ASC

	grouped := make([]*GroupedByDay, 0)
	index := make(map[string]*GroupedByDay)
	for _, activity := range rows {
		if activity.StartAt.IsZero() {
			continue
		}
		key := activity.StartAt.Format(layoutDateOrderable)

		group, ok := index[key]
		if !ok {
			group = &GroupedByDay{
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

	out := make([]GroupedByDay, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}
