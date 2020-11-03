//Package gen generates sites.
package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/wirekang/blogen/cvt"
	"github.com/wirekang/blogen/fl"
)

// BaseInfo is base information of site
type BaseInfo struct {
	Title string
	Addr  string
}

// Tems is structure for template
type Tems struct {
	BaseInfo
	Articles []Article
	Tags     []Tag
}

// Tag is structure for tag
type Tag struct {
	ID    int
	Count int
	Name  string
}

// Article is structure for article
type Article struct {
	Filename   string
	Title      string
	Tags       []Tag
	Date       time.Time
	StringDate string
	HTML       string
}

// GenerateFromTemplate generates static site from template.
func GenerateFromTemplate(bi BaseInfo, cArts []cvt.Article, htmlDir string, templateDir string, outDir string) (ok bool) {
	fmt.Println("Generating...")
	err := fl.MakeIfNotExist(outDir)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tags := getTagsFromCvtArticles(cArts)
	arts := convertArticle(cArts, tags)
	arts = sortArticlesByDate(arts)
	tems := Tems{}
	tems.Title = bi.Title
	tems.Addr = bi.Addr
	tems.Articles = arts
	tems.Tags = tags

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

func convertArticle(cArts []cvt.Article, tags []Tag) []Article {
	arts := make([]Article, len(cArts))
	for i, ca := range cArts {
		for 
		for _, tag := range tags {
			
			arts[i] = Article{
				Filename: ca.Filename,
				Title:    ca.Title,
				Date:     ca.Date,
			}
		}
	}
	return arts
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
	newArt := make([]Article, 0)
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

func executeArticle(art Article, htmlDir string, templateDir string, outDir string) (ok bool) {
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

func getTagsFromCvtArticles(arts []cvt.Article) []Tag {
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

func sortArticlesByDate(arts []Article) []Article {
	var tmp Article
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
