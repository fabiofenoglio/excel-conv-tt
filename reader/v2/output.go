package reader

type OutputRow Row

type Output struct {
	Rows []OutputRow
}

func ToOutput(rows []Row) Output {
	outRows := make([]OutputRow, 0, len(rows))

	for _, row := range rows {
		outRows = append(outRows, OutputRow(row))
	}

	return Output{
		Rows: outRows,
	}
}
