package parser

import (
	"fmt"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/reader/v2"
)

func Execute(ctx config.WorkflowContext, rawInput reader.Output) (Output, error) {
	input := Input{
		Rows: ToInputRows(rawInput.Rows),
	}

	rowsWithRooms, rooms, err := HydrateRooms(ctx, input.Rows)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella lettura delle aule: %w", err)
	}

	rowsWithOperators, operators, err := HydrateOperators(ctx, rowsWithRooms)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella lettura degli educatori: %w", err)
	}

	rowsWithGroups, groups, schools, schoolClasses, err := HydrateGroups(ctx, rowsWithOperators)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella lettura dei gruppi scuola: %w", err)
	}

	rowsWithActivities, activities, activityTypes, err := HydrateActivities(ctx, rowsWithGroups)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella lettura delle attività: %w", err)
	}

	// build anagraphics index for O(1) ammortized lookup

	anagraphics := OutputAnagraphics{
		Rooms:          make(map[string]Room),
		Operators:      make(map[string]Operator),
		VisitingGroups: make(map[string]VisitingGroup),
		Schools:        make(map[string]School),
		SchoolClasses:  make(map[string]SchoolClass),
		Activities:     make(map[string]Activity),
		ActivityTypes:  make(map[string]ActivityType),
	}

	for _, o := range rooms {
		anagraphics.Rooms[o.Code] = o
	}
	for _, o := range operators {
		anagraphics.Operators[o.Code] = o
	}
	for _, o := range groups {
		anagraphics.VisitingGroups[o.Code] = o
	}
	for _, o := range schools {
		anagraphics.Schools[o.Code] = o
	}
	for _, o := range schoolClasses {
		anagraphics.SchoolClasses[o.Code] = o
	}
	for _, o := range activities {
		anagraphics.Activities[o.Code] = o
	}
	for _, o := range activityTypes {
		anagraphics.ActivityTypes[o.Code] = o
	}

	rowsWithWarnings, err := EmitWarnings(ctx, rowsWithActivities, &anagraphics)
	if err != nil {
		return Output{}, fmt.Errorf("errore nella ricerca dei warning: %w", err)
	}

	out := Output{
		Anagraphics: &anagraphics,
		Rows:        ToOutputRows(rowsWithWarnings, &anagraphics),
	}

	return out, nil
}
