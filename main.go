package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/wirekang/blogen/fl"
	"github.com/wirekang/blogen/gen"
	"github.com/wirekang/cfg"
)

func main() {
	init := flag.Bool("i", false, "initialization")
	src := flag.String("src", "", "source directory")

	flag.Parse()

	if *src != "" {
		err := os.Chdir(*src)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if *init {
		err := fl.MakeDirs()
		if err != nil {
			fmt.Printf("Can't make dirs: %s", err)
			os.Exit(1)
		}
		err = fl.CreateExampleFiles()
		if err != nil {
			fmt.Printf("Can't create files: %s", err)
			os.Exit(1)
		}
		fmt.Println("Done")
		os.Exit(0)
	}
	if !fl.IsNecessaryDirsExist() || !fl.IsFilesExist() {
		fmt.Println("Can't find necessary files. Add -i option for initialization.")
		os.Exit(1)
	}
	con, err := cfg.LoadFile(fl.Config())
	if err != nil {
		fmt.Printf("Can't load %s: %s\n", fl.Config(), err)
		os.Exit(1)
	}
	if !con.IsExist("title") || !con.IsExist("addr") {
		fmt.Printf("Can't find necessary config.")
		os.Exit(1)
	}
	title := con.Find("title").String()
	addr := con.Find("addr").String()
	fmt.Printf("Start generating.\ntitle=%s\naddr=%s\n", title, addr)

	mds, err := filepath.Glob(path.Join(fl.MarkdownDir(), "*.md"))
	if err != nil {
		fmt.Printf("Can't read files in %s: %s", fl.MarkdownDir(), err)
		os.Exit(1)
	}
	for _, md := range mds {
		fmt.Printf("%s... ", md)
		err = gen.ParseMD(md, fl.HashDir(), fl.HTMLDir())
		if err != nil {
			fmt.Printf("Fail: %s", err)
			continue
		}
		fmt.Println("Success.")

	}
}
