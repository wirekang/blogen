package main

import (
	"flag"
	"os"

	"github.com/wirekang/blogen/cvt"
	"github.com/wirekang/blogen/gen"
)

func main() {
	meta := flag.Bool("meta", false, "print meta")
	mds := flag.String("mds", "mds", "source directory")
	temp := flag.String("temp", "templates", "templates directory")
	out := flag.String("out", "out", "output directory")
	flag.Parse()

	if *meta {
		println("title: 제목\ntags: 태그1, 태그2, 태그3\ndate: 2020-8-12\n##blogen##\n")
		os.Exit(0)
	}

	arts := cvt.ConvertFiles(*mds, "htmls", "##blogen##", "-")
	ok := gen.GenerateFromTemplate(gen.BaseInfo{Title: "wirekang 블로그", Addr: "localhost"}, arts, "htmls", *temp, *out)
	if !ok {
		println("Failed.")
		os.Exit(1)
	}
	println("Done.")
}
