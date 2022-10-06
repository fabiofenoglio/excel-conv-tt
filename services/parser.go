package services

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/xuri/excelize/v2"
)

const (
	layoutTimeOnlyWithSeconds      = "15:04:05"
	layoutTimeOnlyWithMinutes      = "15:04"
	layoutTimeOnlyWithHours        = "15"
	layoutDateOnlyInITFormat       = "02/01/2006"
	layoutDateOrderable            = "2006-01-02"
	layoutTimeOnlyInReadableFormat = "15:04"
)

var (
	localTimeZone *time.Location
)

func init() {
	var err error
	localTimeZone, err = time.LoadLocation("Europe/Rome")
	if err != nil {
		panic(err)
	}
}

type Cell struct {
	C         uint
	R         uint
	SheetName string
}

func (c Cell) ColumnName() string {
	return toColumnName(c.C)
}
func (c Cell) Code() string {
	return fmt.Sprintf("%s%d", toColumnName(c.C), c.R)
}
func (c Cell) String() string {
	if c.SheetName != "" {
		return fmt.Sprintf("!%s:%s", c.SheetName, c.Code())
	}
	return c.Code()
}
func (c Cell) Copy() Cell {
	return Cell{
		C:         c.C,
		R:         c.R,
		SheetName: c.SheetName,
	}
}

func Parse(input string) (Parsed, error) {
	zero := Parsed{}
	log := GetLogger()

	var err error

	f, err := excelize.OpenFile(input)
	if err != nil {
		return zero, fmt.Errorf("error opening input file: %w", err)
	}

	defer func() {
		// Close the spreadsheet.
		if closeErr := f.Close(); closeErr != nil {
			log.Errorf("error closing input file: %s", closeErr.Error())
		}
	}()

	startingHeaderCell := Cell{1, 4, "Organizzazione"}

	columnNameToFieldNameMap := make(map[string]string)
	sampleRow := &Row{}
	val := reflect.ValueOf(sampleRow).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		columnName := field.Tag.Get("column")
		if columnName != "" {
			columnNameToFieldNameMap[columnName] = field.Name
		}
	}

	currentHeaderCell := startingHeaderCell
	headers := make([]string, 0, 10)
	columnNumberToFieldNameMap := make(map[uint]string)
	fieldNameToColumnNumberMap := make(map[string]uint)

	for {
		cell, err := f.GetCellValue(currentHeaderCell.SheetName, currentHeaderCell.Code())
		if err != nil {
			return zero, fmt.Errorf("error reading header cell %s: %w", currentHeaderCell.Code(), err)
		}
		if cell == "" {
			break
		}
		headers = append(headers, cell)
		if len(headers) > 50 {
			return zero, errors.New("too many headers found")
		}

		if mapsToFieldName, ok := columnNameToFieldNameMap[strings.TrimSpace(cell)]; ok {
			log.Infof("column %s with header '%s' maps to known field '%s'", currentHeaderCell.ColumnName(), cell, mapsToFieldName)

			columnNumberToFieldNameMap[currentHeaderCell.C] = mapsToFieldName
			fieldNameToColumnNumberMap[mapsToFieldName] = currentHeaderCell.C
		} else {
			log.Warnf("column %s with header '%s' does not map to any known field", currentHeaderCell.ColumnName(), cell)
		}

		currentHeaderCell.C++
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		columnName := field.Tag.Get("column")

		required := false
		requiredTag := field.Tag.Get("required")
		if requiredTag != "" {
			required, err = strconv.ParseBool(requiredTag)
			if err != nil {
				panic("invalid value for required tag: " + requiredTag)
			}
		}

		if columnName != "" {
			if mappedByColumnNumer, ok := fieldNameToColumnNumberMap[field.Name]; ok {
				log.Infof("field '%s' is mapped by column '%s'", field.Name, toColumnName(mappedByColumnNumer))
			} else {
				errMsg := fmt.Sprintf("field '%s' is NOT mapped by any column", field.Name)
				if required {
					log.Error(errMsg)
					return zero, errors.New(errMsg)
				} else {
					log.Warn(errMsg)
				}
			}
		}
	}

	log.Info("mapping build completed")
	log.Info("all required fields have a valid mapping column")

	results := make([]ParsedRow, 0, 20)

	currentCell := startingHeaderCell
	currentCell.R++

	for {
		row := Row{}
		anyNonNil := false

		for fieldName, columnNumber := range fieldNameToColumnNumberMap {
			cell := Cell{
				R:         currentCell.R,
				C:         columnNumber,
				SheetName: currentCell.SheetName,
			}
			cellContent, err := f.GetCellValue(cell.SheetName, cell.Code())
			if err != nil {
				return zero, fmt.Errorf("error reading content cell %s: %w", cell.String(), err)
			}

			trimmed := strings.TrimSpace(cellContent)
			if trimmed != cellContent {
				log.Warnf("cell %s has a value that starts or finishes with spaces or newline, this should be avoided (value is '%s')",
					cell.Code(), cellContent)
			}
			cellContent = trimmed

			if cellContent != "" {
				anyNonNil = true
				reflect.ValueOf(&row).Elem().FieldByName(fieldName).SetString(cellContent)

				log.Infof("setting row %d.%s to '%s' by cell %s", len(results), fieldName, cellContent, cell.Code())
			}
		}

		if !anyNonNil {
			break
		}

		err := validate(row)
		if err != nil {
			errMsg := fmt.Sprintf("error validating row %d", currentCell.R)
			log.WithError(err).Errorf(errMsg+": %s", err.Error())
			return zero, fmt.Errorf(errMsg+": %w", err)
		}

		parsed, err := parseRow(row)
		if err != nil {
			errMsg := fmt.Sprintf("error parsing row %d", currentCell.R)
			log.WithError(err).Errorf(errMsg+": %s", err.Error())
			return zero, fmt.Errorf(errMsg+": %w", err)
		}

		results = append(results, parsed)

		currentCell.R++
	}

	roomsMap := make(map[string]Room)
	operatorsMap := make(map[string]Operator)

	for _, a := range results {
		if a.Room.Code != "" {
			roomsMap[a.Room.Code] = a.Room
		}
		if a.Operator.Code != "" {
			operatorsMap[a.Operator.Code] = a.Operator
		}
	}

	rooms := make([]Room, 0, 5)
	for _, r := range roomsMap {
		rooms = append(rooms, r)
	}
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Name < rooms[j].Name
	})

	operators := make([]Operator, 0, 5)
	for _, o := range operatorsMap {
		operators = append(operators, o)
	}
	sort.Slice(operators, func(i, j int) bool {
		return operators[i].Name < operators[j].Name
	})

	return Parsed{
		Rows:         results,
		Rooms:        rooms,
		RoomsMap:     roomsMap,
		Operators:    operators,
		OperatorsMap: operatorsMap,
	}, nil
}

func validate(r Row) error {
	if r.Codice == "" {
		return errors.New("codice is required")
	}

	timeRangeRegexp := regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])\s*\-\s*([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`)
	if !timeRangeRegexp.MatchString(r.Orario) {
		return fmt.Errorf("orario non valido, atteso HH:MM-HH:MM e non '%s'", r.Orario)
	}

	dateRegexp := regexp.MustCompile(`^\d{2}\/\d{2}\/\d{4}$`)
	if !dateRegexp.MatchString(r.Data) {
		return fmt.Errorf("data non valida, atteso GG/MM/YYYY e non '%s'", r.Data)
	}

	return nil
}

func parseRow(r Row) (ParsedRow, error) {

	data, err := time.Parse(layoutDateOnlyInITFormat, r.Data)
	if err != nil {
		return ParsedRow{}, fmt.Errorf("il valore '%s' non e' una data valida nel formato 'GG/MM/YYYY'", r.Data)
	}

	orari := strings.Split(r.Orario, "-")

	v := strings.TrimSpace(orari[0])
	start, err := time.Parse(layoutTimeOnlyWithMinutes, v)
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithSeconds, v)
	}
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithHours, v)
	}
	if err != nil {
		return ParsedRow{}, fmt.Errorf("il valore '%s' non e' un orario di inizio valido nei formati 'HH:MM', 'HH:MM:SS' o 'HH'", v)
	}

	v = strings.TrimSpace(orari[1])
	end, err := time.Parse(layoutTimeOnlyWithMinutes, v)
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithSeconds, v)
	}
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithHours, v)
	}
	if err != nil {
		return ParsedRow{}, fmt.Errorf("il valore '%s' non e' un orario di fine valido nei formati 'HH:MM', 'HH:MM:SS' o 'HH'", v)
	}

	start = time.Date(data.Year(), data.Month(), data.Day(), start.Hour(), start.Minute(), start.Second(), 0, localTimeZone)
	end = time.Date(data.Year(), data.Month(), data.Day(), end.Hour(), end.Minute(), end.Second(), 0, localTimeZone)

	if end.Before(start) {
		end = time.Date(data.Year(), data.Month(), data.Day()+1, end.Hour(), end.Minute(), end.Second(), 0, localTimeZone)
	}

	room := Room{}
	if r.Aula != "" {
		room.Code = strings.ToLower(r.Aula)
		room.Name = r.Aula
	} else {
		room.Code = ""
		room.Name = "???"
	}

	operator := Operator{}
	if r.Educatore != "" {
		operator.Code = strings.ToLower(r.Educatore)
		operator.Name = r.Educatore
	} else {
		operator.Code = ""
		operator.Name = "???"
	}

	numPaganti := 0
	numGratuiti := 0
	numAccompagnatori := 0

	if r.Paganti != "" {
		numPaganti, err = strconv.Atoi(r.Paganti)
		if err != nil {
			return ParsedRow{}, fmt.Errorf("il valore '%s' non e' un numero valido", r.Paganti)
		}
	}
	if r.Gratuiti != "" {
		numGratuiti, err = strconv.Atoi(r.Gratuiti)
		if err != nil {
			return ParsedRow{}, fmt.Errorf("il valore '%s' non e' un numero valido", r.Gratuiti)
		}
	}
	if r.Accompagnatori != "" {
		numAccompagnatori, err = strconv.Atoi(r.Accompagnatori)
		if err != nil {
			return ParsedRow{}, fmt.Errorf("il valore '%s' non e' un numero valido", r.Accompagnatori)
		}
	}

	return ParsedRow{
		Raw:               r,
		StartAt:           start,
		EndAt:             end,
		Duration:          end.Sub(start),
		Room:              room,
		Operator:          operator,
		NumPaganti:        numPaganti,
		NumGratuiti:       numGratuiti,
		NumAccompagnatori: numAccompagnatori,
	}, nil
}
