package excel

import (
	"fmt"

	"github.com/fabiofenoglio/excelconv/aggregator/byactivity"
	"github.com/fabiofenoglio/excelconv/excel"
)

func writeSchoolsForDay(c WriteContext, groups []byactivity.ActivityGroup, startCell excel.Cell) error {
	f := c.outputFile
	cursor := startCell.Copy()

	for _, schoolGroup := range groups {
		cursor.MoveColumn(startCell.Column())

		toWrite := fmt.Sprintf("%s", schoolGroup.Code)
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}
		cursor.MoveRight(1)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(10).Code()); err != nil {
			return err
		}

		toWrite = schoolGroup.School.FullDescription()
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveRight(11)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
			return err
		}
		toWrite = schoolGroup.SchoolClass.FullDescription()
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveRight(3)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(4).Code()); err != nil {
			return err
		}

		toWrite = fmt.Sprintf("%d (%d",
			schoolGroup.Composition.Total,
			schoolGroup.Composition.NumPaying)
		if schoolGroup.Composition.NumAccompanying > 0 {
			toWrite += fmt.Sprintf(", %d acc.", schoolGroup.Composition.NumAccompanying)
		}
		if schoolGroup.Composition.NumFree > 0 {
			toWrite += fmt.Sprintf("+ %d GRAT.", schoolGroup.Composition.NumFree)
		}
		toWrite += ")"

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 25); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}
