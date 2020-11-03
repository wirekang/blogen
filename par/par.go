//Package par parses things.
package par

import (
	"errors"
	"time"

	"github.com/gomarkdown/markdown"
)

// Extract extracts settings and markdown bytes from src.
// Settings and markdown are separated by sep.
func Extract(src []byte, sep byte) (settings []byte, markdown []byte, err error) {
	for i, b := range src {
		if b == sep {
			settings = src[:i]
			markdown = src[i+1:]
			return
		}
	}
	return nil, nil, errors.New("can't sep by " + string(sep))
}

// MarkdownToHTML converts markdown bytes to html bytes.
func MarkdownToHTML(md []byte) (html []byte) {
	return markdown.ToHTML(md, nil, nil)
}

// DateToTime parses date from following format:
//
// - - -
//
// 2020-12-30
func DateToTime(date string) (time.Time, error) {
	return time.Parse("2006-1-2", date)
}
