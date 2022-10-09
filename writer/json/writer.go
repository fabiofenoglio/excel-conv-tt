package json

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/model"
	"github.com/fabiofenoglio/excelconv/writer"
)

type WriterImpl struct{}

var _ writer.Writer = &WriterImpl{}

func (w *WriterImpl) Write(parsed model.ParsedData, args config.Args, log *logrus.Logger) ([]byte, error) {
	log.Debug("writing with JSON writer")

	out := SerializableOutput{
		Activities: parsed.Rows,
		Rooms:      parsed.Rooms,
		Operators:  parsed.Operators,
	}

	serialized, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializing data as JSON: %w", err)
	}

	return serialized, nil
}

func (w *WriterImpl) ComputeDefaultOutputFile(inputFile string) string {
	outPath := filepath.Dir(inputFile)
	inputName := filepath.Base(inputFile)
	inputExt := filepath.Ext(inputFile)
	return outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-parsed.json"
}
