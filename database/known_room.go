package database

type KnownRoom struct {
	Code                 string
	Name                 string
	Slots                uint
	AllowMissingOperator bool
	PreferredOrder       int

	BackgroundColor          string
	SlotPlacementPreferences *SlotPlacementPreferences
	Aliases                  []string
	ShowActivityNames        bool
	Hide                     bool
}
