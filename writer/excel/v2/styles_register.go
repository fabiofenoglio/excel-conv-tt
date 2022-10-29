package excel

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

type Style struct {
	Common         *StyleDefinition
	WarningVariant *StyleDefinition
}

type StyleDefinition struct {
	Style   *excelize.Style
	StyleID int
}

type StyleRegister struct {
	f                *excelize.File
	registeredStyles map[string]*Style
}

func (s *Style) WarningIDOrDefault() int {
	if s.WarningVariant != nil {
		return s.WarningVariant.StyleID
	}
	return s.Common.StyleID
}

func (s *Style) MatchesID(id int) bool {
	if s.WarningVariant != nil && s.WarningVariant.StyleID == id {
		return true
	}
	return s.Common.StyleID == id
}

func NewStyleRegister(f *excelize.File) *StyleRegister {
	return &StyleRegister{
		f:                f,
		registeredStyles: make(map[string]*Style),
	}
}

func (r *StyleRegister) registerIfNeeded(entry *Style) *Style {
	if entry.Common != nil {
		r.registerDefinitionIfNeeded(entry.Common)
	}
	if entry.WarningVariant != nil {
		r.registerDefinitionIfNeeded(entry.WarningVariant)
	}
	return entry
}

func (r *StyleRegister) registerDefinitionIfNeeded(e *StyleDefinition) int {
	if e == nil || e.Style == nil {
		return 0
	}
	if e.StyleID > 0 {
		return e.StyleID
	}
	styleID, err := r.f.NewStyle(e.Style)
	if err != nil {
		panic(err)
	}
	e.StyleID = styleID
	return styleID
}

func (r *StyleRegister) NoOperatorNeededStyle() *Style {
	return r.registerIfNeeded(noOperatorNeededStyle)
}

func (r *StyleRegister) DayBoxStyle() *Style {
	return r.registerIfNeeded(dayBoxStyle)
}

func (r *StyleRegister) DayRoomBoxStyle() *Style {
	return r.registerIfNeeded(dayRoomBoxStyle)
}

func (r *StyleRegister) DayHeaderStyle() *Style {
	return r.registerIfNeeded(dayHeaderStyle)
}

func (r *StyleRegister) InBoxAnnotationOnLeftStyle() *Style {
	return r.registerIfNeeded(inBoxAnnotationOnLeft)
}

func (r *StyleRegister) InBoxAnnotationOnRightStyle() *Style {
	return r.registerIfNeeded(inBoxAnnotationOnRight)
}

func (r *StyleRegister) ToBeFilledStyle() *Style {
	return r.registerIfNeeded(toBeFilledStyle)
}

func (r *StyleRegister) UnusedRoomStyle() *Style {
	return r.registerIfNeeded(unusedRoomStyle)
}

func (r *StyleRegister) NoRoomStyle() *Style {
	return r.registerIfNeeded(noRoomStyle)
}

func (r *StyleRegister) NoOperatorStyle() *Style {
	return r.registerIfNeeded(noOperatorStyle)
}

func (r *StyleRegister) SchoolRecapStyle() *Style {
	return r.registerIfNeeded(schoolRecapStyle)
}

func (r *StyleRegister) OperatorStyle(color string) *Style {
	key := "op/" + strings.ToLower(color)
	if v, ok := r.registeredStyles[key]; ok {
		return v
	}

	style := buildForOperator(color)
	r.registerIfNeeded(style)

	r.registeredStyles[key] = style

	return style
}

func (r *StyleRegister) RoomStyle(color string) *Style {
	key := "room/" + strings.ToLower(color)
	if v, ok := r.registeredStyles[key]; ok {
		return v
	}

	style := buildForRoom(color)
	r.registerIfNeeded(style)

	r.registeredStyles[key] = style

	return style
}
