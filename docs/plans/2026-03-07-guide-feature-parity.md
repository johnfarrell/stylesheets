# Guide Feature Parity & Showcase Improvements

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-extended-cc:executing-plans to implement this plan task-by-task.

**Goal:** Bring all 6 style guides to feature parity on required sections, add unique HTMX patterns beyond hx-post, add tabs to 2 more guides, enrich forms, and add alerts/notices to Glass and Bento.

**Architecture:** Each guide is a self-contained templ file + styles.go. Changes are per-guide: edit the .templ file to add sections, edit styles.go to add CSS classes. Handler routes added in handlers/guides.go only when new server endpoints are needed.

**Tech Stack:** Go + Templ + HTMX + Alpine.js + Tailwind CSS v4

---

### Task 1: Swiss — Add Spacing Scale, Button Toggle, Search Filter

**Files:**
- Modify: `guides/swiss/swiss.templ`

**What to add:**

1. **Spacing Scale section** — insert after Typography (section 2), before Grid System. Use `@swissSection("Spacing Scale", "")` wrapper. Display 4px base unit bars using `var(--color-primary)` background, ruled borders like other Swiss sections. Pattern: same loop as other guides `[]struct{ label, width string }`.

2. **Alpine toggle on Buttons** — add to existing Buttons section. Pattern: `x-data="{ active: false }"`, button swaps between `swiss-btn` and `swiss-btn-outline`, text shows "Active"/"Inactive".

3. **HTMX search filter on Editorial Cards** — unique HTMX pattern. Add `hx-get="/guides/swiss/search"` with `hx-trigger="input changed delay:300ms"` on a text input that filters the editorial cards. Requires new handler endpoint.

**Step 1:** Add spacing scale section to `swiss.templ` after the Typography section (after line ~94, before Grid System section):

```templ
<!-- 3. Spacing Scale -->
@swissSection("Spacing Scale", "") {
    <div class="space-y-3">
        <p class="swiss-label mb-4">Base unit: 8px (matches --baseline)</p>
        for _, s := range []struct{ label, width string }{
            {"8px", "8px"},
            {"16px", "16px"},
            {"24px", "24px"},
            {"32px", "32px"},
            {"48px", "48px"},
            {"64px", "64px"},
        } {
            <div class="flex items-center gap-6">
                <span class="swiss-label w-12 font-mono" style="letter-spacing: 0;">{ s.label }</span>
                <div style={ "width: " + s.width + "; height: 8px; background: var(--color-primary);" }></div>
            </div>
        }
    </div>
}
```

**Step 2:** Add Alpine toggle demo to Buttons section (after the existing button variants, before closing `}`):

```templ
<div class="mt-8">
    <p class="swiss-label mb-4">Toggle Demo [Alpine]</p>
    <div x-data="{ active: false }" class="flex items-center gap-6">
        <button
            class="text-sm"
            { templ.Attributes{
                ":class": "active ? 'swiss-btn' : 'swiss-btn-outline'",
                "@click": "active = !active",
            }... }
        >
            <span x-text="active ? 'ACTIVE' : 'INACTIVE'">INACTIVE</span>
        </button>
        <p class="swiss-label" x-text="active ? 'State: ON' : 'State: OFF'">State: OFF</p>
    </div>
</div>
```

**Step 3:** Add HTMX search filter section after Editorial Cards. Add search input with `hx-get` and a results container. Add handler in `handlers/guides.go`.

**Step 4:** Run `make templ && go build ./...` — verify compilation.

**Step 5:** Commit: `feat(swiss): add spacing scale, button toggle, and HTMX search filter`

---

### Task 2: Bento — Add Expandable Card, Alerts, Form Radios

**Files:**
- Modify: `guides/bento/bento.templ`
- Modify: `guides/bento/styles.go` (add alert CSS classes)

**What to add:**

1. **Cards & Panels section with expandable card** — new section between Status Indicators and Forms. Use `bento-card` class + Alpine `x-data="{ expanded: false }"`. Include a static card and an expandable card following the pattern in brutalist/minimal.

2. **Alert/Notice panels** — new section after Forms. Dashboard-style banners: success (green/accent), warning (amber), error (red/danger). Use existing color vars `--color-accent`, `--color-warning`, `--color-danger`.

3. **Radio buttons in form** — add radio group to existing Forms section.

**Step 1:** Add CSS classes for alerts and cards in `styles.go`:

```go
// Add to guideStyles() return string:
.bento-alert {
    border-radius: var(--radius-md);
    padding: 1rem 1.25rem;
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    font-size: 0.875rem;
}
```

**Step 2:** Add Cards & Panels section to `bento.templ` (after Status Indicators, before Forms):

```templ
@components.Section("Cards & Panels", components.BadgeAlpine) {
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Static card -->
        <div class="bento-card">
            <h3 class="text-sm font-semibold mb-2" style="color: var(--color-text);">Static Card</h3>
            <p class="text-sm" style="color: var(--color-text-muted);">A standard dashboard tile...</p>
            <div class="flex gap-2 mt-4" style="border-top: 1px solid var(--color-border); padding-top: 0.75rem;">
                <button class="bento-btn text-xs">Action</button>
                <button class="bento-btn bento-btn-secondary text-xs">Cancel</button>
            </div>
        </div>
        <!-- Expandable card -->
        <!-- snippet:card-expand -->
        <div class="bento-card" style="padding: 0;" x-data="{ expanded: false }">
            <div class="p-5 flex items-center justify-between cursor-pointer"
                { templ.Attributes{"@click": "expanded = !expanded"}... }>
                <h3 class="text-sm font-semibold" style="color: var(--color-text);">Expandable Panel</h3>
                <span class="text-lg transition-transform duration-200" style="color: var(--color-text-muted);"
                    { templ.Attributes{":class": "expanded ? 'rotate-45' : ''"}... }>+</span>
            </div>
            <div x-show="expanded" x-transition style="border-top: 1px solid var(--color-border);">
                <div class="p-5">
                    <p class="text-sm" style="color: var(--color-text-muted);">
                        Revealed content. Alpine handles toggle with zero server round-trips.
                    </p>
                </div>
            </div>
        </div>
        <!-- /snippet:card-expand -->
        @components.SourceView(snippets["card-expand"])
    </div>
}
```

**Step 3:** Add radio buttons to existing Forms section (after checkbox, before submit button).

**Step 4:** Add Alerts section after Forms:

```templ
@components.Section("Alerts & Notices", components.BadgeNone) {
    <div class="space-y-3">
        <div class="bento-alert" style="background: rgba(16,185,129,0.08); border: 1px solid rgba(16,185,129,0.2); color: var(--color-accent);">
            <span>✓</span><div><p class="font-medium">Success</p><p class="text-xs mt-0.5" style="color: var(--color-text-muted);">Operation completed successfully.</p></div>
        </div>
        <div class="bento-alert" style="background: rgba(245,158,11,0.08); border: 1px solid rgba(245,158,11,0.2); color: var(--color-warning);">
            <span>⚠</span><div><p class="font-medium">Warning</p><p class="text-xs mt-0.5" style="color: var(--color-text-muted);">Resource usage approaching limits.</p></div>
        </div>
        <div class="bento-alert" style="background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2); color: var(--color-danger);">
            <span>✕</span><div><p class="font-medium">Error</p><p class="text-xs mt-0.5" style="color: var(--color-text-muted);">Failed to connect to database.</p></div>
        </div>
    </div>
}
```

**Step 5:** Run `make templ && go build ./...` — verify compilation.

**Step 6:** Commit: `feat(bento): add expandable cards, alerts, and form radio buttons`

---

### Task 3: Brutalist — Add Tabs, HTMX Loading Indicator

**Files:**
- Modify: `guides/brutalist/brutalist.templ`
- Modify: `guides/brutalist/styles.go` (add tab + indicator CSS)

**What to add:**

1. **Tabs section** — Alpine-powered tab bar. Hard black borders, uppercase labels, stark active state (black bg, white text). 3 tabs with sample content.

2. **HTMX loading indicator** — Add `hx-indicator` to the existing form submit. Add a `.brut-spinner` CSS class. Shows a raw loading state during form submission (fits the exposed/functional aesthetic).

**Step 1:** Add CSS for tabs and spinner to `styles.go`:

```go
// Tabs
.brut-tab {
    padding: 0.5rem 1rem;
    font-family: var(--font-body);
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    border: 2px solid var(--color-primary);
    border-bottom: none;
    background: var(--color-bg);
    color: var(--color-primary);
    cursor: pointer;
}
.brut-tab-active {
    background: var(--color-primary);
    color: var(--color-bg);
}
// Loading indicator
.brut-indicator { display: none; }
.htmx-request .brut-indicator { display: inline; }
.htmx-request.brut-indicator { display: inline; }
```

**Step 2:** Add Tabs section to `brutalist.templ` after Cards & Panels.

**Step 3:** Add `hx-indicator` attribute to form submit button and add a spinner element.

**Step 4:** Run `make templ && go build ./...` — verify compilation.

**Step 5:** Commit: `feat(brutalist): add tabs component and HTMX loading indicator`

---

### Task 4: Minimal — Add Lazy-Loaded Section, Form Textarea

**Files:**
- Modify: `guides/minimal/minimal.templ`
- Modify: `handlers/guides.go` (add lazy-load endpoint)

**What to add:**

1. **Lazy-loaded "Design Principles" section** — uses `hx-get="/guides/minimal/principles"` with `hx-trigger="revealed"`. Content loads only when scrolled into view. Fits the progressive disclosure / spacious aesthetic. Requires new handler endpoint.

2. **Textarea in form** — add a `<textarea>` to the existing form for a "message" field.

**Step 1:** Add lazy-load section to `minimal.templ` after Cards & Panels:

```templ
@components.Section("Design Principles", components.BadgeHTMX) {
    <div
        hx-get="/guides/minimal/principles"
        hx-trigger="revealed"
        hx-swap="innerHTML"
        class="min-card p-8"
    >
        <p class="text-sm" style="color: var(--color-text-muted);">Loading...</p>
    </div>
}
```

**Step 2:** Add handler endpoint in `handlers/guides.go` that returns styled HTML content (3 design principle cards).

**Step 3:** Add textarea to form (before submit button).

**Step 4:** Run `make templ && go build ./...` — verify compilation.

**Step 5:** Commit: `feat(minimal): add lazy-loaded principles section and form textarea`

---

### Task 5: Glass — Add Tabs, Alerts, Form Radios + Textarea, Inline Edit

**Files:**
- Modify: `guides/glass/glass.templ`
- Modify: `guides/glass/styles.go` (add tab + alert CSS)
- Modify: `handlers/guides.go` (add inline-edit endpoint)

**What to add:**

1. **Frosted tabs** — pill-shaped tab bar with gradient active indicator. 3 tabs with sample content inside a glass-panel.

2. **Alert/Notice panel** — frosted alert with gradient left border. Info, success, and warning variants.

3. **Form enrichment** — add radio buttons and textarea to existing form.

4. **Inline edit** — HTMX-powered click-to-edit on a card field. `hx-get` loads edit form, `hx-post` saves. Unique HTMX pattern for this guide.

**Step 1:** Add CSS for tabs and alerts to `styles.go`:

```go
// Frosted tabs
.glass-tab {
    padding: 0.5rem 1rem;
    font-size: 0.8rem;
    font-weight: 500;
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    transition: background 0.2s, color 0.2s;
    border: none;
}
.glass-tab-active {
    background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
    color: #fff;
}
```

**Step 2:** Add Tabs section to `glass.templ` after Cards & Panels (before Modal).

**Step 3:** Add Alert section after Tabs.

**Step 4:** Add radio buttons and textarea to form.

**Step 5:** Add inline-edit section with HTMX endpoints. Add handlers to `handlers/guides.go`.

**Step 6:** Run `make templ && go build ./...` — verify compilation.

**Step 7:** Commit: `feat(glass): add tabs, alerts, form radios/textarea, and inline edit`

---

### Task 6: Final Verification

**Step 1:** Run `make build` — full rebuild (templ + tailwind + go build).

**Step 2:** Run `go test ./...` — ensure no regressions.

**Step 3:** Run `make run` and manually verify each guide loads and new sections render.

**Step 4:** Commit any fixes.
