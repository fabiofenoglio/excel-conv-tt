package excel

import (
	"github.com/xuri/excelize/v2"
)

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
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DDDDDD"},
			Pattern: 1,
		},
	}

	dayRoomBoxStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Top: true, Bottom: true, Left: true, Right: true,
		},
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
	}
	schoolRecapHeaderStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Font: &excelize.Font{
			Bold: true,
			Size: 10,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 1,
			Bottom: true,
		},
	}

	dayHeaderStyle = &StyleDefV2{
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
		Border: &StyleDefV2Border{
			Color: "#48752C", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
	}

	inBoxAnnotationOnLeft = &StyleDefV2{
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
		Border: &StyleDefV2Border{
			Color: "#333333", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
	}
	inBoxAnnotationOnRight = &StyleDefV2{
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
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E02222"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: "#FFFFFF",
		},
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
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#EEEEEE"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}

	noOperatorStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
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
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#EEEEEE"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Size: 10,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 1,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}

	toBeFilledStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Bottom: true,
		},
		AsWarning: standardWarningVariant,
	}

	bottomPlaceholderEmojiStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Bottom: true,
		},
	}
	bottomPlaceholderTextStyle = &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 4,
			Bottom: true,
		},
		Font: &excelize.Font{
			Bold: true,
		},
	}

	softDividerStyle = &StyleDefV2{
		Border: &StyleDefV2Border{
			Color: "333333", Style: 1,
			Bottom: true,
		},
	}
)

func buildForRoom(color string) *StyleDefV2 {
	return &StyleDefV2{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFFFF"},
			Pattern: 1,
		},
		Border: &StyleDefV2Border{
			Color: "333333", Style: 2,
			Top: true, Bottom: true, Left: true, Right: true,
		},
		AsWarning: standardWarningVariant,
	}
}

func buildForOperator(color string) *StyleDefV2 {
	return &StyleDefV2{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Fill: excelize.Fill{
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
