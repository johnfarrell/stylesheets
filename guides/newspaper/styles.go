package newspaper

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - multi-column text flow not achievable with Tailwind */
.news-columns-2 { column-count: 2; column-gap: 2rem; column-rule: 1px solid var(--color-border); }
.news-columns-3 { column-count: 3; column-gap: 2rem; column-rule: 1px solid var(--color-border); }
@media (max-width: 768px) {
    .news-columns-2, .news-columns-3 { column-count: 1; }
}
/* [custom] - drop cap not achievable with Tailwind */
.news-dropcap::first-letter {
    float: left;
    font-family: var(--font-display);
    font-size: 4rem;
    line-height: 0.8;
    padding-right: 0.5rem;
    padding-top: 0.25rem;
    color: var(--color-primary);
    font-weight: 900;
}
/* [custom] - section rule line */
.news-rule { border-top: 1px solid var(--color-border); }
.news-rule-thick { border-top: 3px solid var(--color-border); }
.news-rule-red { border-top: 2px solid var(--color-primary); }
/* [custom] - byline with small-caps */
.news-byline {
    font-family: var(--font-body);
    font-variant: small-caps;
    font-size: var(--font-size-caption);
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
}
/* [custom] - pull quote */
.news-pullquote {
    border-left: 3px solid var(--color-primary);
    padding-left: 1.5rem;
    margin: 2rem 0;
    font-family: var(--font-display);
    font-style: italic;
    font-size: 1.5rem;
    line-height: 1.4;
    color: var(--color-secondary);
}
/* [custom] - masthead */
.news-masthead {
    text-align: center;
    border-top: 3px double var(--color-border);
    border-bottom: 3px double var(--color-border);
    padding: 1rem 0;
}
/* [custom] - article card */
.news-card {
    border-top: 2px solid var(--color-border);
    padding-top: 1rem;
}
/* [custom] - button styled as editorial link */
.news-btn {
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 600;
    color: var(--color-secondary);
    background: none;
    border: 1px solid var(--color-border);
    padding: 0.4rem 1rem;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
}
.news-btn:hover {
    background: var(--color-secondary);
    color: var(--color-bg);
}
.news-btn-primary {
    background: var(--color-primary);
    color: #fff;
    border-color: var(--color-primary);
}
.news-btn-primary:hover {
    background: #a01818;
    border-color: #a01818;
}
/* [custom] - input styled as editorial underline field */
.news-input {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--color-border);
    color: var(--color-text);
    padding: 0.375rem 0;
    width: 100%;
}
.news-input:focus {
    outline: none;
    border-bottom: 2px solid var(--color-primary);
}
/* [custom] - breaking news banner */
.news-breaking {
    background: var(--color-primary);
    color: #fff;
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    padding: 0.5rem 1rem;
}
/* [custom] - reading progress bar */
.news-progress {
    position: fixed;
    top: 0;
    left: 0;
    height: 3px;
    background: var(--color-primary);
    z-index: 100;
    transition: width 0.1s linear;
}
/* [custom] - headline sizes */
.news-headline-lg {
    font-family: var(--font-display);
    font-weight: 900;
    font-size: 2.5rem;
    line-height: 1.1;
    color: var(--color-secondary);
}
.news-headline-md {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.5rem;
    line-height: 1.2;
    color: var(--color-secondary);
}
.news-headline-sm {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    line-height: 1.3;
    color: var(--color-secondary);
}
`
}
