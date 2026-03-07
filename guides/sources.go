package guides

import "embed"

// SourceFS holds the embedded .templ source files for snippet extraction.
//
//go:embed brutalist/brutalist.templ minimal/minimal.templ cassette/cassette.templ
var SourceFS embed.FS
