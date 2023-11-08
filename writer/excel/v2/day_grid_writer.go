package excel

import (
	"fmt"
	"strings"
	"time"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/excel"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
)

func writeDayGrid(ctx config.WorkflowContext, c WriteContext, day aggregator2.ScheduleForSingleDayWithRoomsAndGroupSlots, startCell excel.Cell) (DayGridWriteResult, error) {
	log := ctx.Logger

	minGroupWidth := uint(12)
	minGroupWidthPerSlot := uint(5)
	minDayWidthInCells := uint(23)
	boxCellHeight := float64(20)
	moreSlotsAtBottom := 0 // add if you want to show some empty time rows after the last one

	cursor := startCell.Copy()
	f := c.outputFile
	zero := DayGridWriteResult{}

	cursor.MoveBottom(1)

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(),
		fmt.Sprintf("%d", day.Day.Year())); err != nil {
		return zero, err
	}
	cursor.MoveRight(1)
	cursor.MoveBottom(1)

	roomColumnNumbers := make(map[string]uint)
	roomWidths := make(map[string]uint)

	anagraphics := c.anagraphicsRef

	// write the header with the columns / rooms
	for _, group := range day.RoomsSchedule {
		room := anagraphics.Rooms[group.RoomCode]
		if room.Hide {
			continue
		}
		if !room.AlwaysShow && group.NumTotalInAllSlots == 0 {
			continue
		}
		log.Debugf("writing header for room group %s", room.Code)

		numActs := uint(len(group.Slots))

		if group.RoomCode == "" && numActs == 0 {
			continue
		}
		if numActs == 0 && len(group.Slots) == 0 {
			continue
		}

		if len(group.Slots) > 0 && int(numActs) < len(group.Slots) {
			numActs = uint(len(group.Slots))
		}

		toWrite := strings.ToUpper(room.Name)
		if room.Code == "" {
			toWrite = "???"
		}

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return zero, err
		}

		roomColumnNumbers[group.RoomCode] = cursor.Column()
		roomWidths[group.RoomCode] = numActs

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

	relativeToDay := func(startTime time.Time, competenceDate time.Time) aggregator2.TimeOfDay {
		if startTime.YearDay() == competenceDate.YearDay() {
			return aggregator2.TimeOfDay{Hour: startTime.Hour(), Minute: startTime.Minute()}
		}
		return aggregator2.TimeOfDay{Hour: startTime.Hour() + 24, Minute: startTime.Minute()}
	}

	minHourToShowForThisDay := relativeToDay(day.StartAt, day.Day)
	maxHourToShowForThisDay := relativeToDay(day.EndAt, day.Day)

	currentTime := time.Date(day.StartAt.Year(), day.StartAt.Month(), day.StartAt.Day(),
		c.minHour.Hour, c.minHour.Minute, 0, 0, day.StartAt.Location())

	firstTime := currentTime
	breaking := false

	numOfTimeRows := uint(0)
	log.Debugf("writing time columns starting from %s", currentTime.String())
	for {
		effectiveTime := relativeToDay(currentTime, day.Day)

		show := true
		if minHourToShowForThisDay.IsAfter(c.minHour) && effectiveTime.IsBefore(minHourToShowForThisDay) {
			show = false
		}
		if maxHourToShowForThisDay.IsBefore(c.maxHour) && effectiveTime.IsAfter(maxHourToShowForThisDay) {
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

		if !breaking {
			if effectiveTime.IsAfter(c.maxHour) {
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
		c.styleRegister.DayBoxStyle().SingleCell(),
	); err != nil {
		return zero, err
	}

	// draw the box for each room
	for _, group := range day.RoomsSchedule {
		room := anagraphics.Rooms[group.RoomCode]
		if room.Hide {
			continue
		}
		if !room.AlwaysShow && group.NumTotalInAllSlots == 0 {
			continue
		}

		log.Debugf("drawing boxes for room %s", room.Code)

		if _, ok := roomWidths[room.Code]; !ok {
			log.Debugf("skipping room %s", room.Code)
			continue
		}

		columnsForRoomStartAt := roomColumnNumbers[room.Code]

		boxStart := startCell.AtColumn(columnsForRoomStartAt).AtBottom(3)
		boxEnd := boxStart.AtRight(roomWidths[room.Code] - 1).AtBottom(numOfTimeRows - 1)

		styleToUse := c.styleRegister.DayRoomBoxStyle()
		if group.NumTotalInAllSlots == 0 {
			styleToUse = c.styleRegister.UnusedRoomStyle()
		}
		if err := f.SetCellStyle(startCell.SheetName(), boxStart.Code(), boxEnd.Code(), styleToUse.SingleCell()); err != nil {
			return zero, err
		}

		boxHeaderStart := boxStart.AtTop(1)
		boxHeaderEnd := boxHeaderStart.AtColumn(boxEnd.Column())

		var style *RegisteredStyleV2
		if room.Code != "" {
			style = c.styleRegister.RoomStyle(room.BackgroundColor)
		} else {
			style = c.styleRegister.NoRoomStyle()
		}

		if err := f.SetCellStyle(cursor.SheetName(), boxHeaderStart.Code(), boxHeaderEnd.Code(), style.SingleCell()); err != nil {
			return zero, err
		}
	}

	// now place the activities
	cursor.MoveRow(startCell.Row())
	cursor.MoveColumn(startCell.Column())

	// save a map of where each row is placed
	rowPlacementMap := make(map[int]excel.CellBox)

	schoolGroupsIndex := make(map[string]aggregator2.VisitingGroupInDay)
	for _, v := range day.VisitingGroups {
		schoolGroupsIndex[v.VisitingGroupCode] = v
	}

	for _, group := range day.RoomsSchedule {
		room := anagraphics.Rooms[group.RoomCode]
		if room.Hide {
			continue
		}
		if !room.AlwaysShow && group.NumTotalInAllSlots == 0 {
			continue
		}

		log.Debugf("placing activities for room %s", room.Code)

		columnsForRoomStartAt := roomColumnNumbers[group.RoomCode]

		for slotIndex, slot := range group.Slots {
			for _, act := range slot.GroupedActivities {
				if act.StartingSlotIndex != slot.SlotIndex {
					continue
				}
				var operators []parser2.Operator
				operatorCodes := act.OperatorCodes()
				for _, operatorCode := range operatorCodes {
					operators = append(operators, anagraphics.Operators[operatorCode])
				}
				// bookingCodes := act.BookingCodes()
				visitingGroupCodes := act.VisitingGroupCodes()

				cursor.MoveColumn(columnsForRoomStartAt + uint(slotIndex))
				cursor.MoveRow(startCell.Row() + 3)

				inRange := false
				exited := false
				timeCursor := firstTime
				i := 0

				actStartCell := cursor.Copy()
				actEndCell := cursor.Copy().AtRight(uint(act.NumOccupiedSlots - 1))

				for {
					if !inRange {
						if !exited && !timeCursor.Before(act.StartTime) {
							inRange = true
							actStartCell.MoveRow(startCell.Row() + 3 + uint(i))
						}
					} else {
						if !timeCursor.Before(act.EndTime) {
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

				// apply style depending on the operator
				var style *RegisteredStyleV2
				if len(operators) == 0 {
					if room.AllowMissingOperator || !ctx.Config.EnableMissingOperatorsWarning {
						style = c.styleRegister.NoOperatorNeededStyle()
					} else {
						style = c.styleRegister.NoOperatorStyle()
					}
				} else {
					if len(operators) == 1 {
						style = c.styleRegister.OperatorStyle(operators[0].BackgroundColor)
					} else {
						style = c.styleRegister.NoOperatorNeededStyle()
					}
				}
				if len(act.Warnings()) > 0 {
					style = style.WithWarning()
				}

				// decide how to annotate this activity group depending on available rows/columns
				availableRows := actEndCell.Row() - actStartCell.Row() + 1
				availableColumns := int(actEndCell.Column() - actStartCell.Column() + 1)
				haveMultipleLines := availableRows > 1
				haveAtLeastOneColumnPerVisitingGroup := availableColumns >= len(visitingGroupCodes)

				if haveMultipleLines && haveAtLeastOneColumnPerVisitingGroup && room.ShowActivityNamesInside {
					// first row to (lat - 1) row: annotate the activities,
					// last row: one visiting group code per column
					activityNameHeaderBox := excel.NewCellBox(
						actStartCell.Copy(),
						actEndCell.AtTop(1),
					)

					if err := f.MergeCell(
						actStartCell.SheetName(),
						activityNameHeaderBox.TopLeft().Code(),
						activityNameHeaderBox.BottomRight().Code(),
					); err != nil {
						return zero, err
					}

					toWrite := ""
					for _, activityCode := range act.ActivityCodes() {
						if activityRef, ok := anagraphics.Activities[activityCode]; ok {
							toWrite += activityRef.Name
							if activityRef.Language != "" && strings.ToLower(activityRef.Language) != "it" {
								toWrite += "(" + strings.ToUpper(activityRef.Language) + ")"
							}
							toWrite += ", "
						}
					}
					toWrite = strings.TrimSuffix(toWrite, ", ")

					if len(act.Warnings()) > 0 {
						toWrite = "⚠️ " + toWrite
					}

					if err := f.SetCellValue(
						actStartCell.SheetName(),
						activityNameHeaderBox.TopLeft().Code(),
						toWrite,
					); err != nil {
						return zero, err
					}

					var additionalStyles []designatedStyling

					for singleActIndex, singleActivity := range act.Rows {
						var groupRef parser2.VisitingGroup
						toWriteForThisVisitingGroup := singleActivity.BookingCode
						if visitingGroupRef, ok := schoolGroupsIndex[singleActivity.VisitingGroupCode]; ok {
							toWriteForThisVisitingGroup = visitingGroupRef.DisplayCode
							groupRef = anagraphics.VisitingGroups[visitingGroupRef.VisitingGroupCode]
						}
						cellForThis := actStartCell.AtRow(actEndCell.Row()).AtRight(uint(singleActIndex))
						if err := f.SetCellValue(
							actStartCell.SheetName(),
							cellForThis.Code(),
							toWriteForThisVisitingGroup,
						); err != nil {
							return zero, err
						}

						if len(groupRef.Highlights) > 0 {
							if additionalStyle := highlightsToStyle(c, groupRef.Highlights); additionalStyle != nil {
								additionalStyles = append(additionalStyles, designatedStyling{
									additionalStyle,
									excel.NewCellBox(cellForThis, cellForThis),
								})
							}
						}

						cellComment := buildContentOfActivityCommentForSingleGroupOfGroupedActivity(c, singleActivity)
						if cellComment != "" {
							if err := addCommentToCell(f, cellForThis, strings.TrimSpace(cellComment)); err != nil {
								return zero, err
							}
						}
					}

					if err := applyStyleToBox(f, style, excel.NewCellBox(actStartCell, actEndCell)); err != nil {
						return zero, err
					}

					for _, additional := range additionalStyles {
						merged := c.styleRegister.Merge(style, additional.style)

						if err := applyStyleToBox(f, merged, additional.box); err != nil {
							return zero, err
						}
					}

					// write some comment with additional info
					cellComment := buildContentOfActivityComment(c, act)
					if cellComment != "" {
						if err := addCommentToCell(f, actStartCell, strings.TrimSpace(cellComment)); err != nil {
							return zero, err
						}
					}

				} else {
					// write cell-by-cell, normally

					for slotCnt := 0; slotCnt < len(act.Rows); slotCnt++ {
						singleAct := act.Rows[slotCnt]
						var writeInCell string
						var groupRef parser2.VisitingGroup
						if singleAct.VisitingGroupCode != "" {
							if visitingGroup, ok := schoolGroupsIndex[singleAct.VisitingGroupCode]; ok {
								writeInCell = visitingGroup.DisplayCode
								groupRef = anagraphics.VisitingGroups[visitingGroup.VisitingGroupCode]
							} else {
								writeInCell = singleAct.VisitingGroupCode
							}
						} else {
							writeInCell = singleAct.BookingCode
						}

						for r := actStartCell.Row(); r <= actEndCell.Row(); r++ {
							effectiveWrite := writeInCell
							if len(act.Warnings()) > 0 && slotCnt == 0 && actEndCell.Row() > actStartCell.Row() && r == actStartCell.Row() {
								effectiveWrite = "⚠️"
							}

							if err := f.SetCellValue(
								actStartCell.SheetName(),
								actStartCell.AtRow(r).AtRight(uint(slotCnt)).Code(),
								effectiveWrite,
							); err != nil {
								return zero, err
							}
						}

						if len(groupRef.Highlights) > 0 {
							if additionalStyle := highlightsToStyle(c, groupRef.Highlights); additionalStyle != nil {
								style = c.styleRegister.Merge(style, additionalStyle)
							}
						}
					}

					if err := applyStyleToBox(f, style, excel.NewCellBox(actStartCell, actEndCell)); err != nil {
						return zero, err
					}

					// write some comment with additional info
					cellComment := buildContentOfActivityComment(c, act)
					if cellComment != "" {
						if err := addCommentToCell(f, actStartCell, strings.TrimSpace(cellComment)); err != nil {
							return zero, err
						}
					}

					for _, singleAct := range act.Rows {
						// TODO why not done in the other drawing branch?
						rowPlacementMap[singleAct.ID] = excel.NewCellBox(actStartCell, actEndCell)
					}
				}
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
		day.Day.Format("Mon 2 January")); err != nil {
		return zero, err
	}
	if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(), c.styleRegister.DayHeaderStyle().SingleCell()); err != nil {
		return zero, err
	}

	return DayGridWriteResult{
		RowsPlacement: rowPlacementMap,
	}, nil
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
