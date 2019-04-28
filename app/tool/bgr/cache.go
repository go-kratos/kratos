package main

import (
	"go/ast"
	"go/token"
)

var (
	_packageCache = make(map[string]*pkgC)
)

type pkgC struct {
	files []*ast.File
	fset  *token.FileSet
}

func packageCache(dir string) (files []*ast.File, fset *token.FileSet, ok bool) {
	c, ok := _packageCache[dir]
	if ok {
		files = c.files
		fset = c.fset
	}
	return
}

func setPackageCache(dir string, fset *token.FileSet, files []*ast.File) {
	_packageCache[dir] = &pkgC{fset: fset, files: files}
}
