# Stylesheets Expansion Design

Date: 2026-03-07

## Vision

Transform the style guide library into a full reference project — something a developer can open and immediately answer two questions: "What is possible with this stack?" and "How do I replicate this in my project?"

Every interactive component shows its source code on demand. Every guide declares exactly which technologies it showcases and why. New guides expand both aesthetic coverage and technique coverage simultaneously.

---

## Core Feature: Source View (View Source Toggle)

### Mechanism

Go's `//go:embed` embeds all `.templ` source files into the binary at compile time. A `ParseSnippets()` function scans embedded file content for named region markers and extracts the code between them. All snippets are parsed once at startup and cached — zero per-request overhead.

**Marker format:**

In `.templ` files:
```html
<!-- snippet:component-name -->
<div x-data="{ open: false }" @click="...">...</div>
<!-- /snippet:component-name -->
```

In `.go` files (for handler patterns):
```go
// snippet:handler-name
mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
    ...
})
// /snippet:handler-name
```

For HTMX components that span client and server, two snippets can be merged into one `SourceView` display, labeled by technology.

### New Files

- `guides/sources.go` — `//go:embed` directive + exported `SourceFS embed.FS`
- `guides/snippets.go` — `ParseSnippets(data string) map[string]string`, `LoadAll() map[string]map[string]string` (keyed by slug)

### SourceView Component

New `templates/components/sourceview.templ` — standalone, droppable anywhere in any guide's templ.

- Alpine `x-data="{open:false}"` toggle button labeled `< > View Source`
- `<pre x-show="open">` containing the snippet
- Styled entirely with guide CSS vars — looks native to every theme
- No external syntax highlighting library — monospace + color tokens only

### Section Component Update

`Section()` gains `snippets map[string]string` and `sourceKey string` parameters. When `snippets[sourceKey]` is non-empty, `SourceView` renders automatically at the bottom of that section. Static sections (Typography, Spacing) pass an empty `sourceKey` and get no toggle.

---

## TechSummary Component

New `templates/components/techsummary.templ` — renders at the top of every guide before any component sections.

Accepts `[]TechCallout` where each callout has a `Tech string` (HTMX / Alpine / Templ / CSS) and `Description string`. Renders as a structured header:

```
TECHNOLOGIES DEMONSTRATED
──────────────────────────────────────────────
[HTMX]   Polling with hx-trigger="every Ns"
          Out-of-band swaps with hx-swap-oob
[Alpine]  x-transition for blur-in animations
[CSS]     backdrop-filter: blur() — frosted glass
```

Styled with guide CSS vars. Existing guides are retrofitted with this header during the source view retrofit pass.

---

## Retrofit: Existing Guides

Infrastructure (#19) is built first and proven. All three existing guides are then retrofitted (#20) before any new guides are built.

### Brutalist
Snippet markers on: color copy-to-clipboard, button states demo, form HTMX submit, card collapse.

### Minimal
Snippet markers on: color swatches, form submit, card expand.

### Cassette
Snippet markers on: color swatches, button toggle, form submit, status board indicators, instrument readouts, progress trackers, system log HTMX polling (client `hx-trigger` + server handler merged), modal.

---

## New Style Guides

All three guides are built source-view-first — snippet markers written alongside each component, not added after.

### 1. Glassmorphism (`glass`)

**Aesthetic:** Frosted translucent panels over a deep gradient background. Visually striking, increasingly common in modern consumer UIs.

**Technologies Demonstrated:**
| Tech | Technique |
|---|---|
| Alpine | `x-transition` for blur-in panel animations |
| Alpine | Modal with animated backdrop overlay |
| CSS | `backdrop-filter: blur()` for frosted panels |
| CSS | Layered `rgba()` backgrounds + subtle border highlights |
| CSS | CSS vars for consistent frost depth across all components |

**Components:** Color palette, typography, spacing, frosted card variants, luminous buttons (states), form with HTMX submit, blurred modal overlay, gradient hero panel.

---

### 2. Bento Dashboard (`bento`)

**Aesthetic:** Variable-span grid tiles, SaaS dashboard feel. The most technique-dense guide — built to showcase HTMX polling and Alpine reactive state patterns.

**Technologies Demonstrated:**
| Tech | Technique |
|---|---|
| HTMX | `hx-trigger="every Ns"` polling for live metric updates |
| HTMX | `hx-swap="innerHTML"` for in-place number updates |
| HTMX | Form submission with inline response |
| Alpine | Reactive state shared across multiple widgets |
| CSS | CSS grid with variable-span tiles |

**Components:** Color palette, typography, spacing, metric tiles (HTMX live updates), activity feed, sparkline-style progress bars, sortable data table, status indicators, form with HTMX submit.

---

### 3. Swiss International (`swiss`)

**Aesthetic:** Helvetica-era grid discipline. Red/black/white palette. Zero decoration — whitespace, ruled lines, and typographic scale do all the work. The most timeless guide.

**Technologies Demonstrated:**
| Tech | Technique |
|---|---|
| Templ | Component composition with `{ children... }` slots |
| Alpine | Interactive typographic scale demo |
| CSS | Strict modular scale via CSS custom properties |
| CSS | Baseline grid using ruled lines as structural elements |
| HTMX | Form submission with editorial-style response |

**Components:** Color palette, typographic specimen cards (all weights/sizes), information hierarchy demo, spacing/grid layout showcase, editorial pull quotes, form with HTMX submit, cards.

---

## Implementation Order

```
Task #19 — Source view infrastructure (embed + parser + SourceView + Section update)
    ↓
Task #20 — Retrofit Brutalist, Minimal, Cassette (markers + TechSummary + SourceView drops)
    ↓
Task #21 — Build Glass, Bento, Swiss (source-view-first, TechSummary from day one)
```

---

## CSS Conventions

Existing conventions hold. New guides:
- Class prefix: `glass-`, `bento-`, `swiss-`
- All per-guide theming via CSS vars
- `backdrop-filter` requires `/* [custom] - not achievable with Tailwind utilities */` comment
- CSS grid variable-span tiles require inline style or custom class with comment

## New Handler Routes

- `GET /guides/bento/metrics` — returns updated metric tile HTML fragment (HTMX polling)
- `POST /guides/{slug}/demo-form` — existing generic handler covers all new guides
