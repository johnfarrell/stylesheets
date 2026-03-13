# Cassette Local Component Extraction — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Extract 6+ repeated markup patterns from `guides/cassette/cassette.templ` into local templ component files within the `guides/cassette/` package, reducing duplication and improving readability.

**Architecture:** Each component file is a `.templ` file in `guides/cassette/` with unexported templ components. `cassette.templ` calls these components in place of the duplicated inline markup. All components share package `cassette`.

**Tech Stack:** Go 1.25, Templ v0.3.1001, Alpine.js (for interactive attributes)

---

### Task 1: Create `swatches.templ` — color swatch components

**Files:**
- Create: `guides/cassette/swatches.templ`
- Modify: `guides/cassette/cassette.templ:50-137` (Section 1: Color Palette)

**Step 1: Create `guides/cassette/swatches.templ`**

```templ
package cassette

// swatch renders a single clickable color swatch with Alpine copy-to-clipboard.
templ swatch(cssVar, description string) {
	<div class="cursor-pointer" x-data={ "colorSwatch('" + cssVar + "')" } { templ.Attributes{"@click": "copy()"}... }>
		<div style={ "height: 48px; background: var(" + cssVar + "); border: 1px solid var(--color-border);" }></div>
		<p style="font-size: var(--font-size-caption); font-weight: 700; margin-top: 0.25rem;">{ cssVar }</p>
		<p style="font-size: var(--font-size-caption); color: var(--color-text-muted);" x-text="copied ? 'COPIED' : hex"></p>
		<p style="font-size: var(--font-size-caption); color: var(--color-text-muted);">{ description }</p>
	</div>
}

// swatchGroup renders a labeled category header with a grid container for swatches.
templ swatchGroup(title, cols string) {
	<div class="mb-4">
		<p class="mb-2" style="font-size: var(--font-size-caption); color: var(--color-text-muted); font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase;">{ title }</p>
		<div class={ "grid gap-3 mb-6 " + cols }>
			{ children... }
		</div>
	</div>
}
```

**Step 2: Replace Section 1 swatch markup in `cassette.templ`**

Replace the entire color palette section body (lines 50-137) with calls to `swatchGroup` and `swatch`. Example:

```templ
<!-- snippet:color-swatch -->
<div>
	@swatchGroup("DOCUMENT COLORS", "grid-cols-2 sm:grid-cols-3") {
		@swatch("--color-bg", "Page Background")
		@swatch("--color-surface", "Panel / Card Surface")
		@swatch("--color-surface-2", "Table Headers / Inputs")
	}
	@swatchGroup("PRIMARY PALETTE", "grid-cols-2 sm:grid-cols-3") {
		@swatch("--color-primary", "Primary / Rules / Labels")
		@swatch("--color-secondary", "Hover / Secondary Actions")
	}
	@swatchGroup("STATUS COLORS", "grid-cols-2 sm:grid-cols-3") {
		@swatch("--color-danger", "Fault / Warning / KIA")
		@swatch("--color-caution", "Caution / Degraded")
	}
	@swatchGroup("TEXT & BORDERS", "grid-cols-2 sm:grid-cols-4") {
		@swatch("--color-text", "Body Text")
		@swatch("--color-text-muted", "Labels / Metadata")
		@swatch("--color-border", "Borders / Dividers")
		@swatch("--color-rule", "Section Rules")
	}
	<p style="font-size: var(--font-size-caption); color: var(--color-text-muted); margin-top: 0.75rem;">Click any swatch to copy hex value to clipboard.</p>
</div>
<!-- /snippet:color-swatch -->
```

**Step 3: Run templ generate and verify build**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add guides/cassette/swatches.templ guides/cassette/swatches_templ.go guides/cassette/cassette.templ guides/cassette/cassette_templ.go
git commit -m "Extract color swatch components in cassette guide"
```

---

### Task 2: Create `typography.templ` — type specimen row component

**Files:**
- Create: `guides/cassette/typography.templ`
- Modify: `guides/cassette/cassette.templ:150-198` (Section 2: Typography)

**Step 1: Create `guides/cassette/typography.templ`**

```templ
package cassette

// typeSpecimenRow renders one row in the type specimen table.
// last controls whether the bottom border is shown (false for the final row).
templ typeSpecimenRow(fontName, sizeLabel, sampleText, style string, last bool) {
	if last {
		<div style={ "padding: 0.75rem; display: grid; grid-template-columns: 200px 1fr; gap: 1rem; align-items: center;" }>
			<div style="font-size: var(--font-size-caption); color: var(--color-text-muted);">
				<div>{ fontName }</div>
				<div>{ sizeLabel }</div>
			</div>
			<div style={ style }>{ sampleText }</div>
		</div>
	} else {
		<div style={ "border-bottom: 1px solid var(--color-border); padding: 0.75rem; display: grid; grid-template-columns: 200px 1fr; gap: 1rem; align-items: center;" }>
			<div style="font-size: var(--font-size-caption); color: var(--color-text-muted);">
				<div>{ fontName }</div>
				<div>{ sizeLabel }</div>
			</div>
			<div style={ style }>{ sampleText }</div>
		</div>
	}
}
```

**Step 2: Replace Section 2 specimen rows in `cassette.templ`**

Replace lines 150-198 with calls to `typeSpecimenRow`. Example:

```templ
@typeSpecimenRow("Orbitron 700", "2rem / Display", "TECHNICAL REFERENCE MANUAL", "font-family: var(--font-display); font-size: 2rem; font-weight: 700; color: var(--color-primary);", false)
@typeSpecimenRow("IBM Plex Mono 700", "0.875rem / Heading", "SYSTEM STATUS: ALL SUBSYSTEMS NOMINAL", "font-size: 0.875rem; font-weight: 700; letter-spacing: 0.02em;", false)
// ... remaining 5 rows ...
@typeSpecimenRow("IBM Plex Mono 700", "0.6875rem / Caption Uppercase", "CLASSIFICATION: COMPANY CONFIDENTIAL", "font-size: 0.6875rem; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: var(--color-text-muted);", true)
```

**Step 3: Run templ generate and verify build**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add guides/cassette/typography.templ guides/cassette/typography_templ.go guides/cassette/cassette.templ guides/cassette/cassette_templ.go
git commit -m "Extract typography specimen row component in cassette guide"
```

---

### Task 3: Create `datarows.templ` — key-value data row components

**Files:**
- Create: `guides/cassette/datarows.templ`
- Modify: `guides/cassette/cassette.templ` — Sections 6 (panels, lines 467-517 and 562-573)

**Step 1: Create `guides/cassette/datarows.templ`**

```templ
package cassette

// dataRow renders a label-value pair with standard cassette styling.
templ dataRow(label, value string) {
	<div class="flex justify-between" style="padding: 0.25rem 0; border-bottom: 1px solid var(--color-surface-2);">
		<span style="font-size: var(--font-size-caption); color: var(--color-text-muted); font-weight: 700; text-transform: uppercase;">{ label }</span>
		<span style="font-size: var(--font-size-body);">{ value }</span>
	</div>
}

// dataRowColored renders a label-value pair where the value has a custom color.
// If color is empty, falls back to var(--color-text).
templ dataRowColored(label, value, color string) {
	<div class="flex justify-between" style="padding: 0.25rem 0; border-bottom: 1px solid var(--color-surface-2);">
		<span style="font-size: var(--font-size-caption); color: var(--color-text-muted); font-weight: 700; text-transform: uppercase;">{ label }</span>
		if color != "" {
			<span style={ "font-size: var(--font-size-body); font-weight: 700; color: " + color }>{ value }</span>
		} else {
			<span style="font-size: var(--font-size-body); font-weight: 700; color: var(--color-text);">{ value }</span>
		}
	</div>
}
```

**Step 2: Replace data row loops in `cassette.templ`**

Replace the 4 `for` loops in Section 6 with direct calls. Example for the first panel (lines 467-478):

```templ
<div style="padding: 0.75rem;">
	@dataRow("VESSEL", "USCSS NOSTROMO")
	@dataRow("REGISTRATION", "MSV-180286")
	@dataRow("DESTINATION", "LV-426 / ZETA RETICULI")
	@dataRow("MISSION PHASE", "APPROACH VECTOR")
	@dataRow("ETD", "2183-06-01 0600 UTC")
</div>
```

For the danger panel (lines 562-573), use `dataRowColored`:

```templ
<div style="padding: 0.75rem;">
	@dataRowColored("THREAT LEVEL", "EXTREME", "#c0392b")
	@dataRowColored("LOCATION", "DECK C — MEDICAL BAY", "")
	@dataRowColored("PERSONNEL AT RISK", "7 CREW MEMBERS", "#c0392b")
	@dataRowColored("RECOMMENDED ACTION", "INITIATE PROTOCOL", "")
</div>
```

**Step 3: Run templ generate and verify build**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add guides/cassette/datarows.templ guides/cassette/datarows_templ.go guides/cassette/cassette.templ guides/cassette/cassette_templ.go
git commit -m "Extract data row components in cassette guide"
```

---

### Task 4: Create `notices.templ` — dismissible notice component

**Files:**
- Create: `guides/cassette/notices.templ`
- Modify: `guides/cassette/cassette.templ:794-814` (Section 9: dismissible variants)

**Step 1: Create `guides/cassette/notices.templ`**

The severity maps to:
- `"note"` → class `cass-notice-note`, label color `var(--color-primary)`, label text `NOTE`
- `"caution"` → class `cass-notice-caution`, label color `var(--color-caution)`, label text `CAUTION`
- `"warning"` → class `cass-notice-warning`, label color `var(--color-danger)`, label text `WARNING`

```templ
package cassette

import "strings"

// dismissibleNotice renders a notice with an Alpine dismiss button.
// severity must be "note", "caution", or "warning".
templ dismissibleNotice(severity, message string) {
	<div x-data="{ visible: true }" x-show="visible" class={ "cass-notice cass-notice-" + severity }>
		<div class="flex justify-between items-start">
			<div class="cass-notice-label" style={ "color: " + noticeLabelColor(severity) + ";" }>{ strings.ToUpper(severity) }</div>
			<button { templ.Attributes{"@click": "visible=false"}... } style="background: none; border: none; cursor: pointer; color: var(--color-text-muted); font-size: 1rem; line-height: 1; padding: 0;">&#215;</button>
		</div>
		<p>{ message }</p>
	</div>
}
```

Add a helper function (can go in `notices.templ` as a Go function or in `styles.go`):

```go
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
```

**Step 2: Replace Section 9 dismissible blocks in `cassette.templ`**

Replace lines 794-814 with:

```templ
@dismissibleNotice("note", "Hypersleep revival protocol requires a minimum 4-hour monitoring period. Medical Officer must be present for all revivals. Reference: MED-HS-001.")
@dismissibleNotice("caution", "Motion sensor array sector 7G reporting intermittent signal loss. Maintenance crew dispatched. Estimated resolution: 2 hours. Do not rely on sector 7G coverage during this period.")
@dismissibleNotice("warning", "AIRLOCK CYCLE DETECTED — DECK A EMERGENCY AIRLOCK. No crew authorization on record. Investigate immediately. Security to report to Deck A airlock station.")
```

**Step 3: Run templ generate and verify build**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add guides/cassette/notices.templ guides/cassette/notices_templ.go guides/cassette/cassette.templ guides/cassette/cassette_templ.go
git commit -m "Extract dismissible notice component in cassette guide"
```

---

### Task 5: Create `gauges.templ` — bar gauge component

**Files:**
- Create: `guides/cassette/gauges.templ`
- Modify: `guides/cassette/cassette.templ:1122-1157` (Section 12: Bar Gauges)

**Step 1: Create `guides/cassette/gauges.templ`**

```templ
package cassette

// barGauge renders a labeled progress bar bound to an Alpine reactive variable.
templ barGauge(label, alpineVar, fillClass, valueClass string) {
	<div>
		<div class="flex justify-between mb-1">
			<span style="font-size: var(--font-size-caption); font-weight: 700; text-transform: uppercase; letter-spacing: 0.06em;">{ label }</span>
			<span class={ valueClass } x-text={ alpineVar + " + '%'" }></span>
		</div>
		<div class="cass-progress-track">
			<div class={ "cass-progress-fill " + fillClass } { templ.Attributes{":style": "'width:' + " + alpineVar + " + '%'"}... }></div>
		</div>
	</div>
}
```

**Step 2: Replace Section 12 bar gauge blocks in `cassette.templ`**

Replace the 4 gauge blocks (lines 1122-1157) with:

```templ
@barGauge("FUEL CELLS", "fuel", "", "cass-value")
@barGauge("OXYGEN RESERVES", "o2", "cass-progress-fill-green", "cass-value")
@barGauge("COOLANT LEVEL", "coolant", "cass-progress-fill-red", "cass-value-danger")
@barGauge("POWER DISTRIBUTION", "power", "", "cass-value")
```

Note: The coolant gauge has `style="font-weight: 700;"` on the value span in the original. Add `font-weight: 700` to the valueClass-based styling. Since `cass-value-danger` doesn't include `font-weight`, and other `.cass-value` does, check if a wrapper or additional inline style is needed. If the existing classes already handle it (`.cass-value` has `font-weight: 700`), no change needed. `.cass-value-danger` only sets color, so the coolant value needs both: use the `x-text` with inline style. Simplest: add `style="font-weight: 700;"` to the span always since all gauges show bold values.

**Step 3: Run templ generate and verify build**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add guides/cassette/gauges.templ guides/cassette/gauges_templ.go guides/cassette/cassette.templ guides/cassette/cassette_templ.go
git commit -m "Extract bar gauge component in cassette guide"
```

---

### Task 6: Create `logentry.templ` — log entry component + convert static entries to loop

**Files:**
- Create: `guides/cassette/logentry.templ`
- Modify: `guides/cassette/cassette.templ:1371-1430` (Section 14: static log entries)

**Step 1: Create `guides/cassette/logentry.templ`**

```templ
package cassette

// logEntry renders a single static log line in the system event log.
templ logEntry(ts, code, msg string) {
	<div class="flex gap-3" style="font-size: var(--font-size-caption); color: var(--color-text-muted); padding: 2px 0; border-bottom: 1px solid var(--color-surface-2);">
		<span style="min-width: 4.5rem;">{ ts }</span>
		<span style="color: var(--color-primary); font-weight: 700; min-width: 2.5rem;">{ code }</span>
		<span style="color: var(--color-text);">{ msg }</span>
	</div>
}
```

**Step 2: Replace the 11 static log entries with a `for` loop**

Replace lines 1371-1430 with:

```templ
for _, entry := range []struct{ ts, code, msg string }{
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
} {
	@logEntry(entry.ts, entry.code, entry.msg)
}
```

**Step 3: Run templ generate and verify build**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add guides/cassette/logentry.templ guides/cassette/logentry_templ.go guides/cassette/cassette.templ guides/cassette/cassette_templ.go
git commit -m "Extract log entry component and convert static entries to loop in cassette guide"
```

---

### Task 7: Final review — scan for additional optimization opportunities

**Step 1: Re-read the modified `cassette.templ` end-to-end**

Look for any remaining patterns that appear 3+ times and could benefit from extraction.

**Step 2: Run full build and visual verification**

Run: `make build`
Run: `make run` (verify in browser that all 15 sections render identically)

**Step 3: Run linter**

Run: `golangci-lint run`

**Step 4: Final commit if any additional optimizations were made**

```bash
git add -A guides/cassette/
git commit -m "Additional cassette template optimizations"
```
