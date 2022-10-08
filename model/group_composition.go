package model

type GroupComposition struct {
	Total           uint `json:"total"`
	NumPaying       uint `json:"num_paying"`
	NumFree         uint `json:"num_free"`
	NumAccompanying uint `json:"num_accompanying"`
}
