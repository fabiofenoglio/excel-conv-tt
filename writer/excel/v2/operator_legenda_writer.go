package excel

import (
	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/fabiofenoglio/excelconv/parser/v2"
)

func writeOperatorsLegenda(c WriteContext, startCell excel.Cell, operators map[string]parser.Operator) error {
	// may optionally copy rows and sort here

	f := c.outputFile
	cursor := startCell.Copy().AtBottom(2)

	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(3).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "LEGENDA COLORI / EDUCATORI"); err != nil {
		return err
	}

	cursor.MoveBottom(1)

	for _, op := range operators {
		style := c.styleRegister.OperatorStyle(op.BackgroundColor)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(3).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), op.Name); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(),
			style.Common.StyleID); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}
