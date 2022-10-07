package services

import "github.com/xuri/excelize/v2"

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
	Name  string
	Slots uint
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

func init() {
	// https://xuri.me/excelize/en/style.html#border

	knownRoomMap = make(map[string]*KnownRoom)
	knownRoomMap[""] = &KnownRoom{
		Styles: noRoomStyle,
	}

	knownRoomMap["planetario"] = &KnownRoom{
		Name:   "Planetario",
		Styles: buildForRoom("#E0EBF5"),
		Slots:  4,
	}
	knownRoomMap["aula 2"] = &KnownRoom{
		Name:   "Aula 2",
		Styles: buildForRoom("#33CC33"),
		Slots:  2,
	}
	knownRoomMap["aula inutile"] = &KnownRoom{
		Name:   "Aula inutile",
		Styles: buildForRoom("#222222"),
		Slots:  3,
	}

	knownOperatorMap = make(map[string]*KnownOperator)
	knownOperatorMap[""] = &KnownOperator{
		Styles: noOperatorStyle,
	}

	knownOperatorMap["a"] = &KnownOperator{
		Styles: buildForOperator("#E0EBF5"),
	}

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

func RegisterStyles(f *excelize.File) {
	registerTarget = f

	for _, entry := range knownRoomMap {
		registerStyleEntry(&entry.Common, f)
		registerStyleEntry(&entry.Warning, f)
	}

	for _, entry := range knownOperatorMap {
		registerStyleEntry(&entry.Common, f)
		registerStyleEntry(&entry.Warning, f)
	}

	registerStyleEntry(&dayBoxStyle, f)
	registerStyleEntry(&dayRoomBoxStyle, f)
	registerStyleEntry(&dayHeaderStyle, f)
	registerStyleEntry(&toBeFilledStyle, f)
}

func GetEntryForOperator(code string) KnownOperator {
	if entry, ok := knownOperatorMap[code]; ok {
		return *entry
	}

	style := buildForOperator(pickColor())
	registerStyleEntry(&style.Common, registerTarget)
	registerStyleEntry(&style.Warning, registerTarget)

	knownOperatorMap[code] = &KnownOperator{
		Styles: style,
	}

	return *knownOperatorMap[code]
}

func GetEntryForRoom(code string) KnownRoom {
	if entry, ok := knownRoomMap[code]; ok {
		return *entry
	}

	style := buildForRoom(pickColor())
	registerStyleEntry(&style.Common, registerTarget)

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
