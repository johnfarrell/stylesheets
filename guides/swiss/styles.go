package swiss

import "strings"

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

func guideStyles() string {
	return `
/* [custom] - strict typographic grid not achievable with Tailwind utilities alone */
.swiss-rule { border-top: 2px solid var(--color-border); }
.swiss-rule-red { border-top: 3px solid var(--color-primary); }
.swiss-label {
    font-family: var(--font-body);
    font-size: 0.625rem;
    font-weight: 700;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: var(--color-text-muted);
}
.swiss-btn {
    background: var(--color-secondary);
    color: #fff;
    border: 2px solid var(--color-secondary);
    font-family: var(--font-body);
    font-weight: 700;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
    padding: 0.75rem 1.5rem;
}
.swiss-btn:hover { background: var(--color-primary); border-color: var(--color-primary); }
.swiss-btn-outline {
    background: transparent;
    color: var(--color-secondary);
    border: 2px solid var(--color-secondary);
    font-family: var(--font-body);
    font-weight: 700;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
    padding: 0.75rem 1.5rem;
}
.swiss-btn-outline:hover { background: var(--color-secondary); color: #fff; }
.swiss-btn-red {
    background: var(--color-primary);
    color: #fff;
    border: 2px solid var(--color-primary);
    font-family: var(--font-body);
    font-weight: 700;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    cursor: pointer;
    padding: 0.75rem 1.5rem;
}
/* [custom] - CSS grid strict column layout */
.swiss-grid {
    display: grid;
    grid-template-columns: repeat(12, 1fr);
    gap: 0;
    border-left: 2px solid var(--color-border);
}
.swiss-col-4 { grid-column: span 4; border-right: 2px solid var(--color-border); }
.swiss-col-6 { grid-column: span 6; border-right: 2px solid var(--color-border); }
.swiss-col-8 { grid-column: span 8; border-right: 2px solid var(--color-border); }
.swiss-col-12 { grid-column: span 12; border-right: 2px solid var(--color-border); }
@media (max-width: 768px) {
    .swiss-col-4, .swiss-col-6, .swiss-col-8 { grid-column: span 12; }
}
.swiss-input {
    background: var(--color-surface);
    border: 2px solid var(--color-border);
    border-radius: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 1rem;
    padding: 0.6rem 0.75rem;
    width: 100%;
    transition: border-color 0.1s;
}
.swiss-input:focus {
    outline: none;
    border-color: var(--color-primary);
}
`
}
