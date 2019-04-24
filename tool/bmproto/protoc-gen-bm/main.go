package main

import (
	"flag"
	"fmt"
	"github.com/bilibili/kratos/tool/bmproto/pkg/generator"
	bmgen "github.com/bilibili/kratos/tool/bmproto/protoc-gen-bm/generator"
	"github.com/bilibili/kratos/tool/bmproto/pkg/gen"
	"os"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(generator.Version)
		os.Exit(0)
	}

	g := bmgen.BmGenerator()
	gen.Main(g)
}
