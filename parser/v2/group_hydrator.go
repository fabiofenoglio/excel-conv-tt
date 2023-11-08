package parser

import (
	"strings"

	"github.com/fabiofenoglio/excelconv/config"
)

func HydrateGroups(
	_ config.WorkflowContext,
	rows []Row,
) (
	[]Row,
	[]VisitingGroup,
	[]School,
	[]SchoolClass,
	error,
) {
	outRows := make([]Row, 0, len(rows))
	outGroups := make([]VisitingGroup, 0, 10)
	outSchools := make([]School, 0, 10)
	outSchoolClasses := make([]SchoolClass, 0, 10)

	groupsIndex := make(map[string]int)
	schoolsIndex := make(map[string]School)
	schoolClassesIndex := make(map[string]SchoolClass)

	for _, row := range rows {
		mappedRow := row

		groupCode := nameToCode(row.BookingCode)

		if groupCode != "" {
			mappedRow.VisitingGroupCode = groupCode

			groupMapped, groupAlreadyMapped := groupsIndex[groupCode]
			if !groupAlreadyMapped {

				schoolCode := nameToCode(row.schoolType) + "/" + nameToCode(row.schoolName)
				schoolGroupCode := schoolCode + "/" + nameToCode(row.class) + "/" + nameToCode(row.classSection) + "/" + row.BookingCode

				newGroup := VisitingGroup{
					Code:            groupCode,
					SchoolCode:      schoolCode,
					SchoolClassCode: schoolGroupCode,
					Composition: GroupComposition{
						NumPaying:       row.numPaying,
						NumFree:         row.numFree,
						NumAccompanying: row.numAccompanying,
					},
					ClassTeacher:        row.classTeacher,
					ClassRefEmail:       row.classRefEmail,
					BookingNotes:        row.BookingNote,
					OperatorNotes:       row.OperatorNote,
					SpecialProjectNotes: row.SpecialProjectName,
				}
				outGroups = append(outGroups, newGroup)
				groupsIndex[newGroup.Code] = len(outGroups) - 1

				if _, isAlreadyMapped := schoolsIndex[schoolCode]; !isAlreadyMapped {
					newSchool := School{
						Code: schoolCode,
						Type: cleanStringForVisualization(row.schoolType),
						Name: cleanStringForVisualization(row.schoolName),
					}
					schoolsIndex[schoolCode] = newSchool
					outSchools = append(outSchools, newSchool)
				}

				if _, isAlreadyMapped := schoolClassesIndex[schoolGroupCode]; !isAlreadyMapped {
					newClass := SchoolClass{
						Code:       schoolGroupCode,
						SchoolCode: schoolCode,
						Number:     cleanStringForVisualization(row.class),
						Section:    cleanStringForVisualization(row.classSection),
					}
					schoolClassesIndex[schoolGroupCode] = newClass
					outSchoolClasses = append(outSchoolClasses, newClass)
				}

			} else {
				toEnrich := outGroups[groupMapped]
				if toEnrich.Code != groupCode {
					panic("invalid reference")
				}

				if row.OperatorNote != "" && !strings.Contains(toEnrich.OperatorNotes, row.OperatorNote) {
					toEnrich.OperatorNotes = strings.TrimPrefix(toEnrich.OperatorNotes+"\n"+row.OperatorNote, "\n")
				}
				if row.BookingNote != "" && !strings.Contains(toEnrich.BookingNotes, row.BookingNote) {
					toEnrich.BookingNotes = strings.TrimPrefix(toEnrich.BookingNotes+"\n"+row.BookingNote, "\n")
				}
				if row.SpecialProjectName != "" && !strings.Contains(toEnrich.SpecialProjectNotes, row.SpecialProjectName) {
					toEnrich.SpecialProjectNotes = strings.TrimPrefix(toEnrich.SpecialProjectNotes+"\n"+row.SpecialProjectName, "\n")
				}
				if row.classTeacher != "" && toEnrich.ClassTeacher == "" {
					toEnrich.ClassTeacher = row.classTeacher
				}
				if row.classRefEmail != "" && toEnrich.ClassRefEmail == "" {
					toEnrich.ClassRefEmail = row.classRefEmail
				}
				outGroups[groupMapped] = toEnrich
			}
		}

		outRows = append(outRows, mappedRow)
	}

	return outRows, outGroups, outSchools, outSchoolClasses, nil
}
