package main

import (
	"flag"
	"os"
)

func main() {
	meta := flag.Bool("meta", false, "print meta")
	src := flag.String("src", "", "source directory")
	flag.Parse()

	if *meta {
		println("title: 제목\ntags: 태그1, 태그2, 태그3\ndate: 2020.8.12\n##blogen##\n")
		os.Exit(0)
	}

	if *src == "" {
		println("empty src")
		os.Exit(1)
	}

}
