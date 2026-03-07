package guides

import (
	"strings"
	"sync"
)

var (
	snippetCache map[string]map[string]string
	snippetOnce  sync.Once
)

// ParseSnippets extracts named regions from source text.
// Regions are delimited by:
//
//	HTML: <!-- snippet:name --> ... <!-- /snippet:name -->
//	Go:   // snippet:name ... // /snippet:name
//
// Leading/trailing whitespace is trimmed. Unclosed regions are silently ignored.
func ParseSnippets(src string) map[string]string {
	result := map[string]string{}
	lines := strings.Split(src, "\n")
	var current string
	var buf strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if name, ok := extractOpen(trimmed); ok {
			current = name
			buf.Reset()
			continue
		}

		if name, ok := extractClose(trimmed); ok {
			if current == name {
				result[name] = strings.TrimSpace(buf.String())
				current = ""
				buf.Reset()
			}
			continue
		}

		if current != "" {
			if buf.Len() > 0 {
				buf.WriteByte('\n')
			}
			buf.WriteString(line)
		}
	}

	return result
}

func extractOpen(line string) (string, bool) {
	if strings.HasPrefix(line, "<!-- snippet:") && strings.HasSuffix(line, "-->") {
		name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "<!-- snippet:"), "-->"))
		if name != "" {
			return name, true
		}
	}
	if strings.HasPrefix(line, "// snippet:") {
		name := strings.TrimSpace(strings.TrimPrefix(line, "// snippet:"))
		if name != "" {
			return name, true
		}
	}
	return "", false
}

func extractClose(line string) (string, bool) {
	if strings.HasPrefix(line, "<!-- /snippet:") && strings.HasSuffix(line, "-->") {
		name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "<!-- /snippet:"), "-->"))
		if name != "" {
			return name, true
		}
	}
	if strings.HasPrefix(line, "// /snippet:") {
		name := strings.TrimSpace(strings.TrimPrefix(line, "// /snippet:"))
		if name != "" {
			return name, true
		}
	}
	return "", false
}

func loadAll() map[string]map[string]string {
	cache := map[string]map[string]string{}
	for _, g := range All {
		path := g.Slug + "/" + g.Slug + ".templ"
		data, err := SourceFS.ReadFile(path)
		if err != nil {
			// Not yet embedded — skip silently
			continue
		}
		cache[g.Slug] = ParseSnippets(string(data))
	}
	return cache
}

// GetSnippets returns the parsed snippet map for a guide slug.
// Returns a non-nil empty map if the slug is unknown.
func GetSnippets(slug string) map[string]string {
	snippetOnce.Do(func() {
		snippetCache = loadAll()
	})
	if s, ok := snippetCache[slug]; ok {
		return s
	}
	return map[string]string{}
}

var (
	highlightedCache map[string]map[string]string
	highlightOnce    sync.Once
)

func loadHighlighted() map[string]map[string]string {
	snippetOnce.Do(func() { snippetCache = loadAll() })
	raw := snippetCache
	out := make(map[string]map[string]string, len(raw))
	for slug, snippets := range raw {
		out[slug] = make(map[string]string, len(snippets))
		for key, code := range snippets {
			out[slug][key] = Highlight(code, DetectLang(code))
		}
	}
	return out
}

// GetHighlightedSnippets returns syntax-highlighted HTML for each snippet of the named guide.
// Highlighting is computed once at startup and cached. Returns a non-nil map even if the slug is unknown.
func GetHighlightedSnippets(slug string) map[string]string {
	highlightOnce.Do(func() {
		highlightedCache = loadHighlighted()
	})
	if s, ok := highlightedCache[slug]; ok {
		return s
	}
	return map[string]string{}
}
