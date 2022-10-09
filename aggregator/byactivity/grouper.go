package byactivity

import (
	"sort"
	"strings"
	"time"

	"github.com/fabiofenoglio/excelconv/aggregator"
	"github.com/fabiofenoglio/excelconv/model"
)

func GetDifferentActivityGroups(rows []model.ParsedRow) []ActivityGroup {
	index := make(map[string]*ActivityGroup)
	grouped := make([]*ActivityGroup, 0)

	for _, row := range rows {
		if row.Code == "" && row.School.FullDescription() == "" && row.SchoolClass.FullDescription() == "" {
			continue
		}
		keyBuilder := row.Code
		if keyBuilder == "" {
			keyBuilder = row.School.FullDescription() + "|" + row.SchoolClass.FullDescription()
		}
		key := aggregator.Base64Sha([]byte(strings.ToLower(keyBuilder)))

		group, ok := index[key]
		if !ok {
			group = &ActivityGroup{
				Code:             row.Code,
				SequentialNumber: 0,
				School:           row.School,
				SchoolClass:      row.SchoolClass,
				Composition:      row.GroupComposition,
			}
			index[key] = group
			grouped = append(grouped, group)
		}
	}

	sort.SliceStable(grouped, func(i, j int) bool {
		c := strings.Compare(strings.ToLower(grouped[i].Code), strings.ToLower(grouped[j].Code))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].School.Name), strings.ToLower(grouped[j].School.Name))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].SchoolClass.FullDescription()), strings.ToLower(grouped[j].SchoolClass.FullDescription()))
		if c != 0 {
			return c < 0
		}
		return false
	})

	// compute AveragePresence
	for _, actGroup := range grouped {
		var min, max time.Time

		for _, act := range rows {
			if act.Code != actGroup.Code || act.StartAt.IsZero() || act.EndAt.IsZero() {
				continue
			}
			if min.IsZero() || act.StartAt.Before(min) {
				min = act.StartAt
			}
			if max.IsZero() || act.EndAt.After(max) {
				max = act.EndAt
			}
		}

		if min.IsZero() || max.IsZero() || min == max {
			continue
		}

		middle := min.Add(max.Sub(min) / 2)
		actGroup.AveragePresence = middle
	}

	// compute sequential number
	out := make([]ActivityGroup, 0, len(grouped))
	for i, e := range grouped {
		e.SequentialNumber = uint(i + 1)
		out = append(out, *e)
	}

	return out
}
