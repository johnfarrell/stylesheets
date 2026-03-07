package glass

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
/* [custom] - frosted tab pill with gradient active indicator */
.glass-tab {
    padding: 0.5rem 1.25rem;
    font-size: 0.8125rem;
    font-weight: 500;
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    transition: background 0.2s, color 0.2s;
    border: none;
    font-family: var(--font-body);
}
.glass-tab:hover { background: rgba(255,255,255,0.06); color: var(--color-text); }
.glass-tab-active {
    background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
    color: #fff;
}
.glass-tab-active:hover { opacity: 0.9; }
/* [custom] - frosted alert with gradient left border */
.glass-alert {
    border-radius: var(--radius-md);
    padding: 1rem 1.25rem;
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    font-size: 0.875rem;
    font-family: var(--font-body);
    background: var(--frost-bg);
    backdrop-filter: blur(var(--frost-blur));
    -webkit-backdrop-filter: blur(var(--frost-blur));
    border: 1px solid var(--color-border);
}
`
}
