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
		fmt.Fprint(w, demoFormResponse(slug, name))
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
			fmt.Fprintf(w,
				`<div class="bento-card bento-span-6 flex flex-col gap-2"><p class="text-xs font-medium" style="color:var(--color-text-muted)">%s</p><p class="text-2xl font-bold" style="color:var(--color-text)">%s</p><p class="text-xs font-medium" style="color:%s">%s %s</p></div>`,
				templ.EscapeString(m.label), templ.EscapeString(m.value), trendColor, templ.EscapeString(m.change), templ.EscapeString(m.trend),
			)
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
			fmt.Fprintf(w,
				`<div class="flex items-center justify-between"><div><p class="text-xs" style="color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.08em;">Project Name</p><p class="text-lg font-semibold mt-1" style="color: var(--color-text);">%s</p></div><button class="glass-btn-ghost px-3 py-1.5 text-xs" hx-get="/guides/glass/edit-field" hx-target="#glass-editable" hx-swap="innerHTML">Edit</button></div>`,
				templ.EscapeString(name))
			return
		}
		// GET with cancel — return the display view
		if r.URL.Query().Get("cancel") == "true" {
			fmt.Fprint(w,
				`<div class="flex items-center justify-between"><div><p class="text-xs" style="color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.08em;">Project Name</p><p class="text-lg font-semibold mt-1" style="color: var(--color-text);">Aurora Dashboard</p></div><button class="glass-btn-ghost px-3 py-1.5 text-xs" hx-get="/guides/glass/edit-field" hx-target="#glass-editable" hx-swap="innerHTML">Edit</button></div>`)
			return
		}
		// GET — return the edit form
		fmt.Fprint(w,
			`<form hx-post="/guides/glass/edit-field" hx-target="#glass-editable" hx-swap="innerHTML" class="flex items-end gap-3">`+
				`<div class="flex-1"><p class="text-xs mb-1.5" style="color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.08em;">Project Name</p>`+
				`<input type="text" name="name" value="Aurora Dashboard" class="glass-input" autofocus/></div>`+
				`<button type="submit" class="glass-btn-primary px-4 py-2 text-xs">Save</button>`+
				`<button type="button" class="glass-btn-ghost px-4 py-2 text-xs" hx-get="/guides/glass/edit-field?cancel=true" hx-target="#glass-editable" hx-swap="innerHTML">Cancel</button>`+
				`</form>`)
	})

	// Minimal — lazy-loaded Design Principles content
	mux.HandleFunc("/guides/minimal/principles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div class="space-y-6">`+
			`<div><h3 class="text-base font-semibold mb-2" style="color: var(--color-primary);">Reduction</h3>`+
			`<p class="text-sm leading-relaxed" style="color: var(--color-secondary);">Remove until it breaks, then add one thing back. The last element you add is the design.</p></div>`+
			`<hr style="border: none; border-top: 1px solid var(--border-color);"/>`+
			`<div><h3 class="text-base font-semibold mb-2" style="color: var(--color-primary);">Whitespace</h3>`+
			`<p class="text-sm leading-relaxed" style="color: var(--color-secondary);">Space is not emptiness — it is structure. Give every element room to breathe and it will speak more clearly.</p></div>`+
			`<hr style="border: none; border-top: 1px solid var(--border-color);"/>`+
			`<div><h3 class="text-base font-semibold mb-2" style="color: var(--color-primary);">Intention</h3>`+
			`<p class="text-sm leading-relaxed" style="color: var(--color-secondary);">Every choice is deliberate. Color, weight, size, position — nothing is arbitrary. Minimal is not less; it is only what matters.</p></div>`+
			`</div>`)
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

// demoFormResponse returns a guide-styled HTML response for the demo form.
func demoFormResponse(slug, name string) string {
	n := templ.EscapeString(name)
	switch slug {
	case "brutalist":
		return `<div class="border-2 border-black p-3 font-mono" style="background: var(--color-accent); box-shadow: var(--shadow-card);">` +
			`<span class="font-bold uppercase">&#10003; Received:</span> <strong>` + n + `</strong></div>`
	case "minimal":
		return `<div class="p-4" style="background: var(--color-surface); border: var(--border-width) solid var(--border-color); border-radius: var(--radius-lg); box-shadow: var(--shadow-card);">` +
			`<p class="text-sm" style="color: var(--color-accent); font-weight: 500;">&#10003; Submitted successfully</p>` +
			`<p class="text-sm mt-1" style="color: var(--color-secondary);">Thank you, <strong>` + n + `</strong>.</p></div>`
	case "cassette":
		return `<div style="border: 1px solid var(--color-primary); padding: 0.75rem; font-family: var(--font-body); font-size: var(--font-size-caption);">` +
			`<div style="color: var(--color-primary); font-weight: 700; margin-bottom: 0.25rem;">&#9654; TRANSMISSION RECEIVED</div>` +
			`<div style="color: var(--color-text-muted);">OPERATOR: <strong style="color: var(--color-text);">` + n + `</strong> — LOGGED</div></div>`
	case "glass":
		return `<div style="background: var(--frost-bg); backdrop-filter: blur(var(--frost-blur)); -webkit-backdrop-filter: blur(var(--frost-blur)); border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: 1rem;">` +
			`<p class="text-sm font-semibold" style="color: var(--color-primary);">&#10003; Submitted</p>` +
			`<p class="text-sm mt-1" style="color: var(--color-text-muted);">Thank you, <strong style="color: var(--color-text);">` + n + `</strong>.</p></div>`
	case "bento":
		return `<div style="background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: 1rem; display: flex; align-items: flex-start; gap: 0.75rem;">` +
			`<span style="color: var(--color-accent);">&#10003;</span>` +
			`<div><p class="text-sm font-medium" style="color: var(--color-text);">Submitted</p>` +
			`<p class="text-xs mt-0.5" style="color: var(--color-text-muted);">Received from <strong>` + n + `</strong></p></div></div>`
	case "swiss":
		return `<div style="border-top: 3px solid var(--color-primary); padding: 1rem 0; margin-top: 1rem;">` +
			`<p style="font-family: var(--font-body); font-size: 0.625rem; font-weight: 700; letter-spacing: 0.15em; text-transform: uppercase; color: var(--color-primary); margin-bottom: 0.25rem;">&#9654; RECEIVED</p>` +
			`<p style="font-family: var(--font-display); font-weight: 700; color: var(--color-secondary);">` + n + `</p></div>`
	default:
		return `<div class="p-3"><strong>Received:</strong> ` + n + `</div>`
	}
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
