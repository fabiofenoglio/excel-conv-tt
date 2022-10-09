package byactivity

import (
	"time"

	"github.com/fabiofenoglio/excelconv/model"
)

type ActivityGroup struct {
	Code             string
	SequentialNumber uint
	School           model.School
	SchoolClass      model.SchoolClass
	Composition      model.GroupComposition

	AveragePresence time.Time
}
