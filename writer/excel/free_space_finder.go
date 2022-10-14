package excel

import (
	"errors"

	"github.com/fabiofenoglio/excelconv/excel"
)

type AvailableFreeSpace struct {
	box             excel.CellBox
	width           float64
	height          float64
	preference      float64
	arrowAnnotation string
	arrowAtEnd      bool
}

type AvailableFreeSpaceSearchSpec struct {
	steps                []FreeSpaceSearchStep
	minWidth, maxWidth   float64
	minHeight, maxHeight float64
}

type FreeSpaceSearchStep struct {
	Mutators []FreeSpaceSearchStepBoxMutator
}

type FreeSpaceSearchStepBoxMutator struct {
	Left, Right, Top, Bottom uint
}

func findContiguousFreeBoxFromPoint(wc WriteContext, from excel.Cell, specs AvailableFreeSpaceSearchSpec) (*AvailableFreeSpace, error) {
	if len(specs.steps) == 0 {
		return nil, errors.New("mutation steps are required")
	}

	startCell, endCell := from.Copy(), from.Copy()

	var validBox excel.CellBox = nil

	for _, step := range specs.steps {
		if len(step.Mutators) == 0 {
			return nil, errors.New("step mutators are required")
		}

		for {
			attemptingBox := excel.NewCellBox(
				startCell.Copy(),
				endCell.Copy(),
			)
			if !isBoxFree(wc, attemptingBox) {
				break
			}

			attemptingBoxW, attemptingBoxH := measureBox(wc, attemptingBox)
			aboveMinMeasures := attemptingBoxW >= specs.minWidth && attemptingBoxH >= specs.minHeight
			underMaxMeasures := (specs.maxWidth <= 0 || attemptingBoxW <= specs.maxWidth) &&
				(specs.maxHeight <= 0 || attemptingBoxH <= specs.maxHeight)

			if aboveMinMeasures && underMaxMeasures {
				validBox = attemptingBox
			}

			if !underMaxMeasures {
				break
			}

			for _, mutator := range step.Mutators {
				if mutator.Top != 0 {
					startCell.MoveTop(mutator.Top)
				}
				if mutator.Bottom != 0 {
					endCell.MoveBottom(mutator.Bottom)
				}
				if mutator.Left != 0 {
					startCell.MoveLeft(mutator.Left)
				}
				if mutator.Right != 0 {
					endCell.MoveRight(mutator.Right)
				}
			}
		}
	}

	if validBox == nil {
		return nil, nil
	}

	width, height := measureBox(wc, validBox)
	return &AvailableFreeSpace{
		box:    validBox,
		width:  width,
		height: height,
	}, nil
}

func isBoxFree(wc WriteContext, box excel.CellBox) bool {
	sn := box.TopLeft().SheetName()
	r, c := box.TopRow(), box.LeftColumn()
	for ; r <= box.BottomRow(); r++ {
		for ; c <= box.RightColumn(); c++ {
			if !isCellFree(wc, excel.NewCell(sn, c, r)) {
				return false
			}
		}
	}
	return true
}

func isCellFree(c WriteContext, cell excel.Cell) bool {
	f := c.outputFile

	v, err := f.GetCellValue(cell.SheetName(), cell.Code())
	if err != nil {
		return false
	}
	if v != "" {
		return false
	}
	style, err := f.GetCellStyle(cell.SheetName(), cell.Code())
	if err != nil {
		return false
	}
	if style == 0 {
		return true
	}
	if c.styleRegister.DayBoxStyle().MatchesID(style) {
		return true
	}
	if c.styleRegister.DayRoomBoxStyle().MatchesID(style) {
		return true
	}
	if c.styleRegister.UnusedRoomStyle().MatchesID(style) {
		return true
	}

	return false
}

func measureBox(wc WriteContext, box excel.CellBox) (w, h float64) {
	sn := box.TopLeft().SheetName()
	f := wc.outputFile

	width := float64(0)
	height := float64(0)
	for cursor := box.TopLeft().Copy(); cursor.Column() <= box.RightColumn(); cursor.MoveRight(1) {
		w, _ := f.GetColWidth(sn, cursor.ColumnName())
		width += w
	}
	for cursor := box.TopLeft().Copy(); cursor.Row() <= box.BottomRow(); cursor.MoveBottom(1) {
		h, _ := f.GetRowHeight(sn, int(cursor.Row()))
		height += h
	}

	return width, height
}
