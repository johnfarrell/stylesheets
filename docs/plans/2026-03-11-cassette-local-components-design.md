# Cassette Guide — Local Component Extraction

**Date:** 2026-03-11
**Status:** Approved

## Goal

Extract repeated markup patterns from `guides/cassette/cassette.templ` into local templ
components within the `guides/cassette/` package. This reduces duplication, improves
readability, and showcases good templ component practices.

## New Files

All components are unexported (package-private to `cassette`).

| File | Components | Replaces |
|---|---|---|
| `swatches.templ` | `swatch(cssVar, description)`, `swatchGroup(title, cols)` | 10 copy-pasted swatch blocks in Section 1 |
| `typography.templ` | `typeSpecimenRow(fontName, sizeLabel, sampleText, style, last)` | 7 identical grid rows in Section 2 |
| `datarows.templ` | `dataRow(label, value)`, `dataRowColored(label, value, color)` | 4 identical for-loops in Section 6 |
| `notices.templ` | `dismissibleNotice(severity, message)` | 3 copy-pasted dismissible blocks in Section 9 |
| `gauges.templ` | `barGauge(label, alpineVar, fillClass, valueClass)` | 4 identical bar gauge blocks in Section 12 |
| `logentry.templ` | `logEntry(ts, code, msg)` | 11 hand-written static log lines in Section 14 |

## Component Signatures

### swatches.templ

```
templ swatch(cssVar, description string)
```
Renders a single clickable color swatch with Alpine `colorSwatch()` copy-to-clipboard.

```
templ swatchGroup(title, cols string)
```
Renders a labeled category header + grid container. Uses `{ children... }` slot for swatch content.

### typography.templ

```
templ typeSpecimenRow(fontName, sizeLabel, sampleText, style string, last bool)
```
Renders one row in the type specimen table. `last` controls whether border-bottom is rendered.

### datarows.templ

```
templ dataRow(label, value string)
templ dataRowColored(label, value, color string)
```
Renders label-value pairs with standard cassette styling. `dataRowColored` adds a custom
color to the value; empty color falls back to `var(--color-text)`.

### notices.templ

```
templ dismissibleNotice(severity, message string)
```
Renders a notice with Alpine dismiss button. Severity is one of: "note", "caution", "warning".

### gauges.templ

```
templ barGauge(label, alpineVar, fillClass, valueClass string)
```
Renders a labeled progress bar bound to an Alpine reactive variable.

### logentry.templ

```
templ logEntry(ts, code, msg string)
```
Renders a single static log line. The 11 static entries in Section 14 will be converted
from hand-written blocks to a `for` loop using this component.

## Additional Optimizations

During implementation, scan for any other repeated patterns not yet identified and
extract them using the same approach.

## Constraints

- All components stay within `guides/cassette/` (package `cassette`)
- No changes to shared components in `templates/components/`
- Existing visual output must be pixel-identical
- The `snippet` comment markers for `SourceView` must be preserved
