package json

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/aggregator"
	"github.com/fabiofenoglio/excelconv/model"
	"github.com/fabiofenoglio/excelconv/writer"
)

type WriterImpl struct{}

var _ writer.Writer = &WriterImpl{}

func (w *WriterImpl) Write(parsed model.ParsedData, log *logrus.Logger) ([]byte, error) {
	log.Debug("writing with HTML writer")

	groupedByStartDate := aggregator.GroupByStartDay(parsed.Rows)
	log.Infof("found activities spanning %d different days", len(groupedByStartDate))

	pageTitle := fmt.Sprintf(
		"Settimana %s - %s",
		groupedByStartDate[0].Rows[0].StartAt.Format("02/01"),
		groupedByStartDate[len(groupedByStartDate)-1].Rows[0].StartAt.Format("02/01"),
	)

	templateData := TemplateData{
		GenerationTimeStamp: time.Now().UTC().Format("02/01/2006 15:04:05"),
		PageTitle:           pageTitle,
		Activities:          parsed.Rows,
		Rooms:               parsed.Rooms,
		Operators:           parsed.Operators,
	}

	tmpl, err := template.ParseGlob("./writer/html/template/*")
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(out, "main.html", templateData)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (w *WriterImpl) ComputeDefaultOutputFile(inputFile string) string {
	outPath := filepath.Dir(inputFile)
	inputName := filepath.Base(inputFile)
	inputExt := filepath.Ext(inputFile)
	return outPath + "/" + strings.TrimSuffix(inputName, inputExt) + "-parsed.html"
}
