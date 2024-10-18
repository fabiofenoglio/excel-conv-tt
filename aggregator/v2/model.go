package aggregator

import (
	"sort"
	"strings"
	"time"

	"github.com/fabiofenoglio/excelconv/parser/v2"
)

type Row struct {
	InputRow       InputRow
	CompetenceDate time.Time
}

type CommonTimespan struct {
	Start TimeOfDay
	End   TimeOfDay
}

type CommonData struct {
	CommonTimespan          CommonTimespan
	MaxVisitingGroupsPerDay int
}

type TimeOfDay struct {
	Hour, Minute int
}

func (tod TimeOfDay) IsBefore(other TimeOfDay) bool {
	if tod.Hour < other.Hour {
		return true
	}
	if tod.Hour > other.Hour {
		return false
	}
	return tod.Minute < other.Minute
}

func (tod TimeOfDay) IsAfter(other TimeOfDay) bool {
	if tod.Hour > other.Hour {
		return true
	}
	if tod.Hour < other.Hour {
		return false
	}
	return tod.Minute > other.Minute
}

type VisitingGroupInDay struct {
	VisitingGroupCode string
	SequentialCode    string
	DisplayCode       string
	StartsAt          time.Time
}

type scheduleForSingleDay struct {
	Day                   time.Time
	Rows                  []Row
	VisitingGroups        []VisitingGroupInDay
	StartAt               time.Time
	EndAt                 time.Time
	NumeroAttivitaMarkers map[time.Time]int
}

type ScheduleForSingleDayWithRooms struct {
	Day                   time.Time
	VisitingGroups        []VisitingGroupInDay
	RoomsSchedule         []ScheduleForSingleDayAndRoom
	StartAt               time.Time
	EndAt                 time.Time
	NumeroAttivitaMarkers map[time.Time]int
}

type ScheduleForSingleDayAndRoom struct {
	RoomCode string
	Rows     []OutputRow
}

type ScheduleForSingleDayWithRoomsAndGroupedActivities struct {
	Day                                   time.Time
	VisitingGroups                        []VisitingGroupInDay
	RoomsSchedule                         []ScheduleForSingleDayAndRoomWithGroupedActivities
	StartAt                               time.Time
	EndAt                                 time.Time
	NumeroAttivitaMarkers                 map[time.Time]int
	NumeroGruppiAttivitaConfermateMarkers map[time.Time]int
}

type ScheduleForSingleDayAndRoomWithGroupedActivities struct {
	RoomCode          string
	GroupedActivities []GroupedActivity
}

type GroupedActivity struct {
	ID                int
	StartingSlotIndex int
	NumOccupiedSlots  int
	StartTime         time.Time
	EndTime           time.Time
	Rows              []OutputRow
	FitComputationLog []string
	AnyConfirmed      bool
}

func (g *GroupedActivity) distinct(extractor func(OutputRow) string) []string {
	index := make(map[string]bool)
	for _, o := range g.Rows {
		v := extractor(o)
		if v == "" {
			continue
		}
		if _, ok := index[v]; !ok {
			index[v] = true
		}
	}
	out := make([]string, 0, len(index))
	for k := range index {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.Compare(out[i], out[j]) < 0
	})
	return out
}

func (g *GroupedActivity) OperatorCodes() []string {
	return g.distinct(func(row OutputRow) string {
		return row.OperatorCode
	})
}
func (g *GroupedActivity) BookingCodes() []string {
	return g.distinct(func(row OutputRow) string {
		return row.BookingCode
	})
}
func (g *GroupedActivity) VisitingGroupCodes() []string {
	return g.distinct(func(row OutputRow) string {
		return row.VisitingGroupCode
	})
}
func (g *GroupedActivity) ActivityCodes() []string {
	return g.distinct(func(row OutputRow) string {
		return row.ActivityCode
	})
}
func (g *GroupedActivity) Warnings() []parser.Warning {
	index := make(map[string]parser.Warning)
	for _, o := range g.Rows {
		for _, w := range o.Warnings {
			if _, ok := index[w.Code]; !ok {
				index[w.Code] = w
			}
		}
	}
	out := make([]parser.Warning, 0, len(index))
	for _, v := range index {
		out = append(out, v)
	}
	return out
}

// with slots (OLD)

type ScheduleForSingleDayWithRoomsAndSlots struct {
	Day            time.Time
	VisitingGroups []VisitingGroupInDay
	RoomsSchedule  []ScheduleForSingleDayAndRoomWithSlots
	StartAt        time.Time
	EndAt          time.Time
}

type ScheduleForSingleDayAndRoomWithSlots struct {
	RoomCode           string
	Slots              []ScheduleForSingleDayAndRoomSlot
	NumTotalInAllSlots int
}

type ScheduleForSingleDayAndRoomSlot struct {
	SlotIndex int
	Rows      []OutputRow
}

// with slots and groups (NEW)

type ScheduleForSingleDayWithRoomsAndGroupSlots struct {
	Day                                   time.Time
	VisitingGroups                        []VisitingGroupInDay
	RoomsSchedule                         []ScheduleForSingleDayAndRoomWithGroupSlots
	StartAt                               time.Time
	EndAt                                 time.Time
	NumeroAttivitaMarkers                 map[time.Time]int
	NumeroGruppiAttivitaConfermateMarkers map[time.Time]int
}

type ScheduleForSingleDayAndRoomWithGroupSlots struct {
	RoomCode           string
	Slots              []ScheduleForSingleDayAndRoomGroupSlot
	NumTotalInAllSlots int
}

type ScheduleForSingleDayAndRoomGroupSlot struct {
	SlotIndex         int
	GroupedActivities []GroupedActivity
}
