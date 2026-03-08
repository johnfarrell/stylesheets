# Mission Control (OSRS Tracker) Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Add a flagship "Mission Control" style guide themed as a dark-mode OSRS goal tracker, showcasing sidebar tree navigation, HTMX search/lazy-load, Alpine.js dependency graphs, and a combat level calculator.

**Architecture:** Follows the established pattern: `guides/tracker/styles.go` + `guides/tracker/tracker.templ` + handler routes in `handlers/guides.go`. Guide-specific HTMX endpoint responses use small templ components (`search.templ`, `detail.templ`). The guide is data-dense with ~20 canned OSRS items for the interactive demo.

**Tech Stack:** Go 1.25, Templ v0.3.1001, HTMX 2.0.4, Alpine.js 3.14.9, Tailwind CSS v4

**Design doc:** `docs/plans/2026-03-07-osrs-tracker-design.md`

---

### Task 0: Register Mission Control guide (registry + embed + handler wiring)

Infrastructure setup: registry entry, embed directive, handler switch case, placeholder page.

**Files:**
- Modify: `guides/registry.go` (add tracker to `All` slice)
- Modify: `guides/sources.go` (add `tracker` to embed directive)
- Create: `guides/tracker/styles.go`
- Create: `guides/tracker/tracker.templ` (minimal placeholder)
- Modify: `handlers/guides.go` (add import + switch case)
- Modify: `templates/components/formresponse.templ` (add `tracker` case)

**Step 1: Add tracker entry to `guides/registry.go` `All` slice**

After the Newspaper entry, add:

```go
{
    Name:        "Mission Control",
    Slug:        "tracker",
    Description: "Dark-mode OSRS goal tracker — sidebar tree, dependency graphs, progress dashboards.",
    FontURL:     "https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,300;0,9..40,400;0,9..40,500;0,9..40,700&family=Space+Mono:wght@400;700&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#0d1117",
        "--color-surface":     "#161b22",
        "--color-surface-2":   "#1c2128",
        "--color-primary":     "#c8aa6e",
        "--color-secondary":   "#e6edf3",
        "--color-accent":      "#2ea043",
        "--color-warning":     "#d29922",
        "--color-danger":      "#da3633",
        "--color-info":        "#58a6ff",
        "--color-text":        "#e6edf3",
        "--color-text-muted":  "#7d8590",
        "--color-border":      "#30363d",
        "--font-display":      "'Space Mono', monospace",
        "--font-body":         "'DM Sans', sans-serif",
        "--font-mono":         "'Space Mono', monospace",
        "--font-size-display": "2.5rem",
        "--font-size-heading": "1.25rem",
        "--font-size-body":    "0.875rem",
        "--font-size-caption": "0.75rem",
        "--radius-sm":         "4px",
        "--radius-md":         "6px",
        "--radius-lg":         "8px",
        "--shadow-card":       "0 1px 4px rgba(0,0,0,0.4)",
        "--shadow-btn":        "0 1px 2px rgba(0,0,0,0.3)",
        "--border-width":      "1px",
        "--border-color":      "#30363d",
        "--content-max-width": "1200px",
        "--section-padding":   "3rem 2rem",
        "--code-bg":           "#0d1117",
        "--code-text":         "#e6edf3",
        "--code-keyword":      "#c8aa6e",
        "--code-string":       "#2ea043",
        "--code-comment":      "#7d8590",
        "--code-number":       "#d29922",
        "--code-tag":          "#58a6ff",
        "--code-attr":         "#7d8590",
    },
},
```

**Step 2: Update embed directive in `guides/sources.go`**

```go
//go:embed brutalist minimal cassette glass bento swiss terminal retro newspaper tracker
```

**Step 3: Create `guides/tracker/styles.go`**

```go
package tracker

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - dark panel with subtle border and shadow */
.trk-panel {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-card);
}
.trk-panel-header {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--color-text-muted);
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--color-border);
    border-left: 3px solid var(--color-primary);
    background: var(--color-surface-2);
}
.trk-panel-elevated {
    background: var(--color-surface-2);
}
/* [custom] - status indicator dots */
.trk-status-light {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--color-text-muted);
    flex-shrink: 0;
}
.trk-status-complete {
    background: var(--color-accent);
    box-shadow: 0 0 6px var(--color-accent);
}
.trk-status-progress {
    background: var(--color-warning);
    box-shadow: 0 0 6px var(--color-warning);
}
.trk-status-locked {
    background: var(--color-danger);
    box-shadow: 0 0 4px var(--color-danger);
}
/* [custom] - pulsing animation for in-progress items */
@keyframes trk-pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
}
.trk-status-pulse {
    animation: trk-pulse 2s ease-in-out infinite;
}
/* [custom] - progress bar with gold fill */
.trk-progress-bar {
    height: 6px;
    background: var(--color-surface-2);
    border-radius: 3px;
    overflow: hidden;
}
.trk-progress-fill {
    height: 100%;
    background: var(--color-primary);
    border-radius: 3px;
    transition: width 0.3s ease;
}
/* [custom] - sidebar tree navigation */
.trk-tree {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
}
.trk-tree-node {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.35rem 0.5rem;
    cursor: pointer;
    color: var(--color-text);
    border-left: 2px solid transparent;
    transition: background 0.1s, border-color 0.1s;
}
.trk-tree-node:hover {
    background: var(--color-surface-2);
}
.trk-tree-node-active {
    border-left-color: var(--color-primary);
    background: var(--color-surface-2);
    color: var(--color-primary);
}
.trk-tree-toggle {
    display: inline-flex;
    width: 1rem;
    justify-content: center;
    font-size: 0.625rem;
    color: var(--color-text-muted);
    transition: transform 0.15s;
    flex-shrink: 0;
    user-select: none;
}
.trk-tree-toggle-open {
    transform: rotate(90deg);
}
/* [custom] - buttons with gold border accent */
.trk-btn {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
    font-weight: 700;
    letter-spacing: 0.04em;
    color: var(--color-primary);
    background: transparent;
    border: 1px solid var(--color-primary);
    border-radius: var(--radius-sm);
    padding: 0.4rem 1rem;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
}
.trk-btn:hover {
    background: var(--color-primary);
    color: var(--color-bg);
}
.trk-btn-primary {
    background: var(--color-primary);
    color: var(--color-bg);
    border-color: var(--color-primary);
}
.trk-btn-primary:hover {
    background: #b8993e;
    border-color: #b8993e;
}
.trk-btn-danger {
    color: var(--color-danger);
    border-color: var(--color-danger);
}
.trk-btn-danger:hover {
    background: var(--color-danger);
    color: var(--color-text);
}
/* [custom] - dark inset input fields */
.trk-input {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    padding: 0.4rem 0.75rem;
    width: 100%;
    transition: border-color 0.15s;
}
.trk-input:focus {
    outline: none;
    border-color: var(--color-primary);
}
.trk-input::placeholder {
    color: var(--color-text-muted);
}
/* [custom] - search input */
.trk-search {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    padding: 0.4rem 0.75rem;
    width: 100%;
}
.trk-search:focus {
    outline: none;
    border-color: var(--color-primary);
}
.trk-search::placeholder {
    color: var(--color-text-muted);
}
/* [custom] - large monospace readout values */
.trk-readout {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 2rem;
    color: var(--color-primary);
    letter-spacing: 0.02em;
}
/* [custom] - small category tag pills */
.trk-tag {
    display: inline-block;
    font-family: var(--font-display);
    font-size: 0.625rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    padding: 0.15rem 0.5rem;
    border-radius: 2px;
    background: var(--color-surface-2);
    color: var(--color-text-muted);
    border: 1px solid var(--color-border);
}
.trk-tag-skill { color: var(--color-info); border-color: var(--color-info); }
.trk-tag-quest { color: var(--color-primary); border-color: var(--color-primary); }
.trk-tag-diary { color: var(--color-accent); border-color: var(--color-accent); }
.trk-tag-boss { color: var(--color-danger); border-color: var(--color-danger); }
/* [custom] - dependency graph nodes */
.trk-dep-node {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.75rem;
    font-size: var(--font-size-caption);
    min-width: 120px;
    position: relative;
}
.trk-dep-node-complete { border-left: 3px solid var(--color-accent); }
.trk-dep-node-progress { border-left: 3px solid var(--color-warning); }
.trk-dep-node-locked { border-left: 3px solid var(--color-danger); }
.trk-dep-node-dimmed { opacity: 0.3; }
/* [custom] - connecting lines between dep nodes */
.trk-dep-line {
    border-top: 1px dashed var(--color-border);
    width: 2rem;
    align-self: center;
    flex-shrink: 0;
}
/* [custom] - horizontal divider */
.trk-rule {
    border-top: 1px solid var(--color-border);
}
/* [custom] - gold text glow for emphasis */
.trk-glow {
    text-shadow: 0 0 8px rgba(200,170,110,0.4);
}
`
}
```

**Step 4: Create minimal `guides/tracker/tracker.templ`**

```
package tracker

import (
	"github.com/johnfarrell/stylesheets/guides"
)

// Page renders the Mission Control style guide showcase.
templ Page(g guides.Guide, htmxRequest bool) {
	if htmxRequest {
		<span id="font-loader" { templ.Attributes{"hx-swap-oob": "true"}... }>
			<link rel="preconnect" href="https://fonts.googleapis.com"/>
			<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
			<link rel="stylesheet" href={ g.FontURL }/>
		</span>
	}
	@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
	<div style="background: var(--color-bg); min-height: 100%; padding: var(--section-padding); color: var(--color-text); font-family: var(--font-body);">
		<div style="max-width: var(--content-max-width); margin: 0 auto;">
			<h1 class="trk-glow mb-2" style="font-family: var(--font-display); font-size: var(--font-size-display); font-weight: 700; color: var(--color-primary);">
				Mission Control
			</h1>
			<p style="font-size: var(--font-size-caption); color: var(--color-text-muted);">
				OSRS Goal Tracker — Dark-mode dashboard for account progression.
			</p>
		</div>
	</div>
}
```

**Step 5: Update `handlers/guides.go`**

Add import:
```go
trackertempl "github.com/johnfarrell/stylesheets/guides/tracker"
```

Add switch case in `guideContent()`:
```go
case "tracker":
    return trackertempl.Page(g, htmxRequest)
```

**Step 6: Add `tracker` case to `templates/components/formresponse.templ`**

Before the `default:` case, add:
```go
case "tracker":
    <div style="background: var(--color-surface); border: 1px solid var(--color-primary); border-radius: var(--radius-sm); padding: 0.75rem; font-family: var(--font-display); font-size: var(--font-size-caption);">
        <div style="color: var(--color-primary); font-weight: 700;">GOAL CONFIRMED</div>
        <div style="color: var(--color-text-muted); margin-top: 0.25rem;">Tracking: <strong style="color: var(--color-text);">{ name }</strong></div>
    </div>
```

**Step 7: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 8: Commit**

```bash
git add guides/registry.go guides/sources.go guides/tracker/ handlers/guides.go templates/components/formresponse.templ
git commit -m "feat: register Mission Control guide with styles and placeholder page"
```

---

### Task 1: Build full page — required sections 1-6

Build the complete required sections (Color Palette, Typography, Spacing, Buttons, Forms, Cards/Panels).

**Files:**
- Modify: `guides/tracker/tracker.templ` (replace placeholder with full sections 1-6)

**Step 1: Build the full `guides/tracker/tracker.templ` with sections 1-6**

Replace the placeholder with the full page. Include:
- `{{ snippets := guides.GetHighlightedSnippets(g.Slug) }}`
- Import `"github.com/johnfarrell/stylesheets/templates/components"` for Section, TechSummary, SourceView, badges
- Keep the OOB font loader
- Guide header styled as mission control: title "Mission Control" with gold glow, subtitle "OSRS Goal Tracker — Account Progression Dashboard"
- TechSummary callouts listing all interactive features

**Section 1 — Color Palette** [Alpine]:
- Grid of dark color swatches with hex values
- Copy-to-clipboard via Alpine
- Colors: Background #0d1117, Surface #161b22, Gold #c8aa6e, Green #2ea043, Amber #d29922, Red #da3633, Blue #58a6ff, Text #e6edf3, Muted #7d8590, Border #30363d
- Use `trk-panel` for each swatch container
- Gold border highlight on hover

**Section 2 — Typography** [None]:
- Space Mono display/heading samples shown as "mission readout" headers
- DM Sans body/caption samples
- A data panel mixing both: monospace labels ("TOTAL LEVEL:", "QUEST POINTS:") with sans-serif values
- Show all weights

**Section 3 — Spacing Scale** [None]:
- `trk-progress-bar` style bars at each spacing step (4px through 64px)
- Monospace size labels on the left
- Gold fill bars — visually reads like XP bars

**Section 4 — Buttons** [Alpine]:
- Size variants (SM, MD, LG) using `trk-btn`
- Style variants: default (outline gold), primary (solid gold), danger (outline red), disabled
- Toggle demo: "Track" / "Untrack" button. When clicked, text swaps and a `trk-status-light` goes from off to `trk-status-complete` (green glow). Snippet marker around this.

**Section 5 — Forms** [Both]:
- "Add Goal" panel with `trk-panel` + `trk-panel-header`
- Fields: Goal name (text input), Category (select: Skill, Quest, Diary, Boss, Collection), Priority (radio: Low, Normal, High), Notifications (checkbox)
- All using `trk-input` class
- `hx-post="/guides/tracker/demo-form"` with `hx-target="#trk-form-response"` and `hx-swap="innerHTML"`
- Response target div below the form

**Section 6 — Cards/Panels** [Alpine]:
- Basic `trk-panel` with header bar showing a quest summary
- Panel with `trk-status-light` indicators in a grid (mini status board)
- Expandable panel with Alpine x-show + x-transition (collapsible requirements list)
- Snippet markers around the expandable card

**Step 2: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 3: Commit**

```bash
git add guides/tracker/
git commit -m "feat: build Mission Control required sections 1-6

Color palette, typography, spacing, buttons, forms, and cards/panels
with dark-mode tracker aesthetic and gold accent."
```

---

### Task 2: Build Section 7 — Sidebar Tree + HTMX Search & Detail Loading

The flagship interactive section. Create the sidebar tree, HTMX search endpoint, HTMX detail endpoint, and supporting templ components.

**Files:**
- Modify: `guides/tracker/tracker.templ` (add section 7)
- Create: `guides/tracker/search.templ` (search result component)
- Create: `guides/tracker/detail.templ` (detail panel component)
- Modify: `handlers/guides.go` (add search + detail endpoints)
- Modify: `handlers/guides_test.go` (add tests)

**Step 1: Write tests**

Add to `handlers/guides_test.go`:

```go
func TestTrackerSearch(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/tracker/search?q=attack", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html, got %q", ct)
	}
}

func TestTrackerDetail(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/tracker/detail/skill/attack", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
```

**Step 2: Define canned OSRS data in handlers**

Create a package-level data structure in `handlers/guides.go` (or a helper near the endpoints). Include ~20 representative items:

Skills (pick ~8): Attack, Strength, Defence, Ranged, Prayer, Magic, Mining, Cooking
Quests (pick ~6): Cook's Assistant, Dragon Slayer, Monkey Madness, Recipe for Disaster, Desert Treasure, Song of the Elves
Diaries (pick ~3): Lumbridge Easy, Ardougne Medium, Karamja Elite
Bosses (pick ~3): Giant Mole, Zulrah, Vorkath

Each item has: id, category, name, status (complete/progress/locked), level/requirement, description, requirements list (as strings), unlocks list (as strings).

```go
type trackerItem struct {
    ID           string
    Category     string // "skill", "quest", "diary", "boss"
    Name         string
    Status       string // "complete", "progress", "locked"
    Level        int    // current level or 0
    Target       int    // target level or 0
    Description  string
    Requirements []string
    Unlocks      []string
}

var trackerItems = []trackerItem{
    {ID: "attack", Category: "skill", Name: "Attack", Status: "progress", Level: 75, Target: 99, Description: "Determines accuracy with melee weapons.", Requirements: nil, Unlocks: []string{"Abyssal Whip at 70", "Dragon Slayer II requirement"}},
    // ... more items
}
```

**Step 3: Create `guides/tracker/search.templ`**

```
package tracker

// SearchResult renders a single matching item in the search results.
templ SearchResult(id, category, name, status string) {
	<div class="trk-tree-node"
		hx-get={ "/guides/tracker/detail/" + category + "/" + id }
		hx-target="#trk-detail-panel"
		hx-swap="innerHTML">
		<span class={ "trk-status-light " + statusClass(status) }></span>
		<span class={ "trk-tag trk-tag-" + category }>{ category }</span>
		<span>{ name }</span>
	</div>
}
```

Add a helper function `statusClass` in tracker package (can go in styles.go or a new helpers.go):
```go
func statusClass(status string) string {
    switch status {
    case "complete": return "trk-status-complete"
    case "progress": return "trk-status-progress trk-status-pulse"
    case "locked":   return "trk-status-locked"
    default:         return ""
    }
}
```

**Step 4: Create `guides/tracker/detail.templ`**

```
package tracker

// Detail renders the full detail panel for a selected item.
templ Detail(name, category, status, description string, level, target int, requirements, unlocks []string) {
	<div>
		<div class="flex items-center gap-3 mb-4">
			<span class={ "trk-status-light " + statusClass(status) }></span>
			<h3 style="font-family: var(--font-display); font-size: var(--font-size-heading); font-weight: 700; color: var(--color-text);">{ name }</h3>
			<span class={ "trk-tag trk-tag-" + category }>{ category }</span>
		</div>
		<p style="font-size: var(--font-size-body); color: var(--color-text-muted); margin-bottom: 1rem;">{ description }</p>
		if target > 0 {
			<div class="mb-4">
				<div class="flex justify-between mb-1" style="font-family: var(--font-display); font-size: var(--font-size-caption); color: var(--color-text-muted);">
					<span>LEVEL</span>
					<span style="color: var(--color-primary);">{ fmt.Sprintf("%d / %d", level, target) }</span>
				</div>
				<div class="trk-progress-bar">
					<div class="trk-progress-fill" style={ fmt.Sprintf("width: %d%%", level*100/target) }></div>
				</div>
			</div>
		}
		if len(requirements) > 0 {
			<div class="mb-4">
				<p style="font-family: var(--font-display); font-size: var(--font-size-caption); color: var(--color-text-muted); margin-bottom: 0.5rem;">REQUIREMENTS</p>
				for _, req := range requirements {
					<div class="flex items-center gap-2 mb-1" style="font-size: var(--font-size-caption);">
						<span class="trk-status-light trk-status-complete"></span>
						<span>{ req }</span>
					</div>
				}
			</div>
		}
		if len(unlocks) > 0 {
			<div>
				<p style="font-family: var(--font-display); font-size: var(--font-size-caption); color: var(--color-text-muted); margin-bottom: 0.5rem;">UNLOCKS</p>
				for _, u := range unlocks {
					<div class="flex items-center gap-2 mb-1" style="font-size: var(--font-size-caption); color: var(--color-info);">
						<span>{ "→" }</span>
						<span>{ u }</span>
					</div>
				}
			</div>
		}
	</div>
}
```

Import `"fmt"` in the templ file.

**Step 5: Add search endpoint to `handlers/guides.go`**

```go
// Mission Control — search across OSRS items
mux.HandleFunc("/guides/tracker/search", func(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    for _, item := range trackerItems {
        if q != "" && !containsFold(item.Name, q) {
            continue
        }
        trackertempl.SearchResult(item.ID, item.Category, item.Name, item.Status).Render(r.Context(), w)
    }
})
```

**Step 6: Add detail endpoint to `handlers/guides.go`**

```go
// Mission Control — detail panel for selected item
mux.HandleFunc("/guides/tracker/detail/{category}/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    for _, item := range trackerItems {
        if item.ID == id {
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            trackertempl.Detail(item.Name, item.Category, item.Status, item.Description, item.Level, item.Target, item.Requirements, item.Unlocks).Render(r.Context(), w)
            return
        }
    }
    http.NotFound(w, r)
})
```

**Step 7: Build section 7 in tracker.templ**

Add a section after the required 6 sections. Layout:
- Two-column flex layout (sidebar 280px, detail panel flex-1)
- Top: search input with `hx-get="/guides/tracker/search?q=..."` `hx-trigger="input changed delay:300ms"` `hx-target="#trk-search-results"`
- Left: sidebar tree with collapsible categories (Alpine x-data for expand/collapse state)
- Right: detail panel `#trk-detail-panel` (empty initially, "Select an item" placeholder)
- Tree nodes have `hx-get` to load detail on click
- Snippet markers around the interactive section

The sidebar tree data is hardcoded in the templ file using Alpine x-data with nested arrays for categories. Each tree node uses `trk-tree-node` class and fires the HTMX detail load on click.

**Step 8: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 9: Commit**

```bash
git add guides/tracker/ handlers/guides.go handlers/guides_test.go
git commit -m "feat: add sidebar tree with HTMX search and detail loading

Section 7 with collapsible tree navigation, live search filter,
and lazy-loaded detail panels for OSRS items."
```

---

### Task 3: Build Section 8 — Dependency Graph

Add the Alpine-powered dependency graph visualization.

**Files:**
- Modify: `guides/tracker/tracker.templ` (add section 8)

**Step 1: Build section 8 in tracker.templ**

Add a new section for the dependency graph. Key requirements:

- Alpine `x-data` holds graph data for 3 pre-built goal chains:
  1. **Barrows Gloves** — requires: 175 Quest Points, 70 Cooking, 50+ in multiple skills, several subquests
  2. **Quest Cape** — requires: all quests complete, various skill levels
  3. **Ardougne Elite** — requires: 91 Thieving, 90 Smithing, 85 Farming, etc.

- Buttons/dropdown to switch between the 3 examples
- Horizontal left-to-right flow using flexbox
- Each node is a `trk-dep-node` with:
  - Status-colored left border (complete/progress/locked)
  - Item name
  - `trk-tag` for category
  - Level requirement if applicable
- Nodes connected by `trk-dep-line` (dashed border spans)
- Click any node to highlight its prerequisites (dim others with `trk-dep-node-dimmed`)
- Multi-level: some nodes have their own prerequisites shown in a second column

Structure the data as nested objects:
```js
{
    name: 'Barrows Gloves',
    category: 'quest',
    status: 'locked',
    deps: [
        { name: '175 Quest Points', category: 'quest', status: 'progress', deps: [] },
        { name: '70 Cooking', category: 'skill', status: 'complete', level: 70, deps: [] },
        { name: 'Desert Treasure', category: 'quest', status: 'locked', deps: [
            { name: '53 Thieving', category: 'skill', status: 'complete', level: 53, deps: [] },
            { name: '50 Magic', category: 'skill', status: 'complete', level: 50, deps: [] },
        ]},
    ]
}
```

- Snippet markers around the dependency graph section
- `@components.SourceView(snippets["dep-graph"])`

**Step 2: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 3: Commit**

```bash
git add guides/tracker/
git commit -m "feat: add dependency graph visualization

Section 8 with Alpine-powered prerequisite chain viewer for
Barrows Gloves, Quest Cape, and Ardougne Elite goals."
```

---

### Task 4: Build Section 9 — Account Overview Dashboard

Add the summary dashboard with readouts, skill grid, and combat calculator.

**Files:**
- Modify: `guides/tracker/tracker.templ` (add section 9)

**Step 1: Build section 9 in tracker.templ**

Add the final showcase section. Alpine-driven, no server endpoints needed.

**Total Level readout:**
- Large `trk-readout` showing "1847 / 2277"
- Below it, a `trk-progress-bar` showing percentage

**Quest Points readout:**
- `trk-readout` showing "198 / 300"
- Progress bar

**Completion Grid:**
- 23 small `trk-status-light` dots in a grid (one per skill)
- Color-coded: levels 1-49 = red, 50-69 = amber, 70-98 = green, 99 = gold with glow
- Each dot has a tooltip (title attribute) showing skill name + level
- Use Alpine x-data with an array of {name, level} for all 23 skills
- Sample data with varied levels to show all bracket colors

**Achievement Diary Summary:**
- 4 rows: Easy, Medium, Hard, Elite
- Each row: label + count (e.g., "8/12") + `trk-progress-bar`

**Combat Level Calculator:**
- Alpine x-data with the 7 combat skills: Attack, Strength, Defence, Hitpoints, Ranged, Prayer, Magic
- Each skill has a `trk-input` (number type, 1-99) + label
- Combat level computed reactively using the OSRS combat formula:
  ```
  base = 0.25 * (Defence + Hitpoints + floor(Prayer / 2))
  melee = 0.325 * (Attack + Strength)
  ranged = 0.325 * (floor(Ranged / 2) + Ranged)  // = 0.325 * 1.5 * Ranged
  magic = 0.325 * (floor(Magic / 2) + Magic)
  combatLevel = base + max(melee, ranged, magic)
  ```
- Display the computed level as a `trk-readout`
- Snippet markers around the calculator

**Step 2: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 3: Commit**

```bash
git add guides/tracker/
git commit -m "feat: add account overview dashboard with combat calculator

Section 9 with total level/quest point readouts, 23-skill completion
grid, diary progress, and reactive combat level calculator."
```

---

### Task 5: Final verification and tailwind rebuild

Verify everything works end-to-end.

**Files:**
- None (verification only)

**Step 1: Run full build and test**

Run: `make build && go test ./... -v`

**Step 2: Rebuild tailwind**

Run: `make tailwind`

**Step 3: Visual spot-check (manual)**

Run: `make run` — open browser to `http://localhost:8080` and verify:
- Mission Control appears in sidebar as 10th guide
- Dark background with gold accent renders correctly
- Sidebar tree expands/collapses, search filters items, detail panel loads via HTMX
- Dependency graph switches between examples, click-to-highlight works
- Combat calculator computes levels reactively
- All other guides still work

**Step 4: Commit if any fixes needed**

```bash
git commit -m "chore: final verification and tailwind rebuild for Mission Control guide"
```

---

## Task Dependency Summary

```
Task 0: Register guide (independent)
Task 1: Build sections 1-6 (depends on Task 0)
Task 2: Build section 7 — sidebar tree + HTMX (depends on Task 1)
Task 3: Build section 8 — dependency graph (depends on Task 1)
Task 4: Build section 9 — account dashboard (depends on Task 1)
Task 5: Final verification (depends on Tasks 2, 3, 4)
```

Tasks 2, 3, 4 are independent of each other (they each add a new section) but all depend on Task 1 being complete first.
