package tracker

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - dark panel with subtle border and shadow */
.trk-panel {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-card);
}
.trk-panel-header {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--color-text-muted);
    padding: 0.4rem 0.75rem;
    border-bottom: 1px solid var(--color-border);
    border-top: 2px solid var(--color-primary);
    background: var(--color-surface-2);
}
.trk-panel-header-gold {
    background: var(--color-primary);
    color: var(--color-bg);
    border-bottom: none;
    border-top: none;
}
.trk-panel-elevated {
    background: var(--color-surface-2);
}
/* [custom] - status indicator dots */
.trk-status-light {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--color-text-muted);
    flex-shrink: 0;
}
.trk-status-complete {
    background: var(--color-accent);
    box-shadow: 0 0 6px var(--color-accent);
}
.trk-status-progress {
    background: var(--color-warning);
    box-shadow: 0 0 6px var(--color-warning);
}
.trk-status-locked {
    background: var(--color-danger);
    box-shadow: 0 0 4px var(--color-danger);
}
/* [custom] - pulsing animation for in-progress items */
@keyframes trk-pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
}
.trk-status-pulse {
    animation: trk-pulse 2s ease-in-out infinite;
}
/* [custom] - progress bar with bordered track like technical readout */
.trk-progress-bar {
    height: 1rem;
    background: var(--color-surface-2);
    border: 1px solid var(--color-border);
    overflow: hidden;
}
.trk-progress-fill {
    height: 100%;
    background: var(--color-primary);
    transition: width 0.3s ease;
}
/* [custom] - sidebar tree navigation */
.trk-tree {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
}
.trk-tree-node {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.25rem 0.5rem;
    cursor: pointer;
    color: var(--color-text);
    border-left: 2px solid transparent;
    border-bottom: 1px solid var(--color-surface-2);
    transition: background 0.1s, border-color 0.1s;
}
.trk-tree-node:hover {
    background: var(--color-surface-2);
}
.trk-tree-node-active {
    border-left-color: var(--color-primary);
    background: var(--color-surface-2);
    color: var(--color-primary);
}
.trk-tree-toggle {
    display: inline-flex;
    width: 1rem;
    justify-content: center;
    font-size: 0.625rem;
    color: var(--color-text-muted);
    transition: transform 0.15s;
    flex-shrink: 0;
    user-select: none;
}
.trk-tree-toggle-open {
    transform: rotate(90deg);
}
/* [custom] - buttons with gold border accent */
.trk-btn {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-primary);
    background: transparent;
    border: 1px solid var(--color-primary);
    padding: 0.4rem 1rem;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
}
.trk-btn:hover {
    background: var(--color-primary);
    color: var(--color-bg);
}
.trk-btn-primary {
    background: var(--color-primary);
    color: var(--color-bg);
    border-color: var(--color-primary);
}
.trk-btn-primary:hover {
    background: #b8993e;
    border-color: #b8993e;
}
.trk-btn-danger {
    color: var(--color-danger);
    border-color: var(--color-danger);
}
.trk-btn-danger:hover {
    background: var(--color-danger);
    color: var(--color-text);
}
/* [custom] - dark inset input fields */
.trk-input {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    padding: 0.4rem 0.75rem;
    width: 100%;
    transition: border-color 0.15s;
}
.trk-input:focus {
    outline: none;
    border-color: var(--color-primary);
}
.trk-input::placeholder {
    color: var(--color-text-muted);
}
/* [custom] - search input */
.trk-search {
    font-family: var(--font-display);
    font-size: var(--font-size-caption);
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    padding: 0.4rem 0.75rem;
    width: 100%;
}
.trk-search:focus {
    outline: none;
    border-color: var(--color-primary);
}
.trk-search::placeholder {
    color: var(--color-text-muted);
}
/* [custom] - large monospace readout values */
.trk-readout {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.375rem;
    color: var(--color-primary);
    letter-spacing: 0.02em;
    font-variant-numeric: tabular-nums;
    line-height: 1;
}
/* [custom] - small category tag pills */
.trk-tag {
    display: inline-block;
    font-family: var(--font-display);
    font-size: 0.625rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    padding: 0.15rem 0.5rem;
    border-radius: 2px;
    background: var(--color-surface-2);
    color: var(--color-text-muted);
    border: 1px solid var(--color-border);
}
.trk-tag-skill { color: var(--color-info); border-color: var(--color-info); }
.trk-tag-project { color: var(--color-primary); border-color: var(--color-primary); }
.trk-tag-certification { color: var(--color-accent); border-color: var(--color-accent); }
.trk-tag-challenge { color: var(--color-danger); border-color: var(--color-danger); }
/* [custom] - dependency graph nodes */
.trk-dep-node {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.75rem;
    font-size: var(--font-size-caption);
    min-width: 120px;
    position: relative;
}
.trk-dep-node-complete { border-left: 3px solid var(--color-accent); }
.trk-dep-node-progress { border-left: 3px solid var(--color-warning); }
.trk-dep-node-locked { border-left: 3px solid var(--color-danger); }
.trk-dep-node-dimmed { opacity: 0.3; }
/* [custom] - connecting lines between dep nodes */
.trk-dep-line {
    border-top: 1px dashed var(--color-border);
    width: 2rem;
    align-self: center;
    flex-shrink: 0;
}
/* [custom] - horizontal divider */
.trk-rule {
    border-top: 1px solid var(--color-border);
}
/* [custom] - gold top rule for section breaks (cassette-inspired) */
.trk-section-rule {
    border-top: 3px solid var(--color-primary);
    padding-top: 1.5rem;
    margin-top: 2rem;
}
/* [custom] - gold text glow for emphasis */
.trk-glow {
    text-shadow: 0 0 8px rgba(200,170,110,0.4);
}
/* [custom] - bordered input fields like fillable form fields */
.trk-input-underline {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--color-text-muted);
    color: var(--color-text);
    padding: 0.375rem 0;
    width: 100%;
}
.trk-input-underline:focus {
    outline: none;
    border-bottom: 2px solid var(--color-primary);
}
.trk-input-underline::placeholder {
    color: var(--color-text-muted);
}
`
}

// statusClass returns the CSS class for a tracker item status.
func statusClass(status string) string {
	switch status {
	case "complete":
		return "trk-status-complete"
	case "progress":
		return "trk-status-progress trk-status-pulse"
	case "locked":
		return "trk-status-locked"
	default:
		return ""
	}
}
