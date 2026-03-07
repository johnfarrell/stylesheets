# Cassette Futurism Guide Design

Date: 2026-03-06

## Aesthetic

Cassette futurism interpreted as NASA technical documentation — not dark CRT terminals,
but the dense, ruled, deliberate design of Apollo-era mission documents, Weyland-Yutani
corporate technical manuals, and golden-age-of-computing instrument logbooks.

Think: fillable technical forms, numbered sections, blue demarcation rules, IBM Selectric
typewriter output, analog instrument readouts, mission status boards.

---

## Typography

- **Display header (guide masthead only):** Orbitron 700 — the single dramatic scifi moment
- **Everything else:** IBM Plex Mono 300/400/500/700 — closest available Google Font to Berkeley Mono

Font URL: `https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:ital,wght@0,300;0,400;0,500;0,700;1,400&family=Orbitron:wght@400;700&display=swap`

---

## Color System

| Token | Value | Usage |
|---|---|---|
| `--color-bg` | `#f5f4ef` | Warm off-white — document paper |
| `--color-surface` | `#ffffff` | Pure white — data panels |
| `--color-surface-2` | `#e8e7e2` | Light warm gray — panel headers |
| `--color-primary` | `#0b3d91` | NASA deep blue — accent, rules, labels |
| `--color-secondary` | `#1a5276` | Darker blue — secondary elements |
| `--color-danger` | `#c0392b` | Technical red — warnings, critical |
| `--color-caution` | `#c85200` | Orange — caution notices |
| `--color-text` | `#1a1a14` | Near-black — all body text |
| `--color-text-muted` | `#5a5a52` | Medium gray — secondary labels |
| `--color-border` | `#c8c7c0` | Light warm gray — standard borders |
| `--color-rule` | `#0b3d91` | Blue — section demarcation rules |

---

## Layout

- Multi-column (2-col for dense sections, full-width for showcases)
- Numbered sections: `1.0`, `1.1`, etc.
- Heavy horizontal blue rules (`3px solid --color-primary`) between major sections
- Thin gray rules within sections
- Max-width: 1200px
- Section padding: deliberate, not tight but not spacious — technical manual rhythm

---

## Component Sections (15 total)

### Required (6)
1. **Color Palette** [Alpine] — swatches with hex values, copy-to-clipboard
2. **Typography** — IBM Plex Mono specimens across all weights/sizes; Orbitron display
3. **Spacing Scale** — visual ruler with labeled tick marks
4. **Buttons** [Alpine] — technical styles, labeled states, toggle demo
5. **Forms** [HTMX + Alpine] — fillable-form style fields, HTMX submit
6. **Cards/Panels** [Alpine] — document section panels, collapsible data blocks

### Extended (9 more — verbose per user request)
7. **Data Tables** — the defining element of technical documents; sortable headers [Alpine]
8. **Status Board** [Alpine] — system status grid with live indicators, green/amber/red lights
9. **Notices** — WARNING / CAUTION / NOTE blocks (NASA tech manual style)
10. **Document Header Block** — mission doc style metadata header (doc number, revision, date, classification)
11. **Navigation Tabs** [Alpine] — document section tabs, active state
12. **Technical Readouts** [Alpine] — live-updating numerical displays with setInterval
13. **Progress Trackers** [Alpine] — segmented completion bars, percentage displays
14. **System Log** [HTMX] — HTMX-polled scrolling log output every 4s
15. **Modal/Dialog** [Alpine] — technical confirmation dialog with authorization fields

---

## CSS Conventions

- Blue top-border rules (`border-top: 3px solid var(--color-primary)`) for major section breaks
- Corner labels (`position: absolute; top: 0; left: 0`) for panel classifications
- `font-family: var(--font-body)` = IBM Plex Mono everywhere
- Custom CSS only where Tailwind can't reference CSS vars; all marked `/* [custom] - reason */`
- Class prefix: `cass-` to avoid collisions with other guides

---

## New Handler Routes Required

- `GET /guides/cassette/log` — returns one random log entry (HTML fragment) for HTMX polling
- `POST /guides/cassette/demo-form` — handled by existing generic demo-form route
