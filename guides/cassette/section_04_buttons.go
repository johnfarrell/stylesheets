package cassette

type buttonVariant struct {
	Name     string
	Class    string
	Label    string
	UseCase  string
	Disabled bool
}

var buttonVariants = []buttonVariant{
	{"Outline Primary", "cass-btn", "INITIATE SCAN", "Standard actions, secondary operations", false},
	{"Filled Primary", "cass-btn cass-btn-filled", "CONFIRM ACTION", "Primary CTA, high-priority actions", false},
	{"Outline Danger", "cass-btn cass-btn-danger", "ABORT SEQUENCE", "Destructive secondary action", false},
	{"Filled Danger", "cass-btn cass-btn-danger-filled", "JETTISON CARGO", "Irreversible destructive action", false},
	{"Disabled", "cass-btn", "RESTRICTED", "Locked / insufficient clearance", true},
}

type buttonSize struct {
	SizeLabel string
	BtnLabel  string
	SizeStyle string // inline CSS override, empty for default
}

var buttonSizes = []buttonSize{
	{"SM", "INITIATE", "font-size: 0.6875rem; padding: 0.25rem 0.75rem;"},
	{"MD (default)", "TRANSMIT DATA", ""},
	{"LG", "EXECUTE PROTOCOL", "font-size: 0.875rem; padding: 0.6rem 1.5rem;"},
}
