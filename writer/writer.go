package writer

import (
	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/model"
)

type Writer interface {
	ComputeDefaultOutputFile(inputFile string) string

	Write(data model.ParsedData, args config.Args, log *logrus.Logger) ([]byte, error)
}
