package guides

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Highlight returns chroma-highlighted HTML spans for the given code and language.
// Output contains bare <span> tokens with CSS classes — no surrounding <pre> tag.
// Falls back to returning the original code as plain text on any error.
func Highlight(code, lang string) string {
	if code == "" {
		return ""
	}
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	formatter := chromahtml.New(
		chromahtml.WithClasses(true),
		chromahtml.PreventSurroundingPre(true),
	)
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}
	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return code
	}
	return buf.String()
}

// DetectLang returns "html" for template/HTML code and "go" for Go code.
func DetectLang(code string) string {
	trimmed := strings.TrimSpace(code)
	if strings.Contains(code, "<!--") || (len(trimmed) > 0 && trimmed[0] == '<') {
		return "html"
	}
	return "go"
}
