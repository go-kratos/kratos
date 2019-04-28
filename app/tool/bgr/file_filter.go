package main

import (
	"os"
	"strings"
)

// FilterOnlyCodeFile .
var FilterOnlyCodeFile = func(f os.FileInfo) bool {
	if strings.HasSuffix(f.Name(), "_test.go") {
		return false
	}
	if strings.HasSuffix(f.Name(), ".go") {
		return true
	}
	return false
}

// FilterOnlyTestFile .
var FilterOnlyTestFile = func(f os.FileInfo) bool {
	return strings.HasSuffix(f.Name(), "_test.go")
}
