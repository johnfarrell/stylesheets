package brutalist

import (
	"strings"
)

// buildCSSVars generates the :root CSS variable declarations from the guide's CSSVars map.
func buildCSSVars(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(v)
		sb.WriteString(";")
	}
	return sb.String()
}

// guideStyles returns the static CSS classes for this guide.
func guideStyles() string {
	return `
.brut-btn-primary {
    background: var(--color-primary);
    color: var(--color-bg);
    border: var(--border-width) solid var(--border-color);
    /* [custom] - per-guide box-shadow token */
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    transition: transform 0.1s, box-shadow 0.1s;
}
.brut-btn-primary:hover {
    /* [custom] - brutalist shift effect on hover */
    transform: translate(-2px, -2px);
    box-shadow: 5px 5px 0px var(--border-color);
}
.brut-btn-secondary {
    background: var(--color-bg);
    color: var(--color-primary);
    border: var(--border-width) solid var(--border-color);
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    transition: transform 0.1s, box-shadow 0.1s;
}
.brut-btn-secondary:hover {
    transform: translate(-2px, -2px);
    box-shadow: 5px 5px 0px var(--border-color);
}
.brut-card {
    background: var(--color-surface);
    border: var(--border-width) solid var(--border-color);
    /* [custom] - per-guide card shadow token */
    box-shadow: var(--shadow-card);
}
.brut-input {
    background: var(--color-bg);
    border: var(--border-width) solid var(--border-color);
    font-family: var(--font-body);
    border-radius: 0;
}
.brut-input:focus {
    outline: 3px solid var(--color-secondary);
    outline-offset: 0;
}
.brut-btn-accent {
    background: var(--color-accent);
    color: var(--color-primary);
    border: var(--border-width) solid var(--border-color);
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    transition: transform 0.1s, box-shadow 0.1s;
}
.brut-btn-accent:hover {
    transform: translate(-2px, -2px);
    box-shadow: 5px 5px 0px var(--border-color);
}
`
}
