# New Guides (Terminal, Retro OS, Newspaper) Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Add three new style guides — Terminal, Retro OS, and Newspaper — that showcase advanced HTMX (SSE, infinite scroll, view transitions), Alpine.js (keyboard nav, drag-and-drop, Alpine.store), and CSS (CRT effects, beveled borders, newspaper columns).

**Architecture:** Each guide follows the established pattern: `guides/{slug}/styles.go` + `guides/{slug}/{slug}.templ` + handler routes in `handlers/guides.go`. Guide-specific HTMX endpoint responses use small templ components. The SSE extension is loaded globally in `layout.templ`.

**Tech Stack:** Go 1.25, Templ v0.3.1001, HTMX 2.0.4 + SSE extension 2.2.2, Alpine.js 3.14.9, Tailwind CSS v4

**Design doc:** `docs/plans/2026-03-07-new-guides-design.md`

---

### Task 0: Add HTMX SSE extension to layout and add FormResponse cases for new guides

Before building guides, add the SSE extension globally and pre-add the new guide slugs to the shared `FormResponse` component so the demo-form endpoint works for them.

**Files:**
- Modify: `templates/layout.templ:14-16`
- Modify: `templates/components/formresponse.templ`

**Step 1: Add SSE extension script to `templates/layout.templ`**

After the HTMX script tag (line 14), add:
```html
<script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
```

**Step 2: Add new cases to `templates/components/formresponse.templ`**

Add three new `case` blocks before the `default:` case:

```
case "terminal":
    <div style="background: #0a0a0a; border: 1px solid #00ff41; padding: 0.75rem; font-family: 'Fira Code', monospace; font-size: 0.8125rem;">
        <div style="color: #00ff41;">&#9654; RECEIVED</div>
        <div style="color: #00ff41; opacity: 0.7;">operator: <strong style="color: #00ff41;">{ name }</strong> — logged</div>
    </div>
case "retro":
    <div style="background: #c0c0c0; border: 2px outset #c0c0c0; padding: 0.75rem; font-family: 'IBM Plex Sans', sans-serif;">
        <div style="background: #000080; color: #fff; padding: 0.25rem 0.5rem; font-weight: 700; font-size: 0.75rem; margin: -0.75rem -0.75rem 0.5rem; border-bottom: 2px outset #c0c0c0;">Message</div>
        <p style="font-size: 0.875rem;">Form submitted. Thank you, <strong>{ name }</strong>.</p>
    </div>
case "newspaper":
    <div style="border-top: 2px solid #1a1a1a; padding: 0.75rem 0; margin-top: 0.5rem;">
        <p style="font-family: 'Source Serif 4', serif; font-variant: small-caps; font-size: 0.75rem; color: #c41e1e; margin-bottom: 0.25rem;">Received</p>
        <p style="font-family: 'Playfair Display', serif; font-weight: 700; color: #1a1a1a;">{ name }</p>
    </div>
```

**Step 3: Regenerate templ and verify**

Run: `make templ && go build ./...`

**Step 4: Commit**

```bash
git add templates/layout.templ templates/layout_templ.go templates/components/formresponse.templ templates/components/formresponse_templ.go
git commit -m "feat: add HTMX SSE extension and FormResponse cases for new guides"
```

---

### Task 1: Register Terminal guide (registry + embed + handler wiring)

Wire up the Terminal guide's metadata, embed directive, and handler switch case so the infrastructure is in place before building the page.

**Files:**
- Modify: `guides/registry.go` (add Terminal to `All` slice)
- Modify: `guides/sources.go` (add `terminal` to embed directive)
- Create: `guides/terminal/styles.go`
- Create: `guides/terminal/terminal.templ` (minimal placeholder)
- Modify: `handlers/guides.go` (add import + switch case + guide-specific routes)

**Step 1: Add Terminal entry to `guides/registry.go` `All` slice**

Add after the Swiss entry:

```go
{
    Name:        "Terminal",
    Slug:        "terminal",
    Description: "Green-on-black CRT aesthetic. Scanlines, glow, monospace everything.",
    FontURL:     "https://fonts.googleapis.com/css2?family=Fira+Code:wght@300;400;500;600;700&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#0a0a0a",
        "--color-surface":     "#111111",
        "--color-surface-2":   "#1a1a1a",
        "--color-primary":     "#00ff41",
        "--color-secondary":   "#00bfff",
        "--color-accent":      "#ffcc00",
        "--color-danger":      "#ff3333",
        "--color-text":        "#00ff41",
        "--color-text-muted":  "#00994d",
        "--color-border":      "#00ff41",
        "--font-display":      "'Fira Code', monospace",
        "--font-body":         "'Fira Code', monospace",
        "--font-mono":         "'Fira Code', monospace",
        "--font-size-display": "2.5rem",
        "--font-size-heading": "1.25rem",
        "--font-size-body":    "0.875rem",
        "--font-size-caption": "0.75rem",
        "--radius-sm":         "0px",
        "--radius-md":         "0px",
        "--radius-lg":         "2px",
        "--shadow-card":       "0 0 10px rgba(0,255,65,0.15)",
        "--shadow-btn":        "0 0 8px rgba(0,255,65,0.2)",
        "--border-width":      "1px",
        "--border-color":      "#00ff41",
        "--content-max-width": "960px",
        "--section-padding":   "3rem 2rem",
        "--code-bg":           "#050505",
        "--code-text":         "#00ff41",
        "--code-keyword":      "#00bfff",
        "--code-string":       "#ffcc00",
        "--code-comment":      "#00994d",
        "--code-number":       "#ff3333",
        "--code-tag":          "#00bfff",
        "--code-attr":         "#00994d",
    },
},
```

**Step 2: Update embed directive in `guides/sources.go`**

Change:
```go
//go:embed brutalist minimal cassette glass bento swiss
```
To:
```go
//go:embed brutalist minimal cassette glass bento swiss terminal
```

**Step 3: Create `guides/terminal/styles.go`**

```go
package terminal

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - CRT scanline overlay via pseudo-element */
.term-screen { position: relative; }
.term-screen::after {
    content: "";
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: repeating-linear-gradient(
        0deg,
        transparent,
        transparent 2px,
        rgba(0,0,0,0.15) 2px,
        rgba(0,0,0,0.15) 4px
    );
    z-index: 1;
}
/* [custom] - CRT text glow not achievable with Tailwind */
.term-glow { text-shadow: 0 0 5px currentColor; }
.term-glow-strong { text-shadow: 0 0 8px currentColor, 0 0 15px currentColor; }
/* [custom] - blinking cursor animation */
@keyframes term-blink { 0%,49% { opacity: 1; } 50%,100% { opacity: 0; } }
.term-cursor { animation: term-blink 1s step-end infinite; }
/* [custom] - terminal panel */
.term-panel {
    background: var(--color-surface);
    border: var(--border-width) solid var(--color-border);
    box-shadow: var(--shadow-card);
}
/* [custom] - terminal button */
.term-btn {
    background: transparent;
    color: var(--color-primary);
    border: var(--border-width) solid var(--color-primary);
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 500;
    cursor: pointer;
    transition: background 0.1s, color 0.1s, box-shadow 0.1s;
    padding: 0.4rem 1rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
.term-btn:hover {
    background: var(--color-primary);
    color: var(--color-bg);
    box-shadow: var(--shadow-btn);
}
.term-btn-danger { border-color: var(--color-danger); color: var(--color-danger); }
.term-btn-danger:hover { background: var(--color-danger); color: var(--color-bg); }
/* [custom] - terminal input with glow focus */
.term-input {
    background: var(--color-bg);
    border: var(--border-width) solid var(--color-border);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    padding: 0.4rem 0.5rem;
    width: 100%;
    caret-color: var(--color-primary);
}
.term-input:focus {
    outline: none;
    box-shadow: 0 0 6px rgba(0,255,65,0.3);
}
/* [custom] - HTMX loading indicator */
.term-indicator { display: none; }
.htmx-request .term-indicator,
.htmx-request.term-indicator { display: inline; }
/* [custom] - file browser item highlight */
.term-file-active {
    background: var(--color-primary);
    color: var(--color-bg);
}
`
}
```

**Step 4: Create minimal `guides/terminal/terminal.templ`**

```
package terminal

import (
	"github.com/johnfarrell/stylesheets/guides"
)

// Page renders the Terminal style guide showcase.
templ Page(g guides.Guide, htmxRequest bool) {
	if htmxRequest {
		<span id="font-loader" { templ.Attributes{"hx-swap-oob": "true"}... }>
			<link rel="preconnect" href="https://fonts.googleapis.com"/>
			<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
			<link rel="stylesheet" href={ g.FontURL }/>
		</span>
	}
	@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
	<div class="term-screen p-8 min-h-full" style="background: var(--color-bg); color: var(--color-text); font-family: var(--font-body);">
		<div style="max-width: var(--content-max-width); margin: 0 auto;">
			<h1 class="text-3xl font-bold term-glow mb-2" style="font-family: var(--font-display);">Terminal</h1>
			<p style="color: var(--color-text-muted); font-size: var(--font-size-caption);">Guide under construction...</p>
		</div>
	</div>
}
```

**Step 5: Add import and switch case to `handlers/guides.go`**

Add import:
```go
terminaltempl "github.com/johnfarrell/stylesheets/guides/terminal"
```

Add switch case in `guideContent()`:
```go
case "terminal":
    return terminaltempl.Page(g, htmxRequest)
```

**Step 6: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 7: Commit**

```bash
git add guides/registry.go guides/sources.go guides/terminal/ handlers/guides.go
git commit -m "feat: register Terminal guide with styles and placeholder page"
```

---

### Task 2: Register Retro OS guide (registry + embed + handler wiring)

Same infrastructure setup for the Retro OS guide.

**Files:**
- Modify: `guides/registry.go` (add Retro OS to `All` slice)
- Modify: `guides/sources.go` (add `retro` to embed directive)
- Create: `guides/retro/styles.go`
- Create: `guides/retro/retro.templ` (minimal placeholder)
- Modify: `handlers/guides.go` (add import + switch case)

**Step 1: Add Retro OS entry to `guides/registry.go` `All` slice**

```go
{
    Name:        "Retro OS",
    Slug:        "retro",
    Description: "Windowed desktop UI — beveled borders, draggable windows, taskbar.",
    FontURL:     "https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@300;400;500;700&family=VT323&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#008080",
        "--color-surface":     "#c0c0c0",
        "--color-surface-2":   "#dfdfdf",
        "--color-primary":     "#000080",
        "--color-secondary":   "#c0c0c0",
        "--color-accent":      "#ffff00",
        "--color-danger":      "#ff0000",
        "--color-text":        "#000000",
        "--color-text-muted":  "#808080",
        "--color-border":      "#808080",
        "--color-highlight":   "#000080",
        "--font-display":      "'VT323', monospace",
        "--font-body":         "'IBM Plex Sans', sans-serif",
        "--font-mono":         "'VT323', monospace",
        "--font-size-display": "2.5rem",
        "--font-size-heading": "1.25rem",
        "--font-size-body":    "0.8125rem",
        "--font-size-caption": "0.6875rem",
        "--radius-sm":         "0px",
        "--radius-md":         "0px",
        "--radius-lg":         "0px",
        "--shadow-card":       "2px 2px 0px #000000",
        "--shadow-btn":        "1px 1px 0px #000000",
        "--border-width":      "2px",
        "--border-color":      "#808080",
        "--content-max-width": "1100px",
        "--section-padding":   "2rem 2rem",
        "--code-bg":           "#000000",
        "--code-text":         "#c0c0c0",
        "--code-keyword":      "#5555ff",
        "--code-string":       "#ffff55",
        "--code-comment":      "#55ff55",
        "--code-number":       "#ff5555",
        "--code-tag":          "#55ffff",
        "--code-attr":         "#808080",
    },
},
```

**Step 2: Update embed directive in `guides/sources.go`**

```go
//go:embed brutalist minimal cassette glass bento swiss terminal retro
```

**Step 3: Create `guides/retro/styles.go`**

```go
package retro

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - Win95-style 3D beveled borders not achievable with Tailwind */
.retro-raised {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
}
.retro-inset {
    border: 2px solid;
    border-color: #808080 #ffffff #ffffff #808080;
    background: #fff;
}
/* [custom] - window chrome with title bar */
.retro-window {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    box-shadow: var(--shadow-card);
}
.retro-titlebar {
    background: linear-gradient(90deg, var(--color-primary), #1084d0);
    color: #ffffff;
    padding: 0.25rem 0.5rem;
    font-weight: 700;
    font-size: var(--font-size-caption);
    display: flex;
    align-items: center;
    justify-content: space-between;
    user-select: none;
    cursor: default;
}
.retro-titlebar-inactive {
    background: linear-gradient(90deg, #808080, #a0a0a0);
}
/* [custom] - window control buttons (close/minimize) */
.retro-winbtn {
    width: 16px;
    height: 14px;
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-size: 8px;
    line-height: 10px;
    text-align: center;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    font-family: var(--font-body);
}
.retro-winbtn:active {
    border-color: #808080 #ffffff #ffffff #808080;
}
/* [custom] - raised button with press state */
.retro-btn {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    padding: 0.25rem 1rem;
    cursor: pointer;
}
.retro-btn:active {
    border-color: #808080 #ffffff #ffffff #808080;
    padding: 0.3rem 0.95rem 0.2rem 1.05rem;
}
.retro-btn-primary {
    background: var(--color-surface);
    outline: 1px dotted #000;
    outline-offset: -4px;
}
/* [custom] - inset input field */
.retro-input {
    border: 2px solid;
    border-color: #808080 #ffffff #ffffff #808080;
    background: #ffffff;
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    padding: 0.2rem 0.4rem;
    color: #000;
    width: 100%;
}
.retro-input:focus {
    outline: none;
}
/* [custom] - desktop icon */
.retro-icon {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    cursor: pointer;
    padding: 0.5rem;
    border: 1px solid transparent;
    font-size: var(--font-size-caption);
    color: #ffffff;
    text-shadow: 1px 1px 1px #000;
}
.retro-icon:hover {
    border: 1px dotted #ffffff;
}
.retro-icon-img {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
}
/* [custom] - taskbar at bottom */
.retro-taskbar {
    background: var(--color-surface);
    border-top: 2px solid;
    border-color: #ffffff;
    padding: 0.25rem;
    display: flex;
    gap: 0.25rem;
    align-items: center;
}
.retro-taskbar-btn {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    padding: 0.15rem 0.5rem;
    cursor: pointer;
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.retro-taskbar-btn-active {
    border-color: #808080 #ffffff #ffffff #808080;
    background: #dfdfdf;
}
/* [custom] - start button */
.retro-start-btn {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-family: var(--font-body);
    font-weight: 700;
    font-size: var(--font-size-caption);
    padding: 0.15rem 0.5rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.25rem;
}
.retro-start-btn:active {
    border-color: #808080 #ffffff #ffffff #808080;
}
/* [custom] - checkbox and radio with retro styling */
.retro-check {
    width: 13px;
    height: 13px;
    accent-color: var(--color-primary);
}
@media (max-width: 768px) {
    .retro-window { position: static !important; transform: none !important; margin-bottom: 1rem; }
}
`
}
```

**Step 4: Create minimal `guides/retro/retro.templ`**

```
package retro

import (
	"github.com/johnfarrell/stylesheets/guides"
)

// Page renders the Retro OS style guide showcase.
templ Page(g guides.Guide, htmxRequest bool) {
	if htmxRequest {
		<span id="font-loader" { templ.Attributes{"hx-swap-oob": "true"}... }>
			<link rel="preconnect" href="https://fonts.googleapis.com"/>
			<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
			<link rel="stylesheet" href={ g.FontURL }/>
		</span>
	}
	@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
	<div style="background: var(--color-bg); min-height: 100%; padding: 2rem;">
		<div style="max-width: var(--content-max-width); margin: 0 auto;">
			<div class="retro-window" style="width: 400px;">
				<div class="retro-titlebar">
					<span>Welcome</span>
				</div>
				<div style="padding: 1rem;">
					<h1 style="font-family: var(--font-display); font-size: var(--font-size-display);">Retro OS</h1>
					<p style="font-size: var(--font-size-caption); color: var(--color-text-muted);">Guide under construction...</p>
				</div>
			</div>
		</div>
	</div>
}
```

**Step 5: Update `handlers/guides.go`**

Add import:
```go
retrotempl "github.com/johnfarrell/stylesheets/guides/retro"
```

Add switch case:
```go
case "retro":
    return retrotempl.Page(g, htmxRequest)
```

**Step 6: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 7: Commit**

```bash
git add guides/registry.go guides/sources.go guides/retro/ handlers/guides.go
git commit -m "feat: register Retro OS guide with styles and placeholder page"
```

---

### Task 3: Register Newspaper guide (registry + embed + handler wiring)

Same infrastructure setup for the Newspaper guide.

**Files:**
- Modify: `guides/registry.go` (add Newspaper to `All` slice)
- Modify: `guides/sources.go` (add `newspaper` to embed directive)
- Create: `guides/newspaper/styles.go`
- Create: `guides/newspaper/newspaper.templ` (minimal placeholder)
- Modify: `handlers/guides.go` (add import + switch case)

**Step 1: Add Newspaper entry to `guides/registry.go` `All` slice**

```go
{
    Name:        "The Daily Style",
    Slug:        "newspaper",
    Description: "Broadsheet editorial layout — serif headlines, multi-column text, drop caps.",
    FontURL:     "https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,700;0,900;1,400&family=Source+Serif+4:ital,opsz,wght@0,8..60,300;0,8..60,400;0,8..60,600;0,8..60,700;1,8..60,400&display=swap",
    CSSVars: map[string]string{
        "--color-bg":          "#faf9f6",
        "--color-surface":     "#ffffff",
        "--color-surface-2":   "#f0efec",
        "--color-primary":     "#c41e1e",
        "--color-secondary":   "#1a1a1a",
        "--color-text":        "#1a1a1a",
        "--color-text-muted":  "#6b6b6b",
        "--color-border":      "#1a1a1a",
        "--font-display":      "'Playfair Display', serif",
        "--font-body":         "'Source Serif 4', serif",
        "--font-mono":         "'Courier New', Courier, monospace",
        "--font-size-display": "4rem",
        "--font-size-heading": "2rem",
        "--font-size-body":    "1.0625rem",
        "--font-size-caption": "0.75rem",
        "--radius-sm":         "0px",
        "--radius-md":         "0px",
        "--radius-lg":         "0px",
        "--shadow-card":       "none",
        "--shadow-btn":        "none",
        "--border-width":      "1px",
        "--border-color":      "#1a1a1a",
        "--content-max-width": "1000px",
        "--section-padding":   "4rem 2rem",
        "--code-bg":           "#1a1a1a",
        "--code-text":         "#faf9f6",
        "--code-keyword":      "#c41e1e",
        "--code-string":       "#6b6b6b",
        "--code-comment":      "#808080",
        "--code-number":       "#c41e1e",
        "--code-tag":          "#1a1a1a",
        "--code-attr":         "#6b6b6b",
    },
},
```

**Step 2: Update embed directive in `guides/sources.go`**

```go
//go:embed brutalist minimal cassette glass bento swiss terminal retro newspaper
```

**Step 3: Create `guides/newspaper/styles.go`**

```go
package newspaper

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - multi-column text flow not achievable with Tailwind */
.news-columns-2 { column-count: 2; column-gap: 2rem; column-rule: 1px solid var(--color-border); }
.news-columns-3 { column-count: 3; column-gap: 2rem; column-rule: 1px solid var(--color-border); }
@media (max-width: 768px) {
    .news-columns-2, .news-columns-3 { column-count: 1; }
}
/* [custom] - drop cap not achievable with Tailwind */
.news-dropcap::first-letter {
    float: left;
    font-family: var(--font-display);
    font-size: 4rem;
    line-height: 0.8;
    padding-right: 0.5rem;
    padding-top: 0.25rem;
    color: var(--color-primary);
    font-weight: 900;
}
/* [custom] - section rule line */
.news-rule { border-top: 1px solid var(--color-border); }
.news-rule-thick { border-top: 3px solid var(--color-border); }
.news-rule-red { border-top: 2px solid var(--color-primary); }
/* [custom] - byline with small-caps */
.news-byline {
    font-family: var(--font-body);
    font-variant: small-caps;
    font-size: var(--font-size-caption);
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
}
/* [custom] - pull quote */
.news-pullquote {
    border-left: 3px solid var(--color-primary);
    padding-left: 1.5rem;
    margin: 2rem 0;
    font-family: var(--font-display);
    font-style: italic;
    font-size: 1.5rem;
    line-height: 1.4;
    color: var(--color-secondary);
}
/* [custom] - masthead */
.news-masthead {
    text-align: center;
    border-top: 3px double var(--color-border);
    border-bottom: 3px double var(--color-border);
    padding: 1rem 0;
}
/* [custom] - article card */
.news-card {
    border-top: 2px solid var(--color-border);
    padding-top: 1rem;
}
/* [custom] - button styled as editorial link */
.news-btn {
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 600;
    color: var(--color-secondary);
    background: none;
    border: 1px solid var(--color-border);
    padding: 0.4rem 1rem;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
}
.news-btn:hover {
    background: var(--color-secondary);
    color: var(--color-bg);
}
.news-btn-primary {
    background: var(--color-primary);
    color: #fff;
    border-color: var(--color-primary);
}
.news-btn-primary:hover {
    background: #a01818;
    border-color: #a01818;
}
/* [custom] - input styled as editorial underline field */
.news-input {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--color-border);
    color: var(--color-text);
    padding: 0.375rem 0;
    width: 100%;
}
.news-input:focus {
    outline: none;
    border-bottom: 2px solid var(--color-primary);
}
/* [custom] - breaking news banner */
.news-breaking {
    background: var(--color-primary);
    color: #fff;
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    padding: 0.5rem 1rem;
}
/* [custom] - reading progress bar */
.news-progress {
    position: fixed;
    top: 0;
    left: 0;
    height: 3px;
    background: var(--color-primary);
    z-index: 100;
    transition: width 0.1s linear;
}
/* [custom] - headline sizes */
.news-headline-lg {
    font-family: var(--font-display);
    font-weight: 900;
    font-size: 2.5rem;
    line-height: 1.1;
    color: var(--color-secondary);
}
.news-headline-md {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.5rem;
    line-height: 1.2;
    color: var(--color-secondary);
}
.news-headline-sm {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    line-height: 1.3;
    color: var(--color-secondary);
}
`
}
```

**Step 4: Create minimal `guides/newspaper/newspaper.templ`**

```
package newspaper

import (
	"github.com/johnfarrell/stylesheets/guides"
)

// Page renders the Newspaper style guide showcase.
templ Page(g guides.Guide, htmxRequest bool) {
	if htmxRequest {
		<span id="font-loader" { templ.Attributes{"hx-swap-oob": "true"}... }>
			<link rel="preconnect" href="https://fonts.googleapis.com"/>
			<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
			<link rel="stylesheet" href={ g.FontURL }/>
		</span>
	}
	@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
	<div style="background: var(--color-bg); min-height: 100%; padding: var(--section-padding);">
		<div style="max-width: var(--content-max-width); margin: 0 auto;">
			<div class="news-masthead mb-8">
				<h1 style="font-family: var(--font-display); font-size: var(--font-size-display); font-weight: 900;">The Daily Style</h1>
				<p class="news-byline mt-1">A Reference Collection of Design Languages</p>
			</div>
			<p style="color: var(--color-text-muted); font-size: var(--font-size-caption);">Guide under construction...</p>
		</div>
	</div>
}
```

**Step 5: Update `handlers/guides.go`**

Add import:
```go
newspapertempl "github.com/johnfarrell/stylesheets/guides/newspaper"
```

Add switch case:
```go
case "newspaper":
    return newspapertempl.Page(g, htmxRequest)
```

**Step 6: Verify**

Run: `make templ && go build ./... && go test ./...`

**Step 7: Commit**

```bash
git add guides/registry.go guides/sources.go guides/newspaper/ handlers/guides.go
git commit -m "feat: register Newspaper guide with styles and placeholder page"
```

---

### Task 4: Build Terminal guide — full page

Build the complete Terminal guide page with all 9 sections. This is a large task — the templ file will be substantial. Also create the two server endpoints (SSE boot stream and command exec).

**Files:**
- Modify: `guides/terminal/terminal.templ` (full page — replace placeholder)
- Create: `guides/terminal/boot.templ` (SSE boot message templ component)
- Create: `guides/terminal/exec.templ` (command response templ component)
- Modify: `handlers/guides.go` (add `/guides/terminal/boot` SSE endpoint + `/guides/terminal/exec` endpoint)
- Modify: `handlers/guides_test.go` (add tests for new endpoints)

**Step 1: Write tests for the new endpoints**

Add to `handlers/guides_test.go`:

```go
func TestTerminalBootSSE(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/terminal/boot", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Errorf("expected text/event-stream, got %q", ct)
	}
}

func TestTerminalExec(t *testing.T) {
	mux := handlers.NewMux()
	body := strings.NewReader("cmd=help")
	req := httptest.NewRequest(http.MethodPost, "/guides/terminal/exec", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "help") {
		t.Errorf("expected response to contain help output")
	}
}
```

**Step 2: Build the full `guides/terminal/terminal.templ`**

Replace the placeholder with the full page. This file will contain:

- The guide header styled as a boot screen
- TechSummary callouts
- 6 required sections (Color Palette, Typography, Spacing, Buttons, Forms, Cards)
- 3 showcase sections (Live Boot SSE, Command Prompt, File Browser)
- SourceView for key interactive snippets

The guide should:
- Use `term-screen` wrapper for scanline overlay
- Use `term-glow` class on headings
- Use `term-panel` for all card/panel containers
- Use `term-btn` for all buttons
- Use `term-input` for all inputs

Key interactive sections:

**Section 7 — Live System Boot (SSE):**
```html
<div hx-ext="sse" sse-connect="/guides/terminal/boot" sse-swap="message" hx-swap="beforeend" id="term-boot-log"
     class="term-panel p-4 h-64 overflow-y-auto"
     style="font-size: var(--font-size-caption);">
</div>
<button class="term-btn mt-3" hx-get="/guides/terminal/boot" hx-target="#term-boot-log" hx-swap="innerHTML">[REBOOT]</button>
```
NOTE: The reboot button should disconnect the current SSE and reconnect. Use Alpine to toggle a boolean that controls the `sse-connect` attribute visibility.

**Section 8 — Command Prompt (HTMX + Alpine):**
```html
<div x-data="{ history: [], cmd: '' }" class="term-panel p-4">
    <div id="term-history" class="space-y-1 mb-2" style="font-size: var(--font-size-caption);">
        <template x-for="h in history">
            <div>
                <span style="color: var(--color-secondary);">$</span> <span x-text="h.cmd"></span>
                <div x-html="h.output" style="color: var(--color-text-muted);"></div>
            </div>
        </template>
    </div>
    <form hx-post="/guides/terminal/exec" hx-target="#term-cmd-response" hx-swap="innerHTML"
          @submit="history.push({cmd: cmd, output: '...'}); $nextTick(() => { cmd = '' })">
        <div class="flex items-center gap-2">
            <span style="color: var(--color-secondary);">$</span>
            <input type="text" name="cmd" x-model="cmd" class="term-input flex-1" autocomplete="off"
                   placeholder="type help, ls, whoami, or date"/>
        </div>
    </form>
    <div id="term-cmd-response" class="hidden"></div>
</div>
```
NOTE: The Alpine history display and HTMX response need coordination. The HTMX response replaces #term-cmd-response, and an Alpine watcher or hx-on::after-settle callback updates the last history entry's output. This is complex — implement the simpler version first where HTMX response just appends below the input, and Alpine just tracks the command text history.

**Section 9 — File Browser (Alpine keyboard nav):**
```html
<div x-data="{ files: ['README.md','main.go','go.mod','handlers/','guides/','static/','templates/','Makefile'], active: 0, preview: '' }"
     @keydown.up.prevent="active = Math.max(0, active - 1)"
     @keydown.down.prevent="active = Math.min(files.length - 1, active + 1)"
     @keydown.enter="preview = 'Contents of ' + files[active]"
     tabindex="0"
     class="term-panel p-4 focus:outline-none">
    <p style="color: var(--color-text-muted); font-size: var(--font-size-caption);" class="mb-2">
        Use ↑↓ arrows to navigate, Enter to open. Click panel first to focus.
    </p>
    <template x-for="(f, i) in files" :key="f">
        <div class="px-2 py-0.5" style="font-size: var(--font-size-body);"
             :class="i === active ? 'term-file-active' : ''"
             @click="active = i">
            <span x-text="f"></span>
        </div>
    </template>
    <div x-show="preview" class="mt-3 p-2" style="border-top: 1px solid var(--color-border); color: var(--color-text-muted); font-size: var(--font-size-caption);">
        <span x-text="preview"></span>
    </div>
</div>
```

**Step 3: Create `guides/terminal/boot.templ`**

```
package terminal

// BootMessage renders a single SSE boot log message.
templ BootMessage(timestamp, subsystem, message, color string) {
	<div style={ "font-size: var(--font-size-caption); color: " + color }>
		<span style="color: var(--color-text-muted);">[{ timestamp }]</span>
		<span style="font-weight: 700;">{ subsystem }</span>
		{ message }
	</div>
}
```

**Step 4: Create `guides/terminal/exec.templ`**

```
package terminal

// ExecResponse renders the output of a terminal command.
templ ExecResponse(cmd, output string) {
	<div style="font-size: var(--font-size-caption); margin-top: 0.25rem;">
		<div><span style="color: var(--color-secondary);">$</span> { cmd }</div>
		<div style="color: var(--color-text-muted); white-space: pre-wrap;">{ output }</div>
	</div>
}
```

**Step 5: Add SSE boot endpoint to `handlers/guides.go`**

```go
// Terminal — SSE boot sequence
mux.HandleFunc("/guides/terminal/boot", func(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming not supported", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    bootMessages := []struct{ sub, msg, color string }{
        {"BIOS", "POST check... OK", "var(--color-text)"},
        {"BIOS", "Memory: 640K conventional, 3072K extended", "var(--color-text)"},
        {"BOOT", "Loading kernel...", "var(--color-primary)"},
        {"KERN", "Initializing subsystems", "var(--color-primary)"},
        {"NET ", "eth0: link up 1000Mbps", "var(--color-secondary)"},
        {"DISK", "Mounting /dev/sda1 on /", "var(--color-text)"},
        {"DISK", "Filesystem clean — 847392 blocks free", "var(--color-text)"},
        {"AUTH", "Loading user credentials", "var(--color-accent)"},
        {"PROC", "Starting daemon processes", "var(--color-text)"},
        {"PROC", "sshd: listening on port 22", "var(--color-primary)"},
        {"PROC", "httpd: listening on port 8080", "var(--color-primary)"},
        {"NET ", "Firewall rules loaded (47 rules)", "var(--color-secondary)"},
        {"SYS ", "System clock synchronized via NTP", "var(--color-text)"},
        {"SYS ", "All systems nominal", "var(--color-primary)"},
        {"BOOT", "READY. Type 'help' for commands.", "var(--color-primary)"},
    }

    ts := time.Now()
    for i, m := range bootMessages {
        select {
        case <-r.Context().Done():
            return
        default:
        }
        bootTS := ts.Add(time.Duration(i*200) * time.Millisecond).Format("15:04:05.000")
        var buf bytes.Buffer
        terminaltempl.BootMessage(bootTS, m.sub, m.msg, m.color).Render(r.Context(), &buf)
        fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(buf.String(), "\n", ""))
        flusher.Flush()
        time.Sleep(300 * time.Millisecond)
    }
})
```

Add `"bytes"` to the imports in `handlers/guides.go`.

**Step 6: Add command exec endpoint to `handlers/guides.go`**

```go
// Terminal — command execution
mux.HandleFunc("/guides/terminal/exec", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    if err := r.ParseForm(); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    cmd := r.FormValue("cmd")
    var output string
    switch strings.TrimSpace(strings.ToLower(cmd)) {
    case "help":
        output = "Available commands: help, ls, whoami, date, clear\nType any command and press Enter."
    case "ls":
        output = "README.md  main.go  go.mod  go.sum  handlers/  guides/  static/  templates/  Makefile"
    case "whoami":
        output = "guest@stylesheets"
    case "date":
        output = time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
    case "clear":
        output = ""
    default:
        output = fmt.Sprintf("command not found: %s\nType 'help' for available commands.", templ.EscapeString(cmd))
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    terminaltempl.ExecResponse(cmd, output).Render(r.Context(), w)
})
```

**Step 7: Run tests and verify**

Run: `make templ && go build ./... && go test ./...`

**Step 8: Commit**

```bash
git add guides/terminal/ handlers/guides.go handlers/guides_test.go
git commit -m "feat: build complete Terminal guide with SSE boot, command exec, file browser

9 sections including live SSE boot stream, HTMX command prompt,
and Alpine keyboard-navigable file browser."
```

---

### Task 5: Build Retro OS guide — full page

Build the complete Retro OS guide with all 9 sections. Create the HTMX app-loading endpoints.

**Files:**
- Modify: `guides/retro/retro.templ` (full page — replace placeholder)
- Create: `guides/retro/apps.templ` (templ components for the 3 mini-apps: About, Calculator, File Manager)
- Modify: `handlers/guides.go` (add `/guides/retro/app/{name}` endpoint)
- Modify: `handlers/guides_test.go` (add test for app endpoint)

**Step 1: Write test for the app loading endpoint**

Add to `handlers/guides_test.go`:

```go
func TestRetroAppLoad(t *testing.T) {
	mux := handlers.NewMux()
	for _, app := range []string{"about", "calculator", "files"} {
		t.Run(app, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/guides/retro/app/"+app, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("GET /guides/retro/app/%s: expected 200, got %d", app, w.Code)
			}
		})
	}
}
```

**Step 2: Create `guides/retro/apps.templ`**

Three small components for the lazy-loaded window content:

```
package retro

// AppAbout renders the About window content.
templ AppAbout() {
	<div style="padding: 1rem; font-size: var(--font-size-body);">
		<div style="text-align: center; margin-bottom: 1rem;">
			<span style="font-size: 3rem;">🖥️</span>
			<h2 style="font-family: var(--font-display); font-size: 1.5rem; margin-top: 0.5rem;">Retro OS v1.0</h2>
		</div>
		<div class="retro-inset" style="padding: 0.5rem; font-size: var(--font-size-caption);">
			<p>A style guide showcasing the Win95 aesthetic.</p>
			<p style="margin-top: 0.5rem;">Built with Go · Templ · HTMX · Alpine.js</p>
			<p style="margin-top: 0.5rem; color: var(--color-text-muted);">Physical Memory: 640K</p>
		</div>
		<div style="text-align: center; margin-top: 1rem;">
			<button class="retro-btn retro-btn-primary" style="min-width: 80px;"
				{ templ.Attributes{"@click": "$dispatch('close-window', 'about')"}... }>OK</button>
		</div>
	</div>
}

// AppCalculator renders the Calculator window content.
templ AppCalculator() {
	<div style="padding: 0.25rem;" x-data="{ display: '0', op: '', prev: 0, reset: true }">
		<div class="retro-inset" style="padding: 0.25rem 0.5rem; text-align: right; font-family: var(--font-mono); font-size: 1.25rem; margin-bottom: 0.25rem;">
			<span x-text="display">0</span>
		</div>
		<div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 2px;">
			for _, row := range [][]string{{"7","8","9","/"}, {"4","5","6","*"}, {"1","2","3","-"}, {"C","0","=","+"}} {
				for _, btn := range row {
					if btn == "=" {
						<button class="retro-btn" style="padding: 0.4rem;"
							{ templ.Attributes{"@click": `
								let cur = parseFloat(display);
								if (op === '+') display = String(prev + cur);
								else if (op === '-') display = String(prev - cur);
								else if (op === '*') display = String(prev * cur);
								else if (op === '/') display = cur !== 0 ? String(prev / cur) : 'Err';
								op = ''; reset = true;
							`}... }
						>{ btn }</button>
					} else if btn == "C" {
						<button class="retro-btn" style="padding: 0.4rem;"
							{ templ.Attributes{"@click": "display = '0'; op = ''; prev = 0; reset = true"}... }
						>{ btn }</button>
					} else if btn == "+" || btn == "-" || btn == "*" || btn == "/" {
						<button class="retro-btn" style="padding: 0.4rem;"
							{ templ.Attributes{"@click": "prev = parseFloat(display); op = '" + btn + "'; reset = true"}... }
						>{ btn }</button>
					} else {
						<button class="retro-btn" style="padding: 0.4rem;"
							{ templ.Attributes{"@click": "if (reset) { display = '" + btn + "'; reset = false } else { display += '" + btn + "' }"}... }
						>{ btn }</button>
					}
				}
			}
		</div>
	</div>
}

// AppFiles renders the File Manager window content.
templ AppFiles() {
	<div style="font-size: var(--font-size-caption);">
		<div class="retro-raised" style="padding: 0.25rem 0.5rem; margin-bottom: 2px; font-size: 0.625rem;">
			<span style="text-decoration: underline;">F</span>ile
			<span style="margin-left: 0.5rem; text-decoration: underline;">E</span>dit
			<span style="margin-left: 0.5rem; text-decoration: underline;">V</span>iew
			<span style="margin-left: 0.5rem; text-decoration: underline;">H</span>elp
		</div>
		<div class="retro-inset" style="padding: 0.5rem;">
			for _, f := range []struct{icon, name, size string}{
				{"📄", "README.md", "2 KB"},
				{"📄", "main.go", "1 KB"},
				{"📄", "go.mod", "1 KB"},
				{"📁", "handlers", ""},
				{"📁", "guides", ""},
				{"📁", "static", ""},
				{"📁", "templates", ""},
				{"📄", "Makefile", "1 KB"},
			} {
				<div style="display: flex; align-items: center; gap: 0.5rem; padding: 0.15rem 0.25rem; cursor: default;"
					class="hover:bg-[#000080] hover:text-white">
					<span>{ f.icon }</span>
					<span style="flex: 1;">{ f.name }</span>
					<span style="color: var(--color-text-muted); min-width: 3rem; text-align: right;">{ f.size }</span>
				</div>
			}
		</div>
		<div class="retro-inset" style="padding: 0.15rem 0.5rem; margin-top: 2px; font-size: 0.625rem; color: var(--color-text-muted);">
			8 object(s)
		</div>
	</div>
}
```

**Step 3: Build the full `guides/retro/retro.templ`**

Replace the placeholder with the full page. Key sections:

- **The whole page** is styled as a desktop with teal background
- Required sections (1-6) are rendered inside "windows" with title bars
- Use `retro-window` + `retro-titlebar` classes
- Each window has [X] close and [_] minimize buttons via Alpine

**Section 7 — Draggable Windows:**
```html
<div x-data="{ x: 20, y: 20, dragging: false, offsetX: 0, offsetY: 0, z: 10 }"
     @mousedown.self="/* handled on titlebar only */"
     class="retro-window"
     :style="'position: absolute; left:' + x + 'px; top:' + y + 'px; z-index:' + z + '; width: 300px;'"
     @click="z = ++$store.desktop.topZ">
    <div class="retro-titlebar"
         @mousedown="dragging = true; offsetX = $event.clientX - x; offsetY = $event.clientY - y; z = ++$store.desktop.topZ"
         @mousemove.window="if (dragging) { x = $event.clientX - offsetX; y = $event.clientY - offsetY }"
         @mouseup.window="dragging = false">
        <span>Draggable Window</span>
        <div class="flex gap-0.5">
            <button class="retro-winbtn">_</button>
            <button class="retro-winbtn">X</button>
        </div>
    </div>
    <div style="padding: 1rem;">Window content here</div>
</div>
```

Use `Alpine.store('desktop', { topZ: 10 })` initialized in `x-data` on the page wrapper.

**Section 8 — Desktop Icons + HTMX Load:**
Desktop icons that use `@dblclick` to dispatch an event, which opens a window and triggers `hx-get` to load content.

**Section 9 — Taskbar:**
Bottom-fixed bar using `Alpine.store` to track which windows are open.

**Step 4: Add app-loading endpoint to `handlers/guides.go`**

```go
// Retro OS — lazy-load app window content
mux.HandleFunc("/guides/retro/app/{name}", func(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    switch name {
    case "about":
        retrotempl.AppAbout().Render(r.Context(), w)
    case "calculator":
        retrotempl.AppCalculator().Render(r.Context(), w)
    case "files":
        retrotempl.AppFiles().Render(r.Context(), w)
    default:
        http.NotFound(w, r)
    }
})
```

**Step 5: Run tests and verify**

Run: `make templ && go build ./... && go test ./...`

**Step 6: Commit**

```bash
git add guides/retro/ handlers/guides.go handlers/guides_test.go
git commit -m "feat: build complete Retro OS guide with draggable windows, desktop icons, taskbar

9 sections including draggable windows with z-stacking, HTMX lazy-loaded
mini-apps (About, Calculator, File Manager), and Alpine-powered taskbar."
```

---

### Task 6: Build Newspaper guide — full page

Build the complete Newspaper guide with all 9 sections. Create the infinite scroll and article view endpoints.

**Files:**
- Modify: `guides/newspaper/newspaper.templ` (full page — replace placeholder)
- Create: `guides/newspaper/headlines.templ` (headline card templ component for infinite scroll)
- Create: `guides/newspaper/article.templ` (article content templ component)
- Modify: `handlers/guides.go` (add `/guides/newspaper/headlines` and `/guides/newspaper/article/{id}` endpoints)
- Modify: `handlers/guides_test.go` (add tests)

**Step 1: Write tests**

Add to `handlers/guides_test.go`:

```go
func TestNewspaperHeadlines(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/newspaper/headlines?page=0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "news-card") {
		t.Error("expected headline cards in response")
	}
}

func TestNewspaperArticle(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/newspaper/article/0", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
```

**Step 2: Create `guides/newspaper/headlines.templ`**

```
package newspaper

// HeadlineCard renders a single headline in the infinite scroll feed.
templ HeadlineCard(id, category, headline, summary, byline string) {
	<div class="news-card mb-6">
		<p class="news-byline mb-1">{ category }</p>
		<h3 class="news-headline-md mb-2">
			<a style="cursor: pointer; text-decoration: none; color: inherit;"
				hx-get={ "/guides/newspaper/article/" + id }
				hx-target="#news-main-content"
				hx-swap="innerHTML transition:true"
				hx-push-url="false"
			>{ headline }</a>
		</h3>
		<p style="font-family: var(--font-body); font-size: var(--font-size-body); color: var(--color-text-muted); max-width: 55ch;">{ summary }</p>
		<p class="news-byline mt-2">{ byline }</p>
	</div>
}

// HeadlineSentinel renders the scroll sentinel that triggers loading the next page.
templ HeadlineSentinel(nextPage string) {
	<div hx-get={ "/guides/newspaper/headlines?page=" + nextPage }
		hx-trigger="revealed"
		hx-swap="outerHTML"
		style="height: 1px;">
	</div>
}
```

**Step 3: Create `guides/newspaper/article.templ`**

```
package newspaper

// Article renders a full article view.
templ Article(category, headline, byline, body string) {
	<div>
		<button class="news-btn mb-6"
			hx-get="/guides/newspaper/feed"
			hx-target="#news-main-content"
			hx-swap="innerHTML transition:true"
		>← Back to Front Page</button>
		<p class="news-byline mb-2">{ category }</p>
		<h1 class="news-headline-lg mb-3">{ headline }</h1>
		<p class="news-byline mb-6">{ byline }</p>
		<div class="news-columns-2 news-dropcap" style="font-family: var(--font-body); font-size: var(--font-size-body); line-height: 1.7;">
			{ body }
		</div>
	</div>
}
```

**Step 4: Build the full `guides/newspaper/newspaper.templ`**

Replace the placeholder with the full page. Key sections:

- Masthead at top (`news-masthead`)
- Required sections (1-6) with thin rule dividers between them
- Cards section uses newspaper grid layout (1 large + 2-3 smaller)
- Section 7 (infinite scroll) has a `#news-main-content` div that contains the initial headlines
- Section 8 (view transitions) — clicking any headline loads article via HTMX
- Section 9 (reading progress) — Alpine scroll tracker at page top

**Section 9 — Reading Progress Bar:**
```html
<div x-data="{ progress: 0 }"
     @scroll.window="progress = Math.min(100, Math.round(window.scrollY / (document.documentElement.scrollHeight - window.innerHeight) * 100))"
     class="news-progress"
     :style="'width:' + progress + '%'">
</div>
```
This goes at the very top of the page output (fixed position).

**Section 7 — Infinite Scroll Feed:**
```html
<div id="news-main-content">
    <!-- Initial headlines rendered server-side -->
    <!-- Then a sentinel div at the bottom triggers next page load -->
</div>
```

The sentinel uses `hx-trigger="revealed"` and `hx-swap="outerHTML"` — when scrolled into view, it replaces itself with more headlines plus a new sentinel for the next page.

**Step 5: Add headline feed endpoint to `handlers/guides.go`**

The endpoint serves paginated headlines. Use a slice of ~15 canned articles. Each page returns 3 headlines + a sentinel for the next page (until exhausted).

```go
// Newspaper — infinite scroll headlines
mux.HandleFunc("/guides/newspaper/headlines", func(w http.ResponseWriter, r *http.Request) {
    type headline struct{ id int; category, title, summary, byline string }
    allHeadlines := []headline{
        {0, "Design", "The Grid Is Dead, Long Live the Grid", "Modern layout systems have made the rigid grid obsolete — or have they?", "By Jane Chen · 8 min read"},
        {1, "Typography", "Why Your Font Choice Is Wrong", "A provocative look at the assumptions we make about type selection.", "By Marcus Webb · 5 min read"},
        // ... (add ~13 more canned headlines)
    }
    pageStr := r.URL.Query().Get("page")
    page := 0
    if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
        page = p
    }
    perPage := 3
    start := page * perPage
    if start >= len(allHeadlines) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        return // no more headlines
    }
    end := start + perPage
    if end > len(allHeadlines) {
        end = len(allHeadlines)
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    for _, h := range allHeadlines[start:end] {
        newspapertempl.HeadlineCard(strconv.Itoa(h.id), h.category, h.title, h.summary, h.byline).Render(r.Context(), w)
    }
    if end < len(allHeadlines) {
        newspapertempl.HeadlineSentinel(strconv.Itoa(page + 1)).Render(r.Context(), w)
    }
})
```

Add `"strconv"` to imports.

Also add a `/guides/newspaper/feed` endpoint that returns the initial set of headlines (for the "Back to Front Page" button):

```go
// Newspaper — initial feed (for back-to-front-page)
mux.HandleFunc("/guides/newspaper/feed", func(w http.ResponseWriter, r *http.Request) {
    // Redirect to headlines?page=0
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // Render initial 3 headlines + sentinel
    // (reuse the headlines handler logic or just redirect internally)
    http.Redirect(w, r, "/guides/newspaper/headlines?page=0", http.StatusSeeOther)
})
```

**Step 6: Add article view endpoint**

```go
// Newspaper — article view
mux.HandleFunc("/guides/newspaper/article/{id}", func(w http.ResponseWriter, r *http.Request) {
    // Lookup article by ID from the same canned data
    // Return the Article component
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // ... lookup + render
})
```

The article data can be a package-level var in `handlers/guides.go` or a helper function. Keep it simple — same canned data as headlines but with longer body text.

**Step 7: Run tests and verify**

Run: `make templ && go build ./... && go test ./...`

**Step 8: Commit**

```bash
git add guides/newspaper/ handlers/guides.go handlers/guides_test.go
git commit -m "feat: build complete Newspaper guide with infinite scroll, view transitions, reading progress

9 sections including HTMX infinite scroll feed, article view transitions,
and Alpine scroll-tracking reading progress bar."
```

---

### Task 7: Final verification and tailwind rebuild

Verify all guides work, tests pass, and tailwind picks up new classes.

**Files:**
- None (verification only)

**Step 1: Run full build and test**

Run: `make build && go test ./... -v`

**Step 2: Rebuild tailwind to pick up new classes**

Run: `make tailwind`

**Step 3: Visual spot-check (manual)**

Run: `make run` — open browser to `http://localhost:8080` and verify:
- All 9 guides appear in sidebar
- Terminal: scanlines visible, SSE boot streams, command prompt works, file browser has keyboard nav
- Retro OS: windows are draggable, desktop icons load apps via HTMX, taskbar toggles windows
- Newspaper: infinite scroll loads more headlines, clicking headline shows article with transition, reading progress bar moves on scroll

**Step 4: Commit if any fixes needed**

```bash
git commit -m "chore: final verification and tailwind rebuild for new guides"
```

---

## Task Dependency Summary

```
Task 0: SSE extension + FormResponse cases (independent)
Task 1: Register Terminal (depends on Task 0)
Task 2: Register Retro OS (depends on Task 0)
Task 3: Register Newspaper (depends on Task 0)
Task 4: Build Terminal page (depends on Task 1)
Task 5: Build Retro OS page (depends on Task 2)
Task 6: Build Newspaper page (depends on Task 3)
Task 7: Final verification (depends on Tasks 4, 5, 6)
```

Tasks 1-3 can be parallelized (they modify the same files but in additive ways — prefer sequential to avoid conflicts). Tasks 4-6 are independent of each other (separate directories) but each depends on its registration task.
