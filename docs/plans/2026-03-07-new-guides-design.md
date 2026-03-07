# New Guides Design: Terminal, Retro OS, Newspaper

**Goal:** Add 3 new style guides that push deeper into HTMX, Alpine.js, and CSS capabilities while providing visually distinctive aesthetics not covered by existing guides.

**New HTMX features introduced:** SSE streaming, infinite scroll (`hx-trigger="revealed"`), view transitions (`transition:true`), optimistic UI patterns.

**New Alpine features introduced:** Keyboard navigation (`@keydown`), drag-and-drop (`@mousedown`/`@mousemove`/`@mouseup`), `Alpine.store` for cross-component state, typewriter animation, scroll tracking.

**Shared:** SSE extension script (`htmx-ext-sse`) loaded globally in `layout.templ`.

---

## Guide 1: Terminal

**Slug:** `terminal`
**Name:** Terminal
**Font:** Fira Code (Google Fonts)
**Palette:** `#0a0a0a` bg, `#00ff41` primary green, `#ff3333` danger, `#ffcc00` warning, `#00bfff` info. All monospace.
**Prefix:** `term-`

### Required Sections

1. **Color Palette** [Alpine] — Terminal escape code labels (`\033[32m`). Copy-to-clipboard.
2. **Typography** [None] — Fira Code weights in terminal-prompt style (`user@stylesheets:~$`). Typewriter animation on display text via Alpine `x-init` + `setInterval`.
3. **Spacing Scale** [None] — ASCII block character bars (`█`).
4. **Buttons** [Alpine] — Terminal command styling. Primary = `[EXECUTE]`, secondary = `[--dry-run]`, danger = `[sudo rm -rf]`. Toggle = process start/stop with blinking indicator.
5. **Forms** [Both] — Command prompt styling. Blinking cursor CSS animation. HTMX submit to `demo-form`.
6. **Cards/Panels** [Alpine] — Terminal windows with title bars (`── process.log ──`). Expandable.

### Showcase Sections

7. **Live System Boot** [HTMX SSE] — SSE endpoint `/guides/terminal/boot` streams ~15 boot messages with delays. `hx-ext="sse"`, `sse-connect`, `sse-swap`. Auto-scrolling append. "Reboot" button resets via `hx-get`.
8. **Command Prompt** [Both] — Input styled as `$_`. Enter: Alpine instantly shows command in history, HTMX posts to `/guides/terminal/exec` returning styled output. Hardcoded commands: `help`, `ls`, `whoami`, `date`.
9. **File Browser** [Alpine] — Arrow-key navigable file list. `@keydown.up/down` for navigation, `@keydown.enter` to open. Keyboard-driven interaction.

### CSS
- Scanline overlay via `::after` with repeating linear gradient
- `text-shadow: 0 0 5px currentColor` for CRT glow
- `@keyframes` cursor blink
- Dark bg, no border radius, monospace everything

### Server Endpoints
- `GET /guides/terminal/boot` — SSE stream of boot messages
- `POST /guides/terminal/exec` — Returns canned command output
- `POST /guides/terminal/demo-form` — Standard demo form (via shared endpoint)

---

## Guide 2: Retro OS

**Slug:** `retro`
**Name:** Retro OS
**Fonts:** VT323 (headings/pixel) + IBM Plex Sans (UI labels)
**Palette:** `#008080` teal desktop, `#c0c0c0` window chrome, `#000080` title bar blue, `#ffffff` window body, `#000000` text.
**Prefix:** `retro-`

### Required Sections

1. **Color Palette** [Alpine] — Inside a "Color Picker" window. Inset-border swatches. Copy-to-clipboard.
2. **Typography** [None] — Inside a "Notepad" window (`notepad.exe — Typography`). Full hierarchy.
3. **Spacing Scale** [None] — Pixel-art progress bars with beveled `outset` borders.
4. **Buttons** [Alpine] — Beveled `outset` borders, go `inset` on `:active`. Toggle = "Start Menu" button opens/closes menu panel.
5. **Forms** [Both] — Inside "System Properties" dialog. Inset inputs, square checkboxes. HTMX result appears in a "Message Box" dialog.
6. **Cards/Panels** [Alpine] — Multiple windows with [X] close buttons (`x-show`). One minimizable to title bar only.

### Showcase Sections

7. **Draggable Windows** [Alpine] — 2-3 windows draggable by title bar. `@mousedown`/`@mousemove`/`@mouseup`. Tracks `x, y` with `transform: translate()`. Click to bring to front (z-index via `Alpine.store('desktop', { topZ: 10 })`).
8. **Desktop Icons + HTMX Load** [Both] — Row of desktop icons. `@dblclick` opens a window, content lazily loaded via `hx-get="/guides/retro/app/{name}"`. Three mini-apps: "About", "Calculator" (Alpine), "File Manager" (list).
9. **Taskbar** [Alpine] — Bottom taskbar showing open windows as buttons. Click toggles window visibility. `Alpine.store` shares state between taskbar and windows.

### CSS
- `border-style: outset/inset` for 3D bevel
- `box-shadow: inset -1px -1px #0a0a0a, inset 1px 1px #ffffff` for window chrome
- Teal desktop background
- System font stack as fallback
- Window title bars with gradient backgrounds

### Server Endpoints
- `GET /guides/retro/app/about` — Static "About" window content
- `GET /guides/retro/app/calculator` — Calculator UI (Alpine-powered, content is static templ)
- `GET /guides/retro/app/files` — File manager list view
- `POST /guides/retro/demo-form` — Standard demo form (via shared endpoint)

---

## Guide 3: The Daily Style (Newspaper)

**Slug:** `newspaper`
**Name:** The Daily Style
**Fonts:** Playfair Display (headlines) + Source Serif 4 (body text)
**Palette:** `#faf9f6` cream paper, `#1a1a1a` ink black, `#c41e1e` spot red (sparingly). Very restrained.
**Prefix:** `news-`

### Required Sections

1. **Color Palette** [Alpine] — Minimal 3-4 color palette as ink swatches. Copy-to-clipboard.
2. **Typography** [None] — Newspaper front-page layout. Massive serif display, deck heads, body in 2-column CSS `column-count`. Drop cap via `::first-letter`.
3. **Spacing Scale** [None] — Column gutters and leading examples. Thin rule line dividers.
4. **Buttons** [Alpine] — Understated: text links, small bordered "Read More >" buttons. Toggle = "Breaking News" banner slide in/out.
5. **Forms** [Both] — "Letters to the Editor" submission. Bottom-border-only serif inputs. HTMX response in editorial voice.
6. **Cards/Panels** [Alpine] — Newspaper grid: one large feature spanning full width, smaller stories in 2-3 column grid. Expandable.

### Showcase Sections

7. **Breaking News Feed** [HTMX Infinite Scroll] — "Latest Headlines" section. `hx-trigger="revealed"` fires `hx-get="/guides/newspaper/headlines?page=N"`. Returns 3-4 article summaries per batch. `hx-swap="afterend"`. Cycles through ~15 canned headlines.
8. **Article View Transition** [HTMX] — Click headline loads full article with `hx-swap="innerHTML transition:true"`. "Back to Front Page" link returns to listing. Uses HTMX View Transitions API.
9. **Reading Progress** [Alpine] — Red progress bar at top. `@scroll.window` tracks percentage via `scrollY / (documentHeight - windowHeight)`.

### CSS
- `column-count: 2` (3 on wide screens) for newspaper text flow
- `::first-letter` drop cap: large, red, serif, floated
- `font-variant: small-caps` for bylines
- Thin 1px rule lines as section dividers
- Pull quotes with large left-border and italic serif
- `column-rule` for visible gutters

### Server Endpoints
- `GET /guides/newspaper/headlines?page=N` — Returns batch of 3-4 headline cards (infinite scroll)
- `GET /guides/newspaper/article/{id}` — Returns full article content for view transition
- `POST /guides/newspaper/demo-form` — Standard demo form (via shared endpoint)

---

## Tech Stack Coverage After Addition

| Feature | Current Guides | New Guides |
|---|---|---|
| HTMX form submit | Brutalist, Minimal, etc. | All three |
| HTMX polling | Cassette, Bento | — |
| HTMX search filter | Swiss | — |
| HTMX inline edit | Glass | — |
| **HTMX SSE** | — | **Terminal** |
| **HTMX infinite scroll** | — | **Newspaper** |
| **HTMX view transitions** | — | **Newspaper** |
| **HTMX lazy load** | — | **Retro OS** |
| Alpine toggle/tabs | All | All three |
| Alpine copy-to-clipboard | All | All three |
| **Alpine keyboard nav** | — | **Terminal** |
| **Alpine drag & drop** | — | **Retro OS** |
| **Alpine.store** | — | **Retro OS** |
| **Alpine scroll tracking** | — | **Newspaper** |
| **Alpine typewriter** | — | **Terminal** |

## layout.templ Change

Add SSE extension script globally (after HTMX script):
```html
<script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
```
