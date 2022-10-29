package parser

import (
	"github.com/fabiofenoglio/excelconv/config"
)

func HydrateActivities(
	_ config.WorkflowContext,
	rows []Row,
) (
	[]Row,
	[]Activity,
	[]ActivityType,
	error,
) {
	outRows := make([]Row, 0, len(rows))
	outActivities := make([]Activity, 0, 10)
	outActivityTypes := make([]ActivityType, 0, 10)

	activitiesIndex := make(map[string]Activity)
	activityTypesIndex := make(map[string]ActivityType)

	for _, row := range rows {
		mappedRow := row

		// no natural UUID for these fields, computing one on-the-fly based on attributes
		activityTypeCode := nameToCode(row.activityTypeRawString)
		activityCode := activityTypeCode + "/" + nameToCode(row.activityRawString) + "/?lang=" + nameToCode(row.activityLanguageRawString)

		if activityCode != "" {
			mappedRow.ActivityCode = activityCode

			if _, activityAlreadyMapped := activitiesIndex[activityCode]; !activityAlreadyMapped {
				newActivity := Activity{
					Code:     activityCode,
					TypeCode: activityTypeCode,
					Name:     cleanStringForVisualization(row.activityRawString),
					Language: cleanStringForVisualization(row.activityLanguageRawString),
				}
				activitiesIndex[activityCode] = newActivity
				outActivities = append(outActivities, newActivity)

				if _, isAlreadyMapped := activityTypesIndex[activityTypeCode]; !isAlreadyMapped {
					newActivityType := ActivityType{
						Code: activityTypeCode,
						Name: cleanStringForVisualization(row.activityTypeRawString),
					}
					activityTypesIndex[activityTypeCode] = newActivityType
					outActivityTypes = append(outActivityTypes, newActivityType)
				}
			}
		}

		outRows = append(outRows, mappedRow)
	}

	return outRows, outActivities, outActivityTypes, nil
}
