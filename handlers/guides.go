package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	brutalisttempl "github.com/johnfarrell/stylesheets/guides/brutalist"
	"github.com/johnfarrell/stylesheets/guides"
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
		content := guideContent(guide)
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
		partial := guideContent(guide)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templ.Handler(partial).ServeHTTP(w, r)
	})

	// Demo form endpoint for showcasing HTMX form submission
	mux.HandleFunc("/guides/{slug}/demo-form", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		name := r.FormValue("name")
		if name == "" {
			name = "anonymous"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<div class="border-2 border-black p-3 bg-yellow-50 font-mono">✓ Received: <strong>%s</strong></div>`, templ.EscapeString(name))
	})

	return mux
}

// guideContent returns the Templ component for a guide's showcase.
// Add a case here when registering a new guide.
func guideContent(g guides.Guide) templ.Component {
	switch g.Slug {
	case "brutalist":
		return brutalisttempl.Page(g)
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
