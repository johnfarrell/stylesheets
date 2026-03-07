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
