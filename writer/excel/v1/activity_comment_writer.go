package excel

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/fabiofenoglio/excelconv/excel"
	"github.com/fabiofenoglio/excelconv/model"
)

func buildContentOfActivityComment(act model.ParsedRow) string {
	cellComment := ``

	for _, warning := range act.Warnings {
		cellComment += "ATTENZIONE: " + warning.Message + "\n\n"
	}

	if act.Activity.Description != "" {
		cellComment += act.Activity.Description
		if act.Activity.Type != "" {
			cellComment += " (" + act.Activity.Type + ")"
		}
		cellComment += "\n\n"

	} else if act.Activity.Type != "" {
		cellComment += "Tipologia: " + act.Activity.Type + "\n"
	}

	if act.TimeString != "" {
		cellComment += "Orario: " + act.TimeString + "\n"
	}
	if act.Operator.Code != "" {
		cellComment += "Educatore: " + act.Operator.Name + "\n"
	}
	if act.Room.Code != "" {
		cellComment += "Aula: " + act.Room.Name + "\n"
	}

	if act.OperatorNotes != "" {
		cellComment += "Nota operatore: " + act.OperatorNotes + "\n"
	}
	if act.BookingNotes != "" {
		cellComment += "Nota prenotazione: " + act.BookingNotes + "\n"
	}

	if act.SchoolClass.FullDescription() != "" {
		cellComment += "Classe: " + act.SchoolClass.FullDescription() + "\n"
	}

	if act.GroupComposition.Total > 0 {
		c := ""
		entries := 0
		if act.GroupComposition.NumPaying > 0 {
			c += fmt.Sprintf("%d paganti, ", act.GroupComposition.NumPaying)
			entries++
		}
		if act.GroupComposition.NumFree > 0 {
			c += fmt.Sprintf("%d gratuiti, ", act.GroupComposition.NumFree)
			entries++
		}
		if act.GroupComposition.NumAccompanying > 0 {
			c += fmt.Sprintf("%d accompagnatori, ", act.GroupComposition.NumAccompanying)
			entries++
		}
		if entries > 1 {
			c = strings.TrimSuffix(c, ", ") + fmt.Sprintf(" (%d totali)", act.GroupComposition.Total)
		}
		cellComment += strings.TrimSuffix(c, ", ") + "\n"
	}

	if act.School.FullDescription() != "" {
		cellComment += act.School.FullDescription() + "\n"
	}

	if act.Bus != "" {
		cellComment += "Bus: " + act.Bus + "\n"
	}
	if act.PaymentStatus.PaymentAdvance != "" && act.PaymentStatus.PaymentAdvance != "-" {
		cellComment += "Acconti: " + act.PaymentStatus.PaymentAdvance + "\n"
	}
	if act.PaymentStatus.PaymentAdvanceStatus != "" && act.PaymentStatus.PaymentAdvanceStatus != "-" {
		cellComment += "Stato acconti: " + act.PaymentStatus.PaymentAdvanceStatus + "\n"
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
