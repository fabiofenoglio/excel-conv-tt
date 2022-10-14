package excel

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/aggregator/byactivity"
	"github.com/fabiofenoglio/excelconv/aggregator/byday"
	"github.com/fabiofenoglio/excelconv/aggregator/byroom"
	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/fabiofenoglio/excelconv/model"
)

func writeDayGrid(c WriteContext, day byday.GroupedByDay, startCell excel.Cell, log *logrus.Logger) (DayGridWriteResult, error) {
	minGroupWidth := uint(12)
	minGroupWidthPerSlot := uint(5)
	minDayWidthInCells := uint(20)
	boxCellHeight := float64(20)
	moreSlotsAtBottom := 0 // add if you want to show some empty time rows after the last one

	cursor := startCell.Copy()
	f := c.outputFile
	zero := DayGridWriteResult{}

	cursor.MoveBottom(1)

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(),
		fmt.Sprintf("%d", day.Rows[0].StartAt.Year())); err != nil {
		return zero, err
	}
	cursor.MoveRight(1)
	cursor.MoveBottom(1)

	groupedByRoom := byroom.GroupByRoom(c.allData, day.Rows)

	roomColumnNumbers := make(map[string]uint)
	roomWidths := make(map[string]uint)

	// write the header with the columns / rooms
	for _, group := range groupedByRoom {
		log.Debugf("writing header for room group %s", group.Room.Code)

		numActs := uint(len(group.Slots))

		if group.Room.Code == "" && numActs == 0 {
			continue
		}
		if numActs == 0 && group.Room.Slots == 0 {
			continue
		}

		if group.Room.Slots > 0 && numActs < group.Room.Slots {
			numActs = group.Room.Slots
		}

		toWrite := group.Room.Name
		if group.Room.Code == "" {
			toWrite = "???"
		}

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return zero, err
		}

		roomColumnNumbers[group.Room.Code] = cursor.Column()
		roomWidths[group.Room.Code] = numActs

		desiredWidth := minGroupWidth

		if numActs > 1 {
			mergeUpTo := cursor.AtRight(numActs - 1)

			if err := f.MergeCell(cursor.SheetName(), cursor.Code(), mergeUpTo.Code()); err != nil {
				return zero, err
			}

			altWidth := numActs * minGroupWidthPerSlot
			if altWidth > desiredWidth {
				desiredWidth = altWidth
			}

			if err := f.SetColWidth(
				cursor.SheetName(),
				cursor.ColumnName(),
				mergeUpTo.ColumnName(),
				float64(desiredWidth)/float64(numActs),
			); err != nil {
				return zero, err
			}

		} else {
			if err := f.SetColWidth(cursor.SheetName(), cursor.ColumnName(), cursor.ColumnName(), float64(desiredWidth)); err != nil {
				return zero, err
			}
		}

		cursor.MoveRight(numActs)

		// add 1 additional column to leave some space
		if err := f.SetColWidth(cursor.SheetName(), cursor.ColumnName(), cursor.ColumnName(), 2); err != nil {
			return zero, err
		}

		cursor.MoveRight(1)
	}

	// ensure the day block is wide at least as the desired minDayWidthInCells
	actualWidth := cursor.Column() - startCell.Column()
	if actualWidth < minDayWidthInCells {
		diff := uint(minDayWidthInCells - actualWidth)

		cursorTo := cursor.AtRight(diff)
		if err := f.SetColWidth(cursor.SheetName(), cursor.ColumnName(), cursorTo.ColumnName(), float64(minGroupWidthPerSlot)); err != nil {
			return zero, err
		}

		cursor.MoveRight(diff)
	}
	actualWidth = cursor.Column() - startCell.Column()

	// reset cursor to box start
	cursor = startCell.AtBottom(3)

	minHourToShowForThisDay, maxHourToShowForThisDay := GetMaxHoursRangeToDisplay(day.Rows)

	currentTime := day.Rows[0].StartAt.Add(time.Minute * -30)
	if day.Rows[0].StartAt.Hour() < 2 {
		currentTime = day.Rows[0].StartAt
	}

	currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		c.minHour, 0, 0, 0, currentTime.Location())
	rolledOver := false
	firstTime := currentTime
	breaking := false

	numOfTimeRows := uint(0)
	log.Debugf("writing time columns starting from %s", currentTime.String())
	for {
		effectiveHour := currentTime.Hour()
		if rolledOver {
			effectiveHour += 24
		}

		show := true
		if minHourToShowForThisDay > c.minHour && effectiveHour < minHourToShowForThisDay {
			show = false
		}
		if maxHourToShowForThisDay < c.maxHour && effectiveHour > maxHourToShowForThisDay {
			show = false
		}

		toWrite := "--"
		if show {
			toWrite = currentTime.Format(layoutTimeOnlyInReadableFormat)
		}

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return zero, err
		}
		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), boxCellHeight); err != nil {
			return zero, err
		}
		cursor.MoveBottom(1)
		numOfTimeRows++

		currentTime = currentTime.Add(time.Duration(c.minutesStep) * time.Minute)

		if currentTime.Day() != firstTime.Day() {
			rolledOver = true
		}

		if !breaking {
			if effectiveHour > c.maxHour {
				breaking = true
			}
		}

		if breaking {
			if moreSlotsAtBottom > 0 {
				moreSlotsAtBottom--
			} else {
				break
			}
		}
	}

	// finished writing the times on the left

	// draw the box for the whole day
	boxStart := startCell.Copy()
	boxEnd := startCell.Copy().AtRight(actualWidth - 1).AtBottom(numOfTimeRows + 2)
	if err := f.SetCellStyle(
		startCell.SheetName(), boxStart.Code(), boxEnd.Code(),
		c.styleRegister.DayBoxStyle().Common.StyleID,
	); err != nil {
		return zero, err
	}

	// draw the box for each room
	for _, group := range groupedByRoom {
		log.Debugf("drawing boxes for room %s", group.Room.Code)

		if _, ok := roomWidths[group.Room.Code]; !ok {
			log.Debugf("skipping room %s", group.Room.Code)
			continue
		}

		columnsForRoomStartAt := roomColumnNumbers[group.Room.Code]

		boxStart := startCell.AtColumn(columnsForRoomStartAt).AtBottom(3)
		boxEnd := boxStart.AtRight(roomWidths[group.Room.Code] - 1).AtBottom(numOfTimeRows - 1)

		styleToUse := c.styleRegister.DayRoomBoxStyle().Common.StyleID
		if group.TotalInAllSlots == 0 {
			styleToUse = c.styleRegister.UnusedRoomStyle().Common.StyleID
		}
		if err := f.SetCellStyle(startCell.SheetName(), boxStart.Code(), boxEnd.Code(), styleToUse); err != nil {
			return zero, err
		}

		boxHeaderStart := boxStart.AtTop(1)
		boxHeaderEnd := boxHeaderStart.AtColumn(boxEnd.Column())

		var style *Style
		if group.Room != model.NoRoom {
			style = c.styleRegister.RoomStyle(group.Room.BackgroundColor)
		} else {
			style = c.styleRegister.NoRoomStyle()
		}

		if err := f.SetCellStyle(cursor.SheetName(), boxHeaderStart.Code(), boxHeaderEnd.Code(), style.Common.StyleID); err != nil {
			return zero, err
		}
	}

	// now place the activities
	cursor.MoveRow(startCell.Row())
	cursor.MoveColumn(startCell.Column())

	schoolGroups := byactivity.GetDifferentActivityGroups(day.Rows)
	schoolGroupsIndex := make(map[string]byactivity.ActivityGroup)
	for _, sg := range schoolGroups {
		if sg.Code == "" {
			continue
		}
		schoolGroupsIndex[sg.Code] = sg
	}

	// save a map of where each row is placed
	rowPlacementMap := make(map[int]excel.CellBox)

	for _, group := range groupedByRoom {
		log.Debugf("placing activities for room %s", group.Room.Code)

		for slotIndex, slot := range group.Slots {
			for _, act := range slot.Rows {
				var operator model.Operator
				if act.Operator.Code != "" {
					operator, _ = c.allData.OperatorsMap[act.Operator.Code]
				}
				if operator.Code == "" {
					operator.Name = "???"
				}

				columnsForRoomStartAt := roomColumnNumbers[act.Room.Code]
				cursor.MoveColumn(columnsForRoomStartAt + uint(slotIndex))
				cursor.MoveRow(startCell.Row() + 3)

				inRange := false
				exited := false
				timeCursor := firstTime
				i := 0

				actStartCell := cursor.Copy()
				actEndCell := cursor.Copy()

				for {
					if !inRange {
						if !exited && !timeCursor.Before(act.StartAt) {
							inRange = true
							actStartCell.MoveRow(startCell.Row() + 3 + uint(i))
						}
					} else {
						if !timeCursor.Before(act.EndAt) {
							inRange = false
							exited = true
							actEndCell.MoveRow(startCell.Row() + 3 + uint(i-1))
						}
					}

					if exited {
						break
					}

					timeCursor = timeCursor.Add(time.Minute * time.Duration(c.minutesStep))
					i++
				}

				// if merging the cells is NOT desired:
				for r := actStartCell.Row(); r <= actEndCell.Row(); r++ {
					c := actStartCell.AtRow(r)

					/*
						// writeInCell := operator.Name
						writeInCell := act.Raw.Codice
						if decoded, ok := schoolGroupsIndex[act.Raw.Codice]; ok && decoded.NumeroSeq > 0 {
							writeInCell = fmt.Sprintf("%d", decoded.NumeroSeq)
						}
					*/
					writeInCell := shortenGroupCode(act.Code)
					if writeInCell == "" && operator.Code != "" {
						writeInCell = operator.Name
					}

					if err := f.SetCellValue(c.SheetName(), c.Code(), writeInCell); err != nil {
						return zero, err
					}
				}

				// apply style depending on the operator
				var style *Style
				if act.Operator.Code == "" && (group.Room.AllowMissingOperator || !c.args.EnableMissingOperatorsWarning) {
					style = c.styleRegister.NoOperatorNeededStyle()
				} else if act.Operator.Code == "" {
					style = c.styleRegister.NoOperatorStyle()
				} else {
					style = c.styleRegister.OperatorStyle(operator.BackgroundColor)
				}

				styleID := style.Common.StyleID
				if len(act.Warnings) > 0 {
					styleID = style.WarningIDOrDefault()
				}

				if err := f.SetCellStyle(actStartCell.SheetName(), actStartCell.Code(), actEndCell.Code(), styleID); err != nil {
					return zero, err
				}

				// write some comment with additional info
				cellComment := buildContentOfActivityComment(act)
				if cellComment != "" {
					addCommentToCell(f, actStartCell, strings.TrimSpace(cellComment))
				}

				rowPlacementMap[act.ID] = excel.NewCellBox(actStartCell, actEndCell)
			}
		}
	}

	// write the day name in the header
	cursor = startCell.AtRight(1).AtBottom(1)
	toCell := cursor.AtRight(8)
	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), toCell.Code()); err != nil {
		return zero, err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(),
		day.Rows[0].StartAt.Format("Mon 2 January")); err != nil {
		return zero, err
	}
	if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(), c.styleRegister.DayHeaderStyle().Common.StyleID); err != nil {
		return zero, err
	}

	// write "MAT: # POM: #" counters
	groupedByActivityGroup := byactivity.GetDifferentActivityGroups(day.Rows)
	numGroupsMat, numGroupsPom := 0, 0
	for _, actGroup := range groupedByActivityGroup {
		if actGroup.AveragePresence.IsZero() {
			continue
		}
		if actGroup.AveragePresence.Hour() <= 13 {
			numGroupsMat++
		} else {
			numGroupsPom++
		}
	}

	cursor = startCell.AtRight(actualWidth - 8).AtRow(1)
	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
		return zero, err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "MAT:"); err != nil {
		return zero, err
	}
	if err := f.MergeCell(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).AtRight(2).Code()); err != nil {
		return zero, err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.AtBottom(1).Code(), fmt.Sprintf("%d", numGroupsMat)); err != nil {
		return zero, err
	}

	cursor = cursor.AtRight(4)
	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
		return zero, err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "POM:"); err != nil {
		return zero, err
	}
	if err := f.MergeCell(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).AtRight(2).Code()); err != nil {
		return zero, err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.AtBottom(1).Code(), fmt.Sprintf("%d", numGroupsPom)); err != nil {
		return zero, err
	}

	return DayGridWriteResult{
		RowsPlacement: rowPlacementMap,
	}, nil
}

func shortenGroupCode(code string) string {
	if len(code) >= 4 && strings.Contains(code, "-") && strings.HasPrefix(code, "P") {
		return strings.TrimPrefix(code, "P")
	}
	return code
}

type DayGridWriteResult struct {
	RowsPlacement map[int]excel.CellBox
}

func (r *DayGridWriteResult) Aggregate(other DayGridWriteResult) {
	if other.RowsPlacement == nil {
		return
	}
	if r.RowsPlacement == nil {
		r.RowsPlacement = make(map[int]excel.CellBox)
	}
	for k, v := range other.RowsPlacement {
		r.RowsPlacement[k] = v
	}
}
