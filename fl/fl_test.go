package fl

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	os.Mkdir("tmp", 0755)
	os.Chdir("tmp")
	if IsNecessaryDirsExist() {
		t.FailNow()
	}
	if IsFilesExist() {
		t.FailNow()
	}
	err := MakeDirs()
	if err != nil {
		t.Fatal(err)
	}
	err = CreateExampleFiles()
	if err != nil {
		t.Fatal(err)
	}
	if !IsNecessaryDirsExist() {
		t.FailNow()
	}
	if !IsFilesExist() {
		t.FailNow()
	}
	t.Log(Config())
	t.Log(Base())
	t.Log(List())
	t.Log(Single())
	t.Log(Style())
	t.Log(MarkdownDir())
	t.Log(HTMLDir())
	t.Log(TemplateDir())
	t.Log(HashDir())
	t.Log(OutDir())
	os.Chdir("..")
	os.RemoveAll("tmp")
}
