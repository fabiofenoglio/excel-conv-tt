package database

import "strings"

var (
	knownRoomMap      map[string]KnownRoom
	knownRoomAliasMap map[string]string
)

func init() {
	knownRoomMap = make(map[string]KnownRoom)
	knownRoomAliasMap = make(map[string]string)

	registerKnownRooms(KnownRoom{
		Code:            "navetta",
		Name:            "Navetta",
		BackgroundColor: "#D9E9FA",
		Slots:           5,
		PreferredOrder:  -9,
	}, KnownRoom{
		Code:            "museo",
		Name:            "Museo",
		BackgroundColor: "#C8C7F9",
		Slots:           6,
		PreferredOrder:  -8,
	}, KnownRoom{
		Code:                 "pranzo",
		Name:                 "Pranzo",
		BackgroundColor:      "#F1D3F5",
		Slots:                4,
		AllowMissingOperator: true,
		PreferredOrder:       -7,
	}, KnownRoom{
		Code:                     "planetario",
		Name:                     "Planetario",
		BackgroundColor:          "#F6F7D4",
		Slots:                    5,
		PreferredOrder:           -6,
		SlotPlacementPreferences: &planetarioSlotPlacementPreferences,
	}, KnownRoom{
		Code:            "aula 1",
		Aliases:         []string{"aula didattica 1", "lab 1", "laboratorio 1", "laboratorio didattico 1"},
		Name:            "Aula 1",
		BackgroundColor: "#C7E8B5",
		Slots:           1,
		PreferredOrder:  -5,
	}, KnownRoom{
		Code:            "aula 2",
		Aliases:         []string{"aula didattica 2", "lab 2", "laboratorio 2", "laboratorio didattico 2"},
		Name:            "Aula 2",
		BackgroundColor: "#C7E8B5",
		Slots:           1,
		PreferredOrder:  -4,
	})
}

func registerKnownRooms(o ...KnownRoom) {
	for _, obj := range o {
		if obj.Code == "" {
			panic("known room must have a code")
		}
		knownRoomMap[strings.ToLower(obj.Code)] = obj

		for _, alias := range obj.Aliases {
			knownRoomAliasMap[alias] = obj.Code
		}
	}
}

func GetKnownRoom(code string) (KnownRoom, bool) {
	if aliasOf, isAlias := knownRoomAliasMap[code]; isAlias {
		return GetKnownRoom(aliasOf)
	}

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

func GetEffectiveSlotPlacementPreferencesForRoom(roomCode string) SlotPlacementPreferences {
	knownRoom, ok := GetKnownRoom(roomCode)
	if !ok {
		return DefaultSlotPlacementPreferences
	}
	if knownRoom.SlotPlacementPreferences == nil {
		return DefaultSlotPlacementPreferences
	}
	return *knownRoom.SlotPlacementPreferences
}
