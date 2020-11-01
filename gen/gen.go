package gen

import (
	"fmt"
	"io/ioutil"
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
	Tags     []Tag
}

// Tag is structure for tag
type Tag struct {
	ID    int
	Count int
	Name  string
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

	ok = executeIndex(tems, templateDir, outDir)

	if !ok {
		fmt.Println("index.html failed.")
		return false
	}

	for _, art := range arts {
		fmt.Printf("%s...  ", art.Filename)
		ok = executeArticle(art, htmlDir, templateDir, outDir)
		if ok {
			fmt.Println("Ok")
		}
	}

	for _, tag := range tems.Tags {
		fmt.Printf("Tag %s...  ", tag.Name)
		ok := executeTag(tems, tag, templateDir, outDir)
		if ok {
			fmt.Println("Ok")
		}
	}
	return true

}

func executeTag(tems Tems, tag Tag, templateDir string, outDir string) (ok bool) {
	tem := template.New("index.html")
	var err error
	tem, err = tem.ParseFiles(path.Join(templateDir, "index.html"), path.Join(templateDir, "list.html"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	dir := fmt.Sprintf("tag%d", tag.ID)
	os.Mkdir(path.Join(outDir, dir), 0755)
	var f *os.File
	f, err = os.Create(path.Join(outDir, dir, "index.html"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	newArt := make([]model.Article, 0)
	for _, art := range tems.Articles {
	TagLoop:
		for _, t := range art.Tags {
			if t == tag.Name {
				newArt = append(newArt, art)
				break TagLoop
			}
		}
	}
	tems.Articles = newArt
	err = tem.Execute(f, tems)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func executeIndex(tems Tems, templateDir string, outDir string) (ok bool) {
	tem := template.New("index.html")
	files := []string{
		path.Join(templateDir, "index.html"),
		path.Join(templateDir, "list.html"),
	}
	var err error
	tem, err = tem.ParseFiles(files...)
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

func executeArticle(art model.Article, htmlDir string, templateDir string, outDir string) (ok bool) {
	tem := template.New("index.html")
	var err error
	tem, err = tem.ParseFiles(path.Join(templateDir, "index.html"), path.Join(templateDir, "article.html"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	var bytes []byte
	bytes, err = ioutil.ReadFile(path.Join(htmlDir, art.Filename))
	if err != nil {
		fmt.Println(err)
		return false
	}
	art.HTML = string(bytes)
	var f *os.File
	os.Mkdir(path.Join(outDir, art.Filename), 0755)

	f, err = os.Create(path.Join(outDir, art.Filename, "index.html"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = tem.Execute(f, art)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func getTagsFromArticles(arts []model.Article) []Tag {
	m := make(map[string]int, len(arts))
	for _, a := range arts {
		for _, t := range a.Tags {
			m[t]++
		}
	}
	tags := make([]Tag, 0, len(arts))
	id := 0
	for tag, count := range m {
		tags = append(tags, Tag{ID: id, Count: count, Name: tag})
		id++
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
