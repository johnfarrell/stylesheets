package guides

import "embed"

// SourceFS holds the embedded guide subdirectories for snippet extraction.
// All .templ files in each guide directory are scanned for snippet markers.
//
//go:embed brutalist minimal cassette glass bento swiss terminal retro newspaper shelf tracker
var SourceFS embed.FS
