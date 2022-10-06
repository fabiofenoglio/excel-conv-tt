package services

import "github.com/xuri/excelize/v2"

type StyleEntry struct {
	Style   *excelize.Style
	StyleID int
}

type Styles struct {
	Common StyleEntry
}

type KnownRoom struct {
	Styles
}

type KnownOperator struct {
	Styles
}

var availableColors = []string{
	"#7fb574",
	"#8c6a3e",
	"#a16052",
	"#4c6f9c",
	"#65508a",
	"#4bada0",
	"#c7b88f",
	"#a15da0",
	"#d9c7d0",
}

var (
	registerTarget     *excelize.File
	registerColorRange = 0

	knownRoomMap     map[string]*KnownRoom
	knownOperatorMap map[string]*KnownOperator

	dayBoxStyle = StyleEntry{
		Style: &excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#DDDDDD"},
				Pattern: 1,
			},
			Border: []excelize.Border{},
		},
	}
	dayRoomBoxStyle = StyleEntry{
		Style: &excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{},
			Border: []excelize.Border{
				{Type: "left", Color: "333333", Style: 4},
				{Type: "right", Color: "333333", Style: 4},
				{Type: "top", Color: "333333", Style: 4},
				{Type: "bottom", Color: "333333", Style: 4},
			},
		},
	}
	dayHeaderStyle = StyleEntry{
		Style: &excelize.Style{
			Font: &excelize.Font{
				Size: 16,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#947521"},
				Pattern: 1,
			},
			Border: []excelize.Border{
				{Type: "left", Color: "333333", Style: 1},
				{Type: "right", Color: "333333", Style: 1},
				{Type: "top", Color: "333333", Style: 1},
				{Type: "bottom", Color: "333333", Style: 1},
			},
		},
	}
	verticalTextStyle = StyleEntry{
		Style: &excelize.Style{
			Font: &excelize.Font{
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				TextRotation: 90,
			},
		},
	}
)

func commonRoomStyle() Styles {
	return Styles{
		Common: StyleEntry{
			Style: &excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{"#E0EBF5"},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "0000FF", Style: 2},
					{Type: "right", Color: "0000FF", Style: 2},
					{Type: "top", Color: "0000FF", Style: 1},
					{Type: "bottom", Color: "0000FF", Style: 1},
				},
			},
		},
	}
}

func commonOperatorStyle() Styles {
	return Styles{
		Common: StyleEntry{
			Style: &excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{"#E0EBF5"},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "0000FF", Style: 1},
					{Type: "right", Color: "0000FF", Style: 1},
					{Type: "top", Color: "0000FF", Style: 4},
					{Type: "bottom", Color: "0000FF", Style: 4},
				},
			},
		},
	}
}

func init() {
	// https://xuri.me/excelize/en/style.html#border

	knownRoomMap = make(map[string]*KnownRoom)
	knownRoomMap["planetario"] = &KnownRoom{
		Styles: commonRoomStyle(),
	}
	knownRoomMap[""] = &KnownRoom{
		Styles: Styles{
			Common: StyleEntry{
				Style: &excelize.Style{
					Alignment: &excelize.Alignment{
						Horizontal: "center",
						Vertical:   "center",
					},
					Fill: excelize.Fill{
						Type:    "pattern",
						Color:   []string{"#E02222"},
						Pattern: 1,
					},
					Font: &excelize.Font{
						Color: "#FFFFFF",
					},
					Border: []excelize.Border{
						{Type: "left", Color: "FFFFFF", Style: 2},
						{Type: "right", Color: "FFFFFF", Style: 2},
						{Type: "top", Color: "FFFFFF", Style: 1},
						{Type: "bottom", Color: "FFFFFF", Style: 1},
					},
				},
			},
		},
	}

	knownOperatorMap = make(map[string]*KnownOperator)
	knownOperatorMap["a"] = &KnownOperator{
		Styles: commonOperatorStyle(),
	}
	knownOperatorMap[""] = &KnownOperator{
		Styles: Styles{
			Common: StyleEntry{
				Style: &excelize.Style{
					Alignment: &excelize.Alignment{
						Horizontal: "center",
						Vertical:   "center",
					},
					Fill: excelize.Fill{
						Type:    "pattern",
						Color:   []string{"#E02222"},
						Pattern: 1,
					},
					Font: &excelize.Font{
						Color: "#FFFFFF",
					},
					Border: []excelize.Border{
						{Type: "left", Color: "FFFFFF", Style: 1},
						{Type: "right", Color: "FFFFFF", Style: 1},
						{Type: "top", Color: "FFFFFF", Style: 4},
						{Type: "bottom", Color: "FFFFFF", Style: 4},
					},
				},
			},
		},
	}
}
func registerStyleEntry(e *StyleEntry, f *excelize.File) {
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
	}

	for _, entry := range knownOperatorMap {
		registerStyleEntry(&entry.Common, f)
	}

	registerStyleEntry(&dayBoxStyle, f)
	registerStyleEntry(&dayRoomBoxStyle, f)
	registerStyleEntry(&dayHeaderStyle, f)
	registerStyleEntry(&verticalTextStyle, f)
}

func GetEntryForOperator(code string) KnownOperator {
	if entry, ok := knownOperatorMap[code]; ok {
		return *entry
	}

	style := commonOperatorStyle()
	style.Common.Style.Fill.Color = []string{pickColor()}
	registerStyleEntry(&style.Common, registerTarget)

	knownOperatorMap[code] = &KnownOperator{
		Styles: style,
	}

	return *knownOperatorMap[code]
}

func GetEntryForRoom(code string) KnownRoom {
	if entry, ok := knownRoomMap[code]; ok {
		return *entry
	}

	style := commonRoomStyle()
	style.Common.Style.Fill.Color = []string{pickColor()}
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
