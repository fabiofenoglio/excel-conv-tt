package aggregator

import (
	"time"

	"github.com/fabiofenoglio/excelconv/parser/v2"
)

type InputRow struct {
	ID int

	BookingCode  string
	Date         time.Time
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	BookingNote  string
	OperatorNote string
	Payment      parser.PaymentStatus
	Bus          string

	RoomCode          string
	OperatorCode      string
	VisitingGroupCode string
	ActivityCode      string

	Warnings []parser.Warning
}

type Input struct {
	Anagraphics *parser.OutputAnagraphics
	Rows        []InputRow
}

func ToInput(input parser.Output) Input {
	i := Input{
		Anagraphics: input.Anagraphics,
		Rows:        make([]InputRow, 0, len(input.Rows)),
	}

	for _, r := range input.Rows {
		i.Rows = append(i.Rows, InputRow{
			ID:           r.ID,
			BookingCode:  r.BookingCode,
			Date:         r.Date,
			StartTime:    r.StartTime,
			EndTime:      r.EndTime,
			Duration:     r.Duration,
			BookingNote:  r.BookingNote,
			OperatorNote: r.OperatorNote,
			Payment: parser.PaymentStatus{
				PaymentAdvance:       r.Payment.PaymentAdvance,
				PaymentAdvanceStatus: r.Payment.PaymentAdvanceStatus,
			},
			RoomCode:          r.RoomCode,
			OperatorCode:      r.OperatorCode,
			VisitingGroupCode: r.VisitingGroupCode,
			ActivityCode:      r.ActivityCode,
			Bus:               r.Bus,
			Warnings:          r.Warnings,
		})
	}

	return i
}
