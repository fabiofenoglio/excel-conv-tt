package excel

import (
	"fmt"
	"github.com/fabiofenoglio/excelconv/parser/v2"
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

func writeSchoolsForDay(c WriteContext, groups []aggregator2.VisitingGroupInDay, startCell excel.Cell, availableColumns uint) error {
	f := c.outputFile
	commonData := c.allData.CommonData
	cursor := startCell.Copy()
	rightColumn := startCell.Column() + availableColumns - 1

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

		// fetch detailed data for group
		groupRef := c.anagraphicsRef.VisitingGroups[schoolGroup.VisitingGroupCode]
		schoolRef := c.anagraphicsRef.Schools[groupRef.SchoolCode]
		schoolClassRef := c.anagraphicsRef.SchoolClasses[groupRef.SchoolClassCode]

		// write group display code in first cell
		toWrite := schoolGroup.DisplayCode
		displayCodePosition := cursor.Copy()
		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtBottom(1).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}
		cursor.MoveRight(1)

		// merge 7 cells for school name and 2 cells for school type
		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(6).AtBottom(1).Code()); err != nil {
			return err
		}
		if err := f.MergeCell(cursor.SheetName(), cursor.AtRight(7).Code(), cursor.AtRight(8).AtBottom(1).Code()); err != nil {
			return err
		}

		if lastSchoolCodeWrote != "" && groupRef.SchoolCode == lastSchoolCodeWrote {
			// same school of previous line, just merge with top row
			if err := f.MergeCell(cursor.SheetName(),
				cursor.AtTop(2).Code(),
				cursor.AtBottom(1).Code()); err != nil {
				return err
			}
			if err := f.MergeCell(cursor.SheetName(),
				cursor.AtRight(7).AtTop(2).Code(),
				cursor.AtRight(7).AtBottom(1).Code()); err != nil {
				return err
			}
		} else {
			// school is different than previous line, write a new one
			lastSchoolCodeWrote = groupRef.SchoolCode

			if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), rewriteSchoolNameWithRules(schoolRef.Name)); err != nil {
				return err
			}
			if err := f.SetCellValue(cursor.SheetName(), cursor.AtRight(7).Code(), rewriteSchoolNameWithRules(schoolRef.Type)); err != nil {
				return err
			}
		}

		cursor.MoveRight(9)

		// write class name merging 3 cells
		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).AtBottom(1).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), schoolClassRef.FullDescription()); err != nil {
			return err
		}

		cursor.MoveRight(3)

		// write class composition merging 3 cells
		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).AtBottom(1).Code()); err != nil {
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

		cursor.MoveRight(3)

		// merge all remaining columns for notes
		if err := f.MergeCell(cursor.SheetName(),
			cursor.Code(),
			cursor.AtColumn(rightColumn).Code()); err != nil {
			return err
		}
		if err := f.MergeCell(cursor.SheetName(),
			cursor.AtBottom(1).Code(),
			cursor.AtBottom(1).AtColumn(rightColumn).Code()); err != nil {
			return err
		}

		// write notes and contacts in the same cell
		toWrite = ""
		addNote := func(input string) {
			if input != "" && !strings.Contains(toWrite, input) {
				toWrite += "⚠️ " + input + "\n"
			}
		}
		addNote(groupRef.SpecialProjectNotes)
		addNote(groupRef.BookingNotes)
		addNote(groupRef.OperatorNotes)

		toWrite = strings.TrimSuffix(toWrite, "\n")
		didWriteNotes := false
		if toWrite != "" {
			didWriteNotes = true
			if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
				return err
			}
		}

		cursor.MoveBottom(1)

		toWrite = ""
		if groupRef.ClassTeacher != "" {
			toWrite += groupRef.ClassTeacher
		}
		if groupRef.ClassRefEmail != "" {
			toWrite += " - " + groupRef.ClassRefEmail
		}
		toWrite = strings.TrimSuffix(toWrite, "\n")

		writeContactIn := cursor.Copy()
		if !didWriteNotes {
			// merge both lines
			if err := f.MergeCell(cursor.SheetName(),
				cursor.AtTop(1).Code(),
				cursor.Code()); err != nil {
				return err
			}
			writeContactIn = writeContactIn.AtTop(1)
		}

		if err := f.SetCellValue(cursor.SheetName(), writeContactIn.Code(), toWrite); err != nil {
			return err
		}

		// set rows height
		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.AtTop(1).Row()), 25); err != nil {
			return err
		}
		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 15); err != nil {
			return err
		}

		cursor.MoveTop(1)

		if err := f.SetCellStyle(cursor.SheetName(),
			cursor.AtColumn(startCell.Column()).Code(),
			cursor.AtColumn(rightColumn).AtBottom(1).Code(),
			c.styleRegister.SchoolRecapStyle().SingleCell()); err != nil {
			return err
		}

		if didWriteNotes {
			if err := f.SetCellStyle(cursor.SheetName(),
				cursor.AtColumn(startCell.Column()).AtRight(16).Code(),
				cursor.AtColumn(rightColumn).Code(),
				c.styleRegister.SchoolRecapNotesStyle().SingleCell()); err != nil {
				return err
			}

			if err := f.SetCellStyle(cursor.SheetName(),
				cursor.AtColumn(startCell.Column()).AtRight(16).AtBottom(1).Code(),
				cursor.AtColumn(rightColumn).AtBottom(1).Code(),
				c.styleRegister.SchoolRecapContactStyle().SingleCell()); err != nil {
				return err
			}
		} else {
			if err := f.SetCellStyle(cursor.SheetName(),
				cursor.AtColumn(startCell.Column()).AtRight(16).Code(),
				cursor.AtColumn(rightColumn).AtBottom(1).Code(),
				c.styleRegister.SchoolRecapContactStyle().SingleCell()); err != nil {
				return err
			}
		}

		if len(groupRef.Highlights) > 0 {
			styleForHighlight := highlightsToStyle(c, groupRef.Highlights)
			if styleForHighlight != nil {
				merged := c.styleRegister.Merge(c.styleRegister.SchoolRecapStyle(), styleForHighlight)
				if err := applyStyleToBox(f, merged, excel.NewCellBox(displayCodePosition, displayCodePosition.AtBottom(1))); err != nil {
					return err
				}
			}
		}

		cursor.MoveBottom(2)
	}

	currentNum := len(groups) * 2
	for currentNum < commonData.MaxVisitingGroupsPerDay*2 {
		cursor.MoveBottom(1)
		currentNum++
	}

	return nil
}

func highlightsToStyle(c WriteContext, highlights []parser.HighlightReason) *RegisteredStyleV2 {
	if len(highlights) == 0 {
		return nil
	}

	for _, highlight := range highlights {
		if highlight == parser.HighlightSpecialProject {
			return c.styleRegister.HighlightForSpecialProjectStyle()
		}
		if highlight == parser.HighlightSpecialNotes {
			return c.styleRegister.HighlightForSpecialNotesStyle()
		}
	}

	return c.styleRegister.SchoolRecapNotesStyle()
}
