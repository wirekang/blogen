package model

import "time"

// Article contains basic information of one article
type Article struct {
	Filename string
	Title    string
	Tags     []string
	Date     time.Time
}
