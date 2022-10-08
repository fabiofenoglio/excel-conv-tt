package model

import "strings"

type SchoolClass struct {
	Number  string `json:"number"`
	Section string `json:"section"`
}

func (s SchoolClass) FullDescription() string {
	out := ""
	if s.Number != "" {
		out += s.Number + " "
	}
	if s.Section != "" {
		out += s.Section
	}
	return strings.TrimSpace(out)
}

func (s SchoolClass) String() string {
	return s.FullDescription()
}
