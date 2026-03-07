package terminal

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - CRT scanline overlay via pseudo-element */
.term-screen { position: relative; }
.term-screen::after {
    content: "";
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: repeating-linear-gradient(
        0deg,
        transparent,
        transparent 2px,
        rgba(0,0,0,0.15) 2px,
        rgba(0,0,0,0.15) 4px
    );
    z-index: 1;
}
/* [custom] - CRT text glow not achievable with Tailwind */
.term-glow { text-shadow: 0 0 5px currentColor; }
.term-glow-strong { text-shadow: 0 0 8px currentColor, 0 0 15px currentColor; }
/* [custom] - blinking cursor animation */
@keyframes term-blink { 0%,49% { opacity: 1; } 50%,100% { opacity: 0; } }
.term-cursor { animation: term-blink 1s step-end infinite; }
/* [custom] - terminal panel */
.term-panel {
    background: var(--color-surface);
    border: var(--border-width) solid var(--color-border);
    box-shadow: var(--shadow-card);
}
/* [custom] - terminal button */
.term-btn {
    background: transparent;
    color: var(--color-primary);
    border: var(--border-width) solid var(--color-primary);
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 500;
    cursor: pointer;
    transition: background 0.1s, color 0.1s, box-shadow 0.1s;
    padding: 0.4rem 1rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
.term-btn:hover {
    background: var(--color-primary);
    color: var(--color-bg);
    box-shadow: var(--shadow-btn);
}
.term-btn-danger { border-color: var(--color-danger); color: var(--color-danger); }
.term-btn-danger:hover { background: var(--color-danger); color: var(--color-bg); }
/* [custom] - terminal input with glow focus */
.term-input {
    background: var(--color-bg);
    border: var(--border-width) solid var(--color-border);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    padding: 0.4rem 0.5rem;
    width: 100%;
    caret-color: var(--color-primary);
}
.term-input:focus {
    outline: none;
    box-shadow: 0 0 6px rgba(0,255,65,0.3);
}
/* [custom] - HTMX loading indicator */
.term-indicator { display: none; }
.htmx-request .term-indicator,
.htmx-request.term-indicator { display: inline; }
/* [custom] - file browser item highlight */
.term-file-active {
    background: var(--color-primary);
    color: var(--color-bg);
}
`
}
