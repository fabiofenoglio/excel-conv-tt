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

func (s SchoolClass) Hash() string {
	out := ""
	if s.Number != "" {
		out += strings.ToLower(strings.TrimSpace(s.Number)) + "/"
	}
	if s.Section != "" {
		out += strings.ToLower(strings.TrimSpace(s.Section))
	}
	return out
}

func (s SchoolClass) SortableIdentifier() string {
	out := ""
	if s.Number != "" {
		out += strings.ToLower(strings.TrimSpace(s.Number)) + "/"
	}
	if s.Section != "" {
		out += strings.ToLower(strings.TrimSpace(s.Section))
	}
	return out
}
