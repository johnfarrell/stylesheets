package bento

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
.bento-card {
    background: var(--color-surface);
    border: var(--border-width) solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-card);
    padding: 1.5rem;
}
.bento-btn {
    background: var(--color-primary);
    color: #fff;
    border: none;
    border-radius: var(--radius-sm);
    font-family: var(--font-body);
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.15s;
    padding: 0.5rem 1rem;
}
.bento-btn:hover { opacity: 0.9; }
/* [custom] - CSS grid variable-span tiles not achievable with static Tailwind classes */
.bento-grid {
    display: grid;
    grid-template-columns: repeat(12, 1fr);
    gap: 1rem;
}
.bento-span-4 { grid-column: span 4; }
.bento-span-6 { grid-column: span 6; }
.bento-span-8 { grid-column: span 8; }
.bento-span-12 { grid-column: span 12; }
@media (max-width: 768px) {
    .bento-span-4, .bento-span-6, .bento-span-8, .bento-span-12 { grid-column: span 12; }
}
.bento-input {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    font-family: var(--font-body);
    padding: 0.5rem 0.75rem;
    width: 100%;
    transition: border-color 0.15s;
}
.bento-input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px rgba(99,102,241,0.1);
}
`
}
