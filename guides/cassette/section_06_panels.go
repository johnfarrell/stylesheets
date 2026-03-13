package cassette

type keyValuePair struct {
	Label string
	Value string
}

type coloredKeyValuePair struct {
	Label string
	Value string
	Color string // CSS color or empty
}

type panelData struct {
	Header string
	Rows   []keyValuePair
}

var missionDataPanel = panelData{
	Header: "MISSION DATA",
	Rows: []keyValuePair{
		{"VESSEL", "USCSS NOSTROMO"},
		{"REGISTRATION", "MSV-180286"},
		{"DESTINATION", "LV-426 / ZETA RETICULI"},
		{"MISSION PHASE", "APPROACH VECTOR"},
		{"ETD", "2183-06-01 0600 UTC"},
	},
}

var specialOrderPanel = panelData{
	Header: "SPECIAL ORDER 937 — ACTIVE",
	Rows: []keyValuePair{
		{"PRIORITY", "WEYLAND-YUTANI CORP."},
		{"CLASSIFICATION", "EYES ONLY"},
		{"OBJECTIVE", "RETRIEVE SPECIMEN"},
		{"CREW AWARENESS", "SCIENCE OFFICER ONLY"},
		{"ACTIVATED BY", "MOTHER — MU/TH/UR 6000"},
	},
}

var engineeringPanel = panelData{
	Header: "ENGINEERING SUBSYSTEMS",
	Rows: []keyValuePair{
		{"REACTOR STATUS", "ONLINE — 98.7%"},
		{"COOLANT TEMP.", "487°C — NOMINAL"},
		{"FUEL RESERVE", "67.4%"},
		{"THRUST CAPACITY", "100% AVAILABLE"},
	},
}

type accentNote struct {
	Title string
	Body  string
}

var navNote = accentNote{
	Title: "NAVIGATIONAL NOTE",
	Body:  "The Nostromo has been diverted from its registered course by automated subroutine activation. Navigation override is locked pending Science Officer authorization. ETA to LV-426: 18 hours at current velocity.",
}

var vesselStatusCells = []docCell{
	{"HULL INTEGRITY", "100%", "cass-value-ok"},
	{"SHIELDS", "N/A", ""},
	{"LIFE SUPPORT", "NOMINAL", "cass-value-ok"},
	{"REACTOR", "98.7%", "cass-value-ok"},
	{"NAVIGATION", "OVERRIDE", "cass-value-warn"},
	{"CREW STATUS", "7/7", "cass-value-ok"},
}

type dangerPanelData struct {
	Header string
	Rows   []coloredKeyValuePair
}

var dangerPanel = dangerPanelData{
	Header: "CRITICAL ALERT — XENOMORPH CONTAINMENT BREACH",
	Rows: []coloredKeyValuePair{
		{"THREAT LEVEL", "EXTREME", "#c0392b"},
		{"LOCATION", "DECK C — MEDICAL BAY", ""},
		{"PERSONNEL AT RISK", "7 CREW MEMBERS", "#c0392b"},
		{"RECOMMENDED ACTION", "INITIATE PROTOCOL", ""},
	},
}
