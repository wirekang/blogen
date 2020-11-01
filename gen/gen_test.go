package gen

import (
	"fmt"
	"os"
	"testing"

	"github.com/wirekang/blogen/model"
)

func TestCheckTemplateFiles(t *testing.T) {
	os.Chdir("..")
	if !checkTemplateFiles("example/templates") {
		t.FailNow()
	}
}

func TestGetTags(t *testing.T) {
	arts := make([]model.Article, 10)
	tags := make([]string, 12)
	for i := range tags {
		tags[i] = fmt.Sprintf("tag%d", i)
	}
	for i := range arts {
		arts[i] = model.Article{Tags: []string{tags[i], tags[i+1], tags[i+2]}}
	}
	re := getTagsFromArticles(arts)
	if len(re) != len(tags) {
		t.FailNow()
	}
	fmt.Println(re)
	fmt.Println(tags)
}
