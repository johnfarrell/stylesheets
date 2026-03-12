package handlers

import (
	"bytes"
	"fmt"
	"log/slog"
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
	retrotempl "github.com/johnfarrell/stylesheets/guides/retro"
	swisstempl "github.com/johnfarrell/stylesheets/guides/swiss"
	terminaltempl "github.com/johnfarrell/stylesheets/guides/terminal"
	trackertempl "github.com/johnfarrell/stylesheets/guides/tracker"
	"github.com/johnfarrell/stylesheets/templates"
	"github.com/johnfarrell/stylesheets/templates/components"
)

func init() {
	for i := range guides.All {
		switch guides.All[i].Slug {
		case "brutalist":
			guides.All[i].PageFunc = brutalisttempl.Page
		case "cassette":
			guides.All[i].PageFunc = cassettetempl.Page
		case "minimal":
			guides.All[i].PageFunc = minimaltempl.Page
		case "glass":
			guides.All[i].PageFunc = glasstempl.Page
		case "bento":
			guides.All[i].PageFunc = bentotempl.Page
		case "swiss":
			guides.All[i].PageFunc = swisstempl.Page
		case "terminal":
			guides.All[i].PageFunc = terminaltempl.Page
		case "retro":
			guides.All[i].PageFunc = retrotempl.Page
		case "newspaper":
			guides.All[i].PageFunc = newspapertempl.Page
		case "tracker":
			guides.All[i].PageFunc = trackertempl.Page
		}
	}
}

// NewMux creates and returns the application HTTP mux with all routes registered.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	staticFS := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	mux.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		staticFS.ServeHTTP(w, r)
	}))

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
		var content templ.Component
		if guide.PageFunc != nil {
			content = guide.PageFunc(guide, false)
		} else {
			content = templ.Raw(fmt.Sprintf(`<div class="p-8"><h1 class="text-2xl font-bold">%s</h1><p class="text-gray-500 mt-2">%s</p></div>`,
				templ.EscapeString(guide.Name), templ.EscapeString(guide.Description)))
		}
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
		var partial templ.Component
		if guide.PageFunc != nil {
			partial = guide.PageFunc(guide, isHTMX)
		} else {
			partial = templ.Raw(fmt.Sprintf(`<div class="p-8"><h1 class="text-2xl font-bold">%s</h1><p class="text-gray-500 mt-2">%s</p></div>`,
				templ.EscapeString(guide.Name), templ.EscapeString(guide.Description)))
		}
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
		if err := components.FormResponse(slug, name).Render(r.Context(), w); err != nil {
			slog.Error("render failed", "error", err)
		}
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
			if err := bentotempl.MetricTile(m.label, m.value, m.change, m.trend, trendColor).Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
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
			if err := glasstempl.EditFieldDisplay(name).Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
			return
		}
		// GET with cancel — return the display view
		if r.URL.Query().Get("cancel") == "true" {
			if err := glasstempl.EditFieldDisplay("Aurora Dashboard").Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
			return
		}
		// GET — return the edit form
		if err := glasstempl.EditFieldForm("Aurora Dashboard").Render(r.Context(), w); err != nil {
			slog.Error("render failed", "error", err)
		}
	})

	// Minimal — lazy-loaded Design Principles content
	mux.HandleFunc("/guides/minimal/principles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := minimaltempl.Principles().Render(r.Context(), w); err != nil {
			slog.Error("render failed", "error", err)
		}
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
			if err := swisstempl.SearchResult(a.eyebrow, a.headline, a.body).Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
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
			if err := terminaltempl.BootMessage(bootTS, m.sub, m.msg, m.color).Render(r.Context(), &buf); err != nil {
				slog.Error("render boot message", "error", err)
				continue
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(buf.String(), "\n", "")); err != nil {
				slog.Error("write SSE event", "error", err)
				return
			}
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
		if err := terminaltempl.ExecResponse(cmd, output).Render(r.Context(), w); err != nil {
			slog.Error("render failed", "error", err)
		}
	})

	// Retro OS — lazy-load app window content
	mux.HandleFunc("/guides/retro/app/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch name {
		case "about":
			if err := retrotempl.AppAbout().Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
		case "calculator":
			if err := retrotempl.AppCalculator().Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
		case "files":
			if err := retrotempl.AppFiles().Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
		default:
			http.NotFound(w, r)
		}
	})

	// Newspaper — infinite scroll headlines
	mux.HandleFunc("/guides/newspaper/headlines", func(w http.ResponseWriter, r *http.Request) {
		allHeadlines := newspapertempl.Headlines
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
			if err := newspapertempl.HeadlineCard(strconv.Itoa(h.ID), h.Category, h.Title, h.Summary, h.Byline).Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
		}
		if end < len(allHeadlines) {
			if err := newspapertempl.HeadlineSentinel(strconv.Itoa(page+1)).Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
		}
	})

	// Newspaper — article view
	mux.HandleFunc("/guides/newspaper/article/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		a, ok := newspapertempl.Articles[id]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := newspapertempl.Article(a.Category, a.Title, a.Byline, a.Body).Render(r.Context(), w); err != nil {
			slog.Error("render failed", "error", err)
		}
	})

	// Newspaper — initial feed (back to front page)
	mux.HandleFunc("/guides/newspaper/feed", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/guides/newspaper/headlines?page=0", http.StatusSeeOther)
	})

	// Mission Control — search across tracker items
	mux.HandleFunc("/guides/tracker/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		for _, item := range trackertempl.Items {
			if q != "" && !containsFold(item.Name, q) {
				continue
			}
			if err := trackertempl.SearchResult(item.ID, item.Category, item.Name, item.Status).Render(r.Context(), w); err != nil {
				slog.Error("render failed", "error", err)
			}
		}
	})

	// Mission Control — detail panel for selected item
	mux.HandleFunc("/guides/tracker/detail/{category}/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		for _, item := range trackertempl.Items {
			if item.ID == id {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if err := trackertempl.Detail(item.Name, item.Category, item.Status, item.Description, item.Level, item.Target, item.Requirements, item.Unlocks).Render(r.Context(), w); err != nil {
					slog.Error("render failed", "error", err)
				}
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
		if err := cassettetempl.LogEntry(ts, e.sub, e.msg).Render(r.Context(), w); err != nil {
			slog.Error("render failed", "error", err)
		}
	})

	return mux
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
