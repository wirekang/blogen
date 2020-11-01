package gen

import (
	"fmt"
	"os"
	"path"

	"github.com/wirekang/blogen/fl"
	"github.com/wirekang/blogen/model"
)

// TemplateFiles is files that should exists in template directory
var TemplateFiles = []string{"base.html", "main.html", "item.html", "single.html"}

// GenerateFromTemplate generates static site from template.
func GenerateFromTemplate(arts []model.Article, htmlDir string, templateDir string, outDir string) (ok bool) {
	fmt.Println("Generating...")
	if !checkTemplateFiles(templateDir) {
		return false
	}
	if !fl.IsExists(outDir) {
		os.Mkdir(outDir, 0644)
	}

	tags := getTagsFromArticles(arts)

	return true
}

func checkTemplateFiles(templateDir string) (ok bool) {
	for _, file := range TemplateFiles {
		if !fl.IsExists(path.Join(templateDir, file)) {
			return false
		}
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
