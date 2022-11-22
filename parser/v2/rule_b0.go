package parser

import (
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/fabiofenoglio/excelconv/config"
)

func ApplyRuleB0Level(
	ctx config.WorkflowContext,
	rows []Row,
) ([]Row, error) {
	rows, err := applyRuleB0LevelRule0(ctx, rows)
	if err != nil {
		return nil, err
	}

	/*
		rows, err = applyRuleB0LevelRule1(ctx, rows)
		if err != nil {
			return nil, err
		}
	*/

	return rows, nil
}

func applyRuleB0LevelRule0(
	ctx config.WorkflowContext,
	rows []Row,
) ([]Row, error) {
	logger := ctx.Logger
	activityCoupleRegex := regexp.MustCompile(`(?m)^\s*[\s\*]*\s*([\s0-9hmHM\,\.]{2,})\s*\+\s*(.*)`)

	updatedIndex := make(map[int]bool)

	for i, row := range rows {
		if updatedIndex[i] {
			continue
		}
		if row.activityRawString == "" {
			continue
		}
		if !activityCoupleRegex.MatchString(row.activityRawString) {
			continue
		}
		if row.RoomCode != "museo" && row.RoomCode != "planetario" {
			continue
		}

		originalActivityName := row.activityRawString
		matches := activityCoupleRegex.FindStringSubmatch(originalActivityName)
		if len(matches) != 3 {
			continue
		}

		foundOtherRow := false
		matchingRow := Row{}
		matchingIndex := 0

		for j, otherRow := range rows {
			if i == j {
				continue
			}
			if updatedIndex[j] {
				continue
			}
			if otherRow.Date != row.Date {
				continue
			}
			if otherRow.BookingCode != row.BookingCode {
				continue
			}
			if otherRow.activityRawString != row.activityRawString {
				continue
			}
			if otherRow.RoomCode != "museo" && otherRow.RoomCode != "planetario" {
				continue
			}
			if otherRow.RoomCode == row.RoomCode {
				continue
			}

			if !foundOtherRow {
				foundOtherRow = true
				matchingRow = otherRow
				matchingIndex = j
			} else {
				// too many matches
				foundOtherRow = false
				break
			}
		}

		if !foundOtherRow {
			continue
		}

		// found another row that matches the criteria:
		// - same group, same day
		// - same booked activity name
		// - bookd activity name in the form "aaaa + bbbb"
		// - different room
		// - rooms are museo & planetario (or inverse)
		if row.RoomCode == "museo" {
			row.activityRawString = originalActivityName
			matchingRow.activityRawString = strings.TrimSpace(matches[2])
		} else {
			matchingRow.activityRawString = originalActivityName
			row.activityRawString = strings.TrimSpace(matches[2])
		}

		logger.Infof("rewrote activity %v from [%s] in room [%s] to [%s]",
			rows[i].ID, rows[i].activityRawString, rows[i].RoomCode, row.activityRawString)
		logger.Infof("rewrote activity %v from [%s] in room [%s] to [%s]",
			rows[matchingIndex].ID, rows[matchingIndex].activityRawString, rows[matchingIndex].RoomCode, matchingRow.activityRawString)

		rows[i] = row
		rows[matchingIndex] = matchingRow

		updatedIndex[i] = true
		updatedIndex[matchingIndex] = true
	}

	return rows, nil
}

//nolint:unused,deadcode
func applyRuleB0LevelRule1(
	ctx config.WorkflowContext,
	rows []Row,
) ([]Row, error) {
	logger := ctx.Logger
	distincts := make(map[string]string)

	stringSimilarityMetrics := metrics.NewJaroWinkler()

	clean := func(raw string) string {
		return strings.ToUpper(strings.TrimSpace(raw))
	}

	for i, row := range rows {
		cleanedActivityName := clean(row.activityRawString)
		if cleanedActivityName == "" {
			continue
		}

		rewrote := false
		for otherKey, otherName := range distincts {
			if otherKey == cleanedActivityName {
				continue
			}

			similarity := strutil.Similarity(cleanedActivityName, otherKey, stringSimilarityMetrics)
			if similarity >= 0.90 {
				row.activityRawString = otherName
				logger.Infof("normalized activity %v name from [%s] to [%s] (confidence: %v)",
					row.ID, rows[i].activityRawString, row.activityRawString, similarity)
				rows[i] = row
				rewrote = true
				break
			}
		}

		if !rewrote {
			distincts[cleanedActivityName] = row.activityRawString
		}
	}

	return rows, nil
}
