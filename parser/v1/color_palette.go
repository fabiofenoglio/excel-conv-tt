package parser

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
	registerColorRange  = 0
	assignedColorsCache map[string]string
)

func init() {
	assignedColorsCache = make(map[string]string)
}

func pickColor(key string) string {
	if v, assignedAlready := assignedColorsCache[key]; assignedAlready {
		return v
	}

	v := availableColors[registerColorRange]
	registerColorRange++
	if registerColorRange >= len(availableColors) {
		registerColorRange = 0
	}

	assignedColorsCache[key] = v
	return v
}
