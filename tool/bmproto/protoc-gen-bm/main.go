package main

import (
	"flag"
	"fmt"
	"kratos/tool/bmproto/pkg/generator"
	bmgen "kratos/tool/bmproto/protoc-gen-bm/generator"
	"kratos/tool/bmproto/protoc-gen-liverpc/gen"
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
