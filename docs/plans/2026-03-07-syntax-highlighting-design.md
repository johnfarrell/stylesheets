# Syntax Highlighting Design

Date: 2026-03-07

## Goal

Add token-accurate syntax highlighting to every SourceView code block, with each guide using its own color scheme derived from its design language.

## Architecture

Server-side only. Zero client-side JS. Zero per-request overhead.

`github.com/alecthomas/chroma/v2` parses code into tokens and emits HTML `<span>` elements with CSS class names. Guide CSS vars map those class names to per-guide colors. Highlighting is cached at startup alongside the existing snippet cache.

**New file: `guides/highlight.go`**

```go
func Highlight(code, lang string) string
func detectLang(code string) string   // "html" if contains "<!--" or starts with "<", else "go"
```

Chroma formatter options: `html.WithClasses(true)`, `html.PreventSurroundingPre(true)` — outputs bare `<span>` tokens only, no wrapping `<pre>`.

**Modified: `guides/snippets.go`**

Add `highlightedCache map[string]map[string]string` and `highlightOnce sync.Once`. New `GetHighlightedSnippets(slug string) map[string]string` populates the cache by calling `Highlight(raw, detectLang(raw))` for each snippet in `GetSnippets`.

**Modified: `templates/components/sourceview.templ`**

`<pre>` gains `class="code-block"`. Inner content changes from `{ code }` to `@templ.Raw(code)`.

**Modified: All six guide `.templ` files**

`guides.GetSnippets(g.Slug)` → `guides.GetHighlightedSnippets(g.Slug)`.

**Modified: `guides/cassette/cassette.templ`**

The combined system-log SourceView (HTML client + Go handler in one string) is split into two separate labeled SourceView calls:
- HTML snippet: `snippets["system-log"]` (pre-highlighted as HTML)
- Go snippet: `guides.Highlight(cassette.LogHandlerSnippet, "go")` (highlighted inline)

---

## CSS Token Mapping

Added to `static/css/input.css`:

```css
.code-block { background: var(--code-bg); color: var(--code-text); border-radius: var(--radius-sm); }
.code-block .k, .code-block .kd, .code-block .kn, .code-block .kp { color: var(--code-keyword); font-weight: bold; }
.code-block .s, .code-block .s1, .code-block .s2, .code-block .sa { color: var(--code-string); }
.code-block .c, .code-block .c1, .code-block .cm                   { color: var(--code-comment); font-style: italic; }
.code-block .m, .code-block .mi, .code-block .mf                   { color: var(--code-number); }
.code-block .nt                                                      { color: var(--code-tag); }
.code-block .na                                                      { color: var(--code-attr); }
```

---

## Per-Guide CSS Vars

Eight new vars added to each guide's `CSSVars` in `guides/registry.go`:

`--code-bg`, `--code-text`, `--code-keyword`, `--code-string`, `--code-comment`, `--code-number`, `--code-tag`, `--code-attr`

| Guide     | bg          | text      | keyword   | string    | comment   | number    | tag       | attr      |
|-----------|-------------|-----------|-----------|-----------|-----------|-----------|-----------|-----------|
| Brutalist | `#0a0a0a`   | `#f0f0f0` | `#ff0000` | `#ffff00` | `#555555` | `#ff8800` | `#ff0000` | `#cccccc` |
| Minimal   | `#f6f8fa`   | `#24292f` | `#0550ae` | `#0a3069` | `#6e7781` | `#953800` | `#116329` | `#953800` |
| Cassette  | `#1a1814`   | `#c8c7c0` | `#4a9eff` | `#98c379` | `#5a6a52` | `#d19a66` | `#e06c75` | `#abb2bf` |
| Glass     | `rgba(10,10,25,0.9)` | `#f1f5f9` | `#a78bfa` | `#86efac` | `#475569` | `#f472b6` | `#60a5fa` | `#94a3b8` |
| Bento     | `#1e1e2e`   | `#e2e8f0` | `#6366f1` | `#86efac` | `#64748b` | `#f59e0b` | `#8b5cf6` | `#94a3b8` |
| Swiss     | `#1a1a1a`   | `#ffffff` | `#e63329` | `#767676` | `#444444` | `#e63329` | `#ffffff` | `#aaaaaa` |

---

## Language Detection

Simple heuristic in `detectLang`:
- Contains `<!--` or first non-whitespace char is `<` → `"html"`
- Otherwise → `"go"`

All `.templ` snippet file content detects as HTML. `LogHandlerSnippet` and any pure Go snippets detect as Go.
