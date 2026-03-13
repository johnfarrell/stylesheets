package cassette

type mastheadDef struct {
	Tagline  string
	Title    string
	Subtitle string
	DocNo    string
	Revision string
	Date     string
}

var masthead = mastheadDef{
	Tagline:  "USCSS TECHNICAL REFERENCE // WCYPD COLONY SYSTEMS",
	Title:    "CASSETTE FUTURISM",
	Subtitle: "Design System Reference — NASA/TM-2026-CSS-001 — Revision D",
	DocNo:    "CSS-4.2.1",
	Revision: "D",
	Date:     "2026-03-06",
}
