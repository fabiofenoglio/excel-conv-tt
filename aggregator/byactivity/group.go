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

	StartsAt        time.Time
	AveragePresence time.Time
}

func (g *ActivityGroup) SequentialCode() string {
	return fmt.Sprintf("%d-%s", g.SequentialNumberForSchool,
		numToChars(g.SequentialNumberInsideSchool))
}
