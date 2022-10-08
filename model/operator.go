package model

type Operator struct {
	Code string `json:"code"`
	Name string `json:"name"`

	Known           bool   `json:"is_known"`
	BackgroundColor string `json:"-"`
}

var (
	NoOperator = Operator{
		Code:            "",
		Name:            "???",
		Known:           false,
		BackgroundColor: "",
	}
)
