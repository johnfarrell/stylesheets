package guides_test

import (
	"strings"
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
)

func TestParseSnippets_BasicExtraction(t *testing.T) {
	input := "before\n<!-- snippet:foo -->\n<div>hello</div>\n<!-- /snippet:foo -->\nafter"
	got := guides.ParseSnippets(input)
	want := "<div>hello</div>"
	if got["foo"] != want {
		t.Errorf("ParseSnippets[foo] = %q, want %q", got["foo"], want)
	}
}

func TestParseSnippets_MultipleRegions(t *testing.T) {
	input := `<!-- snippet:a -->
line-a
<!-- /snippet:a -->
middle
<!-- snippet:b -->
line-b
<!-- /snippet:b -->`
	got := guides.ParseSnippets(input)
	if got["a"] != "line-a" {
		t.Errorf("a = %q, want %q", got["a"], "line-a")
	}
	if got["b"] != "line-b" {
		t.Errorf("b = %q, want %q", got["b"], "line-b")
	}
}

func TestParseSnippets_GoStyleMarkers(t *testing.T) {
	input := "// snippet:handler\nfunc foo() {}\n// /snippet:handler"
	got := guides.ParseSnippets(input)
	if got["handler"] != "func foo() {}" {
		t.Errorf("handler = %q, want %q", got["handler"], "func foo() {}")
	}
}

func TestParseSnippets_MissingClose(t *testing.T) {
	input := "<!-- snippet:foo -->\n<div>hello</div>"
	got := guides.ParseSnippets(input)
	if _, ok := got["foo"]; ok {
		t.Error("expected no entry for unclosed snippet")
	}
}

func TestParseSnippets_EmptyInput(t *testing.T) {
	got := guides.ParseSnippets("")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestParseSnippets_NoSpaceMarkerVariant(t *testing.T) {
	input := "<!-- snippet:bar-->\n<p>content</p>\n<!-- /snippet:bar-->"
	got := guides.ParseSnippets(input)
	if got["bar"] != "<p>content</p>" {
		t.Errorf("bar = %q, want %q", got["bar"], "<p>content</p>")
	}
}

func TestGetSnippets_ReturnsEmptyMapForUnknownSlug(t *testing.T) {
	got := guides.GetSnippets("does-not-exist")
	if got == nil {
		t.Error("GetSnippets must return non-nil map")
	}
	if len(got) != 0 {
		t.Errorf("expected empty map for unknown slug, got %v", got)
	}
}

func TestGetHighlightedSnippets_ReturnsNonNilForUnknownSlug(t *testing.T) {
	got := guides.GetHighlightedSnippets("does-not-exist")
	if got == nil {
		t.Error("GetHighlightedSnippets must return non-nil map")
	}
}

func TestGetHighlightedSnippets_HighlightedHTMLContainsSpans(t *testing.T) {
	// Verify that Highlight applied to a real HTML snippet produces span tags.
	raw := guides.ParseSnippets("<!-- snippet:demo -->\n<div x-data=\"{}\">hello</div>\n<!-- /snippet:demo -->")
	code, ok := raw["demo"]
	if !ok {
		t.Fatal("test fixture snippet not parsed")
	}
	hl := guides.Highlight(code, guides.DetectLang(code))
	if !strings.Contains(hl, "<span") {
		t.Errorf("highlighted HTML snippet expected to contain <span, got: %s", hl)
	}
}
