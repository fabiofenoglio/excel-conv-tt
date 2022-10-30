package excel

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
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

var _ writer.WriterV2 = &WriterImpl{}

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
		log.Debugf("writing day %v", groupByDay.Day)

		dayWriteResult, err := writeDayWithDetails(ctx, wc, groupByDay, daysGridCursor)
		if err != nil {
			return nil, fmt.Errorf("error writing day: %w", err)
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
			err := writeOperatorsLegenda(wc, daysGridCursor, anagraphicsRef.Operators)
			if err != nil {
				return nil, fmt.Errorf("error writing operator colors: %w", err)
			}
		}
	}

	for _, groupByDay := range groupedByStartDate {
		if err := writeAnnotationsOnActivities(ctx, wc, groupByDay, aggregateDayWriter); err != nil {
			return nil, fmt.Errorf("error writing activity annotations: %w", err)
		}
	}

	out, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error writing output to buffer: %w", err)
	}
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
			return zero, fmt.Errorf("error writing row: %w", err)
		}
		out = outFromDayGridWriter
	}
	tracker.MoveAtBottomLeftOfCoveredArea()

	// WRITE SCHOOL/GROUPS FOR THE DAY
	schoolGroupsForThisDay := groupByDay.VisitingGroups
	if len(schoolGroupsForThisDay) > 0 {
		sort.SliceStable(schoolGroupsForThisDay, func(i, j int) bool {
			r := strings.Compare(strings.ToLower(schoolGroupsForThisDay[i].SequentialCode), strings.ToLower(schoolGroupsForThisDay[j].SequentialCode))
			return r < 0
		})

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
