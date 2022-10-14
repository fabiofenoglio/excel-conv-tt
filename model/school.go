package model

import "strings"

type School struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (s School) FullDescription() string {
	out := ""
	if s.Type != "" {
		out += s.Type + " "
	}
	if s.Name != "" {
		out += s.Name
	}
	return strings.TrimSpace(out)
}

func (s School) String() string {
	return s.FullDescription()
}

func (s School) Hash() string {
	out := ""
	if s.Type != "" {
		out += strings.ToLower(strings.TrimSpace(s.Type)) + "/"
	}
	if s.Name != "" {
		out += strings.ToLower(strings.TrimSpace(s.Name))
	}
	return out
}

func (s School) SortableIdentifier() string {
	out := ""
	if s.Type != "" {
		out += strings.ToLower(strings.TrimSpace(s.Type)) + "/"
	}
	if s.Name != "" {
		out += strings.ToLower(strings.TrimSpace(s.Name))
	}
	return out
}
