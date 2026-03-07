# Production Polish Design

Date: 2026-03-07

## Goal

Add the minimum production-ready surface to the site: a landing page, GitHub link, copyright notice, and a custom 404 page.

## Architecture

Four self-contained changes. No new routes, no new packages, no structural changes to the guide system.

1. **`/` route** — replace the auto-redirect with a real handler rendering a new `templates/home.templ` inside the existing `Layout`
2. **`templates/home.templ`** — new file: project tagline + responsive guide card grid
3. **`templates/sidebar.templ`** — add GitHub SVG link + copyright line to the existing footer `<div>`
4. **404 handler** — replace bare `http.NotFound` calls with a styled `templates/notfound.templ` rendered inside `Layout`

---

## Landing Page (`templates/home.templ`)

**Header:**
- Title: "Stylesheets"
- Tagline: "A reference collection of UI design languages built with Go, Templ, HTMX, Alpine.js, and Tailwind CSS."

**Guide cards:** Responsive 2–3 column grid. One card per guide showing:
- Guide name
- Guide description
- "View Guide" button — navigates via HTMX (`hx-get`, `hx-target="#content"`, `hx-push-url`) same as sidebar links

**Styling:** Neutral gray/white palette — consistent with sidebar shell, no guide-specific theming.

**Footer line:** Small GitHub link + copyright at the bottom of the content area.

---

## Sidebar Changes (`templates/sidebar.templ`)

The existing footer `<div>` gains two additions:

```
[GitHub SVG icon] GitHub     ← href: https://github.com/johnfarrell/stylesheets
Go · Templ · HTMX · Alpine · Tailwind   (existing line)
© 2026 John Farrell
```

All `text-xs` gray — understated, matches existing sidebar style.

---

## 404 Page (`templates/notfound.templ`)

Rendered inside `Layout` (sidebar remains visible and navigable). Content:

- Large "404" in gray
- "Page not found" subhead
- "The page you're looking for doesn't exist."
- "← Back to guides" link to `/`

Neutral styling — no guide color scheme.

The `/` handler, all `/guides/{slug}` handlers, and any other `http.NotFound` calls are replaced with a `renderNotFound` helper that renders this component.
