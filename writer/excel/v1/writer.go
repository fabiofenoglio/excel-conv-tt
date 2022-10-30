package excel

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
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
			return nil, errors.Wrap(err, "error writing day")
		}
		aggregateDayWriter.Aggregate(dayWriteResult)

		// MOVE CURSOR IN POSITION FOR NEXT DAY
		daysGridCursor.
			MoveRow(startingCell.Row()).
			MoveColumn(daysGridCursor.CoveredArea().RightColumn() + 1)
	}

	if false { // do not write operators legenda
		daysGridCursor.MoveAtRightTopOfCoveredArea()
		{
			// write operator colours
			err := writeOperatorsLegenda(wc, daysGridCursor, parsed.Operators)
			if err != nil {
				return nil, errors.Wrap(err, "error writing operator colors")
			}
		}
	}

	for _, groupByDay := range groupedByStartDate {
		if err := writeAnnotationsOnActivities(wc, groupByDay, aggregateDayWriter, log); err != nil {
			return nil, errors.Wrap(err, "error writing activity annotations")
		}
	}

	// Save spreadsheet by the given path.
	out, err := f.WriteToBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "error writing output to buffer")
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
			return zero, errors.Wrap(err, "error writing row")
		}
		out = outFromDayGridWriter
	}
	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE SCHOOL/GROUPS FOR THE DAY
	schoolGroupsForThisDay := byactivity.GetDifferentActivityGroups(groupByDay.Rows)
	if len(schoolGroupsForThisDay) > 0 {
		sort.SliceStable(schoolGroupsForThisDay, func(i, j int) bool {
			diff := int(schoolGroupsForThisDay[i].SequentialNumberForSchool) - int(schoolGroupsForThisDay[j].SequentialNumberForSchool)
			if diff != 0 {
				return diff < 0
			}
			diff = int(schoolGroupsForThisDay[i].SequentialNumberInsideSchool) - int(schoolGroupsForThisDay[j].SequentialNumberInsideSchool)
			if diff != 0 {
				return diff < 0
			}
			r := strings.Compare(strings.ToLower(schoolGroupsForThisDay[i].Code), strings.ToLower(schoolGroupsForThisDay[j].Code))
			return r < 0
		})

		err := writeSchoolsForDay(c, schoolGroupsForThisDay, tracker)
		if err != nil {
			return zero, errors.Wrap(err, "error writing schools for day")
		}
	}

	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE PLACEHOLDERS FOR ORDER OF THE DAY
	{
		err := writePlaceholdersForDay(c, tracker)
		if err != nil {
			return zero, errors.Wrap(err, "error writing placeholders for OOD")
		}
	}

	return out, nil
}
