package glass

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
/* [custom] - backdrop-filter not achievable with Tailwind utilities */
.glass-panel {
    background: var(--frost-bg);
    backdrop-filter: blur(var(--frost-blur));
    -webkit-backdrop-filter: blur(var(--frost-blur));
    border: var(--border-width) solid var(--color-border);
    border-radius: var(--radius-md);
}
.glass-btn-primary {
    background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
    color: #fff;
    border: none;
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-btn);
    font-family: var(--font-body);
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.2s, transform 0.1s;
}
.glass-btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
.glass-btn-ghost {
    background: var(--frost-bg);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    border: var(--border-width) solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    font-family: var(--font-body);
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
}
.glass-btn-ghost:hover { background: rgba(255,255,255,0.14); }
/* [custom] - radial gradient background not achievable with Tailwind alone */
.glass-bg {
    background: radial-gradient(ellipse at 20% 50%, rgba(167,139,250,0.15) 0%, transparent 60%),
                radial-gradient(ellipse at 80% 20%, rgba(96,165,250,0.1) 0%, transparent 50%),
                var(--color-bg);
    min-height: 100%;
}
.glass-input {
    background: rgba(255,255,255,0.06);
    border: 1px solid rgba(255,255,255,0.15);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    font-family: var(--font-body);
    padding: 0.6rem 0.75rem;
    width: 100%;
    transition: border-color 0.2s, background 0.2s;
}
.glass-input:focus {
    outline: none;
    border-color: var(--color-primary);
    background: rgba(255,255,255,0.10);
}
`
}
