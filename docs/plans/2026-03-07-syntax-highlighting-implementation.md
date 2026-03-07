# Syntax Highlighting Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:subagent-driven-development to implement this plan task-by-task.

**Goal:** Add token-accurate, per-guide syntax highlighting to every SourceView code block using chroma server-side with CSS vars.

**Architecture:** Chroma parses snippets into tokens at startup (cached via `sync.Once`), outputs `<span class="k">` etc. (no surrounding `<pre>`). Each guide defines `--code-*` CSS vars. One shared CSS block in `input.css` maps chroma class names to those vars.

**Tech Stack:** `github.com/alecthomas/chroma/v2` (already in go.mod), Go, Templ, CSS custom properties.

---

## Context for implementer

The project is a style guide showcase: `guides/` holds guide data + snippet parsing, `guides/{slug}/{slug}.templ` are the guide pages, `templates/components/sourceview.templ` is the shared "View Source" toggle rendered in every guide. `static/css/input.css` is the Tailwind input file (compiled to `output.css`).

**Key files:**
- `guides/snippets.go` — snippet cache (`GetSnippets`), you will add `GetHighlightedSnippets` alongside it
- `guides/snippets_test.go` — existing tests use `package guides_test`
- `templates/components/sourceview.templ` — renders `{ code }` inside a `<pre>`, change to `@templ.Raw(code)`
- `guides/registry.go` — all 6 guides' `CSSVars` maps, add 8 `--code-*` vars each
- `static/css/input.css` — add chroma token → CSS var mappings here
- `guides/{slug}/{slug}.templ` — each calls `guides.GetSnippets(g.Slug)`, change to `GetHighlightedSnippets`
- `guides/cassette/styles.go` — contains `LogHandlerSnippet` const (Go source, must be highlighted separately)

**Chroma was already added:** `go get github.com/alecthomas/chroma/v2` is done. It is in `go.mod`.

**Build commands:**
- `make templ` — regenerate `*_templ.go` files after any `.templ` change
- `make tailwind` — rebuild `static/css/output.css`
- `go test ./...` — run all tests
- `git commit -m "message"` — GPG signing required; if it fails, stop and ask the user

---

## Task 1: guides/highlight.go — Highlight function and detectLang

**Files:**
- Create: `guides/highlight.go`
- Create: `guides/highlight_test.go`

**Step 1: Write the failing tests**

Create `guides/highlight_test.go`:

```go
package guides_test

import (
	"strings"
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
)

func TestHighlight_HTML(t *testing.T) {
	out := guides.Highlight(`<div class="foo">hello</div>`, "html")
	if !strings.Contains(out, "<span") {
		t.Errorf("expected span tags in output, got: %s", out)
	}
}

func TestHighlight_Go(t *testing.T) {
	out := guides.Highlight(`func main() {}`, "go")
	if !strings.Contains(out, "<span") {
		t.Errorf("expected span tags in output, got: %s", out)
	}
}

func TestHighlight_Empty(t *testing.T) {
	if out := guides.Highlight("", "go"); out != "" {
		t.Errorf("expected empty string for empty input, got: %q", out)
	}
}

func TestHighlight_NoPre(t *testing.T) {
	out := guides.Highlight(`func main() {}`, "go")
	if strings.Contains(out, "<pre") {
		t.Errorf("output must not contain a <pre> tag, got: %s", out)
	}
}

func TestDetectLang_HTML(t *testing.T) {
	cases := []string{
		`<div>hello</div>`,
		`  <span x-data="{}">`,
		`<!-- snippet:foo -->`,
	}
	for _, c := range cases {
		if lang := guides.DetectLang(c); lang != "html" {
			t.Errorf("DetectLang(%q) = %q, want html", c, lang)
		}
	}
}

func TestDetectLang_Go(t *testing.T) {
	cases := []string{
		`func main() {}`,
		`package guides`,
		`// handler comment`,
	}
	for _, c := range cases {
		if lang := guides.DetectLang(c); lang != "go" {
			t.Errorf("DetectLang(%q) = %q, want go", c, lang)
		}
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./guides/... -run TestHighlight -v
```
Expected: FAIL — `guides.Highlight` and `guides.DetectLang` undefined.

**Step 3: Implement guides/highlight.go**

```go
package guides

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Highlight returns chroma-highlighted HTML spans for the given code and language.
// Output contains bare <span> tokens with CSS classes — no surrounding <pre> tag.
// Falls back to returning the original code as plain text on any error.
func Highlight(code, lang string) string {
	if code == "" {
		return ""
	}

	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// WithClasses emits class names instead of inline styles.
	// PreventSurroundingPre removes the outer <pre> — SourceView provides its own.
	formatter := chromahtml.New(
		chromahtml.WithClasses(true),
		chromahtml.PreventSurroundingPre(true),
	)

	// Style choice does not matter when using WithClasses — colors come from CSS vars.
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return code
	}
	return buf.String()
}

// DetectLang returns "html" for template/HTML code and "go" for Go code.
// Detection is based on simple content heuristics.
func DetectLang(code string) string {
	trimmed := strings.TrimSpace(code)
	if strings.Contains(code, "<!--") || (len(trimmed) > 0 && trimmed[0] == '<') {
		return "html"
	}
	return "go"
}
```

Note: `DetectLang` is exported (capital D) so tests in `package guides_test` can call it.

**Step 4: Run tests to verify they pass**

```bash
go test ./guides/... -run "TestHighlight|TestDetectLang" -v
```
Expected: all 6 tests PASS.

**Step 5: Commit**

```bash
git add guides/highlight.go guides/highlight_test.go go.mod go.sum
git commit -m "feat: add chroma-based Highlight and DetectLang to guides package"
```

---

## Task 2: GetHighlightedSnippets — highlighted snippet cache

**Files:**
- Modify: `guides/snippets.go`
- Modify: `guides/snippets_test.go`

**Step 1: Write the failing test**

Add to `guides/snippets_test.go` (inside the existing file, after the last test):

```go
func TestGetHighlightedSnippets_ReturnsNonNilForUnknownSlug(t *testing.T) {
	got := guides.GetHighlightedSnippets("does-not-exist")
	if got == nil {
		t.Error("GetHighlightedSnippets must return non-nil map")
	}
}

func TestGetHighlightedSnippets_HighlightedHTMLContainsSpans(t *testing.T) {
	// Verify that Highlight applied to a real HTML snippet produces span tags.
	raw := guides.ParseSnippets(`<!-- snippet:demo --><div x-data="{}">hello</div><!-- /snippet:demo -->`)
	code, ok := raw["demo"]
	if !ok {
		t.Fatal("test fixture snippet not parsed")
	}
	hl := guides.Highlight(code, guides.DetectLang(code))
	if !strings.Contains(hl, "<span") {
		t.Errorf("highlighted HTML snippet expected to contain <span, got: %s", hl)
	}
}
```

Also add `"strings"` to the imports in `guides/snippets_test.go` if not already present.

**Step 2: Run tests to verify they fail**

```bash
go test ./guides/... -run "TestGetHighlightedSnippets" -v
```
Expected: FAIL — `guides.GetHighlightedSnippets` undefined.

**Step 3: Add GetHighlightedSnippets to guides/snippets.go**

Add after the existing `GetSnippets` function (at the bottom of the file):

```go
var (
	highlightedCache map[string]map[string]string
	highlightOnce   sync.Once
)

func loadHighlighted() map[string]map[string]string {
	raw := loadAll()
	out := make(map[string]map[string]string, len(raw))
	for slug, snippets := range raw {
		out[slug] = make(map[string]string, len(snippets))
		for key, code := range snippets {
			out[slug][key] = Highlight(code, DetectLang(code))
		}
	}
	return out
}

// GetHighlightedSnippets returns syntax-highlighted HTML for each snippet of the named guide.
// Highlighting is computed once at startup and cached. Returns a non-nil map even if the slug is unknown.
func GetHighlightedSnippets(slug string) map[string]string {
	highlightOnce.Do(func() {
		highlightedCache = loadHighlighted()
	})
	if s, ok := highlightedCache[slug]; ok {
		return s
	}
	return map[string]string{}
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./guides/... -v
```
Expected: all tests PASS (including existing ParseSnippets tests).

**Step 5: Commit**

```bash
git add guides/snippets.go guides/snippets_test.go
git commit -m "feat: add GetHighlightedSnippets with startup-cached chroma highlighting"
```

---

## Task 3: Add --code-* CSS vars to all 6 guides in registry.go

**Files:**
- Modify: `guides/registry.go`

**Step 1: No test to write** — CSS vars flow through to the browser; this is a data change verified visually and by the build passing.

**Step 2: Edit guides/registry.go**

Add the following 8 vars to each guide's `CSSVars` map. Find each guide by its `Slug` field.

**Brutalist** (black terminal, red keywords — raw and harsh):
```go
"--code-bg":      "#0a0a0a",
"--code-text":    "#f0f0f0",
"--code-keyword": "#ff0000",
"--code-string":  "#ffff00",
"--code-comment": "#555555",
"--code-number":  "#ff8800",
"--code-tag":     "#ff0000",
"--code-attr":    "#cccccc",
```

**Cassette** (warm dark, NASA blue keywords — technical document feel):
```go
"--code-bg":      "#1a1814",
"--code-text":    "#c8c7c0",
"--code-keyword": "#4a9eff",
"--code-string":  "#98c379",
"--code-comment": "#5a6a52",
"--code-number":  "#d19a66",
"--code-tag":     "#e06c75",
"--code-attr":    "#abb2bf",
```

**Minimal** (GitHub light — clean and professional):
```go
"--code-bg":      "#f6f8fa",
"--code-text":    "#24292f",
"--code-keyword": "#0550ae",
"--code-string":  "#0a3069",
"--code-comment": "#6e7781",
"--code-number":  "#953800",
"--code-tag":     "#116329",
"--code-attr":    "#953800",
```

**Glass** (deep dark purple/blue — matches the dark frosted aesthetic):
```go
"--code-bg":      "rgba(10,10,25,0.9)",
"--code-text":    "#f1f5f9",
"--code-keyword": "#a78bfa",
"--code-string":  "#86efac",
"--code-comment": "#475569",
"--code-number":  "#f472b6",
"--code-tag":     "#60a5fa",
"--code-attr":    "#94a3b8",
```

**Bento** (One Dark-ish — modern SaaS dark mode):
```go
"--code-bg":      "#1e1e2e",
"--code-text":    "#e2e8f0",
"--code-keyword": "#6366f1",
"--code-string":  "#86efac",
"--code-comment": "#64748b",
"--code-number":  "#f59e0b",
"--code-tag":     "#8b5cf6",
"--code-attr":    "#94a3b8",
```

**Swiss** (stark black, red as signal — consistent with the guide's red/black/white philosophy):
```go
"--code-bg":      "#1a1a1a",
"--code-text":    "#ffffff",
"--code-keyword": "#e63329",
"--code-string":  "#767676",
"--code-comment": "#444444",
"--code-number":  "#e63329",
"--code-tag":     "#ffffff",
"--code-attr":    "#aaaaaa",
```

**Step 3: Verify build**

```bash
go build ./...
```
Expected: SUCCESS (registry.go is pure data, no compilation issues).

**Step 4: Commit**

```bash
git add guides/registry.go
git commit -m "feat: add --code-* CSS vars to all 6 guides for syntax highlighting"
```

---

## Task 4: CSS token mappings + SourceView component update

**Files:**
- Modify: `static/css/input.css`
- Modify: `templates/components/sourceview.templ`

**Step 1: Add chroma token → CSS var mappings to static/css/input.css**

Append to the end of `static/css/input.css` (after the `[x-cloak]` rule):

```css
/* Syntax highlighting — chroma CSS class tokens mapped to per-guide CSS vars */
.code-block {
    background: var(--code-bg, #1a1a1a);
    color: var(--code-text, #f0f0f0);
    border-radius: var(--radius-sm, 0);
}
/* Keywords: func, package, if, for, return, var, etc. */
.code-block .k, .code-block .kd, .code-block .kn,
.code-block .kp, .code-block .kr { color: var(--code-keyword); font-weight: bold; }
/* String literals */
.code-block .s,  .code-block .s1, .code-block .s2,
.code-block .sa, .code-block .sb, .code-block .sc { color: var(--code-string); }
/* Comments */
.code-block .c,  .code-block .c1, .code-block .cm,
.code-block .cs, .code-block .cp { color: var(--code-comment); font-style: italic; }
/* Numeric literals */
.code-block .m,  .code-block .mi, .code-block .mf,
.code-block .mo { color: var(--code-number); }
/* HTML/Templ tag names */
.code-block .nt { color: var(--code-tag); }
/* HTML/Templ attribute names */
.code-block .na { color: var(--code-attr); }
/* Builtins and function names — use keyword color */
.code-block .nb, .code-block .nf, .code-block .nx { color: var(--code-keyword); }
```

**Step 2: Update templates/components/sourceview.templ**

Current `<pre>` block (lines 18-21):
```go
<pre
    class="mt-2 p-4 overflow-x-auto text-xs rounded border"
    style="font-family: var(--font-body, monospace); background: var(--color-surface, #f8f8f8); border-color: var(--color-border, #e5e7eb); color: var(--color-text, #1a1a1a); line-height: 1.6; white-space: pre; tab-size: 2;"
>{ code }</pre>
```

Replace with:
```go
<pre
    class="code-block mt-2 p-4 overflow-x-auto text-xs rounded border"
    style="font-family: monospace; border-color: var(--color-border, #e5e7eb); line-height: 1.6; white-space: pre; tab-size: 2;"
>@templ.Raw(code)</pre>
```

Key changes:
- Added `code-block` class (picks up the chroma CSS var mappings)
- Removed `background`, `color` from inline style (now driven by `.code-block` in CSS)
- Changed `{ code }` to `@templ.Raw(code)` (code is now pre-highlighted HTML, not plain text)

**Step 3: Regenerate templ and rebuild Tailwind**

```bash
make templ
make tailwind
```
Expected: both complete without errors.

**Step 4: Run all tests**

```bash
go test ./...
```
Expected: all tests PASS (existing handler tests check HTTP 200, still valid).

**Step 5: Commit**

```bash
git add static/css/input.css templates/components/sourceview.templ templates/components/sourceview_templ.go
git commit -m "feat: add chroma CSS token mappings and update SourceView to render highlighted HTML"
```

---

## Task 5: Wire all 6 guide templ files + fix Cassette system-log

**Files:**
- Modify: `guides/brutalist/brutalist.templ`
- Modify: `guides/cassette/cassette.templ`
- Modify: `guides/minimal/minimal.templ`
- Modify: `guides/glass/glass.templ`
- Modify: `guides/bento/bento.templ`
- Modify: `guides/swiss/swiss.templ`

**Context:** Every guide templ file has this near the top of `Page()`:
```go
{{ snippets := guides.GetSnippets(g.Slug) }}
```
Change every occurrence to:
```go
{{ snippets := guides.GetHighlightedSnippets(g.Slug) }}
```
That is the only change needed in 5 of the 6 guide files. Cassette needs an additional fix described below.

**Step 1: Update the 5 non-Cassette guides**

In each of these files, replace `guides.GetSnippets(g.Slug)` with `guides.GetHighlightedSnippets(g.Slug)`:
- `guides/brutalist/brutalist.templ`
- `guides/minimal/minimal.templ`
- `guides/glass/glass.templ`
- `guides/bento/bento.templ`
- `guides/swiss/swiss.templ`

**Step 2: Fix guides/cassette/cassette.templ**

First, the same `GetSnippets` → `GetHighlightedSnippets` change.

Second, find the system-log SourceView. It currently looks like:
```go
@components.SourceView("<!-- HTMX (client) -->\n" + snippets["system-log"] + "\n\n<!-- Go handler (server) -->\n" + LogHandlerSnippet)
```

Replace it with two separate labeled SourceView calls:
```go
<p class="text-xs mb-1" style="color: var(--color-text-muted); font-family: var(--font-body);">HTMX client — cassette.templ</p>
@components.SourceView(snippets["system-log"])
<p class="text-xs mt-4 mb-1" style="color: var(--color-text-muted); font-family: var(--font-body);">Go handler — handlers/guides.go</p>
@components.SourceView(guides.Highlight(LogHandlerSnippet, "go"))
```

Note: `LogHandlerSnippet` is defined in `guides/cassette/styles.go` in the same `cassette` package, so no import needed. `guides.Highlight` is already accessible since `cassette.templ` imports `github.com/johnfarrell/stylesheets/guides`.

**Step 3: Regenerate templ**

```bash
make templ
```
Expected: generates all `*_templ.go` files without errors.

**Step 4: Run all tests**

```bash
go test ./...
```
Expected: all tests PASS.

**Step 5: Build the full binary**

```bash
make build
```
Expected: clean build with no errors.

**Step 6: Commit**

```bash
git add guides/brutalist/brutalist_templ.go guides/brutalist/brutalist.templ \
        guides/cassette/cassette_templ.go guides/cassette/cassette.templ \
        guides/minimal/minimal_templ.go guides/minimal/minimal.templ \
        guides/glass/glass_templ.go guides/glass/glass.templ \
        guides/bento/bento_templ.go guides/bento/bento.templ \
        guides/swiss/swiss_templ.go guides/swiss/swiss.templ
git commit -m "feat: wire GetHighlightedSnippets into all guides, split Cassette system-log view"
```
