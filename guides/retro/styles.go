package retro

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - Win95-style 3D beveled borders not achievable with Tailwind */
.retro-raised {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
}
.retro-inset {
    border: 2px solid;
    border-color: #808080 #ffffff #ffffff #808080;
    background: #fff;
}
/* [custom] - window chrome with title bar */
.retro-window {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    box-shadow: var(--shadow-card);
}
.retro-titlebar {
    background: linear-gradient(90deg, var(--color-primary), #1084d0);
    color: #ffffff;
    padding: 0.25rem 0.5rem;
    font-weight: 700;
    font-size: var(--font-size-caption);
    display: flex;
    align-items: center;
    justify-content: space-between;
    user-select: none;
    cursor: default;
}
.retro-titlebar-inactive {
    background: linear-gradient(90deg, #808080, #a0a0a0);
}
/* [custom] - window control buttons (close/minimize) */
.retro-winbtn {
    width: 16px;
    height: 14px;
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-size: 8px;
    line-height: 10px;
    text-align: center;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    font-family: var(--font-body);
}
.retro-winbtn:active {
    border-color: #808080 #ffffff #ffffff #808080;
}
/* [custom] - raised button with press state */
.retro-btn {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    padding: 0.25rem 1rem;
    cursor: pointer;
}
.retro-btn:active {
    border-color: #808080 #ffffff #ffffff #808080;
    padding: 0.3rem 0.95rem 0.2rem 1.05rem;
}
.retro-btn-primary {
    background: var(--color-surface);
    outline: 1px dotted #000;
    outline-offset: -4px;
}
/* [custom] - inset input field */
.retro-input {
    border: 2px solid;
    border-color: #808080 #ffffff #ffffff #808080;
    background: #ffffff;
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    padding: 0.2rem 0.4rem;
    color: #000;
    width: 100%;
}
.retro-input:focus {
    outline: none;
}
/* [custom] - desktop icon */
.retro-icon {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    cursor: pointer;
    padding: 0.5rem;
    border: 1px solid transparent;
    font-size: var(--font-size-caption);
    color: #ffffff;
    text-shadow: 1px 1px 1px #000;
}
.retro-icon:hover {
    border: 1px dotted #ffffff;
}
.retro-icon-img {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
}
/* [custom] - taskbar at bottom */
.retro-taskbar {
    background: var(--color-surface);
    border-top: 2px solid;
    border-color: #ffffff;
    padding: 0.25rem;
    display: flex;
    gap: 0.25rem;
    align-items: center;
}
.retro-taskbar-btn {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    padding: 0.15rem 0.5rem;
    cursor: pointer;
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.retro-taskbar-btn-active {
    border-color: #808080 #ffffff #ffffff #808080;
    background: #dfdfdf;
}
/* [custom] - start button */
.retro-start-btn {
    border: 2px solid;
    border-color: #ffffff #808080 #808080 #ffffff;
    background: var(--color-surface);
    font-family: var(--font-body);
    font-weight: 700;
    font-size: var(--font-size-caption);
    padding: 0.15rem 0.5rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.25rem;
}
.retro-start-btn:active {
    border-color: #808080 #ffffff #ffffff #808080;
}
/* [custom] - checkbox and radio with retro styling */
.retro-check {
    width: 13px;
    height: 13px;
    accent-color: var(--color-primary);
}
@media (max-width: 768px) {
    .retro-window { position: static !important; transform: none !important; margin-bottom: 1rem; }
}
`
}
