package gen

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/wirekang/cfg"
	"github.com/wirekang/errutil"
	"github.com/wirekang/fileutil"
)

var (
	Sep      = "##blogen##"
	articles []Article
	tags     []Tag
)

type Article struct {
	ID    string
	Title string
	Tags  string
	Date  string
	Time  time.Time
	MD    string
	HTML  string
}

type Tag struct {
	ID    int
	Name  string
	Count int
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

	err = parseTag(config.Find("tags").StringArray())
	es.Push(err)

	aid, err := getID(filename)
	es.Push(err)

	if isHashed(aid, mdString, hashDir) {
		return es.First()
	}

	if es.First() != nil {
		return es.First()
	}
	es.Clear()

	err = writeHash(aid, mdString, hashDir)
	es.Push(err)

	err = writeHTML(aid, mdString, htmlDir)
	es.Push(err)

	return es.First()
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
				Count: 0,
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
