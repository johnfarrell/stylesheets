# Cassette Design System — Integration Guide

A guide for integrating the Cassette design language (NASA technical document aesthetic) into an existing Go + Templ application. This covers the design token system, CSS architecture, component patterns, data modeling, and interactive behavior with Alpine.js and HTMX.

## Prerequisites

- Go 1.21+ with [Templ](https://templ.guide)
- Tailwind CSS v4 (layout only — all visual styling uses CSS variables)
- Alpine.js 3.x (client-side interactivity)
- HTMX 2.x (server-driven partial updates)
- Google Fonts: **IBM Plex Mono** (body) + **Orbitron** (display headings)

```html
<link rel="preconnect" href="https://fonts.googleapis.com"/>
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:ital,wght@0,300;0,400;0,500;0,700;1,400&family=Orbitron:wght@400;700&display=swap"/>
```

---

## 1. Design Token System

The entire visual language is driven by CSS custom properties injected into `:root`. This is the foundation — every `cass-*` class and inline style references these variables. Override them to re-theme the whole UI without touching templates.

### Token Map

Inject these as a `<style>` block in your layout's `<head>`, or via `templ.Raw()` inside a Templ component:

```css
:root {
  /* Document colors */
  --color-bg:         #f5f4ef;  /* warm off-white page */
  --color-surface:    #ffffff;  /* panel/card background */
  --color-surface-2:  #e8e7e2;  /* table headers, input backgrounds */

  /* Primary palette */
  --color-primary:    #0b3d91;  /* technical blue — rules, labels, CTAs */
  --color-secondary:  #1a5276;  /* darker blue — hover states */

  /* Status colors */
  --color-danger:     #c0392b;  /* faults, critical warnings */
  --color-caution:    #c85200;  /* degraded, caution */

  /* Text & borders */
  --color-text:       #1a1a14;  /* body text */
  --color-text-muted: #5a5a52;  /* labels, metadata, captions */
  --color-border:     #c8c7c0;  /* panel borders, dividers */
  --color-rule:       #0b3d91;  /* section rule lines */

  /* Typography */
  --font-display:      'Orbitron', sans-serif;
  --font-body:         'IBM Plex Mono', monospace;
  --font-mono:         'IBM Plex Mono', monospace;
  --font-size-display: 2rem;
  --font-size-heading: 0.875rem;
  --font-size-body:    0.8125rem;
  --font-size-caption: 0.6875rem;

  /* Shape */
  --radius-sm: 0px;
  --radius-md: 2px;
  --radius-lg: 2px;

  /* Elevation */
  --shadow-card: 0 1px 3px rgba(0,0,0,0.12);
  --shadow-btn:  none;

  /* Layout */
  --layout-columns:    2;
  --layout-gap:        1.5rem;
  --content-max-width: 1200px;
  --section-padding:   2rem 2rem;
}
```

### Dark Mode

Override only the tokens that change. Typography, layout, and radius carry over:

```css
[data-theme="dark"] {
  --color-bg:         #1a1814;
  --color-surface:    #2a2820;
  --color-surface-2:  #3a3830;
  --color-primary:    #4a9eff;
  --color-secondary:  #6ba3d6;
  --color-danger:     #ff6b5b;
  --color-caution:    #f0a030;
  --color-text:       #d4d3cc;
  --color-text-muted: #8a8a82;
  --color-border:     #4a4940;
  --color-rule:       #4a9eff;
  --shadow-card:      0 1px 3px rgba(0,0,0,0.4);
  --border-color:     #4a4940;
}
```

Toggle with JS: `document.documentElement.setAttribute('data-theme', 'dark')`.

### Storing Tokens in Go

Rather than hardcoding CSS in a template, store tokens as a `map[string]string` and generate the `:root` block at render time. This lets you swap themes programmatically:

```go
type Theme struct {
    Name    string
    Tokens  map[string]string
    Dark    map[string]string // nil = no dark mode
    FontURL string
}

// BuildCSSVars generates sorted "key:value;" pairs for injection into :root.
func BuildCSSVars(vars map[string]string) string {
    keys := make([]string, 0, len(vars))
    for k := range vars {
        keys = append(keys, k)
    }
    slices.Sort(keys)
    var sb strings.Builder
    for _, k := range keys {
        sb.WriteString(k + ":" + vars[k] + ";")
    }
    return sb.String()
}
```

In your Templ layout, inject with:

```templ
@templ.Raw("<style>:root{" + BuildCSSVars(theme.Tokens) + "}</style>")
```

> **Why `templ.Raw`?** Templ treats `<style>` block content as raw text — Go expressions and loops don't work inside them. You must build the CSS string in Go and inject it.

---

## 2. CSS Classes

All custom classes are prefixed `cass-` to avoid collisions. Include this CSS in a `<style>` block or a dedicated stylesheet. Every property references CSS variables from Section 1 — no hardcoded colors.

### Body

```css
.cass-body {
  background: var(--color-bg);
  color: var(--color-text);
  font-family: var(--font-body);
  font-size: var(--font-size-body);
}
```

### Panels

Three variants: default (gray header), blue header, and danger (red border + header).

```css
.cass-panel { background: var(--color-surface); border: 1px solid var(--color-border); }
.cass-panel-header {
  background: var(--color-surface-2);
  border-bottom: 1px solid var(--color-border);
  padding: 0.4rem 0.75rem;
  font-size: var(--font-size-caption);
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-text-muted);
}
.cass-panel-header-blue { background: var(--color-primary); color: #fff; border-bottom: none; }
.cass-panel-header-danger { background: var(--color-danger); color: #fff; border-color: var(--color-danger); }
.cass-panel-danger { border-color: var(--color-danger); }
.cass-panel-body { padding: 0.75rem; }
```

### Buttons

Outline by default, fill on hover. Variants: filled, danger, danger-filled.

```css
.cass-btn {
  font-family: var(--font-body);
  font-size: var(--font-size-caption);
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border: 1px solid var(--color-primary);
  color: var(--color-primary);
  background: transparent;
  padding: 0.4rem 1rem;
  cursor: pointer;
  transition: background 0.1s, color 0.1s;
}
.cass-btn:hover { background: var(--color-primary); color: #fff; }
.cass-btn:disabled { border-color: var(--color-border); color: var(--color-text-muted); cursor: not-allowed; background: transparent; }
.cass-btn-filled { background: var(--color-primary); color: #fff; }
.cass-btn-filled:hover { background: var(--color-secondary); border-color: var(--color-secondary); }
.cass-btn-danger { border-color: var(--color-danger); color: var(--color-danger); }
.cass-btn-danger:hover { background: var(--color-danger); color: #fff; }
.cass-btn-danger-filled { background: var(--color-danger); color: #fff; border-color: var(--color-danger); }
```

### Form Inputs

Underline-style text fields with a boxed variant:

```css
.cass-input {
  font-family: var(--font-body);
  font-size: var(--font-size-body);
  background: var(--color-surface);
  border: none;
  border-bottom: 1px solid var(--color-text);
  color: var(--color-text);
  padding: 0.375rem 0;
  width: 100%;
}
.cass-input:focus { outline: none; border-bottom: 2px solid var(--color-primary); }
.cass-input::placeholder { color: var(--color-text-muted); }
.cass-input-box { border: 1px solid var(--color-border); padding: 0.375rem 0.5rem; border-bottom-width: 1px; }
.cass-input-box:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 2px rgba(11,61,145,0.12); }
.cass-label {
  font-size: var(--font-size-caption);
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--color-text-muted);
  display: block;
  margin-bottom: 0.25rem;
}
.cass-field-group { border: 1px solid var(--color-border); padding: 0.75rem; position: relative; }
.cass-field-group-label {
  position: absolute; top: -0.6rem; left: 0.5rem;
  background: var(--color-bg); padding: 0 0.25rem;
  font-size: var(--font-size-caption); font-weight: 700;
  color: var(--color-primary); letter-spacing: 0.06em; text-transform: uppercase;
}
.cass-check { width: 0.875rem; height: 0.875rem; border: 1px solid var(--color-text); accent-color: var(--color-primary); cursor: pointer; }
```

### Tables

```css
.cass-table { width: 100%; border-collapse: collapse; font-size: var(--font-size-body); }
.cass-table th {
  background: var(--color-surface-2); border: 1px solid var(--color-border);
  padding: 0.4rem 0.6rem; font-size: var(--font-size-caption); font-weight: 700;
  letter-spacing: 0.06em; text-transform: uppercase; text-align: left; color: var(--color-text-muted);
}
.cass-table td { border: 1px solid var(--color-border); padding: 0.4rem 0.6rem; vertical-align: top; }
.cass-table tr:hover td { background: rgba(11,61,145,0.03); }
.cass-table .cass-td-num { text-align: right; font-variant-numeric: tabular-nums; }
```

### Notices

Left-border alert blocks at three severity levels:

```css
.cass-notice { border-left: 4px solid; padding: 0.75rem 1rem; font-size: var(--font-size-body); }
.cass-notice-note { border-color: var(--color-primary); background: rgba(11,61,145,0.04); }
.cass-notice-caution { border-color: var(--color-caution); background: rgba(200,82,0,0.05); }
.cass-notice-warning { border-color: var(--color-danger); background: rgba(192,57,43,0.05); }
.cass-notice-label {
  font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase;
  font-size: var(--font-size-caption); margin-bottom: 0.25rem;
}
```

### Progress Bars

```css
.cass-progress-track { background: var(--color-surface-2); border: 1px solid var(--color-border); height: 1.125rem; overflow: hidden; }
.cass-progress-fill { height: 100%; background: var(--color-primary); transition: width 0.5s ease; }
.cass-progress-fill-green { background: #27ae60; }
.cass-progress-fill-red { background: var(--color-danger); }
```

### Readout Displays

Instrument-style numeric displays:

```css
.cass-readout { background: var(--color-surface-2); border: 1px solid var(--color-border); padding: 0.5rem 0.75rem; }
.cass-readout-value { font-size: 1.375rem; font-weight: 700; color: var(--color-primary); line-height: 1; font-variant-numeric: tabular-nums; }
.cass-readout-value-danger { color: var(--color-danger); }
.cass-readout-unit { font-size: var(--font-size-caption); color: var(--color-text-muted); }
```

### Status Indicator Lights

Glowing 10px dots with animation:

```css
.cass-light { width: 10px; height: 10px; border-radius: 50%; display: inline-block; flex-shrink: 0; }
.cass-light-green { background: #27ae60; box-shadow: 0 0 4px #27ae60; }
.cass-light-amber { background: #d4a017; box-shadow: 0 0 4px #d4a017; }
.cass-light-red { background: var(--color-danger); box-shadow: 0 0 4px var(--color-danger); }
.cass-light-off { background: var(--color-border); box-shadow: none; }
@keyframes cass-blink { 0%,49% { opacity: 1 } 50%,100% { opacity: 0.2 } }
.cass-blink { animation: cass-blink 1.2s step-end infinite; }
@keyframes cass-pulse { 0%,100% { opacity: 1 } 50% { opacity: 0.4 } }
.cass-pulse { animation: cass-pulse 2s ease-in-out infinite; }
```

### Modals

```css
.cass-overlay {
  position: fixed; inset: 0; background: rgba(26,26,20,0.65);
  z-index: 50; display: flex; align-items: center; justify-content: center;
}
.cass-modal { background: var(--color-surface); border: 2px solid var(--color-primary); max-width: 480px; width: 100%; }
.cass-modal-header {
  background: var(--color-primary); color: #fff; padding: 0.6rem 1rem;
  font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; font-size: var(--font-size-caption);
}
.cass-modal-header-danger { background: var(--color-danger); }
```

### Value Helpers

```css
.cass-value { font-weight: 700; color: var(--color-primary); font-variant-numeric: tabular-nums; }
.cass-value-danger { color: var(--color-danger); }
.cass-value-ok { color: #27ae60; }
.cass-value-warn { color: var(--color-caution); }
```

### Section Rules

```css
.cass-section-rule { border-top: 3px solid var(--color-primary); padding-top: 1.5rem; margin-top: 2rem; }
```

---

## 3. Templ Component Patterns

The design is composed of small, reusable Templ components. Each one accepts data via parameters and uses `{ children... }` for nested content.

### Panel (Wrapper Component)

Panels are the primary content container. The pattern uses `{ children... }` to wrap arbitrary content:

```templ
templ Panel(header string) {
    <div class="cass-panel">
        <div class="cass-panel-header">{ header }</div>
        <div class="cass-panel-body">
            { children... }
        </div>
    </div>
}

templ PanelBlue(header string) {
    <div class="cass-panel">
        <div class="cass-panel-header cass-panel-header-blue">{ header }</div>
        <div class="cass-panel-body">
            { children... }
        </div>
    </div>
}

templ PanelDanger(header string) {
    <div class="cass-panel cass-panel-danger">
        <div class="cass-panel-header cass-panel-header-danger">{ header }</div>
        <div class="cass-panel-body">
            { children... }
        </div>
    </div>
}
```

Usage:

```templ
@PanelBlue("SYSTEM CONFIGURATION") {
    <p>Panel body content goes here.</p>
}
```

### Notice (Data-Driven Component)

Notices take severity + message and derive their styling from it. The Go helper functions keep the template clean:

```go
// notice.go
const (
    SeverityNote    = "note"
    SeverityCaution = "caution"
    SeverityWarning = "warning"
)

func noticeLabelColor(severity string) string {
    switch severity {
    case SeverityCaution:
        return "var(--color-caution)"
    case SeverityWarning:
        return "var(--color-danger)"
    default:
        return "var(--color-primary)"
    }
}
```

```templ
templ Notice(severity, message string) {
    <div class={ "cass-notice cass-notice-" + severity }>
        <div class="cass-notice-label"
             style={ "color: " + noticeLabelColor(severity) + ";" }>
            { strings.ToUpper(severity) }
        </div>
        <p>{ message }</p>
    </div>
}
```

To add Alpine-powered dismissal:

```templ
templ DismissibleNotice(severity, message string) {
    <div x-data="{ visible: true }" x-show="visible"
         class={ "cass-notice cass-notice-" + severity }>
        <div class="flex justify-between items-start">
            <div class="cass-notice-label"
                 style={ "color: " + noticeLabelColor(severity) + ";" }>
                { strings.ToUpper(severity) }
            </div>
            <button { templ.Attributes{"@click": "visible=false"}... }
                    style="background:none;border:none;cursor:pointer;color:var(--color-text-muted);font-size:1rem;">
                &#215;
            </button>
        </div>
        <p>{ message }</p>
    </div>
}
```

### Progress Bar (Parameterized Component)

```templ
templ ProgressBar(label string, percent int, fillClass string) {
    <div>
        <div class="flex justify-between mb-1">
            <span style="font-size:var(--font-size-caption);font-weight:700;text-transform:uppercase;letter-spacing:0.06em;">
                { label }
            </span>
            <span class="cass-value">{ fmt.Sprintf("%d%%", percent) }</span>
        </div>
        <div class="cass-progress-track">
            if fillClass != "" {
                <div class={ "cass-progress-fill " + fillClass }
                     style={ fmt.Sprintf("width: %d%%", percent) }></div>
            } else {
                <div class="cass-progress-fill"
                     style={ fmt.Sprintf("width: %d%%", percent) }></div>
            }
        </div>
    </div>
}
```

---

## 4. Data Modeling

A key pattern in the cassette guide: **define structs in Go, attach methods that generate Alpine.js expressions, and pass them to templates.** This keeps templates declarative and pushes logic into testable Go code.

### The Readout Pattern (Recommended)

This is the most instructive example. A readout is an instrument display that can be static or dynamically updating:

```go
type Readout struct {
    Label      string
    Unit       string
    Baseline   float64  // starting value
    Variance   float64  // random fluctuation range
    IntervalMs int      // 0 = static, >0 = live update interval
    FaultBelow float64  // value below which we show FAULT (red)
    WarnBelow  float64  // value below which we show LOW (orange)
    // Static-only fields (used when IntervalMs == 0)
    StaticValue string
    StaticLabel string
}

func (r Readout) IsDynamic() bool { return r.IntervalMs > 0 }

// AlpineData returns the x-data expression: { v: 101.3 }
func (r Readout) AlpineData() string {
    return fmt.Sprintf("{ v: %g }", r.Baseline)
}

// AlpineInit returns the x-init expression that starts a setInterval loop.
func (r Readout) AlpineInit() string {
    return fmt.Sprintf(
        "setInterval(()=>{ v=parseFloat((%g+(Math.random()-0.5)*%g).toFixed(1)) },%d)",
        r.Baseline, r.Variance*2, r.IntervalMs,
    )
}

// ValueClassExpr returns the Alpine :class expression for threshold-based coloring.
func (r Readout) ValueClassExpr() string {
    if r.FaultBelow > 0 {
        return fmt.Sprintf(
            "v < %g ? 'cass-readout-value cass-readout-value-danger' : 'cass-readout-value'",
            r.FaultBelow,
        )
    }
    return ""
}
```

The template then stays minimal — it just calls the methods:

```templ
templ ReadoutView(r Readout) {
    if r.IsDynamic() {
        <div class="cass-readout" x-data={ r.AlpineData() } x-init={ r.AlpineInit() }>
            <div class="cass-readout-unit" style="margin-bottom:0.25rem;">{ r.Label }</div>
            <div class="flex items-baseline gap-1">
                <span { templ.Attributes{":class": r.ValueClassExpr()}... }
                      x-text="v">{ fmt.Sprintf("%g", r.Baseline) }</span>
                <span class="cass-readout-unit">{ r.Unit }</span>
            </div>
        </div>
    } else {
        <div class="cass-readout">
            <div class="cass-readout-unit" style="margin-bottom:0.25rem;">{ r.Label }</div>
            <div class="flex items-baseline gap-1">
                <span class="cass-readout-value">{ r.StaticValue }</span>
                <span class="cass-readout-unit">{ r.Unit }</span>
            </div>
            <div style="margin-top:0.25rem;font-size:var(--font-size-caption);font-weight:700;"
                 class="cass-value-ok">{ r.StaticLabel }</div>
        </div>
    }
}
```

**Why this pattern works:**
- The struct is your single source of truth for behavior
- Alpine expressions are built in Go — testable, type-checked
- Templates stay declarative with no string concatenation gymnastics
- Adding a new readout is just adding a struct literal to a slice

### Data Slice Pattern

Group related data as slices of structs rather than passing individual arguments:

```go
var dashboardReadouts = []Readout{
    {Label: "CABIN PRESSURE", Unit: "kPa", Baseline: 101.3, Variance: 0.2, IntervalMs: 2500},
    {Label: "REACTOR OUTPUT", Unit: "%", Baseline: 98.7, Variance: 1.0, IntervalMs: 2500, FaultBelow: 80},
    {Label: "VELOCITY", Unit: "km/s", StaticValue: "12.4", StaticLabel: "CRUISE"},
}
```

Then loop in the template:

```templ
<div class="grid grid-cols-2 md:grid-cols-3 gap-3">
    for _, r := range dashboardReadouts {
        @ReadoutView(r)
    }
</div>
```

### Key-Value Metadata

For labeled data displays:

```go
type DocCell struct {
    Label      string
    Value      string
    ValueClass string // e.g., "cass-value-ok", "cass-value-danger"
}
```

```templ
templ DocCellView(c DocCell) {
    <div class="cass-doc-cell">
        <div class="cass-doc-cell-label">{ c.Label }</div>
        if c.ValueClass != "" {
            <div class={ "cass-doc-cell-value " + c.ValueClass }>{ c.Value }</div>
        } else {
            <div class="cass-doc-cell-value">{ c.Value }</div>
        }
    </div>
}
```

---

## 5. Alpine.js Patterns

Alpine is used for all client-side interactivity. A few critical Templ syntax rules:

### Event Handlers (`@click`, `@input`, etc.)

**Never** write `@click="..."` directly on a Templ element — Templ parses `@` as a component call. Use the spread syntax:

```templ
// WRONG — templ parse error
<button @click="open = !open">Toggle</button>

// CORRECT
<button { templ.Attributes{"@click": "open = !open"}... }>Toggle</button>
```

### Reactive Bindings (`:class`, `:style`, `:disabled`)

Same spread pattern:

```templ
<div { templ.Attributes{":class": "active ? 'cass-light-green' : 'cass-light-off'"}... }></div>
<div { templ.Attributes{":style": "'width:' + percent + '%'"}... }></div>
<button { templ.Attributes{":disabled": "code !== 'CONFIRM'"}... }>Submit</button>
```

### Attributes That Work Directly

`x-data`, `x-show`, `x-text`, `x-model`, `x-init`, `x-ref`, `x-transition` — these work as normal attributes:

```templ
<div x-data="{ open: true }" x-show="open">
    <span x-text="message"></span>
</div>
```

### Status Board (State Machine)

A common pattern — cycle through states on click, update indicator lights:

```go
type SystemStatus struct {
    Name   string `json:"n"`
    Status string `json:"s"` // "ok", "warn", "err"
}

func StatusBoardJSON(systems []SystemStatus) string {
    b, _ := json.Marshal(map[string]any{"systems": systems})
    return string(b)
}
```

```templ
<div x-data={ StatusBoardJSON(systems) }>
    <template x-for="sys in systems">
        <div class="flex items-center gap-2 cursor-pointer"
             { templ.Attributes{"@click": "sys.s = sys.s==='ok'?'warn':sys.s==='warn'?'err':'ok'"}... }>
            <div { templ.Attributes{":class": `
                sys.s==='ok'   ? 'cass-light cass-light-green cass-pulse' :
                sys.s==='warn' ? 'cass-light cass-light-amber cass-blink' :
                                 'cass-light cass-light-red cass-blink'
            `}... }></div>
            <span x-text="sys.n" style="font-size:var(--font-size-caption);"></span>
        </div>
    </template>
</div>
```

### Modal Controller

Single state variable controls multiple modal types:

```templ
<div x-data="{ modal: null }">
    <button { templ.Attributes{"@click": "modal='auth'"}... } class="cass-btn">Open Auth</button>
    <button { templ.Attributes{"@click": "modal='confirm'"}... } class="cass-btn cass-btn-danger">Open Confirm</button>

    <!-- Overlay -->
    <div class="cass-overlay" x-show="modal !== null" x-transition
         { templ.Attributes{"@click.self": "modal=null"}... }>

        <!-- Auth modal -->
        <div class="cass-modal" x-show="modal === 'auth'">
            <div class="cass-modal-header">Authorization Required</div>
            <div class="cass-panel-body">
                <!-- form content -->
            </div>
        </div>

        <!-- Confirm modal -->
        <div class="cass-modal" x-show="modal === 'confirm'">
            <div class="cass-modal-header cass-modal-header-danger">Confirm Action</div>
            <div class="cass-panel-body">
                <!-- confirm content -->
            </div>
        </div>
    </div>
</div>
```

---

## 6. HTMX Patterns

### Polling (Append-Mode Event Log)

Server pushes new HTML fragments on an interval. The handler returns a rendered Templ component:

**Template:**

```templ
// LogEntry is exported — the HTTP handler renders it as a partial.
templ LogEntry(timestamp, subsystem, message string) {
    <div style="font-size:var(--font-size-caption);color:var(--color-text-muted);">
        <span>[{ timestamp }]</span>
        <span style="color:var(--color-primary);font-weight:700;margin:0 0.25rem;">{ subsystem }</span>
        <span style="color:var(--color-text);">{ message }</span>
    </div>
}
```

**Page markup:**

```templ
<div id="event-log"
     hx-get="/my-app/log"
     hx-trigger="every 4s"
     hx-swap="beforeend">
    <!-- initial static entries rendered here -->
    for _, e := range bootLog {
        @LogEntry(e.Timestamp, e.Code, e.Message)
    }
</div>
```

**Handler:**

```go
mux.HandleFunc("/my-app/log", func(w http.ResponseWriter, r *http.Request) {
    ts := time.Now().Format("15:04:05")
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // Render the exported LogEntry templ component as a partial
    mytempl.LogEntry(ts, "SYS", "Heartbeat nominal").Render(r.Context(), w)
})
```

**Key points:**
- `hx-swap="beforeend"` appends without replacing existing content
- `hx-trigger="every 4s"` polls automatically
- The handler returns a Templ component fragment, not a full page

### Form Submission

```templ
<form hx-post="/my-app/submit"
      hx-target="#form-response"
      hx-swap="innerHTML">
    <label class="cass-label">Mission ID</label>
    <input type="text" name="mission_id" class="cass-input" placeholder="Enter ID"/>
    <button type="submit" class="cass-btn cass-btn-filled mt-3">Submit</button>
</form>
<div id="form-response"></div>
```

---

## 7. Layout Structure

### Page Shell

The outermost structure for a cassette-styled page:

```templ
templ Page(theme Theme) {
    @templ.Raw("<style>:root{" + BuildCSSVars(theme.Tokens) + "}" + BuildDarkCSS(theme.Dark) + "</style>")
    @templ.Raw("<style>" + cassetteCSS + "</style>")
    <div class="cass-body min-h-screen">
        <div style="padding:var(--section-padding);max-width:var(--content-max-width);margin:0 auto;">
            { children... }
        </div>
    </div>
}
```

### Section Dividers

Sections are separated by a 3px blue rule:

```templ
templ Section(title string) {
    <section class="cass-section-rule">
        <h2 style="font-family:var(--font-display);font-size:var(--font-size-heading);font-weight:700;letter-spacing:0.08em;text-transform:uppercase;color:var(--color-primary);margin-bottom:1rem;">
            { title }
        </h2>
        { children... }
    </section>
}
```

### Masthead (Document Header)

The NASA-style title block:

```go
type Masthead struct {
    Tagline  string
    Title    string
    Subtitle string
    DocNo    string
    Revision string
    Date     string
}
```

```templ
templ MastheadView(m Masthead) {
    <div class="mb-8 pb-4" style="border-bottom:2px solid var(--color-text);">
        <div class="flex items-start justify-between">
            <div>
                <p class="font-bold tracking-widest mb-1"
                   style="color:var(--color-text-muted);font-size:var(--font-size-caption);">
                    { m.Tagline }
                </p>
                <h1 class="font-bold"
                    style="font-family:var(--font-display);font-size:var(--font-size-display);color:var(--color-primary);letter-spacing:0.05em;">
                    { m.Title }
                </h1>
                <p style="font-size:var(--font-size-body);color:var(--color-text-muted);margin-top:0.25rem;">
                    { m.Subtitle }
                </p>
            </div>
            <div class="text-right"
                 style="font-size:var(--font-size-caption);color:var(--color-text-muted);line-height:1.8;">
                <div>DOC NO. { m.DocNo }</div>
                <div>REV. { m.Revision }</div>
                <div>{ m.Date }</div>
                <div style="color:var(--color-primary);font-weight:700;margin-top:0.25rem;">
                    &#9654; APPROVED
                </div>
            </div>
        </div>
        <div class="mt-3" style="height:3px;background:var(--color-primary);"></div>
    </div>
}
```

---

## 8. Tailwind's Role

Tailwind handles **layout and spacing only** — never colors, fonts, or visual styling. This keeps the theme portable.

**Use Tailwind for:**
- `flex`, `grid`, `grid-cols-*`, `gap-*`
- `mb-*`, `mt-*`, `p-*`, spacing utilities
- `items-center`, `justify-between`, `text-right`
- Responsive breakpoints: `md:grid-cols-3`, `sm:grid-cols-2`
- `min-h-screen`, `overflow-hidden`

**Never use Tailwind for:**
- Colors (`text-blue-500`, `bg-gray-100`) — use CSS variables
- Font families or sizes — use CSS variables
- Borders or shadows — use `cass-*` classes

This separation means you can swap from Tailwind to plain CSS for layout without touching any visual styles, and you can re-theme the entire UI by changing only the CSS variable map.

---

## 9. Suggested File Organization

For integrating into an existing project, organize the cassette design assets alongside your application code:

```
yourapp/
├── ui/
│   ├── theme.go            # Theme struct, token map, BuildCSSVars()
│   ├── cassette.css        # All cass-* classes (from Section 2)
│   ├── components/
│   │   ├── panel.templ     # Panel, PanelBlue, PanelDanger
│   │   ├── notice.templ    # Notice, DismissibleNotice
│   │   ├── notice.go       # Severity constants, noticeLabelColor()
│   │   ├── progress.templ  # ProgressBar
│   │   ├── readout.templ   # ReadoutView
│   │   ├── readout.go      # Readout struct + Alpine expression methods
│   │   ├── table.templ     # Table wrapper (if needed)
│   │   ├── modal.templ     # Modal, Overlay
│   │   └── logentry.templ  # LogEntry (exported for HTMX partials)
│   └── layout.templ        # Page shell, font loading, CSS injection
├── handlers/
│   └── ...                 # Your HTTP handlers, referencing ui/components
└── main.go
```

**Key principles:**
- One file per component, `.go` companion for data structs + helper methods
- Keep exported Templ components (like `LogEntry`) that handlers render as partials
- CSS lives in one file — it's small (~90 lines) and all references CSS variables
- The theme token map lives in Go, not in CSS, so you can swap it at runtime

---

## 10. Integration Checklist

1. **Add fonts** — Google Fonts `<link>` for IBM Plex Mono + Orbitron in your `<head>`
2. **Inject CSS variables** — Build the `:root` block from a Go map via `templ.Raw()`
3. **Include `cassette.css`** — Either as a static file or embedded in a `<style>` block
4. **Add Tailwind** — For layout utilities only (grid, flex, spacing, breakpoints)
5. **Add Alpine.js** — `<script defer src="...alpinejs..."></script>` before `</head>`
6. **Add HTMX** — `<script src="...htmx..."></script>` in `<head>` (only if using server partials)
7. **Create components** — Start with Panel and Notice, add others as needed
8. **Wire handlers** — For any HTMX endpoints (log polling, form submission)
9. **Apply `.cass-body`** — On your outermost content wrapper

The design is intentionally modular. You don't need all components — pick the ones relevant to your application and ignore the rest.
