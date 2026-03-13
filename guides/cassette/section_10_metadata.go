package cassette

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
