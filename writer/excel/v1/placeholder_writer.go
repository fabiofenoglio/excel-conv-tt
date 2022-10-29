package excel

import (
	"github.com/fabiofenoglio/excelconv/excel"
)

func writePlaceholdersForDay(c WriteContext, startCell excel.Cell) error {
	f := c.outputFile
	cursor := startCell.Copy()

	// write placeholder for info at the bottom
	rowsToWrite := []string{
		"Piano 0 / Accoglienza",
		"Bookshop / Cassa",
		"Cambio Stefano",
		"On / off museo",
		"Allest. / disallest.",
		"Assenti",
		"Appuntamenti / note",
		"Responsabile emergenza",
		"Addetto antincendio / impianti",
		"Addetto antincendio / primo soccorso",
		"Orari navetta dalle - alle",
	}

	cursor = cursor.AtBottom(1)

	for _, rowToWrite := range rowsToWrite {
		valueCursors := cursor.AtRight(10)
		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), valueCursors.AtLeft(1).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), rowToWrite); err != nil {
			return err
		}

		if err := f.SetCellValue(cursor.SheetName(), valueCursors.Code(), "?"); err != nil {
			return err
		}
		if err := f.MergeCell(cursor.SheetName(), valueCursors.Code(), valueCursors.AtRight(9).Code()); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), valueCursors.Code(), valueCursors.Code(),
			c.styleRegister.ToBeFilledStyle().Common.StyleID); err != nil {
			return err
		}

		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 25); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}
