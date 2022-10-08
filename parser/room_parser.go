package parser

import (
	"strings"

	"github.com/fabiofenoglio/excelconv/database"
	"github.com/fabiofenoglio/excelconv/model"
)

func roomFromName(name string) model.Room {
	code := strings.ToLower(strings.TrimSpace(name))

	if code == "" {
		return model.NoRoom
	}

	knownRoom, isKnown := database.GetKnownRoom(code)

	backgroundColor := knownRoom.BackgroundColor

	if !isKnown {
		backgroundColor = pickColor("room/" + code)
	}

	return model.Room{
		Code:                 code,
		Name:                 name,
		Known:                isKnown,
		AllowMissingOperator: knownRoom.AllowMissingOperator,
		BackgroundColor:      backgroundColor,
		Slots:                knownRoom.Slots,
		PreferredOrder:       knownRoom.PreferredOrder,
	}
}
