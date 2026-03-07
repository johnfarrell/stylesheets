package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/johnfarrell/stylesheets/handlers"
)

func TestIndexRedirects(t *testing.T) {
	mux := handlers.NewMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/guides/") {
		t.Errorf("expected redirect to /guides/*, got %q", loc)
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
}
