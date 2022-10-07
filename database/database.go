package database

import (
	"github.com/xuri/excelize/v2"
)

type StyleEntry struct {
	Style   *excelize.Style
	StyleID int
}

type Styles struct {
	Common  StyleEntry
	Warning StyleEntry
}

type KnownRoom struct {
	Styles
	Name                 string
	Slots                uint
	AllowMissingOperator bool
}

type KnownOperator struct {
	Styles
}

var (
	registerTarget     *excelize.File
	registerColorRange = 0

	knownRoomMap     map[string]*KnownRoom
	knownOperatorMap map[string]*KnownOperator
)

func KnownRoomMap() map[string]*KnownRoom {
	return knownRoomMap
}

func KnownOperatorMap() map[string]*KnownOperator {
	return knownOperatorMap
}

func init() {
	// https://xuri.me/excelize/en/style.html#border

	knownRoomMap = make(map[string]*KnownRoom)
	knownRoomMap[""] = &KnownRoom{
		Styles: noRoomStyle,
	}

	knownRoomMap["navetta"] = &KnownRoom{
		Name:   "Navetta",
		Styles: buildForRoom("#D9E9FA"),
		Slots:  5,
	}
	knownRoomMap["museo"] = &KnownRoom{
		Name:   "Museo",
		Styles: buildForRoom("#C8C7F9"),
		Slots:  6,
	}
	knownRoomMap["pranzo"] = &KnownRoom{
		Name:                 "Pranzo",
		Styles:               buildForRoom("#F1D3F5"),
		Slots:                4,
		AllowMissingOperator: true,
	}
	knownRoomMap["planetario"] = &KnownRoom{
		Name:   "Planetario",
		Styles: buildForRoom("#F6F7D4"),
		Slots:  5,
	}
	knownRoomMap["lab"] = &KnownRoom{
		Name:   "Lab",
		Styles: buildForRoom("#C7E8B5"),
		Slots:  2,
	}

	knownOperatorMap = make(map[string]*KnownOperator)
	knownOperatorMap[""] = &KnownOperator{
		Styles: noOperatorStyle,
	}

	knownOperatorMap["a"] = &KnownOperator{
		Styles: buildForOperator("#E0EBF5"),
	}

}

func NoOperatorNeededStyleID() int {
	if noOperatorNeededStyle.StyleID == 0 {
		registerStyleEntry(&noOperatorNeededStyle, registerTarget)
	}
	return noOperatorNeededStyle.StyleID
}

func DayBoxStyleID() int {
	if dayBoxStyle.StyleID == 0 {
		registerStyleEntry(&dayBoxStyle, registerTarget)
	}
	return dayBoxStyle.StyleID
}

func DayRoomBoxStyleID() int {
	if dayRoomBoxStyle.StyleID == 0 {
		registerStyleEntry(&dayRoomBoxStyle, registerTarget)
	}
	return dayRoomBoxStyle.StyleID
}

func DayHeaderStyleID() int {
	if dayHeaderStyle.StyleID == 0 {
		registerStyleEntry(&dayHeaderStyle, registerTarget)
	}
	return dayHeaderStyle.StyleID
}

func ToBeFilledStyleID() int {
	if toBeFilledStyle.StyleID == 0 {
		registerStyleEntry(&toBeFilledStyle, registerTarget)
	}
	return toBeFilledStyle.StyleID
}

func UnusedRoomStyleID() int {
	if unusedRoomStyle.StyleID == 0 {
		registerStyleEntry(&unusedRoomStyle, registerTarget)
	}
	return unusedRoomStyle.StyleID
}

func RegisterStyles(f *excelize.File) {
	registerTarget = f
}

func GetEntryForOperator(code string) KnownOperator {
	if entry, ok := knownOperatorMap[code]; ok {
		if entry.Common.StyleID == 0 && registerTarget != nil {
			registerStyleEntry(&entry.Common, registerTarget)
			registerStyleEntry(&entry.Warning, registerTarget)
		}
		return *entry
	}

	style := buildForOperator(pickColor())
	if registerTarget != nil {
		registerStyleEntry(&style.Common, registerTarget)
		registerStyleEntry(&style.Warning, registerTarget)
	}
	knownOperatorMap[code] = &KnownOperator{
		Styles: style,
	}

	return *knownOperatorMap[code]
}

func GetEntryForRoom(code string) KnownRoom {
	if entry, ok := knownRoomMap[code]; ok {
		if entry.Common.StyleID == 0 && registerTarget != nil {
			registerStyleEntry(&entry.Common, registerTarget)
			registerStyleEntry(&entry.Warning, registerTarget)
		}
		return *entry
	}

	style := buildForRoom(pickColor())
	if registerTarget != nil {
		registerStyleEntry(&style.Common, registerTarget)
	}

	knownRoomMap[code] = &KnownRoom{
		Styles: style,
	}

	return *knownRoomMap[code]
}

func pickColor() string {
	v := availableColors[registerColorRange]
	registerColorRange++
	if registerColorRange >= len(availableColors) {
		registerColorRange = 0
	}
	return v
}

func registerStyleEntry(e *StyleEntry, f *excelize.File) {
	if e == nil || e.Style == nil {
		return
	}
	if e.StyleID > 0 {
		return
	}
	styleID, err := f.NewStyle(e.Style)
	if err != nil {
		panic(err)
	}
	e.StyleID = styleID
}
