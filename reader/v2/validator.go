package reader

import (
	"errors"
	"fmt"
	"regexp"
	_ "time/tzdata"
)

func Validate(rows []Row) error {
	for _, row := range rows {
		if err := validateRow(row); err != nil {
			return fmt.Errorf("la riga %d non è valida: %w", row.rowNumber, err)
		}
	}
	return nil
}

func validateRow(r Row) error {
	if r.BookingCode == "" {
		return errors.New("codice mancante")
	}

	timeRangeRegexp := regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])\s*\-\s*([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`)
	if !timeRangeRegexp.MatchString(r.TimesRawString) {
		return fmt.Errorf("orario non valido, atteso HH:MM-HH:MM e non '%s'", r.TimesRawString)
	}

	dateRegexp := regexp.MustCompile(`^\d{2}\/\d{2}\/\d{4}$`)
	if !dateRegexp.MatchString(r.DateRawString) {
		return fmt.Errorf("data non valida, atteso GG/MM/YYYY e non '%s'", r.DateRawString)
	}

	return nil
}
