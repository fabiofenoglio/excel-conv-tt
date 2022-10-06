package services

func GetMaxHoursRangeToDisplay(rows []ParsedRow) (start, end int) {

	// compute the max time range to be shown between all days

	minHourToShow := 12
	maxHourToShow := 13
	for _, activity := range rows {
		if !activity.StartAt.IsZero() && activity.StartAt.Hour() < minHourToShow {
			minHourToShow = activity.StartAt.Hour()
		}
		if !activity.StartAt.IsZero() && !activity.EndAt.IsZero() && activity.StartAt.Day() != activity.EndAt.Day() {
			// rolls over
			maxHourToShow = 24 + activity.EndAt.Hour()
		} else {
			if !activity.EndAt.IsZero() && activity.EndAt.Hour() > maxHourToShow {
				maxHourToShow = activity.EndAt.Hour()
			}
		}
	}

	return minHourToShow, maxHourToShow
}
