package gen

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/wirekang/cfg"
	"github.com/wirekang/errutil"
	"github.com/wirekang/fileutil"
)

var (
	Sep   = "##blogen##"
	posts []Post
	tags  []Tag
)

type TemplateBase struct {
	Title string
	Addr  string
	Posts []Post
	Tags  []Tag

	//list
	InList       bool
	IsFiltered   bool
	FilteredTags []Tag

	//single
	InSingle     bool
	Post         Post
	HTML         template.HTML
	RelatedPosts []Post
}

type Post struct {
	ID    string
	Title string
	Tags  []Tag
	Time  time.Time
}

type Tag struct {
	ID    int
	Name  string
	Count int
}

// Generate generates static files.
func Generate(title string, addr string, templateDir string, htmlDir string, outDir string) error {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	tem, err := template.ParseFiles(path.Join(templateDir, "base.html"),
		path.Join(templateDir, "list.html"), path.Join(templateDir, "style.css"))
	if err != nil {
		return err
	}
	err = generateList(title, addr, tem, path.Join(outDir, "index.html"), nil)
	es := errutil.NewErrorStack(err)

	for _, tag := range tags {
		err = generateList(title, addr, tem, path.Join(outDir, fmt.Sprintf("tag%d.html", tag.ID)), []Tag{tag})
		es.Push(err)
	}

	if es.First() != nil {
		return es.First()
	}

	tem, err = template.ParseFiles(path.Join(templateDir, "base.html"), path.Join(templateDir, "single.html"),
		path.Join(templateDir, "style.css"))
	if err != nil {
		return err
	}

	for _, post := range posts {
		err = generateSingle(title, addr, tem, outDir, htmlDir, post)
		es.Push(err)
	}
	return es.First()
}

func generateList(title string, addr string, tem *template.Template, file string, filteredTags []Tag) error {
	templateBase := TemplateBase{
		Title:        title,
		Addr:         addr,
		Tags:         tags,
		InList:       true,
		IsFiltered:   filteredTags != nil,
		FilteredTags: filteredTags,
	}

	if filteredTags == nil {
		templateBase.Posts = posts
	} else {

		psts := make([]Post, 0)
		for _, pst := range posts {
			contain := false
		Loop:
			for _, ft := range filteredTags {
				for _, t := range pst.Tags {
					if ft.ID == t.ID {
						contain = true
						break Loop
					}
				}
			}
			if contain {
				psts = append(psts, pst)
			}
		}
		templateBase.Posts = psts
	}
	wr, err := os.Create(file)
	if err != nil {
		return err
	}
	return tem.Execute(wr, templateBase)
}

func generateSingle(title string, addr string, tem *template.Template, outDir string, htmlDir string, post Post) error {
	wr, err := os.Create(path.Join(outDir, fmt.Sprintf("%s.html", post.ID)))
	if err != nil {
		return err
	}
	html, err := ioutil.ReadFile(path.Join(htmlDir, post.ID))
	if err != nil {
		return err
	}

	templateBase := TemplateBase{
		Title:    title,
		Addr:     addr,
		Tags:     tags,
		Post:     post,
		InSingle: true,
		HTML:     template.HTML(html),
	}
	rel := make([]Post, 0)
	for _, pst := range posts {
	Loop:
		for _, t1 := range pst.Tags {
			for _, t2 := range post.Tags {
				if t1.ID == t2.ID {
					rel = append(rel, pst)
					break Loop
				}
			}
		}
	}
	templateBase.RelatedPosts = rel
	return tem.Execute(wr, templateBase)
}

// ParseMD get tags from md file and write html file to htmlDir.
func ParseMD(filename string, hashDir string, htmlDir string) error {
	src, err := ioutil.ReadFile(filename)
	es := errutil.NewErrorStack(err)

	configString, mdString, err := split(string(src))
	es.Push(err)

	config, err := cfg.Load(configString)
	es.Push(err)

	if !config.IsExist("tags") {
		es.Push(errors.New("no tags"))
	}

	if es.First() != nil {
		return es.First()
	}

	aid, err := getID(filename)
	if err != nil {
		return err
	}

	err = parseTag(config.Find("tags").StringArray())
	es.Push(err)

	post, err := parsePost(config)
	if err != nil {
		return err
	}
	post.ID = aid
	ts := make([]Tag, 0)
	for _, t := range config.Find("tags").StringArray() {
		tag, err := findTag(t)
		if err != nil {
			continue
		}
		ts = append(ts, tag)
	}
	post.Tags = ts
	posts = append(posts, post)

	if isHashed(aid, mdString, hashDir) {
		return nil
	}

	err = writeHash(aid, mdString, hashDir)
	es.Push(err)

	err = writeHTML(aid, mdString, htmlDir)
	es.Push(err)

	return es.First()
}

// parsePost parses title, date from config to Post.
func parsePost(con cfg.Config) (Post, error) {
	var err error
	post := Post{}
	post.Title = con.Find("title").String()
	post.Time, err = con.Find("date").Date()
	if err != nil {
		return post, err
	}
	return post, err
}

// split splits string by Sep.
func split(src string) (config string, md string, err error) {
	arr := strings.Split(src, Sep)
	if len(arr) < 2 {
		return "", "", errors.New("no seperator")
	}
	return arr[0], arr[1], nil
}

// parseTag parses tags from string.
func parseTag(src []string) error {
	if len(src) == 0 {
		return errors.New("no tags")
	}
	for _, str := range src {
		new := true
		for i, tag := range tags {
			if tag.Name == strings.TrimSpace(str) {
				new = false
				tags[i].Count++
				break
			}
		}
		if new {
			tags = append(tags, Tag{
				Name:  strings.TrimSpace(str),
				ID:    len(tags) + 1,
				Count: 1,
			})
		}
	}
	return nil
}

// // isHashed returns true if hash of md is already stored.
func isHashed(aid, md string, hashDir string) bool {
	hashFile := path.Join(hashDir, aid)
	if !fileutil.IsExist(hashFile) {
		return false
	}

	b, err := ioutil.ReadFile(hashFile)
	es := errutil.NewErrorStack(err)

	hash, err := getHash(md)
	es.Push(err)

	if es.First() != nil {
		return false
	}
	return bytes.Equal(b, hash)
}

// getID returns filename without .md extension.
func getID(filename string) (string, error) {
	b := path.Base(filename)
	if !strings.Contains(b, ".md") {
		return "", fmt.Errorf("%s is not .md file", b)
	}
	return strings.Replace(b, ".md", "", 1), nil
}

// writeHash writes hash of the md to hashDir/aid file.
func writeHash(aid string, md string, hashDir string) error {
	hashFile := path.Join(hashDir, aid)
	hash, err := getHash(md)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(hashFile, hash, 0755)
	return err
}

// getHash returns md5 hash of the string.
func getHash(str string) ([]byte, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(str))
	return hasher.Sum(nil), err
}

// writeHTML writes html file to htmlDir.
func writeHTML(aid string, md string, htmlDir string) error {
	file := path.Join(htmlDir, aid)
	err := ioutil.WriteFile(file, markdown.ToHTML([]byte(md), nil, nil), 0755)
	return err
}

// findTag returns tag.
func findTag(tag string) (Tag, error) {
	for _, t := range tags {
		if t.Name == tag {
			return t, nil
		}
	}
	return Tag{}, fmt.Errorf("no tag %s", tag)
}
