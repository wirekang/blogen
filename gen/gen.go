package gen

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/wirekang/blogen/fl"
	"github.com/wirekang/blogen/model"
)

// BaseInfo is base information of site
type BaseInfo struct {
	Title string
	Addr  string
}

// Tems is structure for template
type Tems struct {
	BaseInfo
	Articles []model.Article
	Tags     []string
}

// GenerateFromTemplate generates static site from template.
func GenerateFromTemplate(bi BaseInfo, arts []model.Article, htmlDir string, templateDir string, outDir string) (ok bool) {
	fmt.Println("Generating...")
	if !fl.IsExists(outDir) {
		os.Mkdir(outDir, 0755)
	}
	arts = sortArticlesByDate(arts)
	tems := Tems{}
	tems.Title = bi.Title
	tems.Addr = bi.Addr
	tems.Articles = arts
	tems.Tags = getTagsFromArticles(arts)
	tem := template.New("base.html")
	var err error
	tem, err = tem.ParseGlob(path.Join(templateDir, "*.html"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	var f *os.File
	f, err = os.Create(path.Join(outDir, "index.html"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = tem.Execute(f, tems)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

func getTagsFromArticles(arts []model.Article) []string {
	m := make(map[string]bool, len(arts))
	for _, a := range arts {
		for _, t := range a.Tags {
			m[t] = true
		}
	}
	tags := make([]string, 0, len(arts))
	for k := range m {
		tags = append(tags, k)
	}
	return tags
}

func sortArticlesByDate(arts []model.Article) []model.Article {
	var tmp model.Article
	for i := len(arts) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if arts[j].Date.Unix() < arts[j+1].Date.Unix() {
				tmp = arts[j]
				arts[j] = arts[j+1]
				arts[j+1] = tmp
			}
		}
	}
	return arts
}
