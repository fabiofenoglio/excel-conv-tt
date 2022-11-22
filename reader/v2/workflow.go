package reader

import (
	"math/rand"
	"time"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

func Execute(ctx config.WorkflowContext, input Input) (Output, error) {

	span := sentry.StartSpan(ctx.Context, "read rows")
	rows, err := ReadFromFile(ctx.ForContext(span.Context()), input.FilePath)
	span.Finish()
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella lettura dei dati dal file di input")
	}

	span = sentry.StartSpan(ctx.Context, "filter rows")
	rows, err = Filter(rows)
	span.Finish()
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella scrematura iniziale delle righe dal file di input")
	}

	span = sentry.StartSpan(ctx.Context, "validate rows integrity")
	err = Validate(rows)
	span.Finish()
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella validazione dei dati di input")
	}

	span = sentry.StartSpan(ctx.Context, "convert rows data")
	rows, err = Convert(rows)
	span.Finish()
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella conversione dei dati di input")
	}

	span = sentry.StartSpan(ctx.Context, "apply A0-level rules")
	rows, err = ApplyRuleA0Level(ctx, rows)
	span.Finish()
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nell'applicazione delle regole di livello A0")
	}

	// randomizer: randomize rows to enforce full sorting
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rows), func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })

	return ToOutput(rows), nil
}
