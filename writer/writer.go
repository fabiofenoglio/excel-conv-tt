package writer

import (
	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/model"
)

type Writer interface {
	ComputeDefaultOutputFile(inputFile string) string

	Write(data model.ParsedData, args config.Args, log *logrus.Logger) ([]byte, error)
}

type WriterV2 interface {
	ComputeDefaultOutputFile(inputFile string) string

	Write(ctx config.WorkflowContext, parsed aggregator2.Output, anagraphicsRef *parser2.OutputAnagraphics) ([]byte, error)
}
