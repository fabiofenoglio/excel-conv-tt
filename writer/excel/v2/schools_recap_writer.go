package excel

import (
	"fmt"
	"regexp"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/excel"
)

type rewriteRule struct {
	find    *regexp.Regexp
	replace string
}

var schoolNameRewriteRules = []rewriteRule{
	{find: regexp.MustCompile(`(?mi)^scuola\s*\:*\s*`), replace: ""},
	{find: regexp.MustCompile(`(?mi)(^|\s)((?:ISTITUTO\s+COMPRENSIVO|I\.C\.|(?:IC))[\s\-]+)+`), replace: "${1}I.C. "},
	{find: regexp.MustCompile(`(?mi)istruzione\s+secondaria\s*superiore`), replace: "I.S.S."},
	{find: regexp.MustCompile(`(?mi)(?:^|\s)(primaria)(?:$|\s)`), replace: " EL "},
	{find: regexp.MustCompile(`(?mi)(?:^|\s)(secondaria\s+(I°?|primo)\s+grado)(?:$|\s)`), replace: " SM "},
	{find: regexp.MustCompile(`(?mi)(?:^|\s)(secondaria\s+(II°?|secondo)\s+grado)(?:$|\s)`), replace: " SUP "},
	{find: regexp.MustCompile(`[\r\n]+`), replace: " "},
	{find: regexp.MustCompile(`[\s\t]+`), replace: " "},
	{find: regexp.MustCompile(`^\s*`), replace: ""},
	{find: regexp.MustCompile(`\s*$`), replace: ""},
}

func rewriteSchoolNameWithRules(raw string) string {
	for _, r := range schoolNameRewriteRules {
		raw = r.find.ReplaceAllString(raw, r.replace)
	}
	return raw
}

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

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(6).Code()); err != nil {
			return err
		}

		if err := f.MergeCell(cursor.SheetName(), cursor.AtRight(7).Code(), cursor.AtRight(8).Code()); err != nil {
			return err
		}

		if lastSchoolCodeWrote != "" && groupRef.SchoolCode == lastSchoolCodeWrote {
			if err := f.MergeCell(cursor.SheetName(), cursor.AtTop(1).Code(), cursor.Code()); err != nil {
				return err
			}
			if err := f.MergeCell(cursor.SheetName(), cursor.AtRight(7).AtTop(1).Code(), cursor.AtRight(7).Code()); err != nil {
				return err
			}
		} else {
			lastSchoolCodeWrote = groupRef.SchoolCode

			toWrite = rewriteSchoolNameWithRules(schoolRef.Name)

			if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
				return err
			}

			toWrite = rewriteSchoolNameWithRules(schoolRef.Type)

			if err := f.SetCellValue(cursor.SheetName(), cursor.AtRight(7).Code(), toWrite); err != nil {
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

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), groupRef.Composition.NumTotal()); err != nil {
			return err
		}

		if groupRef.Composition.NumFree > 0 {
			toWrite = fmt.Sprintf("Paganti: %d", groupRef.Composition.NumPaying)
			if groupRef.Composition.NumAccompanying > 0 {
				toWrite += fmt.Sprintf("\nAccompagnatori: %d", groupRef.Composition.NumAccompanying)
			}
			if groupRef.Composition.NumFree > 0 {
				toWrite += fmt.Sprintf("\nGRATUITI: %d", groupRef.Composition.NumFree)
			}
			if err := addCommentToCell(f, cursor, toWrite); err != nil {
				return err
			}
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
