package excel

import (
	"github.com/xuri/excelize/v2"
)

type StyleDefV2 struct {
	Alignment *excelize.Alignment
	Fill      *excelize.Fill
	Border    *StyleDefV2Border
	Font      *excelize.Font
	AsWarning func(s *excelize.Style)
}

type StyleDefV2Border struct {
	Color  string `json:"color"`
	Style  int    `json:"style"`
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
}

func (d *StyleDefV2) base() *excelize.Style {
	fill := excelize.Fill{}
	if d.Fill != nil {
		fill = *d.Fill
	}

	return &excelize.Style{
		Alignment: d.Alignment,
		Fill:      fill,
		Font:      d.Font,
	}
}

func (d *StyleDefV2) Middle() *excelize.Style {
	return d.base()
}

func (d *StyleDefV2) Top() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 1)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) Bottom() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 1)
	if d.Border != nil {
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) Left() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 1)
	if d.Border != nil {
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) Right() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 1)
	if d.Border != nil {
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) TopLeft() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 2)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) TopRight() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 2)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) BottomLeft() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 2)
	if d.Border != nil {
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) BottomRight() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 2)
	if d.Border != nil {
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) TopLeftBottom() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 4)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) TopRightBottom() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 4)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) TopLeftRight() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 4)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) BottomLeftRight() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 4)
	if d.Border != nil {
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) LeftRight() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 2)
	if d.Border != nil {
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) TopBottom() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 2)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}

func (d *StyleDefV2) SingleCell() *excelize.Style {
	o := d.base()
	borders := make([]excelize.Border, 0, 4)
	if d.Border != nil {
		if d.Border.Top {
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Left {
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Bottom {
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
		if d.Border.Right {
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: d.Border.Color,
				Style: d.Border.Style,
			})
		}
	}
	o.Border = borders
	return o
}
