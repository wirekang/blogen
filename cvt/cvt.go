// Package cvt converts markdown files to html.
package cvt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/wirekang/blogen/fl"
)

// Article is structure for article
type Article struct {
	Filename string
	Title    string
	Tags     []string
	Date     time.Time
}

// ConvertFiles read markdown files from directory src and parse them to html.
// The result will be written to directory dst.
//
// Written articles are returned.
func ConvertFiles(srcDir string, dstDir string, metaSep string, dateSep string) []Article {
	infos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var r *regexp.Regexp
	r, err = regexp.Compile("md$")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = fl.MakeIfNotExist(dstDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	articles := make([]Article, 0, len(infos))
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if !r.MatchString(info.Name()) {
			continue
		}
		filename := strings.Split(info.Name(), ".")[0]
		fmt.Printf("%s ... ", filename)
		ok, art := convertFile(path.Join(srcDir, info.Name()), path.Join(dstDir, filename), metaSep, dateSep)
		if !ok {
			continue
		}
		fmt.Println("Ok")
		art.Filename = filename
		articles = append(articles, art)
	}
	return articles
}

func convertFile(src string, dst string, metaSep string, dateSep string) (ok bool, art Article) {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println("Reading file failed.")
		return
	}
	ok, meta, md := seperateMeta(string(bytes), metaSep)
	if !ok {
		fmt.Println("Seperating meta failed.")
		return
	}
	var article Article
	ok, article = metaToArticle(meta, dateSep)
	if !ok {
		fmt.Println("Parsing meta failed.")
		return
	}
	html := parseMarkdown(md)
	err = ioutil.WriteFile(dst, html, 0644)
	if err != nil {
		fmt.Println(err)
		return false, article
	}
	return true, article
}

func seperateMeta(str string, metaSep string) (ok bool, meta string, md string) {
	a := strings.Split(str, metaSep)
	if len(a) != 2 {
		return false, "", ""
	}
	return true, a[0], a[1]
}

func metaToArticle(meta string, dateSep string) (ok bool, ar Article) {
	lines := strings.Split(meta, "\n")
	result := Article{}
	if len(lines) < 3 {
		return false, result
	}

	for _, line := range lines {
		t := strings.Split(line, ":")

		if len(t) != 2 {
			continue
		}

		t[0] = strings.TrimSpace(t[0])
		t[1] = strings.TrimSpace(t[1])
		switch t[0] {
		case "title":
			result.Title = t[1]
		case "tags":
			tags := strings.Split(t[1], ",")
			result.Tags = make([]string, len(tags))
			for i, t := range tags {
				tags[i] = strings.TrimSpace(t)
			}
			copy(result.Tags, tags)
		case "date":
			t, err := parseDate(t[1], dateSep)
			if err != nil {
				fmt.Println(err)
				continue
			}
			result.Date = t
		}
	}
	return true, result
}

func parseMarkdown(md string) (html []byte) {
	return markdown.ToHTML([]byte(md), nil, nil)
}

func parseDate(str string, sep string) (time.Time, error) {
	ds := strings.Split(str, sep)
	result := time.Time{}
	if len(ds) != 3 {
		return result, errors.New("date parsing error")
	}
	year, err := strconv.Atoi(ds[0])
	if err != nil {
		return result, err
	}
	month, err := strconv.Atoi(ds[1])
	if err != nil {
		return result, err
	}
	day, err := strconv.Atoi(ds[2])
	if err != nil {
		return result, err
	}
	result = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return result, nil
}
