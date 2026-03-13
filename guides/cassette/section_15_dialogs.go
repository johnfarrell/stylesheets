package cassette

type modalTrigger struct {
	Label     string
	BtnClass  string
	ResetExpr string // Alpine expression
}

var modalTriggers = []modalTrigger{
	{"REQUEST ACCESS", "cass-btn cass-btn-filled", "modal='auth'"},
	{"CONFIRM JETTISON", "cass-btn cass-btn-danger", "modal='jettison'; authCode=''; confirmed=false"},
	{"SYSTEM NOTIFICATION", "cass-btn", "modal='alert'"},
}

var jettisonCargo = []string{
	"Refinery Module RFN-7 (8,400 MT)",
	"Processing Equipment Set A (12,000 MT)",
	"Cargo Bay 3 Contents (2,200 MT)",
}

var alertMetadataCells = []docCell{
	{"MSG REF", "SYS-MAINT-2026-003", ""},
	{"PRIORITY", "ROUTINE", ""},
	{"ISSUED BY", "MU/TH/UR 6000", ""},
	{"TIMESTAMP", "2026-03-06 01:42", ""},
}
