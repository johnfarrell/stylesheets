# Project Overview

stylesheets is a collection and showcase for the specified tech stack to be used as a reference.

# Tech Stack

- Go 1.25
- No database
- net/http
- Templ
- HTMX
- Alpine.js
- Tailwind.css

# Architecture

```
main.go               # Entrypoint
guides/               # Individual style guides. Each guide has its own folder. Shared components are in the base path.
handlers/             # HTTP handlers
static/               # Static assets used
templates/            # Main location for Templ files
  components/         # Individual re-useable components across all style guides.
```

# Build & Run

```sh
make build          # Build
make run            # Run locally
make docker-build   # Build docker image
make docker-run     # Build and run docker image
make clean          # Clean egenerated files
golangci-lint run                 # Lint
```

# Project-Specific Rules

- This project is self-contained and showcases styles as a single page. Multiple styles can be viewed by navigating to other pages.
- Rely on Templ/HTMX as much as possible for the foundation.
- Use liberal use of Tailwind.css for styling
- Google Fonts should be used for all fonts. Loaded via `<link>` in `<head>`.
- All styles should be responsive.


## Style Guide Definition

A style guide demonstrates a visual design language -- Colors, typography, spacing, buttons, forms, cards, alerts, and other UI components --
all styled to a specific aesthetic.

### Required Sections

Every style guide should include at minimum:

1. Color Palette: Named swatches with hex values and usage philosophy.
2. Typography: Display, heading, body, and caption font samples with sizes/weights.
3. Spacing: Visual scale demonstration (base unit typically 4px or 8px).
4. Buttons: Primary, secondary, disabled states, multiple sizes.
5. Forms: Text inputs, selects, checkboxes, radio buttons with labels.
6. Cards/Panels: Content containers styled to the aesthetic.

Optional but common: Alerts, Navigation, Modals, Grid system, Design Principles section.

### Rules
1. Responsive: Include mobile breakpoints
2. Utilize the tech stack defined in thsi document.
3. Single Page: All components and style showcase should be viewable in one page.