package cassette

// System status values used in section 08 (status board).
const (
	statusOK   = "ok"
	statusWarn = "warn"
	statusErr  = "err"
)

// Notice severity values used in section 09 (notices).
const (
	severityNote    = "note"
	severityCaution = "caution"
	severityWarning = "warning"
)

// Shared types used across multiple sections.

// docCell is used by sections 06 (panels), 10 (metadata), and 15 (dialogs).
type docCell struct {
	Label      string
	Value      string
	ValueClass string
}

// docBlock is used by section 10 (metadata) and embeds docCell.
type docBlock struct {
	Title   string
	TitleBg string
	Rows    [][]docCell
	Summary *docCell
}

// staticLogEntry is used by sections 11 (navigation tabs) and 14 (event log).
type staticLogEntry struct {
	Timestamp string
	Code      string
	Message   string
}
