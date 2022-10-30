package parser

import (
	"time"

	"github.com/fabiofenoglio/excelconv/reader/v2"
)

type Input struct {
	Rows []InputRow
}

type InputRow struct {
	ID int

	BookingCode          string
	Date                 time.Time
	StartTime            time.Time
	EndTime              time.Time
	Duration             time.Duration
	BookingNote          string
	OperatorNote         string
	Bus                  string
	PaymentAdvance       string
	PaymentAdvanceStatus string

	// campi che verranno processati prima di essere esposti

	schoolType      string
	schoolName      string
	class           string
	classSection    string
	classTeacher    string
	classRefEmail   string
	numPaying       int
	numFree         int
	numAccompanying int

	activityRawString         string
	activityLanguageRawString string
	activityTypeRawString     string
	operatorRawString         string
	roomRawString             string
}

func ToInputRows(rows []reader.OutputRow) []InputRow {
	out := make([]InputRow, 0, len(rows))

	for _, input := range rows {
		out = append(out, InputRow{
			ID:                        input.ID,
			BookingCode:               input.BookingCode,
			Date:                      input.Date,
			operatorRawString:         input.Operator,
			roomRawString:             input.Room,
			activityRawString:         input.Activity,
			StartTime:                 input.StartTime,
			EndTime:                   input.EndTime,
			Duration:                  input.Duration,
			numPaying:                 input.NumPaying,
			numFree:                   input.NumFree,
			numAccompanying:           input.NumAccompanying,
			activityLanguageRawString: input.ActivityLanguage,
			BookingNote:               input.BookingNote,
			OperatorNote:              input.OperatorNote,
			schoolType:                input.SchoolType,
			schoolName:                input.SchoolName,
			class:                     input.Class,
			classSection:              input.ClassSection,
			classTeacher:              input.ClassTeacher,
			classRefEmail:             input.ClassRefEmail,
			activityTypeRawString:     input.Type,
			Bus:                       input.Bus,
			PaymentAdvance:            input.PaymentAdvance,
			PaymentAdvanceStatus:      input.PaymentAdvanceStatus,
		})
	}

	return out
}
