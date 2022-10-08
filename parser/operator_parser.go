package parser

import (
	"strings"

	"github.com/fabiofenoglio/excelconv/database"
	"github.com/fabiofenoglio/excelconv/model"
)

func operatorFromName(name string) model.Operator {
	code := strings.ToLower(strings.TrimSpace(name))

	if code == "" {
		return model.NoOperator
	}

	knownOperator, isKnown := database.GetKnownOperator(code)

	backgroundColor := knownOperator.BackgroundColor

	if !isKnown {
		backgroundColor = pickColor("op/" + code)
	} else {
		if len(knownOperator.Name) > 0 {
			name = knownOperator.Name
		}
	}

	return model.Operator{
		Code:            code,
		Name:            name,
		Known:           isKnown,
		BackgroundColor: backgroundColor,
	}
}
