package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/johnfarrell/stylesheets/guides"
	bentotempl "github.com/johnfarrell/stylesheets/guides/bento"
	brutalisttempl "github.com/johnfarrell/stylesheets/guides/brutalist"
	cassettetempl "github.com/johnfarrell/stylesheets/guides/cassette"
	glasstempl "github.com/johnfarrell/stylesheets/guides/glass"
	minimaltempl "github.com/johnfarrell/stylesheets/guides/minimal"
	newspapertempl "github.com/johnfarrell/stylesheets/guides/newspaper"
	trackertempl "github.com/johnfarrell/stylesheets/guides/tracker"
	retrotempl "github.com/johnfarrell/stylesheets/guides/retro"
	swisstempl "github.com/johnfarrell/stylesheets/guides/swiss"
	terminaltempl "github.com/johnfarrell/stylesheets/guides/terminal"
	"github.com/johnfarrell/stylesheets/templates"
	"github.com/johnfarrell/stylesheets/templates/components"
)

type trackerItem struct {
	ID           string
	Category     string // "skill", "quest", "diary", "boss"
	Name         string
	Status       string // "complete", "progress", "locked"
	Level        int    // current level or 0
	Target       int    // target level or 0
	Description  string
	Requirements []string
	Unlocks      []string
}

var trackerItems = []trackerItem{
	{ID: "attack", Category: "skill", Name: "Attack", Status: "progress", Level: 75, Target: 99, Description: "Determines accuracy with melee weapons.", Requirements: nil, Unlocks: []string{"Abyssal Whip at 70", "Dragon Claws at 60"}},
	{ID: "strength", Category: "skill", Name: "Strength", Status: "progress", Level: 82, Target: 99, Description: "Determines max hit with melee weapons.", Requirements: nil, Unlocks: []string{"Bandos Godsword spec at 70"}},
	{ID: "defence", Category: "skill", Name: "Defence", Status: "complete", Level: 70, Target: 70, Description: "Determines armour effectiveness.", Requirements: nil, Unlocks: []string{"Barrows equipment at 70"}},
	{ID: "ranged", Category: "skill", Name: "Ranged", Status: "progress", Level: 80, Target: 99, Description: "Determines accuracy and damage with ranged weapons.", Requirements: nil, Unlocks: []string{"Toxic Blowpipe at 75", "Armadyl Crossbow at 70"}},
	{ID: "prayer", Category: "skill", Name: "Prayer", Status: "progress", Level: 52, Target: 77, Description: "Unlocks combat prayers.", Requirements: nil, Unlocks: []string{"Protect from Melee at 43", "Rigour at 74", "Augury at 77"}},
	{ID: "magic", Category: "skill", Name: "Magic", Status: "progress", Level: 85, Target: 99, Description: "Determines accuracy with magic spells.", Requirements: nil, Unlocks: []string{"Ice Barrage at 94", "Trident of the Swamp at 75"}},
	{ID: "mining", Category: "skill", Name: "Mining", Status: "progress", Level: 72, Target: 85, Description: "Allows mining ore from rocks.", Requirements: nil, Unlocks: []string{"Amethyst at 92", "Runite at 85"}},
	{ID: "cooking", Category: "skill", Name: "Cooking", Status: "complete", Level: 70, Target: 70, Description: "Cook food to restore hitpoints.", Requirements: nil, Unlocks: []string{"Sharks at 80", "Anglerfish at 84"}},
	{ID: "cooks-assistant", Category: "quest", Name: "Cook's Assistant", Status: "complete", Level: 0, Target: 0, Description: "Help the Lumbridge cook gather ingredients for a cake.", Requirements: nil, Unlocks: []string{"300 Cooking XP", "Use of Cook-o-matic"}},
	{ID: "dragon-slayer", Category: "quest", Name: "Dragon Slayer", Status: "complete", Level: 0, Target: 0, Description: "Prove yourself a champion by slaying Elvarg.", Requirements: []string{"32 Quest Points", "8 Crafting", "34 Cooking (or good food)"}, Unlocks: []string{"Rune Platebody", "Access to Crandor"}},
	{ID: "monkey-madness", Category: "quest", Name: "Monkey Madness", Status: "progress", Level: 0, Target: 0, Description: "Rescue a squadron of Royal Guard from Ape Atoll.", Requirements: []string{"Tree Gnome Village", "Grand Tree", "35 Agility"}, Unlocks: []string{"Dragon Scimitar", "Access to Ape Atoll"}},
	{ID: "recipe-for-disaster", Category: "quest", Name: "Recipe for Disaster", Status: "locked", Level: 0, Target: 0, Description: "Save the Lumbridge Council from the Culinaromancer.", Requirements: []string{"175 Quest Points", "70 Cooking", "53 Thieving", "53 Fishing", "50 Mining", "50 Smithing"}, Unlocks: []string{"Barrows Gloves", "Culinaromancer's Chest"}},
	{ID: "desert-treasure", Category: "quest", Name: "Desert Treasure", Status: "locked", Level: 0, Target: 0, Description: "Uncover the secrets of the ancient element altars.", Requirements: []string{"50 Magic", "53 Thieving", "50 Firemaking", "10 Slayer", "The Digsite quest"}, Unlocks: []string{"Ancient Magicks", "Ice spells"}},
	{ID: "song-of-the-elves", Category: "quest", Name: "Song of the Elves", Status: "locked", Level: 0, Target: 0, Description: "Restore the city of Prifddinas to its former glory.", Requirements: []string{"70 Agility", "70 Construction", "70 Farming", "70 Herblore", "70 Hunter", "70 Mining", "70 Smithing", "70 Woodcutting"}, Unlocks: []string{"Access to Prifddinas", "Crystal equipment"}},
	{ID: "lumbridge-easy", Category: "diary", Name: "Lumbridge Easy Diary", Status: "complete", Level: 0, Target: 0, Description: "Complete easy tasks around Lumbridge and Draynor.", Requirements: []string{"15 Fishing", "15 Mining", "Rune Mysteries quest"}, Unlocks: []string{"Explorer's Ring 1", "Free cabbage teleports"}},
	{ID: "ardougne-medium", Category: "diary", Name: "Ardougne Medium Diary", Status: "progress", Level: 0, Target: 0, Description: "Complete medium tasks around Ardougne.", Requirements: []string{"51 Thieving", "35 Woodcutting", "Tribal Totem quest"}, Unlocks: []string{"Ardougne Cloak 2", "Improved pickpocketing"}},
	{ID: "karamja-elite", Category: "diary", Name: "Karamja Elite Diary", Status: "locked", Level: 0, Target: 0, Description: "Complete elite tasks on Karamja.", Requirements: []string{"91 Runecraft", "87 Herblore", "86 Smithing"}, Unlocks: []string{"Karamja Gloves 4", "Unlimited gem mine teleports"}},
	{ID: "giant-mole", Category: "boss", Name: "Giant Mole", Status: "complete", Level: 0, Target: 0, Description: "A large mole dwelling beneath Falador Park.", Requirements: []string{"Dharok's set recommended", "Falador Hard Diary (locator)"}, Unlocks: []string{"Mole parts (Clingy mole pet)", "Baby mole pet"}},
	{ID: "zulrah", Category: "boss", Name: "Zulrah", Status: "progress", Level: 0, Target: 0, Description: "A solo snake boss with multiple phases.", Requirements: []string{"Regicide quest", "75+ Magic/Ranged recommended"}, Unlocks: []string{"Tanzanite Fang → Blowpipe", "Magic Fang → Trident", "Serpentine Visage → Helm"}},
	{ID: "vorkath", Category: "boss", Name: "Vorkath", Status: "locked", Level: 0, Target: 0, Description: "An undead dragon boss on Ungael.", Requirements: []string{"Dragon Slayer II quest", "90+ Ranged recommended"}, Unlocks: []string{"Dragonbone Necklace", "Vorkath's Head → Assembler", "Vorki pet"}},
}

// NewMux creates and returns the application HTTP mux with all routes registered.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			renderNotFound(w, r)
			return
		}
		page := templates.Layout(guides.All, "", "", templates.Home(guides.All))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templ.Handler(page).ServeHTTP(w, r)
	})

	// Full page guide render
	mux.HandleFunc("/guides/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		guide, ok := guides.BySlug(slug)
		if !ok {
			renderNotFound(w, r)
			return
		}
		content := guideContent(guide, false)
		page := templates.Layout(guides.All, guide.Slug, guide.FontURL, content)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templ.Handler(page).ServeHTTP(w, r)
	})

	// HTMX partial content swap
	mux.HandleFunc("/guides/{slug}/content", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		guide, ok := guides.BySlug(slug)
		if !ok {
			renderNotFound(w, r)
			return
		}
		isHTMX := r.Header.Get("HX-Request") == "true"
		partial := guideContent(guide, isHTMX)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templ.Handler(partial).ServeHTTP(w, r)
	})

	// Demo form endpoint for showcasing HTMX form submission
	mux.HandleFunc("/guides/{slug}/demo-form", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		slug := r.PathValue("slug")
		name := r.FormValue("name")
		if name == "" {
			name = "anonymous"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		components.FormResponse(slug, name).Render(r.Context(), w)
	})

	// Bento Dashboard — live metric tiles (HTMX polling every 3s)
	mux.HandleFunc("/guides/bento/metrics", func(w http.ResponseWriter, r *http.Request) {
		type metric struct{ label, value, change, trend string }
		metrics := []metric{
			{"Active Users", fmt.Sprintf("%d", 1200+int(time.Now().Unix())%300), "+12%", "↑"},
			{"Revenue", fmt.Sprintf("$%.1fK", 48.2+float64(int(time.Now().Unix())%20)/10), "+8%", "↑"},
			{"Error Rate", fmt.Sprintf("%.1f%%", 0.3+float64(int(time.Now().Unix())%10)/10), "-0.1%", "↓"},
			{"Response Time", fmt.Sprintf("%dms", 120+int(time.Now().Unix())%80), "+5ms", "→"},
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for _, m := range metrics {
			trendColor := "var(--color-accent)"
			if m.trend == "↓" {
				trendColor = "var(--color-danger)"
			}
			if m.trend == "→" {
				trendColor = "var(--color-text-muted)"
			}
			bentotempl.MetricTile(m.label, m.value, m.change, m.trend, trendColor).Render(r.Context(), w)
		}
	})

	// Glass — inline edit: show edit form
	mux.HandleFunc("/guides/glass/edit-field", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			name := r.FormValue("name")
			if name == "" {
				name = "Aurora Dashboard"
			}
			glasstempl.EditFieldDisplay(name).Render(r.Context(), w)
			return
		}
		// GET with cancel — return the display view
		if r.URL.Query().Get("cancel") == "true" {
			glasstempl.EditFieldDisplay("Aurora Dashboard").Render(r.Context(), w)
			return
		}
		// GET — return the edit form
		glasstempl.EditFieldForm("Aurora Dashboard").Render(r.Context(), w)
	})

	// Minimal — lazy-loaded Design Principles content
	mux.HandleFunc("/guides/minimal/principles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		minimaltempl.Principles().Render(r.Context(), w)
	})

	// Swiss — HTMX search filter for editorial cards
	mux.HandleFunc("/guides/swiss/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		type article struct{ eyebrow, headline, body string }
		articles := []article{
			{"Design Systems", "Grid as Foundation", "The grid is not a cage — it is a liberation from chaos."},
			{"Typography", "Weight Creates Hierarchy", "Bold speaks first. Regular speaks second. Light speaks last."},
			{"Color", "Red as Signal", "In Swiss design, red is never decoration. It is a signal."},
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for _, a := range articles {
			if q != "" && !containsFold(a.eyebrow+a.headline+a.body, q) {
				continue
			}
			swisstempl.SearchResult(a.eyebrow, a.headline, a.body).Render(r.Context(), w)
		}
	})

	// Terminal — SSE boot sequence stream
	mux.HandleFunc("/guides/terminal/boot", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		bootMessages := []struct{ sub, msg, color string }{
			{"BIOS", "POST check... OK", "var(--color-text)"},
			{"BIOS", "Memory: 640K conventional, 3072K extended", "var(--color-text)"},
			{"BOOT", "Loading kernel...", "var(--color-primary)"},
			{"KERN", "Initializing subsystems", "var(--color-primary)"},
			{"NET ", "eth0: link up 1000Mbps", "var(--color-secondary)"},
			{"DISK", "Mounting /dev/sda1 on /", "var(--color-text)"},
			{"DISK", "Filesystem clean — 847392 blocks free", "var(--color-text)"},
			{"AUTH", "Loading user credentials", "var(--color-accent)"},
			{"PROC", "Starting daemon processes", "var(--color-text)"},
			{"PROC", "sshd: listening on port 22", "var(--color-primary)"},
			{"PROC", "httpd: listening on port 8080", "var(--color-primary)"},
			{"NET ", "Firewall rules loaded (47 rules)", "var(--color-secondary)"},
			{"SYS ", "System clock synchronized via NTP", "var(--color-text)"},
			{"SYS ", "All systems nominal", "var(--color-primary)"},
			{"BOOT", "READY. Type 'help' for commands.", "var(--color-primary)"},
		}

		ts := time.Now()
		for i, m := range bootMessages {
			select {
			case <-r.Context().Done():
				return
			default:
			}
			bootTS := ts.Add(time.Duration(i*200) * time.Millisecond).Format("15:04:05.000")
			var buf bytes.Buffer
			terminaltempl.BootMessage(bootTS, m.sub, m.msg, m.color).Render(r.Context(), &buf)
			fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(buf.String(), "\n", ""))
			flusher.Flush()
			time.Sleep(300 * time.Millisecond)
		}
	})

	// Terminal — command exec endpoint
	mux.HandleFunc("/guides/terminal/exec", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		cmd := r.FormValue("cmd")
		var output string
		switch strings.TrimSpace(strings.ToLower(cmd)) {
		case "help":
			output = "Available commands: help, ls, whoami, date, clear\nType any command and press Enter."
		case "ls":
			output = "README.md  main.go  go.mod  go.sum  handlers/  guides/  static/  templates/  Makefile"
		case "whoami":
			output = "guest@stylesheets"
		case "date":
			output = time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
		case "clear":
			output = ""
		default:
			output = "command not found: " + templ.EscapeString(cmd) + "\nType 'help' for available commands."
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		terminaltempl.ExecResponse(cmd, output).Render(r.Context(), w)
	})

	// Retro OS — lazy-load app window content
	mux.HandleFunc("/guides/retro/app/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch name {
		case "about":
			retrotempl.AppAbout().Render(r.Context(), w)
		case "calculator":
			retrotempl.AppCalculator().Render(r.Context(), w)
		case "files":
			retrotempl.AppFiles().Render(r.Context(), w)
		default:
			http.NotFound(w, r)
		}
	})

	// Newspaper — infinite scroll headlines
	mux.HandleFunc("/guides/newspaper/headlines", func(w http.ResponseWriter, r *http.Request) {
		type headline struct {
			id                                int
			category, title, summary, byline string
		}
		allHeadlines := []headline{
			{0, "Design", "The Grid Is Dead, Long Live the Grid", "Modern layout systems have made the rigid grid obsolete — or have they? A look at the evolution of page structure.", "By Jane Chen · 8 min read"},
			{1, "Typography", "Why Your Font Choice Is Wrong", "A provocative look at the assumptions designers make about typeface selection and readability.", "By Marcus Webb · 5 min read"},
			{2, "Color Theory", "The Case Against Color", "When restraint becomes the most powerful tool in a designer's arsenal.", "By Sarah Kim · 6 min read"},
			{3, "CSS", "Container Queries Changed Everything", "How the newest CSS specification is reshaping component-driven design.", "By Dev Patel · 7 min read"},
			{4, "Editorial", "Print Is Not Dead, It Evolved", "The newspaper aesthetic finds new life in digital interfaces.", "By The Editors · 4 min read"},
			{5, "Architecture", "Whitespace Is Not Empty Space", "Understanding the active role of negative space in visual hierarchy.", "By Yuki Tanaka · 5 min read"},
			{6, "Web Standards", "The Semantic Web We Were Promised", "Two decades later, are we any closer to the original vision?", "By Alex Rivera · 9 min read"},
			{7, "Design Systems", "One Component to Rule Them All", "The pursuit of the perfect reusable component — and why it's a trap.", "By Priya Sharma · 6 min read"},
			{8, "Typography", "The Golden Ratio Is Overrated", "Mathematical beauty does not always equal visual beauty.", "By Marcus Webb · 4 min read"},
			{9, "Accessibility", "Designing for Everyone Means Designing for No One", "A counterpoint to universal design — and why specificity matters.", "By Jordan Lee · 7 min read"},
			{10, "CSS", "Flexbox vs Grid: The Final Answer", "Spoiler: the answer is both. But knowing when to use which is the real skill.", "By Dev Patel · 5 min read"},
			{11, "Editorial", "The Attention Economy Broke Design", "How metrics-driven design is undermining craft.", "By The Editors · 3 min read"},
			{12, "Color Theory", "Red Means Stop (Except When It Doesn't)", "Cultural context and the unreliability of color as communication.", "By Sarah Kim · 6 min read"},
			{13, "Architecture", "Every Layout Is a Compromise", "The tensions between content, aesthetics, and engineering.", "By Yuki Tanaka · 8 min read"},
			{14, "Web Standards", "HTML Is a Programming Language", "A deliberately provocative position, rigorously defended.", "By Alex Rivera · 5 min read"},
		}
		pageStr := r.URL.Query().Get("page")
		page := 0
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			page = p
		}
		perPage := 3
		start := page * perPage
		if start >= len(allHeadlines) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			return
		}
		end := start + perPage
		if end > len(allHeadlines) {
			end = len(allHeadlines)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for _, h := range allHeadlines[start:end] {
			newspapertempl.HeadlineCard(strconv.Itoa(h.id), h.category, h.title, h.summary, h.byline).Render(r.Context(), w)
		}
		if end < len(allHeadlines) {
			newspapertempl.HeadlineSentinel(strconv.Itoa(page + 1)).Render(r.Context(), w)
		}
	})

	// Newspaper — article view
	mux.HandleFunc("/guides/newspaper/article/{id}", func(w http.ResponseWriter, r *http.Request) {
		type article struct {
			category, title, byline, body string
		}
		articles := map[string]article{
			"0":  {"Design", "The Grid Is Dead, Long Live the Grid", "By Jane Chen · March 7, 2026", "The grid has been the backbone of graphic design since the Bauhaus movement. For nearly a century, designers have relied on invisible lines to create order from chaos. But as digital interfaces have grown more fluid and responsive, the rigid grid has begun to feel like a constraint rather than a tool. Modern CSS layout systems — Flexbox, Grid, and now container queries — have given designers unprecedented freedom. Yet paradoxically, this freedom has led many back to the grid, not as a cage, but as a starting point. The best modern layouts use the grid as a foundation, then deliberately break it to create visual tension and hierarchy. The grid is dead. Long live the grid."},
			"1":  {"Typography", "Why Your Font Choice Is Wrong", "By Marcus Webb · March 6, 2026", "Every designer has a favorite typeface. For some it is Helvetica, that Swiss army knife of type. For others, it is something more expressive — a Didot, perhaps, or a carefully crafted variable font. But here is the uncomfortable truth: your font choice probably matters less than you think. Research consistently shows that readers adapt to virtually any well-set typeface within seconds. What matters far more is the typographic system — the relationships between sizes, weights, and spacing. A mediocre font set beautifully will always outperform a beautiful font set poorly. Stop agonizing over the typeface. Start obsessing over the system."},
			"2":  {"Color Theory", "The Case Against Color", "By Sarah Kim · March 5, 2026", "In a world of vibrant gradients and bold color palettes, there is something radical about restraint. The most powerful designs often use color sparingly — a single accent against a field of neutrals. This newspaper-inspired aesthetic proves the point: with just cream, black, and a touch of red, we can create hierarchy, emphasis, and emotional resonance. Color is not decoration. It is signal. And when everything is colorful, nothing stands out. The next time you reach for a rainbow palette, ask yourself: what if I used just one color instead?"},
			"3":  {"CSS", "Container Queries Changed Everything", "By Dev Patel · March 4, 2026", "For years, responsive design meant media queries — asking the viewport how wide it was, then making decisions based on that answer. But components do not live in viewports. They live in containers. A card might appear in a sidebar, a main column, or a modal, each with different available widths. Container queries finally let us ask the right question: how much space does my parent give me? This changes everything about how we think about component design. No more breakpoint gymnastics. No more wrapper divs to simulate container awareness. Just components that know their context and respond accordingly."},
			"4":  {"Editorial", "Print Is Not Dead, It Evolved", "By The Editors · March 3, 2026", "Every few years, someone declares print dead. And every few years, print proves them wrong — not by staying the same, but by evolving. The newspaper aesthetic you see on this page is not nostalgia. It is a recognition that centuries of typographic refinement produced principles that transcend medium. Column layouts, drop caps, pull quotes, careful leading — these are not print artifacts. They are solutions to the universal problem of making text readable and engaging. The web did not kill print. It gave print new life."},
		}
		id := r.PathValue("id")
		a, ok := articles[id]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		newspapertempl.Article(a.category, a.title, a.byline, a.body).Render(r.Context(), w)
	})

	// Newspaper — initial feed (back to front page)
	mux.HandleFunc("/guides/newspaper/feed", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/guides/newspaper/headlines?page=0", http.StatusSeeOther)
	})

	// Mission Control — search across OSRS items
	mux.HandleFunc("/guides/tracker/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for _, item := range trackerItems {
			if q != "" && !containsFold(item.Name, q) {
				continue
			}
			trackertempl.SearchResult(item.ID, item.Category, item.Name, item.Status).Render(r.Context(), w)
		}
	})

	// Mission Control — detail panel for selected item
	mux.HandleFunc("/guides/tracker/detail/{category}/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		for _, item := range trackerItems {
			if item.ID == id {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				trackertempl.Detail(item.Name, item.Category, item.Status, item.Description, item.Level, item.Target, item.Requirements, item.Unlocks).Render(r.Context(), w)
				return
			}
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/guides/cassette/log", func(w http.ResponseWriter, r *http.Request) {
		entries := []struct{ sub, msg string }{
			{"SYS", "WCYPD COLONY SYSTEMS — HEARTBEAT NOMINAL"},
			{"NET", "NETWORK NODE 3 — PACKET LOSS 0.1% — WITHIN TOLERANCE"},
			{"ATM", "ATMOSPHERIC PROCESSOR — PRESSURE STABLE AT 101.3 kPa"},
			{"SEC", "MOTION SENSOR ARRAY — SECTOR 7G — NO CONTACTS"},
			{"PWR", "POWER GRID — OUTPUT 98.7% — NOMINAL"},
			{"NAV", "NAVIGATION ARRAY — COURSE HEADING CONFIRMED"},
			{"SCI", "SCIENCE LAB — ACCESS RESTRICTED — SPECIAL ORDER 937"},
			{"MED", "HYPERSLEEP UNITS — ALL OCCUPANT VITALS STABLE"},
			{"COM", "LONG-RANGE COMMS — SIGNAL RELAY B — ACTIVE"},
			{"ENG", "REACTOR COOLANT — TEMP 487°C — NOMINAL RANGE"},
			{"SEC", "BULKHEAD DOOR 14A — SEALED — VERIFIED"},
			{"SYS", "EMERGENCY LIGHTING — STANDBY MODE — READY"},
		}
		idx := int(time.Now().Unix()) % len(entries)
		e := entries[idx]
		ts := time.Now().Format("15:04:05")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		cassettetempl.LogEntry(ts, e.sub, e.msg).Render(r.Context(), w)
	})

	return mux
}

// guideContent returns the Templ component for a guide's showcase.
// Add a case here when registering a new guide.
func guideContent(g guides.Guide, htmxRequest bool) templ.Component {
	switch g.Slug {
	case "brutalist":
		return brutalisttempl.Page(g, htmxRequest)
	case "cassette":
		return cassettetempl.Page(g, htmxRequest)
	case "minimal":
		return minimaltempl.Page(g, htmxRequest)
	case "glass":
		return glasstempl.Page(g, htmxRequest)
	case "bento":
		return bentotempl.Page(g, htmxRequest)
	case "swiss":
		return swisstempl.Page(g, htmxRequest)
	case "terminal":
		return terminaltempl.Page(g, htmxRequest)
	case "retro":
		return retrotempl.Page(g, htmxRequest)
	case "newspaper":
		return newspapertempl.Page(g, htmxRequest)
	case "tracker":
		return trackertempl.Page(g, htmxRequest)
	default:
		return templ.Raw(fmt.Sprintf(`<div class="p-8"><h1 class="text-2xl font-bold">%s</h1><p class="text-gray-500 mt-2">%s</p></div>`,
			templ.EscapeString(g.Name), templ.EscapeString(g.Description)))
	}
}

// renderNotFound serves a styled 404 page inside the main layout.
// For HTMX partial requests, it redirects the whole page to the landing page instead.
func renderNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	page := templates.Layout(guides.All, "", "", templates.NotFound())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templ.Handler(page, templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
}

// containsFold checks if s contains substr, case-insensitive.
func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

