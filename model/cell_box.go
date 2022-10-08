package model

import "fmt"

type CellBox interface {
	TopLeft() Cell
	TopRight() Cell
	BottomLeft() Cell
	BottomRight() Cell

	TopRow() uint
	BottomRow() uint
	LeftColumn() uint
	RightColumn() uint

	String() string
}

type CellBoxImpl struct {
	topLeft     Cell
	bottomRight Cell
}

var _ CellBox = &CellBoxImpl{}

func NewCellBox(topLeft, bottomRight Cell) *CellBoxImpl {
	return &CellBoxImpl{
		topLeft:     NewCell(topLeft.SheetName(), topLeft.Column(), topLeft.Row()),
		bottomRight: NewCell(bottomRight.SheetName(), bottomRight.Column(), bottomRight.Row()),
	}
}

func (c *CellBoxImpl) Width() uint {
	return c.bottomRight.Column() - c.topLeft.Column() + 1
}

func (c *CellBoxImpl) Height() uint {
	return c.bottomRight.Row() - c.topLeft.Row() + 1
}

func (c *CellBoxImpl) TopLeft() Cell {
	return c.topLeft
}

func (c *CellBoxImpl) TopRight() Cell {
	return c.topLeft.AtRight(c.Width() - 1)
}

func (c *CellBoxImpl) BottomLeft() Cell {
	return c.bottomRight.AtLeft(c.Width() - 1)
}

func (c *CellBoxImpl) BottomRight() Cell {
	return c.bottomRight
}

func (c *CellBoxImpl) TopRow() uint {
	return c.topLeft.Row()
}

func (c *CellBoxImpl) BottomRow() uint {
	return c.bottomRight.Row()
}

func (c *CellBoxImpl) LeftColumn() uint {
	return c.topLeft.Column()
}

func (c *CellBoxImpl) RightColumn() uint {
	return c.bottomRight.Column()
}

func (c *CellBoxImpl) String() string {
	return fmt.Sprintf("[%v - %v]", c.TopLeft().String(), c.BottomRight().String())
}
