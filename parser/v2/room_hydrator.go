package parser

import (
	"strings"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/database"
)

func HydrateRooms(_ config.WorkflowContext, rows []InputRow) ([]Row, []Room, error) {
	outRows := make([]Row, 0, len(rows))
	outRooms := make([]Room, 0, 10)

	roomsIndex := make(map[string]Room)

	for _, row := range rows {
		mappedRow := Row{
			InputRow: row,
			RoomCode: "",
		}

		roomCode := nameToCode(row.roomRawString)

		if roomCode != "" {
			if alreadyMapped, ok := roomsIndex[roomCode]; ok {
				mappedRow.RoomCode = alreadyMapped.Code
			} else {
				newRoom := buildNewRoom(roomCode, strings.TrimSpace(row.roomRawString))
				mappedRow.RoomCode = newRoom.Code

				if _, alreadyMappedAsAlias := roomsIndex[newRoom.Code]; !alreadyMappedAsAlias {
					roomsIndex[newRoom.Code] = newRoom
					outRooms = append(outRooms, newRoom)
				}
			}
		}

		outRows = append(outRows, mappedRow)
	}

	// add known rooms that were not referenced
	for _, room := range database.GetKnownRooms() {
		if _, ok := roomsIndex[room.Code]; !ok {
			fromName := buildNewRoom(room.Code, room.Name)
			roomsIndex[fromName.Code] = fromName
			outRooms = append(outRooms, fromName)
		}
	}

	return outRows, outRooms, nil
}

func buildNewRoom(code string, name string) Room {

	knownRoom, isKnown := database.GetKnownRoom(code)

	if isKnown {
		if len(knownRoom.Name) > 0 {
			name = knownRoom.Name
		}
		code = knownRoom.Code
	}

	return Room{
		Code:                           code,
		Name:                           cleanStringForVisualization(name),
		Known:                          isKnown,
		AllowMissingOperator:           knownRoom.AllowMissingOperator,
		BackgroundColor:                "",
		Slots:                          knownRoom.Slots,
		PreferredOrder:                 knownRoom.PreferredOrder,
		ShowActivityNamesAsAnnotations: knownRoom.ShowActivityNamesAsAnnotations,
		Hide:                           knownRoom.Hide,
		GroupActivities:                knownRoom.GroupActivities,
		ShowActivityNamesInside:        knownRoom.ShowActivityNamesInside,
	}
}
