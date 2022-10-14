package parser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/database"
	"github.com/fabiofenoglio/excelconv/model"
)

const (
	layoutTimeOnlyWithSeconds = "15:04:05"
	layoutTimeOnlyWithMinutes = "15:04"
	layoutTimeOnlyWithHours   = "15"
	layoutDateOnlyInITFormat  = "02/01/2006"
)

func Parse(input string, args config.Args, log *logrus.Logger) (model.ParsedData, error) {
	zero := model.ParsedData{}

	excelRows, err := ReadFromFile(input, log)
	if err != nil {
		return zero, fmt.Errorf("error reading from file: %w", err)
	}

	results := make([]model.ParsedRow, 0, len(excelRows))

	for i, row := range excelRows {
		err := validate(row)
		if err != nil {
			errMsg := fmt.Sprintf("error validating row %d", i)
			log.WithError(err).Errorf(errMsg+": %s", err.Error())
			return zero, fmt.Errorf(errMsg+": %w", err)
		}

		parsed, err := parseRow(row, args)
		if err != nil {
			errMsg := fmt.Sprintf("error parsing row %d", i)
			log.WithError(err).Errorf(errMsg+": %s", err.Error())
			return zero, fmt.Errorf(errMsg+": %w", err)
		}

		results = append(results, parsed)
	}

	roomsMap := make(map[string]model.Room)
	operatorsMap := make(map[string]model.Operator)

	for _, a := range results {
		if a.Room != model.NoRoom {
			roomsMap[a.Room.Code] = a.Room
		}
		if a.Operator != model.NoOperator {
			operatorsMap[a.Operator.Code] = a.Operator
		}
	}

	for _, room := range database.GetKnownRooms() {
		fromName := roomFromName(room.Code)
		if _, ok := roomsMap[fromName.Code]; !ok {
			roomsMap[fromName.Code] = fromName
		}
	}
	for _, operator := range database.GetKnownOperators() {
		fromName := operatorFromName(operator.Code)
		if _, ok := operatorsMap[fromName.Code]; !ok {
			operatorsMap[fromName.Code] = fromName
		}
	}

	rooms := make([]model.Room, 0, 5)
	for _, r := range roomsMap {
		rooms = append(rooms, r)
	}
	sort.Slice(rooms, func(i, j int) bool {
		preferredSortRes := rooms[i].PreferredOrder - rooms[j].PreferredOrder
		if preferredSortRes != 0 {
			return preferredSortRes < 0
		}
		return rooms[i].Name < rooms[j].Name
	})

	operators := make([]model.Operator, 0, 5)
	for _, o := range operatorsMap {
		operators = append(operators, o)
	}
	sort.Slice(operators, func(i, j int) bool {
		return operators[i].Name < operators[j].Name
	})

	return model.ParsedData{
		Rows:         results,
		Rooms:        rooms,
		RoomsMap:     roomsMap,
		Operators:    operators,
		OperatorsMap: operatorsMap,
	}, nil
}

func parseRow(r ExcelRow, args config.Args) (model.ParsedRow, error) {
	localTimeZone := config.TimeZone()

	data, err := time.Parse(layoutDateOnlyInITFormat, r.Data)
	if err != nil {
		return model.ParsedRow{}, fmt.Errorf("il valore '%s' non e' una data valida nel formato 'GG/MM/YYYY'", r.Data)
	}

	orari := strings.Split(r.Orario, "-")

	v := strings.TrimSpace(orari[0])
	start, err := time.Parse(layoutTimeOnlyWithMinutes, v)
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithSeconds, v)
	}
	if err != nil {
		start, err = time.Parse(layoutTimeOnlyWithHours, v)
	}
	if err != nil {
		return model.ParsedRow{}, fmt.Errorf("il valore '%s' non e' un orario di inizio valido nei formati 'HH:MM', 'HH:MM:SS' o 'HH'", v)
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
		return model.ParsedRow{}, fmt.Errorf("il valore '%s' non e' un orario di fine valido nei formati 'HH:MM', 'HH:MM:SS' o 'HH'", v)
	}

	start = time.Date(data.Year(), data.Month(), data.Day(), start.Hour(), start.Minute(), start.Second(), 0, localTimeZone)
	end = time.Date(data.Year(), data.Month(), data.Day(), end.Hour(), end.Minute(), end.Second(), 0, localTimeZone)

	if end.Before(start) {
		end = time.Date(data.Year(), data.Month(), data.Day()+1, end.Hour(), end.Minute(), end.Second(), 0, localTimeZone)
	}

	room := roomFromName(r.Aula)
	operator := operatorFromName(r.Educatore)

	numPaganti := 0
	numGratuiti := 0
	numAccompagnatori := 0

	if r.Paganti != "" {
		numPaganti, err = strconv.Atoi(r.Paganti)
		if err != nil {
			return model.ParsedRow{}, fmt.Errorf("il valore '%s' non e' un numero valido", r.Paganti)
		}
	}
	if r.Gratuiti != "" {
		numGratuiti, err = strconv.Atoi(r.Gratuiti)
		if err != nil {
			return model.ParsedRow{}, fmt.Errorf("il valore '%s' non e' un numero valido", r.Gratuiti)
		}
	}
	if r.Accompagnatori != "" {
		numAccompagnatori, err = strconv.Atoi(r.Accompagnatori)
		if err != nil {
			return model.ParsedRow{}, fmt.Errorf("il valore '%s' non e' un numero valido", r.Accompagnatori)
		}
	}

	r2 := model.ParsedRow{
		ID:       database.NextSequenceValue(),
		Code:     r.Codice,
		StartAt:  start,
		EndAt:    end,
		Duration: end.Sub(start),
		Room:     room,
		Operator: operator,
		Activity: model.Activity{
			Description: r.Attivita,
			Language:    r.LinguaAttivita,
			Type:        r.Tipologia,
		},
		School: model.School{
			Type: r.TipologiaScuola,
			Name: r.NomeScuola,
		},
		SchoolClass: model.SchoolClass{
			Number:  r.Classe,
			Section: r.Sezione,
		},
		GroupComposition: model.GroupComposition{
			Total:           uint(numPaganti + numGratuiti + numAccompagnatori),
			NumPaying:       uint(numPaganti),
			NumFree:         uint(numGratuiti),
			NumAccompanying: uint(numAccompagnatori),
		},
		DateString:    r.Data,
		BookingNotes:  r.NotaPrenotazione,
		OperatorNotes: r.NotaOperatore,
		Bus:           r.Bus,
		PaymentStatus: model.PaymentStatus{
			PaymentAdvance:       r.Acconti,
			PaymentAdvanceStatus: r.StatoAcconti,
		},
	}

	r2.Warnings = getWarnings(r2, args)

	return r2, nil
}
