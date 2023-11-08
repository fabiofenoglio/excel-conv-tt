package excel

import (
	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/xuri/excelize/v2"
)

type designatedStyling struct {
	style *RegisteredStyleV2
	box   excel.CellBox
}

func applyStyleToBox(f *excelize.File, style *RegisteredStyleV2, box excel.CellBox) error {
	sn := box.TopLeft().SheetName()
	if box.Width() == 1 && box.Height() == 1 {
		// is a single 1x1 cell
		return f.SetCellStyle(sn, box.TopLeft().Code(), box.BottomRight().Code(), style.SingleCell())
	}

	if box.Width() == 1 {
		// is a column
		if e := f.SetCellStyle(sn, box.TopLeft().Code(), box.TopLeft().Code(), style.TopLeftRight()); e != nil {
			return e
		}
		if box.Height() > 2 {
			if e := f.SetCellStyle(sn, box.TopLeft().AtBottom(1).Code(), box.BottomRight().AtTop(1).Code(), style.LeftRight()); e != nil {
				return e
			}
		}
		if e := f.SetCellStyle(sn, box.BottomRight().Code(), box.BottomRight().Code(), style.BottomLeftRight()); e != nil {
			return e
		}
		return nil
	}

	if box.Height() == 1 {
		// is a row
		if e := f.SetCellStyle(sn, box.TopLeft().Code(), box.TopLeft().Code(), style.TopLeftBottom()); e != nil {
			return e
		}
		if box.Width() > 2 {
			if e := f.SetCellStyle(sn, box.TopLeft().AtRight(1).Code(), box.BottomRight().AtLeft(1).Code(), style.TopBottom()); e != nil {
				return e
			}
		}
		if e := f.SetCellStyle(sn, box.BottomRight().Code(), box.BottomRight().Code(), style.TopRightBottom()); e != nil {
			return e
		}
		return nil
	}

	// is a NxM box with N,M > 1
	if e := f.SetCellStyle(sn, box.TopLeft().Code(), box.TopLeft().Code(), style.TopLeft()); e != nil {
		return e
	}
	if e := f.SetCellStyle(sn, box.TopRight().Code(), box.TopRight().Code(), style.TopRight()); e != nil {
		return e
	}
	if e := f.SetCellStyle(sn, box.BottomLeft().Code(), box.BottomLeft().Code(), style.BottomLeft()); e != nil {
		return e
	}
	if e := f.SetCellStyle(sn, box.BottomRight().Code(), box.BottomRight().Code(), style.BottomRight()); e != nil {
		return e
	}
	if box.Height() > 2 {
		if e := f.SetCellStyle(sn, box.TopLeft().AtBottom(1).Code(), box.BottomLeft().AtTop(1).Code(), style.Left()); e != nil {
			return e
		}
		if e := f.SetCellStyle(sn, box.TopRight().AtBottom(1).Code(), box.BottomRight().AtTop(1).Code(), style.Right()); e != nil {
			return e
		}
	}
	if box.Width() > 2 {
		if e := f.SetCellStyle(sn, box.TopLeft().AtRight(1).Code(), box.TopRight().AtLeft(1).Code(), style.Top()); e != nil {
			return e
		}
		if e := f.SetCellStyle(sn, box.BottomLeft().AtRight(1).Code(), box.BottomRight().AtLeft(1).Code(), style.Bottom()); e != nil {
			return e
		}
	}
	if box.Height() > 2 && box.Width() > 2 {
		if e := f.SetCellStyle(
			sn,
			box.TopLeft().AtRight(1).AtBottom(1).Code(),
			box.BottomRight().AtLeft(1).AtTop(1).Code(),
			style.Middle(),
		); e != nil {
			return e
		}
	}

	return nil
}
