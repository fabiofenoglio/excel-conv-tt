package model

type Room struct {
	Code string `json:"code"`
	Name string `json:"name"`

	Known                bool   `json:"is_known"`
	AllowMissingOperator bool   `json:"-"`
	BackgroundColor      string `json:"-"`
	Slots                uint   `json:"slots,omitempty"`
	PreferredOrder       int    `json:"-"`
	ShowActivityNames    bool   `json:"-"`
	Hide                 bool   `json:"-"`
}

var (
	NoRoom = Room{
		Code:                 "",
		Name:                 "???",
		Known:                false,
		AllowMissingOperator: false,
		BackgroundColor:      "",
		Slots:                0,
		PreferredOrder:       0,
	}
)
