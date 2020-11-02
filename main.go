package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/wirekang/blogen/cvt"
	"github.com/wirekang/blogen/gen"
)

func main() {
	meta := flag.Bool("meta", false, "print meta")
	mdsDir := flag.String("mds", "mds", "source directory")
	templatesDir := flag.String("temp", "templates", "templates directory")
	outDir := flag.String("out", "out", "output directory")
	flag.Parse()

	if *meta {
		println("title: Title of Article\ntags: tag1, tag2, tag3\ndate: 2020-8-12\n##blogen##\n")
		os.Exit(0)
	}
	tmp, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmp)

	os.RemoveAll(*outDir)
	arts := cvt.ConvertFiles(*mdsDir, tmp, "##blogen##", "-")
	ok := gen.GenerateFromTemplate(gen.BaseInfo{Title: "wirekang 블로그", Addr: "//127.0.0.1:5500/example/out"}, arts, tmp, *templatesDir, *outDir)
	if !ok {
		println("Failed.")
		os.Exit(1)
	}
	println("Done.")
}
