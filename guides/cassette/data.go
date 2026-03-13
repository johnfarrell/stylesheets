package cassette

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
