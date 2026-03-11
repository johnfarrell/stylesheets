# OSRS Goal Tracker Style Guide Design: Mission Control

**Goal:** A flagship style guide themed as a dark-mode goal and progress tracker for Old School RuneScape. "Mission control for your account." Demonstrates the full tech stack with a grounded use case — browsing, planning, and tracking progression across the entire OSRS ecosystem (skills, quests, diaries, bosses, collection log).

**Aesthetic:** Dark background, gold accent, retro mission-control flair with monospace labels and status indicators — but grounded in modern dark-mode tracker UX rather than full retro cosplay.

**Key Tech Showcases:**
- HTMX search filter (live search across game content)
- HTMX lazy-loaded detail panels (click item, server returns detail card)
- Alpine.js dependency graph (complex nested client-side state)
- Alpine.js sidebar tree (expand/collapse with hundreds of items)
- Alpine.js combat level calculator (reactive computed values)

---

## Identity

**Slug:** `tracker`
**Name:** Mission Control
**Fonts:** Space Mono (headers, labels, data readouts) + DM Sans (body text, descriptions)
**Prefix:** `trk-`

---

## Color Palette

| Name | Hex | CSS Var | Usage |
|---|---|---|---|
| Background | `#0d1117` | `--color-bg` | Page background |
| Surface | `#161b22` | `--color-surface` | Cards, panels |
| Surface 2 | `#1c2128` | `--color-surface-2` | Elevated surfaces, hover |
| Primary (Gold) | `#c8aa6e` | `--color-primary` | Accent, active states |
| Success (Green) | `#2ea043` | `--color-accent` | Completed, met requirements |
| Warning (Amber) | `#d29922` | `--color-warning` | In-progress, partial |
| Danger (Red) | `#da3633` | `--color-danger` | Locked, unmet |
| Info (Blue) | `#58a6ff` | `--color-info` | Links, informational |
| Text | `#e6edf3` | `--color-text` | Primary text |
| Text Muted | `#7d8590` | `--color-text-muted` | Secondary, labels |
| Border | `#30363d` | `--color-border` | Panel borders, dividers |

OSRS-inspired: gold is the dominant accent (gold text in-game), green/amber/red for status states.

---

## CSS Classes

### Panels & Cards
- `trk-panel` — Dark bordered card (surface bg, border stroke, subtle shadow)
- `trk-panel-header` — Uppercase monospace label bar with gold left-border accent
- `trk-panel-elevated` — Lighter surface for nested/hover states

### Status System
- `trk-status-light` — 8px dot indicator (base)
- `trk-status-complete` / `trk-status-progress` / `trk-status-locked` — Green/amber/red with subtle glow
- `trk-status-pulse` — Keyframe animation for in-progress (slow amber pulse)
- `trk-progress-bar` — Thin horizontal bar, gold fill on dark track

### Sidebar Tree
- `trk-tree` — Tree navigation container
- `trk-tree-node` — Single item with indent, monospace label, status light, hover highlight
- `trk-tree-node-active` — Gold left-border + elevated background
- `trk-tree-toggle` — Expand/collapse arrow, rotates on open

### Buttons & Inputs
- `trk-btn` — Bordered button, monospace text, gold border. Solid gold on hover
- `trk-btn-primary` — Solid gold background, dark text
- `trk-input` — Dark inset input, border-bottom highlight on focus
- `trk-search` — Search input styling

### Data Display
- `trk-readout` — Large monospace value display (XP, levels, completion %)
- `trk-tag` — Small pill label for categories (skill, quest, diary, etc.)
- `trk-dep-node` — Dependency graph node with status-colored left border
- `trk-dep-line` — Connecting line between dependency nodes (CSS borders/pseudo-elements)

### Utility
- `trk-rule` — Thin divider in border color
- `trk-glow` — Subtle gold text-shadow

---

## Sections

### Required Sections (1-6)

**1. Color Palette** [Alpine]
Copy-to-clipboard swatches. Dark swatches with hex values and usage labels. Gold border on hover.

**2. Typography** [None]
Space Mono at display/heading sizes (shown as "mission readout" headers). DM Sans at body/caption. Sample data panel mixing both fonts — monospace labels with sans-serif descriptions.

**3. Spacing Scale** [None]
Gold progress-bar-style bars at each spacing step. Monospace labels. Feels like a loading/XP bar.

**4. Buttons** [Alpine]
`trk-btn` variants (default, primary, danger, disabled). Sizes. Toggle demo: "Track" / "Untrack" button that swaps state and shows a status light going from off to green.

**5. Forms** [Both]
"Add Goal" form styled as a mission control input panel. Dark inset inputs, select dropdown for goal category (Skill, Quest, Diary, Boss, Collection). HTMX submit to `/guides/tracker/demo-form`. Response rendered as a confirmed-goal card.

**6. Cards/Panels** [Alpine]
Multiple `trk-panel` variants: basic data panel, panel with header bar, expandable panel, status summary panel with mini grid of status lights.

### Showcase Sections (7-9)

**7. Sidebar Tree + HTMX Search & Detail Loading** [Both]

Flagship section — a working content browser for OSRS progression.

**Left side: Sidebar tree** with top-level categories:
- Skills (all 23: Attack, Strength, Defence, Ranged, Prayer, Magic, Runecraft, Construction, Hitpoints, Agility, Herblore, Thieving, Crafting, Fletching, Slayer, Hunter, Mining, Smithing, Fishing, Cooking, Firemaking, Woodcutting, Farming)
- Quests (Free-to-play, Members — representative subset)
- Achievement Diaries (Ardougne, Desert, Falador, Fremennik, Kandarin, Karamja, Kourend & Kebos, Lumbridge & Draynor, Morytania, Varrock, Western Provinces, Wilderness)
- Bosses (representative subset)

Each node: status light (complete/in-progress/locked) + level/completion indicator.

**Top: Search bar** — `hx-get="/guides/tracker/search?q=..."` with `hx-trigger="input changed delay:300ms"`. Returns matching items. Replaces tree content temporarily.

**Right side: Detail panel** — Click any tree item fires `hx-get="/guides/tracker/detail/{category}/{id}"`. Detail shows:
- Name, category tag, status
- Requirements list with met/unmet indicators
- What this item unlocks (reverse dependencies)
- Progress bar if applicable

Alpine manages tree expand/collapse and active-node highlighting. HTMX handles search and detail loading.

Pre-populated with ~15-20 representative items across categories.

**8. Dependency Graph** [Alpine]

Visual prerequisite chain viewer. Pick a goal and see its full requirement tree as connected nodes.

**Layout:** Horizontal left-to-right flow. Each `trk-dep-node` shows:
- Item name, category tag, status light
- Level requirement if applicable
- Connected via `trk-dep-line` CSS borders

**Interactivity:**
- Alpine x-data holds graph data for 2-3 pre-built goal chains
- Buttons to switch between examples (Barrows Gloves, Quest Cape, Ardougne Elite)
- Click any node to highlight its own prerequisites (dims the rest)
- Skill nodes show level needed vs. "current" level

Fully client-side — no server round-trip.

**9. Account Overview Dashboard** [Alpine]

Summary dashboard — the "mission control" payoff.

- **Total Level readout** — Large monospace `trk-readout` (e.g. "1847 / 2277")
- **Quest Points** — "198 / 300" with progress bar
- **Completion grid** — 23 status lights for all skills, color-coded by bracket (1-49 red, 50-69 amber, 70-98 green, 99 gold glow)
- **Achievement summary row** — Diary completion counts (Easy/Medium/Hard/Elite) as mini progress bars
- **Combat Level calculator** — Alpine-computed combat level from editable skill inputs. Change a skill level, combat level updates live.

All Alpine-driven, dense data display.

---

## Server Endpoints

- `GET /guides/tracker/search?q=...` — Returns filtered tree nodes matching query
- `GET /guides/tracker/detail/{category}/{id}` — Returns detail card for selected item
- `POST /guides/tracker/demo-form` — Standard demo form (via shared FormResponse component)

---

## layout.templ Change

None — no new global scripts needed.

---

## Tech Stack Coverage

| Feature | Existing Guides | This Guide |
|---|---|---|
| HTMX form submit | All guides | Yes |
| HTMX search filter | Swiss | **Yes (game content search)** |
| HTMX lazy load | Retro OS | **Yes (detail panels)** |
| Alpine toggle/tabs | All | Yes |
| Alpine copy-to-clipboard | All | Yes |
| Alpine complex nested state | — | **Yes (dependency graph)** |
| Alpine computed values | — | **Yes (combat calculator)** |
| Alpine tree navigation | — | **Yes (sidebar tree)** |
| Dark mode design | — | **Yes (first dark guide)** |
