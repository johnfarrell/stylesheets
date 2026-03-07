# Stylesheets Expansion Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Add a "View Source" toggle to every interactive component and three new style guides (Glassmorphism, Bento Dashboard, Swiss International), making the app a complete reference for what the stack can do and how to replicate it.

**Architecture:** Source snippets are extracted from real `.templ` and `.go` files at startup using `//go:embed` + a marker-comment parser — no hand-written duplicates, always in sync. A new `SourceView` component drops inline anywhere in a guide. A `TechSummary` component at the top of each guide declares what technologies it demonstrates. New guides are built with snippets from day one; existing guides are retrofitted first.

**Tech Stack:** Go 1.26, Templ v0.3.1001, HTMX 2.0.4, Alpine.js 3.14.9, Tailwind CSS v4, `embed.FS`

---

## Critical Rules (read before touching any file)

- **Commits:** Always `git commit -m "message"` with GPG signing. If a commit fails due to GPG, stop and ask the user — do NOT bypass with `--no-gpg-sign`
- **Templ Alpine attributes:** Never write `@click=` directly on an element. Use `{ templ.Attributes{"@click": "..."}... }` spread syntax
- **CSS injection:** Never loop inside a `<style>` block in templ. Use `@templ.Raw(...)` with Go helper functions
- **Build order:** After editing any `.templ` file run `make templ` before `go build` or `go test`
- **Test:** `go test ./...` must pass before every commit
- **Full build verification:** `make build` (templ + tailwind + go build) before final commit of each task

---

## Task 1: Source View Infrastructure

**Native task ID:** #19

**Files:**
- Create: `guides/sources.go`
- Create: `guides/snippets.go`
- Create: `guides/snippets_test.go`
- Create: `templates/components/sourceview.templ`
- Create: `templates/components/techsummary.templ`

---

### Step 1: Write failing tests for ParseSnippets

Create `guides/snippets_test.go`:

```go
package guides_test

import (
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
)

func TestParseSnippets_BasicExtraction(t *testing.T) {
	input := "before\n<!-- snippet:foo -->\n<div>hello</div>\n<!-- /snippet:foo -->\nafter"
	got := guides.ParseSnippets(input)
	want := "<div>hello</div>"
	if got["foo"] != want {
		t.Errorf("ParseSnippets[foo] = %q, want %q", got["foo"], want)
	}
}

func TestParseSnippets_MultipleRegions(t *testing.T) {
	input := `<!-- snippet:a -->
line-a
<!-- /snippet:a -->
middle
<!-- snippet:b -->
line-b
<!-- /snippet:b -->`
	got := guides.ParseSnippets(input)
	if got["a"] != "line-a" {
		t.Errorf("a = %q, want %q", got["a"], "line-a")
	}
	if got["b"] != "line-b" {
		t.Errorf("b = %q, want %q", got["b"], "line-b")
	}
}

func TestParseSnippets_GoStyleMarkers(t *testing.T) {
	input := "// snippet:handler\nfunc foo() {}\n// /snippet:handler"
	got := guides.ParseSnippets(input)
	if got["handler"] != "func foo() {}" {
		t.Errorf("handler = %q, want %q", got["handler"], "func foo() {}")
	}
}

func TestParseSnippets_MissingClose(t *testing.T) {
	// Unclosed snippet — should not panic, returns empty map
	input := "<!-- snippet:foo -->\n<div>hello</div>"
	got := guides.ParseSnippets(input)
	if _, ok := got["foo"]; ok {
		t.Error("expected no entry for unclosed snippet")
	}
}

func TestParseSnippets_EmptyInput(t *testing.T) {
	got := guides.ParseSnippets("")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestGetSnippets_ReturnsEmptyMapForUnknownSlug(t *testing.T) {
	got := guides.GetSnippets("does-not-exist")
	if got == nil {
		t.Error("GetSnippets must return non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("expected empty map for unknown slug, got %v", got)
	}
}
```

### Step 2: Run tests — verify they fail

```bash
cd /home/john/projects/stylesheets && go test ./guides/... -run TestParseSnippets -v
```

Expected: `FAIL` — `guides.ParseSnippets` undefined

### Step 3: Create `guides/sources.go` — embed directive

```go
package guides

import "embed"

// SourceFS holds the embedded .templ and .go source files for snippet extraction.
// The guides/ directory is embedded so snippet markers in .templ files are accessible at runtime.
//
//go:embed brutalist/brutalist.templ minimal/minimal.templ cassette/cassette.templ
var SourceFS embed.FS
```

### Step 4: Create `guides/snippets.go` — parser + cache

```go
package guides

import (
	"strings"
	"sync"
)

// snippetCache holds parsed snippets for all guides, keyed by slug then name.
var (
	snippetCache map[string]map[string]string
	snippetOnce  sync.Once
)

// ParseSnippets extracts named regions from source text.
// Regions are delimited by:
//   HTML: <!-- snippet:name --> ... <!-- /snippet:name -->
//   Go:   // snippet:name ... // /snippet:name
//
// Leading/trailing whitespace is trimmed from each extracted region.
// Unclosed regions are silently ignored.
func ParseSnippets(src string) map[string]string {
	result := map[string]string{}
	lines := strings.Split(src, "\n")
	var current string
	var buf strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect open marker — both HTML and Go comment styles
		if name, ok := extractOpen(trimmed); ok {
			current = name
			buf.Reset()
			continue
		}

		// Detect close marker
		if name, ok := extractClose(trimmed); ok {
			if current == name {
				result[name] = strings.TrimSpace(buf.String())
				current = ""
				buf.Reset()
			}
			continue
		}

		// Accumulate content when inside a region
		if current != "" {
			if buf.Len() > 0 {
				buf.WriteByte('\n')
			}
			buf.WriteString(line)
		}
	}

	return result
}

// extractOpen returns the snippet name if line is an open marker, else ("", false).
func extractOpen(line string) (string, bool) {
	// HTML: <!-- snippet:name -->
	if strings.HasPrefix(line, "<!-- snippet:") && strings.HasSuffix(line, "-->") {
		name := strings.TrimSuffix(strings.TrimPrefix(line, "<!-- snippet:"), " -->")
		name = strings.TrimSpace(name)
		if name != "" {
			return name, true
		}
	}
	// Go: // snippet:name
	if strings.HasPrefix(line, "// snippet:") {
		name := strings.TrimSpace(strings.TrimPrefix(line, "// snippet:"))
		if name != "" {
			return name, true
		}
	}
	return "", false
}

// extractClose returns the snippet name if line is a close marker, else ("", false).
func extractClose(line string) (string, bool) {
	// HTML: <!-- /snippet:name -->
	if strings.HasPrefix(line, "<!-- /snippet:") && strings.HasSuffix(line, "-->") {
		name := strings.TrimSuffix(strings.TrimPrefix(line, "<!-- /snippet:"), " -->")
		name = strings.TrimSpace(name)
		if name != "" {
			return name, true
		}
	}
	// Go: // /snippet:name
	if strings.HasPrefix(line, "// /snippet:") {
		name := strings.TrimSpace(strings.TrimPrefix(line, "// /snippet:"))
		if name != "" {
			return name, true
		}
	}
	return "", false
}

// loadAll parses snippets for all registered guides from SourceFS.
// Called once via sync.Once at first GetSnippets call.
func loadAll() map[string]map[string]string {
	cache := map[string]map[string]string{}
	files := map[string]string{
		"brutalist": "brutalist/brutalist.templ",
		"minimal":   "minimal/minimal.templ",
		"cassette":  "cassette/cassette.templ",
	}
	for slug, path := range files {
		data, err := SourceFS.ReadFile(path)
		if err != nil {
			// File not embedded — skip silently (guide may not exist yet)
			continue
		}
		cache[slug] = ParseSnippets(string(data))
	}
	return cache
}

// GetSnippets returns the parsed snippet map for a guide slug.
// Returns a non-nil empty map if the slug is unknown.
func GetSnippets(slug string) map[string]string {
	snippetOnce.Do(func() {
		snippetCache = loadAll()
	})
	if s, ok := snippetCache[slug]; ok {
		return s
	}
	return map[string]string{}
}
```

### Step 5: Run tests — verify they pass

```bash
cd /home/john/projects/stylesheets && go test ./guides/... -v
```

Expected: All `TestParseSnippets_*` and `TestGetSnippets_*` PASS. Existing registry tests also PASS.

### Step 6: Create `templates/components/sourceview.templ`

```go
package components

// SourceView renders an inline "View Source" toggle for a code snippet.
// Pass the snippet string from guides.GetSnippets(slug)["snippet-name"].
// Renders nothing if code is empty — safe to call unconditionally.
templ SourceView(code string) {
	if code != "" {
		<div x-data="{ open: false }" class="mt-4">
			<button
				class="flex items-center gap-1.5 text-xs font-mono px-2 py-1 border rounded transition-colors cursor-pointer"
				style="border-color: var(--color-border, #e5e7eb); color: var(--color-text-muted, #6b7280); background: var(--color-surface, #fff);"
				{ templ.Attributes{"@click": "open = !open"}... }
			>
				<span class="font-bold" style="color: var(--color-primary, #000);">&lt;/&gt;</span>
				<span x-text="open ? 'Hide Source' : 'View Source'">View Source</span>
				<span x-text="open ? '▴' : '▾'">▾</span>
			</button>
			<div x-show="open" x-cloak>
				<pre
					class="mt-2 p-4 overflow-x-auto text-xs rounded border"
					style="font-family: var(--font-body, monospace); background: var(--color-surface, #f8f8f8); border-color: var(--color-border, #e5e7eb); color: var(--color-text, #1a1a1a); line-height: 1.6; white-space: pre; tab-size: 2;"
				>{ code }</pre>
			</div>
		</div>
	}
}
```

### Step 7: Create `templates/components/techsummary.templ`

```go
package components

// TechCallout describes one technology used in a guide.
type TechCallout struct {
	Tech        string // "HTMX", "Alpine", "Templ", "CSS"
	Description string // e.g. "hx-trigger=\"every 3s\" polling for live data"
}

// TechSummary renders the "Technologies Demonstrated" header block at the
// top of a guide. Each guide passes its own []TechCallout slice.
// Renders nothing if callouts is empty.
templ TechSummary(callouts []TechCallout) {
	if len(callouts) > 0 {
		<div
			class="mb-10 p-5 border rounded"
			style="border-color: var(--color-border, #e5e7eb); background: var(--color-surface, #fff);"
		>
			<p
				class="text-xs font-bold uppercase tracking-widest mb-4"
				style="color: var(--color-text-muted, #6b7280); font-family: var(--font-body, monospace);"
			>Technologies Demonstrated</p>
			<div class="space-y-2">
				for _, c := range callouts {
					<div class="flex gap-3 text-xs" style="font-family: var(--font-body, monospace);">
						<span
							class="font-bold shrink-0 w-16"
							style="color: var(--color-primary, #000);"
						>{ "[" + c.Tech + "]" }</span>
						<span style="color: var(--color-text, #1a1a1a);">{ c.Description }</span>
					</div>
				}
			</div>
		</div>
	}
}
```

### Step 8: Generate templ files and verify build

```bash
cd /home/john/projects/stylesheets && make templ && go build ./... && go test ./...
```

Expected: All pass. Two new `*_templ.go` files generated in `templates/components/`.

### Step 9: Commit

```bash
cd /home/john/projects/stylesheets && git add guides/sources.go guides/snippets.go guides/snippets_test.go templates/components/sourceview.templ templates/components/techsummary.templ templates/components/sourceview_templ.go templates/components/techsummary_templ.go && git commit -m "feat: add source view infrastructure — embed, snippet parser, SourceView, TechSummary"
```

---

## Task 2: Retrofit Existing Guides

**Native task ID:** #20
**Blocked by:** Task 1 (#19)

**Files:**
- Modify: `guides/brutalist/brutalist.templ`
- Modify: `guides/minimal/minimal.templ`
- Modify: `guides/cassette/cassette.templ`
- Modify: `guides/sources.go` (already embeds all three — no change needed if Task 1 done correctly)

**Pattern for every guide:**

1. Add `<!-- snippet:name -->` / `<!-- /snippet:name -->` markers around each interactive component
2. At the top of `Page()`, load snippets: `{{ snippets := guides.GetSnippets(g.Slug) }}`
3. After each interactive component block, drop `@components.SourceView(snippets["name"])`
4. Add `@components.TechSummary([]components.TechCallout{...})` after the guide header `<div>`

---

### Brutalist Guide

**Snippets to mark and their names:**

| Name | What it wraps |
|---|---|
| `color-swatch` | The `colorSwatch` private templ component (the entire `templ colorSwatch(...)` block) |
| `button-toggle` | The Alpine toggle demo `<div x-data=...>` block |
| `form-htmx` | The `<form hx-post=...>` element |
| `card-collapse` | The Alpine-powered collapsible card `x-data` block |

**Step 1: Open `guides/brutalist/brutalist.templ` and add snippets**

At the top of `Page()` body (first line inside the function, before the `if htmxRequest` block — use a Go expression):

```templ
templ Page(g guides.Guide, htmxRequest bool) {
	{{ snippets := guides.GetSnippets(g.Slug) }}
	...
```

Add TechSummary after the guide header `<div class="mb-12 border-b-4...">`:

```templ
@components.TechSummary([]components.TechCallout{
    {Tech: "Alpine", Description: "Copy-to-clipboard with x-data + @click, state reset with setTimeout"},
    {Tech: "Alpine", Description: "Toggle demo — x-bind:class for conditional styling"},
    {Tech: "HTMX",   Description: "hx-post form submission with hx-target + hx-swap"},
    {Tech: "Alpine", Description: "Collapsible card panels with x-show + x-transition"},
})
```

Wrap the `colorSwatch` private component with snippet markers. Find the `templ colorSwatch(...)` definition (it's below `Page`) and add markers:

```templ
templ colorSwatch(label, cssVar, hex string) {
```

Since `colorSwatch` is a sub-component, mark its usage call site instead. In the Color Palette section, wrap the grid div:

```html
<!-- snippet:color-swatch -->
<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-4">
    @colorSwatch("Primary", "var(--color-primary)", "#000000")
    ...
</div>
<!-- /snippet:color-swatch -->
```

Then add after that section's content: `@components.SourceView(snippets["color-swatch"])`

For the button toggle, mark the `<div x-data="{ active: false }">` block:

```html
<!-- snippet:button-toggle -->
<div x-data="{ active: false }" class="flex items-center gap-6">
    <button
        class="..."
        x-bind:class="active ? 'bg-black text-white ...' : '...'"
        { templ.Attributes{"@click": "active = !active"}... }
    >
        <span x-text="active ? 'ACTIVE' : 'INACTIVE'">INACTIVE</span>
    </button>
    <p class="..." x-text="active ? 'State: ON' : 'State: OFF'">State: OFF</p>
</div>
<!-- /snippet:button-toggle -->
```

Then add: `@components.SourceView(snippets["button-toggle"])`

For the form, mark the `<form hx-post=...>` element:

```html
<!-- snippet:form-htmx -->
<form
    hx-post="/guides/brutalist/demo-form"
    hx-target="#form-response"
    hx-swap="innerHTML"
    class="space-y-6"
>
    ...
</form>
<!-- /snippet:form-htmx -->
```

Then add: `@components.SourceView(snippets["form-htmx"])`

For the card collapse section, mark the collapsible card `x-data` block (the first card that has `x-show`).

**Step 2: Verify build**

```bash
cd /home/john/projects/stylesheets && make templ && go test ./... && make build
```

Expected: All pass.

**Step 3: Manual smoke test** — run `make run`, open `http://localhost:8080/guides/brutalist`, verify:
- TechSummary block appears at top
- "View Source" buttons appear on Color Palette, Buttons, Forms, Cards sections
- Clicking toggles expands/collapses correctly
- No toggle appears on Typography or Spacing sections

**Step 4: Commit Brutalist**

```bash
cd /home/john/projects/stylesheets && git add guides/brutalist/ && git commit -m "feat(brutalist): add TechSummary and SourceView snippet toggles"
```

---

### Minimal Guide

**Snippets to mark:**

| Name | What it wraps |
|---|---|
| `color-swatch` | The color swatch grid + `minColorSwatch` calls |
| `form-htmx` | The `<form hx-post=...>` element |
| `card-expand` | The Alpine expandable card block |

Apply same pattern: `{{ snippets := guides.GetSnippets(g.Slug) }}` at top of `Page()`, add `TechSummary`, add markers, add `SourceView` drops.

```templ
@components.TechSummary([]components.TechCallout{
    {Tech: "Alpine", Description: "Copy-to-clipboard with navigator.clipboard.writeText"},
    {Tech: "HTMX",   Description: "hx-post form submission with hx-target response swap"},
    {Tech: "Alpine", Description: "x-show card expand with smooth x-transition"},
})
```

**Step 5: Verify and commit**

```bash
cd /home/john/projects/stylesheets && make templ && go test ./... && git add guides/minimal/ && git commit -m "feat(minimal): add TechSummary and SourceView snippet toggles"
```

---

### Cassette Guide

**Snippets to mark (Cassette has more interactive components):**

| Name | What it wraps |
|---|---|
| `color-swatch` | Color swatch grid |
| `button-toggle` | Alpine button state toggle |
| `form-htmx` | HTMX form |
| `status-board` | Status indicator lights Alpine block |
| `instrument-readout` | Numeric readout `setInterval` Alpine block |
| `progress-tracker` | Segmented progress bar Alpine block |
| `system-log` | The HTMX polling div (`hx-trigger="every 4s"`) |
| `modal` | The Alpine modal open/close block |

For `system-log`, the snippet should show BOTH the HTMX client attribute AND the Go handler. Combine them with a label comment:

```
<!-- CLIENT: HTMX polling -->
<div hx-get="/guides/cassette/log" hx-trigger="every 4s" hx-swap="beforeend" ...></div>

<!-- SERVER: Go handler -->
mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
    // time-based log entry selection
    fmt.Fprintf(w, `<div ...>[%s] %s %s</div>`, ts, sub, msg)
})
```

Since the handler lives in `handlers/guides.go`, create a separate snippet there using Go comment markers:

In `handlers/guides.go`, wrap the cassette log handler:
```go
// snippet:cassette-log-handler
mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
    ...
})
// /snippet:cassette-log-handler
```

Then update `guides/sources.go` to also embed the handler file:

```go
//go:embed brutalist/brutalist.templ minimal/minimal.templ cassette/cassette.templ
var SourceFS embed.FS
```

Wait — `handlers/guides.go` is outside the `guides/` package directory. `//go:embed` can only embed files within the package directory or subdirectories.

**Solution:** Copy the minimal handler snippet into the cassette templ as a HTML comment (for display purposes). In the cassette templ, the `system-log` snippet shows just the HTMX element. Add a second snippet called `cassette-log-server` as a hardcoded display-only string constant in `guides/cassette/styles.go`:

In `guides/cassette/styles.go`, add:

```go
// LogHandlerSnippet is the server-side handler shown in the system log SourceView.
// This is a display copy — the real handler lives in handlers/guides.go.
const LogHandlerSnippet = `mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
    entries := []struct{ sub, msg string }{ ... }
    idx := int(time.Now().Unix()) % len(entries)
    e := entries[idx]
    ts := time.Now().Format("15:04:05")
    fmt.Fprintf(w, ` + "`" + `<div ...>[%s] %s %s</div>` + "`" + `, ts, e.sub, e.msg)
})`
```

Then in the cassette templ, pass a combined string to SourceView:

```templ
@components.SourceView("<!-- HTMX (client) -->\n" + snippets["system-log"] + "\n\n// Go handler (server)\n" + cassette.LogHandlerSnippet)
```

Add TechSummary for Cassette:

```templ
@components.TechSummary([]components.TechCallout{
    {Tech: "Alpine", Description: "Copy-to-clipboard on color swatches"},
    {Tech: "Alpine", Description: "Button state toggle with x-bind:class"},
    {Tech: "HTMX",   Description: "Form submission with hx-post + hx-target"},
    {Tech: "Alpine", Description: "Status board lights with reactive x-data state"},
    {Tech: "Alpine", Description: "setInterval instrument readouts (live number updates)"},
    {Tech: "Alpine", Description: "Segmented progress tracker with dynamic width binding"},
    {Tech: "HTMX",   Description: "System log polling: hx-trigger=\"every 4s\" + hx-swap=\"beforeend\""},
    {Tech: "Alpine", Description: "Modal open/close with x-show + x-transition overlay"},
})
```

**Step 6: Verify and commit**

```bash
cd /home/john/projects/stylesheets && make templ && go test ./... && make build && git add guides/cassette/ && git commit -m "feat(cassette): add TechSummary and SourceView snippet toggles"
```

---

## Task 3: Three New Style Guides

**Native task ID:** #21
**Blocked by:** Tasks 1 and 2 (#19, #20)

Build each guide in order: Glass → Bento → Swiss. Each guide follows the established pattern from the existing guides. Build one at a time, verify and commit after each.

**For each guide, these files are created:**
- `guides/{slug}/{slug}.templ`
- `guides/{slug}/styles.go`
- Update `guides/registry.go` — add entry to `All` slice
- Update `guides/sources.go` — add embed path
- Update `guides/snippets.go` `loadAll()` — add slug→path entry
- Update `handlers/guides.go` — add case to `guideContent()` switch

---

### Guide A: Glassmorphism (`glass`)

**Font URL:**
```
https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700&display=swap
```

**Registry entry** (`guides/registry.go`, add to `All` slice after Minimal):

```go
{
    Name:        "Glassmorphism",
    Slug:        "glass",
    Description: "Frosted translucent panels over deep gradients. Modern, layered, luminous.",
    FontURL:     "https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#0f0f1a",
        "--color-surface":     "rgba(255,255,255,0.08)",
        "--color-surface-2":   "rgba(255,255,255,0.14)",
        "--color-primary":     "#a78bfa",
        "--color-secondary":   "#60a5fa",
        "--color-accent":      "#f472b6",
        "--color-text":        "#f1f5f9",
        "--color-text-muted":  "#94a3b8",
        "--color-border":      "rgba(255,255,255,0.15)",
        "--font-display":      "'Plus Jakarta Sans', sans-serif",
        "--font-body":         "'Plus Jakarta Sans', sans-serif",
        "--font-size-display": "3.5rem",
        "--font-size-heading": "1.5rem",
        "--font-size-body":    "1rem",
        "--font-size-caption": "0.8rem",
        "--radius-sm":         "8px",
        "--radius-md":         "16px",
        "--radius-lg":         "24px",
        "--shadow-card":       "0 8px 32px rgba(0,0,0,0.4)",
        "--shadow-btn":        "0 4px 15px rgba(167,139,250,0.3)",
        "--border-width":      "1px",
        "--border-color":      "rgba(255,255,255,0.15)",
        "--content-max-width": "1100px",
        "--section-padding":   "4rem 2rem",
        "--frost-blur":        "16px",
        "--frost-bg":          "rgba(255,255,255,0.08)",
        "--glow-primary":      "0 0 20px rgba(167,139,250,0.4)",
    },
},
```

**`guides/glass/styles.go`:**

```go
package glass

import (
	"strings"
)

func buildCSSVars(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(v)
		sb.WriteString(";")
	}
	return sb.String()
}

func guideStyles() string {
	return `
/* [custom] - backdrop-filter not achievable with Tailwind utilities */
.glass-panel {
    background: var(--frost-bg);
    backdrop-filter: blur(var(--frost-blur));
    -webkit-backdrop-filter: blur(var(--frost-blur));
    border: var(--border-width) solid var(--color-border);
    border-radius: var(--radius-md);
}
.glass-btn-primary {
    background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
    color: #fff;
    border: none;
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.2s, transform 0.1s;
}
.glass-btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
.glass-btn-ghost {
    background: var(--frost-bg);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    border: var(--border-width) solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    font-family: var(--font-body);
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
}
.glass-btn-ghost:hover { background: rgba(255,255,255,0.14); }
/* [custom] - radial gradient background not achievable with Tailwind alone */
.glass-bg {
    background: radial-gradient(ellipse at 20% 50%, rgba(167,139,250,0.15) 0%, transparent 60%),
                radial-gradient(ellipse at 80% 20%, rgba(96,165,250,0.1) 0%, transparent 50%),
                var(--color-bg);
    min-height: 100%;
}
`
}
```

**`guides/glass/glass.templ`** — full Page component:

The page wraps everything in `<div class="glass-bg">`. Structure:

1. OOB font loader (same pattern as other guides)
2. CSS injection via `@templ.Raw(...)`
3. `{{ snippets := guides.GetSnippets(g.Slug) }}`
4. Guide header with gradient text title
5. `@components.TechSummary(...)` with Glass callouts
6. Sections: Color Palette, Typography, Spacing, Buttons, Forms, Cards/Panels, Modal

**TechSummary callouts for Glass:**
```go
[]components.TechCallout{
    {Tech: "CSS",    Description: "backdrop-filter: blur() — frosted glass panels"},
    {Tech: "CSS",    Description: "rgba() layered backgrounds with border highlights"},
    {Tech: "CSS",    Description: "CSS vars for consistent frost depth (--frost-blur, --frost-bg)"},
    {Tech: "Alpine", Description: "x-transition for modal blur-in animation"},
    {Tech: "Alpine", Description: "Copy-to-clipboard on color swatches"},
    {Tech: "HTMX",   Description: "hx-post form submission with response swap"},
}
```

**Snippet markers to add in glass.templ:**

| Name | Component |
|---|---|
| `color-swatch` | Swatch grid with Alpine copy |
| `glass-panel` | A frosted card example showing the CSS class usage |
| `button-primary` | Gradient button markup |
| `modal` | Alpine modal with `x-transition` |
| `form-htmx` | HTMX form |

**Color Palette section** — swatches with Alpine copy-to-clipboard. Use hex values: Primary `#a78bfa`, Secondary `#60a5fa`, Accent `#f472b6`, BG `#0f0f1a`, Surface `rgba(255,255,255,0.08)`, Text `#f1f5f9`, Muted `#94a3b8`.

Each swatch is a glass-panel div showing the color + hex label. Alpine `@click` copies to clipboard with visual feedback.

**Typography section** — display, heading, body, caption at each weight. White/light text on dark bg.

**Spacing section** — same ruler pattern as other guides, white bars on dark.

**Buttons section** — Primary (gradient), Ghost (frosted), Disabled. Plus Alpine toggle demo.

**Forms section** — HTMX post to `/guides/glass/demo-form` (handled by generic demo-form route). Glass-styled inputs with `rgba` backgrounds and border highlights.

**Cards/Panels section** — Three frost depth variants: light (`rgba(255,255,255,0.06)`), medium (`rgba(255,255,255,0.10)`), heavy (`rgba(255,255,255,0.16)`). Each shows backdrop-filter in action.

**Modal section** — Alpine modal with `x-transition:enter` / `x-transition:leave` for smooth blur-in. The modal overlay itself is `backdrop-filter: blur(4px)` over the page content.

```templ
<!-- snippet:modal -->
<div x-data="{ open: false }">
    <button class="glass-btn-primary px-4 py-2" { templ.Attributes{"@click": "open = true"}... }>
        Open Modal
    </button>
    <!-- Overlay -->
    <div
        x-show="open"
        class="fixed inset-0 flex items-center justify-center z-50"
        style="backdrop-filter: blur(4px); -webkit-backdrop-filter: blur(4px); background: rgba(0,0,0,0.5);"
        { templ.Attributes{
            "x-transition:enter": "transition ease-out duration-200",
            "x-transition:enter-start": "opacity-0",
            "x-transition:enter-end": "opacity-100",
            "x-transition:leave": "transition ease-in duration-150",
            "x-transition:leave-start": "opacity-100",
            "x-transition:leave-end": "opacity-0",
            "@click.self": "open = false",
        }... }
    >
        <div class="glass-panel p-8 max-w-md w-full mx-4">
            <h3 class="text-lg font-semibold mb-2" style="color: var(--color-text);">Frosted Dialog</h3>
            <p class="text-sm mb-6" style="color: var(--color-text-muted);">
                This overlay uses backdrop-filter: blur() on both the modal panel and the scrim.
            </p>
            <button class="glass-btn-primary px-4 py-2 text-sm" { templ.Attributes{"@click": "open = false"}... }>Close</button>
        </div>
    </div>
</div>
<!-- /snippet:modal -->
```

**Update `guides/sources.go`** — add glass embed:
```go
//go:embed brutalist/brutalist.templ minimal/minimal.templ cassette/cassette.templ glass/glass.templ
var SourceFS embed.FS
```

**Update `guides/snippets.go` `loadAll()`** — add:
```go
"glass": "glass/glass.templ",
```

**Update `handlers/guides.go` `guideContent()` switch** — add:
```go
case "glass":
    return glasstempl.Page(g, htmxRequest)
```

And add import:
```go
glasstempl "github.com/johnfarrell/stylesheets/guides/glass"
```

**Step 1: Write registry test for glass**

Add to `guides/registry_test.go`:

```go
func TestGlassGuideRegistered(t *testing.T) {
	_, ok := guides.BySlug("glass")
	if !ok {
		t.Fatal("expected 'glass' guide to be registered")
	}
}
```

**Step 2: Run — verify fails**

```bash
cd /home/john/projects/stylesheets && go test ./guides/... -run TestGlassGuideRegistered -v
```

Expected: FAIL

**Step 3: Implement — registry entry + styles.go + glass.templ + handler + sources.go + snippets.go**

**Step 4: Build and test**

```bash
cd /home/john/projects/stylesheets && make templ && go test ./... && make build
```

Expected: All pass.

**Step 5: Commit glass**

```bash
cd /home/john/projects/stylesheets && git add guides/glass/ guides/registry.go guides/sources.go guides/snippets.go handlers/guides.go guides/registry_test.go && git commit -m "feat: add Glassmorphism style guide with frosted panels, Alpine modal, SourceView"
```

---

### Guide B: Bento Dashboard (`bento`)

**Font URL:**
```
https://fonts.googleapis.com/css2?family=DM+Sans:wght@300;400;500;700&display=swap
```

**Registry entry:**

```go
{
    Name:        "Bento Dashboard",
    Slug:        "bento",
    Description: "Variable-span grid tiles, live HTMX metrics, SaaS dashboard patterns.",
    FontURL:     "https://fonts.googleapis.com/css2?family=DM+Sans:wght@300;400;500;700&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#f8fafc",
        "--color-surface":     "#ffffff",
        "--color-surface-2":   "#f1f5f9",
        "--color-primary":     "#6366f1",
        "--color-secondary":   "#8b5cf6",
        "--color-accent":      "#10b981",
        "--color-danger":      "#ef4444",
        "--color-warning":     "#f59e0b",
        "--color-text":        "#0f172a",
        "--color-text-muted":  "#64748b",
        "--color-border":      "#e2e8f0",
        "--font-display":      "'DM Sans', sans-serif",
        "--font-body":         "'DM Sans', sans-serif",
        "--font-size-display": "2.5rem",
        "--font-size-heading": "1.25rem",
        "--font-size-body":    "0.9375rem",
        "--font-size-caption": "0.8125rem",
        "--radius-sm":         "8px",
        "--radius-md":         "12px",
        "--radius-lg":         "16px",
        "--shadow-card":       "0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.04)",
        "--shadow-btn":        "0 1px 2px rgba(0,0,0,0.06)",
        "--border-width":      "1px",
        "--border-color":      "#e2e8f0",
        "--content-max-width": "1200px",
        "--section-padding":   "3rem 2rem",
    },
},
```

**`guides/bento/styles.go`:**

```go
package bento

import "strings"

func buildCSSVars(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		sb.WriteString(k); sb.WriteString(":"); sb.WriteString(v); sb.WriteString(";")
	}
	return sb.String()
}

func guideStyles() string {
	return `
.bento-card {
    background: var(--color-surface);
    border: var(--border-width) solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-card);
    padding: 1.5rem;
}
.bento-btn {
    background: var(--color-primary);
    color: #fff;
    border: none;
    border-radius: var(--radius-sm);
    font-family: var(--font-body);
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.15s;
    padding: 0.5rem 1rem;
}
.bento-btn:hover { opacity: 0.9; }
/* [custom] - CSS grid variable-span not achievable with static Tailwind classes */
.bento-grid {
    display: grid;
    grid-template-columns: repeat(12, 1fr);
    gap: 1rem;
}
.bento-span-4 { grid-column: span 4; }
.bento-span-6 { grid-column: span 6; }
.bento-span-8 { grid-column: span 8; }
.bento-span-12 { grid-column: span 12; }
@media (max-width: 768px) {
    .bento-span-4, .bento-span-6, .bento-span-8, .bento-span-12 { grid-column: span 12; }
}
`
}
```

**New handler route** — add to `handlers/guides.go` `NewMux()`:

```go
// Bento Dashboard — live metric tiles (HTMX polling)
mux.HandleFunc("/guides/bento/metrics", func(w http.ResponseWriter, r *http.Request) {
    metrics := []struct {
        label, value, unit, trend string
    }{
        {"Active Users",   fmt.Sprintf("%d", 1200+int(time.Now().Unix())%300),  "", "↑"},
        {"Revenue",        fmt.Sprintf("$%.0fK", 48.2+float64(int(time.Now().Unix())%20)/10), "", "↑"},
        {"Error Rate",     fmt.Sprintf("%.1f%%", 0.3+float64(int(time.Now().Unix())%10)/10), "", "↓"},
        {"Response Time",  fmt.Sprintf("%dms", 120+int(time.Now().Unix())%80),  "", "→"},
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    for _, m := range metrics {
        fmt.Fprintf(w,
            `<div class="bento-card flex flex-col gap-1"><p class="text-xs" style="color:var(--color-text-muted)">%s</p><p class="text-2xl font-bold" style="color:var(--color-text)">%s <span class="text-sm font-normal">%s</span></p><p class="text-xs" style="color:var(--color-accent)">%s</p></div>`,
            templ.EscapeString(m.label), templ.EscapeString(m.value), templ.EscapeString(m.unit), templ.EscapeString(m.trend),
        )
    }
})
```

**TechSummary callouts for Bento:**
```go
[]components.TechCallout{
    {Tech: "HTMX",   Description: "hx-trigger=\"every 3s\" polling — live metric tile updates"},
    {Tech: "HTMX",   Description: "hx-swap=\"innerHTML\" for in-place number replacement"},
    {Tech: "HTMX",   Description: "hx-post form submission with inline response"},
    {Tech: "Alpine", Description: "Reactive state shared across sortable table columns"},
    {Tech: "CSS",    Description: "grid-column: span N — variable-width bento tiles"},
}
```

**Sections:**

1. **Color Palette** [Alpine] — swatches with copy
2. **Typography** — DM Sans weights
3. **Spacing** — ruler bars
4. **Metric Tiles** [HTMX] — `<div hx-get="/guides/bento/metrics" hx-trigger="every 3s" hx-swap="innerHTML">` containing 4 tiles in a 2×2 bento grid. Initial render shows static values, HTMX updates every 3s.

   ```templ
   <!-- snippet:metric-tiles -->
   <div
       hx-get="/guides/bento/metrics"
       hx-trigger="every 3s"
       hx-swap="innerHTML"
       class="bento-grid"
   >
       <!-- Initial values (replaced immediately by HTMX) -->
       <div class="bento-card bento-span-6">...</div>
       <div class="bento-card bento-span-6">...</div>
       <div class="bento-card bento-span-6">...</div>
       <div class="bento-card bento-span-6">...</div>
   </div>
   <!-- /snippet:metric-tiles -->
   ```

5. **Data Table** [Alpine] — sortable table with `x-data="{sort:'name', dir:'asc'}"`. 6 sample rows. Column header `@click` toggles sort.

   ```templ
   <!-- snippet:sortable-table -->
   <div x-data="{ sort: 'name', dir: 'asc' }" class="bento-card overflow-x-auto">
       <table class="w-full text-sm">
           <thead>
               <tr>
                   <th
                       class="text-left p-3 cursor-pointer select-none"
                       { templ.Attributes{"@click": "sort='name'; dir = dir==='asc'?'desc':'asc'"}... }
                   >Name <span x-text="sort==='name' ? (dir==='asc' ? '↑' : '↓') : '↕'"></span></th>
                   ...
               </tr>
           </thead>
           <tbody>...</tbody>
       </table>
   </div>
   <!-- /snippet:sortable-table -->
   ```

6. **Status Indicators** [Alpine] — row of status pills with Alpine toggle for each
7. **Forms** [HTMX] — post to `/guides/bento/demo-form`

**Step 1: Write registry test**

```go
func TestBentoGuideRegistered(t *testing.T) {
    _, ok := guides.BySlug("bento")
    if !ok {
        t.Fatal("expected 'bento' guide to be registered")
    }
}
```

**Step 2: Run — verify fails**

```bash
cd /home/john/projects/stylesheets && go test ./guides/... -run TestBentoGuideRegistered -v
```

**Step 3: Implement all bento files**

**Step 4: Build and test**

```bash
cd /home/john/projects/stylesheets && make templ && go test ./... && make build
```

**Step 5: Commit**

```bash
cd /home/john/projects/stylesheets && git add guides/bento/ guides/registry.go guides/sources.go guides/snippets.go handlers/guides.go guides/registry_test.go && git commit -m "feat: add Bento Dashboard style guide with HTMX polling metrics and sortable table"
```

---

### Guide C: Swiss International (`swiss`)

**Font URL:**
```
https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@300;400;500;700&display=swap
```

**Registry entry:**

```go
{
    Name:        "Swiss International",
    Slug:        "swiss",
    Description: "Helvetica-era grid discipline. Red/black/white. Typography as the only decoration.",
    FontURL:     "https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@300;400;500;700&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#ffffff",
        "--color-surface":     "#ffffff",
        "--color-surface-2":   "#f5f5f5",
        "--color-primary":     "#e63329",
        "--color-secondary":   "#1a1a1a",
        "--color-text":        "#1a1a1a",
        "--color-text-muted":  "#767676",
        "--color-border":      "#1a1a1a",
        "--font-display":      "'IBM Plex Sans', sans-serif",
        "--font-body":         "'IBM Plex Sans', sans-serif",
        "--font-size-display": "5rem",
        "--font-size-heading": "2rem",
        "--font-size-body":    "1rem",
        "--font-size-caption": "0.75rem",
        "--radius-sm":         "0px",
        "--radius-md":         "0px",
        "--radius-lg":         "0px",
        "--shadow-card":       "none",
        "--shadow-btn":        "none",
        "--border-width":      "2px",
        "--border-color":      "#1a1a1a",
        "--content-max-width": "1100px",
        "--section-padding":   "5rem 2rem",
        "--grid-columns":      "12",
        "--baseline":          "8px",
    },
},
```

**`guides/swiss/styles.go`:**

```go
package swiss

import "strings"

func buildCSSVars(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		sb.WriteString(k); sb.WriteString(":"); sb.WriteString(v); sb.WriteString(";")
	}
	return sb.String()
}

func guideStyles() string {
	return `
/* [custom] - strict typographic grid not achievable with Tailwind utilities alone */
.swiss-rule { border-top: 2px solid var(--color-border); }
.swiss-rule-red { border-top: 2px solid var(--color-primary); }
.swiss-rule-thin { border-top: 1px solid var(--color-border); }
.swiss-label {
    font-family: var(--font-body);
    font-size: 0.625rem;
    font-weight: 700;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: var(--color-text-muted);
}
.swiss-btn {
    background: var(--color-secondary);
    color: #fff;
    border: 2px solid var(--color-secondary);
    font-family: var(--font-body);
    font-weight: 700;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
}
.swiss-btn:hover { background: var(--color-primary); border-color: var(--color-primary); }
.swiss-btn-outline {
    background: transparent;
    color: var(--color-secondary);
    border: 2px solid var(--color-secondary);
    font-family: var(--font-body);
    font-weight: 700;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
}
.swiss-btn-outline:hover { background: var(--color-secondary); color: #fff; }
/* [custom] - CSS grid strict column layout */
.swiss-grid {
    display: grid;
    grid-template-columns: repeat(12, 1fr);
    gap: 0;
}
.swiss-col-4 { grid-column: span 4; }
.swiss-col-6 { grid-column: span 6; }
.swiss-col-8 { grid-column: span 8; }
@media (max-width: 768px) {
    .swiss-col-4, .swiss-col-6, .swiss-col-8 { grid-column: span 12; }
}
`
}
```

**TechSummary callouts for Swiss:**
```go
[]components.TechCallout{
    {Tech: "Templ",  Description: "children... slot composition — reusable editorial card components"},
    {Tech: "Alpine", Description: "Interactive typographic scale demo with live size adjustment"},
    {Tech: "CSS",    Description: "Strict modular scale via --font-size-* CSS custom properties"},
    {Tech: "CSS",    Description: "Baseline grid using ruled lines (border-top) as structural elements"},
    {Tech: "HTMX",   Description: "hx-post form submission with editorial-style inline response"},
}
```

**Sections:**

1. **Color Palette** [Alpine] — 3 colors only: Red primary, Black, White. Each swatch with copy.
2. **Typography** [Alpine] — Specimen cards at all weights. PLUS an interactive scale demo: Alpine `x-data="{size: 1}"` with a slider (`<input type="range">`) that adjusts a `style` binding for font-size. Shows CSS var-driven modular scale in action.

   ```templ
   <!-- snippet:type-scale-demo -->
   <div x-data="{ size: 1 }" class="border-t-2 border-black pt-6 mt-6">
       <div class="flex items-center gap-4 mb-4">
           <span class="swiss-label">Scale Factor</span>
           <input type="range" min="0.5" max="3" step="0.1" x-model="size" class="w-48"/>
           <span class="swiss-label" x-text="size + 'x'">1x</span>
       </div>
       <p
           class="font-bold leading-tight"
           style="font-family: var(--font-display);"
           :style="'font-size: calc(var(--font-size-display) * ' + size + ')'"
       >Helvetica Neue</p>
   </div>
   <!-- /snippet:type-scale-demo -->
   ```

3. **Spacing / Grid** — Visual demonstration of the 12-column strict grid using `swiss-grid` CSS class. Show column spans (4, 6, 8, 12) with filled red/black cells.

4. **Typographic Hierarchy** [Templ] — This section showcases Templ component composition. Define a private `templ articleCard(eyebrow, headline, body string)` component that wraps text in a structured layout. Call it multiple times with different content. The snippet shows the component definition + usage.

   ```templ
   <!-- snippet:templ-composition -->
   @articleCard("Design Systems", "Grid as Foundation", "The grid is not a cage — it is a liberation...")
   @articleCard("Typography", "Weight Creates Hierarchy", "Bold speaks first. Light whispers. ...")
   <!-- /snippet:templ-composition -->
   ```

5. **Pull Quotes** — Editorial-style pull quotes with large red left border rule and oversized quotation mark.
6. **Forms** [HTMX] — Clean Swiss-styled form posting to `/guides/swiss/demo-form`.
7. **Cards** [Alpine] — Bordered cards (no radius, no shadow — just `border: 2px solid black`). Alpine expand/collapse.

**Step 1: Write registry test**

```go
func TestSwissGuideRegistered(t *testing.T) {
    _, ok := guides.BySlug("swiss")
    if !ok {
        t.Fatal("expected 'swiss' guide to be registered")
    }
}
```

**Step 2: Run — verify fails**

```bash
cd /home/john/projects/stylesheets && go test ./guides/... -run TestSwissGuideRegistered -v
```

**Step 3: Implement all swiss files**

**Step 4: Final build and full test suite**

```bash
cd /home/john/projects/stylesheets && make templ && go test ./... && make build
```

Expected: All pass. Six guides registered. All SourceView toggles working.

**Step 5: Commit Swiss and final**

```bash
cd /home/john/projects/stylesheets && git add guides/swiss/ guides/registry.go guides/sources.go guides/snippets.go handlers/guides.go guides/registry_test.go && git commit -m "feat: add Swiss International style guide with Templ composition and typographic scale demo"
```
