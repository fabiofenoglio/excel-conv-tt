package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/database"
	"github.com/fabiofenoglio/excelconv/logger"
	"github.com/fabiofenoglio/excelconv/model"
	"github.com/fabiofenoglio/excelconv/parser"
)

type WriteContext struct {
	minHour     int
	maxHour     int
	minutesStep int
	allData     model.ParsedData
	outputFile  *excelize.File
}

func Run(args model.Args) error {
	input := args.PositionalArgs.InputFile
	if input == "" {
		return errors.New("missing input file")
	}
	log := logger.GetLogger()
	if args.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	// parse the specified input file

	parsed, err := parser.Parse(input)
	if err != nil {
		return fmt.Errorf("error reading from input file: %w", err)
	}
	log.Infof("found %d rows", len(parsed.Rows))

	// compute the max time range to be shown between all days

	minHourToShow, maxHourToShow := GetMaxHoursRangeToDisplay(parsed.Rows)
	log.Infof("will show the range %d to %d (inclusive)", minHourToShow, maxHourToShow)

	// group by start date, ordering by start time ASC

	groupedByStartDate := GroupByStartDay(parsed.Rows)
	log.Infof("found activities spanning %d different days", len(groupedByStartDate))

	f := excelize.NewFile()
	database.RegisterStyles(f)

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

	startingCell := model.NewCell(sheetName, 2, 1)

	// Rename the default sheet and set it as active
	f.SetSheetName("Sheet1", startingCell.SheetName())
	f.SetActiveSheet(f.GetSheetIndex(startingCell.SheetName()))

	if err := f.SetRowHeight(startingCell.SheetName(), 1, 20); err != nil {
		return err
	}
	if err := f.SetRowHeight(startingCell.SheetName(), 2, 30); err != nil {
		return err
	}
	if err := f.SetRowHeight(startingCell.SheetName(), 3, 20); err != nil {
		return err
	}
	if err := f.SetColWidth(startingCell.SheetName(), "A", "A", 3); err != nil {
		return err
	}

	// write day by day
	dayCursor := startingCell.Copy()

	for _, groupByDay := range groupedByStartDate {
		log.Infof("writing day %v", groupByDay.Key)

		dayCursor.MoveRow(startingCell.Row())

		tracker := model.NewCellWithTracker(dayCursor)

		// code block for variable scoping
		{
			err := writeDayGrid(wc, groupByDay, tracker, log)
			if err != nil {
				return fmt.Errorf("error writing row: %w", err)
			}
		}
		tracker.MoveAtBottomLeftOfCoveredArea()

		schoolGroupsForThisDay := GetDifferentSchoolGroups(groupByDay.Rows)
		if len(schoolGroupsForThisDay) > 0 {
			tracker.MoveBottom(1)

			err := writeSchoolsForDay(wc, schoolGroupsForThisDay, tracker)
			if err != nil {
				return fmt.Errorf("error writing schools for day: %w", err)
			}
			tracker.MoveAtBottomLeftOfCoveredArea()
		}

		// code block for variable scoping
		{
			err := writePlaceholdersForDay(wc, tracker)
			if err != nil {
				return fmt.Errorf("error writing placeholders for OOD: %w", err)
			}
			tracker.MoveAtBottomLeftOfCoveredArea()
		}

		dayCursor.MoveColumn(tracker.CoveredArea().TopRight().Column() + 1)
	}

	dayCursor.MoveRow(startingCell.Row()).MoveBottom(2).MoveRight(1)
	{
		// write operator colours
		err := writeOperatorsLegenda(wc, dayCursor)
		if err != nil {
			return fmt.Errorf("error writing operator colors: %w", err)
		}
	}

	// Save spreadsheet by the given path.
	if err := saveToOutputFile(computeOutputFile(input), f); err != nil {
		return err
	}
	return nil
}

func writeDayGrid(c WriteContext, day model.GroupedRows, startCell model.CellWithTracker, log *logrus.Logger) error {
	minGroupWidth := uint(12)
	minGroupWidthPerSlot := uint(5)
	minDayWidthInCells := uint(20)
	boxCellHeight := float64(20)
	moreSlotsAtBottom := 0 // add if you want to show some empty time rows after the last one

	cursor := startCell.Copy()
	f := c.outputFile

	cursor.MoveBottom(1)

	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(),
		fmt.Sprintf("%d", day.Rows[0].StartAt.Year())); err != nil {
		return err
	}
	cursor.MoveRight(1)
	cursor.MoveBottom(1)

	groupedByRoom := GroupByRoom(c.allData, day.Rows)

	roomColumnNumbers := make(map[string]uint)
	roomWidths := make(map[string]uint)

	// write the header with the columns / rooms

	for _, group := range groupedByRoom {
		log.Infof("writing header for room group %s", group.Key)

		room := c.allData.RoomsMap[group.Key]
		numActs := uint(len(group.Rows))

		if room.Code == "" && numActs == 0 {
			continue
		}

		roomInfo := database.GetEntryForRoom(room.Code)
		if numActs == 0 && roomInfo.Slots == 0 {
			continue
		}

		if roomInfo.Slots > 0 && numActs < roomInfo.Slots {
			numActs = roomInfo.Slots
		}

		toWrite := room.Name
		if room.Code == "" {
			toWrite = "???"
		}

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		roomColumnNumbers[room.Code] = cursor.Column()
		roomWidths[room.Code] = numActs

		desiredWidth := minGroupWidth

		if numActs > 1 {
			mergeUpTo := cursor.AtRight(numActs - 1)

			if err := f.MergeCell(cursor.SheetName(), cursor.Code(), mergeUpTo.Code()); err != nil {
				return err
			}

			altWidth := numActs * minGroupWidthPerSlot
			if altWidth > desiredWidth {
				desiredWidth = altWidth
			}

			if err := f.SetColWidth(
				cursor.SheetName(),
				cursor.ColumnName(),
				mergeUpTo.ColumnName(),
				float64(desiredWidth)/float64(numActs),
			); err != nil {
				return err
			}

		} else {
			if err := f.SetColWidth(cursor.SheetName(), cursor.ColumnName(), cursor.ColumnName(), float64(desiredWidth)); err != nil {
				return err
			}
		}

		cursor.MoveRight(numActs)

		// add 1 additional column to leave some space
		if err := f.SetColWidth(cursor.SheetName(), cursor.ColumnName(), cursor.ColumnName(), 2); err != nil {
			return err
		}

		cursor.MoveRight(1)
	}

	// ensure the day block is wide at least as the desired minDayWidthInCells
	actualWidth := cursor.Column() - startCell.Column()
	if actualWidth < minDayWidthInCells {
		diff := uint(minDayWidthInCells - actualWidth)

		cursorTo := cursor.AtRight(diff)
		if err := f.SetColWidth(cursor.SheetName(), cursor.ColumnName(), cursorTo.ColumnName(), float64(minGroupWidthPerSlot)); err != nil {
			return err
		}

		cursor.MoveRight(diff)
	}
	actualWidth = cursor.Column() - startCell.Column()

	// reset cursor to box start
	cursor = startCell.AtBottom(3)

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

	numOfTimeRows := uint(0)
	log.Infof("writing time columns starting from %s", currentTime.String())
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

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}
		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), boxCellHeight); err != nil {
			return err
		}
		cursor.MoveBottom(1)
		numOfTimeRows++

		currentTime = currentTime.Add(time.Duration(c.minutesStep) * time.Minute)

		if currentTime.Day() != firstTime.Day() {
			rolledOver = true
		}

		if !breaking {
			if effectiveHour > c.maxHour {
				breaking = true
			}
		}

		if breaking {
			if moreSlotsAtBottom > 0 {
				moreSlotsAtBottom--
			} else {
				break
			}
		}
	}

	// finished writing the times on the left

	// draw the box for the whole day
	boxStart := startCell.Copy()
	boxEnd := startCell.Copy().AtRight(actualWidth - 1).AtBottom(numOfTimeRows + 2)
	if err := f.SetCellStyle(startCell.SheetName(), boxStart.Code(), boxEnd.Code(), database.DayBoxStyleID()); err != nil {
		return err
	}

	// draw the box for each room
	for _, group := range groupedByRoom {
		log.Infof("drawing boxes for room %s", group.Key)

		if _, ok := roomWidths[group.Key]; !ok {
			log.Infof("skipping room %s", group.Key)
			continue
		}

		columnsForRoomStartAt := roomColumnNumbers[group.Key]

		boxStart := startCell.AtColumn(columnsForRoomStartAt).AtBottom(3)
		boxEnd := boxStart.AtRight(roomWidths[group.Key] - 1).AtBottom(numOfTimeRows - 1)

		styleToUse := database.DayRoomBoxStyleID()
		if len(group.Rows) == 0 {
			styleToUse = database.UnusedRoomStyleID()
		}
		if err := f.SetCellStyle(startCell.SheetName(), boxStart.Code(), boxEnd.Code(), styleToUse); err != nil {
			return err
		}

		boxHeaderStart := boxStart.AtTop(1)
		boxHeaderEnd := boxHeaderStart.AtColumn(boxEnd.Column())

		roomInfo := database.GetEntryForRoom(group.Key)
		if err := f.SetCellStyle(cursor.SheetName(), boxHeaderStart.Code(), boxHeaderEnd.Code(), roomInfo.Styles.Common.StyleID); err != nil {
			return err
		}
	}

	// now place the activities
	cursor.MoveRow(startCell.Row())
	cursor.MoveColumn(startCell.Column())

	schoolGroups := GetDifferentSchoolGroups(day.Rows)
	schoolGroupsIndex := make(map[string]model.SchoolGroup)
	for _, sg := range schoolGroups {
		if sg.Codice == "" {
			continue
		}
		schoolGroupsIndex[sg.Codice] = sg
	}

	for _, group := range groupedByRoom {
		log.Infof("placing activities for room %s", group.Key)
		roomInfo := database.GetEntryForRoom(group.Key)

		for actIndex, act := range group.Rows {
			var operator model.Operator
			if act.Operator.Code != "" {
				operator, _ = c.allData.OperatorsMap[act.Operator.Code]
			}
			if operator.Code == "" {
				operator.Name = "???"
			}

			columnsForRoomStartAt := roomColumnNumbers[act.Room.Code]
			cursor.MoveColumn(columnsForRoomStartAt + uint(actIndex))
			cursor.MoveRow(startCell.Row() + 3)

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
						actStartCell.MoveRow(startCell.Row() + 3 + uint(i))
					}
				} else {
					if !timeCursor.Before(act.EndAt) {
						inRange = false
						exited = true
						actEndCell.MoveRow(startCell.Row() + 3 + uint(i-1))
					}
				}

				if exited {
					break
				}

				timeCursor = timeCursor.Add(time.Minute * time.Duration(c.minutesStep))
				i++
			}

			// if merging the cells is NOT desired:
			for r := actStartCell.Row(); r <= actEndCell.Row(); r++ {
				c := actStartCell.AtRow(r)

				/*
					// writeInCell := operator.Name
					writeInCell := act.Raw.Codice
					if decoded, ok := schoolGroupsIndex[act.Raw.Codice]; ok && decoded.NumeroSeq > 0 {
						writeInCell = fmt.Sprintf("%d", decoded.NumeroSeq)
					}
				*/
				writeInCell := shortenGroupCode(act.Raw.Codice)
				if writeInCell == "" && operator.Code != "" {
					writeInCell = operator.Name
				}

				if err := f.SetCellValue(c.SheetName(), c.Code(), writeInCell); err != nil {
					return err
				}
			}

			// apply style depending on the operator
			var style int
			if act.Operator.Code == "" && roomInfo.AllowMissingOperator {
				style = database.NoOperatorNeededStyleID()
			} else {
				operatorInfo := database.GetEntryForOperator(act.Operator.Code)
				style = operatorInfo.Styles.Common.StyleID
				if len(act.Warnings) > 0 && operatorInfo.Styles.Warning.Style != nil {
					style = operatorInfo.Styles.Warning.StyleID
				}
			}

			if err := f.SetCellStyle(actStartCell.SheetName(), actStartCell.Code(), actEndCell.Code(), style); err != nil {
				return err
			}

			// write some comment with additional info
			cellComment := buildComment(act)
			if cellComment != "" {
				comment(f, actStartCell, strings.TrimSpace(cellComment))
			}
		}
	}

	// write the day name in the header
	cursor = startCell.AtRight(1).AtBottom(1)
	toCell := cursor.AtRight(8)
	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), toCell.Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(),
		day.Rows[0].StartAt.Format("Mon 2 January")); err != nil {
		return err
	}
	if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(), database.DayHeaderStyleID()); err != nil {
		return err
	}

	// write "MAT: ? POM: ?" placeholder
	cursor = startCell.AtRight(actualWidth - 8).AtRow(1)
	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "MAT:"); err != nil {
		return err
	}
	if err := f.MergeCell(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).AtRight(2).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.AtBottom(1).Code(), "?"); err != nil {
		return err
	}
	if err := f.SetCellStyle(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).Code(), database.ToBeFilledStyleID()); err != nil {
		return err
	}

	cursor = cursor.AtRight(4)
	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "POM:"); err != nil {
		return err
	}
	if err := f.MergeCell(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).AtRight(2).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.AtBottom(1).Code(), "?"); err != nil {
		return err
	}
	if err := f.SetCellStyle(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).Code(), database.ToBeFilledStyleID()); err != nil {
		return err
	}

	return nil
}

func writeSchoolsForDay(c WriteContext, groups []model.SchoolGroup, startCell model.Cell) error {
	f := c.outputFile
	cursor := startCell.Copy()

	for _, schoolGroup := range groups {
		cursor.MoveColumn(startCell.Column())

		toWrite := fmt.Sprintf("%s", schoolGroup.Codice)
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}
		cursor.MoveRight(1)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(5).Code()); err != nil {
			return err
		}

		toWrite = schoolGroup.NomeScuola
		if schoolGroup.TipologiaScuola != "" {
			toWrite += "\n" + schoolGroup.TipologiaScuola
		}

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveRight(6)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(2).Code()); err != nil {
			return err
		}
		toWrite = ""
		if schoolGroup.Classe != "" {
			toWrite = schoolGroup.Classe
		}
		if schoolGroup.Sezione != "" {
			toWrite += " " + schoolGroup.Sezione
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		cursor.MoveRight(3)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(4).Code()); err != nil {
			return err
		}

		toWrite = fmt.Sprintf("%d (%d",
			schoolGroup.NumPaganti+schoolGroup.NumGratuiti+schoolGroup.NumAccompagnatori,
			schoolGroup.NumPaganti)
		if schoolGroup.NumAccompagnatori > 0 {
			toWrite += fmt.Sprintf(", %d acc.", schoolGroup.NumAccompagnatori)
		}
		if schoolGroup.NumGratuiti > 0 {
			toWrite += fmt.Sprintf("+ %d GRAT.", schoolGroup.NumGratuiti)
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

func writePlaceholdersForDay(c WriteContext, startCell model.Cell) error {
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
	}

	cursor = startCell.AtBottom(1)

	for _, rowToWrite := range rowsToWrite {
		valueCursors := cursor.AtRight(7)
		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), valueCursors.AtLeft(1).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), rowToWrite); err != nil {
			return err
		}

		if err := f.SetCellValue(cursor.SheetName(), valueCursors.Code(), "?"); err != nil {
			return err
		}
		if err := f.MergeCell(cursor.SheetName(), valueCursors.Code(), valueCursors.AtRight(7).Code()); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), valueCursors.Code(), valueCursors.Code(),
			database.ToBeFilledStyleID()); err != nil {
			return err
		}

		if err := f.SetRowHeight(cursor.SheetName(), int(cursor.Row()), 25); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}

func writeOperatorsLegenda(c WriteContext, startCell model.Cell) error {
	f := c.outputFile
	cursor := startCell.Copy()

	if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(3).Code()); err != nil {
		return err
	}
	if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), "LEGENDA COLORI / EDUCATORI"); err != nil {
		return err
	}

	cursor.MoveBottom(1)

	for _, op := range c.allData.Operators {
		operatorInfo := database.GetEntryForOperator(op.Code)

		if err := f.MergeCell(cursor.SheetName(), cursor.Code(), cursor.AtRight(3).Code()); err != nil {
			return err
		}
		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), op.Name); err != nil {
			return err
		}
		if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(),
			operatorInfo.Common.StyleID); err != nil {
			return err
		}

		cursor.MoveBottom(1)
	}

	return nil
}

func buildComment(act model.ParsedRow) string {
	cellComment := ``

	for _, warning := range act.Warnings {
		cellComment += "ATTENZIONE: " + warning + "\n\n"
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

func comment(f *excelize.File, cell model.Cell, content string) error {
	type serializable struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	v := serializable{
		Author: "planner",
		Text:   content,
	}
	serialized, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return f.AddComment(cell.SheetName(), cell.Code(), string(serialized))
}

func shortenGroupCode(code string) string {
	if len(code) >= 4 && strings.Contains(code, "-") && strings.HasPrefix(code, "P") {
		return strings.TrimPrefix(code, "P")
	}
	return code
}

func computeOutputFile(inputFile string) string {
	outPath := filepath.Dir(inputFile)
	inputName := filepath.Base(inputFile)
	inputExt := filepath.Ext(inputFile)
	return outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-convertito" + inputExt
}

func saveToOutputFile(outputFile string, f *excelize.File) error {
	if err := f.SaveAs(outputFile); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}
	return nil
}
