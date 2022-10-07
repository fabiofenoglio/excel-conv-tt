package model

import "fmt"

type Cell struct {
	c         uint
	r         uint
	sheetName string
}

func NewCell(sheetName string, c, r uint) Cell {
	return Cell{
		c:         c,
		r:         r,
		sheetName: sheetName,
	}
}

func (c Cell) SheetName() string {
	return c.sheetName
}
func (c Cell) Column() uint {
	return c.c
}
func (c Cell) Row() uint {
	return c.r
}
func (c Cell) ColumnName() string {
	return toColumnName(c.c)
}
func (c Cell) Code() string {
	return fmt.Sprintf("%s%d", toColumnName(c.c), c.r)
}
func (c Cell) String() string {
	if c.sheetName != "" {
		return fmt.Sprintf("!%s:%s", c.sheetName, c.Code())
	}
	return c.Code()
}
func (c Cell) Copy() Cell {
	return Cell{
		c:         c.c,
		r:         c.r,
		sheetName: c.sheetName,
	}
}
func (c Cell) AtLeft(cells uint) Cell {
	cp := c.Copy()
	cp.c -= cells
	return cp
}
func (c Cell) AtRight(cells uint) Cell {
	cp := c.Copy()
	cp.c += cells
	return cp
}
func (c Cell) AtTop(cells uint) Cell {
	cp := c.Copy()
	cp.r -= cells
	return cp
}
func (c Cell) AtBottom(cells uint) Cell {
	cp := c.Copy()
	cp.r += cells
	return cp
}
func (c Cell) AtColumn(to uint) Cell {
	cp := c.Copy()
	cp.c = to
	return cp
}
func (c Cell) AtRow(to uint) Cell {
	cp := c.Copy()
	cp.r = to
	return cp
}
func (c *Cell) MoveLeft(cells uint) {
	c.c -= cells
}
func (c *Cell) MoveRight(cells uint) {
	c.c += cells
}
func (c *Cell) MoveTop(cells uint) {
	c.r -= cells
}
func (c *Cell) MoveBottom(cells uint) {
	c.r += cells
}
func (c *Cell) MoveColumn(to uint) {
	c.c = to
}
func (c *Cell) MoveRow(to uint) {
	c.r = to
}
