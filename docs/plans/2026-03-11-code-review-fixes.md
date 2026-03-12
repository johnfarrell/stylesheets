# Code Review Fixes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Fix all issues found in the comprehensive code review — 5 important issues and 5 suggestions.

**Architecture:** Mostly targeted edits across handlers, templates, registry, and Dockerfile. The biggest structural change is extracting static data from handlers into guide packages and adding a `PageFunc` to the Guide struct to eliminate the `guideContent()` switch. The responsive sidebar requires a new Alpine-based mobile hamburger menu.

**Tech Stack:** Go 1.25, Templ, HTMX, Alpine.js, Tailwind CSS v4

---

### Task 1: Fix unchecked `.Render()` errors in handlers/guides.go

**Files:**
- Modify: `handlers/guides.go` (all `.Render()` call sites)

**Step 1: Add `log/slog` import and wrap all Render calls**

Add `log/slog` to imports. Then wrap every unchecked `.Render(r.Context(), w)` call with error logging. There are ~18 call sites across lines 128, 149, 165, 170, 174, 180, 197, 239, 273, 282, 284, 286, 332, 335, 358, 374, 384, 410.

Pattern — replace:
```go
component.Render(r.Context(), w)
```
with:
```go
if err := component.Render(r.Context(), w); err != nil {
    slog.Error("render failed", "error", err)
}
```

For the SSE boot handler (line 239), the render writes to a `bytes.Buffer`, so it should be:
```go
if err := terminaltempl.BootMessage(bootTS, m.sub, m.msg, m.color).Render(r.Context(), &buf); err != nil {
    slog.Error("render boot message", "error", err)
    continue
}
```

**Step 2: Run linter to verify errcheck findings are resolved**

Run: `golangci-lint run ./handlers/...`
Expected: No errcheck violations on Render calls.

**Step 3: Run tests**

Run: `go test ./handlers/... -v`
Expected: All tests pass.

**Step 4: Commit**

```bash
git add handlers/guides.go
git commit -m "Fix unchecked Render errors in handlers"
```

---

### Task 2: Remove leftover test copy from home.templ

**Files:**
- Modify: `templates/home.templ:16-18`

**Step 1: Remove the test paragraph**

Delete lines 16-18:
```templ
<p class="mt-2 text-gray-500 max-w-xl">
This is a test commit to see if cloud build works.
</p>
```

**Step 2: Regenerate templ and run tests**

Run: `make templ && go test ./handlers/... -v`
Expected: All tests pass. Landing page no longer shows test text.

**Step 3: Commit**

```bash
git add templates/home.templ templates/home_templ.go
git commit -m "Remove leftover test copy from landing page"
```

---

### Task 3: Add SRI hash to HTMX SSE extension script

**Files:**
- Modify: `templates/layout.templ:15`

**Step 1: Generate SRI hash for htmx-ext-sse@2.2.2/sse.js**

Run: `curl -s https://unpkg.com/htmx-ext-sse@2.2.2/sse.js | openssl dgst -sha384 -binary | openssl base64 -A`

**Step 2: Add integrity + crossorigin attributes**

Replace line 15:
```templ
<script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
```
with:
```templ
<script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js" integrity="sha384-{HASH}" crossorigin="anonymous"></script>
```

**Step 3: Regenerate templ and run tests**

Run: `make templ && go test ./handlers/... -v`
Expected: All pass. The SSE-using terminal boot test still works.

**Step 4: Commit**

```bash
git add templates/layout.templ templates/layout_templ.go
git commit -m "Add SRI hash to HTMX SSE extension script"
```

---

### Task 4: Align Dockerfile Go version with go.mod

**Files:**
- Modify: `Dockerfile:25`

**Step 1: Change golang base image version**

Replace:
```dockerfile
FROM golang:1.26-alpine AS go-builder
```
with:
```dockerfile
FROM golang:1.25-alpine AS go-builder
```

**Step 2: Verify Docker build**

Run: `make docker-build` (if Docker available, otherwise visual check is fine)

**Step 3: Commit**

```bash
git add Dockerfile
git commit -m "Align Dockerfile Go version with go.mod (1.25)"
```

---

### Task 5: Extract static data from handlers into guide packages

**Files:**
- Create: `guides/tracker/data.go`
- Create: `guides/newspaper/data.go`
- Modify: `handlers/guides.go`

**Step 1: Create guides/tracker/data.go**

Move the `trackerItem` struct and `trackerItems` slice from `handlers/guides.go` (lines 27-64) into `guides/tracker/data.go`, exported as `Item` struct and `Items` slice.

```go
package tracker

// Item represents a single tracker entry (skill, project, certification, or challenge).
type Item struct {
	ID           string
	Category     string
	Name         string
	Status       string
	Level        int
	Target       int
	Description  string
	Requirements []string
	Unlocks      []string
}

// Items is the full list of tracker entries displayed in the guide.
var Items = []Item{
    // ... (move all 20 items here)
}
```

**Step 2: Create guides/newspaper/data.go**

Move headline and article data into `guides/newspaper/data.go`:

```go
package newspaper

// Headline represents a newspaper headline entry.
type Headline struct {
	ID                               int
	Category, Title, Summary, Byline string
}

// Headlines is the full list of newspaper headlines for infinite scroll.
var Headlines = []Headline{
    // ... (move all 15 headlines here)
}

// Article represents a full newspaper article.
type Article struct {
	Category, Title, Byline, Body string
}

// Articles maps article IDs to their full content.
var Articles = map[string]Article{
    // ... (move all 5 articles here)
}
```

**Step 3: Update handlers/guides.go to use the new packages**

- Remove `trackerItem` struct and `trackerItems` var (lines 27-64)
- Remove inline `allHeadlines` and `articles` from their respective handlers
- Update references: `trackerItems` → `trackertempl.Items` (since it's already aliased), and headlines/articles similarly
- Note: `trackertempl` alias already exists for the tracker package, so `trackertempl.Items` works. For newspaper, `newspapertempl` alias already exists, so `newspapertempl.Headlines` and `newspapertempl.Articles` work.

**Step 4: Run linter and tests**

Run: `golangci-lint run ./... && go test ./... -v`
Expected: All pass. No lint issues.

**Step 5: Commit**

```bash
git add guides/tracker/data.go guides/newspaper/data.go handlers/guides.go
git commit -m "Extract inline static data to guide packages"
```

---

### Task 6: Add PageFunc to Guide struct and eliminate guideContent() switch

**Files:**
- Modify: `guides/registry.go` (add PageFunc field + populate in All)
- Modify: `handlers/guides.go` (remove guideContent, remove 10 import aliases, use PageFunc)

**Step 1: Add PageFunc to Guide struct**

In `guides/registry.go`, add to the struct:
```go
type Guide struct {
	Name        string
	Slug        string
	Description string
	FontURL     string
	CSSVars     map[string]string
	// PageFunc renders the guide's full showcase page.
	// Nil means the guide uses a placeholder.
	PageFunc    func(Guide, bool) templ.Component
}
```

Add `"github.com/a-h/templ"` to imports.

**Step 2: Write a test that verifies all guides have PageFunc set**

Add to `guides/registry_test.go` (or existing test file):
```go
func TestAllGuidesHavePageFunc(t *testing.T) {
	for _, g := range All {
		if g.PageFunc == nil {
			t.Errorf("guide %q has nil PageFunc", g.Slug)
		}
	}
}
```

Run: `go test ./guides/... -run TestAllGuidesHavePageFunc -v`
Expected: FAIL (PageFunc not yet set).

**Step 3: Populate PageFunc in each guide's All entry**

This requires importing each guide's templ package into `registry.go`. Since registry.go is in the `guides` package and guide templates are in sub-packages (`guides/brutalist`, etc.), this creates an import cycle.

**Alternative approach:** Instead of putting PageFunc in registry.go directly, keep the function mapping in handlers where the imports already exist. Create a `PageFunc` field on Guide, but set it in `handlers/guides.go` via an `init()` or in `NewMux()`.

Better: Create a `RegisterPages()` function in handlers that sets PageFunc on each guide:

In `handlers/guides.go`:
```go
func init() {
	for i := range guides.All {
		switch guides.All[i].Slug {
		case "brutalist":
			guides.All[i].PageFunc = brutalisttempl.Page
		case "cassette":
			guides.All[i].PageFunc = cassettetempl.Page
		// ... all 10
		}
	}
}
```

Then replace `guideContent(g, htmxRequest)` calls with `g.PageFunc(g, htmxRequest)` and delete the `guideContent` function.

**Step 4: Run test again**

Run: `go test ./guides/... -run TestAllGuidesHavePageFunc -v`
Expected: PASS (handlers init() sets all PageFuncs).

Wait — `guides` tests won't import `handlers`, so the init() won't run during `guides` package tests. The test should live in `handlers_test` instead, or we accept that PageFunc is set at handler init time and test it there.

Move the test to `handlers/guides_test.go`:
```go
func TestAllGuidesHavePageFunc(t *testing.T) {
	for _, g := range guides.All {
		if g.PageFunc == nil {
			t.Errorf("guide %q has nil PageFunc", g.Slug)
		}
	}
}
```

Run: `go test ./handlers/... -run TestAllGuidesHavePageFunc -v`
Expected: PASS.

**Step 5: Remove guideContent() function and simplify imports**

Delete the `guideContent` function (lines 416-444). Replace usage with `g.PageFunc(g, htmxRequest)` — handle nil PageFunc with the existing fallback.

**Step 6: Run full test suite + lint**

Run: `golangci-lint run ./... && go test ./... -v`
Expected: All pass.

**Step 7: Commit**

```bash
git add guides/registry.go handlers/guides.go handlers/guides_test.go
git commit -m "Add PageFunc to Guide struct, eliminate guideContent switch"
```

---

### Task 7: Make sidebar responsive with mobile hamburger menu

**Files:**
- Modify: `templates/sidebar.templ`
- Modify: `templates/layout.templ`

**Step 1: Add mobile hamburger overlay to layout.templ**

Wrap the sidebar + main in an Alpine x-data scope that manages `sidebarOpen` state. Add a hamburger button visible only on small screens (`md:hidden`), and make the sidebar hidden on small screens by default (`hidden md:flex`), shown as an overlay when `sidebarOpen` is true.

In `templates/layout.templ`, modify the `<body>`:
```templ
<body class="flex h-screen bg-gray-100 overflow-hidden" x-data="{ sidebarOpen: false }">
    <!-- Mobile hamburger -->
    <button
        class="fixed top-4 left-4 z-50 md:hidden bg-white border border-gray-200 rounded-lg p-2 shadow-sm"
        { templ.Attributes{"@click": "sidebarOpen = !sidebarOpen"}... }
    >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
        </svg>
    </button>
    <!-- Mobile overlay backdrop -->
    <div
        class="fixed inset-0 bg-black/50 z-30 md:hidden"
        x-show="sidebarOpen"
        { templ.Attributes{"@click": "sidebarOpen = false"}... }
    ></div>
    <!-- Font loader -->
    <span id="font-loader" hidden>
        if fontURL != "" {
            <link rel="preconnect" href="https://fonts.googleapis.com"/>
            <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
            <link rel="stylesheet" href={ fontURL }/>
        }
    </span>
    <!-- Sidebar -->
    @Sidebar(allGuides, activeSlug)
    <!-- Main content -->
    <main id="content" class="flex-1 overflow-y-auto">
        @content
    </main>
</body>
```

**Step 2: Make sidebar responsive in sidebar.templ**

Change the `<nav>` class from:
```
class="w-64 bg-white border-r border-gray-200 flex flex-col h-full shrink-0"
```
to:
```
class="fixed inset-y-0 left-0 z-40 w-64 bg-white border-r border-gray-200 flex flex-col h-full shrink-0 -translate-x-full md:translate-x-0 md:relative md:z-auto transition-transform duration-200"
```

Add Alpine bindings to show/hide on mobile:
```templ
:class="sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
```

Also close sidebar when a guide link is clicked on mobile — add `@click` handler to close: modify the existing @click to also set `sidebarOpen = false`.

**Step 3: Regenerate templ and test**

Run: `make templ && go test ./handlers/... -v`
Expected: All pass. Manually verify at different viewport widths.

**Step 4: Commit**

```bash
git add templates/sidebar.templ templates/sidebar_templ.go templates/layout.templ templates/layout_templ.go
git commit -m "Add responsive mobile sidebar with hamburger menu"
```

---

### Task 8: Add Cache-Control headers to static file handler

**Files:**
- Modify: `handlers/guides.go:71`

**Step 1: Wrap FileServer with cache middleware**

Replace line 71:
```go
mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
```
with:
```go
staticFS := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
mux.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "public, max-age=86400")
    staticFS.ServeHTTP(w, r)
}))
```

**Step 2: Add a test for Cache-Control header**

Add to `handlers/guides_test.go`:
```go
func TestStaticFilesCacheControl(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/static/css/output.css", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	cc := w.Header().Get("Cache-Control")
	if cc != "public, max-age=86400" {
		t.Errorf("expected Cache-Control header, got %q", cc)
	}
}
```

**Step 3: Run tests**

Run: `go test ./handlers/... -v`
Expected: All pass.

**Step 4: Commit**

```bash
git add handlers/guides.go handlers/guides_test.go
git commit -m "Add Cache-Control headers to static file serving"
```

---

### Task 9: Replace sort.Strings with slices.Sort in BuildCSSVars

**Files:**
- Modify: `guides/registry.go:470-475`

**Step 1: Update import and function**

Replace `"sort"` import with `"slices"`. Then replace:
```go
sort.Strings(keys)
```
with:
```go
slices.Sort(keys)
```

**Step 2: Run tests**

Run: `go test ./guides/... -v`
Expected: All pass.

**Step 3: Commit**

```bash
git add guides/registry.go
git commit -m "Use slices.Sort instead of sort.Strings"
```

---

### Task 10: Add safety comment for templ.Raw in sourceview.templ

**Files:**
- Modify: `templates/components/sourceview.templ`

**Step 1: Add comment documenting the safety invariant**

Read the file first, then add a comment above the `@templ.Raw(code)` line:
```templ
// code is pre-escaped HTML from chroma syntax highlighter (class-based formatter).
// Safe to use templ.Raw because chroma HTML-escapes all user content.
```

**Step 2: Commit**

```bash
git add templates/components/sourceview.templ
git commit -m "Document safety invariant for templ.Raw in source view"
```
