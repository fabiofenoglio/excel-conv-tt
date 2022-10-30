package json

import (
	"encoding/json"
	"path/filepath"
	"strings"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	"github.com/fabiofenoglio/excelconv/config"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
	"github.com/fabiofenoglio/excelconv/writer"
	"github.com/pkg/errors"
)

type WriterImpl struct{}

var _ writer.WriterV2 = &WriterImpl{}

func (w *WriterImpl) Write(ctx config.WorkflowContext, parsed aggregator2.Output, anagraphicsRef *parser2.OutputAnagraphics) ([]byte, error) {
	log := ctx.Logger
	log.Debug("writing with JSON writer")

	out := SerializableOutput{
		CommonData:     parsed.CommonData,
		Days:           parsed.Days,
		AnagraphicsRef: anagraphicsRef,
	}

	serialized, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "error serializing data as JSON")
	}

	return serialized, nil
}

func (w *WriterImpl) ComputeDefaultOutputFile(inputFile string) string {
	outPath := filepath.Dir(inputFile)
	inputName := filepath.Base(inputFile)
	inputExt := filepath.Ext(inputFile)
	return outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-parsed.json"
}
