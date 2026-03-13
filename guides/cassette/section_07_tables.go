package cassette

var crewTableHeaders = []string{"ID", "DESIGNATION", "RANK", "DEPARTMENT", "STATUS", "CLEARANCE"}
var envTableHeaders = []string{"SYSTEM", "CURRENT READING", "NOMINAL RANGE", "STATUS"}

type crewMember struct {
	ID        string
	Name      string
	Rank      string
	Dept      string
	Status    string
	Clearance string
}

var crewManifest = []crewMember{
	{"NOS-001", "DALLAS, A.J.", "CAPTAIN", "COMMAND", "ACTIVE", "L-5"},
	{"NOS-002", "RIPLEY, E.", "WARRANT OFFICER", "COMMAND", "ACTIVE", "L-4"},
	{"NOS-003", "KANE, G.", "EXECUTIVE OFFICER", "COMMAND", "KIA", "L-4"},
	{"NOS-004", "ASH", "SCIENCE OFFICER", "SCIENCE", "DECOMMISSIONED", "UNRESTRICTED"},
	{"NOS-005", "LAMBERT, J.", "NAVIGATOR", "COMMAND", "KIA", "L-3"},
	{"NOS-006", "BRETT, S.", "ENGINEER", "ENGINEERING", "KIA", "L-2"},
	{"NOS-007", "PARKER, D.", "CHIEF ENGINEER", "ENGINEERING", "KIA", "L-2"},
}

type envReading struct {
	System  string
	Reading string
	Range   string
	Status  string
}

var envReadings = []envReading{
	{"Atmospheric Pressure", "101.3 kPa", "99.0–103.0 kPa", "NOMINAL"},
	{"Cabin Temperature", "19.7 °C", "18.0–22.0 °C", "NOMINAL"},
	{"Oxygen Partial Press.", "21.3 kPa", "20.5–22.0 kPa", "NOMINAL"},
	{"Reactor Coolant", "94.1%", "< 95%", "CAUTION"},
	{"Fuel Cell Capacity", "67.4%", "> 20%", "NOMINAL"},
	{"Hypersleep Units", "7/7", "7/7", "NOMINAL"},
}
