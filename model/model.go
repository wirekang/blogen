package model

import (
	"time"
)

// Article contains basic information of one article. HTML is empty except when parsing article.html
type Article struct {
	Filename   string
	Title      string
	Tags       []string
	Date       time.Time
	StringDate string
	HTML       string
}
