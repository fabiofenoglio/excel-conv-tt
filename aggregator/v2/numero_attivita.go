package aggregator

func SplitByNumeroAttivitaPlaceholder(input []Row) ([]Row, []Row) {
	var regular []Row
	var placeholder []Row
	for _, row := range input {
		if row.InputRow.IsPlaceholderNumeroAttivita {
			placeholder = append(placeholder, row)
		} else {
			regular = append(regular, row)
		}
	}

	return regular, placeholder
}
