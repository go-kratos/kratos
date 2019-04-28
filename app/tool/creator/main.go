package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	err          error
	_mode, _func string
	files        []string
	parses       []*parse
)

func main() {
	flag.StringVar(&_mode, "m", "test", "Generating code by Working mode. [test|interface|mock|upgrade...]")
	flag.StringVar(&_func, "func", "", "Generating code by function.")
	flag.Parse()
	if len(os.Args) == 1 {
		println("Creater is a tool for generating code.\n\nUsage: creater [-m]")
		flag.PrintDefaults()
		return
	}
	if err = parseArgs(os.Args[1:], &files, 0); err != nil {
		panic(err)
	}
	switch _mode {
	case "monkey":
		if parses, err = parseFile(files...); err != nil {
			panic(err)
		}
		if err = genMonkey(parses); err != nil {
			panic(err)
		}
	case "test":
		if parses, err = parseFile(files...); err != nil {
			panic(err)
		}
		if err = genTest(parses); err != nil {
			panic(err)
		}
	case "interface":
		if parses, err = parseFile(files...); err != nil {
			panic(err)
		}
		if err = genInterface(parses); err != nil {
			panic(err)
		}
	case "mock":
		if err = genMock(files...); err != nil {
			panic(err)
		}
	case "upgrade":
		if err = upBladeMaster(files); err != nil {
			panic(err)
		}
	default:
	}
	fmt.Println(moha)
}
