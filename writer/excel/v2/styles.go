package excel

import (
	"github.com/xuri/excelize/v2"
)

const defaultFontSize = float64(11)

type FontOverride struct {
	SizeOffset   float64
	ForceSize    float64
	ForceBold    bool
	ForceNotBold bool
	Color        string
}

func defaultFontBuilder(settings *FontOverride) *excelize.Font {
	out := &excelize.Font{
		Size: defaultFontSize,
		Bold: false,
	}

	if settings == nil {
		return out
	}

	if settings.ForceSize != 0 {
		out.Size = settings.ForceSize
	}
	if settings.SizeOffset != 0 {
		out.Size += settings.SizeOffset
	}
	if settings.ForceBold {
		out.Bold = true
	}
	if settings.ForceNotBold {
		out.Bold = false
	}
	if settings.Color != "" {
		out.Color = settings.Color
	}

	return out
}

var (
	standardWarningVariant = func(s *excelize.Style) {
		if s.Border != nil {
			for i := range s.Border {
				s.Border[i].Color = "#CC2222"
				s.Border[i].Style = 5
			}
		}
	}

	dayBoxStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DDDDDD"},
			Pattern: 1,
		},
		Font: defaultFontBuilder(nil),
	}

	dayRoomBoxStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold: true,
			ForceSize: 10,
		}),
	}

	schoolRecapStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Bottom: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold:  true,
			SizeOffset: 3,
		}),
	}
	schoolRecapNotesStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Border: &StyleDefV2Border{
			Color:  "#CC2222",
			Style:  1,
			Bottom: true,
			Left:   true,
			Top:    true,
			Right:  true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold:  true,
			SizeOffset: 0,
		}),
	}
	schoolRecapContactStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Border: &StyleDefV2Border{
			Color:  "333333",
			Style:  4,
			Bottom: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold:  false,
			SizeOffset: -1,
		}),
	}
	schoolRecapHeaderStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Font: defaultFontBuilder(nil),
		Border: &StyleDefV2Border{
			Color:  "333333",
			Style:  1,
			Bottom: true,
		},
	}
	highlightForSpecialProjectStyle = &StyleDefV2{
		Border: &StyleDefV2Border{
			Color:  "#9900cc",
			Style:  5,
			Bottom: true,
			Left:   true,
			Top:    true,
			Right:  true,
		},
	}
	highlightForSpecialNotesStyle = &StyleDefV2{
		Border: &StyleDefV2Border{
			Color:  "#0066ff",
			Style:  5,
			Bottom: true,
			Left:   true,
			Top:    true,
			Right:  true,
		},
	}
	highlightForUnconfirmedStyle = &StyleDefV2{
		Font: defaultFontBuilder(&FontOverride{
			Color: "#aaaaaa",
		}),
	}

	dayHeaderStyle = &StyleDefV2{
		Font: defaultFontBuilder(&FontOverride{
			SizeOffset: 3,
			Color:      "#FFFFFF",
		}),
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#48752C"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "#48752C", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
	}

	inBoxAnnotationOnLeft = &StyleDefV2{
		Font: defaultFontBuilder(&FontOverride{
			SizeOffset: 0,
			Color:      "#333333",
		}),
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFFFF"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "#333333", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
	}
	inBoxAnnotationOnRight = &StyleDefV2{
		Font: defaultFontBuilder(&FontOverride{
			SizeOffset: 11,
			Color:      "#333333",
		}),
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFFFF"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "#333333", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
	}

	noRoomStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E02222"},
			Pattern: 1,
		},
		Font: defaultFontBuilder(&FontOverride{
			Color:     "#FFFFFF",
			ForceBold: true,
			ForceSize: 10,
		}),
		Border: &StyleDefV2Border{
			Color: "FFFFFF", Style: 2,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}

	unusedRoomStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#EEEEEE"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold: true,
			ForceSize: 10,
		}),
		AsWarning: standardWarningVariant,
	}

	noOperatorStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E02222"},
			Pattern: 1,
		},
		Font: defaultFontBuilder(&FontOverride{
			Color:     "#FFFFFF",
			ForceBold: true,
			ForceSize: 10,
		}),
		Border: &StyleDefV2Border{
			Color: "FFFFFF", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}

	noOperatorNeededStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#EEEEEE"},
			Pattern: 1,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold: true,
			ForceSize: 10,
		}),
		Border: &StyleDefV2Border{
			Color: "333333", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}

	bottomPlaceholderEmojiStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: &StyleDefV2Border{
			Color:  "333333",
			Style:  1,
			Bottom: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold:  true,
			SizeOffset: 2,
		}),
	}
	bottomPlaceholderTextStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: &StyleDefV2Border{
			Color:  "333333",
			Style:  1,
			Bottom: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold:  true,
			SizeOffset: 1,
		}),
	}
	toBeFilledStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: &StyleDefV2Border{
			Color:  "333333",
			Style:  1,
			Bottom: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceBold:  false,
			SizeOffset: 2,
		}),
		AsWarning: standardWarningVariant,
	}

	softDividerStyle = &StyleDefV2{
		Border: &StyleDefV2Border{
			Color:  "333333",
			Style:  1,
			Bottom: true,
		},
		Font: defaultFontBuilder(nil),
	}
)

func buildForRoom(color string) *StyleDefV2 {
	return &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFFFF"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 2,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		Font: defaultFontBuilder(&FontOverride{
			ForceSize: 10,
		}),
		AsWarning: standardWarningVariant,
	}
}

func buildForOperator(color string) *StyleDefV2 {
	return &StyleDefV2{
		Font: defaultFontBuilder(&FontOverride{
			ForceSize: 10,
		}),
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Fill: &excelize.Fill{
			Type:    "pattern",
			Color:   []string{color},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "555555", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}
}
