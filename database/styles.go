package database

import (
	"github.com/xuri/excelize/v2"
)

var availableColors = []string{
	"#D9E9FA",
	"#C8C7F9",
	"#F1D3F5",
	"#F6F7D4",
	"#C7E8B5",
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

	noRoomStyle = Styles{
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
	}

	unusedRoomStyle = StyleEntry{
		Style: &excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#EEEEEE"},
				Pattern: 1,
			},
			Border: []excelize.Border{
				{Type: "left", Color: "333333", Style: 4},
				{Type: "right", Color: "333333", Style: 4},
				{Type: "top", Color: "333333", Style: 4},
				{Type: "bottom", Color: "333333", Style: 4},
			},
		},
	}

	noOperatorStyle = Styles{
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
					Size:  10,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "FFFFFF", Style: 1},
					{Type: "right", Color: "FFFFFF", Style: 1},
					{Type: "top", Color: "FFFFFF", Style: 4},
					{Type: "bottom", Color: "FFFFFF", Style: 4},
				},
			},
		},
	}

	noOperatorNeededStyle = StyleEntry{
		Style: &excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#EEEEEE"},
				Pattern: 1,
			},
			Font: &excelize.Font{
				Size: 10,
			},
			Border: []excelize.Border{
				{Type: "left", Color: "EEEEEE", Style: 1},
				{Type: "right", Color: "EEEEEE", Style: 1},
				{Type: "top", Color: "EEEEEE", Style: 4},
				{Type: "bottom", Color: "EEEEEE", Style: 4},
			},
		},
	}

	toBeFilledStyle = StyleEntry{
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
	}
)

func buildForRoom(color string) Styles {
	return Styles{
		Common: StyleEntry{
			Style: &excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{color},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: color, Style: 2},
					{Type: "right", Color: color, Style: 2},
					{Type: "top", Color: color, Style: 1},
					{Type: "bottom", Color: color, Style: 1},
				},
			},
		},
	}
}

func buildForOperator(color string) Styles {
	return Styles{
		Common: StyleEntry{
			Style: &excelize.Style{
				Font: &excelize.Font{
					Size: 10,
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{color},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: color, Style: 1},
					{Type: "right", Color: color, Style: 1},
					{Type: "top", Color: color, Style: 4},
					{Type: "bottom", Color: color, Style: 4},
				},
			},
		},
		Warning: StyleEntry{
			Style: &excelize.Style{
				Font: &excelize.Font{
					Size: 10,
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{color},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "#CC2222", Style: 5},
					{Type: "right", Color: "#CC2222", Style: 5},
					{Type: "top", Color: "#CC2222", Style: 1},
					{Type: "bottom", Color: "#CC2222", Style: 1},
				},
			},
		},
	}
}
