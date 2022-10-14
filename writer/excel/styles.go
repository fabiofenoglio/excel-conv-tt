package excel

import (
	"github.com/xuri/excelize/v2"
)

var (
	dayBoxStyle = &Style{
		Common: &StyleDefinition{
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
		},
	}

	dayRoomBoxStyle = &Style{
		Common: &StyleDefinition{
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
		},
	}

	dayHeaderStyle = &Style{
		Common: &StyleDefinition{
			Style: &excelize.Style{
				Font: &excelize.Font{
					Size:  14,
					Color: "#FFFFFF",
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{"#48752C"},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "#48752C", Style: 1},
					{Type: "right", Color: "#48752C", Style: 1},
					{Type: "top", Color: "#48752C", Style: 1},
					{Type: "bottom", Color: "#48752C", Style: 1},
				},
			},
		},
	}

	inBoxAnnotationOnLeft = &Style{
		Common: &StyleDefinition{
			Style: &excelize.Style{
				Font: &excelize.Font{
					Size:  11,
					Color: "#333333",
				},
				Alignment: &excelize.Alignment{
					Horizontal: "left",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{"#FFFFFF"},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "#333333", Style: 1},
					{Type: "right", Color: "#333333", Style: 1},
					{Type: "top", Color: "#333333", Style: 1},
					{Type: "bottom", Color: "#333333", Style: 1},
				},
			},
		},
	}
	inBoxAnnotationOnRight = &Style{
		Common: &StyleDefinition{
			Style: &excelize.Style{
				Font: &excelize.Font{
					Size:  11,
					Color: "#333333",
				},
				Alignment: &excelize.Alignment{
					Horizontal: "right",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{"#FFFFFF"},
					Pattern: 1,
				},
				Border: []excelize.Border{
					{Type: "left", Color: "#333333", Style: 1},
					{Type: "right", Color: "#333333", Style: 1},
					{Type: "top", Color: "#333333", Style: 1},
					{Type: "bottom", Color: "#333333", Style: 1},
				},
			},
		},
	}

	noRoomStyle = &Style{
		Common: &StyleDefinition{
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

	unusedRoomStyle = &Style{
		Common: &StyleDefinition{
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
		},
	}

	noOperatorStyle = &Style{
		Common: &StyleDefinition{
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

	noOperatorNeededStyle = &Style{
		Common: &StyleDefinition{
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
		},
	}

	toBeFilledStyle = &Style{
		Common: &StyleDefinition{
			Style: &excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{"#DD3333"},
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
)

func buildForRoom(color string) *Style {
	return &Style{
		Common: &StyleDefinition{
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

func buildForOperator(color string) *Style {
	return &Style{
		Common: &StyleDefinition{
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
		WarningVariant: &StyleDefinition{
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
