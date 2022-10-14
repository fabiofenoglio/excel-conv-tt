package model

import (
	"time"
)

type ParsedRow struct {
	ID int

	Code string `json:"code"`

	StartAt  time.Time     `json:"start_at"`
	EndAt    time.Time     `json:"end_at"`
	Duration time.Duration `json:"duration"`
	Room     Room          `json:"room"`
	Operator Operator      `json:"operator"`

	Activity         Activity         `json:"activity"`
	School           School           `json:"school"`
	SchoolClass      SchoolClass      `json:"school_class"`
	GroupComposition GroupComposition `json:"group_composition"`

	TimeString    string        `json:"time_as_string"`
	DateString    string        `json:"date_as_string"`
	BookingNotes  string        `json:"booking_notes"`
	OperatorNotes string        `json:"operator_notes"`
	Bus           string        `json:"bus"`
	PaymentStatus PaymentStatus `json:"payment_status"`

	Warnings []Warning `json:"warnings"`
}
