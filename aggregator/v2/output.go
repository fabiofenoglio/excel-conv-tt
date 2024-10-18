package aggregator

import (
	"time"

	"github.com/fabiofenoglio/excelconv/parser/v2"
)

type OutputRow struct {
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

	Confirmed                   *bool
	IsPlaceholderNumeroAttivita bool

	Warnings []parser.Warning

	CompetenceDate time.Time
}

type Output struct {
	CommonData CommonData
	Days       []ScheduleForSingleDayWithRoomsAndGroupSlots
}

func ToOutputRow(input Row) OutputRow {
	return OutputRow{
		ID:           input.InputRow.ID,
		BookingCode:  input.InputRow.BookingCode,
		Date:         input.InputRow.Date,
		StartTime:    input.InputRow.StartTime,
		EndTime:      input.InputRow.EndTime,
		Duration:     input.InputRow.Duration,
		BookingNote:  input.InputRow.BookingNote,
		OperatorNote: input.InputRow.OperatorNote,
		Payment: parser.PaymentStatus{
			PaymentAdvance:       input.InputRow.Payment.PaymentAdvance,
			PaymentAdvanceStatus: input.InputRow.Payment.PaymentAdvanceStatus,
		},
		RoomCode:                    input.InputRow.RoomCode,
		OperatorCode:                input.InputRow.OperatorCode,
		VisitingGroupCode:           input.InputRow.VisitingGroupCode,
		ActivityCode:                input.InputRow.ActivityCode,
		CompetenceDate:              input.CompetenceDate,
		Bus:                         input.InputRow.Bus,
		Warnings:                    input.InputRow.Warnings,
		IsPlaceholderNumeroAttivita: input.InputRow.IsPlaceholderNumeroAttivita,
		Confirmed:                   input.InputRow.Confirmed,
	}
}
