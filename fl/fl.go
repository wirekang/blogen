// Package fl handle files.
package fl

import (
	"os"
	"path"

	"github.com/wirekang/fileutil"
)

type directory struct {
	path        string
	isNecessary bool
}

type file struct {
	path    string
	example string
}

var (
	dirs  = make([]directory, 0)
	files = make([]file, 0)
)

var (
	root       = directory{"", false}
	configFile = newFile(root, "blogen.cfg",
		"title= Title of blog\ndescription= Description\naddr= //example.com")
)

var mdDir = newDirectory(root, "md", true)

var (
	templateDir = newDirectory(root, "template", true)
	baseFile    = newFile(templateDir, "base.html", " ")
	listFile    = newFile(templateDir, "list.html", " ")
	singleFile  = newFile(templateDir, "single.html", "")
	styleFile   = newFile(templateDir, "style.css", "h1 {\n\n}")
)

var outDir = newDirectory(root, "out", false)

var (
	genDir  = newDirectory(root, "blogen-cache", false)
	hashDir = newDirectory(genDir, "hash", false)
	htmlDir = newDirectory(genDir, "html", false)
)

func newDirectory(parent directory, name string, isNecessary bool) directory {
	d := directory{
		path:        path.Join(parent.path, name),
		isNecessary: isNecessary,
	}
	dirs = append(dirs, d)
	return d
}

func newFile(parent directory, name string, example string) file {
	f := file{
		path:    path.Join(parent.path, name),
		example: example,
	}
	files = append(files, f)
	return f
}

// IsNecessaryDirsExist returns true if all necessary directories exist.
func IsNecessaryDirsExist() bool {
	for _, dir := range dirs {
		if !dir.isNecessary {
			continue
		}
		if !fileutil.IsExist(dir.path) {
			return false
		}
	}
	return true
}

// MakeDirs makes all directories.
func MakeDirs() error {
	for _, dir := range dirs {
		err := fileutil.MakeIfNotExist(dir.path)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsFilesExist returns true if all files exist.
func IsFilesExist() bool {
	for _, file := range files {
		if !fileutil.IsExist(file.path) {
			return false
		}
	}
	return true
}

// CreateExampleFiles creates all files with example content.
func CreateExampleFiles() error {
	for _, file := range files {
		if fileutil.IsExist(file.path) {
			continue
		}
		f, err := os.Create(file.path)
		if err != nil {
			return err
		}
		_, err = f.WriteString(file.example)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Config returns file path of config file.
func Config() string {
	return configFile.path
}

// Base returns file path of base.html.
func Base() string {
	return baseFile.path
}

// List returns file path of list.html.
func List() string {
	return listFile.path
}

// Single returns file path of single.html.
func Single() string {
	return singleFile.path
}

// Style returns file path of style.html.
func Style() string {
	return styleFile.path
}

// MarkdownDir returns path of markdown directory.
func MarkdownDir() string {
	return mdDir.path
}

// TemplateDir returns path of template directory.
func TemplateDir() string {
	return templateDir.path
}

// HashDir returns path of hash generated from markdown directory.
func HashDir() string {
	return hashDir.path
}

// HTMLDir returns path of html parsed from markdown directory.
func HTMLDir() string {
	return htmlDir.path
}

// OutDir returns path of output directory.
func OutDir() string {
	return outDir.path
}
