package reader

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fabiofenoglio/excelconv/config"
)

func Execute(ctx config.WorkflowContext, input Input) (Output, error) {

	rows, err := ReadFromFile(ctx, input.FilePath)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella lettura dei dati dal file di input: %w", err)
	}

	rows, err = Filter(rows)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella scrematura iniziale delle righe dal file di input: %w", err)
	}

	err = Validate(rows)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella validazione dei dati di input: %w", err)
	}

	rows, err = Convert(rows)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella conversione dei dati di input: %w", err)
	}

	// randomizer: randomize rows to enforce full sorting
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rows), func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })

	return ToOutput(rows), nil
}
