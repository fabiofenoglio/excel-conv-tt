package model

type ParsedData struct {
	Rows []ParsedRow

	Rooms    []Room
	RoomsMap map[string]Room

	Operators    []Operator
	OperatorsMap map[string]Operator
}
