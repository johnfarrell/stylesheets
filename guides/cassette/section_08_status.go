package cassette

import "strings"

type systemStatus struct {
	Name   string
	Status string // statusOK, statusWarn, statusErr
}

var initialSystemStatuses = []systemStatus{
	{"LIFE SUPPORT", statusOK},
	{"PROPULSION", statusOK},
	{"NAVIGATION", statusWarn},
	{"COMMUNICATIONS", statusOK},
	{"POWER GRID", statusOK},
	{"HYPERSLEEP", statusOK},
	{"MOTION SENSORS", statusErr},
	{"ATMOSPHERIC", statusOK},
	{"FIRE SUPPRESSION", statusOK},
	{"CARGO LOCKS", statusWarn},
}

// systemStatusJSON returns an Alpine-compatible x-data object literal.
func systemStatusJSON() string {
	var b strings.Builder
	b.WriteString("{ systems: [")
	for i, s := range initialSystemStatuses {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString("{n:'")
		b.WriteString(s.Name)
		b.WriteString("',s:'")
		b.WriteString(s.Status)
		b.WriteString("'}")
	}
	b.WriteString("] }")
	return b.String()
}

type personnelEntry struct {
	Name     string
	Location string
	Status   string
}

var personnelTracking = []personnelEntry{
	{"RIPLEY, E.", "DECK C — MEDICAL", "ACTIVE"},
	{"BISHOP", "SCIENCE LAB", "ACTIVE"},
	{"HICKS, D.", "ARMORY", "ACTIVE"},
}
