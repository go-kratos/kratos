package main

import (
	"io/ioutil"
	"regexp"
)

var (
	_BmPattern = `(\w+)\s*(:?=)\s*bm\.Default\(\)((?:\n.*?)+)\s*\w+\.Serve\(.*?\,\s*c\..*?\);\s*(\w+)\s*!=\s*nil\s*{(\n*)(.*)\.Error\(.*\,\s*\w+\)`
	_BmReplace = `${1} ${2} bm.DefaultServer(c.BM)${3} ${1}.Start(); ${4} != nil {${5}${6}.Error("bm.DefaultServer error(%v)", ${4})`
	// _BmConfPattern = `type\sConfig\sstruct\s*{((?:.*?\n)+?)(?:(?:\s*BM\s+\*\w+\n((?:.*?\n)+?)^})|(?:(.*?)^}))`
	// _BmConfReplace = `type Config struct {\{1}\tBM *bm.ServerConfig\n\{2}}`
	// _confPath      = "/../conf/conf.go"
)

func upBladeMaster(files []string) (err error) {
	for _, file := range files {
		var bs []byte
		if bs, err = ioutil.ReadFile(file); err != nil {
			return
		}
		var reg *regexp.Regexp
		if reg, err = regexp.Compile(_BmPattern); err != nil {
			return
		}
		if !reg.Match(bs) {
			continue
		}
		bs = reg.ReplaceAll(bs, []byte(_BmReplace))
		if err = ioutil.WriteFile(file, bs, 0644); err != nil {
			return
		}
	}
	return
}
