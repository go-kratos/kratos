package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bilibili/kratos/tool/bmproto/pkg/generator"
	"github.com/bilibili/kratos/tool/bmproto/pkg/gen"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(generator.Version)
		os.Exit(0)
	}

	g := NewSwaggerGenerator()
	gen.Main(g)
}
