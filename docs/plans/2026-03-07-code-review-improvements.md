# Code Review Improvements Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Address all findings from the code review — eliminate duplication, fix bugs, improve security and maintainability, and move inline HTML to templ components.

**Architecture:** Shared utilities go into the `guides` package. Guide-specific handler endpoints get templ component counterparts. The handler file shrinks as rendering moves to templates. All changes are backward-compatible; no route changes.

**Tech Stack:** Go 1.25, Templ v0.3.1001, Tailwind CSS v4, HTMX 2.0.4, Alpine.js 3.14.9

---

### Task 0: Extract shared `BuildCSSVars()` into `guides` package

The identical `buildCSSVars()` function is duplicated across all 6 guide `styles.go` files. Extract it once into the `guides` package and update all callers.

**Files:**
- Modify: `guides/registry.go` (add `BuildCSSVars` function)
- Modify: `guides/brutalist/styles.go` (remove `buildCSSVars`)
- Modify: `guides/minimal/styles.go` (remove `buildCSSVars`)
- Modify: `guides/cassette/styles.go` (remove `buildCSSVars`)
- Modify: `guides/glass/styles.go` (remove `buildCSSVars`)
- Modify: `guides/bento/styles.go` (remove `buildCSSVars`)
- Modify: `guides/swiss/styles.go` (remove `buildCSSVars`)
- Modify: `guides/brutalist/brutalist.templ` (update call)
- Modify: `guides/minimal/minimal.templ` (update call)
- Modify: `guides/cassette/cassette.templ` (update call)
- Modify: `guides/glass/glass.templ` (update call)
- Modify: `guides/bento/bento.templ` (update call)
- Modify: `guides/swiss/swiss.templ` (update call)
- Test: `guides/registry_test.go`

**Step 1: Write the test for `BuildCSSVars`**

Add to `guides/registry_test.go`:

```go
func TestBuildCSSVars(t *testing.T) {
	vars := map[string]string{
		"--color-bg":   "#fff",
		"--color-text": "#000",
	}
	result := guides.BuildCSSVars(vars)
	if !strings.Contains(result, "--color-bg:#fff;") {
		t.Errorf("expected --color-bg:#fff; in result, got %q", result)
	}
	if !strings.Contains(result, "--color-text:#000;") {
		t.Errorf("expected --color-text:#000; in result, got %q", result)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./guides/ -run TestBuildCSSVars -v`
Expected: FAIL — `BuildCSSVars` not defined

**Step 3: Add `BuildCSSVars` to `guides/registry.go`**

Add at end of file:

```go
// BuildCSSVars generates ":root" CSS variable declarations from a map.
// Keys are sorted for deterministic output.
func BuildCSSVars(vars map[string]string) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(vars[k])
		sb.WriteString(";")
	}
	return sb.String()
}
```

Add `"sort"` and `"strings"` to the imports in `registry.go`.

**Step 4: Run test to verify it passes**

Run: `go test ./guides/ -run TestBuildCSSVars -v`
Expected: PASS

**Step 5: Remove `buildCSSVars` from all 6 guide `styles.go` files**

In each file, delete the entire `buildCSSVars` function. Also remove the `"strings"` import if `guideStyles()` doesn't use it (it won't — `guideStyles` returns a raw string literal). For `cassette/styles.go`, keep the `"strings"` import only if other code uses it (it doesn't — `strings` is only used by `buildCSSVars` there, but check if `guides` import covers it via the `guides.Highlight` call; `cassette/styles.go` already imports `guides` so removing `strings` is safe).

**Step 6: Update all 6 `.templ` files to call `guides.BuildCSSVars`**

In each templ file, find the line like:
```
@templ.Raw("<style>:root{" + buildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
```
Change `buildCSSVars(` to `guides.BuildCSSVars(`:
```
@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
```

The `guides` package is already imported in every templ file (for `guides.Guide` type).

**Step 7: Regenerate templ and verify build**

Run: `make templ && go build ./... && go test ./...`
Expected: All pass, no errors

**Step 8: Commit**

```bash
git add guides/registry.go guides/registry_test.go \
  guides/brutalist/styles.go guides/brutalist/brutalist.templ \
  guides/minimal/styles.go guides/minimal/minimal.templ \
  guides/cassette/styles.go guides/cassette/cassette.templ \
  guides/glass/styles.go guides/glass/glass.templ \
  guides/bento/styles.go guides/bento/bento.templ \
  guides/swiss/styles.go guides/swiss/swiss.templ \
  guides/brutalist/brutalist_templ.go guides/minimal/minimal_templ.go \
  guides/cassette/cassette_templ.go guides/glass/glass_templ.go \
  guides/bento/bento_templ.go guides/swiss/swiss_templ.go
git commit -m "refactor: extract shared BuildCSSVars into guides package

Removes 6 identical copies of buildCSSVars from guide styles.go files.
New version sorts keys for deterministic CSS output."
```

---

### Task 1: Fix Makefile `docker-build` / `docker-run` dependency inversion and remove stale `dev` phony

The `docker-build` target depends on `docker-run`, which is backwards. Also `dev` is declared in `.PHONY` but has no rule.

**Files:**
- Modify: `Makefile`

**Step 1: Fix Makefile**

Change line 1:
```makefile
.PHONY: build templ tailwind clean run docker-build docker-run
```

Change lines 32-37:
```makefile
# Build Docker image
docker-build:
	docker build -t $(IMAGE) .

# Run the Docker image locally (builds first if needed)
docker-run: docker-build
	docker run --rm -p 8080:8080 $(IMAGE)
```

**Step 2: Verify make targets parse correctly**

Run: `make -n build`
Expected: Shows `templ generate`, `tailwindcss ...`, `go build ...` — no errors

**Step 3: Commit**

```bash
git add Makefile
git commit -m "fix: correct docker-build/docker-run dependency order in Makefile

docker-run now depends on docker-build (was inverted).
Remove stale 'dev' from .PHONY since it has no rule."
```

---

### Task 2: Add SRI hash to Alpine.js CDN script tag

HTMX has an integrity hash but Alpine.js does not. Add one for consistency and supply-chain security.

**Files:**
- Modify: `templates/layout.templ:16`

**Step 1: Fetch the SRI hash for Alpine.js 3.14.9**

Run: `curl -s https://cdn.jsdelivr.net/npm/alpinejs@3.14.9/dist/cdn.min.js | openssl dgst -sha384 -binary | openssl base64 -A`

Use the resulting hash.

**Step 2: Update `templates/layout.templ` line 16**

Replace:
```html
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.14.9/dist/cdn.min.js" crossorigin="anonymous"></script>
```
With:
```html
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.14.9/dist/cdn.min.js" integrity="sha384-XXXXXX" crossorigin="anonymous"></script>
```
(substitute actual hash from step 1)

**Step 3: Regenerate templ and verify build**

Run: `make templ && go build ./...`
Expected: No errors

**Step 4: Commit**

```bash
git add templates/layout.templ templates/layout_templ.go
git commit -m "security: add SRI integrity hash to Alpine.js CDN script

Matches existing pattern used for HTMX script tag."
```

---

### Task 3: Mark `chroma/v2` as direct dependency in `go.mod`

`chroma/v2` is directly imported in `guides/highlight.go` but marked `// indirect` in `go.mod`.

**Files:**
- Modify: `go.mod`

**Step 1: Run `go mod tidy`**

Run: `go mod tidy`

This should automatically move `chroma/v2` from indirect to direct. Verify by checking `go.mod`:

Run: `grep chroma go.mod`
Expected: Line without `// indirect`

**Step 2: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: mark chroma/v2 as direct dependency in go.mod"
```

---

### Task 4: Add test verifying all registered guides have handler coverage

The `guideContent()` switch can silently fall through to placeholder for a new guide if someone forgets to add a case. Add a test that catches this.

**Files:**
- Modify: `handlers/guides_test.go`

**Step 1: Write the test**

Replace `TestAllGuidePagesOK` and add a comprehensive test that iterates `guides.All`:

```go
func TestAllRegisteredGuidesReturnOK(t *testing.T) {
	mux := handlers.NewMux()
	for _, g := range guides.All {
		t.Run(g.Slug, func(t *testing.T) {
			// Full page
			req := httptest.NewRequest(http.MethodGet, "/guides/"+g.Slug, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("GET /guides/%s: expected 200, got %d", g.Slug, w.Code)
			}
			if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
				t.Errorf("GET /guides/%s: expected text/html, got %q", g.Slug, ct)
			}
			// Verify guide renders its own content (not placeholder)
			body := w.Body.String()
			if strings.Contains(body, "placeholder") {
				t.Errorf("GET /guides/%s: appears to render placeholder instead of real content", g.Slug)
			}
		})
	}
}
```

Add `"github.com/johnfarrell/stylesheets/guides"` to the import block.

**Step 2: Remove the old `TestAllGuidePagesOK`**

Delete the old `TestAllGuidePagesOK` function (lines 89-105 of `handlers/guides_test.go`) since the new test supersedes it.

**Step 3: Run tests**

Run: `go test ./handlers/ -run TestAllRegisteredGuidesReturnOK -v`
Expected: PASS for all 6 guides

**Step 4: Commit**

```bash
git add handlers/guides_test.go
git commit -m "test: verify all registered guides have handler coverage

Replaces hardcoded slug list with dynamic iteration over guides.All.
Catches missing guideContent() switch cases by detecting placeholder fallback."
```

---

### Task 5: Add `PORT` env var support to `main.go`

The server address is hardcoded to `:8080`. Add an env var fallback for container deployments.

**Files:**
- Modify: `main.go`

**Step 1: Update `main.go`**

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/johnfarrell/stylesheets/handlers"
)

func main() {
	mux := handlers.NewMux()
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	addr = ":" + addr
	log.Printf("Starting server on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
```

**Step 2: Verify build**

Run: `go build -o /dev/null .`
Expected: No errors

**Step 3: Commit**

```bash
git add main.go
git commit -m "feat: support PORT env var for server address

Defaults to 8080 when unset. Useful for container deployments."
```

---

### Task 6: Add graceful shutdown to `main.go`

Replace bare `http.ListenAndServe` with `signal.NotifyContext` + `server.Shutdown(ctx)` for clean exits.

**Files:**
- Modify: `main.go`

**Step 1: Update `main.go`**

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/johnfarrell/stylesheets/handlers"
)

func main() {
	mux := handlers.NewMux()
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	addr = ":" + addr

	srv := &http.Server{Addr: addr, Handler: mux}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting server on http://localhost%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
```

**Step 2: Verify build**

Run: `go build -o /dev/null .`
Expected: No errors

**Step 3: Commit**

```bash
git add main.go
git commit -m "feat: add graceful shutdown with signal handling

Server cleanly shuts down on SIGINT/SIGTERM instead of hard exit."
```

---

### Task 7: Move demo form response HTML to templ components

`demoFormResponse()` in `handlers/guides.go` builds HTML via string concatenation. Move each variant into a small templ component.

**Files:**
- Create: `templates/components/formresponse.templ`
- Modify: `handlers/guides.go:244-274` (replace `demoFormResponse` with templ call)
- Test: `handlers/guides_test.go`

**Step 1: Write test for demo form endpoint**

Add to `handlers/guides_test.go`:

```go
func TestDemoFormPostReturnsHTML(t *testing.T) {
	mux := handlers.NewMux()
	for _, g := range guides.All {
		t.Run(g.Slug, func(t *testing.T) {
			body := strings.NewReader("name=TestUser")
			req := httptest.NewRequest(http.MethodPost, "/guides/"+g.Slug+"/demo-form", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("POST demo-form %s: expected 200, got %d", g.Slug, w.Code)
			}
			if !strings.Contains(w.Body.String(), "TestUser") {
				t.Errorf("POST demo-form %s: expected body to contain 'TestUser'", g.Slug)
			}
		})
	}
}
```

**Step 2: Run test to confirm it passes with current code**

Run: `go test ./handlers/ -run TestDemoFormPostReturnsHTML -v`
Expected: PASS (this is a characterization test for the existing behavior)

**Step 3: Create `templates/components/formresponse.templ`**

```
package components

// FormResponse renders the HTMX demo form success message styled per guide.
templ FormResponse(slug, name string) {
	switch slug {
	case "brutalist":
		<div class="border-2 border-black p-3 font-mono" style="background: var(--color-accent); box-shadow: var(--shadow-card);">
			<span class="font-bold uppercase">&#10003; Received:</span> <strong>{ name }</strong>
		</div>
	case "minimal":
		<div class="p-4" style="background: var(--color-surface); border: var(--border-width) solid var(--border-color); border-radius: var(--radius-lg); box-shadow: var(--shadow-card);">
			<p class="text-sm" style="color: var(--color-accent); font-weight: 500;">&#10003; Submitted successfully</p>
			<p class="text-sm mt-1" style="color: var(--color-secondary);">Thank you, <strong>{ name }</strong>.</p>
		</div>
	case "cassette":
		<div style="border: 1px solid var(--color-primary); padding: 0.75rem; font-family: var(--font-body); font-size: var(--font-size-caption);">
			<div style="color: var(--color-primary); font-weight: 700; margin-bottom: 0.25rem;">&#9654; TRANSMISSION RECEIVED</div>
			<div style="color: var(--color-text-muted);">OPERATOR: <strong style="color: var(--color-text);">{ name }</strong> — LOGGED</div>
		</div>
	case "glass":
		<div style="background: var(--frost-bg); backdrop-filter: blur(var(--frost-blur)); -webkit-backdrop-filter: blur(var(--frost-blur)); border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: 1rem;">
			<p class="text-sm font-semibold" style="color: var(--color-primary);">&#10003; Submitted</p>
			<p class="text-sm mt-1" style="color: var(--color-text-muted);">Thank you, <strong style="color: var(--color-text);">{ name }</strong>.</p>
		</div>
	case "bento":
		<div style="background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: 1rem; display: flex; align-items: flex-start; gap: 0.75rem;">
			<span style="color: var(--color-accent);">&#10003;</span>
			<div>
				<p class="text-sm font-medium" style="color: var(--color-text);">Submitted</p>
				<p class="text-xs mt-0.5" style="color: var(--color-text-muted);">Received from <strong>{ name }</strong></p>
			</div>
		</div>
	case "swiss":
		<div style="border-top: 3px solid var(--color-primary); padding: 1rem 0; margin-top: 1rem;">
			<p style="font-family: var(--font-body); font-size: 0.625rem; font-weight: 700; letter-spacing: 0.15em; text-transform: uppercase; color: var(--color-primary); margin-bottom: 0.25rem;">&#9654; RECEIVED</p>
			<p style="font-family: var(--font-display); font-weight: 700; color: var(--color-secondary);">{ name }</p>
		</div>
	default:
		<div class="p-3"><strong>Received:</strong> { name }</div>
	}
}
```

**Step 4: Update the handler to use the templ component**

In `handlers/guides.go`, replace the demo-form handler body (lines 79-85) with:

```go
name := r.FormValue("name")
if name == "" {
    name = "anonymous"
}
w.Header().Set("Content-Type", "text/html; charset=utf-8")
templ.Handler(components.FormResponse(slug, name)).ServeHTTP(w, r)
```

Add `"github.com/johnfarrell/stylesheets/templates/components"` to imports. Delete the `demoFormResponse()` function entirely.

**Step 5: Regenerate templ and run tests**

Run: `make templ && go build ./... && go test ./handlers/ -v`
Expected: All pass including `TestDemoFormPostReturnsHTML`

**Step 6: Commit**

```bash
git add templates/components/formresponse.templ templates/components/formresponse_templ.go handlers/guides.go handlers/guides_test.go
git commit -m "refactor: move demo form response HTML to templ component

Replaces string-concatenation HTML in demoFormResponse() with a
proper templ component. Automatic XSS escaping via templ's { name }."
```

---

### Task 8: Move bento metrics HTML to templ component

The `/guides/bento/metrics` endpoint builds HTML tiles via `fmt.Fprintf`.

**Files:**
- Create: `guides/bento/metrics.templ`
- Modify: `handlers/guides.go:88-110`

**Step 1: Create `guides/bento/metrics.templ`**

```
package bento

// MetricTile renders a single live metric card for the bento dashboard.
templ MetricTile(label, value, change, trend, trendColor string) {
	<div class="bento-card bento-span-6 flex flex-col gap-2">
		<p class="text-xs font-medium" style="color:var(--color-text-muted)">{ label }</p>
		<p class="text-2xl font-bold" style="color:var(--color-text)">{ value }</p>
		<p class="text-xs font-medium" style={ "color:" + trendColor }>{ change } { trend }</p>
	</div>
}
```

**Step 2: Update the handler**

Replace the inline `fmt.Fprintf` loop in the bento metrics handler with:

```go
mux.HandleFunc("/guides/bento/metrics", func(w http.ResponseWriter, r *http.Request) {
    type metric struct{ label, value, change, trend string }
    metrics := []metric{
        {"Active Users", fmt.Sprintf("%d", 1200+int(time.Now().Unix())%300), "+12%", "↑"},
        {"Revenue", fmt.Sprintf("$%.1fK", 48.2+float64(int(time.Now().Unix())%20)/10), "+8%", "↑"},
        {"Error Rate", fmt.Sprintf("%.1f%%", 0.3+float64(int(time.Now().Unix())%10)/10), "-0.1%", "↓"},
        {"Response Time", fmt.Sprintf("%dms", 120+int(time.Now().Unix())%80), "+5ms", "→"},
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    for _, m := range metrics {
        trendColor := "var(--color-accent)"
        if m.trend == "↓" {
            trendColor = "var(--color-danger)"
        }
        if m.trend == "→" {
            trendColor = "var(--color-text-muted)"
        }
        bentotempl.MetricTile(m.label, m.value, m.change, m.trend, trendColor).Render(r.Context(), w)
    }
})
```

**Step 3: Regenerate templ and run tests**

Run: `make templ && go build ./... && go test ./handlers/ -run TestBentoMetricsOK -v`
Expected: PASS

**Step 4: Commit**

```bash
git add guides/bento/metrics.templ guides/bento/metrics_templ.go handlers/guides.go
git commit -m "refactor: move bento metrics tiles to templ component"
```

---

### Task 9: Move glass edit-field HTML to templ component

The `/guides/glass/edit-field` endpoint builds 3 different HTML states via raw strings.

**Files:**
- Create: `guides/glass/editfield.templ`
- Modify: `handlers/guides.go:112-143`

**Step 1: Create `guides/glass/editfield.templ`**

```
package glass

// EditFieldDisplay renders the read-only view of the editable field.
templ EditFieldDisplay(name string) {
	<div class="flex items-center justify-between">
		<div>
			<p class="text-xs" style="color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.08em;">Project Name</p>
			<p class="text-lg font-semibold mt-1" style="color: var(--color-text);">{ name }</p>
		</div>
		<button class="glass-btn-ghost px-3 py-1.5 text-xs" hx-get="/guides/glass/edit-field" hx-target="#glass-editable" hx-swap="innerHTML">Edit</button>
	</div>
}

// EditFieldForm renders the inline edit form.
templ EditFieldForm(currentName string) {
	<form hx-post="/guides/glass/edit-field" hx-target="#glass-editable" hx-swap="innerHTML" class="flex items-end gap-3">
		<div class="flex-1">
			<p class="text-xs mb-1.5" style="color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.08em;">Project Name</p>
			<input type="text" name="name" value={ currentName } class="glass-input" autofocus/>
		</div>
		<button type="submit" class="glass-btn-primary px-4 py-2 text-xs">Save</button>
		<button type="button" class="glass-btn-ghost px-4 py-2 text-xs" hx-get="/guides/glass/edit-field?cancel=true" hx-target="#glass-editable" hx-swap="innerHTML">Cancel</button>
	</form>
}
```

**Step 2: Update the handler**

Replace the glass edit-field handler body with templ component calls:

```go
mux.HandleFunc("/guides/glass/edit-field", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if r.Method == http.MethodPost {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        name := r.FormValue("name")
        if name == "" {
            name = "Aurora Dashboard"
        }
        glasstempl.EditFieldDisplay(name).Render(r.Context(), w)
        return
    }
    if r.URL.Query().Get("cancel") == "true" {
        glasstempl.EditFieldDisplay("Aurora Dashboard").Render(r.Context(), w)
        return
    }
    glasstempl.EditFieldForm("Aurora Dashboard").Render(r.Context(), w)
})
```

**Step 3: Regenerate templ and run tests**

Run: `make templ && go build ./... && go test ./handlers/ -v`
Expected: All pass

**Step 4: Commit**

```bash
git add guides/glass/editfield.templ guides/glass/editfield_templ.go handlers/guides.go
git commit -m "refactor: move glass edit-field HTML to templ components"
```

---

### Task 10: Move minimal principles HTML to templ component

The `/guides/minimal/principles` endpoint builds HTML via raw string.

**Files:**
- Create: `guides/minimal/principles.templ`
- Modify: `handlers/guides.go:146-158`

**Step 1: Create `guides/minimal/principles.templ`**

```
package minimal

// Principles renders the lazy-loaded Design Principles content.
templ Principles() {
	<div class="space-y-6">
		<div>
			<h3 class="text-base font-semibold mb-2" style="color: var(--color-primary);">Reduction</h3>
			<p class="text-sm leading-relaxed" style="color: var(--color-secondary);">Remove until it breaks, then add one thing back. The last element you add is the design.</p>
		</div>
		<hr style="border: none; border-top: 1px solid var(--border-color);"/>
		<div>
			<h3 class="text-base font-semibold mb-2" style="color: var(--color-primary);">Whitespace</h3>
			<p class="text-sm leading-relaxed" style="color: var(--color-secondary);">Space is not emptiness — it is structure. Give every element room to breathe and it will speak more clearly.</p>
		</div>
		<hr style="border: none; border-top: 1px solid var(--border-color);"/>
		<div>
			<h3 class="text-base font-semibold mb-2" style="color: var(--color-primary);">Intention</h3>
			<p class="text-sm leading-relaxed" style="color: var(--color-secondary);">Every choice is deliberate. Color, weight, size, position — nothing is arbitrary. Minimal is not less; it is only what matters.</p>
		</div>
	</div>
}
```

**Step 2: Update the handler**

```go
mux.HandleFunc("/guides/minimal/principles", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    minimaltempl.Principles().Render(r.Context(), w)
})
```

**Step 3: Regenerate templ and run tests**

Run: `make templ && go build ./... && go test ./...`
Expected: All pass

**Step 4: Commit**

```bash
git add guides/minimal/principles.templ guides/minimal/principles_templ.go handlers/guides.go
git commit -m "refactor: move minimal principles HTML to templ component"
```

---

### Task 11: Move swiss search HTML to templ component

The `/guides/swiss/search` endpoint builds article cards via `fmt.Fprintf`.

**Files:**
- Create: `guides/swiss/search.templ`
- Modify: `handlers/guides.go:161-179`

**Step 1: Create `guides/swiss/search.templ`**

```
package swiss

// SearchResult renders a single editorial search result card.
templ SearchResult(eyebrow, headline, body string) {
	<div class="border-t-2 border-black pt-4 pb-8">
		<p class="swiss-label mb-3">{ eyebrow }</p>
		<h3 class="text-2xl font-bold mb-3" style="font-family: var(--font-display); color: var(--color-secondary);">{ headline }</h3>
		<p class="text-base leading-relaxed" style="color: var(--color-text); max-width: 55ch;">{ body }</p>
	</div>
}
```

**Step 2: Update the handler**

```go
mux.HandleFunc("/guides/swiss/search", func(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    type article struct{ eyebrow, headline, body string }
    articles := []article{
        {"Design Systems", "Grid as Foundation", "The grid is not a cage — it is a liberation from chaos."},
        {"Typography", "Weight Creates Hierarchy", "Bold speaks first. Regular speaks second. Light speaks last."},
        {"Color", "Red as Signal", "In Swiss design, red is never decoration. It is a signal."},
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    for _, a := range articles {
        if q != "" && !containsFold(a.eyebrow+a.headline+a.body, q) {
            continue
        }
        swisstempl.SearchResult(a.eyebrow, a.headline, a.body).Render(r.Context(), w)
    }
})
```

**Step 3: Regenerate templ and run tests**

Run: `make templ && go build ./... && go test ./...`
Expected: All pass

**Step 4: Commit**

```bash
git add guides/swiss/search.templ guides/swiss/search_templ.go handlers/guides.go
git commit -m "refactor: move swiss search result HTML to templ component"
```

---

### Task 12: Move cassette log HTML to templ component

The `/guides/cassette/log` endpoint builds a log entry via `fmt.Fprintf`.

**Files:**
- Create: `guides/cassette/logentry.templ`
- Modify: `handlers/guides.go:181-204`

**Step 1: Create `guides/cassette/logentry.templ`**

```
package cassette

// LogEntry renders a single system log line for the live feed.
templ LogEntry(timestamp, subsystem, message string) {
	<div class="flex gap-3" style="font-size:0.6875rem;color:var(--color-text-muted);padding:2px 0;border-bottom:1px solid var(--color-surface-2);font-family:var(--font-body)">
		<span style="min-width:4.5rem">[{ timestamp }]</span>
		<span style="color:var(--color-primary);font-weight:700;min-width:2.5rem">{ subsystem }</span>
		<span style="color:var(--color-text)">{ message }</span>
	</div>
}
```

**Step 2: Update the handler**

```go
mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
    entries := []struct{ sub, msg string }{
        {"SYS", "WCYPD COLONY SYSTEMS — HEARTBEAT NOMINAL"},
        {"NET", "NETWORK NODE 3 — PACKET LOSS 0.1% — WITHIN TOLERANCE"},
        {"ATM", "ATMOSPHERIC PROCESSOR — PRESSURE STABLE AT 101.3 kPa"},
        {"SEC", "MOTION SENSOR ARRAY — SECTOR 7G — NO CONTACTS"},
        {"PWR", "POWER GRID — OUTPUT 98.7% — NOMINAL"},
        {"NAV", "NAVIGATION ARRAY — COURSE HEADING CONFIRMED"},
        {"SCI", "SCIENCE LAB — ACCESS RESTRICTED — SPECIAL ORDER 937"},
        {"MED", "HYPERSLEEP UNITS — ALL OCCUPANT VITALS STABLE"},
        {"COM", "LONG-RANGE COMMS — SIGNAL RELAY B — ACTIVE"},
        {"ENG", "REACTOR COOLANT — TEMP 487°C — NOMINAL RANGE"},
        {"SEC", "BULKHEAD DOOR 14A — SEALED — VERIFIED"},
        {"SYS", "EMERGENCY LIGHTING — STANDBY MODE — READY"},
    }
    idx := int(time.Now().Unix()) % len(entries)
    e := entries[idx]
    ts := time.Now().Format("15:04:05")
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    cassettetempl.LogEntry(ts, e.sub, e.msg).Render(r.Context(), w)
})
```

**Step 3: Regenerate templ and run tests**

Run: `make templ && go build ./... && go test ./handlers/ -run TestCassetteLogOK -v`
Expected: PASS

**Step 4: Commit**

```bash
git add guides/cassette/logentry.templ guides/cassette/logentry_templ.go handlers/guides.go
git commit -m "refactor: move cassette log entry HTML to templ component"
```

---

### Task 13: Clean up handler — remove unused `demoFormResponse` helper and verify final state

After all templ migrations, clean up any remaining dead code and do a final verification pass.

**Files:**
- Modify: `handlers/guides.go` (remove `placeholderContent` if unused, remove `containsFold` if unused — note: `containsFold` is still used by swiss search handler, keep it)

**Step 1: Verify no dead code remains**

Run: `go vet ./... && go build ./...`
Check for unused imports or functions.

**Step 2: Run full test suite**

Run: `go test ./... -v`
Expected: All pass

**Step 3: Run tailwind build to verify CSS coverage**

Run: `make tailwind`
Expected: No errors (new templ files are picked up by `@source` directive)

**Step 4: Commit if any cleanup was needed**

```bash
git add handlers/guides.go
git commit -m "chore: remove dead code after templ migration"
```

---

## Task Dependency Summary

```
Task 0: Extract BuildCSSVars (independent)
Task 1: Fix Makefile (independent)
Task 2: Alpine SRI hash (independent)
Task 3: go mod tidy (independent)
Task 4: Test coverage (independent)
Task 5: PORT env var (independent)
Task 6: Graceful shutdown (blocked by Task 5)
Task 7: Demo form templ (independent)
Task 8: Bento metrics templ (independent)
Task 9: Glass edit-field templ (independent)
Task 10: Minimal principles templ (independent)
Task 11: Swiss search templ (independent)
Task 12: Cassette log templ (independent)
Task 13: Final cleanup (blocked by Tasks 7-12)
```

Tasks 0-5 and 7-12 are all independent and can be parallelized. Task 6 depends on Task 5. Task 13 depends on all templ migration tasks (7-12).
