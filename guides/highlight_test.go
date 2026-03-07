package guides_test

import (
	"strings"
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
)

func TestHighlight_HTML(t *testing.T) {
	out := guides.Highlight(`<div class="foo">hello</div>`, "html")
	if !strings.Contains(out, "<span") {
		t.Errorf("expected span tags in output, got: %s", out)
	}
}

func TestHighlight_Go(t *testing.T) {
	out := guides.Highlight(`func main() {}`, "go")
	if !strings.Contains(out, "<span") {
		t.Errorf("expected span tags in output, got: %s", out)
	}
}

func TestHighlight_Empty(t *testing.T) {
	if out := guides.Highlight("", "go"); out != "" {
		t.Errorf("expected empty string for empty input, got: %q", out)
	}
}

func TestHighlight_NoPre(t *testing.T) {
	out := guides.Highlight(`func main() {}`, "go")
	if strings.Contains(out, "<pre") {
		t.Errorf("output must not contain a <pre> tag, got: %s", out)
	}
}

func TestDetectLang_HTML(t *testing.T) {
	cases := []string{
		`<div>hello</div>`,
		`  <span x-data="{}">`,
		`<!-- snippet:foo -->`,
	}
	for _, c := range cases {
		if lang := guides.DetectLang(c); lang != "html" {
			t.Errorf("DetectLang(%q) = %q, want html", c, lang)
		}
	}
}

func TestDetectLang_Go(t *testing.T) {
	cases := []string{
		`func main() {}`,
		`package guides`,
		`// handler comment`,
	}
	for _, c := range cases {
		if lang := guides.DetectLang(c); lang != "go" {
			t.Errorf("DetectLang(%q) = %q, want go", c, lang)
		}
	}
}
