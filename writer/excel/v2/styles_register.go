package excel

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

type RegisteredStyleV2 struct {
	styleDef *StyleDefV2
	f        *excelize.File
	warning  bool
}

type StyleRegister struct {
	f                *excelize.File
	registeredStyles map[string]*RegisteredStyleV2
}

func NewStyleRegister(f *excelize.File) *StyleRegister {
	return &StyleRegister{
		f:                f,
		registeredStyles: make(map[string]*RegisteredStyleV2),
	}
}

func (r *StyleRegister) Get(entry *StyleDefV2) *RegisteredStyleV2 {
	return r.registerIfNeeded(entry)
}

func (r *StyleRegister) registerIfNeeded(entry *StyleDefV2) *RegisteredStyleV2 {
	return &RegisteredStyleV2{
		f:        r.f,
		styleDef: entry,
		warning:  false,
	}
}

func (r *RegisteredStyleV2) WithWarning() *RegisteredStyleV2 {
	if r.warning {
		return r
	}
	return &RegisteredStyleV2{
		styleDef: r.styleDef,
		f:        r.f,
		warning:  true,
	}
}

func (r *RegisteredStyleV2) registerWithVariants(s *excelize.Style) int {
	if r.warning && r.styleDef.AsWarning != nil {
		r.styleDef.AsWarning(s)
	}

	styleID, err := r.f.NewStyle(s)
	if err != nil {
		panic(err)
	}
	return styleID
}

func (r *RegisteredStyleV2) Middle() int {
	return r.registerWithVariants(r.styleDef.Middle())
}
func (r *RegisteredStyleV2) Top() int {
	return r.registerWithVariants(r.styleDef.Top())
}
func (r *RegisteredStyleV2) Bottom() int {
	return r.registerWithVariants(r.styleDef.Bottom())
}
func (r *RegisteredStyleV2) Left() int {
	return r.registerWithVariants(r.styleDef.Left())
}
func (r *RegisteredStyleV2) Right() int {
	return r.registerWithVariants(r.styleDef.Right())
}
func (r *RegisteredStyleV2) TopLeft() int {
	return r.registerWithVariants(r.styleDef.TopLeft())
}
func (r *RegisteredStyleV2) BottomLeft() int {
	return r.registerWithVariants(r.styleDef.BottomLeft())
}
func (r *RegisteredStyleV2) TopRight() int {
	return r.registerWithVariants(r.styleDef.TopRight())
}
func (r *RegisteredStyleV2) BottomRight() int {
	return r.registerWithVariants(r.styleDef.BottomRight())
}
func (r *RegisteredStyleV2) TopLeftRight() int {
	return r.registerWithVariants(r.styleDef.TopLeftRight())
}
func (r *RegisteredStyleV2) BottomLeftRight() int {
	return r.registerWithVariants(r.styleDef.BottomLeftRight())
}
func (r *RegisteredStyleV2) TopLeftBottom() int {
	return r.registerWithVariants(r.styleDef.TopLeftBottom())
}
func (r *RegisteredStyleV2) TopRightBottom() int {
	return r.registerWithVariants(r.styleDef.TopRightBottom())
}
func (r *RegisteredStyleV2) LeftRight() int {
	return r.registerWithVariants(r.styleDef.LeftRight())
}
func (r *RegisteredStyleV2) TopBottom() int {
	return r.registerWithVariants(r.styleDef.TopBottom())
}
func (r *RegisteredStyleV2) SingleCell() int {
	return r.registerWithVariants(r.styleDef.SingleCell())
}

func (r *StyleRegister) NoOperatorNeededStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(noOperatorNeededStyle)
}

func (r *StyleRegister) DayBoxStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(dayBoxStyle)
}

func (r *StyleRegister) DayRoomBoxStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(dayRoomBoxStyle)
}

func (r *StyleRegister) DayHeaderStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(dayHeaderStyle)
}

func (r *StyleRegister) InBoxAnnotationOnLeftStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(inBoxAnnotationOnLeft)
}

func (r *StyleRegister) InBoxAnnotationOnRightStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(inBoxAnnotationOnRight)
}

func (r *StyleRegister) ToBeFilledStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(toBeFilledStyle)
}

func (r *StyleRegister) UnusedRoomStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(unusedRoomStyle)
}

func (r *StyleRegister) NoRoomStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(noRoomStyle)
}

func (r *StyleRegister) NoOperatorStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(noOperatorStyle)
}

func (r *StyleRegister) SchoolRecapStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(schoolRecapStyle)
}

func (r *StyleRegister) SchoolRecapNotesStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(schoolRecapNotesStyle)
}

func (r *StyleRegister) SchoolRecapContactStyle() *RegisteredStyleV2 {
	return r.registerIfNeeded(schoolRecapContactStyle)
}

func (r *StyleRegister) OperatorStyle(color string) *RegisteredStyleV2 {
	key := "op/" + strings.ToLower(color)
	if v, ok := r.registeredStyles[key]; ok {
		return v
	}

	style := buildForOperator(color)
	reg := r.registerIfNeeded(style)

	r.registeredStyles[key] = reg

	return reg
}

func (r *StyleRegister) RoomStyle(color string) *RegisteredStyleV2 {
	key := "room/" + strings.ToLower(color)
	if v, ok := r.registeredStyles[key]; ok {
		return v
	}

	style := buildForRoom(color)
	reg := r.registerIfNeeded(style)

	r.registeredStyles[key] = reg

	return reg
}
