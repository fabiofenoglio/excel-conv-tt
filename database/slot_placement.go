package database

type SlotPlacementPreferences struct {
	PointsForFittingWithSmallClearance        int
	PointsForFittingWithLotClearance          int
	PenaltyForActivityImmediatelyToTheRight   int
	PenaltyForActivity2ndToTheRight           int
	PenaltyForActivityImmediatelyToTheLeft    int
	PointsForOperatorHasOtherActivitiesInSlot int
	PointsForGroupHasOtherActivitiesInSlot    int
	PenaltyForEachOtherActivityInSlot         int
}

var (
	DefaultSlotPlacementPreferences = SlotPlacementPreferences{
		PointsForFittingWithSmallClearance:        50,
		PointsForFittingWithLotClearance:          10,
		PenaltyForActivityImmediatelyToTheRight:   25,
		PenaltyForActivity2ndToTheRight:           5,
		PenaltyForActivityImmediatelyToTheLeft:    0,
		PointsForOperatorHasOtherActivitiesInSlot: 15,
		PointsForGroupHasOtherActivitiesInSlot:    40,
		PenaltyForEachOtherActivityInSlot:         5,
	}

	planetarioSlotPlacementPreferences = SlotPlacementPreferences{
		PointsForFittingWithSmallClearance:        10,
		PointsForFittingWithLotClearance:          0,
		PenaltyForActivityImmediatelyToTheRight:   -50,
		PenaltyForActivity2ndToTheRight:           0,
		PenaltyForActivityImmediatelyToTheLeft:    -50,
		PointsForOperatorHasOtherActivitiesInSlot: 15,
		PointsForGroupHasOtherActivitiesInSlot:    30,
		PenaltyForEachOtherActivityInSlot:         5,
	}
)
