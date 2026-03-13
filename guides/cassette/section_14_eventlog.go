package cassette

import "github.com/johnfarrell/stylesheets/guides"

// logHandlerSnippet is the server-side handler shown in the system log SourceView.
// Display copy only — the live handler is in handlers/guides.go.
const logHandlerSnippet = `mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
    entries := []struct{ sub, msg string }{
        {"SYS", "WCYPD COLONY SYSTEMS — HEARTBEAT NOMINAL"},
        // ... more entries ...
    }
    idx := int(time.Now().Unix()) % len(entries)
    e := entries[idx]
    ts := time.Now().Format("15:04:05")
    fmt.Fprintf(w,
        ` + "`" + `<div ...>[%s] %s %s</div>` + "`" + `,
        ts, templ.EscapeString(e.sub), templ.EscapeString(e.msg),
    )
})`

// highlightedLogHandler is logHandlerSnippet pre-highlighted as Go, cached at startup.
var highlightedLogHandler = guides.Highlight(logHandlerSnippet, "go")

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
