package guides

// Guide defines a style guide's metadata and theme tokens.
type Guide struct {
	Name        string
	Slug        string
	Description string
	FontURL     string
	// CSSVars holds all per-guide CSS custom properties.
	// Any visual property that differs between guides belongs here:
	// colors, typography, radius, shadows, borders, layout tokens, etc.
	CSSVars map[string]string
}

// All is the ordered list of registered style guides.
// Add new guides here to register them with the application.
var All = []Guide{
	{
		Name:        "Brutalist",
		Slug:        "brutalist",
		Description: "Raw, functional, unapologetic design with heavy borders and stark contrast.",
		FontURL:     "https://fonts.googleapis.com/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400&display=swap",
		CSSVars: map[string]string{
			// Colors
			"--color-primary":    "#000000",
			"--color-secondary":  "#FF0000",
			"--color-accent":     "#FFFF00",
			"--color-bg":         "#FFFFFF",
			"--color-surface":    "#F5F5F5",
			"--color-text":       "#000000",
			"--color-text-muted": "#555555",
			// Typography
			"--font-display":      "'Space Mono', monospace",
			"--font-body":         "'Space Mono', monospace",
			"--font-size-display": "3.5rem",
			"--font-size-heading": "1.75rem",
			"--font-size-body":    "1rem",
			"--font-size-caption": "0.75rem",
			// Shape
			"--radius-sm": "0px",
			"--radius-md": "0px",
			"--radius-lg": "0px",
			// Elevation/Shadows
			"--shadow-card": "4px 4px 0px #000000",
			"--shadow-btn":  "3px 3px 0px #000000",
			// Borders
			"--border-width": "2px",
			"--border-color": "#000000",
			// Layout
			"--layout-columns":    "1",
			"--layout-gap":        "2rem",
			"--content-max-width": "900px",
			"--section-padding":   "3rem 2rem",
		},
	},
	{
		Name:        "Cassette",
		Slug:        "cassette",
		Description: "NASA technical document aesthetic — dense, ruled, deliberate.",
		FontURL:     "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:ital,wght@0,300;0,400;0,500;0,700;1,400&family=Orbitron:wght@400;700&display=swap",
		CSSVars: map[string]string{
			"--color-bg":          "#f5f4ef",
			"--color-surface":     "#ffffff",
			"--color-surface-2":   "#e8e7e2",
			"--color-primary":     "#0b3d91",
			"--color-secondary":   "#1a5276",
			"--color-danger":      "#c0392b",
			"--color-caution":     "#c85200",
			"--color-text":        "#1a1a14",
			"--color-text-muted":  "#5a5a52",
			"--color-border":      "#c8c7c0",
			"--color-rule":        "#0b3d91",
			"--font-display":      "'Orbitron', sans-serif",
			"--font-body":         "'IBM Plex Mono', monospace",
			"--font-size-display": "2rem",
			"--font-size-heading": "0.875rem",
			"--font-size-body":    "0.8125rem",
			"--font-size-caption": "0.6875rem",
			"--radius-sm":         "0px",
			"--radius-md":         "2px",
			"--radius-lg":         "2px",
			"--shadow-card":       "0 1px 3px rgba(0,0,0,0.12)",
			"--shadow-btn":        "none",
			"--border-width":      "1px",
			"--border-color":      "#c8c7c0",
			"--layout-columns":    "2",
			"--layout-gap":        "1.5rem",
			"--content-max-width": "1200px",
			"--section-padding":   "2rem 2rem",
		},
	},
	{
		Name:        "Minimal",
		Slug:        "minimal",
		Description: "Calm, spacious, single-column design with generous whitespace.",
		FontURL:     "https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&display=swap",
		CSSVars: map[string]string{
			// Colors
			"--color-primary":    "#1a1a1a",
			"--color-secondary":  "#6b7280",
			"--color-accent":     "#3b82f6",
			"--color-bg":         "#fafafa",
			"--color-surface":    "#ffffff",
			"--color-text":       "#1a1a1a",
			"--color-text-muted": "#9ca3af",
			// Typography
			"--font-display":      "'Inter', sans-serif",
			"--font-body":         "'Inter', sans-serif",
			"--font-size-display": "3rem",
			"--font-size-heading": "1.5rem",
			"--font-size-body":    "1rem",
			"--font-size-caption": "0.8rem",
			// Shape
			"--radius-sm": "4px",
			"--radius-md": "8px",
			"--radius-lg": "16px",
			// Elevation
			"--shadow-card": "0 1px 3px rgba(0,0,0,0.08)",
			"--shadow-btn":  "0 1px 2px rgba(0,0,0,0.05)",
			// Borders
			"--border-width": "1px",
			"--border-color": "#e5e7eb",
			// Layout — single column, generous spacing
			"--layout-columns":    "1",
			"--layout-gap":        "4rem",
			"--content-max-width": "640px",
			"--section-padding":   "5rem 2rem",
		},
	},
}

// BySlug looks up a guide by its URL slug.
func BySlug(slug string) (Guide, bool) {
	for _, g := range All {
		if g.Slug == slug {
			return g, true
		}
	}
	return Guide{}, false
}
