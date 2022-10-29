package parser

import "time"

type OutputAnagraphics struct {
	Rooms          map[string]Room
	Operators      map[string]Operator
	VisitingGroups map[string]VisitingGroup
	Schools        map[string]School
	SchoolClasses  map[string]SchoolClass
	Activities     map[string]Activity
	ActivityTypes  map[string]ActivityType
}

type Output struct {
	Anagraphics *OutputAnagraphics
	Rows        []OutputRow
}

type OutputRow struct {
	ID int

	BookingCode  string
	Date         time.Time
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	BookingNote  string
	OperatorNote string
	Payment      PaymentStatus
	Bus          string

	RoomCode          string
	OperatorCode      string
	VisitingGroupCode string
	ActivityCode      string

	Warnings []Warning

	anagraphicsRef *OutputAnagraphics
}

func ToOutputRows(rows []Row, anagraphicsRef *OutputAnagraphics) []OutputRow {
	out := make([]OutputRow, 0, len(rows))

	for _, input := range rows {
		out = append(out, OutputRow{
			ID:                input.ID,
			BookingCode:       input.BookingCode,
			Date:              input.Date,
			StartTime:         input.StartTime,
			EndTime:           input.EndTime,
			Duration:          input.Duration,
			BookingNote:       input.BookingNote,
			OperatorNote:      input.OperatorNote,
			RoomCode:          input.RoomCode,
			OperatorCode:      input.OperatorCode,
			VisitingGroupCode: input.VisitingGroupCode,
			ActivityCode:      input.ActivityCode,
			Payment: PaymentStatus{
				PaymentAdvance:       input.PaymentAdvance,
				PaymentAdvanceStatus: input.PaymentAdvanceStatus,
			},
			Warnings:       input.Warnings,
			Bus:            input.Bus,
			anagraphicsRef: anagraphicsRef,
		})
	}

	return out
}

func (r *OutputRow) Room() Room {
	return r.anagraphicsRef.Rooms[r.RoomCode]
}
func (r *OutputRow) Operator() Operator {
	return r.anagraphicsRef.Operators[r.OperatorCode]
}
func (r *OutputRow) Activity() Activity {
	return r.anagraphicsRef.Activities[r.ActivityCode]
}
func (r *OutputRow) ActivityType() ActivityType {
	act := r.Activity()
	return r.anagraphicsRef.ActivityTypes[act.TypeCode]
}
func (r *OutputRow) VisitingGroup() VisitingGroup {
	return r.anagraphicsRef.VisitingGroups[r.VisitingGroupCode]
}
func (r *OutputRow) School() School {
	group := r.VisitingGroup()
	return r.anagraphicsRef.Schools[group.SchoolCode]
}
func (r *OutputRow) SchoolClass() SchoolClass {
	group := r.VisitingGroup()
	return r.anagraphicsRef.SchoolClasses[group.SchoolClassCode]
}
