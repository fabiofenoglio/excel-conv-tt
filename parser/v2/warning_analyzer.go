package parser

import (
	"fmt"
	"strings"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/pkg/errors"
)

func EmitWarnings(ctx config.WorkflowContext, rows []Row, anagraphicsRef *OutputAnagraphics) ([]Row, error) {
	out := make([]Row, 0, len(rows))

	for _, row := range rows {
		outCopy := row
		warnings, err := emitWarningsForRow(ctx, outCopy, anagraphicsRef)
		if err != nil {
			return nil, errors.Wrap(err, "errore nell'analisi di coerenza")
		}
		outCopy.Warnings = warnings
		out = append(out, outCopy)
	}

	return out, nil
}

func emitWarningsForRow(ctx config.WorkflowContext, row Row, anagraphicsRef *OutputAnagraphics) ([]Warning, error) {
	out := make([]Warning, 0)

	activity := anagraphicsRef.Activities[row.ActivityCode]
	room := anagraphicsRef.Rooms[row.RoomCode]
	group := anagraphicsRef.VisitingGroups[row.VisitingGroupCode]

	if strings.Contains(activity.Name, "??") {
		out = append(out, Warning{
			Code:    "activity-question-marks",
			Message: "QUESTA ATTIVITA' SEMBRA INDETERMINATA",
		})
	}

	if row.RoomCode == "" {
		out = append(out, Warning{
			Code:    "no-room",
			Message: "NESSUNA AULA O RISORSA ASSEGNATA",
		})
	}

	if ctx.Config.EnableMissingOperatorsWarning && row.OperatorCode == "" && !room.AllowMissingOperator {
		out = append(out, Warning{
			Code:    "no-operator",
			Message: "NESSUN EDUCATORE ASSEGNATO",
		})
	}

	if activity.Language != "" && strings.ToLower(activity.Language) != "it" {
		out = append(out, Warning{
			Code:    "non-it-lang",
			Message: "ATTIVITA' PREVISTA IN LINGUA: " + activity.Language,
		})
	}

	if group.Composition.NumFree > 0 {
		out = append(out, Warning{
			Code:    "free-admission-presence",
			Message: fmt.Sprintf("SONO PRESENTI %d INGRESSI GRATUITI", group.Composition.NumFree),
		})
	}

	return out, nil
}
