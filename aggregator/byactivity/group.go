package byactivity

import (
	"fmt"
	"time"

	"github.com/fabiofenoglio/excelconv/model"
)

type ActivityGroup struct {
	Code string

	SequentialNumber             uint
	SequentialNumberForSchool    uint
	SequentialNumberInsideSchool uint

	School      model.School
	SchoolClass model.SchoolClass
	Composition model.GroupComposition
	Notes       string

	AveragePresence time.Time
}

func (g *ActivityGroup) SequentialCode() string {
	return fmt.Sprintf("%d-%d", g.SequentialNumberForSchool,
		g.SequentialNumberInsideSchool)
}
