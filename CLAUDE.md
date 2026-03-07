# Style Guide and Reference Project
---
This is a collection of style guides, references, and showcases for the following tech stack:
 - `Golang`: The backend language of choice
 - [Templ](https://github.com/a-h/templ): HTML templating library for Golang
 - [HTMX](https://htmx.org/): HTML tools for server interactions
 - [Alpine.js](https://alpinejs.dev/): Lightweight JavaScript framework for additional interactivity where HTMX can't work.
 - [Tailwind.css](https://tailwindcss.com/): Lightweight CSS framework providing many utility and styling classes

# Core Principles
---

 - This project is self-contained and showcases styles as a single page. Multiple styles can be viewed by navigating to other pages.
 - Rely on Templ/HTMX as much as possible for the foundation.
 - Use liberal use of Tailwind.css for styling
 - Google Fonts should be used for all fonts. Loaded via `<link>` in `<head>`.
 - All styles should be responsive.


# Style Guide Definition
---

A style guide demonstrates a visual design language -- Colors, typography, spacing, buttons, forms, cards, alerts, and other UI components --
all styled to a specific aesthetic.

## Required Sections

Every style guide should include at minimum:

1. Color Palette: Named swatches with hex values and usage philosophy.
2. Typography: Display, heading, body, and caption font samples with sizes/weights.
3. Spacing: Visual scale demonstration (base unit typically 4px or 8px).
4. Buttons: Primary, secondary, disabled states, multiple sizes.
5. Forms: Text inputs, selects, checkboxes, radio buttons with labels.
6. Cards/Panels: Content containers styled to the aesthetic.

Optional but common: Alerts, Navigation, Modals, Grid system, Design Principles section.

## Rules
1. Responsive: Include mobile breakpoints
2. Utilize the tech stack defined in thsi document.
3. Single Page: All components and style showcase should be viewable in one page.