// Package blogen controls the main data flow.
package blogen

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/wirekang/blogen/fl"
)

// Tag is keyword for article.
type Tag struct {
	ID    int
	Name  string
	Count int
}

// Article is a piece of writing.
type Article struct {
	ID    string
	Title string
	Time  time.Time
	Date  string
	Tags  []Tag
	HTML  string
}

const settingsFile = "settings.txt"
const mdDir = "md"
const templateDir = "template"
const templateBaseFile = "base.html"
const templateMainFile = "main.html"
const templateListFile = "list.html"
const templateSingleFile = "single.html"
const templateStyleFile = "style.css"

const checksumDir = "cache"
const outDir = "out"
const genDir = "gen"

// CheckDirs returns true if all necessary directories exist.
// If a directory doesn't exist, an empty directory is created.
func CheckDirs() bool {
	necessaryDirs := []string{mdDir, templateDir}
	ok := true
	for _, d := range necessaryDirs {
		if !fl.IsExist(d) {
			ok = false
			os.MkdirAll(d, 0755)
		}
	}
	return ok
}

// CheckFiles returns true if all necessary files exist.
// If a file doesn't exist, an empty file is created.
func CheckFiles() bool {
	necessaryFiles := []string{settingsFile,
		path.Join(templateDir, templateBaseFile),
		path.Join(templateDir, templateMainFile),
		path.Join(templateDir, templateListFile),
		path.Join(templateDir, templateSingleFile),
		path.Join(templateDir, templateStyleFile)}
	ok := true
	for _, f := range necessaryFiles {
		if !fl.IsExist(f) {
			ok = false
			f, _ := os.Create(f)
			f.Close()
		}
	}
	return ok
}

// IsModified returns true if given article file differs from the cached checksum.
func IsModified(aid string) (bool, error) {
	ch, err := ioutil.ReadFile(pathOfChecksum(aid))
	if err != nil {
		return true, nil
	}
	md, err := ioutil.ReadFile(pathOfMarkdown(aid))
	if err != nil {
		return true, err
	}
	return !bytes.Equal(ch, md), nil

}

func pathOfMarkdown(aid string) string {
	return path.Join(mdDir, aid+".md")
}

func pathOfChecksum(aid string) string {
	return path.Join(checksumDir, aid)
}
