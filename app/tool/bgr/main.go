package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"go-common/app/tool/bgr/log"
)

var (
	_flagType   string
	_flagScript string
	_flagDebug  bool
	_flagHit    string

	_log *log.Logger
)

func init() {
	flag.StringVar(&_flagType, "type", "file", "args type, file or dir")
	flag.StringVar(&_flagScript, "script", defaultDir(), "input script dir")
	flag.BoolVar(&_flagDebug, "debug", false, "set true, if need print debug info")
	flag.StringVar(&_flagHit, "hit", "", "filter hit key")
	flag.Parse()

	_log = log.New(os.Stdout, _flagDebug)
}

func defaultDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func main() {
	targets := flag.Args()
	switch _flagType {
	case "file":
		targets = filterFiles(targets)
		targets = combineDirs(targets)
	}

	_log.Debugf("check targets: %+v", targets)
	walkScript(_flagScript)

	for _, dir := range targets {
		if strings.HasSuffix(dir, "...") {
			walkDir(strings.TrimRight(dir, "..."))
		} else {
			if err := AstInspect(dir); err != nil {
				_log.Fatalf("%+v", err)
			}
		}
	}

	for _, desc := range _warns {
		_log.Warn(desc)
	}
	for _, desc := range _errors {
		_log.Error(desc)
	}
	if len(_errors) > 0 {
		os.Exit(1)
	}
}

func walkDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if err := AstInspect(path); err != nil {
				_log.Fatalf("%+v", err)
			}
		}
		return nil
	})
}

func combineDirs(files []string) (fs []string) {
	fmap := make(map[string]struct{})
	for _, f := range files {
		index := strings.LastIndex(f, "/")
		if index > 0 {
			fmap[f[:index]] = struct{}{}
		}
	}
	for k := range fmap {
		fs = append(fs, k)
	}
	return
}

func filterFiles(files []string) (fs []string) {
	for _, f := range files {
		if strings.Contains(f, _flagHit) {
			fs = append(fs, f)
		}
	}
	return
}
