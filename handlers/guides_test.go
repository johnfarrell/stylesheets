package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
	"github.com/johnfarrell/stylesheets/handlers"
)

func TestIndexRendersLandingPage(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html, got %q", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Brutalist") {
		t.Error("expected guide names in landing page body")
	}
}

func TestGuidePageOK(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/brutalist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html content type, got %q", ct)
	}
}

func TestGuideContentPartialOK(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/brutalist/content", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGuideNotFound(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/does-not-exist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html for 404 page, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "404") {
		t.Error("expected '404' in response body")
	}
}

func TestUnknownPathIs404(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/this-does-not-exist", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html for 404 page, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "404") {
		t.Error("expected styled 404 page body, got plain response")
	}
}

func TestAllRegisteredGuidesReturnOK(t *testing.T) {
	mux := handlers.NewMux()
	for _, g := range guides.All {
		t.Run(g.Slug, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/guides/"+g.Slug, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("GET /guides/%s: expected 200, got %d", g.Slug, w.Code)
			}
			if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
				t.Errorf("GET /guides/%s: expected text/html, got %q", g.Slug, ct)
			}
			body := w.Body.String()
			if strings.Contains(body, `class="text-gray-500 mt-2"`) {
				t.Errorf("GET /guides/%s: appears to render placeholder instead of real content", g.Slug)
			}
		})
	}
}

func TestBentoMetricsOK(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/bento/metrics", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html, got %q", ct)
	}
	if !strings.Contains(w.Body.String(), "bento-card") {
		t.Errorf("expected bento-card tiles in response body")
	}
}

func TestHTMXContentNotFoundRedirects(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/does-not-exist/content", nil)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if w.Header().Get("HX-Redirect") != "/" {
		t.Errorf("expected HX-Redirect: /, got %q", w.Header().Get("HX-Redirect"))
	}
}

func TestCassetteLogOK(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/guides/cassette/log", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html, got %q", ct)
	}
}
