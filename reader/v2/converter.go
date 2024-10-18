package reader

import (
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/pkg/errors"
)

const (
	layoutTimeOnlyWithSeconds = "15:04:05"
	layoutTimeOnlyWithMinutes = "15:04"
	layoutTimeOnlyWithHours   = "15"
	layoutDateOnlyInITFormat  = "02/01/2006"
)

func Convert(rows []Row) ([]Row, error) {
	out := make([]Row, 0, len(rows))

	for _, row := range rows {
		converted, err := convertRow(row)
		if err != nil {
			return nil, errors.Wrapf(err, "errore nella riga %d", row.rowNumber)
		}
		out = append(out, converted)
	}

	return out, nil
}

func convertRow(r Row) (Row, error) {

	dataHalfDay, start, end, err := getStartAndEndTimes(r)
	if err != nil {
		return Row{}, err
	}

	r.Date = dataHalfDay
	r.StartTime = start
	r.EndTime = end
	r.Duration = r.EndTime.Sub(r.StartTime)

	if r.NumPayingRawString != "" {
		r.NumPaying, err = strconv.Atoi(r.NumPayingRawString)
		if err != nil {
			return Row{}, errors.Wrapf(err, "il valore del numero paganti '%s' non e' una numero valido", r.NumPayingRawString)
		}
	}
	if r.NumAccompanyingRawString != "" {
		r.NumAccompanying, err = strconv.Atoi(r.NumAccompanyingRawString)
		if err != nil {
			return Row{}, errors.Wrapf(err, "il valore del numero accompagnatori '%s' non e' una numero valido", r.NumAccompanyingRawString)
		}
	}
	if r.NumFreeRawString != "" {
		r.NumFree, err = strconv.Atoi(r.NumFreeRawString)
		if err != nil {
			return Row{}, errors.Wrapf(err, "il valore del numero gratuiti '%s' non e' una numero valido", r.NumFreeRawString)
		}
	}

	if strings.TrimSpace(r.ConfirmedRawString) != "" {
		confirmed, err := parseLocalizedBoolean(r.ConfirmedRawString)
		if err != nil {
			return Row{}, errors.Wrapf(err, "il valore del campo 'confermata' ('%s') non e' valido", r.ConfirmedRawString)
		}
		r.Confirmed = confirmed
	}

	return r, nil
}

func getStartAndEndTimes(r Row) (time.Time, time.Time, time.Time, error) {
	localTimeZone := config.TimeZone()

	data, err := time.Parse(layoutDateOnlyInITFormat, r.DateRawString)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, errors.Wrapf(err, "il valore '%s' non e' una data valida nel formato 'GG/MM/YYYY'", r.DateRawString)
	}

	dataHalfDay := time.Date(data.Year(), data.Month(), data.Day(), 12, 0, 0, 0, localTimeZone)

	orari := strings.Split(r.TimesRawString, "-")

	v := strings.TrimSpace(orari[0])
	start, err := time.Parse(layoutTimeOnlyWithMinutes, v)
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithSeconds, v)
	}
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithHours, v)
	}
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, errors.Wrapf(err,
			"il valore '%s' non e' un orario di inizio valido nei formati 'HH:MM', 'HH:MM:SS' o 'HH'", v)
	}

	v = strings.TrimSpace(orari[1])
	end, err := time.Parse(layoutTimeOnlyWithMinutes, v)
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithSeconds, v)
	}
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithHours, v)
	}
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, errors.Wrapf(err,
			"il valore '%s' non e' un orario di fine valido nei formati 'HH:MM', 'HH:MM:SS' o 'HH'", v)
	}

	start = time.Date(data.Year(), data.Month(), data.Day(), start.Hour(), start.Minute(), start.Second(), 0, localTimeZone)
	end = time.Date(data.Year(), data.Month(), data.Day(), end.Hour(), end.Minute(), end.Second(), 0, localTimeZone)

	if end.Before(start) {
		end = time.Date(data.Year(), data.Month(), data.Day()+1, end.Hour(), end.Minute(), end.Second(), 0, localTimeZone)
	}

	return dataHalfDay, start, end, nil
}

func parseLocalizedBoolean(raw string) (*bool, error) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return nil, nil
	}

	if strings.HasPrefix(raw, "s") {
		v := true
		return &v, nil
	}
	if strings.HasPrefix(raw, "n") {
		v := false
		return &v, nil
	}

	return nil, errors.Errorf("il valore '%s' non è un flag booleano valido", raw)
}
