package cvt

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestConvertFiles(t *testing.T) {
	os.Chdir("..")
	arts := ConvertFiles("example/mds", "example/htmls", "##blogen##", ".")
	if arts == nil {
		t.FailNow()
	}
	for _, a := range arts {
		fmt.Println(a)
	}
}

func TestConvertFile(t *testing.T) {
	os.Chdir("..")
	ok, art := convertFile("example/mds/md1.md", "example/htmls/md1.html", "##blogen##", ".")
	if !ok {
		t.FailNow()
	}
	fmt.Println(art)
}

func TestParseDate(t *testing.T) {
	ti, e := parseDate("2028-12-03", "-")
	if e != nil {
		t.FailNow()
	}
	year, month, day := ti.Date()
	if year != 2028 || month != 12 || day != 3 {
		fmt.Println(ti)
		t.FailNow()
	}
}
func TestSeperateMeta(t *testing.T) {
	META := 100
	MD := 4000
	meta := make([]byte, META)
	md := make([]byte, MD)
	for i := 0; i < META; i++ {
		meta[i] = byte(rand.Intn(255))
	}
	for i := 0; i < MD; i++ {
		md[i] = byte(rand.Intn(255))
	}
	sep := "##blogen##"
	str := string(append(append(meta, sep...), md...))
	o, i, m := seperateMeta(str, "##blogen##")
	if !o {
		t.Fail()
	}
	if i != string(meta) {
		t.Fail()
	}
	if m != string(md) {
		t.Fail()
	}
}

func TestParseMeta(t *testing.T) {
	str := "title: 제목제목\ntags: 태그1, 태그2 \ndate:2018.02.02"

	ok, ar := metaToArticle(str, ".")
	if !ok {
		t.FailNow()
	}
	if ar.Title != "제목제목" {
		t.FailNow()
	}
	if len(ar.Tags) != 2 || ar.Tags[0] != "태그1" || ar.Tags[1] != "태그2" {
		t.FailNow()
	}
}

func TestParseMarkdown(t *testing.T) {
	md := "# H1\n## H2"
	html := "<h1>H1</h1>\n\n<h2>H2</h2>\n"
	if html != string(parseMarkdown(md)) {
		t.FailNow()
	}
}
