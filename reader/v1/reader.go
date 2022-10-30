package reader

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/sirupsen/logrus"
)

func ReadFromFile(input string, log *logrus.Logger) ([]ExcelRow, error) {
	var err error

	f, err := excelize.OpenFile(input)
	if err != nil {
		return nil, errors.Wrap(err, "error opening input file")
	}

	defer func() {
		// Close the spreadsheet.
		if closeErr := f.Close(); closeErr != nil {
			log.Errorf("error closing input file: %s", closeErr.Error())
		}
	}()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		sheetName = "Organizzazione"
		log.Warnf("reverting to default source sheet name '%s'", sheetName)
	}

	startingHeaderCell := excel.NewCell(sheetName, 1, 4)

	columnNameToFieldNameMap := make(map[string]string)
	sampleRow := &ExcelRow{}
	val := reflect.ValueOf(sampleRow).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		columnName := strings.ToLower(strings.TrimSpace(field.Tag.Get("column")))
		if columnName != "" {
			columnNameToFieldNameMap[columnName] = field.Name
		}
	}

	currentHeaderCell := startingHeaderCell
	headers := make([]string, 0, 10)
	fieldNameToColumnNumberMap := make(map[string]uint)

	for {
		cell, err := f.GetCellValue(currentHeaderCell.SheetName(), currentHeaderCell.Code())
		if err != nil {
			return nil, errors.Wrapf(err, "error reading header cell %s", currentHeaderCell.Code())
		}
		cellCode := strings.ToLower(strings.TrimSpace(cell))
		if cellCode == "" {
			break
		}

		headers = append(headers, cell)
		if len(headers) > 5000 {
			return nil, errors.New("too many headers found")
		}

		if mapsToFieldName, ok := columnNameToFieldNameMap[cellCode]; ok {

			if _, isDuplicated := fieldNameToColumnNumberMap[mapsToFieldName]; isDuplicated {
				log.Warnf("multiple mapping columns for field '%s'", mapsToFieldName)

			} else {
				log.Debugf("column %s with header '%s' maps to known field '%s'", currentHeaderCell.ColumnName(), cell, mapsToFieldName)
				fieldNameToColumnNumberMap[mapsToFieldName] = currentHeaderCell.Column()
			}

		} else {
			log.Warnf("column %s with header '%s' does not map to any known field", currentHeaderCell.ColumnName(), cell)
		}

		currentHeaderCell.MoveRight(1)
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
				log.Debugf("field '%s' is mapped by column '%d'", field.Name, mappedByColumnNumer)
			} else {
				errMsg := fmt.Sprintf("field '%s' is NOT mapped by any column", field.Name)
				if required {
					log.Error(errMsg)
					return nil, errors.New(errMsg)
				} else {
					log.Warn(errMsg)
				}
			}
		}
	}

	log.Debug("mapping build completed")
	log.Debug("all required fields have a valid mapping column")

	results := make([]ExcelRow, 0, 20)

	currentCell := startingHeaderCell.AtBottom(1)

	for {
		row := ExcelRow{}
		anyNonNil := false

		for fieldName, columnNumber := range fieldNameToColumnNumberMap {
			cell := currentCell.AtColumn(columnNumber)

			cellContent, err := f.GetCellValue(cell.SheetName(), cell.Code())
			if err != nil {
				return nil, errors.Wrapf(err, "error reading content cell %v", cell)
			}

			trimmed := strings.TrimSpace(cellContent)
			if trimmed != cellContent {
				log.Warnf("cell %s has a value that starts or finishes with spaces or newline, this should be avoided",
					cell.Code())
			}
			cellContent = trimmed

			if cellContent != "" {
				anyNonNil = true
				reflect.ValueOf(&row).Elem().FieldByName(fieldName).SetString(cellContent)

				log.Debugf("setting row %d.%s to '%s' by cell %s", len(results), fieldName, cellContent, cell.Code())
			}
		}

		if !anyNonNil {
			break
		}

		results = append(results, row)

		currentCell.MoveBottom(1)
	}

	return results, nil
}
