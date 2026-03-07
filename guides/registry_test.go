package guides_test

import (
	"testing"

	"github.com/johnfarrell/stylesheets/guides"
)

func TestRegistryNotEmpty(t *testing.T) {
	if len(guides.All) == 0 {
		t.Fatal("guide registry must not be empty")
	}
}

func TestGuideBySlug(t *testing.T) {
	guide, ok := guides.BySlug("brutalist")
	if !ok {
		t.Fatal("expected to find 'brutalist' guide")
	}
	if guide.Slug != "brutalist" {
		t.Errorf("expected slug 'brutalist', got %q", guide.Slug)
	}
}

func TestGuideBySlugNotFound(t *testing.T) {
	_, ok := guides.BySlug("does-not-exist")
	if ok {
		t.Fatal("expected BySlug to return false for unknown slug")
	}
}

func TestGuideHasRequiredFields(t *testing.T) {
	for _, g := range guides.All {
		if g.Name == "" {
			t.Errorf("guide with slug %q has empty Name", g.Slug)
		}
		if g.Slug == "" {
			t.Error("guide has empty Slug")
		}
		if g.FontURL == "" {
			t.Errorf("guide %q has empty FontURL", g.Slug)
		}
	}
}
