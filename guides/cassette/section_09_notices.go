package cassette

import "strings"

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
