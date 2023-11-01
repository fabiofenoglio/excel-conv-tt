package excel

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/fabiofenoglio/excelconv/writer"
)

const (
	layoutTimeOnlyInReadableFormat = "15:04"
)

type WriteContext struct {
	minHour        aggregator2.TimeOfDay
	maxHour        aggregator2.TimeOfDay
	minutesStep    int
	allData        aggregator2.Output
	anagraphicsRef *parser2.OutputAnagraphics
	outputFile     *excelize.File
	styleRegister  *StyleRegister
}

type WriterImpl struct{}

var _ writer.Writer = &WriterImpl{}

func (w *WriterImpl) ComputeDefaultOutputFile(inputFile string) string {
	outPath := filepath.Dir(inputFile)
	inputName := filepath.Base(inputFile)
	inputExt := filepath.Ext(inputFile)
	return outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-convertito.xlsx"
}

func (w *WriterImpl) Write(ctx config.WorkflowContext, parsed aggregator2.Output, anagraphicsRef *parser2.OutputAnagraphics) ([]byte, error) {
	log := ctx.Logger
	log.Debug("writing with excel writer")

	// compute the max time range to be shown between all days

	minHourToShow, maxHourToShow := parsed.CommonData.CommonTimespan.Start, parsed.CommonData.CommonTimespan.End
	log.Debugf("will show the range %d to %d (inclusive)", minHourToShow, maxHourToShow)

	// group by start date, ordering by start time ASC

	groupedByStartDate := parsed.Days
	log.Infof("found activities spanning %d different days", len(groupedByStartDate))
	if len(groupedByStartDate) == 0 {
		return nil, errors.New("nessun dato valido presente nel file di input")
	}

	f := excelize.NewFile()

	wc := WriteContext{
		minHour:        minHourToShow,
		maxHour:        maxHourToShow,
		minutesStep:    15,
		allData:        parsed,
		anagraphicsRef: anagraphicsRef,
		outputFile:     f,
		styleRegister:  NewStyleRegister(f),
	}

	sheetName := fmt.Sprintf(
		"settimana %s-%s",
		groupedByStartDate[0].Day.Format("0201"),
		groupedByStartDate[len(groupedByStartDate)-1].Day.Format("0201"),
	)

	startingCell := excel.NewCell(sheetName, 2, 1)

	// Rename the default sheet and set it as active
	if err := f.SetSheetName("Sheet1", startingCell.SheetName()); err != nil {
		return nil, fmt.Errorf("error setting sheet name: %w", err)
	}

	index, err := f.GetSheetIndex(startingCell.SheetName())
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

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

	for i, groupByDay := range groupedByStartDate {
		span := sentry.StartSpan(ctx.Context, fmt.Sprintf("write day %d", i))
		log.Debugf("writing day %v", groupByDay.Day)

		dayWriteResult, err := writeDayWithDetails(ctx, wc, groupByDay, daysGridCursor)
		if err != nil {
			span.Finish()
			return nil, errors.Wrap(err, "error writing day")
		}
		aggregateDayWriter.Aggregate(dayWriteResult)

		// MOVE CURSOR IN POSITION FOR NEXT DAY
		daysGridCursor.
			MoveRow(startingCell.Row()).
			MoveColumn(daysGridCursor.CoveredArea().RightColumn() + 1)

		span.Finish()
	}

	if false { // do not write operators legenda
		span := sentry.StartSpan(ctx.Context, "write operator legenda")
		daysGridCursor.MoveAtRightTopOfCoveredArea()
		// write operator colours
		err := writeOperatorsLegenda(wc, daysGridCursor, anagraphicsRef.Operators)
		if err != nil {
			span.Finish()
			return nil, errors.Wrap(err, "error writing operator colors")
		}
		span.Finish()
	}

	span := sentry.StartSpan(ctx.Context, "write annotations")
	for _, groupByDay := range groupedByStartDate {
		if err := writeAnnotationsOnActivities(ctx, wc, groupByDay, aggregateDayWriter); err != nil {
			span.Finish()
			return nil, errors.Wrap(err, "error writing activity annotations")
		}
	}
	span.Finish()

	span = sentry.StartSpan(ctx.Context, "write to buffer")
	out, err := f.WriteToBuffer()
	if err != nil {
		span.Finish()
		return nil, errors.Wrap(err, "error writing output to buffer")
	}
	span.Finish()
	return out.Bytes(), nil
}

func writeDayWithDetails(ctx config.WorkflowContext, c WriteContext, groupByDay aggregator2.ScheduleForSingleDayWithRoomsAndGroupSlots, startCell excel.Cell) (DayGridWriteResult, error) {
	zero := DayGridWriteResult{}

	tracker := excel.NewCellWithTracker(startCell.Copy())

	// WRITE DAY GRID W/ ROOMS AND ACTIVITIES
	var out DayGridWriteResult
	{
		outFromDayGridWriter, err := writeDayGrid(ctx, c, groupByDay, tracker)
		if err != nil {
			return zero, errors.Wrap(err, "error writing row")
		}
		out = outFromDayGridWriter
	}
	tracker.MoveAtBottomLeftOfCoveredArea()
	tracker.MoveTop(1)

	if err := c.outputFile.SetRowHeight(tracker.SheetName(), int(tracker.Row()), 10); err != nil {
		return zero, err
	}
	if err := c.outputFile.SetCellStyle(tracker.SheetName(), tracker.Code(), tracker.AtColumn(tracker.CoveredArea().RightColumn()-1).Code(),
		c.styleRegister.Get(softDividerStyle).Bottom()); err != nil {
		return zero, err
	}
	tracker.MoveBottom(1)
	if err := c.outputFile.SetRowHeight(tracker.SheetName(), int(tracker.Row()), 10); err != nil {
		return zero, err
	}
	tracker.MoveBottom(1)

	numAvailableColumns := tracker.CoveredArea().RightColumn() - tracker.Column()

	// WRITE SCHOOL/GROUPS FOR THE DAY
	schoolGroupsForThisDay := groupByDay.VisitingGroups
	if len(schoolGroupsForThisDay) > 0 {
		sort.SliceStable(schoolGroupsForThisDay, func(i, j int) bool {
			r := strings.Compare(strings.ToLower(schoolGroupsForThisDay[i].SequentialCode), strings.ToLower(schoolGroupsForThisDay[j].SequentialCode))
			return r < 0
		})

		err := writeSchoolsForDay(c, schoolGroupsForThisDay, tracker, numAvailableColumns)
		if err != nil {
			return zero, errors.Wrap(err, "error writing schools for day")
		}
	}

	tracker.MoveAtBottomLeftOfCoveredArea()
	tracker.MoveTop(1)

	if err := c.outputFile.SetCellStyle(tracker.SheetName(), tracker.Code(), tracker.AtColumn(tracker.CoveredArea().RightColumn()-1).Code(),
		c.styleRegister.Get(softDividerStyle).Bottom()); err != nil {
		return zero, err
	}
	tracker.MoveBottom(1)

	// WRITE PLACEHOLDERS FOR ORDER OF THE DAY
	{
		err := writePlaceholdersForDay(c, tracker, numAvailableColumns)
		if err != nil {
			return zero, errors.Wrap(err, "error writing placeholders for OOD")
		}
	}

	return out, nil
}
