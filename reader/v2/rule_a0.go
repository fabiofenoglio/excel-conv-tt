package reader

import (
	_ "time/tzdata"

	"github.com/fabiofenoglio/excelconv/config"
)

func ApplyRuleA0Level(
	_ config.WorkflowContext,
	rows []Row,
) ([]Row, error) {
	return rows, nil
}
