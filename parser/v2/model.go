package parser

import (
	"strings"
)

type HighlightReason string

var (
	HighlightSpecialProject HighlightReason = "special_project"
	HighlightSpecialNotes   HighlightReason = "special_notes"
)

type Room struct {
	Code                           string `json:"code"`
	Name                           string `json:"name"`
	Known                          bool   `json:"is_known"`
	AllowMissingOperator           bool   `json:"-"`
	BackgroundColor                string `json:"-"`
	Slots                          uint   `json:"slots,omitempty"`
	PreferredOrder                 int    `json:"-"`
	ShowActivityNamesAsAnnotations bool   `json:"-"`
	Hide                           bool   `json:"-"`
	GroupActivities                bool   `json:"-"`
	ShowActivityNamesInside        bool   `json:"-"`
	AlwaysShow                     bool   `json:"-"`
}

type Row struct {
	InputRow
	RoomCode          string
	OperatorCode      string
	VisitingGroupCode string
	ActivityCode      string
	Warnings          []Warning
}

type Operator struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Known           bool   `json:"is_known"`
	BackgroundColor string `json:"-"`
}

type Activity struct {
	Code     string `json:"code"`
	TypeCode string `json:"type_code"`
	Name     string `json:"name"`
	Language string `json:"lang"`
}

type ActivityType struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VisitingGroup struct {
	Code                string
	SchoolCode          string
	SchoolClassCode     string
	Composition         GroupComposition
	ClassTeacher        string
	ClassRefEmail       string
	BookingNotes        string
	OperatorNotes       string
	SpecialProjectNotes string
	Highlights          []HighlightReason
}

type School struct {
	Code string `json:"code"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (s School) FullDescription() string {
	out := ""
	if s.Type != "" {
		out += s.Type + " "
	}
	if s.Name != "" {
		out += s.Name
	}
	return strings.TrimSpace(out)
}

func (s School) SortableIdentifier() string {
	out := ""
	if s.Type != "" {
		out += strings.ToLower(strings.TrimSpace(s.Type)) + "/"
	}
	if s.Name != "" {
		out += strings.ToLower(strings.TrimSpace(s.Name)) + "/"
	}
	if s.Code != "" {
		out += strings.ToLower(strings.TrimSpace(s.Code))
	}
	return out
}

type SchoolClass struct {
	Code       string `json:"code"`
	SchoolCode string `json:"school_code"`
	Number     string `json:"number"`
	Section    string `json:"section"`
}

func (s SchoolClass) FullDescription() string {
	out := ""
	if s.Number != "" {
		out += s.Number + " "
	}
	if s.Section != "" {
		out += s.Section
	}
	return strings.TrimSpace(out)
}

func (s SchoolClass) SortableIdentifier() string {
	out := ""
	if s.Number != "" {
		out += strings.ToLower(strings.TrimSpace(s.Number)) + "/"
	}
	if s.Section != "" {
		out += strings.ToLower(strings.TrimSpace(s.Section)) + "/"
	}
	if s.Code != "" {
		out += strings.ToLower(strings.TrimSpace(s.Code))
	}
	return out
}

type GroupComposition struct {
	NumPaying       int `json:"num_paying"`
	NumFree         int `json:"num_free"`
	NumAccompanying int `json:"num_accompanying"`
}

func (c GroupComposition) NumTotal() int {
	return c.NumPaying + c.NumFree + c.NumAccompanying
}

type PaymentStatus struct {
	PaymentAdvance       string `json:"advance"`
	PaymentAdvanceStatus string `json:"advance_status"`
}

type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
