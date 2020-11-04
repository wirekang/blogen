// Package blogen controls the main data flow.
package blogen

import (
	"bytes"
	"crypto/md5"
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

// Checksum returns md5 hash of given bytes.
func Checksum(b []byte) ([]byte, error) {
	hasher := md5.New()
	_, err := hasher.Write(b)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// IsModified returns true if given md bytes are newer than saved.
func IsModified(md []byte) (bool, error) {
	saved, err := ioutil.ReadFile(PathOfChecksum(aid))
	if err != nil {
		return true, nil
	}
	sum, err := Checksum(md)
	if err != nil {
		return true, err
	}
	return !bytes.Equal(saved, sum), nil
}

// PathOfMarkdown returns md/{aid}.md
func PathOfMarkdown(aid string) string {
	return path.Join(mdDir, aid+".md")
}

// PathOfChecksum returns cache/{aid}
func PathOfChecksum(aid string) string {
	return path.Join(checksumDir, aid)
}
