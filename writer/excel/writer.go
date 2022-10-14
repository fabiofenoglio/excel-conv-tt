package excel

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/aggregator/byactivity"
	"github.com/fabiofenoglio/excelconv/aggregator/byday"
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
	if err := f.SetRowHeight(startingCell.SheetName(), 2, 20); err != nil {
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
	aggregateDayWriter := DayGridWriteResult{}

	for _, groupByDay := range groupedByStartDate {
		log.Debugf("writing day %v", groupByDay.Key)

		dayWriteResult, err := writeDayWithDetails(wc, groupByDay, daysGridCursor, log)
		if err != nil {
			return nil, fmt.Errorf("error writing day: %w", err)
		}
		aggregateDayWriter.Aggregate(dayWriteResult)

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

	for _, groupByDay := range groupedByStartDate {
		if err := writeAnnotationsOnActivities(wc, groupByDay, aggregateDayWriter, log); err != nil {
			return nil, fmt.Errorf("error writing activity annotations: %w", err)
		}
	}

	// Save spreadsheet by the given path.
	out, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error writing output to buffer: %w", err)
	}
	return out.Bytes(), nil
}

func writeDayWithDetails(c WriteContext, groupByDay byday.GroupedByDay, startCell excel.Cell, log *logrus.Logger) (DayGridWriteResult, error) {
	zero := DayGridWriteResult{}

	tracker := excel.NewCellWithTracker(startCell.Copy())

	// WRITE DAY GRID W/ ROOMS AND ACTIVITIES
	var out DayGridWriteResult
	{
		outFromDayGridWriter, err := writeDayGrid(c, groupByDay, tracker, log)
		if err != nil {
			return zero, fmt.Errorf("error writing row: %w", err)
		}
		out = outFromDayGridWriter
	}
	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE SCHOOL/GROUPS FOR THE DAY
	schoolGroupsForThisDay := byactivity.GetDifferentActivityGroups(groupByDay.Rows)
	if len(schoolGroupsForThisDay) > 0 {
		err := writeSchoolsForDay(c, schoolGroupsForThisDay, tracker)
		if err != nil {
			return zero, fmt.Errorf("error writing schools for day: %w", err)
		}
	}
	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE PLACEHOLDERS FOR ORDER OF THE DAY
	{
		err := writePlaceholdersForDay(c, tracker)
		if err != nil {
			return zero, fmt.Errorf("error writing placeholders for OOD: %w", err)
		}
	}

	return out, nil
}
