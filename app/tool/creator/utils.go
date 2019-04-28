package main

import (
	"fmt"
	"strings"

	"golang.org/x/tools/imports"
)

// GoImport Use golang.org/x/tools/imports auto import pkg
func GoImport(file string, bytes []byte) (res []byte, err error) {
	options := &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}
	if res, err = imports.Process(file, bytes, options); err != nil {
		fmt.Printf("GoImport(%s) error(%v)", file, err)
		res = bytes
		return
	}
	return
}

// IsService checkout the file belongs to service  or not
func IsService(pName string) bool {
	return pName == "service"
}

//ConvertHump convert words to hump style
func ConvertHump(words string) string {
	return strings.ToUpper(words[0:1]) + words[1:]
}
