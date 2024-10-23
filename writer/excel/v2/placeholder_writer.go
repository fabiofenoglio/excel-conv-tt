package excel

import (
	"strings"

	"github.com/fabiofenoglio/excelconv/excel"
)

func writePlaceholdersForDay(c WriteContext, startCell excel.Cell, availableColumns uint) error {
	f := c.outputFile
	cursor := startCell.Copy()

	type placeholder struct {
		emoji string
		text  string
	}

	// write placeholder for info at the bottom
	rowsToWrite := []placeholder{
		{emoji: "🛂", text: "Piano 0 / Accoglienza"},
		{emoji: "💶", text: "Bookshop / Cassa"},
		{emoji: "🔀", text: "Cambio Stefano"},
		{emoji: "🔌", text: "On / off museo"},
		{emoji: "🛠", text: "Allest. / disallest."},
		{emoji: "🚷", text: "Assenti"},
		{emoji: "📝", text: "Appuntamenti / note"},
		{emoji: "🚨", text: "Responsabile emergenza / antincendio"},
		{emoji: "🧯", text: "Addetto antincendio / impianti"},
		{emoji: "⛑", text: "Primo soccorso"},
		{emoji: "🚌", text: "Orari navetta dalle - alle"},
	}

	cursor = cursor.AtBottom(1)

	lastColumn := startCell.Column() + availableColumns - 1

	for _, rowToWrite := range rowsToWrite {
		valueCursors := cursor.AtRight(11)

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), rowToWrite.emoji); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(),
			c.styleRegister.Get(bottomPlaceholderEmojiStyle).SingleCell()); err != nil {
			return err
		}

		if err := f.MergeCell(cursor.SheetName(), cursor.AtRight(1).Code(), valueCursors.AtLeft(1).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.AtRight(1).Code(), strings.ToUpper(rowToWrite.text)); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), cursor.AtRight(1).Code(), valueCursors.AtLeft(1).Code(),
			c.styleRegister.Get(bottomPlaceholderTextStyle).SingleCell()); err != nil {
			return err
		}

		if err := f.SetCellValue(cursor.SheetName(), valueCursors.Code(), ""); err != nil {
			return err
		}
		if err := f.MergeCell(cursor.SheetName(), valueCursors.Code(), valueCursors.AtColumn(lastColumn).Code()); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), valueCursors.Code(), valueCursors.AtColumn(lastColumn).Code(),
			c.styleRegister.ToBeFilledStyle().SingleCell()); err != nil {
			return err
		}

		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 25); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}
