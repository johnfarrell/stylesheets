package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/johnfarrell/stylesheets/guides"
	brutalisttempl "github.com/johnfarrell/stylesheets/guides/brutalist"
	cassettetempl "github.com/johnfarrell/stylesheets/guides/cassette"
	minimaltempl "github.com/johnfarrell/stylesheets/guides/minimal"
	"github.com/johnfarrell/stylesheets/templates"
)

// NewMux creates and returns the application HTTP mux with all routes registered.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Root redirect to first guide
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if len(guides.All) == 0 {
			http.Error(w, "no guides registered", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/guides/"+guides.All[0].Slug, http.StatusFound)
	})

	// Full page guide render
	mux.HandleFunc("/guides/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		guide, ok := guides.BySlug(slug)
		if !ok {
			http.NotFound(w, r)
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
			http.NotFound(w, r)
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
		name := r.FormValue("name")
		if name == "" {
			name = "anonymous"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<div class="border-2 border-black p-3 bg-yellow-50 font-mono">✓ Received: <strong>%s</strong></div>`, templ.EscapeString(name))
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
	default:
		return placeholderContent(g)
	}
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
