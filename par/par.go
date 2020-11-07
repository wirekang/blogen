//Package par parses things.
package par

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/gomarkdown/markdown"
)

// Extract extracts config and markdown bytes from src.
// Config and markdown are separated by sep.
func Extract(file string, sep byte) (config []byte, markdown []byte, err error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	for i, b := range src {
		if b == sep {
			config = src[:i]
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
