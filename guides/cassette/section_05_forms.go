package cassette

type selectOption struct {
	Value string
	Label string
}

var classificationOptions = []selectOption{
	{"routine", "ROUTINE"},
	{"priority", "PRIORITY"},
	{"classified", "CLASSIFIED"},
	{"eyes-only", "EYES ONLY"},
}

type checkboxOption struct {
	Value string
	Label string
}

var systemsCheckboxes = []checkboxOption{
	{"life-support", "LIFE SUPPORT"},
	{"propulsion", "PROPULSION"},
	{"comms", "COMMUNICATIONS"},
}

type radioOption struct {
	Value       string
	Label       string
	StatusClass string
	Checked     bool
}

var crewStatusRadios = []radioOption{
	{"nominal", "NOMINAL", "cass-value-ok", true},
	{"degraded", "DEGRADED", "cass-value-warn", false},
	{"critical", "CRITICAL", "cass-value-danger", false},
}

var thrustTickLabels = []string{"0% — OFFLINE", "50% — CRUISE", "100% — MAX"}
