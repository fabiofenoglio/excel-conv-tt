package excel

import (
	"fmt"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/excel"
)

func writeSchoolsForDay(c WriteContext, groups []aggregator2.VisitingGroupInDay, startCell excel.Cell) error {
	f := c.outputFile
	commonData := c.allData.CommonData
	cursor := startCell.Copy()
	rightColumn := startCell.Column() + 21

	// write header
	cursor.MoveColumn(startCell.Column())

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "#"); err != nil {
		return err
	}
	cursor.MoveRight(1)

	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(8).Code()); err != nil {
		return err
	}

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "SCUOLA"); err != nil {
		return err
	}

	cursor.MoveRight(9)

	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "CLASSE"); err != nil {
		return err
	}

	cursor.MoveRight(3)

	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
		return err
	}

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "NUM."); err != nil {
		return err
	}

	if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 35); err != nil {
		return err
	}

	cursor.MoveRight(3)

	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(6).Code()); err != nil {
		return err
	}

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "NOTE E REFERENTI"); err != nil {
		return err
	}

	if err := f.SetCellStyle(cursor.SheetName(), cursor.AtColumn(startCell.Column()).Code(), cursor.AtColumn(rightColumn).Code(),
		c.styleRegister.Get(schoolRecapHeaderStyle).SingleCell()); err != nil {
		return err
	}

	cursor.MoveBottom(1)

	// write group by group
	lastSchoolCodeWrote := ""

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

		if lastSchoolCodeWrote != "" && groupRef.SchoolCode == lastSchoolCodeWrote {
			if err := f.MergeCell(cursor.SheetName(), cursor.AtTop(1).Code(), cursor.Code()); err != nil {
				return err
			}
		} else {
			lastSchoolCodeWrote = groupRef.SchoolCode

			toWrite = schoolRef.FullDescription()
			if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
				return err
			}
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

		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 35); err != nil {
			return err
		}

		cursor.MoveRight(3)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(6).Code()); err != nil {
			return err
		}

		toWrite = ""
		if groupRef.BookingNotes != "" {
			toWrite += "⚠️ " + groupRef.BookingNotes + "\n"
		}
		if groupRef.OperatorNotes != "" {
			toWrite += "⚠️ " + groupRef.OperatorNotes + "\n"
		}

		teacherLine := ""
		if groupRef.ClassTeacher != "" {
			teacherLine += groupRef.ClassTeacher
		}
		if groupRef.ClassRefEmail != "" {
			teacherLine += " - " + groupRef.ClassRefEmail
		}
		if teacherLine != "" {
			toWrite += strings.TrimPrefix(teacherLine, " - ") + "\n"
		}

		toWrite = strings.TrimSuffix(toWrite, "\n")

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	currentNum := len(groups)
	for currentNum < commonData.MaxVisitingGroupsPerDay {
		cursor.MoveBottom(1)
		currentNum++
	}

	if err := f.SetCellStyle(cursor.SheetName(), startCell.AtBottom(1).Code(), cursor.AtColumn(rightColumn).Code(),
		c.styleRegister.SchoolRecapStyle().SingleCell()); err != nil {
		return err
	}

	return nil
}
