# Framework Design: Stylesheets Reference Application

Date: 2026-03-06

## Overview

A personal library of interactive style guides built with Go, Templ, HTMX, Alpine.js, and Tailwind CSS.
Each style guide demonstrates a distinct visual design language with fully interactive components,
clearly labeled by which technology powers each interaction. The app serves as a living reference
for use in future projects.

---

## Architecture: Flat Guide Registry (Approach A)

Each style guide is a Go struct in a central registry. The router auto-generates routes from the
registry. No database — all guides are defined in code.

### Project Structure

```
stylesheets/
├── main.go                  # Entry point, server setup, route registration
├── go.mod / go.sum
├── guides/
│   ├── registry.go          # Guide struct + slice of all registered guides
│   ├── brutalist/
│   │   └── brutalist.templ  # Component showcase for this guide
│   ├── material/
│   │   └── material.templ
│   └── ...                  # One directory per guide
├── templates/
│   ├── layout.templ         # Root HTML shell (head, sidebar, #content swap target)
│   ├── sidebar.templ        # Nav list rendered from registry
│   └── components/          # Optional shared wrappers (section headers, badges)
├── static/
│   └── css/
│       └── base.css         # Tailwind directives + global custom CSS
└── docs/
    └── plans/               # Design documents
```

---

## Routing & HTMX Navigation

Three route patterns:

```
GET /                          # Redirects to first guide in registry
GET /guides/{slug}             # Full page render (layout + sidebar + content)
GET /guides/{slug}/content     # HTMX partial — content area only
```

**Navigation flow:**
1. First visit → redirect to first registered guide
2. Full page renders `layout.templ` with sidebar + initial content
3. Sidebar links use `hx-get="/guides/{slug}/content"` + `hx-target="#content"` + `hx-push-url="/guides/{slug}"`
4. Only `#content` swaps; sidebar stays fixed; URL stays in sync
5. Each guide's CSS variables are injected as a `<style>` block inside the swapped content
6. The Google Font `<link>` is swapped via HTMX out-of-band (`hx-swap-oob`) targeting `#font-loader` in the layout

**Active state:** Alpine.js `x-data` on the sidebar nav tracks the active slug and applies highlight classes.

---

## Guide Registry

Each guide is defined as a Go struct:

```go
type Guide struct {
    Name        string
    Slug        string
    Description string
    FontURL     string
    CSSVars     map[string]string
}
```

The `CSSVars` map supports any CSS property that varies between guides — not just colors:

```go
Guide{
    Name:        "Brutalist",
    Slug:        "brutalist",
    Description: "Raw, functional, unapologetic design",
    FontURL:     "https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&display=swap",
    CSSVars: map[string]string{
        // Colors
        "--color-primary":    "#000000",
        "--color-secondary":  "#FF0000",
        "--color-bg":         "#FFFFFF",
        "--color-text":       "#000000",
        // Typography
        "--font-display":     "'Space Mono', monospace",
        "--font-body":        "'Space Mono', monospace",
        "--font-size-display":"4rem",
        // Shape
        "--radius-sm":        "0px",
        "--radius-md":        "0px",
        "--radius-lg":        "0px",
        // Elevation
        "--shadow-card":      "4px 4px 0px #000000",
        "--shadow-btn":       "2px 2px 0px #000000",
        // Borders
        "--border-width":     "2px",
        "--border-color":     "#000000",
        // Layout
        "--layout-columns":   "1",
        "--layout-gap":       "2rem",
        "--content-max-width":"1200px",
        "--section-padding":  "3rem 2rem",
    },
}
```

---

## Component Showcase Structure

Each guide's Templ file has **full control over its internal layout**. The shared `section.templ`
wrapper is optional — a guide can use it, extend it, or ignore it entirely.

**Required sections per guide:**
1. Color Palette — swatches with hex values; copy-to-clipboard via Alpine.js `[Alpine]`
2. Typography — live font samples at all weights/sizes
3. Spacing — visual scale demonstration
4. Buttons — primary/secondary/disabled/sizes with hover & active states
5. Forms — inputs, selects, checkboxes, radios; HTMX form submit with server response `[HTMX]`
6. Cards/Panels — content containers

**Optional sections:** Alerts, Navigation, Modals, Grid system, Design Principles

**Tech labeling convention:** Every interactive element displays a badge indicating which technology
drives it: `[HTMX]` for server-driven interactions, `[Alpine]` for client-side UI. This is a
core feature of the reference library — always explicit about what tech does what.

---

## Theming System

**Rule:** If a visual property changes between guides, it becomes a CSS variable. If it is
consistent across all guides (base layout grid, utility spacing), it stays as a Tailwind class.

**Custom CSS convention:** Any custom CSS must include a comment explaining why Tailwind cannot
handle it:
```css
/* [custom] - CSS var reference for per-guide theming, not possible with static Tailwind classes */
box-shadow: var(--shadow-card);
```

**Layout flexibility:** Guides are not constrained to a single layout. A minimalist guide may use
a single narrow column with large spacing; a retro-console guide may use a dense multi-column grid.
CSS vars like `--layout-columns`, `--layout-gap`, and `--content-max-width` provide starting points,
but guide Templ files can write any HTML structure that fits the aesthetic.

---

## Technology Responsibilities

| Technology  | Responsibility |
|-------------|----------------|
| Go          | HTTP server, routing, registry, template rendering |
| Templ       | All HTML — layout, sidebar, guide showcases |
| HTMX        | Content area swaps, form submissions, out-of-band font loading |
| Alpine.js   | Client-side UI state — active nav, copy-to-clipboard, dropdowns, modals, toggles |
| Tailwind    | Layout, spacing, base utilities — everything that doesn't vary per guide |
| CSS vars    | Per-guide theming — colors, fonts, radius, shadows, borders, layout tokens |
| Google Fonts| All display and body fonts, loaded via `<link>` in `<head>` |
