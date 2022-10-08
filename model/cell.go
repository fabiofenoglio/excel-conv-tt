package model

import "fmt"

type Cell interface {
	SheetName() string
	Row() uint
	Column() uint
	ColumnName() string
	Code() string
	Copy() Cell
	AtLeft(cells uint) Cell
	AtRight(cells uint) Cell
	AtTop(cells uint) Cell
	AtBottom(cells uint) Cell
	AtColumn(to uint) Cell
	AtRow(to uint) Cell
	MoveLeft(cells uint) Cell
	MoveRight(cells uint) Cell
	MoveTop(cells uint) Cell
	MoveBottom(cells uint) Cell
	MoveColumn(to uint) Cell
	MoveRow(to uint) Cell

	String() string
}

type CellImpl struct {
	c         uint
	r         uint
	sheetName string
}

var _ Cell = &CellImpl{}

func NewCell(sheetName string, c, r uint) *CellImpl {
	return &CellImpl{
		c:         c,
		r:         r,
		sheetName: sheetName,
	}
}

func (c *CellImpl) SheetName() string {
	return c.sheetName
}
func (c *CellImpl) Column() uint {
	return c.c
}
func (c *CellImpl) Row() uint {
	return c.r
}
func (c *CellImpl) ColumnName() string {
	return toColumnName(c.c)
}
func (c *CellImpl) Code() string {
	return fmt.Sprintf("%s%d", toColumnName(c.c), c.r)
}

func (c *CellImpl) Copy() Cell {
	return c.copy()
}

func (c *CellImpl) MoveLeft(cells uint) Cell {
	c.c -= cells
	return c
}
func (c *CellImpl) MoveRight(cells uint) Cell {
	c.c += cells
	return c
}
func (c *CellImpl) MoveTop(cells uint) Cell {
	c.r -= cells
	return c
}
func (c *CellImpl) MoveBottom(cells uint) Cell {
	c.r += cells
	return c
}
func (c *CellImpl) MoveColumn(to uint) Cell {
	c.c = to
	return c
}
func (c *CellImpl) MoveRow(to uint) Cell {
	c.r = to
	return c
}

func (c *CellImpl) AtLeft(cells uint) Cell {
	return c.copy().MoveLeft(cells)
}
func (c *CellImpl) AtRight(cells uint) Cell {
	return c.copy().MoveRight(cells)
}
func (c *CellImpl) AtTop(cells uint) Cell {
	return c.copy().MoveTop(cells)
}
func (c *CellImpl) AtBottom(cells uint) Cell {
	return c.copy().MoveBottom(cells)
}
func (c *CellImpl) AtColumn(to uint) Cell {
	return c.copy().MoveColumn(to)
}
func (c *CellImpl) AtRow(to uint) Cell {
	return c.copy().MoveRow(to)
}

func (c *CellImpl) String() string {
	if c.sheetName != "" {
		return fmt.Sprintf("!%s:%s", c.sheetName, c.Code())
	}
	return c.Code()
}

func (c *CellImpl) copy() *CellImpl {
	return &CellImpl{
		c:         c.c,
		r:         c.r,
		sheetName: c.sheetName,
	}
}
