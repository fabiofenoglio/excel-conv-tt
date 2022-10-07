package model

import (
	"time"
)

type ParsedRow struct {
	Raw ExcelRow

	StartAt           time.Time
	EndAt             time.Time
	Duration          time.Duration
	Room              Room
	Operator          Operator
	NumPaganti        int
	NumGratuiti       int
	NumAccompagnatori int

	Warnings []string
}
