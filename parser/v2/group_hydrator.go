package parser

import (
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

	groupsIndex := make(map[string]VisitingGroup)
	schoolsIndex := make(map[string]School)
	schoolClassesIndex := make(map[string]SchoolClass)

	for _, row := range rows {
		mappedRow := row

		groupCode := nameToCode(row.BookingCode)

		if groupCode != "" {
			mappedRow.VisitingGroupCode = groupCode

			if _, groupAlreadyMapped := groupsIndex[groupCode]; !groupAlreadyMapped {

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
				}
				groupsIndex[newGroup.Code] = newGroup
				outGroups = append(outGroups, newGroup)

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

			}
		}

		outRows = append(outRows, mappedRow)
	}

	return outRows, outGroups, outSchools, outSchoolClasses, nil
}
