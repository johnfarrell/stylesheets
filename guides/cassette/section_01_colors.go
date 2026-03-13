package cassette

type swatchDef struct {
	CSSVar string
	Desc   string
}

type swatchGroupDef struct {
	Title    string
	Cols     string // Tailwind grid-cols class
	Swatches []swatchDef
}

var colorGroups = []swatchGroupDef{
	{"DOCUMENT COLORS", "grid-cols-2 sm:grid-cols-3", []swatchDef{
		{"--color-bg", "Page Background"},
		{"--color-surface", "Panel / Card Surface"},
		{"--color-surface-2", "Table Headers / Inputs"},
	}},
	{"PRIMARY PALETTE", "grid-cols-2 sm:grid-cols-3", []swatchDef{
		{"--color-primary", "Primary / Rules / Labels"},
		{"--color-secondary", "Hover / Secondary Actions"},
	}},
	{"STATUS COLORS", "grid-cols-2 sm:grid-cols-3", []swatchDef{
		{"--color-danger", "Fault / Warning / KIA"},
		{"--color-caution", "Caution / Degraded"},
	}},
	{"TEXT & BORDERS", "grid-cols-2 sm:grid-cols-4", []swatchDef{
		{"--color-text", "Body Text"},
		{"--color-text-muted", "Labels / Metadata"},
		{"--color-border", "Borders / Dividers"},
		{"--color-rule", "Section Rules"},
	}},
}
