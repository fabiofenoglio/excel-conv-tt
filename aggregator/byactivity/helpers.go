package byactivity

func numToChars(columnNumber uint) string {

	// To store result (Excel column name)
	var columnName = ""
	for ok := true; ok; ok = columnNumber > 0 {

		// Find remainder
		rem := columnNumber % 26

		// If remainder is 0, then a
		// 'Z' must be there in output
		if rem == 0 {
			columnName += "z"
			columnNumber = (columnNumber / 26) - 1
		} else // If remainder is non-zero
		{
			columnName += string(rune((rem - 1) + uint('a'))) //nolint:govet
			columnNumber = columnNumber / 26
		}
	}

	// Reverse the string
	columnName = reverse(columnName)
	return columnName
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
