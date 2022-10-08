package database

import "strings"

var (
	knownRoomMap map[string]KnownRoom
)

func init() {
	knownRoomMap = make(map[string]KnownRoom)

	registerKnownRooms(KnownRoom{
		Code:            "navetta",
		Name:            "Navetta",
		BackgroundColor: "#D9E9FA",
		Slots:           5,
		PreferredOrder:  -5,
	}, KnownRoom{
		Code:            "museo",
		Name:            "Museo",
		BackgroundColor: "#C8C7F9",
		Slots:           6,
		PreferredOrder:  -4,
	}, KnownRoom{
		Code:                 "pranzo",
		Name:                 "Pranzo",
		BackgroundColor:      "#F1D3F5",
		Slots:                4,
		AllowMissingOperator: true,
		PreferredOrder:       -3,
	}, KnownRoom{
		Code:            "planetario",
		Name:            "Planetario",
		BackgroundColor: "#F6F7D4",
		Slots:           5,
		PreferredOrder:  -2,
	}, KnownRoom{
		Code:            "lab",
		Name:            "Lab",
		BackgroundColor: "#C7E8B5",
		Slots:           2,
		PreferredOrder:  -1,
	})
}

func registerKnownRooms(o ...KnownRoom) {
	for _, obj := range o {
		if obj.Code == "" {
			panic("known room must have a code")
		}
		knownRoomMap[strings.ToLower(obj.Code)] = obj
	}
}

func GetKnownRoom(code string) (KnownRoom, bool) {
	res, ok := knownRoomMap[code]
	return res, ok
}

func GetKnownRooms() []KnownRoom {
	out := make([]KnownRoom, 0, len(knownRoomMap))
	for _, v := range knownRoomMap {
		out = append(out, v)
	}
	return out
}
