package guides

import "embed"

// SourceFS holds the embedded guide subdirectories for snippet extraction.
// Each guide's .templ file is read by path ({slug}/{slug}.templ) at runtime.
//
//go:embed brutalist minimal cassette glass bento swiss
var SourceFS embed.FS
