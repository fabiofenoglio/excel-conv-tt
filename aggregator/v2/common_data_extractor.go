package aggregator

import (
	"time"

	"github.com/fabiofenoglio/excelconv/config"
)

func ExtractCommonData(_ config.WorkflowContext, rows []Row) CommonData {

	return CommonData{
		CommonTimespan: extractCommonTimespan(rows),
	}
}

func extractCommonTimespan(rows []Row) CommonTimespan {

	// compute the max time range to be shown between all days
	minHourToShow := TimeOfDay{12, 0}
	maxHourToShow := TimeOfDay{14, 0}

	relativeToDay := func(startTime time.Time, competenceDate time.Time) TimeOfDay {
		if startTime.YearDay() == competenceDate.YearDay() {
			return TimeOfDay{startTime.Hour(), startTime.Minute()}
		}
		return TimeOfDay{startTime.Hour() + 24, startTime.Minute()}
	}

	for _, row := range rows {
		activity := row.InputRow
		if !activity.StartTime.IsZero() {
			startRelative := relativeToDay(activity.StartTime, row.CompetenceDate)
			if startRelative.IsBefore(minHourToShow) {
				minHourToShow = startRelative
			}
			if startRelative.IsAfter(maxHourToShow) {
				maxHourToShow = startRelative
			}
		}

		if !activity.EndTime.IsZero() {
			endRelative := relativeToDay(activity.EndTime, row.CompetenceDate)
			if endRelative.IsBefore(minHourToShow) {
				minHourToShow = endRelative
			}
			if endRelative.IsAfter(maxHourToShow) {
				maxHourToShow = endRelative
			}
		}
	}

	return CommonTimespan{
		Start: minHourToShow,
		End:   maxHourToShow,
	}

}
