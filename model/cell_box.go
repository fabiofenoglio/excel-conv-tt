package model

type CellBox struct {
	topLeft     Cell
	bottomRight Cell
}

func NewCellBox(topLeft, bottomRight Cell) CellBox {
	return CellBox{
		topLeft:     topLeft,
		bottomRight: bottomRight,
	}
}

func (c *CellBox) Width() uint {
	return c.bottomRight.Column() - c.topLeft.Column() + 1
}

func (c *CellBox) Height() uint {
	return c.bottomRight.Row() - c.topLeft.Row() + 1
}

func (c *CellBox) TopLeft() Cell {
	return c.topLeft
}

func (c *CellBox) TopRight() Cell {
	return c.topLeft.AtRight(c.Width() - 1)
}

func (c *CellBox) BottomLeft() Cell {
	return c.bottomRight.AtLeft(c.Width() - 1)
}

func (c *CellBox) BottomRight() Cell {
	return c.bottomRight
}
