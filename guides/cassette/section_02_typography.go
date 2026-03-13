package cassette

type typeSpecimen struct {
	FontName   string
	SizeLabel  string
	SampleText string
	Style      string // inline CSS
}

var typeSpecimens = []typeSpecimen{
	{"Orbitron 700", "2rem / Display", "TECHNICAL REFERENCE MANUAL", "font-family: var(--font-display); font-size: 2rem; font-weight: 700; color: var(--color-primary);"},
	{"IBM Plex Mono 700", "0.875rem / Heading", "SYSTEM STATUS: ALL SUBSYSTEMS NOMINAL", "font-size: 0.875rem; font-weight: 700; letter-spacing: 0.02em;"},
	{"IBM Plex Mono 500", "0.8125rem / Body Medium", "Crew manifest updated 2183-06-03 at 0347 UTC", "font-size: 0.8125rem; font-weight: 500;"},
	{"IBM Plex Mono 400", "0.8125rem / Body Regular", "The atmospheric processor regulates colony life support across all 72 decks. Nominal operating pressure is 101.3 kPa with a tolerance of ±2.5 kPa.", "font-size: 0.8125rem; font-weight: 400; line-height: 1.6;"},
	{"IBM Plex Mono 300", "0.8125rem / Body Light", "Secondary system documentation, footnotes, and supplementary reference materials.", "font-size: 0.8125rem; font-weight: 300; color: var(--color-text-muted);"},
	{"IBM Plex Mono 400i", "0.8125rem / Body Italic", "Note: This document supersedes all prior revisions. Consult engineering before making any modifications.", "font-size: 0.8125rem; font-weight: 400; font-style: italic; color: var(--color-text-muted);"},
	{"IBM Plex Mono 700", "0.6875rem / Caption Uppercase", "CLASSIFICATION: COMPANY CONFIDENTIAL", "font-size: 0.6875rem; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: var(--color-text-muted);"},
}

const typeTableHeader = "TYPE SPECIMEN TABLE — IBM PLEX MONO + ORBITRON"

type paragraphSpecimen struct {
	Header string
	Body   string
}

var paragraphSpec = paragraphSpecimen{
	Header: "PARAGRAPH SPECIMEN — MISSION BRIEFING DOCUMENT",
	Body:   "WEYLAND-YUTANI CORP. — MISSION BRIEFING — USCSS NOSTROMO (MSV-180286) — The crew of the Nostromo has been diverted to planetoid LV-426 to investigate a transmission of unknown origin. Special Order 937 has been activated per standing directives from the Science Division. All non-essential crew members are to be kept uninformed of the nature of the mission objective. The Science Officer will assume command authority for all specimen-related activities. Crew safety is secondary to retrieval of specimen. This order supersedes all standing protocols established under the Colonial Marine Corps Charter, Articles 1 through 18. Any violation of this order constitutes breach of contract under ICC Charter and is subject to immediate punitive action.",
}
