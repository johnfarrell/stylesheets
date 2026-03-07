package cassette

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
/* [custom] - document body */
.cass-body{background:var(--color-bg);color:var(--color-text);font-family:var(--font-body);font-size:var(--font-size-body);}
/* [custom] - blue top rule for section breaks */
.cass-section-rule{border-top:3px solid var(--color-primary);padding-top:1.5rem;margin-top:2rem;}
/* [custom] - panel */
.cass-panel{background:var(--color-surface);border:1px solid var(--color-border);}
.cass-panel-header{background:var(--color-surface-2);border-bottom:1px solid var(--color-border);padding:0.4rem 0.75rem;font-size:var(--font-size-caption);font-weight:700;letter-spacing:0.08em;text-transform:uppercase;color:var(--color-text-muted);}
.cass-panel-header-blue{background:var(--color-primary);color:#fff;border-bottom:none;}
/* [custom] - buttons */
.cass-btn{font-family:var(--font-body);font-size:var(--font-size-caption);font-weight:700;letter-spacing:0.08em;text-transform:uppercase;border:1px solid var(--color-primary);color:var(--color-primary);background:transparent;padding:0.4rem 1rem;cursor:pointer;transition:background 0.1s,color 0.1s;}
.cass-btn:hover{background:var(--color-primary);color:#fff;}
.cass-btn:disabled{border-color:var(--color-border);color:var(--color-text-muted);cursor:not-allowed;background:transparent;}
.cass-btn-filled{background:var(--color-primary);color:#fff;}
.cass-btn-filled:hover{background:var(--color-secondary);border-color:var(--color-secondary);}
.cass-btn-danger{border-color:var(--color-danger);color:var(--color-danger);}
.cass-btn-danger:hover{background:var(--color-danger);color:#fff;}
.cass-btn-danger-filled{background:var(--color-danger);color:#fff;border-color:var(--color-danger);}
/* [custom] - form inputs styled as fillable fields */
.cass-input{font-family:var(--font-body);font-size:var(--font-size-body);background:var(--color-surface);border:none;border-bottom:1px solid var(--color-text);color:var(--color-text);padding:0.375rem 0;width:100%;}
.cass-input:focus{outline:none;border-bottom:2px solid var(--color-primary);}
.cass-input::placeholder{color:var(--color-text-muted);}
.cass-input-box{border:1px solid var(--color-border);padding:0.375rem 0.5rem;border-bottom-width:1px;}
.cass-input-box:focus{outline:none;border-color:var(--color-primary);box-shadow:0 0 0 2px rgba(11,61,145,0.12);}
.cass-label{font-size:var(--font-size-caption);font-weight:700;letter-spacing:0.06em;text-transform:uppercase;color:var(--color-text-muted);display:block;margin-bottom:0.25rem;}
.cass-field-group{border:1px solid var(--color-border);padding:0.75rem;position:relative;}
.cass-field-group-label{position:absolute;top:-0.6rem;left:0.5rem;background:var(--color-bg);padding:0 0.25rem;font-size:var(--font-size-caption);font-weight:700;color:var(--color-primary);letter-spacing:0.06em;text-transform:uppercase;}
.cass-check{width:0.875rem;height:0.875rem;border:1px solid var(--color-text);accent-color:var(--color-primary);cursor:pointer;}
/* [custom] - status indicator lights */
.cass-light{width:10px;height:10px;border-radius:50%;display:inline-block;flex-shrink:0;}
.cass-light-green{background:#27ae60;box-shadow:0 0 4px #27ae60;}
.cass-light-amber{background:#d4a017;box-shadow:0 0 4px #d4a017;}
.cass-light-red{background:var(--color-danger);box-shadow:0 0 4px var(--color-danger);}
.cass-light-blue{background:var(--color-primary);box-shadow:0 0 4px var(--color-primary);}
.cass-light-off{background:var(--color-border);box-shadow:none;}
@keyframes cass-blink{0%,49%{opacity:1}50%,100%{opacity:0.2}}
.cass-blink{animation:cass-blink 1.2s step-end infinite;}
@keyframes cass-pulse{0%,100%{opacity:1}50%{opacity:0.4}}
.cass-pulse{animation:cass-pulse 2s ease-in-out infinite;}
/* [custom] - NASA notice blocks */
.cass-notice{border-left:4px solid;padding:0.75rem 1rem;font-size:var(--font-size-body);}
.cass-notice-note{border-color:var(--color-primary);background:rgba(11,61,145,0.04);}
.cass-notice-caution{border-color:var(--color-caution);background:rgba(200,82,0,0.05);}
.cass-notice-warning{border-color:var(--color-danger);background:rgba(192,57,43,0.05);}
.cass-notice-label{font-weight:700;letter-spacing:0.08em;text-transform:uppercase;font-size:var(--font-size-caption);margin-bottom:0.25rem;}
/* [custom] - data table */
.cass-table{width:100%;border-collapse:collapse;font-size:var(--font-size-body);}
.cass-table th{background:var(--color-surface-2);border:1px solid var(--color-border);padding:0.4rem 0.6rem;font-size:var(--font-size-caption);font-weight:700;letter-spacing:0.06em;text-transform:uppercase;text-align:left;color:var(--color-text-muted);}
.cass-table td{border:1px solid var(--color-border);padding:0.4rem 0.6rem;vertical-align:top;}
.cass-table tr:hover td{background:rgba(11,61,145,0.03);}
.cass-table .cass-td-num{text-align:right;font-variant-numeric:tabular-nums;}
/* [custom] - progress bar */
.cass-progress-track{background:var(--color-surface-2);border:1px solid var(--color-border);height:1.125rem;overflow:hidden;}
.cass-progress-fill{height:100%;background:var(--color-primary);transition:width 0.5s ease;}
.cass-progress-fill-green{background:#27ae60;}
.cass-progress-fill-red{background:var(--color-danger);}
/* [custom] - tab nav */
.cass-tab{font-family:var(--font-body);font-size:var(--font-size-caption);font-weight:700;letter-spacing:0.06em;text-transform:uppercase;padding:0.5rem 1rem;border:1px solid var(--color-border);border-bottom:none;cursor:pointer;background:var(--color-surface-2);color:var(--color-text-muted);transition:all 0.1s;}
.cass-tab-active{background:var(--color-surface);color:var(--color-primary);border-top:2px solid var(--color-primary);}
/* [custom] - readout display */
.cass-readout{background:var(--color-surface-2);border:1px solid var(--color-border);padding:0.5rem 0.75rem;}
.cass-readout-value{font-size:1.375rem;font-weight:700;color:var(--color-primary);line-height:1;font-variant-numeric:tabular-nums;}
.cass-readout-value-danger{color:var(--color-danger);}
.cass-readout-unit{font-size:var(--font-size-caption);color:var(--color-text-muted);}
/* [custom] - modal overlay */
.cass-overlay{position:fixed;inset:0;background:rgba(26,26,20,0.65);z-index:50;display:flex;align-items:center;justify-content:center;}
.cass-modal{background:var(--color-surface);border:2px solid var(--color-primary);max-width:480px;width:100%;}
.cass-modal-header{background:var(--color-primary);color:#fff;padding:0.6rem 1rem;font-weight:700;letter-spacing:0.08em;text-transform:uppercase;font-size:var(--font-size-caption);}
.cass-modal-header-danger{background:var(--color-danger);}
/* [custom] - value color helpers */
.cass-value{font-weight:700;color:var(--color-primary);font-variant-numeric:tabular-nums;}
.cass-value-danger{color:var(--color-danger);}
.cass-value-ok{color:#27ae60;}
.cass-value-warn{color:var(--color-caution);}
/* [custom] - doc header block */
.cass-doc-header{border:2px solid var(--color-text);}
.cass-doc-cell{border:1px solid var(--color-border);padding:0.35rem 0.6rem;}
.cass-doc-cell-label{font-size:var(--font-size-caption);font-weight:700;text-transform:uppercase;letter-spacing:0.06em;color:var(--color-text-muted);}
.cass-doc-cell-value{font-size:var(--font-size-body);font-weight:500;}
/* [custom] - step tracker */
.cass-step-line{height:2px;background:var(--color-border);flex:1;}
.cass-step-line-done{background:var(--color-primary);}
.cass-step-circle{width:1.75rem;height:1.75rem;border-radius:50%;border:2px solid var(--color-border);display:flex;align-items:center;justify-content:center;font-size:var(--font-size-caption);font-weight:700;background:var(--color-surface);flex-shrink:0;}
.cass-step-circle-done{border-color:var(--color-primary);background:var(--color-primary);color:#fff;}
.cass-step-circle-active{border-color:var(--color-primary);color:var(--color-primary);}
`
}
