package services

import (
	"crypto/sha1"
	"encoding/base64"
	"sort"
	"strings"

	"github.com/fabiofenoglio/excelconv/model"
)

const (
	layoutDateOrderable            = "2006-01-02"
	layoutTimeOnlyInReadableFormat = "15:04"
)

func GroupByStartDay(rows []model.ParsedRow) []model.GroupedRows {

	// group by start date, ordering each group by start time ASC, end time ASC

	grouped := make([]*model.GroupedRows, 0)
	index := make(map[string]*model.GroupedRows)
	for _, activity := range rows {
		if activity.StartAt.IsZero() {
			continue
		}
		key := activity.StartAt.Format(layoutDateOrderable)

		group, ok := index[key]
		if !ok {
			group = &model.GroupedRows{
				Key: key,
			}
			index[key] = group
			grouped = append(grouped, group)
		}
		group.Rows = append(group.Rows, activity)
	}

	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Key < grouped[j].Key
	})

	for _, group := range grouped {
		sort.Slice(group.Rows, func(i, j int) bool {
			if group.Rows[i].StartAt.UnixMilli() < group.Rows[j].StartAt.UnixMilli() {
				return true
			}
			if group.Rows[i].StartAt.UnixMilli() > group.Rows[j].StartAt.UnixMilli() {
				return false
			}
			return group.Rows[i].EndAt.UnixMilli() < group.Rows[j].EndAt.UnixMilli()
		})
	}

	out := make([]model.GroupedRows, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func GroupByRoom(allData model.ParsedData, rows []model.ParsedRow) []model.GroupedRows {

	// group by room, ordering each group by start time ASC, end time ASC

	grouped := make([]*model.GroupedRows, 0)
	index := make(map[string]*model.GroupedRows)
	for _, activity := range rows {
		key := activity.Room.Code

		group, ok := index[key]
		if !ok {
			group = &model.GroupedRows{
				Key: key,
			}
			index[key] = group
			grouped = append(grouped, group)
		}
		group.Rows = append(group.Rows, activity)
	}

	for _, knownRoom := range allData.Rooms {
		if _, ok := index[knownRoom.Code]; !ok {
			group := &model.GroupedRows{
				Key: knownRoom.Code,
			}
			index[knownRoom.Code] = group
			grouped = append(grouped, group)
		}
	}

	sort.Slice(grouped, func(i, j int) bool {
		return grouped[i].Key < grouped[j].Key
	})

	for _, group := range grouped {
		sort.Slice(group.Rows, func(i, j int) bool {
			if group.Rows[i].StartAt.UnixMilli() < group.Rows[j].StartAt.UnixMilli() {
				return true
			}
			if group.Rows[i].StartAt.UnixMilli() > group.Rows[j].StartAt.UnixMilli() {
				return false
			}

			if group.Rows[i].EndAt.UnixMilli() < group.Rows[j].EndAt.UnixMilli() {
				return true
			}
			if group.Rows[i].EndAt.UnixMilli() > group.Rows[j].EndAt.UnixMilli() {
				return false
			}

			c := strings.Compare(group.Rows[i].Operator.Code, group.Rows[j].Operator.Code)
			return c < 0
		})
	}

	out := make([]model.GroupedRows, 0, len(grouped))
	for _, e := range grouped {
		out = append(out, *e)
	}
	return out
}

func GetDifferentSchoolGroups(rows []model.ParsedRow) []model.SchoolGroup {
	index := make(map[string]*model.SchoolGroup)
	grouped := make([]*model.SchoolGroup, 0)

	for _, row := range rows {
		if row.Raw.NomeScuola == "" && row.Raw.Classe == "" && row.Raw.Codice == "" {
			continue
		}
		keyBuilder := row.Raw.Codice
		if keyBuilder == "" {
			keyBuilder = row.Raw.NomeScuola + "|" + row.Raw.Classe + "|" + row.Raw.Sezione
		}
		key := Base64Sha([]byte(strings.ToLower(keyBuilder)))

		group, ok := index[key]
		if !ok {
			group = &model.SchoolGroup{
				Code:              key,
				Codice:            row.Raw.Codice,
				TipologiaScuola:   row.Raw.TipologiaScuola,
				NomeScuola:        row.Raw.NomeScuola,
				Classe:            row.Raw.Classe,
				Sezione:           row.Raw.Sezione,
				NumPaganti:        row.NumPaganti,
				NumGratuiti:       row.NumGratuiti,
				NumAccompagnatori: row.NumAccompagnatori,
			}
			index[key] = group
			grouped = append(grouped, group)
		}

		if row.Raw.TipologiaScuola != "" {
			group.TipologiaScuola = row.Raw.TipologiaScuola
		}
		if row.Raw.NomeScuola != "" {
			group.NomeScuola = row.Raw.NomeScuola
		}
		if row.Raw.Classe != "" {
			group.Classe = row.Raw.Classe
		}
		if row.Raw.Sezione != "" {
			group.Sezione = row.Raw.Sezione
		}
		if row.NumPaganti > 0 {
			group.NumPaganti = row.NumPaganti
		}
		if row.NumGratuiti > 0 {
			group.NumGratuiti = row.NumGratuiti
		}
		if row.NumAccompagnatori > 0 {
			group.NumAccompagnatori = row.NumAccompagnatori
		}
	}

	sort.SliceStable(grouped, func(i, j int) bool {
		c := strings.Compare(strings.ToLower(grouped[i].Codice), strings.ToLower(grouped[j].Codice))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].NomeScuola), strings.ToLower(grouped[j].NomeScuola))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].Classe), strings.ToLower(grouped[j].Classe))
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(strings.ToLower(grouped[i].Sezione), strings.ToLower(grouped[j].Sezione))
		if c != 0 {
			return c < 0
		}
		return false
	})

	out := make([]model.SchoolGroup, 0, len(grouped))
	for i, e := range grouped {
		e.NumeroSeq = i + 1
		out = append(out, *e)
	}
	return out
}

func Base64Sha(content []byte) string {
	h := sha1.New()
	h.Write(content)
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}
