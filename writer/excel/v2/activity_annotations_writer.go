package excel

import (
	"fmt"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/excel"
)

func writeAnnotationsOnActivities(
	ctx config.WorkflowContext,
	c WriteContext,
	groupByDay aggregator2.ScheduleForSingleDayWithRoomsAndGroupSlots,
	aggregateDayWriter DayGridWriteResult,
) error {
	rowPlacementMap := aggregateDayWriter.RowsPlacement
	f := c.outputFile

	alreadyAnnotatedIndex := make(map[int]bool)

	// write annotations where possible
	for _, roomSchedule := range groupByDay.RoomsSchedule {
		for _, slot := range roomSchedule.Slots {
			for _, group := range slot.GroupedActivities {
				for _, act := range group.Rows {
					roomRef := c.anagraphicsRef.Rooms[act.RoomCode]
					if !roomRef.ShowActivityNamesAsAnnotations {
						continue
					}

					if _, ok := alreadyAnnotatedIndex[act.ID]; ok {
						// already annotated by other things
						continue
					}

					actBox, ok := rowPlacementMap[act.ID]
					if !ok {
						ctx.Logger.Warningf(
							"could not annotate activity %v because no cell placement was registered", act.ID)
						continue
					}

					activityRef := c.anagraphicsRef.Activities[act.ActivityCode]

					if activityRef.Name == "" {
						continue
					}

					var bestSpace AvailableFreeSpace

					foundSpaces, err := findContiguousFreeBoxes(c, actBox, AvailableFreeSpaceSearchSpec{
						minWidth:  10,
						maxWidth:  30,
						minHeight: 15,
						maxHeight: 45,
					})
					if err != nil {
						return fmt.Errorf("error looking up spaces to write annotations: %w", err)
					}

					if len(foundSpaces) == 0 {
						// TODO do something else,
						// like for instance adding to the cell text and increasing the column width
						continue
					}

					bestSpace = findBestFreeBox(c, foundSpaces)

					if err := f.MergeCell(actBox.TopLeft().SheetName(), bestSpace.box.TopLeft().Code(), bestSpace.box.BottomRight().Code()); err != nil {
						return err
					}
					toWrite := activityRef.Name
					if bestSpace.arrowAnnotation != "" {
						if bestSpace.arrowAtEnd {
							toWrite = toWrite + " " + bestSpace.arrowAnnotation
						} else {
							toWrite = bestSpace.arrowAnnotation + " " + toWrite
						}
					}
					if err := f.SetCellValue(actBox.TopLeft().SheetName(), bestSpace.box.TopLeft().Code(), toWrite); err != nil {
						return err
					}

					style := c.styleRegister.InBoxAnnotationOnLeftStyle()
					if bestSpace.arrowAtEnd {
						style = c.styleRegister.InBoxAnnotationOnRightStyle()
					}
					if err := f.SetCellStyle(actBox.TopLeft().SheetName(), bestSpace.box.TopLeft().Code(), bestSpace.box.BottomRight().Code(),
						style.Common.StyleID); err != nil {
						return err
					}

					alreadyAnnotatedIndex[act.ID] = true

					// explore all connected activities to avoid repeated annotations for adjacent groups of activities

					var explore func(from aggregator2.OutputRow, box excel.CellBox)

					explore = func(from aggregator2.OutputRow, box excel.CellBox) {

						fromOperatorRef := c.anagraphicsRef.Operators[from.OperatorCode]
						fromRoomRef := c.anagraphicsRef.Rooms[from.RoomCode]
						fromActivityRef := c.anagraphicsRef.Activities[from.ActivityCode]

						for _, otherRoom := range groupByDay.RoomsSchedule {
							for _, otherSlot := range otherRoom.Slots {
								for _, otherActGroup := range otherSlot.GroupedActivities {
									for _, otherAct := range otherActGroup.Rows {
										if otherAct.ID == from.ID {
											continue
										}
										if _, ok := alreadyAnnotatedIndex[otherAct.ID]; ok {
											continue
										}

										otherOperatorRef := c.anagraphicsRef.Operators[otherAct.OperatorCode]
										otherRoomRef := c.anagraphicsRef.Rooms[otherAct.RoomCode]
										otherActivityRef := c.anagraphicsRef.Activities[otherAct.ActivityCode]

										if otherRoomRef.Code != fromRoomRef.Code {
											continue
										}
										if otherOperatorRef.Code != fromOperatorRef.Code {
											continue
										}
										if otherActivityRef.Code != fromActivityRef.Code {
											continue
										}
										otherActivityBox, ok := rowPlacementMap[otherAct.ID]
										if !ok {
											continue
										}
										if !isAdjacent(box, otherActivityBox) {
											continue
										}

										alreadyAnnotatedIndex[otherAct.ID] = true

										explore(otherAct, otherActivityBox)
									}
								}
							}
						}
					}

					explore(act, actBox)
				}
			}

		}
	}

	return nil
}

func isAdjacent(rectangle, otherRectangle excel.CellBox) bool {
	if rectangle.LeftColumn() == (otherRectangle.RightColumn()+1) || otherRectangle.LeftColumn() == (rectangle.RightColumn()+1) {
		return rectangle.TopRow() == otherRectangle.TopRow() && rectangle.BottomRow() == otherRectangle.BottomRow()
	}
	if rectangle.TopRow() == (otherRectangle.BottomRow()+1) || otherRectangle.TopRow() == (rectangle.BottomRow()+1) {
		return rectangle.LeftColumn() == otherRectangle.LeftColumn() && rectangle.RightColumn() == otherRectangle.RightColumn()
	}
	return false
}

func findBestFreeBox(c WriteContext, boxes []AvailableFreeSpace) AvailableFreeSpace {
	bestScore := float64(-1)
	best := AvailableFreeSpace{}

	for _, box := range boxes {
		score := float64(0)

		// count num of free cells around this annotation
		numFreeCells := countFreeCellsAroundBox(c, box.box)
		perimeter := 2*(box.box.RightColumn()-box.box.LeftColumn()+3) +
			2*(box.box.BottomRow()-box.box.TopRow()+3)
		fractionOfFreeCells := float64(numFreeCells) / float64(perimeter)

		score += (fractionOfFreeCells * 10) * 10000
		score += box.preference * 10000
		score += box.width*100 + box.height

		if score > bestScore {
			bestScore = score
			best = box
		}
	}

	return best
}

func findContiguousFreeBoxes(wc WriteContext, box excel.CellBox, specs AvailableFreeSpaceSearchSpec) ([]AvailableFreeSpace, error) {
	// try at the top
	out := make([]AvailableFreeSpace, 0)

	attempt := func(from excel.Cell, with []FreeSpaceSearchStep, settings AvailableFreeSpace) error {
		top, err := findContiguousFreeBoxFromPoint(
			wc,
			from,
			AvailableFreeSpaceSearchSpec{
				steps:     with,
				minWidth:  specs.minWidth,
				maxWidth:  specs.maxWidth,
				minHeight: specs.minHeight,
				maxHeight: specs.maxHeight,
			},
		)
		if err != nil {
			return err
		}
		if top != nil {
			top.preference = settings.preference
			top.arrowAnnotation = settings.arrowAnnotation
			top.arrowAtEnd = settings.arrowAtEnd
			out = append(out, *top)
		}
		return nil
	}

	if err := attempt(box.TopRight().AtRight(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Right: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Bottom: 1}}},
	}, AvailableFreeSpace{preference: 10, arrowAnnotation: "◀"}); err != nil {
		return out, err
	}
	if err := attempt(box.TopLeft().AtTop(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Right: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Top: 1}}},
	}, AvailableFreeSpace{preference: 5, arrowAnnotation: "▼"}); err != nil {
		return out, err
	}
	if err := attempt(box.TopRight().AtTop(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Left: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Top: 1}}},
	}, AvailableFreeSpace{preference: 1, arrowAnnotation: "▼", arrowAtEnd: true}); err != nil {
		return out, err
	}
	if err := attempt(box.BottomLeft().AtBottom(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Right: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Bottom: 1}}},
	}, AvailableFreeSpace{preference: 0, arrowAnnotation: "▲"}); err != nil {
		return out, err
	}
	if err := attempt(box.BottomRight().AtRight(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Right: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Top: 1}}},
	}, AvailableFreeSpace{preference: 3, arrowAnnotation: "◀"}); err != nil {
		return out, err
	}
	if err := attempt(box.BottomRight().AtBottom(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Left: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Bottom: 1}}},
	}, AvailableFreeSpace{preference: 0, arrowAnnotation: "▲", arrowAtEnd: true}); err != nil {
		return out, err
	}
	if err := attempt(box.TopLeft().AtLeft(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Left: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Bottom: 1}}},
	}, AvailableFreeSpace{preference: 1, arrowAnnotation: "▶", arrowAtEnd: true}); err != nil {
		return out, err
	}
	if err := attempt(box.BottomLeft().AtLeft(1), []FreeSpaceSearchStep{
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Left: 1}}},
		{Mutators: []FreeSpaceSearchStepBoxMutator{{Bottom: 1}}},
	}, AvailableFreeSpace{preference: 0, arrowAnnotation: "▶", arrowAtEnd: true}); err != nil {
		return out, err
	}

	return out, nil
}

func countFreeCellsAroundBox(c WriteContext, box excel.CellBox) uint {
	numFreeCells := uint(0)
	row := int(box.TopRow()) - 1
	if row >= 1 {
		for col := int(box.LeftColumn()) - 1; col <= int(box.RightColumn())+1; col++ {
			if col < 1 {
				continue
			}
			if isCellFree(c, excel.NewCell(box.TopLeft().SheetName(), uint(col), uint(row))) {
				numFreeCells++
			}
		}
	}
	row = int(box.BottomRow()) + 1
	{
		for col := int(box.LeftColumn()) - 1; col <= int(box.RightColumn())+1; col++ {
			if col < 1 {
				continue
			}
			if isCellFree(c, excel.NewCell(box.TopLeft().SheetName(), uint(col), uint(row))) {
				numFreeCells++
			}
		}
	}

	col := int(box.LeftColumn()) - 1
	if col >= 1 {
		for row := int(box.TopRow()) - 1; row <= int(box.BottomRow())+1; row++ {
			if row < 1 {
				continue
			}
			if isCellFree(c, excel.NewCell(box.TopLeft().SheetName(), uint(col), uint(row))) {
				numFreeCells++
			}
		}
	}
	col = int(box.LeftColumn()) + 1
	{
		for row := int(box.TopRow()) - 1; row <= int(box.BottomRow())+1; row++ {
			if row < 1 {
				continue
			}
			if isCellFree(c, excel.NewCell(box.TopLeft().SheetName(), uint(col), uint(row))) {
				numFreeCells++
			}
		}
	}

	return numFreeCells
}
