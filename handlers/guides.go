package handlers

import (
	"context"
	"fmt"
	"io"
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
	swisstempl "github.com/johnfarrell/stylesheets/guides/swiss"
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
			fmt.Fprintf(w,
				`<div class="border-t-2 border-black pt-4 pb-8"><p class="swiss-label mb-3">%s</p><h3 class="text-2xl font-bold mb-3" style="font-family: var(--font-display); color: var(--color-secondary);">%s</h3><p class="text-base leading-relaxed" style="color: var(--color-text); max-width: 55ch;">%s</p></div>`,
				templ.EscapeString(a.eyebrow), templ.EscapeString(a.headline), templ.EscapeString(a.body),
			)
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
		fmt.Fprintf(w,
			`<div class="flex gap-3" style="font-size:0.6875rem;color:var(--color-text-muted);padding:2px 0;border-bottom:1px solid var(--color-surface-2);font-family:var(--font-body)"><span style="min-width:4.5rem">[%s]</span><span style="color:var(--color-primary);font-weight:700;min-width:2.5rem">%s</span><span style="color:var(--color-text)">%s</span></div>`,
			ts, templ.EscapeString(e.sub), templ.EscapeString(e.msg),
		)
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
	default:
		return placeholderContent(g)
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

// placeholderContent renders a minimal placeholder until guide packages are implemented.
func placeholderContent(g guides.Guide) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<div class="p-8"><h1 class="text-2xl font-bold">%s</h1><p class="text-gray-500 mt-2">%s</p></div>`,
			templ.EscapeString(g.Name),
			templ.EscapeString(g.Description),
		)
		return err
	})
}
