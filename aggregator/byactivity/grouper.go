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

	notesIndex := make(map[string]bool)

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
				StartsAt:         row.StartAt,
			}
			if group.StartsAt.IsZero() || (!row.StartAt.IsZero() && row.StartAt.Before(group.StartsAt)) {
				group.StartsAt = row.StartAt
			}
			index[key] = group
			grouped = append(grouped, group)
		}

		if row.BookingNotes != "" {
			codeForNote := row.Code + "/bk/" + strings.ToUpper(strings.TrimSpace(row.BookingNotes))
			if _, alreadyAdded := notesIndex[codeForNote]; !alreadyAdded {
				notesIndex[codeForNote] = true
				if group.Notes != "" {
					group.Notes += "\n"
				}
				group.Notes += row.BookingNotes
			}
		}

		if row.OperatorNotes != "" {
			codeForNote := row.Code + "/op/" + strings.ToUpper(strings.TrimSpace(row.OperatorNotes))
			if _, alreadyAdded := notesIndex[codeForNote]; !alreadyAdded {
				notesIndex[codeForNote] = true
				if group.Notes != "" {
					group.Notes += "\n"
				}
				group.Notes += row.OperatorNotes
			}
		}
	}

	sort.SliceStable(grouped, func(i, j int) bool {
		if grouped[i].StartsAt.Before(grouped[j].StartsAt) {
			return true
		} else if grouped[i].StartsAt.After(grouped[j].StartsAt) {
			return false
		}
		c := strings.Compare(strings.ToLower(grouped[i].Code), strings.ToLower(grouped[j].Code))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].School.SortableIdentifier()), strings.ToLower(grouped[j].School.SortableIdentifier()))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].SchoolClass.SortableIdentifier()), strings.ToLower(grouped[j].SchoolClass.SortableIdentifier()))
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
	indexForNumberingSchools := make(map[string]uint)
	indexForNumberingClasses := make(map[string]uint)
	counterForSchool := uint(1)
	countersForClassInsideSchool := make(map[string]uint)

	out := make([]ActivityGroup, 0, len(grouped))
	for i, e := range grouped {
		// check if this school has a number already
		k := e.School.Hash()
		numForSchool, ok := indexForNumberingSchools[k]
		if !ok {
			numForSchool = counterForSchool
			counterForSchool++
			indexForNumberingSchools[k] = numForSchool
		}

		// check if this school group has a number already
		k = e.Code + "/" + e.SchoolClass.Hash()
		numForGroupInsideSchool, ok := indexForNumberingClasses[k]
		if !ok {
			var countingAlready bool
			numForGroupInsideSchool, countingAlready = countersForClassInsideSchool[e.School.Hash()]
			if !countingAlready {
				numForGroupInsideSchool = 1
			}
			countersForClassInsideSchool[e.School.Hash()] = numForGroupInsideSchool + 1
			indexForNumberingClasses[k] = numForGroupInsideSchool
		}

		e.SequentialNumberInsideSchool = numForGroupInsideSchool
		e.SequentialNumberForSchool = numForSchool
		e.SequentialNumber = uint(i + 1)

		out = append(out, *e)
	}

	return out
}
