package excel

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/aggregator/byactivity"
	"github.com/fabiofenoglio/excelconv/aggregator/byday"
	"github.com/fabiofenoglio/excelconv/aggregator/byroom"
	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/fabiofenoglio/excelconv/model"
	"github.com/fabiofenoglio/excelconv/writer"
)

const (
	layoutTimeOnlyInReadableFormat = "15:04"
)

type WriteContext struct {
	args          config.Args
	minHour       int
	maxHour       int
	minutesStep   int
	allData       model.ParsedData
	outputFile    *excelize.File
	styleRegister *StyleRegister
}

type WriterImpl struct{}

var _ writer.Writer = &WriterImpl{}

func (w *WriterImpl) ComputeDefaultOutputFile(inputFile string) string {
	outPath := filepath.Dir(inputFile)
	inputName := filepath.Base(inputFile)
	inputExt := filepath.Ext(inputFile)
	return outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-convertito.xlsx"
}

func (w *WriterImpl) Write(parsed model.ParsedData, args config.Args, log *logrus.Logger) ([]byte, error) {
	log.Debug("writing with excel writer")

	// compute the max time range to be shown between all days

	minHourToShow, maxHourToShow := GetMaxHoursRangeToDisplay(parsed.Rows)
	log.Debugf("will show the range %d to %d (inclusive)", minHourToShow, maxHourToShow)

	// group by start date, ordering by start time ASC

	groupedByStartDate := byday.GroupByStartDay(parsed.Rows)
	log.Infof("found activities spanning %d different days", len(groupedByStartDate))

	f := excelize.NewFile()

	wc := WriteContext{
		args:          args,
		minHour:       minHourToShow,
		maxHour:       maxHourToShow,
		minutesStep:   15,
		allData:       parsed,
		outputFile:    f,
		styleRegister: NewStyleRegister(f),
	}

	sheetName := fmt.Sprintf(
		"settimana %s-%s",
		groupedByStartDate[0].Rows[0].StartAt.Format("0201"),
		groupedByStartDate[len(groupedByStartDate)-1].Rows[0].StartAt.Format("0201"),
	)

	startingCell := excel.NewCell(sheetName, 2, 1)

	// Rename the default sheet and set it as active
	f.SetSheetName("Sheet1", startingCell.SheetName())
	f.SetActiveSheet(f.GetSheetIndex(startingCell.SheetName()))

	if err := f.SetRowHeight(startingCell.SheetName(), 1, 20); err != nil {
		return nil, err
	}
	if err := f.SetRowHeight(startingCell.SheetName(), 2, 30); err != nil {
		return nil, err
	}
	if err := f.SetRowHeight(startingCell.SheetName(), 3, 20); err != nil {
		return nil, err
	}
	if err := f.SetColWidth(startingCell.SheetName(), "A", "A", 3); err != nil {
		return nil, err
	}

	// write day by day
	daysGridCursor := excel.NewCellWithTracker(startingCell)

	for _, groupByDay := range groupedByStartDate {
		log.Debugf("writing day %v", groupByDay.Key)

		err := writeDayWithDetails(wc, groupByDay, daysGridCursor, log)
		if err != nil {
			return nil, fmt.Errorf("error writing day: %w", err)
		}

		// MOVE CURSOR IN POSITION FOR NEXT DAY
		daysGridCursor.
			MoveRow(startingCell.Row()).
			MoveColumn(daysGridCursor.CoveredArea().RightColumn() + 1)
	}

	daysGridCursor.MoveAtRightTopOfCoveredArea()
	{
		// write operator colours
		err := writeOperatorsLegenda(wc, daysGridCursor, parsed.Operators)
		if err != nil {
			return nil, fmt.Errorf("error writing operator colors: %w", err)
		}
	}

	// Save spreadsheet by the given path.
	out, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error writing output to buffer: %w", err)
	}
	return out.Bytes(), nil
}

func writeDayWithDetails(c WriteContext, groupByDay byday.GroupedByDay, startCell excel.Cell, log *logrus.Logger) error {

	tracker := excel.NewCellWithTracker(startCell.Copy())

	// WRITE DAY GRID W/ ROOMS AND ACTIVITIES
	{
		err := writeDayGrid(c, groupByDay, tracker, log)
		if err != nil {
			return fmt.Errorf("error writing row: %w", err)
		}
	}
	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE SCHOOL/GROUPS FOR THE DAY
	schoolGroupsForThisDay := byactivity.GetDifferentActivityGroups(groupByDay.Rows)
	if len(schoolGroupsForThisDay) > 0 {
		err := writeSchoolsForDay(c, schoolGroupsForThisDay, tracker)
		if err != nil {
			return fmt.Errorf("error writing schools for day: %w", err)
		}
	}
	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE PLACEHOLDERS FOR ORDER OF THE DAY
	{
		err := writePlaceholdersForDay(c, tracker)
		if err != nil {
			return fmt.Errorf("error writing placeholders for OOD: %w", err)
		}
	}

	return nil
}

func writeDayGrid(c WriteContext, day byday.GroupedByDay, startCell excel.Cell, log *logrus.Logger) error {
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

	groupedByRoom := byroom.GroupByRoom(c.allData, day.Rows)

	roomColumnNumbers := make(map[string]uint)
	roomWidths := make(map[string]uint)

	// write the header with the columns / rooms

	for _, group := range groupedByRoom {
		log.Debugf("writing header for room group %s", group.Room.Code)

		numActs := uint(len(group.Slots))

		if group.Room.Code == "" && numActs == 0 {
			continue
		}
		if numActs == 0 && group.Room.Slots == 0 {
			continue
		}

		if group.Room.Slots > 0 && numActs < group.Room.Slots {
			numActs = group.Room.Slots
		}

		toWrite := group.Room.Name
		if group.Room.Code == "" {
			toWrite = "???"
		}

		if err := f.SetCellValue(cursor.SheetName(), cursor.Code(), toWrite); err != nil {
			return err
		}

		roomColumnNumbers[group.Room.Code] = cursor.Column()
		roomWidths[group.Room.Code] = numActs

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
	log.Debugf("writing time columns starting from %s", currentTime.String())
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
	if err := f.SetCellStyle(
		startCell.SheetName(), boxStart.Code(), boxEnd.Code(),
		c.styleRegister.DayBoxStyle().Common.StyleID,
	); err != nil {
		return err
	}

	// draw the box for each room
	for _, group := range groupedByRoom {
		log.Debugf("drawing boxes for room %s", group.Room.Code)

		if _, ok := roomWidths[group.Room.Code]; !ok {
			log.Debugf("skipping room %s", group.Room.Code)
			continue
		}

		columnsForRoomStartAt := roomColumnNumbers[group.Room.Code]

		boxStart := startCell.AtColumn(columnsForRoomStartAt).AtBottom(3)
		boxEnd := boxStart.AtRight(roomWidths[group.Room.Code] - 1).AtBottom(numOfTimeRows - 1)

		styleToUse := c.styleRegister.DayRoomBoxStyle().Common.StyleID
		if group.TotalInAllSlots == 0 {
			styleToUse = c.styleRegister.UnusedRoomStyle().Common.StyleID
		}
		if err := f.SetCellStyle(startCell.SheetName(), boxStart.Code(), boxEnd.Code(), styleToUse); err != nil {
			return err
		}

		boxHeaderStart := boxStart.AtTop(1)
		boxHeaderEnd := boxHeaderStart.AtColumn(boxEnd.Column())

		var style *Style
		if group.Room != model.NoRoom {
			style = c.styleRegister.RoomStyle(group.Room.BackgroundColor)
		} else {
			style = c.styleRegister.NoRoomStyle()
		}

		if err := f.SetCellStyle(cursor.SheetName(), boxHeaderStart.Code(), boxHeaderEnd.Code(), style.Common.StyleID); err != nil {
			return err
		}
	}

	// now place the activities
	cursor.MoveRow(startCell.Row())
	cursor.MoveColumn(startCell.Column())

	schoolGroups := byactivity.GetDifferentActivityGroups(day.Rows)
	schoolGroupsIndex := make(map[string]byactivity.ActivityGroup)
	for _, sg := range schoolGroups {
		if sg.Code == "" {
			continue
		}
		schoolGroupsIndex[sg.Code] = sg
	}

	for _, group := range groupedByRoom {
		log.Debugf("placing activities for room %s", group.Room.Code)

		for slotIndex, slot := range group.Slots {
			for _, act := range slot.Rows {
				var operator model.Operator
				if act.Operator.Code != "" {
					operator, _ = c.allData.OperatorsMap[act.Operator.Code]
				}
				if operator.Code == "" {
					operator.Name = "???"
				}

				columnsForRoomStartAt := roomColumnNumbers[act.Room.Code]
				cursor.MoveColumn(columnsForRoomStartAt + uint(slotIndex))
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
					writeInCell := shortenGroupCode(act.Code)
					if writeInCell == "" && operator.Code != "" {
						writeInCell = operator.Name
					}

					if err := f.SetCellValue(c.SheetName(), c.Code(), writeInCell); err != nil {
						return err
					}
				}

				// apply style depending on the operator
				var style *Style
				if act.Operator.Code == "" && (group.Room.AllowMissingOperator || !c.args.EnableMissingOperatorsWarning) {
					style = c.styleRegister.NoOperatorNeededStyle()
				} else if act.Operator.Code == "" {
					style = c.styleRegister.NoOperatorStyle()
				} else {
					style = c.styleRegister.OperatorStyle(operator.BackgroundColor)
				}

				styleID := style.Common.StyleID
				if len(act.Warnings) > 0 {
					styleID = style.WarningIDOrDefault()
				}

				if err := f.SetCellStyle(actStartCell.SheetName(), actStartCell.Code(), actEndCell.Code(), styleID); err != nil {
					return err
				}

				// write some comment with additional info
				cellComment := buildContentOfActivityComment(act)
				if cellComment != "" {
					addCommentToCell(f, actStartCell, strings.TrimSpace(cellComment))
				}
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
	if err := f.SetCellStyle(cursor.SheetName(), cursor.Code(), cursor.Code(), c.styleRegister.DayHeaderStyle().Common.StyleID); err != nil {
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
	if err := f.SetCellStyle(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).Code(),
		c.styleRegister.ToBeFilledStyle().Common.StyleID); err != nil {
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
	if err := f.SetCellStyle(cursor.SheetName(), cursor.AtBottom(1).Code(), cursor.AtBottom(1).Code(),
		c.styleRegister.ToBeFilledStyle().Common.StyleID); err != nil {
		return err
	}

	return nil
}

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
	}

	cursor = startCell.AtBottom(1)

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

func writeOperatorsLegenda(c WriteContext, startCell excel.Cell, operators []model.Operator) error {
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

func buildContentOfActivityComment(act model.ParsedRow) string {
	cellComment := ``

	for _, warning := range act.Warnings {
		cellComment += "ATTENZIONE: " + warning.Message + "\n\n"
	}

	if act.Activity.Description != "" {
		cellComment += act.Activity.Description
		if act.Activity.Type != "" {
			cellComment += " (" + act.Activity.Type + ")"
		}
		cellComment += "\n\n"

	} else if act.Activity.Type != "" {
		cellComment += "Tipologia: " + act.Activity.Type + "\n"
	}

	if act.TimeString != "" {
		cellComment += "Orario: " + act.TimeString + "\n"
	}
	if act.Operator.Code != "" {
		cellComment += "Educatore: " + act.Operator.Name + "\n"
	}
	if act.Room.Code != "" {
		cellComment += "Aula: " + act.Room.Name + "\n"
	}

	if act.OperatorNotes != "" {
		cellComment += "Nota operatore: " + act.OperatorNotes + "\n"
	}
	if act.BookingNotes != "" {
		cellComment += "Nota prenotazione: " + act.BookingNotes + "\n"
	}

	if act.SchoolClass.FullDescription() != "" {
		cellComment += "Classe: " + act.SchoolClass.FullDescription() + "\n"
	}

	if act.GroupComposition.Total > 0 {
		c := ""
		entries := 0
		if act.GroupComposition.NumPaying > 0 {
			c += fmt.Sprintf("%d paganti, ", act.GroupComposition.NumPaying)
			entries++
		}
		if act.GroupComposition.NumFree > 0 {
			c += fmt.Sprintf("%d gratuiti, ", act.GroupComposition.NumFree)
			entries++
		}
		if act.GroupComposition.NumAccompanying > 0 {
			c += fmt.Sprintf("%d accompagnatori, ", act.GroupComposition.NumAccompanying)
			entries++
		}
		if entries > 1 {
			c = strings.TrimSuffix(c, ", ") + fmt.Sprintf(" (%d totali)", act.GroupComposition.Total)
		}
		cellComment += strings.TrimSuffix(c, ", ") + "\n"
	}

	if act.School.FullDescription() != "" {
		cellComment += act.School.FullDescription() + "\n"
	}

	if act.Bus != "" {
		cellComment += "Bus: " + act.Bus + "\n"
	}
	if act.PaymentStatus.PaymentAdvance != "" && act.PaymentStatus.PaymentAdvance != "-" {
		cellComment += "Acconti: " + act.PaymentStatus.PaymentAdvance + "\n"
	}
	if act.PaymentStatus.PaymentAdvanceStatus != "" && act.PaymentStatus.PaymentAdvanceStatus != "-" {
		cellComment += "Stato acconti: " + act.PaymentStatus.PaymentAdvanceStatus + "\n"
	}

	return strings.TrimSpace(cellComment)
}

func addCommentToCell(f *excelize.File, cell excel.Cell, content string) error {
	type serializable struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	v := serializable{
		Author: "planner ",
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
