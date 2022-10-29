package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	_ "time/tzdata"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/model"
	"github.com/fabiofenoglio/excelconv/reader/v1"
)

func validate(r reader.ExcelRow) error {
	if r.Codice == "" {
		return errors.New("codice is required")
	}

	timeRangeRegexp := regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])\s*\-\s*([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`)
	if !timeRangeRegexp.MatchString(r.Orario) {
		return fmt.Errorf("orario non valido, atteso HH:MM-HH:MM e non '%s'", r.Orario)
	}

	dateRegexp := regexp.MustCompile(`^\d{2}\/\d{2}\/\d{4}$`)
	if !dateRegexp.MatchString(r.Data) {
		return fmt.Errorf("data non valida, atteso GG/MM/YYYY e non '%s'", r.Data)
	}

	return nil
}

func getWarnings(act model.ParsedRow, args config.Args) []model.Warning {
	out := make([]model.Warning, 0)

	if strings.Contains(act.Activity.Description, "??") {
		out = append(out, model.Warning{
			Code:    "activity-question-marks",
			Message: "QUESTA ATTIVITA' SEMBRA INDETERMINATA",
		})
	}

	if act.Room.Code == "" {
		out = append(out, model.Warning{
			Code:    "no-room",
			Message: "NESSUNA AULA O RISORSA ASSEGNATA",
		})
	}

	if args.EnableMissingOperatorsWarning && act.Operator.Code == "" && !act.Room.AllowMissingOperator {
		out = append(out, model.Warning{
			Code:    "no-operator",
			Message: "NESSUN EDUCATORE ASSEGNATO",
		})
	}

	if act.Activity.Language != "" && strings.ToLower(act.Activity.Language) != "it" {
		out = append(out, model.Warning{
			Code:    "non-it-lang",
			Message: "ATTIVITA' PREVISTA IN LINGUA: " + act.Activity.Language,
		})
	}

	if act.GroupComposition.NumFree > 0 {
		out = append(out, model.Warning{
			Code:    "free-admission-presence",
			Message: fmt.Sprintf("SONO PRESENTI %d INGRESSI GRATUITI", act.GroupComposition.NumFree),
		})
	}

	return out
}
