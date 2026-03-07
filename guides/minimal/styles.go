package minimal

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
.min-btn-primary {
    background: var(--color-accent);
    color: #ffffff;
    border: var(--border-width) solid var(--color-accent);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    font-weight: 500;
    transition: opacity 0.15s, box-shadow 0.15s, transform 0.15s;
}
.min-btn-primary:hover {
    opacity: 0.88;
    box-shadow: 0 4px 12px rgba(59,130,246,0.25);
    transform: translateY(-1px);
}
.min-btn-secondary {
    background: var(--color-surface);
    color: var(--color-primary);
    border: var(--border-width) solid var(--border-color);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    font-weight: 500;
    transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
}
.min-btn-secondary:hover {
    border-color: var(--color-secondary);
    box-shadow: 0 2px 8px rgba(0,0,0,0.08);
    transform: translateY(-1px);
}
.min-btn-ghost {
    background: transparent;
    color: var(--color-accent);
    border: var(--border-width) solid transparent;
    border-radius: var(--radius-md);
    font-family: var(--font-body);
    font-weight: 500;
    transition: background 0.15s;
}
.min-btn-ghost:hover {
    background: rgba(59,130,246,0.06);
}
.min-card {
    background: var(--color-surface);
    border: var(--border-width) solid var(--border-color);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-card);
}
.min-input {
    background: var(--color-surface);
    border: var(--border-width) solid var(--border-color);
    border-radius: var(--radius-md);
    font-family: var(--font-body);
    color: var(--color-text);
    transition: border-color 0.15s, box-shadow 0.15s;
}
.min-input:focus {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px rgba(59,130,246,0.12);
}
.min-divider {
    border: none;
    border-top: var(--border-width) solid var(--border-color);
}
`
}
