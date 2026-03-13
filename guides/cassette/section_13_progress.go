package cassette

var authSequenceSteps = []string{"INITIATE", "VERIFY", "AUTHORIZE", "EXECUTE", "CONFIRM"}

type progressItem struct {
	Label      string
	Percent    int
	FillClass  string
	ValueClass string
}

var missionPhases = []progressItem{
	{"PRE-LAUNCH CHECKLIST", 100, "cass-progress-fill-green", "cass-value-ok"},
	{"TRANSIT TO LV-426", 67, "", "cass-value"},
	{"SURFACE SURVEY", 12, "", "cass-value"},
	{"SPECIMEN RECOVERY", 0, "", ""},
}

type taskItem struct {
	Task        string
	Owner       string
	Percent     int
	FillClass   string
	Status      string
	StatusClass string
}

var missionTasks = []taskItem{
	{"Atmospheric Processor Repair", "HICKS, D.", 100, "cass-progress-fill-green", "COMPLETE", "cass-value-ok"},
	{"Colonial Lab Sweep", "VASQUEZ, J.", 73, "", "IN PROGRESS", "cass-value"},
	{"Alien Nest Survey", "BISHOP", 40, "", "IN PROGRESS", "cass-value"},
	{"Drop Ship Maintenance", "FERRO, C.", 20, "", "DELAYED", "cass-value-warn"},
	{"Perimeter Defense Setup", "APONE, SGT.", 0, "cass-progress-fill-red", "NOT STARTED", "cass-value-danger"},
}
