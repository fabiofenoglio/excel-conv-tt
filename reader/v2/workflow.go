package reader

import (
	"math/rand"
	"time"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/pkg/errors"
)

func Execute(ctx config.WorkflowContext, input Input) (Output, error) {

	rows, err := ReadFromFile(ctx, input.FilePath)
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella lettura dei dati dal file di input")
	}

	rows, err = Filter(rows)
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella scrematura iniziale delle righe dal file di input")
	}

	err = Validate(rows)
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella validazione dei dati di input")
	}

	rows, err = Convert(rows)
	if err != nil {
		return Output{}, errors.Wrap(err, "errore nella conversione dei dati di input")
	}

	// randomizer: randomize rows to enforce full sorting
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rows), func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })

	return ToOutput(rows), nil
}
