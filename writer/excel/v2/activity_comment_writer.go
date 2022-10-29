package excel

import (
	"encoding/json"
	"fmt"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/excel"
)

func buildContentOfActivityComment(c WriteContext, groupedActivities aggregator2.GroupedActivity) string {
	cellComment := ``

	for _, warning := range groupedActivities.Warnings() {
		cellComment += "ATTENZIONE: " + warning.Message + "\n\n"
	}

	room := c.anagraphicsRef.Rooms[groupedActivities.Rows[0].RoomCode]

	if room.Name != "" {
		cellComment += "Aula: " + room.Name + "\n"
	}

	if !groupedActivities.StartTime.IsZero() {
		cellComment += "Orario: " + groupedActivities.StartTime.Format(layoutTimeOnlyInReadableFormat) + " - " +
			groupedActivities.EndTime.Format(layoutTimeOnlyInReadableFormat) + "\n"
	}

	for _, act := range groupedActivities.Rows {

		operator := c.anagraphicsRef.Operators[act.OperatorCode]
		activityRef := c.anagraphicsRef.Activities[act.ActivityCode]
		activityTypeRef := c.anagraphicsRef.ActivityTypes[activityRef.TypeCode]

		groupRef := c.anagraphicsRef.VisitingGroups[act.VisitingGroupCode]
		school := c.anagraphicsRef.Schools[groupRef.SchoolCode]
		schoolClass := c.anagraphicsRef.SchoolClasses[groupRef.SchoolClassCode]

		if activityRef.Name != "" {
			cellComment += activityRef.Name
			if activityTypeRef.Name != "" {
				cellComment += " (" + activityTypeRef.Name + ")"
			}
			cellComment += "\n\n"

		} else if activityTypeRef.Name != "" {
			cellComment += "Tipologia: " + activityTypeRef.Name + "\n"
		}

		if operator.Name != "" {
			cellComment += "Educatore: " + operator.Name + "\n"
		}

		if act.OperatorNote != "" {
			cellComment += "Nota operatore: " + act.OperatorNote + "\n"
		}
		if act.BookingNote != "" {
			cellComment += "Nota prenotazione: " + act.BookingNote + "\n"
		}

		if schoolClass.FullDescription() != "" {
			cellComment += "Classe: " + schoolClass.FullDescription() + "\n"
		}

		if school.FullDescription() != "" {
			cellComment += school.FullDescription() + "\n"
		}

		if groupRef.Composition.NumTotal() > 0 {
			c := ""
			entries := 0
			if groupRef.Composition.NumPaying > 0 {
				c += fmt.Sprintf("%d paganti, ", groupRef.Composition.NumPaying)
				entries++
			}
			if groupRef.Composition.NumFree > 0 {
				c += fmt.Sprintf("%d gratuiti, ", groupRef.Composition.NumFree)
				entries++
			}
			if groupRef.Composition.NumAccompanying > 0 {
				c += fmt.Sprintf("%d accompagnatori, ", groupRef.Composition.NumAccompanying)
				entries++
			}
			if entries > 1 {
				c = strings.TrimSuffix(c, ", ") + fmt.Sprintf(" (%d totali)", groupRef.Composition.NumTotal())
			}
			cellComment += strings.TrimSuffix(c, ", ") + "\n"
		}

		if act.Bus != "" {
			cellComment += "Bus: " + act.Bus + "\n"
		}
		if act.Payment.PaymentAdvance != "" && act.Payment.PaymentAdvance != "-" {
			cellComment += "Acconti: " + act.Payment.PaymentAdvance + "\n"
		}
		if act.Payment.PaymentAdvanceStatus != "" && act.Payment.PaymentAdvanceStatus != "-" {
			cellComment += "Stato acconti: " + act.Payment.PaymentAdvanceStatus + "\n"
		}

		if len(groupedActivities.Rows) > 1 {
			cellComment += "--------------------------\n"
		}
	}

	if len(groupedActivities.FitComputationLog) > 0 {
		cellComment += "Fattori di piazzamento: \n"
		for _, log := range groupedActivities.FitComputationLog {
			cellComment += "  " + log + "\n"
		}
	}

	return strings.TrimSpace(cellComment)
}

func buildContentOfActivityCommentForSingleGroupOfGroupedActivity(c WriteContext, act aggregator2.OutputRow) string {
	cellComment := ``

	operator := c.anagraphicsRef.Operators[act.OperatorCode]
	activityRef := c.anagraphicsRef.Activities[act.ActivityCode]
	activityTypeRef := c.anagraphicsRef.ActivityTypes[activityRef.TypeCode]

	groupRef := c.anagraphicsRef.VisitingGroups[act.VisitingGroupCode]
	school := c.anagraphicsRef.Schools[groupRef.SchoolCode]
	schoolClass := c.anagraphicsRef.SchoolClasses[groupRef.SchoolClassCode]

	if activityRef.Name != "" {
		cellComment += activityRef.Name
		if activityTypeRef.Name != "" {
			cellComment += " (" + activityTypeRef.Name + ")"
		}
		cellComment += "\n\n"

	} else if activityTypeRef.Name != "" {
		cellComment += "Tipologia: " + activityTypeRef.Name + "\n"
	}

	if operator.Name != "" {
		cellComment += "Educatore: " + operator.Name + "\n"
	}

	if act.OperatorNote != "" {
		cellComment += "Nota operatore: " + act.OperatorNote + "\n"
	}
	if act.BookingNote != "" {
		cellComment += "Nota prenotazione: " + act.BookingNote + "\n"
	}

	if schoolClass.FullDescription() != "" {
		cellComment += "Classe: " + schoolClass.FullDescription() + "\n"
	}

	if school.FullDescription() != "" {
		cellComment += school.FullDescription() + "\n"
	}

	if groupRef.Composition.NumTotal() > 0 {
		c := ""
		entries := 0
		if groupRef.Composition.NumPaying > 0 {
			c += fmt.Sprintf("%d paganti, ", groupRef.Composition.NumPaying)
			entries++
		}
		if groupRef.Composition.NumFree > 0 {
			c += fmt.Sprintf("%d gratuiti, ", groupRef.Composition.NumFree)
			entries++
		}
		if groupRef.Composition.NumAccompanying > 0 {
			c += fmt.Sprintf("%d accompagnatori, ", groupRef.Composition.NumAccompanying)
			entries++
		}
		if entries > 1 {
			c = strings.TrimSuffix(c, ", ") + fmt.Sprintf(" (%d totali)", groupRef.Composition.NumTotal())
		}
		cellComment += strings.TrimSuffix(c, ", ") + "\n"
	}

	if act.Bus != "" {
		cellComment += "Bus: " + act.Bus + "\n"
	}
	if act.Payment.PaymentAdvance != "" && act.Payment.PaymentAdvance != "-" {
		cellComment += "Acconti: " + act.Payment.PaymentAdvance + "\n"
	}
	if act.Payment.PaymentAdvanceStatus != "" && act.Payment.PaymentAdvanceStatus != "-" {
		cellComment += "Stato acconti: " + act.Payment.PaymentAdvanceStatus + "\n"
	}

	return strings.TrimSpace(cellComment)
}

func addCommentToCell(f *excelize.File, cell excel.Cell, content string) error {
	type serializable struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	}
	v := serializable{
		Author: "planner ",
		Text:   content,
	}
	serialized, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return f.AddComment(cell.SheetName(), cell.Code(), string(serialized))
}
