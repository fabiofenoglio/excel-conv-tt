package services

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type WriteContext struct {
	minHour     int
	maxHour     int
	minutesStep int
	allData     Parsed
	outputFile  *excelize.File
}

func Run(input string) error {
	log := GetLogger()

	// parse the specified input file

	parsed, err := Parse(input)
	if err != nil {
		return fmt.Errorf("error reading from input file: %w", err)
	}
	log.Infof("found %d rows", len(parsed.Rows))

	// analyze the data

	// compute the max time range to be shown between all days

	minHourToShow, maxHourToShow := GetMaxHoursRangeToDisplay(parsed.Rows)
	log.Infof("will show the range %d to %d (inclusive)", minHourToShow, maxHourToShow)

	// group by start date, ordering by start time ASC

	groupedByStartDate := GroupByStartDay(parsed.Rows)
	log.Infof("found activities spanning %d different days", len(groupedByStartDate))

	groupedByRoom := GroupByRoom(parsed.Rows)
	log.Infof("found activities occupying %d different rooms", len(groupedByRoom))

	groupedByOperator := GroupByOperator(parsed.Rows)
	log.Infof("found activities occupying %d different operators", len(groupedByOperator))

	f := excelize.NewFile()
	RegisterStyles(f)

	wc := WriteContext{
		minHour:     minHourToShow,
		maxHour:     maxHourToShow,
		minutesStep: 15,
		allData:     parsed,
		outputFile:  f,
	}

	sheetName := fmt.Sprintf(
		"settimana %s-%s",
		groupedByStartDate[0].Rows[0].StartAt.Format("0201"),
		groupedByStartDate[len(groupedByStartDate)-1].Rows[0].StartAt.Format("0201"),
	)

	startingCell := Cell{C: 2, R: 2, SheetName: sheetName}
	// Create a new sheet.

	f.SetSheetName("Sheet1", startingCell.SheetName)

	// Set active sheet of the workbook.
	f.SetActiveSheet(f.GetSheetIndex(startingCell.SheetName))

	if err := f.SetRowHeight(startingCell.SheetName, 1, 10); err != nil {
		return err
	}
	if err := f.SetRowHeight(startingCell.SheetName, 2, 30); err != nil {
		return err
	}
	if err := f.SetRowHeight(startingCell.SheetName, 3, 20); err != nil {
		return err
	}

	if err := f.SetColWidth(startingCell.SheetName, "A", "A", 3); err != nil {
		return err
	}

	// MEGA OUTPUT WRITE
	cell := startingCell
	for _, groupByDay := range groupedByStartDate {
		cell.R = startingCell.R
		finalCell, err := writeDay(wc, groupByDay, cell)
		if err != nil {
			return fmt.Errorf("error writing row: %w", err)
		}
		cell.C = finalCell.C + 2
		cell.R = finalCell.R
	}

	for i := 0; i < 20; i++ {
		r := cell.Copy()
		r.R += uint(i + 1)
		if err := f.SetRowHeight(r.SheetName, int(r.R), 3); err != nil {
			return err
		}
	}

	// Save spreadsheet by the given path.
	outPath := filepath.Dir(input)
	inputName := filepath.Base(input)
	inputExt := filepath.Ext(input)
	if err := f.SaveAs(outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-convertito" + inputExt); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}

func writeDay(c WriteContext, day GroupedRows, startCell Cell) (Cell, error) {
	cursor := startCell
	f := c.outputFile

	if err := f.SetCellValue(cursor.SheetName, cursor.Code(),
		fmt.Sprintf("%d", day.Rows[0].StartAt.Year())); err != nil {
		return cursor, err
	}
	cursor.C++
	cursor.R++

	maxC := cursor.C
	maxR := cursor.R

	groupedByRoom := GroupByRoom(day.Rows)

	minGroupWidth := uint(12)
	minGroupWidthPerSlot := uint(4)

	roomColumnNumbers := make(map[string]uint)

	for _, group := range groupedByRoom {
		room := c.allData.RoomsMap[group.Key]
		numActs := uint(len(group.Rows))

		toWrite := room.Name
		if room.Code == "" {
			toWrite = "???"
		}

		if err := f.SetCellValue(cursor.SheetName, cursor.Code(), toWrite); err != nil {
			return cursor, err
		}

		roomColumnNumbers[room.Code] = cursor.C

		desiredWidth := minGroupWidth

		if numActs > 1 {
			mergeUpTo := cursor.Copy()
			mergeUpTo.C += (numActs - 1)

			if err := f.MergeCell(cursor.SheetName, cursor.Code(), mergeUpTo.Code()); err != nil {
				return cursor, err
			}

			altWidth := numActs * minGroupWidthPerSlot
			if altWidth > desiredWidth {
				desiredWidth = altWidth
			}

			if err := f.SetColWidth(
				cursor.SheetName,
				cursor.ColumnName(),
				mergeUpTo.ColumnName(),
				float64(desiredWidth)/float64(numActs),
			); err != nil {
				return cursor, err
			}

		} else {
			if err := f.SetColWidth(cursor.SheetName, cursor.ColumnName(), cursor.ColumnName(), float64(desiredWidth)); err != nil {
				return cursor, err
			}
		}

		cursor.C += numActs
		maxC = cursor.C

		// add 1 additional column to leave some space
		if err := f.SetColWidth(cursor.SheetName, cursor.ColumnName(), cursor.ColumnName(), 2); err != nil {
			return cursor, err
		}
		cursor.C++
		maxC++
	}

	// reset cursor to box start
	cursor.C = startCell.C
	cursor.R = startCell.R + 2

	minHourToShowForThisDay, maxHourToShowForThisDay := GetMaxHoursRangeToDisplay(day.Rows)

	currentTime := day.Rows[0].StartAt.Add(time.Minute * -30)
	if day.Rows[0].StartAt.Hour() < 2 {
		currentTime = day.Rows[0].StartAt
	}

	currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		c.minHour, 0, 0, 0, currentTime.Location())
	rolledOver := false
	firstTime := currentTime
	breaking := false
	margin := 0 // add if you want to show some empty time rows after the last one

	boxCellHeight := float64(20)

	for {
		effectiveHour := currentTime.Hour()
		if rolledOver {
			effectiveHour += 24
		}

		show := true
		if minHourToShowForThisDay > c.minHour && effectiveHour < minHourToShowForThisDay {
			show = false
		}
		if maxHourToShowForThisDay < c.maxHour && effectiveHour > maxHourToShowForThisDay {
			show = false
		}

		toWrite := "--"
		if show {
			toWrite = currentTime.Format(layoutTimeOnlyInReadableFormat)
		}

		if err := f.SetCellValue(cursor.SheetName, cursor.Code(), toWrite); err != nil {
			return cursor, err
		}
		if err := f.SetRowHeight(cursor.SheetName, int(cursor.R), boxCellHeight); err != nil {
			return cursor, err
		}
		cursor.R++
		maxR = cursor.R

		currentTime = currentTime.Add(time.Duration(c.minutesStep) * time.Minute)

		if currentTime.Day() != firstTime.Day() {
			rolledOver = true
		}

		if !breaking {
			if effectiveHour > c.maxHour {
				breaking = true
			}
		} else {
			if margin > 0 {
				margin--
			} else {
				break
			}
		}
	}

	// finished writing the times on the left

	// draw the box for the whole day
	boxStart := startCell.Copy()
	boxEnd := startCell.Copy()
	boxEnd.C = maxC - 1
	boxEnd.R = maxR - 1
	if err := f.SetCellStyle(startCell.SheetName, boxStart.Code(), boxEnd.Code(), dayBoxStyle.StyleID); err != nil {
		return cursor, err
	}

	// draw the box for each room
	for _, group := range groupedByRoom {
		columnsForRoomStartAt := roomColumnNumbers[group.Key]

		boxStart := startCell.Copy()
		boxStart.C = columnsForRoomStartAt
		boxStart.R = startCell.R + 2

		boxEnd := boxStart.Copy()
		boxEnd.C += uint(len(group.Rows) - 1)
		boxEnd.R = maxR - 1

		if err := f.SetCellStyle(startCell.SheetName, boxStart.Code(), boxEnd.Code(), dayRoomBoxStyle.StyleID); err != nil {
			return cursor, err
		}

		boxHeaderStart := boxStart.Copy()
		boxHeaderStart.R--
		boxHeaderEnd := boxHeaderStart.Copy()
		boxHeaderEnd.C = boxEnd.C

		roomInfo := GetEntryForRoom(group.Key)
		if err := f.SetCellStyle(cursor.SheetName, boxHeaderStart.Code(), boxHeaderEnd.Code(), roomInfo.Styles.Common.StyleID); err != nil {
			return cursor, err
		}

	}

	// now place the activities
	cursor.R = startCell.R
	cursor.C = startCell.C

	for _, group := range groupedByRoom {
		// TODO add explicit sorting here?
		for actIndex, act := range group.Rows {
			var operator Operator
			if act.Operator.Code != "" {
				operator, _ = c.allData.OperatorsMap[act.Operator.Code]
			}
			if operator.Code == "" {
				operator.Name = "???"
			}

			columnsForRoomStartAt := roomColumnNumbers[act.Room.Code]
			cursor.C = columnsForRoomStartAt + uint(actIndex)
			cursor.R = startCell.R + 2

			inRange := false
			exited := false
			timeCursor := firstTime
			i := 0

			actStartCell := cursor.Copy()
			actEndCell := cursor.Copy()

			for {
				if !inRange {
					if !exited && !timeCursor.Before(act.StartAt) {
						inRange = true
						actStartCell.R = startCell.R + 2 + uint(i)
					}
				} else {
					if !timeCursor.Before(act.EndAt) {
						inRange = false
						exited = true
						actEndCell.R = startCell.R + 2 + uint(i-1)
					}
				}

				if exited {
					break
				}

				timeCursor = timeCursor.Add(time.Minute * time.Duration(c.minutesStep))
				i++
			}

			/*
				// if merging the cells is desired:

				if actStartCell.R != actEndCell.R {
					if err := f.MergeCell(actStartCell.SheetName, actStartCell.Code(), actEndCell.Code()); err != nil {
						return cursor, err
					}
				}
				if err := f.SetCellValue(actStartCell.SheetName, actStartCell.Code(), operator.Name); err != nil {
					return cursor, err
				}
			*/
			// if merging the cells is NOT desired:
			for r := actStartCell.R; r <= actEndCell.R; r++ {
				c := actStartCell.Copy()
				c.R = r
				if err := f.SetCellValue(c.SheetName, c.Code(), operator.Name); err != nil {
					return cursor, err
				}
			}

			// apply style depending from operator
			operatorInfo := GetEntryForOperator(act.Operator.Code)
			if err := f.SetCellStyle(actStartCell.SheetName, actStartCell.Code(), actEndCell.Code(), operatorInfo.Styles.Common.StyleID); err != nil {
				return cursor, err
			}

			// write some comment with additional info
			cellComment := buildComment(act)
			if cellComment != "" {
				comment(f, actStartCell, strings.TrimSpace(cellComment))
			}

			// write activity name under everything
			if act.Raw.Attivita != "" {
				actDescCell := actStartCell.Copy()
				actDescCell.R = maxR + 1

				actDescFinalCell := actDescCell.Copy()
				actDescFinalCell.R += 10

				if err := f.MergeCell(actDescCell.SheetName, actDescCell.Code(), actDescFinalCell.Code()); err != nil {
					return cursor, err
				}
				if err := f.SetCellValue(actDescCell.SheetName, actDescCell.Code(), act.Raw.Attivita); err != nil {
					return cursor, err
				}
				if err := f.SetCellStyle(actDescCell.SheetName, actDescCell.Code(), actDescCell.Code(), verticalTextStyle.StyleID); err != nil {
					return cursor, err
				}

			}
		}
	}

	// write the day name in the header
	cursor = startCell.Copy()
	cursor.C++
	// check if we can merge some cells
	if maxC-startCell.C >= 11 {
		toCell := cursor.Copy()
		toCell.C += 8

		if err := f.MergeCell(cursor.SheetName, cursor.Code(), toCell.Code()); err != nil {
			return cursor, err
		}
	}

	if err := f.SetCellValue(cursor.SheetName, cursor.Code(),
		day.Rows[0].StartAt.Format("Mon 2 January")); err != nil {
		return cursor, err
	}
	if err := f.SetCellStyle(cursor.SheetName, cursor.Code(), cursor.Code(), dayHeaderStyle.StyleID); err != nil {
		return cursor, err
	}

	// compute the final cursor position
	finalCursor := startCell.Copy()
	finalCursor.C = maxC - 1
	finalCursor.R = maxR - 1

	return finalCursor, nil
}

func buildComment(act ParsedRow) string {
	cellComment := ``
	if act.Room.Code == "" {
		cellComment += "ATTENZIONE: NESSUNA AULA O RISORSA ASSEGNATA\n\n"
	}
	if act.Operator.Code == "" {
		cellComment += "ATTENZIONE: NESSUN EDUCATORE ASSEGNATO\n\n"
	}

	if act.Raw.LinguaAttivita != "" && strings.ToLower(act.Raw.LinguaAttivita) != "it" {
		cellComment += "ATTENZIONE: ATTIVITA PREVISTA IN LINGUA: " + act.Raw.LinguaAttivita + "\n\n"
	}

	if act.Raw.Attivita != "" {
		cellComment += act.Raw.Attivita

		if act.Raw.Tipologia != "" {
			cellComment += " (" + act.Raw.Tipologia + ")"
		}
		if act.Raw.Codice != "" {
			cellComment += " [" + act.Raw.Codice + "]"
		}

		cellComment += "\n\n"

	} else if act.Raw.Tipologia != "" {
		cellComment += "Tipologia: " + act.Raw.Tipologia + "\n"
	}

	if act.Raw.Orario != "" {
		cellComment += "Orario: " + act.Raw.Orario + "\n"
	}
	if act.Raw.Educatore != "" {
		cellComment += "Educatore: " + act.Raw.Educatore + "\n"
	}
	if act.Raw.Aula != "" {
		cellComment += "Aula: " + act.Raw.Aula + "\n"
	}

	if act.Raw.NotaOperatore != "" {
		cellComment += "Nota operatore: " + act.Raw.NotaOperatore + "\n"
	}
	if act.Raw.NotaPrenotazione != "" {
		cellComment += "Nota prenotazione: " + act.Raw.NotaPrenotazione + "\n"
	}

	if act.Raw.Classe != "" || act.Raw.Sezione != "" {
		if act.Raw.Classe != "" {
			cellComment += "Classe: " + act.Raw.Classe + " "
			if act.Raw.Sezione != "" {
				cellComment += act.Raw.Sezione
			}
		} else {
			cellComment += "Sezione: " + act.Raw.NomeScuola
		}
		cellComment += "\n"
	}

	if act.NumPaganti > 0 || act.NumAccompagnatori > 0 || act.NumGratuiti > 0 {
		c := ""
		tot := 0
		entries := 0
		if act.NumPaganti > 0 {
			c += fmt.Sprintf("%d paganti, ", act.NumPaganti)
			tot += act.NumPaganti
			entries++
		}
		if act.NumGratuiti > 0 {
			c += fmt.Sprintf("%d gratuiti, ", act.NumGratuiti)
			tot += act.NumGratuiti
			entries++
		}
		if act.NumAccompagnatori > 0 {
			c += fmt.Sprintf("%d accompagnatori, ", act.NumAccompagnatori)
			tot += act.NumAccompagnatori
			entries++
		}

		if entries > 1 {
			c = strings.TrimSuffix(c, ", ") + fmt.Sprintf(" (%d totali)", tot)
		}

		cellComment += strings.TrimSuffix(c, ", ") + "\n"
	}

	if act.Raw.NomeScuola != "" || act.Raw.TipologiaScuola != "" {
		if act.Raw.TipologiaScuola != "" {
			cellComment += act.Raw.TipologiaScuola + " "
		}
		if act.Raw.NomeScuola != "" {
			cellComment += act.Raw.NomeScuola
		}
		cellComment += "\n"
	}

	if act.Raw.Bus != "" {
		cellComment += "Bus: " + act.Raw.Bus + "\n"
	}
	if act.Raw.Acconti != "" && act.Raw.Acconti != "-" {
		cellComment += "Acconti: " + act.Raw.Acconti + "\n"
	}
	if act.Raw.StatoAcconti != "" && act.Raw.StatoAcconti != "-" {
		cellComment += "Acconti: " + act.Raw.StatoAcconti + "\n"
	}

	return strings.TrimSpace(cellComment)
}

func comment(f *excelize.File, cell Cell, content string) error {
	type serializable struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	v := serializable{
		Author: "excel converter",
		Text:   content,
	}
	serialized, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return f.AddComment(cell.SheetName, cell.Code(), string(serialized))
}
