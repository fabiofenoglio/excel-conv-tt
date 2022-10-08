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

	knownRoom, isKnown := database.GetKnownOperator(code)

	backgroundColor := knownRoom.BackgroundColor

	if !isKnown {
		backgroundColor = pickColor("op/" + code)
	}

	return model.Operator{
		Code:            code,
		Name:            name,
		Known:           isKnown,
		BackgroundColor: backgroundColor,
	}
}
