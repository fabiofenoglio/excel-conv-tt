package model

import "fmt"

type CellWithTracker interface {
	Cell
	CoveredArea() CellBox
	MoveAtBottomLeftOfCoveredArea() CellWithTracker
	MoveAtRightTopOfCoveredArea() CellWithTracker
}

type CellWithTrackerImpl struct {
	cell Cell

	minColumn uint
	minRow    uint
	maxColumn uint
	maxRow    uint

	parent   *CellWithTrackerImpl
	children []*CellWithTrackerImpl
}

func NewCellWithTracker(cell Cell) *CellWithTrackerImpl {
	return &CellWithTrackerImpl{
		cell:      cell,
		minColumn: cell.Column(),
		minRow:    cell.Row(),
		maxColumn: cell.Column(),
		maxRow:    cell.Row(),
	}
}

var _ Cell = &CellWithTrackerImpl{}
var _ CellWithTracker = &CellWithTrackerImpl{}

func (c *CellWithTrackerImpl) SheetName() string {
	return c.cell.SheetName()
}

func (c *CellWithTrackerImpl) Row() uint {
	return c.cell.Row()
}

func (c *CellWithTrackerImpl) Column() uint {
	return c.cell.Column()
}

func (c *CellWithTrackerImpl) ColumnName() string {
	return c.cell.ColumnName()
}

func (c *CellWithTrackerImpl) Code() string {
	return c.cell.Code()
}

func (c *CellWithTrackerImpl) Copy() Cell {
	return c.copy()
}

func (c *CellWithTrackerImpl) MoveLeft(cells uint) Cell {
	c.cell.MoveLeft(cells)
	c.trackNewCol(c.cell.Column())
	return c
}

func (c *CellWithTrackerImpl) MoveRight(cells uint) Cell {
	c.cell.MoveRight(cells)
	c.trackNewCol(c.cell.Column())
	return c
}

func (c *CellWithTrackerImpl) MoveTop(cells uint) Cell {
	c.cell.MoveTop(cells)
	c.trackNewRow(c.cell.Row())
	return c
}

func (c *CellWithTrackerImpl) MoveBottom(cells uint) Cell {
	c.cell.MoveBottom(cells)
	c.trackNewRow(c.cell.Row())
	return c
}

func (c *CellWithTrackerImpl) MoveColumn(to uint) Cell {
	c.cell.MoveColumn(to)
	c.trackNewCol(c.cell.Column())
	return c
}

func (c *CellWithTrackerImpl) MoveRow(to uint) Cell {
	c.cell.MoveRow(to)
	c.trackNewRow(c.cell.Row())
	return c
}

func (c *CellWithTrackerImpl) AtLeft(cells uint) Cell {
	return c.copy().MoveLeft(cells)
}

func (c *CellWithTrackerImpl) AtRight(cells uint) Cell {
	return c.copy().MoveRight(cells)
}

func (c *CellWithTrackerImpl) AtTop(cells uint) Cell {
	return c.copy().MoveTop(cells)
}

func (c *CellWithTrackerImpl) AtBottom(cells uint) Cell {
	return c.copy().MoveBottom(cells)
}

func (c *CellWithTrackerImpl) AtColumn(to uint) Cell {
	return c.copy().MoveColumn(to)
}

func (c *CellWithTrackerImpl) AtRow(to uint) Cell {
	return c.copy().MoveRow(to)
}

func (c *CellWithTrackerImpl) MoveAtBottomLeftOfCoveredArea() CellWithTracker {
	ca := c.coveredArea()
	c.MoveRow(ca.BottomRow() + 1).MoveColumn(ca.LeftColumn())
	return c
}

func (c *CellWithTrackerImpl) MoveAtRightTopOfCoveredArea() CellWithTracker {
	ca := c.coveredArea()
	c.MoveRow(ca.TopRow()).MoveColumn(ca.RightColumn() + 1)
	return c
}

func (c *CellWithTrackerImpl) CoveredArea() CellBox {
	return c.coveredArea()
}

func (c *CellWithTrackerImpl) coveredArea() CellBox {
	topLeft := c.cell.Copy().AtColumn(c.minColumn).AtRow(c.minRow)
	bottomRight := c.cell.Copy().AtColumn(c.maxColumn).AtRow(c.maxRow)

	for _, children := range c.children {
		childrenCellBox := children.coveredArea()
		if childrenCellBox.TopLeft().Column() < topLeft.Column() {
			topLeft.MoveColumn(childrenCellBox.TopLeft().Column())
		}
		if childrenCellBox.TopLeft().Row() < topLeft.Row() {
			topLeft.MoveRow(childrenCellBox.TopLeft().Row())
		}
		if childrenCellBox.BottomRight().Column() > bottomRight.Column() {
			bottomRight.MoveColumn(childrenCellBox.BottomRight().Column())
		}
		if childrenCellBox.BottomRight().Row() > bottomRight.Row() {
			bottomRight.MoveRow(childrenCellBox.BottomRight().Row())
		}
	}

	cellBox := &CellBoxImpl{
		topLeft:     topLeft,
		bottomRight: bottomRight,
	}

	return cellBox
}

func (c *CellWithTrackerImpl) copy() *CellWithTrackerImpl {
	cp := &CellWithTrackerImpl{
		cell:      c.cell.Copy(),
		minColumn: c.minColumn,
		minRow:    c.minRow,
		maxColumn: c.maxColumn,
		maxRow:    c.maxRow,
		parent:    c,
	}

	c.children = append(c.children, cp)

	return cp
}

func (c *CellWithTrackerImpl) trackNewCol(col uint) *CellWithTrackerImpl {
	if col < c.minColumn {
		c.minColumn = col
	}
	if col > c.maxColumn {
		c.maxColumn = col
	}
	if c.parent != nil {
		c.parent.trackNewCol(col)
	}
	return c
}

func (c *CellWithTrackerImpl) trackNewRow(row uint) *CellWithTrackerImpl {
	if row < c.minRow {
		c.minRow = row
	}
	if row > c.maxRow {
		c.maxRow = row
	}
	if c.parent != nil {
		c.parent.trackNewRow(row)
	}
	return c
}

func (c *CellWithTrackerImpl) String() string {
	return fmt.Sprintf("%s*", c.cell.String())
}
