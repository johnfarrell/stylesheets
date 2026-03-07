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
		name := strings.TrimSuffix(strings.TrimPrefix(line, "<!-- snippet:"), " -->")
		name = strings.TrimSpace(name)
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
		name := strings.TrimSuffix(strings.TrimPrefix(line, "<!-- /snippet:"), " -->")
		name = strings.TrimSpace(name)
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
	files := map[string]string{
		"brutalist": "brutalist/brutalist.templ",
		"minimal":   "minimal/minimal.templ",
		"cassette":  "cassette/cassette.templ",
	}
	for slug, path := range files {
		data, err := SourceFS.ReadFile(path)
		if err != nil {
			continue
		}
		cache[slug] = ParseSnippets(string(data))
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
