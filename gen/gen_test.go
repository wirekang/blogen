package gen

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/wirekang/blogen/model"
)

func TestTemplate(t *testing.T) {
	os.Chdir("..")
	tems := Tems{}
	tems.Addr = "localhost"
	tems.Title = "Title of tems"
	tems.Tags = []string{"Tag1", "태그2", "ㅁㄴㅇㄹ"}
	tem := template.New("base.html")
	var err error
	tem, err = tem.ParseGlob("example/templates/*.html")
	if err != nil {
		fmt.Println(err)
	}
	tem.Execute(os.Stdout, tems)
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
