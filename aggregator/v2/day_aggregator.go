package aggregator

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

const (
	layoutDateOrderable = "2006-01-02"
)

func AggregateByCompetenceDay(_ config.WorkflowContext, rows []Row, anagraphicsRef *parser.OutputAnagraphics) []scheduleForSingleDay {
	// group by competence date, ordering each group by start time ASC, end time ASC

	grouped := make([]*scheduleForSingleDay, 0)
	index := make(map[string]*scheduleForSingleDay)

	// create a group for each competence day
	for _, row := range rows {
		rowCopy := row
		if rowCopy.CompetenceDate.IsZero() {
			continue
		}
		key := rowCopy.CompetenceDate.Format(layoutDateOrderable)

		group, ok := index[key]
		if !ok {
			group = &scheduleForSingleDay{
				Day:                             rowCopy.CompetenceDate,
				Rows:                            nil,
				VisitingGroups:                  nil,
				NumeroAttivitaMarkers:           make(map[time.Time]int),
				NumeroAttivitaConfermateMarkers: make(map[time.Time]int),
			}
			index[key] = group
			grouped = append(grouped, group)
		}

		if row.InputRow.IsPlaceholderNumeroAttivita {
			if !row.InputRow.StartTime.IsZero() {
				group.NumeroAttivitaMarkers[row.InputRow.StartTime]++
				if row.InputRow.Confirmed != nil && *row.InputRow.Confirmed {
					group.NumeroAttivitaConfermateMarkers[row.InputRow.StartTime]++
				}
			}
			continue
		}

		if !row.InputRow.StartTime.IsZero() && (group.StartAt.IsZero() || row.InputRow.StartTime.Before(group.StartAt)) {
			group.StartAt = row.InputRow.StartTime
		}
		if !row.InputRow.EndTime.IsZero() && (group.StartAt.IsZero() || row.InputRow.EndTime.Before(group.StartAt)) {
			group.StartAt = row.InputRow.EndTime
		}

		if !row.InputRow.StartTime.IsZero() && (group.EndAt.IsZero() || row.InputRow.StartTime.After(group.EndAt)) {
			group.EndAt = row.InputRow.StartTime
		}
		if !row.InputRow.EndTime.IsZero() && (group.EndAt.IsZero() || row.InputRow.EndTime.After(group.EndAt)) {
			group.EndAt = row.InputRow.EndTime
		}

		group.Rows = append(group.Rows, rowCopy)
	}

	// populate the day groups' rows VisitingGroups field
	// and compute its StartsAt attribute in the same scan
	for _, group := range grouped {
		groupsIndex := make(map[string]*VisitingGroupInDay)

		for _, row := range group.Rows {
			groupCode := row.InputRow.VisitingGroupCode
			if groupCode == "" {
				continue
			}

			visitingGroup, ok := groupsIndex[groupCode]
			if !ok {
				visitingGroup = &VisitingGroupInDay{
					VisitingGroupCode: groupCode,
				}
				groupsIndex[groupCode] = visitingGroup
			}

			if !row.InputRow.StartTime.IsZero() && (visitingGroup.StartsAt.IsZero() || row.InputRow.StartTime.Before(visitingGroup.StartsAt)) {
				visitingGroup.StartsAt = row.InputRow.StartTime
			}
		}

		for _, vg := range groupsIndex {
			group.VisitingGroups = append(group.VisitingGroups, *vg)
		}

	}

	// sort day groups by competence day
	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Day.Before(grouped[j].Day)
	})

	for _, group := range grouped {
		// sort each day group's rows by StartTime, EndTime, BookingCode, ID
		sort.Slice(group.Rows, func(i, j int) bool {
			if group.Rows[i].InputRow.StartTime.UnixMilli() < group.Rows[j].InputRow.StartTime.UnixMilli() {
				return true
			}
			if group.Rows[i].InputRow.StartTime.UnixMilli() > group.Rows[j].InputRow.StartTime.UnixMilli() {
				return false
			}
			if group.Rows[i].InputRow.EndTime.UnixMilli() < group.Rows[j].InputRow.EndTime.UnixMilli() {
				return true
			}
			if group.Rows[i].InputRow.EndTime.UnixMilli() > group.Rows[j].InputRow.EndTime.UnixMilli() {
				return false
			}
			diff := strings.Compare(group.Rows[i].InputRow.BookingCode, group.Rows[j].InputRow.BookingCode)
			if diff != 0 {
				return diff < 0
			}
			return group.Rows[i].InputRow.ID < group.Rows[j].InputRow.ID
		})

		// sort each day group's visiting groups by StartsAt, School, SchoolClass Code
		sort.Slice(group.VisitingGroups, func(i, j int) bool {
			diff := group.VisitingGroups[i].StartsAt.UnixMilli() - group.VisitingGroups[j].StartsAt.UnixMilli()
			if diff != 0 {
				return diff < 0
			}

			groupRefI := anagraphicsRef.VisitingGroups[group.VisitingGroups[i].VisitingGroupCode]
			schoolRefI := anagraphicsRef.Schools[groupRefI.SchoolCode]
			schoolClassRefI := anagraphicsRef.SchoolClasses[groupRefI.SchoolClassCode]

			groupRefJ := anagraphicsRef.VisitingGroups[group.VisitingGroups[j].VisitingGroupCode]
			schoolRefJ := anagraphicsRef.Schools[groupRefJ.SchoolCode]
			schoolClassRefJ := anagraphicsRef.SchoolClasses[groupRefJ.SchoolClassCode]

			diffStr := strings.Compare(schoolRefI.SortableIdentifier(), schoolRefJ.SortableIdentifier())
			if diffStr != 0 {
				return diffStr < 0
			}
			diffStr = strings.Compare(schoolRefI.Code, schoolRefJ.Code)
			if diffStr != 0 {
				return diffStr < 0
			}

			diffStr = strings.Compare(schoolClassRefI.SortableIdentifier(), schoolClassRefJ.SortableIdentifier())
			if diffStr != 0 {
				return diffStr < 0
			}
			diffStr = strings.Compare(schoolClassRefI.Code, schoolClassRefJ.Code)
			if diffStr != 0 {
				return diffStr < 0
			}

			return strings.Compare(group.VisitingGroups[i].VisitingGroupCode, group.VisitingGroups[j].VisitingGroupCode) < 0
		})

		// assign sequential numbers to school groups based on the order just sorted
		schoolNumbers := make(map[string]int)
		schoolNumberProgressive := 0
		schoolGroupsNumbers := make(map[string]int)
		schoolGroupsNumbersProgressive := make(map[string]int)
		for i, visitingGroup := range group.VisitingGroups {

			groupRef := anagraphicsRef.VisitingGroups[visitingGroup.VisitingGroupCode]
			schoolRef := anagraphicsRef.Schools[groupRef.SchoolCode]
			schoolClassRef := anagraphicsRef.SchoolClasses[groupRef.SchoolClassCode]

			schoolNumAssigned, ok := schoolNumbers[schoolRef.Code]
			if !ok {
				schoolNumAssigned = schoolNumberProgressive + 1
				schoolNumberProgressive++
				schoolNumbers[schoolRef.Code] = schoolNumAssigned
			}

			schoolGroupNumberAssigned, ok := schoolGroupsNumbers[schoolClassRef.Code]
			if !ok {
				schoolGroupNumberAssigned = schoolGroupsNumbersProgressive[schoolRef.Code] + 1
				schoolGroupsNumbersProgressive[schoolRef.Code] = schoolGroupNumberAssigned
				schoolGroupsNumbers[schoolClassRef.Code] = schoolGroupNumberAssigned
			}

			group.VisitingGroups[i].SequentialCode = fmt.Sprintf("%09d-%09d", schoolNumAssigned, schoolGroupNumberAssigned)
			group.VisitingGroups[i].DisplayCode = fmt.Sprintf("%d-%s", schoolNumAssigned, numToChars(uint(schoolGroupNumberAssigned)))
		}
	}

	out := make([]scheduleForSingleDay, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}
