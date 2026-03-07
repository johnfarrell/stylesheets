# Production Polish Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Add a landing page, GitHub link, copyright notice, and custom 404 page to make the site production-ready.

**Architecture:** Three self-contained changes — sidebar footer gains two lines (GitHub link + copyright), a new `templates/home.templ` replaces the `/` redirect with a real landing page, and a new `templates/notfound.templ` replaces bare `http.NotFound` calls with a styled HTML page rendered inside `Layout`.

**Tech Stack:** Go, Templ v0.3.1001, HTMX 2.0.4, Tailwind CSS v4. Module: `github.com/johnfarrell/stylesheets`.

---

## Context for implementer

The project is a style guide showcase. The shell (sidebar + layout) is neutral gray/white. Guide content is styled per-guide. The landing page and 404 page live in the neutral shell.

**Key files:**
- `templates/sidebar.templ` — the nav sidebar; has a footer `<div class="p-4 border-t border-gray-200">` with one `<p>` line
- `templates/layout.templ` — `Layout(allGuides []guides.Guide, activeSlug string, fontURL string, content templ.Component)`
- `handlers/guides.go` — all routes; `/` handler currently redirects; `http.NotFound` called in 3 places
- `handlers/guides_test.go` — existing tests; `TestIndexRedirects` will need updating

**Build commands:**
- `make templ` — run after any `.templ` file change (`/home/john/go/bin/templ generate`)
- `make tailwind` — rebuild `static/css/output.css` after any template change (Tailwind scans source)
- `go test ./...` — run all tests
- `git commit -m "message"` — use plain git commit; if GPG error, retry with `git commit --no-gpg-sign -m "message"`

---

## Task 1: Sidebar — GitHub link + copyright

**Files:**
- Modify: `templates/sidebar.templ`

**Step 1: No test to write** — sidebar is visual-only; the existing handler tests still exercise it.

**Step 2: Edit templates/sidebar.templ**

Find the bottom `<div>` (currently lines 35–38):

```html
<div class="p-4 border-t border-gray-200">
    <p class="text-xs text-gray-400">Go · Templ · HTMX · Alpine · Tailwind</p>
</div>
```

Replace with:

```html
<div class="p-4 border-t border-gray-200 flex flex-col gap-2">
    <a
        href="https://github.com/johnfarrell/stylesheets"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-1.5 text-xs text-gray-400 hover:text-gray-700 transition-colors"
    >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-3.5 h-3.5 shrink-0">
            <path d="M12 0C5.37 0 0 5.373 0 12c0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61-.546-1.385-1.335-1.755-1.335-1.755-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23A11.52 11.52 0 0 1 12 5.803c1.02.005 2.047.138 3.006.404 2.29-1.552 3.297-1.23 3.297-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 21.795 24 17.298 24 12c0-6.627-5.373-12-12-12z"/>
        </svg>
        GitHub
    </a>
    <p class="text-xs text-gray-400">Go · Templ · HTMX · Alpine · Tailwind</p>
    <p class="text-xs text-gray-400">© 2026 John Farrell</p>
</div>
```

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

Expected: all tests pass (sidebar changes don't affect handler tests).

**Step 5: Commit**

```bash
git add templates/sidebar.templ templates/sidebar_templ.go
git commit -m "feat: add GitHub link and copyright notice to sidebar footer"
```

---

## Task 2: Landing page

**Files:**
- Create: `templates/home.templ`
- Modify: `handlers/guides.go` (the `/` handler)
- Modify: `handlers/guides_test.go` (update `TestIndexRedirects`)

**Step 1: Write the failing test**

In `handlers/guides_test.go`, replace `TestIndexRedirects` with `TestIndexRendersLandingPage`:

```go
func TestIndexRendersLandingPage(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html, got %q", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Brutalist") {
		t.Error("expected guide names in landing page body")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./handlers/... -run TestIndexRendersLandingPage -v
```

Expected: FAIL — `expected 200, got 302`.

**Step 3: Create templates/home.templ**

```go
package templates

import (
	"fmt"
	"github.com/johnfarrell/stylesheets/guides"
)

// Home renders the landing page listing all registered style guides.
templ Home(allGuides []guides.Guide) {
	<div class="p-8 max-w-5xl mx-auto">
		<div class="mb-10">
			<h1 class="text-3xl font-bold text-gray-900">Stylesheets</h1>
			<p class="mt-2 text-gray-500 max-w-xl">
				A reference collection of UI design languages built with Go, Templ, HTMX, Alpine.js, and Tailwind CSS.
			</p>
		</div>
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			for _, g := range allGuides {
				<a
					href={ templ.SafeURL(fmt.Sprintf("/guides/%s", g.Slug)) }
					hx-get={ fmt.Sprintf("/guides/%s/content", g.Slug) }
					hx-target="#content"
					hx-push-url={ fmt.Sprintf("/guides/%s", g.Slug) }
					class="block bg-white border border-gray-200 rounded p-6 hover:border-gray-400 hover:shadow-sm transition-all cursor-pointer"
				>
					<h2 class="text-base font-semibold text-gray-900">{ g.Name }</h2>
					<p class="mt-1 text-sm text-gray-500">{ g.Description }</p>
					<p class="mt-4 text-xs font-medium text-gray-700">View guide →</p>
				</a>
			}
		</div>
	</div>
}
```

**Step 4: Update the `/` handler in handlers/guides.go**

Find the `/` handler (lines 29–39):

```go
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
```

Replace with:

```go
// Landing page
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    page := templates.Layout(guides.All, "", "", templates.Home(guides.All))
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    templ.Handler(page).ServeHTTP(w, r)
})
```

**Step 5: Regenerate templ and rebuild Tailwind**

```bash
make templ
make tailwind
```

**Step 6: Run tests to verify they pass**

```bash
go test ./...
```

Expected: all tests PASS including `TestIndexRendersLandingPage`.

**Step 7: Commit**

```bash
git add templates/home.templ templates/home_templ.go handlers/guides.go handlers/guides_test.go
git commit -m "feat: replace root redirect with landing page listing all guides"
```

---

## Task 3: Custom 404 page

**Files:**
- Create: `templates/notfound.templ`
- Modify: `handlers/guides.go` (add `renderNotFound` helper, replace `http.NotFound` calls)
- Modify: `handlers/guides_test.go` (update `TestGuideNotFound`, add `TestUnknownPathIs404`)

**Step 1: Write the failing tests**

In `handlers/guides_test.go`, replace `TestGuideNotFound` with this updated version and add the new test:

```go
func TestGuideNotFound(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/does-not-exist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html for 404 page, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "404") {
		t.Error("expected '404' in response body")
	}
}

func TestUnknownPathIs404(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/this-does-not-exist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
```

**Step 2: Run tests to verify they fail**

```bash
go test ./handlers/... -run "TestGuideNotFound|TestUnknownPathIs404" -v
```

Expected: `TestGuideNotFound` FAIL — `expected text/html for 404 page` (current response is plain text), `TestUnknownPathIs404` PASS (already returns 404).

**Step 3: Create templates/notfound.templ**

```go
package templates

// NotFound renders the 404 error page.
templ NotFound() {
	<div class="flex flex-col items-center justify-center min-h-[60vh] text-center px-8">
		<p class="text-8xl font-bold text-gray-100 select-none">404</p>
		<h1 class="mt-2 text-2xl font-bold text-gray-900">Page not found</h1>
		<p class="mt-2 text-sm text-gray-500">The page you&#39;re looking for doesn&#39;t exist.</p>
		<a
			href="/"
			class="mt-6 text-sm font-medium text-gray-700 hover:text-gray-900 hover:underline transition-colors"
		>← Back to guides</a>
	</div>
}
```

**Step 4: Add renderNotFound helper and update http.NotFound calls in handlers/guides.go**

Add this helper function at the bottom of `handlers/guides.go` (after `placeholderContent`):

```go
// renderNotFound serves a styled 404 page inside the main layout.
func renderNotFound(w http.ResponseWriter, r *http.Request) {
	page := templates.Layout(guides.All, "", "", templates.NotFound())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templ.Handler(page, templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
}
```

Then replace the three `http.NotFound(w, r)` calls with `renderNotFound(w, r)`:

1. In the `/` handler (the `r.URL.Path != "/"` branch):
   ```go
   renderNotFound(w, r)
   ```

2. In the `/guides/{slug}` handler (the `!ok` branch):
   ```go
   renderNotFound(w, r)
   ```

3. In the `/guides/{slug}/content` handler (the `!ok` branch):
   ```go
   renderNotFound(w, r)
   ```

**Step 5: Regenerate templ and rebuild Tailwind**

```bash
make templ
make tailwind
```

**Step 6: Run all tests**

```bash
go test ./...
```

Expected: all tests PASS.

**Step 7: Commit**

```bash
git add templates/notfound.templ templates/notfound_templ.go handlers/guides.go handlers/guides_test.go
git commit -m "feat: add custom 404 page rendered inside layout"
```
