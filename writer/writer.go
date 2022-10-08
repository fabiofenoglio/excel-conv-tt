package writer

import (
	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/model"
)

type Writer interface {
	ComputeDefaultOutputFile(inputFile string) string

	Write(data model.ParsedData, log *logrus.Logger) ([]byte, error)
}
