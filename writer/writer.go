package writer

import (
	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/config"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
)

type Writer interface {
	ComputeDefaultOutputFile(inputFile string) string

	Write(ctx config.WorkflowContext, parsed aggregator2.Output, anagraphicsRef *parser2.OutputAnagraphics) ([]byte, error)
}
