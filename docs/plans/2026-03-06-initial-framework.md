# Initial Framework Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Build the foundational framework for a data-driven style guide reference app using Go, Templ, HTMX, Alpine.js, and Tailwind CSS.

**Architecture:** A central guide registry holds all style guides as Go structs. The HTTP server auto-generates routes from the registry. HTMX swaps only the content area when navigating between guides; the sidebar stays fixed. Each guide's CSS variables are injected with the swapped content to apply its theme.

**Tech Stack:** Go 1.26 (stdlib net/http), Templ (HTML templating), HTMX (CDN), Alpine.js (CDN), Tailwind CSS CLI, Google Fonts

---

## Task 1: Initialize Go Module and Install Dependencies

**Files:**
- Create: `go.mod`
- Create: `go.sum` (auto-generated)

**Step 1: Initialize the Go module**

Run from `/home/john/projects/stylesheets`:
```bash
go mod init github.com/johnfarrell/stylesheets
```
Expected: `go.mod` created with `module github.com/johnfarrell/stylesheets` and `go 1.26`

**Step 2: Install Templ**

```bash
go get github.com/a-h/templ@latest
```
Expected: `go.mod` updated, `go.sum` created

**Step 3: Install the Templ CLI**

```bash
go install github.com/a-h/templ/cmd/templ@latest
```
Expected: `templ` binary available at `$(go env GOPATH)/bin/templ`
Verify: `templ version` prints a version number

**Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: initialize Go module with templ dependency"
```

---

## Task 2: Set Up Tailwind CSS CLI

**Files:**
- Create: `tailwind.config.js`
- Create: `static/css/input.css`
- Modify: `.gitignore` (add generated CSS)

**Step 1: Download the Tailwind CSS standalone CLI**

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
mv tailwindcss-linux-x64 tailwindcss
```
Expected: `tailwindcss` binary in project root

**Step 2: Create the Tailwind config**

Create `tailwind.config.js`:
```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ",
    "./**/*.go",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

**Step 3: Create the CSS input file**

Create `static/css/input.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**Step 4: Create the output directory and run initial build**

```bash
mkdir -p static/css
./tailwindcss -i static/css/input.css -o static/css/output.css
```
Expected: `static/css/output.css` created with compiled CSS

**Step 5: Add generated files and binary to .gitignore**

Add these lines to `.gitignore`:
```
# Tailwind
static/css/output.css
tailwindcss
```

**Step 6: Commit**

```bash
git add tailwind.config.js static/css/input.css .gitignore
git commit -m "feat: add Tailwind CSS CLI setup"
```

---

## Task 3: Create Makefile for Dev Workflow

**Files:**
- Create: `Makefile`

**Step 1: Create the Makefile**

Create `Makefile`:
```makefile
.PHONY: dev build templ tailwind clean

# Generate templ files and build
build: templ tailwind
	go build -o ./bin/stylesheets ./...

# Run templ code generator
templ:
	templ generate

# Build Tailwind CSS
tailwind:
	./tailwindcss -i static/css/input.css -o static/css/output.css --minify

# Watch mode for development (run each in a separate terminal)
watch-templ:
	templ generate --watch

watch-tailwind:
	./tailwindcss -i static/css/input.css -o static/css/output.css --watch

# Run the server
run: build
	./bin/stylesheets

# Clean generated files
clean:
	rm -f ./bin/stylesheets
	rm -f static/css/output.css
	find . -name "*_templ.go" -delete
```

**Step 2: Verify the Makefile works**

```bash
make templ
```
Expected: No error (nothing to generate yet — no .templ files exist)

**Step 3: Commit**

```bash
git add Makefile
git commit -m "feat: add Makefile for build and dev workflow"
```

---

## Task 4: Create Guide Registry

**Files:**
- Create: `guides/registry.go`

**Step 1: Write the test first**

Create `guides/registry_test.go`:
```go
package guides_test

import (
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
)

func TestRegistryNotEmpty(t *testing.T) {
	if len(guides.All) == 0 {
		t.Fatal("guide registry must not be empty")
	}
}

func TestGuideBySlug(t *testing.T) {
	guide, ok := guides.BySlug("brutalist")
	if !ok {
		t.Fatal("expected to find 'brutalist' guide")
	}
	if guide.Slug != "brutalist" {
		t.Errorf("expected slug 'brutalist', got %q", guide.Slug)
	}
}

func TestGuideBySlugNotFound(t *testing.T) {
	_, ok := guides.BySlug("does-not-exist")
	if ok {
		t.Fatal("expected BySlug to return false for unknown slug")
	}
}

func TestGuideHasRequiredFields(t *testing.T) {
	for _, g := range guides.All {
		if g.Name == "" {
			t.Errorf("guide with slug %q has empty Name", g.Slug)
		}
		if g.Slug == "" {
			t.Error("guide has empty Slug")
		}
		if g.FontURL == "" {
			t.Errorf("guide %q has empty FontURL", g.Slug)
		}
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./guides/...
```
Expected: FAIL — package does not exist yet

**Step 3: Create the registry**

Create `guides/registry.go`:
```go
package guides

// Guide defines a style guide's metadata and theme tokens.
type Guide struct {
	Name        string
	Slug        string
	Description string
	FontURL     string
	// CSSVars holds all per-guide CSS custom properties.
	// Any visual property that differs between guides belongs here:
	// colors, typography, radius, shadows, borders, layout tokens, etc.
	CSSVars map[string]string
}

// All is the ordered list of registered style guides.
// Add new guides here to register them with the application.
var All = []Guide{
	{
		Name:        "Brutalist",
		Slug:        "brutalist",
		Description: "Raw, functional, unapologetic design with heavy borders and stark contrast.",
		FontURL:     "https://fonts.googleapis.com/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400&display=swap",
		CSSVars: map[string]string{
			// Colors
			"--color-primary":    "#000000",
			"--color-secondary":  "#FF0000",
			"--color-accent":     "#FFFF00",
			"--color-bg":         "#FFFFFF",
			"--color-surface":    "#F5F5F5",
			"--color-text":       "#000000",
			"--color-text-muted": "#555555",
			// Typography
			"--font-display":      "'Space Mono', monospace",
			"--font-body":         "'Space Mono', monospace",
			"--font-size-display": "3.5rem",
			"--font-size-heading": "1.75rem",
			"--font-size-body":    "1rem",
			"--font-size-caption": "0.75rem",
			// Shape
			"--radius-sm": "0px",
			"--radius-md": "0px",
			"--radius-lg": "0px",
			// Elevation/Shadows
			"--shadow-card": "4px 4px 0px #000000",
			"--shadow-btn":  "3px 3px 0px #000000",
			// Borders
			"--border-width": "2px",
			"--border-color": "#000000",
			// Layout
			"--layout-columns":    "1",
			"--layout-gap":        "2rem",
			"--content-max-width": "900px",
			"--section-padding":   "3rem 2rem",
		},
	},
}

// BySlug looks up a guide by its URL slug.
func BySlug(slug string) (Guide, bool) {
	for _, g := range All {
		if g.Slug == slug {
			return g, true
		}
	}
	return Guide{}, false
}
```

**Step 4: Run tests to verify they pass**

```bash
go test ./guides/... -v
```
Expected: All 4 tests PASS

**Step 5: Commit**

```bash
git add guides/registry.go guides/registry_test.go
git commit -m "feat: add guide registry with Guide struct and BySlug lookup"
```

---

## Task 5: Create Layout and Sidebar Templates

**Files:**
- Create: `templates/layout.templ`
- Create: `templates/sidebar.templ`
- Create: `templates/layout_templ.go` (auto-generated — do not edit manually)
- Create: `templates/sidebar_templ.go` (auto-generated — do not edit manually)

**Step 1: Create the layout template**

Create `templates/layout.templ`:
```go
package templates

import "github.com/johnfarrell/stylesheets/guides"

templ Layout(allGuides []guides.Guide, activeSlug string, fontURL string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>Stylesheets — Style Guide Reference</title>
			<link rel="stylesheet" href="/static/css/output.css"/>
			<!-- Per-guide Google Font: swapped via HTMX out-of-band on navigation -->
			<div id="font-loader" hx-swap-oob="true">
				if fontURL != "" {
					<link rel="preconnect" href="https://fonts.googleapis.com"/>
					<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
					<link rel="stylesheet" href={ fontURL }/>
				}
			</div>
			<!-- HTMX -->
			<script src="https://unpkg.com/htmx.org@2.0.4" integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+" crossorigin="anonymous"></script>
			<!-- Alpine.js -->
			<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
		</head>
		<body class="flex h-screen bg-gray-100 overflow-hidden">
			<!-- Sidebar -->
			@Sidebar(allGuides, activeSlug)
			<!-- Main content area — HTMX swaps this on navigation -->
			<main id="content" class="flex-1 overflow-y-auto">
				{ children... }
			</main>
		</body>
	</html>
}
```

**Step 2: Create the sidebar template**

Create `templates/sidebar.templ`:
```go
package templates

import (
	"fmt"
	"github.com/johnfarrell/stylesheets/guides"
)

templ Sidebar(allGuides []guides.Guide, activeSlug string) {
	<nav
		class="w-64 bg-white border-r border-gray-200 flex flex-col h-full shrink-0"
		x-data={ fmt.Sprintf(`{ active: '%s' }`, activeSlug) }
	>
		<div class="p-6 border-b border-gray-200">
			<h1 class="text-lg font-bold text-gray-900 leading-tight">Stylesheets</h1>
			<p class="text-xs text-gray-500 mt-1">Style Guide Reference</p>
		</div>
		<ul class="flex-1 overflow-y-auto py-4">
			for _, g := range allGuides {
				<li>
					<a
						href={ templ.SafeURL(fmt.Sprintf("/guides/%s", g.Slug)) }
						hx-get={ fmt.Sprintf("/guides/%s/content", g.Slug) }
						hx-target="#content"
						hx-push-url={ fmt.Sprintf("/guides/%s", g.Slug) }
						hx-on:click={ fmt.Sprintf("active = '%s'", g.Slug) }
						class="flex flex-col px-6 py-3 text-sm transition-colors"
						:class={ fmt.Sprintf(`active === '%s' ? 'bg-gray-100 text-gray-900 font-semibold border-r-2 border-gray-900' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'`, g.Slug) }
					>
						<span>{ g.Name }</span>
						<span class="text-xs text-gray-400 font-normal mt-0.5">{ g.Description }</span>
					</a>
				</li>
			}
		</ul>
		<div class="p-4 border-t border-gray-200">
			<p class="text-xs text-gray-400">Go · Templ · HTMX · Alpine · Tailwind</p>
		</div>
	</nav>
}
```

**Step 3: Generate the Go code from templates**

```bash
make templ
```
Expected: `templates/layout_templ.go` and `templates/sidebar_templ.go` created. No errors.

**Step 4: Verify templates compile**

```bash
go build ./...
```
Expected: Builds successfully (no main package yet — that's fine, it will warn but not error if run as `go build ./templates/...`)

**Step 5: Commit**

```bash
git add templates/layout.templ templates/sidebar.templ templates/layout_templ.go templates/sidebar_templ.go
git commit -m "feat: add layout and sidebar Templ templates"
```

---

## Task 6: Create Shared Section Component and Tech Badges

**Files:**
- Create: `templates/components/section.templ`
- Create: `templates/components/section_templ.go` (auto-generated)

**Step 1: Create the shared section component**

Create `templates/components/section.templ`:
```go
package components

// TechBadge identifies which technology powers an interactive element.
type TechBadge int

const (
	BadgeNone   TechBadge = iota
	BadgeHTMX             // Server-driven interaction via HTMX
	BadgeAlpine           // Client-side UI via Alpine.js
	BadgeBoth             // Uses both HTMX and Alpine.js
)

// Section wraps a guide section with a consistent header and optional tech badge.
// Guides are free to use this wrapper or write their own layout entirely.
templ Section(title string, badge TechBadge) {
	<section class="mb-12">
		<div class="flex items-center gap-3 mb-6">
			<h2 class="text-lg font-semibold text-gray-700 uppercase tracking-widest">{ title }</h2>
			switch badge {
			case BadgeHTMX:
				<span class="px-2 py-0.5 text-xs font-mono bg-blue-100 text-blue-700 rounded border border-blue-200">[HTMX]</span>
			case BadgeAlpine:
				<span class="px-2 py-0.5 text-xs font-mono bg-green-100 text-green-700 rounded border border-green-200">[Alpine]</span>
			case BadgeBoth:
				<span class="px-2 py-0.5 text-xs font-mono bg-blue-100 text-blue-700 rounded border border-blue-200">[HTMX]</span>
				<span class="px-2 py-0.5 text-xs font-mono bg-green-100 text-green-700 rounded border border-green-200">[Alpine]</span>
			}
		</div>
		{ children... }
	</section>
}
```

**Step 2: Generate Go code from the template**

```bash
make templ
```
Expected: `templates/components/section_templ.go` created. No errors.

**Step 3: Verify compilation**

```bash
go build ./templates/...
```
Expected: Builds successfully.

**Step 4: Commit**

```bash
git add templates/components/section.templ templates/components/section_templ.go
git commit -m "feat: add shared Section component with HTMX/Alpine tech badges"
```

---

## Task 7: Create HTTP Server and Routes

**Files:**
- Create: `main.go`
- Create: `handlers/guides.go`
- Create: `handlers/guides_test.go`

**Step 1: Write handler tests first**

Create `handlers/guides_test.go`:
```go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/johnfarrell/stylesheets/handlers"
)

func TestIndexRedirects(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/guides/") {
		t.Errorf("expected redirect to /guides/*, got %q", loc)
	}
}

func TestGuidePageOK(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/brutalist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html content type, got %q", ct)
	}
}

func TestGuideContentPartialOK(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/brutalist/content", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGuideNotFound(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/does-not-exist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./handlers/...
```
Expected: FAIL — package does not exist yet

**Step 3: Create the handlers package**

Create `handlers/guides.go`:
```go
package handlers

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/johnfarrell/stylesheets/guides"
	guidetemplates "github.com/johnfarrell/stylesheets/guides/brutalist"
	"github.com/johnfarrell/stylesheets/templates"
)

// NewMux creates and returns the application HTTP mux with all routes registered.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Root redirect to first guide
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if len(guides.All) == 0 {
			http.Error(w, "no guides registered", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/guides/"+guides.All[0].Slug, http.StatusFound)
	})

	// Full page guide render
	mux.HandleFunc("/guides/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		guide, ok := guides.BySlug(slug)
		if !ok {
			http.NotFound(w, r)
			return
		}
		content := guideContent(guide)
		page := templates.Layout(guides.All, guide.Slug, guide.FontURL, content)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templ.Handler(page).ServeHTTP(w, r)
	})

	// HTMX partial content swap
	mux.HandleFunc("/guides/{slug}/content", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		guide, ok := guides.BySlug(slug)
		if !ok {
			http.NotFound(w, r)
			return
		}
		partial := guideContent(guide)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templ.Handler(partial).ServeHTTP(w, r)
	})

	return mux
}

// guideContent returns the Templ component for the given guide's showcase.
// Add a case here when registering a new guide.
func guideContent(g guides.Guide) templ.Component {
	switch g.Slug {
	case "brutalist":
		return guidetemplates.Page(g)
	default:
		return guidetemplates.Page(g) // fallback to brutalist until more guides exist
	}
}

// cssVarsBlock generates an inline <style> block injecting a guide's CSS variables.
func cssVarsBlock(vars map[string]string) string {
	var sb strings.Builder
	sb.WriteString(":root {\n")
	for k, v := range vars {
		sb.WriteString("  ")
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString(";\n")
	}
	sb.WriteString("}")
	return sb.String()
}
```

**Note:** The `layout.templ` signature needs a `content templ.Component` parameter. Update `templates/layout.templ` to accept it:

```go
templ Layout(allGuides []guides.Guide, activeSlug string, fontURL string, content templ.Component) {
    ...
    <main id="content" class="flex-1 overflow-y-auto">
        @content
    </main>
    ...
}
```

Regenerate after editing: `make templ`

**Step 4: Create main.go**

Create `main.go`:
```go
package main

import (
	"log"
	"net/http"

	"github.com/johnfarrell/stylesheets/handlers"
)

func main() {
	mux := handlers.NewMux()
	addr := ":8080"
	log.Printf("Starting server on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
```

**Step 5: Run tests to verify they pass**

```bash
go test ./handlers/... -v
```
Expected: All 4 tests PASS

**Step 6: Commit**

```bash
git add handlers/guides.go handlers/guides_test.go main.go
git commit -m "feat: add HTTP server, routing, and guide handlers"
```

---

## Task 8: Create the Brutalist Guide Showcase

**Files:**
- Create: `guides/brutalist/brutalist.templ`
- Create: `guides/brutalist/brutalist_templ.go` (auto-generated)

This is the first full guide implementation. It validates the entire stack end-to-end and serves as the reference template for future guides.

**Step 1: Create the brutalist guide template**

Create `guides/brutalist/brutalist.templ`:
```go
package brutalist

import (
	"fmt"
	"github.com/johnfarrell/stylesheets/guides"
	"github.com/johnfarrell/stylesheets/templates/components"
)

// Page renders the full Brutalist style guide showcase.
templ Page(g guides.Guide) {
	<!-- Inject this guide's CSS variables and font -->
	<!-- [custom] - CSS vars for per-guide theming cannot be expressed as static Tailwind classes -->
	<style>
		:root {
			for k, v := range g.CSSVars {
				{ k }: { v };
			}
		}
		.font-display { font-family: var(--font-display); }
		.font-body    { font-family: var(--font-body); }
		.btn-primary {
			background: var(--color-primary);
			color: var(--color-bg);
			border: var(--border-width) solid var(--border-color);
			/* [custom] - per-guide box-shadow token */
			box-shadow: var(--shadow-btn);
			font-family: var(--font-body);
		}
		.btn-primary:hover {
			/* [custom] - brutalist shift effect on hover */
			transform: translate(-2px, -2px);
			box-shadow: 5px 5px 0px var(--border-color);
		}
		.btn-secondary {
			background: var(--color-bg);
			color: var(--color-primary);
			border: var(--border-width) solid var(--border-color);
			box-shadow: var(--shadow-btn);
			font-family: var(--font-body);
		}
		.card {
			background: var(--color-surface);
			border: var(--border-width) solid var(--border-color);
			/* [custom] - per-guide card shadow token */
			box-shadow: var(--shadow-card);
		}
	</style>
	<!-- Out-of-band font loader swap -->
	<div id="font-loader" hx-swap-oob="true">
		<link rel="preconnect" href="https://fonts.googleapis.com"/>
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
		<link rel="stylesheet" href={ g.FontURL }/>
	</div>
	<div class="p-8 max-w-[900px] mx-auto font-body" style="font-family: var(--font-body)">
		<!-- Guide header -->
		<div class="mb-12 pb-8 border-b-2 border-black">
			<h1 class="font-display text-5xl font-bold uppercase tracking-tight mb-2" style="font-family: var(--font-display)">
				{ g.Name }
			</h1>
			<p class="text-gray-600">{ g.Description }</p>
		</div>
		<!-- 1. Color Palette [Alpine] -->
		@components.Section("Color Palette", components.BadgeAlpine) {
			<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
				@colorSwatch("Primary", "var(--color-primary)", "#000000")
				@colorSwatch("Secondary", "var(--color-secondary)", "#FF0000")
				@colorSwatch("Accent", "var(--color-accent)", "#FFFF00")
				@colorSwatch("Background", "var(--color-bg)", "#FFFFFF")
				@colorSwatch("Surface", "var(--color-surface)", "#F5F5F5")
				@colorSwatch("Text", "var(--color-text)", "#000000")
				@colorSwatch("Text Muted", "var(--color-text-muted)", "#555555")
			</div>
		}
		<!-- 2. Typography -->
		@components.Section("Typography", components.BadgeNone) {
			<div class="space-y-6">
				<div>
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-1">Display / 3.5rem / Bold</p>
					<p class="font-bold border-b-2 border-black pb-2" style="font-family: var(--font-display); font-size: var(--font-size-display)">The Quick Brown Fox</p>
				</div>
				<div>
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-1">Heading / 1.75rem / Bold</p>
					<p class="font-bold" style="font-family: var(--font-display); font-size: var(--font-size-heading)">Style Guide Reference</p>
				</div>
				<div>
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-1">Body / 1rem / Regular</p>
					<p style="font-family: var(--font-body); font-size: var(--font-size-body)">The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs.</p>
				</div>
				<div>
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-1">Caption / 0.75rem / Regular</p>
					<p class="text-gray-500" style="font-family: var(--font-body); font-size: var(--font-size-caption)">Caption text — used for labels, metadata, and supplementary information.</p>
				</div>
			</div>
		}
		<!-- 3. Spacing -->
		@components.Section("Spacing Scale", components.BadgeNone) {
			<div class="space-y-3">
				for _, s := range []struct{ label, size string }{
					{"4px (base)", "4px"}, {"8px", "8px"}, {"16px", "16px"},
					{"24px", "24px"}, {"32px", "32px"}, {"48px", "48px"}, {"64px", "64px"},
				} {
					<div class="flex items-center gap-4">
						<span class="text-xs font-mono text-gray-500 w-24">{ s.label }</span>
						<div class="bg-black h-4 border border-black" style={ fmt.Sprintf("width: %s", s.size) }></div>
					</div>
				}
			</div>
		}
		<!-- 4. Buttons [Alpine] -->
		@components.Section("Buttons", components.BadgeAlpine) {
			<div class="space-y-6">
				<div>
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-3">States</p>
					<div class="flex flex-wrap gap-4">
						<button class="btn-primary px-6 py-3 font-bold uppercase tracking-wide transition-all">Primary</button>
						<button class="btn-secondary px-6 py-3 font-bold uppercase tracking-wide transition-all">Secondary</button>
						<button disabled class="px-6 py-3 font-bold uppercase tracking-wide bg-gray-200 text-gray-400 border-2 border-gray-300 cursor-not-allowed">Disabled</button>
					</div>
				</div>
				<div>
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-3">Sizes</p>
					<div class="flex flex-wrap items-center gap-4">
						<button class="btn-primary px-3 py-1.5 text-sm font-bold uppercase tracking-wide transition-all">Small</button>
						<button class="btn-primary px-6 py-3 font-bold uppercase tracking-wide transition-all">Medium</button>
						<button class="btn-primary px-8 py-4 text-lg font-bold uppercase tracking-wide transition-all">Large</button>
					</div>
				</div>
				<!-- Alpine.js toggle example -->
				<div x-data="{ toggled: false }">
					<p class="text-xs uppercase tracking-widest text-gray-400 mb-3">Alpine.js Toggle State</p>
					<button
						class="btn-primary px-6 py-3 font-bold uppercase tracking-wide transition-all"
						x-on:click="toggled = !toggled"
						:class="toggled ? 'bg-red-600 border-red-600' : ''"
					>
						<span x-text="toggled ? 'ACTIVE' : 'TOGGLE ME'">TOGGLE ME</span>
					</button>
					<p class="text-xs text-gray-400 mt-2 font-mono">x-data / x-on:click / x-text — no server round-trip</p>
				</div>
			</div>
		}
		<!-- 5. Forms [HTMX + Alpine] -->
		@components.Section("Forms", components.BadgeBoth) {
			<div class="space-y-6 max-w-lg">
				<form
					hx-post="/guides/brutalist/demo-form"
					hx-target="#form-response"
					hx-swap="innerHTML"
					class="space-y-4"
				>
					<div>
						<label class="block text-xs uppercase tracking-widest text-gray-600 mb-1">Text Input</label>
						<input
							type="text"
							name="name"
							placeholder="Enter something raw..."
							class="w-full px-4 py-3 border-2 border-black bg-white font-mono focus:outline-none focus:ring-2 focus:ring-black transition-shadow"
						/>
					</div>
					<div>
						<label class="block text-xs uppercase tracking-widest text-gray-600 mb-1">Select</label>
						<select class="w-full px-4 py-3 border-2 border-black bg-white font-mono focus:outline-none">
							<option>Option Alpha</option>
							<option>Option Beta</option>
							<option>Option Gamma</option>
						</select>
					</div>
					<div class="space-y-2">
						<label class="block text-xs uppercase tracking-widest text-gray-600">Checkboxes</label>
						<label class="flex items-center gap-3 cursor-pointer">
							<input type="checkbox" class="w-4 h-4 border-2 border-black accent-black"/>
							<span class="font-mono text-sm">Option One</span>
						</label>
						<label class="flex items-center gap-3 cursor-pointer">
							<input type="checkbox" class="w-4 h-4 border-2 border-black accent-black"/>
							<span class="font-mono text-sm">Option Two</span>
						</label>
					</div>
					<div class="space-y-2">
						<label class="block text-xs uppercase tracking-widest text-gray-600">Radio Buttons</label>
						<label class="flex items-center gap-3 cursor-pointer">
							<input type="radio" name="choice" class="w-4 h-4 border-2 border-black accent-black"/>
							<span class="font-mono text-sm">Choice A</span>
						</label>
						<label class="flex items-center gap-3 cursor-pointer">
							<input type="radio" name="choice" class="w-4 h-4 border-2 border-black accent-black"/>
							<span class="font-mono text-sm">Choice B</span>
						</label>
					</div>
					<button type="submit" class="btn-primary w-full py-3 font-bold uppercase tracking-wide transition-all">
						Submit via HTMX
					</button>
				</form>
				<div id="form-response" class="font-mono text-sm text-gray-600"></div>
				<p class="text-xs text-gray-400 font-mono">hx-post / hx-target / hx-swap — server handles submission, updates #form-response</p>
			</div>
		}
		<!-- 6. Cards [Alpine] -->
		@components.Section("Cards & Panels", components.BadgeAlpine) {
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div class="card p-6">
					<h3 class="font-bold text-lg uppercase mb-2" style="font-family: var(--font-display)">Basic Card</h3>
					<p class="text-sm text-gray-600">Content container with border, shadow, and surface color defined by CSS variables.</p>
				</div>
				<div class="card p-6" x-data="{ expanded: false }">
					<div class="flex justify-between items-start">
						<h3 class="font-bold text-lg uppercase" style="font-family: var(--font-display)">Expandable</h3>
						<button
							class="text-xs font-mono border-2 border-black px-2 py-1 hover:bg-black hover:text-white transition-colors"
							x-on:click="expanded = !expanded"
							x-text="expanded ? '− COLLAPSE' : '+ EXPAND'"
						>+ EXPAND</button>
					</div>
					<div x-show="expanded" x-collapse class="mt-4 text-sm text-gray-600 border-t-2 border-black pt-4">
						Hidden content revealed via Alpine.js x-show and x-collapse. No server round-trip.
					</div>
				</div>
			</div>
		}
	</div>
}

// colorSwatch renders a single color swatch with copy-to-clipboard via Alpine.js.
templ colorSwatch(name, cssVar, hex string) {
	<div
		x-data={ fmt.Sprintf(`{ copied: false, hex: '%s' }`, hex) }
		class="cursor-pointer group"
		x-on:click="navigator.clipboard.writeText(hex); copied = true; setTimeout(() => copied = false, 1500)"
	>
		<div
			class="h-16 border-2 border-black mb-2 transition-transform group-hover:-translate-y-1"
			style={ fmt.Sprintf("background: %s", cssVar) }
		></div>
		<p class="text-xs font-bold uppercase">{ name }</p>
		<p class="text-xs font-mono text-gray-500" x-text="copied ? 'Copied!' : hex">{ hex }</p>
	</div>
}
```

**Step 2: Add the demo form route to handlers**

Add to `handlers/guides.go` inside `NewMux()`:
```go
// Demo form endpoint for showcasing HTMX form submission
mux.HandleFunc("/guides/{slug}/demo-form", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    r.ParseForm()
    name := r.FormValue("name")
    if name == "" {
        name = "anonymous"
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, `<div class="border-2 border-black p-3 bg-yellow-50 font-mono">✓ Received: <strong>%s</strong></div>`, templ.EscapeString(name))
})
```

**Step 3: Generate templates**

```bash
make templ
```
Expected: `guides/brutalist/brutalist_templ.go` created. No errors.

**Step 4: Build the full project**

```bash
go build ./...
```
Expected: Builds successfully with no errors.

**Step 5: Build Tailwind CSS**

```bash
make tailwind
```
Expected: `static/css/output.css` generated with all Tailwind classes used in templates.

**Step 6: Run the server and verify manually**

```bash
go run main.go
```
Then open `http://localhost:8080` in a browser. Verify:
- [ ] Redirects to `/guides/brutalist`
- [ ] Sidebar shows "Brutalist" entry, highlighted as active
- [ ] Full guide showcase renders with correct styles
- [ ] Color swatches copy hex value to clipboard on click
- [ ] Buttons have hover shift effect
- [ ] Toggle button switches state with Alpine.js
- [ ] Form submits via HTMX and shows response without page reload
- [ ] Expandable card expands/collapses with Alpine.js

**Step 7: Run all tests**

```bash
go test ./... -v
```
Expected: All tests PASS.

**Step 8: Commit**

```bash
git add guides/brutalist/brutalist.templ guides/brutalist/brutalist_templ.go handlers/guides.go
git commit -m "feat: add Brutalist style guide with full interactive showcase"
```

---

## Task 9: Add a Second Guide to Validate Data-Driven Approach

**Files:**
- Create: `guides/minimal/minimal.templ`
- Create: `guides/minimal/minimal_templ.go` (auto-generated)
- Modify: `guides/registry.go` (add Minimal guide entry)
- Modify: `handlers/guides.go` (add case in `guideContent`)

**Goal:** Prove that adding a second guide only requires: (1) a registry entry, (2) a Templ file, (3) a case in the switch. No routing changes, no layout changes.

**Step 1: Add Minimal guide to registry**

In `guides/registry.go`, append to `All`:
```go
{
    Name:        "Minimal",
    Slug:        "minimal",
    Description: "Calm, spacious, single-column design with generous whitespace.",
    FontURL:     "https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&display=swap",
    CSSVars: map[string]string{
        // Colors
        "--color-primary":    "#1a1a1a",
        "--color-secondary":  "#6b7280",
        "--color-accent":     "#3b82f6",
        "--color-bg":         "#fafafa",
        "--color-surface":    "#ffffff",
        "--color-text":       "#1a1a1a",
        "--color-text-muted": "#9ca3af",
        // Typography
        "--font-display":      "'Inter', sans-serif",
        "--font-body":         "'Inter', sans-serif",
        "--font-size-display": "3rem",
        "--font-size-heading": "1.5rem",
        "--font-size-body":    "1rem",
        "--font-size-caption": "0.8rem",
        // Shape
        "--radius-sm": "4px",
        "--radius-md": "8px",
        "--radius-lg": "16px",
        // Elevation
        "--shadow-card": "0 1px 3px rgba(0,0,0,0.08)",
        "--shadow-btn":  "0 1px 2px rgba(0,0,0,0.05)",
        // Borders
        "--border-width": "1px",
        "--border-color": "#e5e7eb",
        // Layout — single column, generous spacing
        "--layout-columns":    "1",
        "--layout-gap":        "4rem",
        "--content-max-width": "640px",
        "--section-padding":   "5rem 2rem",
    },
},
```

**Step 2: Create the minimal guide Templ file**

Create `guides/minimal/minimal.templ` — follow the same structure as brutalist but styled for the minimal aesthetic. The layout should be single column, wide spacing, soft colors, rounded corners. All the same required sections apply.

**Step 3: Wire up in handlers**

In `handlers/guides.go`, add an import and a case:
```go
import minimaltempl "github.com/johnfarrell/stylesheets/guides/minimal"

// in guideContent switch:
case "minimal":
    return minimaltempl.Page(g)
```

**Step 4: Generate and build**

```bash
make templ && go build ./...
```
Expected: No errors.

**Step 5: Run all tests**

```bash
go test ./... -v
```
Expected: All tests PASS (registry test now finds 2 guides, slug tests still pass).

**Step 6: Manual verification**

Run `go run main.go` and verify:
- [ ] Sidebar shows both "Brutalist" and "Minimal"
- [ ] Clicking "Minimal" swaps content via HTMX (no full page reload)
- [ ] Font changes to Inter
- [ ] CSS variables change (rounded corners, soft shadows, generous spacing)
- [ ] URL updates to `/guides/minimal`
- [ ] Browser back button returns to Brutalist correctly

**Step 7: Commit**

```bash
git add guides/registry.go guides/minimal/ handlers/guides.go
git commit -m "feat: add Minimal style guide, validate data-driven guide registration"
```

---

## Summary

After completing all 9 tasks:

- Working Go web server with auto-generated routes from the guide registry
- HTMX-powered navigation with content-only swaps and URL history
- Per-guide theming via CSS custom properties (colors, fonts, radius, shadows, borders, layout)
- Two complete style guides (Brutalist, Minimal) demonstrating full stack contrast
- Shared Section component with HTMX/Alpine tech badges
- All tests passing
- Adding a new guide = registry entry + Templ file + one switch case
