package handlers

import (
	"bytes"
	"fmt"
	"net/http"
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
	"github.com/johnfarrell/stylesheets/templates"
	"github.com/johnfarrell/stylesheets/templates/components"
)

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

