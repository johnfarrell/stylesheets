package cassette

// Tab content data for navigation section.

type navItem struct {
	ID    string
	Label string
}

var navTabs = []navItem{
	{"overview", "OVERVIEW"},
	{"systems", "SYSTEMS"},
	{"crew", "CREW"},
	{"cargo", "CARGO"},
	{"logs", "LOGS"},
}

var tocItems = []navItem{
	{"1", "1.0  Color Palette"},
	{"2", "2.0  Typography"},
	{"3", "3.0  Spacing Scale"},
	{"4", "4.0  Button Components"},
	{"5", "5.0  Form Components"},
	{"6", "6.0  Panel Components"},
}

var tabOverviewParagraphs = []string{
	"The USCSS Nostromo (registration MSV-180286) is a modified Lockmart CM-88B Bison M-Class starfreighter, commissioned 2116 and currently assigned commercial haulage routes under contract ref. CMO-180286. The vessel has been diverted from its registered transit route by automated command subroutine Mother, acting under Special Order 937.",
	"All crew have been revived from hypersleep at grid reference Zeta II Reticuli. Navigation shows current position 34 light years from Earth. Mission duration estimate: 10 months transit to LV-426. All systems reading nominal except Navigation (override engaged) and Science Lab (access restricted under SO-937).",
}

type tabSystemRow struct {
	System  string
	Status  string
	Reading string
	Notes   string
}

var tabSystems = []tabSystemRow{
	{"Life Support", "NOMINAL", "101.3 kPa", "All within tolerances"},
	{"Propulsion", "NOMINAL", "98.7%", "Cruise configuration"},
	{"Navigation", "OVERRIDE", "SO-937", "Locked by Mother"},
	{"Communications", "ACTIVE", "Relay B", "Long-range active"},
}

type tabCrewRow struct {
	Name   string
	Role   string
	Status string
}

var tabCrew = []tabCrewRow{
	{"RIPLEY, E.", "Warrant Officer / Command", "ACTIVE"},
	{"BISHOP", "Synthetic / Science Division", "ACTIVE"},
	{"HICKS, D.", "Corporal / Colonial Marines", "ACTIVE"},
}

type tabCargoRow struct {
	Item   string
	Qty    string
	Mass   string
	Status string
}

var tabCargo = []tabCargoRow{
	{"Ore Processing Equipment", "1 set", "20,000 MT", "SECURED"},
	{"Refinery Consumables", "42 units", "8,400 MT", "SECURED"},
	{"Personnel Equipment", "7 kits", "350 kg", "SECURED"},
	{"Science Lab Samples", "CLASSIFIED", "CLASSIFIED", "SO-937"},
}

var tabLogs = []staticLogEntry{
	{"2183-06-01 04:12", "SYS", "Mother initiated hypersleep revival sequence"},
	{"2183-06-01 04:47", "NAV", "Course correction applied — new heading LV-426"},
	{"2183-06-01 05:03", "SCI", "Special Order 937 activated — Science Lab sealed"},
	{"2183-06-01 06:22", "COM", "Long-range signal detected — grid ref: Zeta Reticuli"},
	{"2183-06-01 07:55", "SYS", "Crew briefing complete — descent vehicle prepped"},
}
