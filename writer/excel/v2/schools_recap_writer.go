package excel

import (
	"fmt"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/excel"
)

func writeSchoolsForDay(c WriteContext, groups []aggregator2.VisitingGroupInDay, startCell excel.Cell) error {
	f := c.outputFile
	cursor := startCell.Copy()

	for _, schoolGroup := range groups {
		cursor.MoveColumn(startCell.Column())

		groupRef := c.anagraphicsRef.VisitingGroups[schoolGroup.VisitingGroupCode]
		schoolRef := c.anagraphicsRef.Schools[groupRef.SchoolCode]
		schoolClassRef := c.anagraphicsRef.SchoolClasses[groupRef.SchoolClassCode]

		toWrite := schoolGroup.DisplayCode
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}
		cursor.MoveRight(1)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(8).Code()); err != nil {
			return err
		}

		toWrite = schoolRef.FullDescription()
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveRight(9)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
			return err
		}
		toWrite = schoolClassRef.FullDescription()
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveRight(3)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
			return err
		}

		toWrite = fmt.Sprintf("%d (%d",
			groupRef.Composition.NumTotal(),
			groupRef.Composition.NumPaying)
		if groupRef.Composition.NumAccompanying > 0 {
			toWrite += fmt.Sprintf(", %d acc.", groupRef.Composition.NumAccompanying)
		}
		if groupRef.Composition.NumFree > 0 {
			toWrite += fmt.Sprintf("+ %d GRAT.", groupRef.Composition.NumFree)
		}
		toWrite += ")"

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 25); err != nil {
			return err
		}

		cursor.MoveRight(3)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(5).Code()); err != nil {
			return err
		}

		toWrite = ""
		if groupRef.BookingNotes != "" {
			toWrite += groupRef.BookingNotes
		}
		if groupRef.OperatorNotes != "" {
			toWrite += "\n" + groupRef.OperatorNotes
		}
		toWrite = strings.TrimPrefix(toWrite, "\n")
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(),
			c.styleRegister.SchoolRecapStyle().Common.StyleID); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}
