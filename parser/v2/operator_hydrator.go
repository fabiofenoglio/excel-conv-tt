package parser

import (
	"strings"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/database"
)

func HydrateOperators(_ config.WorkflowContext, rows []Row) ([]Row, []Operator, error) {
	outRows := make([]Row, 0, len(rows))
	outOperators := make([]Operator, 0, 10)

	operatorsIndex := make(map[string]Operator)

	for _, row := range rows {
		mappedRow := row

		operatorCode := nameToCode(row.operatorRawString)

		if operatorCode != "" {
			if alreadyMapped, ok := operatorsIndex[operatorCode]; ok {
				mappedRow.OperatorCode = alreadyMapped.Code
			} else {
				newOperator := buildNewOperator(operatorCode, strings.TrimSpace(row.operatorRawString))
				mappedRow.OperatorCode = newOperator.Code

				if _, alreadyMappedAsAlias := operatorsIndex[newOperator.Code]; !alreadyMappedAsAlias {
					operatorsIndex[newOperator.Code] = newOperator
					outOperators = append(outOperators, newOperator)
				}
			}
		}

		outRows = append(outRows, mappedRow)
	}

	// add known opeartors that were not referenced
	for _, operator := range database.GetKnownOperators() {
		if _, ok := operatorsIndex[operator.Code]; !ok {
			fromName := buildNewOperator(operator.Code, operator.Name)
			operatorsIndex[fromName.Code] = fromName
			outOperators = append(outOperators, fromName)
		}
	}

	return outRows, outOperators, nil
}

func buildNewOperator(code string, name string) Operator {

	knownOperator, isKnown := database.GetKnownOperator(code)

	backgroundColor := knownOperator.BackgroundColor

	if !isKnown {
		backgroundColor = pickColor("op/" + code)
	} else {
		if len(knownOperator.Name) > 0 {
			name = knownOperator.Name
		}
		code = knownOperator.Code
	}

	return Operator{
		Code:            code,
		Name:            cleanStringForVisualization(name),
		Known:           isKnown,
		BackgroundColor: backgroundColor,
	}
}
