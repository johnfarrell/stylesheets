# Templ Modularization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Extract two identical code patterns (OOB font loader, CSS injection) from all 10 guide templ files into shared components.

**Architecture:** Two new templ components in `templates/components/`. Each guide's `Page()` replaces inline boilerplate with a one-line component call. Guide-specific styling (swatches, forms, cards, buttons, typography, spacing) stays per-guide.

**Tech Stack:** Go 1.25, Templ v0.3.1001

---

### Task 1: Create OOBFontLoader component

**Files:**
- Create: `templates/components/fontloader.templ`

**Step 1: Create the component**

```templ
package components

// OOBFontLoader renders the OOB font swap for HTMX navigation.
// Call inside `if htmxRequest { }` in each guide's Page().
templ OOBFontLoader(fontURL string) {
	<span id="font-loader" { templ.Attributes{"hx-swap-oob": "true"}... }>
		<link rel="preconnect" href="https://fonts.googleapis.com"/>
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
		<link rel="stylesheet" href={ fontURL }/>
	</span>
}
```

**Step 2: Generate templ and build**

Run: `make templ && go build ./...`
Expected: PASS — component compiles with no errors

**Step 3: Commit**

```bash
git add templates/components/fontloader.templ templates/components/fontloader_templ.go
git commit -m "Add OOBFontLoader shared component"
```

---

### Task 2: Create GuideStyles component

**Files:**
- Create: `templates/components/guidestyles.templ`

**Step 1: Create the component**

```templ
package components

import "github.com/johnfarrell/stylesheets/guides"

// GuideStyles injects the CSS custom properties, dark mode overrides,
// and guide-specific styles into the page via a <style> block.
templ GuideStyles(g guides.Guide, extraCSS string) {
	@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guides.BuildDarkCSS(g.DarkCSSVars) + extraCSS + "</style>")
}
```

Note: `BuildDarkCSS(nil)` returns `""`, so guides without dark mode (Glass, Terminal, Retro) work without changes.

**Step 2: Generate templ and build**

Run: `make templ && go build ./...`
Expected: PASS

**Step 3: Commit**

```bash
git add templates/components/guidestyles.templ templates/components/guidestyles_templ.go
git commit -m "Add GuideStyles shared component"
```

---

### Task 3: Update all 10 guides to use new components

**Files to modify (templ files only — _templ.go files are regenerated):**
- `guides/brutalist/brutalist.templ`
- `guides/minimal/minimal.templ`
- `guides/cassette/cassette.templ`
- `guides/glass/glass.templ`
- `guides/bento/bento.templ`
- `guides/swiss/swiss.templ`
- `guides/terminal/terminal.templ`
- `guides/retro/retro.templ`
- `guides/newspaper/newspaper.templ`
- `guides/tracker/tracker.templ`

**Pattern A — OOB Font Loader replacement**

In each guide, replace the 7-line block:
```templ
if htmxRequest {
    <!-- OOB font loader: swaps the #font-loader span in <body> on HTMX navigation -->
    <span id="font-loader" { templ.Attributes{"hx-swap-oob": "true"}... }>
        <link rel="preconnect" href="https://fonts.googleapis.com"/>
        <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
        <link rel="stylesheet" href={ g.FontURL }/>
    </span>
}
```

With:
```templ
if htmxRequest {
    @components.OOBFontLoader(g.FontURL)
}
```

**Pattern B — CSS injection replacement**

Replace the CSS injection line. There are two variants:

Variant 1 (guides with dark mode — brutalist, cassette, minimal, bento, swiss, newspaper, tracker):
```templ
@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guides.BuildDarkCSS(g.DarkCSSVars) + guideStyles() + "</style>")
```

Variant 2 (guides without dark mode — glass, terminal, retro):
```templ
@templ.Raw("<style>:root{" + guides.BuildCSSVars(g.CSSVars) + "}" + guideStyles() + "</style>")
```

Both become:
```templ
@components.GuideStyles(g, guideStyles())
```

**Guide-specific notes:**
- **Cassette** and **Tracker**: CSS injection comes BEFORE the OOB font loader (opposite of other guides). Preserve this order.
- **Glass, Terminal, Retro**: Don't have `DarkCSSVars` — `BuildDarkCSS(nil)` returns `""` so `GuideStyles` handles this correctly.
- **Glass**: Uses `guides.` package import but not for `BuildDarkCSS` — after this change, the `guides` import may become unused in the templ file if `guideStyles()` is the only remaining reference. Check: `guideStyles()` is defined in the same package so it doesn't need a `guides.` import. However, `g guides.Guide` in the function signature keeps the import needed. No change required.

**Step 1: Update all 10 guide templ files**

Apply both patterns A and B to each file. For guides that don't import `components` yet (unlikely — all do), add the import.

**Step 2: Generate templ and build**

Run: `make templ && go build ./...`
Expected: PASS — all guides compile

**Step 3: Commit**

```bash
git add guides/*/
git commit -m "Replace inline OOB loader and CSS injection with shared components"
```

---

### Task 4: Build, test, and verify

**Step 1: Full build**

Run: `make build`
Expected: PASS

**Step 2: Lint**

Run: `go vet ./... && golangci-lint run`
Expected: PASS — no unused imports, no lint errors

**Step 3: Run server and spot-check**

Run: `make run`
Manually verify: navigate to each guide, confirm styles load, dark mode toggle works, HTMX navigation between guides swaps fonts correctly.

**Step 4: Final commit if any fixups needed**
