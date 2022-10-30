package reader

import (
	"regexp"
	_ "time/tzdata"

	"github.com/pkg/errors"
)

func Validate(rows []Row) error {
	for _, row := range rows {
		if err := validateRow(row); err != nil {
			return errors.Wrapf(err, "la riga %d non è valida", row.rowNumber)
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
		return errors.Errorf("orario non valido, atteso HH:MM-HH:MM e non '%s'", r.TimesRawString)
	}

	dateRegexp := regexp.MustCompile(`^\d{2}\/\d{2}\/\d{4}$`)
	if !dateRegexp.MatchString(r.DateRawString) {
		return errors.Errorf("data non valida, atteso GG/MM/YYYY e non '%s'", r.DateRawString)
	}

	return nil
}
