package cassette

import (
	"fmt"
	"strings"
)

// --- Section 01: Colors ---

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

// --- Section 02: Typography ---

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

// --- Section 04: Buttons ---

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

// --- Section 09: Notices ---

type noticeDef struct {
	Severity string // "note", "caution", "warning"
	Message  string
}

var staticNotices = []noticeDef{
	{"note", "All personnel must complete atmospheric decompression protocol before entering Deck C. Standard decompression time is 15 minutes. Reference: SOP-ATM-042, Revision 3."},
	{"caution", "Reactor coolant pressure is operating at 94.1% of nominal limits. Monitor pressure gauge readings every 15 minutes. Initiate venting procedure if readings exceed 97% nominal. Contact Chief Engineer immediately."},
	{"warning", "UNAUTHORIZED ACCESS TO SCIENCE DIVISION LABORATORY IS STRICTLY PROHIBITED. Special Order 937 is in effect. Violation is subject to immediate contract termination and criminal prosecution under ICC Charter Article 14, Section 9."},
}

var dismissibleNotices = []noticeDef{
	{"note", "Hypersleep revival protocol requires a minimum 4-hour monitoring period. Medical Officer must be present for all revivals. Reference: MED-HS-001."},
	{"caution", "Motion sensor array sector 7G reporting intermittent signal loss. Maintenance crew dispatched. Estimated resolution: 2 hours. Do not rely on sector 7G coverage during this period."},
	{"warning", "AIRLOCK CYCLE DETECTED — DECK A EMERGENCY AIRLOCK. No crew authorization on record. Investigate immediately. Security to report to Deck A airlock station."},
}

// noticeLabelColor returns the CSS color for a given notice severity.
func noticeLabelColor(severity string) string {
	switch severity {
	case "caution":
		return "var(--color-caution)"
	case "warning":
		return "var(--color-danger)"
	default:
		return "var(--color-primary)"
	}
}

// noticeLabel returns the uppercase label for a notice severity.
func noticeLabel(severity string) string {
	return strings.ToUpper(severity)
}

// --- Section 06: Panels ---

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

// --- Section 07: Tables ---

var crewTableHeaders = []string{"ID", "DESIGNATION", "RANK", "DEPARTMENT", "STATUS", "CLEARANCE"}
var envTableHeaders = []string{"SYSTEM", "CURRENT READING", "NOMINAL RANGE", "STATUS"}

// --- Section 08: Status Monitor ---

type systemStatus struct {
	Name   string
	Status string // "ok", "warn", "err"
}

var initialSystemStatuses = []systemStatus{
	{"LIFE SUPPORT", "ok"},
	{"PROPULSION", "ok"},
	{"NAVIGATION", "warn"},
	{"COMMUNICATIONS", "ok"},
	{"POWER GRID", "ok"},
	{"HYPERSLEEP", "ok"},
	{"MOTION SENSORS", "err"},
	{"ATMOSPHERIC", "ok"},
	{"FIRE SUPPRESSION", "ok"},
	{"CARGO LOCKS", "warn"},
}

// systemStatusJSON returns an Alpine-compatible x-data object literal.
func systemStatusJSON() string {
	var b strings.Builder
	b.WriteString("{ systems: [")
	for i, s := range initialSystemStatuses {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString("{n:'")
		b.WriteString(s.Name)
		b.WriteString("',s:'")
		b.WriteString(s.Status)
		b.WriteString("'}")
	}
	b.WriteString("] }")
	return b.String()
}

// --- Section 13: Progress ---

var authSequenceSteps = []string{"INITIATE", "VERIFY", "AUTHORIZE", "EXECUTE", "CONFIRM"}

// --- Section 05: Forms ---

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

// --- Section 12: Readouts ---

type readoutDef struct {
	Label       string
	Unit        string
	Baseline    float64
	Variance    float64
	IntervalMs  int     // 0 = static (no Alpine update)
	FaultBelow  float64 // 0 = no fault threshold
	WarnBelow   float64 // 0 = no warn threshold
	StaticValue string  // used when IntervalMs == 0
	StaticLabel string  // status label for static readouts
}

var instrumentReadouts = []readoutDef{
	{Label: "CABIN PRESSURE", Unit: "kPa", Baseline: 101.3, Variance: 0.2, IntervalMs: 2500},
	{Label: "O\u2082 PARTIAL PRESS.", Unit: "kPa", Baseline: 21.3, Variance: 0.15, IntervalMs: 2500},
	{Label: "CABIN TEMP.", Unit: "\u00b0C", Baseline: 19.7, Variance: 0.25, IntervalMs: 2500},
	{Label: "REACTOR OUTPUT", Unit: "%", Baseline: 98.7, Variance: 1.0, IntervalMs: 2500, FaultBelow: 80},
	{Label: "VELOCITY", Unit: "km/s", StaticValue: "12.4", StaticLabel: "CRUISE"},
	{Label: "FUEL REMAINING", Unit: "%", Baseline: 67.4, Variance: 0.1, IntervalMs: 3000, WarnBelow: 30},
}

func (r readoutDef) isDynamic() bool { return r.IntervalMs > 0 }

func (r readoutDef) alpineData() string {
	if !r.isDynamic() {
		return ""
	}
	return fmt.Sprintf("{ v: %g }", r.Baseline)
}

func (r readoutDef) alpineInit() string {
	if !r.isDynamic() {
		return ""
	}
	return fmt.Sprintf("setInterval(()=>{ v=parseFloat((%g+(Math.random()-0.5)*%g).toFixed(1)) },%d)",
		r.Baseline, r.Variance*2, r.IntervalMs)
}

func (r readoutDef) hasThreshold() bool {
	return r.FaultBelow > 0 || r.WarnBelow > 0
}

// valueClassExpr returns the Alpine :class expression for the readout value.
func (r readoutDef) valueClassExpr() string {
	if r.FaultBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-readout-value cass-readout-value-danger' : 'cass-readout-value'", r.FaultBelow)
	}
	if r.WarnBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-readout-value cass-readout-value-danger' : 'cass-readout-value'", r.WarnBelow)
	}
	return ""
}

// statusClassExpr returns the Alpine :class expression for the status label.
func (r readoutDef) statusClassExpr() string {
	if r.FaultBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-value-danger' : 'cass-value-ok'", r.FaultBelow)
	}
	if r.WarnBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-value-warn' : 'cass-value-ok'", r.WarnBelow)
	}
	return ""
}

// statusTextExpr returns the Alpine x-text expression for the status label.
func (r readoutDef) statusTextExpr() string {
	if r.FaultBelow > 0 {
		return fmt.Sprintf("v < %g ? 'FAULT' : 'NOMINAL'", r.FaultBelow)
	}
	if r.WarnBelow > 0 {
		return fmt.Sprintf("v < %g ? 'LOW' : 'NOMINAL'", r.WarnBelow)
	}
	return ""
}

type barGaugeDef struct {
	Label      string
	AlpineVar  string
	InitialVal int
	FillClass  string
	ValueClass string
}

var barGauges = []barGaugeDef{
	{"FUEL CELLS", "fuel", 67, "", "cass-value"},
	{"OXYGEN RESERVES", "o2", 89, "cass-progress-fill-green", "cass-value"},
	{"COOLANT LEVEL", "coolant", 23, "cass-progress-fill-red", "cass-value-danger"},
	{"POWER DISTRIBUTION", "power", 78, "", "cass-value"},
}

// barGaugeAlpineData returns the x-data expression for the bar gauge panel.
func barGaugeAlpineData() string {
	var b strings.Builder
	b.WriteString("{ ")
	for i, g := range barGauges {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s: %d", g.AlpineVar, g.InitialVal)
	}
	b.WriteString(" }")
	return b.String()
}

// barGaugeRefreshExpr returns the Alpine @click expression to randomize all gauges.
func barGaugeRefreshExpr() string {
	var b strings.Builder
	for i, g := range barGauges {
		if i > 0 {
			b.WriteString("; ")
		}
		fmt.Fprintf(&b, "%s=Math.floor(Math.random()*60+30)", g.AlpineVar)
	}
	return b.String()
}

// --- Section 03: Spacing ---

type spacingStep struct {
	Px      int
	Label   string
	TwClass string
}

var spacingSteps = []spacingStep{
	{4, "4px / 0.25rem", "p-1"},
	{8, "8px / 0.5rem", "p-2"},
	{12, "12px / 0.75rem", "p-3"},
	{16, "16px / 1rem", "p-4"},
	{20, "20px / 1.25rem", "p-5"},
	{24, "24px / 1.5rem", "p-6"},
	{32, "32px / 2rem", "p-8"},
	{40, "40px / 2.5rem", "p-10"},
	{48, "48px / 3rem", "p-12"},
	{64, "64px / 4rem", "p-16"},
	{80, "80px / 5rem", "p-20"},
	{96, "96px / 6rem", "p-24"},
}

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

type personnelEntry struct {
	Name     string
	Location string
	Status   string
}

var personnelTracking = []personnelEntry{
	{"RIPLEY, E.", "DECK C — MEDICAL", "ACTIVE"},
	{"BISHOP", "SCIENCE LAB", "ACTIVE"},
	{"HICKS, D.", "ARMORY", "ACTIVE"},
}

type staticLogEntry struct {
	Timestamp string
	Code      string
	Message   string
}

var bootLogEntries = []staticLogEntry{
	{"[00:00:01]", "SYS", "WCYPD COLONY SYSTEMS v4.2.1 — BOOT SEQUENCE COMPLETE"},
	{"[00:00:03]", "NET", "NETWORK INTERFACES INITIALIZED — 4 NODES ACTIVE"},
	{"[00:00:07]", "ATM", "ATMOSPHERIC PROCESSOR — NOMINAL — 101.3 kPa"},
	{"[00:00:12]", "NAV", "NAVIGATION ARRAY — CALIBRATION COMPLETE"},
	{"[00:00:15]", "SCI", "SCIENCE LAB — ACCESS RESTRICTED — SPECIAL ORDER 937 ACTIVE"},
	{"[00:00:18]", "PWR", "POWER GRID — OUTPUT 98.7% NOMINAL"},
	{"[00:00:22]", "SEC", "MOTION SENSOR ARRAY — ARMED — 24 SECTORS ACTIVE"},
	{"[00:00:25]", "MED", "HYPERSLEEP UNITS 1-7 — OCCUPANTS STABLE"},
	{"[00:00:31]", "COM", "LONG-RANGE COMMS — SIGNAL LOCK CONFIRMED — RELAY B"},
	{"[00:00:38]", "ENG", "REACTOR TEMP — 487°C — WITHIN NOMINAL RANGE"},
	{"[00:00:44]", "SYS", "SPECIAL ORDER 937 — ACTIVATED — SCIENCE DEPT NOTIFIED"},
	{"[00:00:51]", "SEC", "BULKHEAD DOORS — ALL SEALED — OVERRIDE DISABLED"},
}

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

type docCell struct {
	Label      string
	Value      string
	ValueClass string
}

type docBlock struct {
	Title   string
	TitleBg string
	Rows    [][]docCell
	Summary *docCell
}

var docBlocks = []docBlock{
	{
		Title:   "WEYLAND-YUTANI CORPORATION // TECHNICAL OPERATIONS DIVISION",
		TitleBg: "var(--color-primary)",
		Rows: [][]docCell{
			{
				{"DOC NO.", "WY-TECH-2183-CSS-001", ""},
				{"REVISION", "D", ""},
				{"DATE", "2026-03-06", ""},
				{"CLASSIFICATION", "COMPANY CONFIDENTIAL", "cass-value-danger"},
			},
			{
				{"AUTHOR", "M.BISHOP, SYNTHETIC DIV", ""},
				{"APPROVED", "CARTER J. BURKE", ""},
				{"PROJECT", "HADLEY'S HOPE", ""},
				{"STATUS", "ACTIVE", "cass-value-ok"},
			},
		},
		Summary: &docCell{
			Label: "DESCRIPTION",
			Value: "Comprehensive technical reference and design system specification for Weyland-Yutani colony management interfaces. Covers all UI components, color systems, and interaction patterns for colony administration terminals.",
		},
	},
	{
		Title:   "COLONIAL MARINE CORPS // OPERATIONS COMMAND — MISSION ORDERS",
		TitleBg: "var(--color-text)",
		Rows: [][]docCell{
			{
				{"ORDER NO.", "CMC-OPS-2183-447", ""},
				{"PRIORITY", "CLASSIFIED", "cass-value-danger"},
				{"ISSUED", "2183-05-28", ""},
				{"EXPIRY", "ON COMPLETION", ""},
			},
			{
				{"UNIT", "2ND BN / ALPHA CO", ""},
				{"CO", "LT. GORMAN, S.", ""},
				{"OBJECTIVE", "RESCUE & INVESTIGATE", ""},
				{"STATUS", "ACTIVE", "cass-value-ok"},
			},
		},
		Summary: &docCell{
			Label: "MISSION SUMMARY",
			Value: "Proceed to LV-426, Hadley's Hope colony. Investigate loss of contact with 158 colonists. Secure perimeter, assess threat, and recover any survivors. Coordinate with Weyland-Yutani civilian representative.",
		},
	},
}

// Tab content data for navigation section.
type tabDef struct {
	ID    string
	Label string
}

var navTabs = []tabDef{
	{"overview", "OVERVIEW"},
	{"systems", "SYSTEMS"},
	{"crew", "CREW"},
	{"cargo", "CARGO"},
	{"logs", "LOGS"},
}

type navItem struct {
	ID    string
	Label string
}

var tocItems = []navItem{
	{"1", "1.0  Color Palette"},
	{"2", "2.0  Typography"},
	{"3", "3.0  Spacing Scale"},
	{"4", "4.0  Button Components"},
	{"5", "5.0  Form Components"},
	{"6", "6.0  Panel Components"},
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

// --- Section 11: Navigation (overview tab) ---

var tabOverviewParagraphs = []string{
	"The USCSS Nostromo (registration MSV-180286) is a modified Lockmart CM-88B Bison M-Class starfreighter, commissioned 2116 and currently assigned commercial haulage routes under contract ref. CMO-180286. The vessel has been diverted from its registered transit route by automated command subroutine Mother, acting under Special Order 937.",
	"All crew have been revived from hypersleep at grid reference Zeta II Reticuli. Navigation shows current position 34 light years from Earth. Mission duration estimate: 10 months transit to LV-426. All systems reading nominal except Navigation (override engaged) and Science Lab (access restricted under SO-937).",
}

// --- Section 15: Dialogs ---

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

// --- Masthead ---

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

// --- Tab logs ---

var tabLogs = []staticLogEntry{
	{"2183-06-01 04:12", "SYS", "Mother initiated hypersleep revival sequence"},
	{"2183-06-01 04:47", "NAV", "Course correction applied — new heading LV-426"},
	{"2183-06-01 05:03", "SCI", "Special Order 937 activated — Science Lab sealed"},
	{"2183-06-01 06:22", "COM", "Long-range signal detected — grid ref: Zeta Reticuli"},
	{"2183-06-01 07:55", "SYS", "Crew briefing complete — descent vehicle prepped"},
}
