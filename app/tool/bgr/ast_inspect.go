package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	_defaultFileFilter = FilterOnlyCodeFile

	_errors = make([]string, 0)
	_warns  = make([]string, 0)
)

// AstInspect .
func AstInspect(dir string) (err error) {
	var (
		files []*ast.File
		fset  *token.FileSet
	)
	if dir, err = filepath.Abs(dir); err != nil {
		err = errors.WithStack(err)
		return
	}
	if files, fset, err = parsePackageFiles(dir); err != nil {
		return
	}
	for _, f := range files {
		handler := defaultNodeHandler(f, dir, fset)
		ast.Inspect(f, handler)
	}
	return
}

func defaultNodeHandler(f *ast.File, dir string, fileset *token.FileSet) func(ast.Node) bool {
	return func(node ast.Node) bool {
		if node == nil {
			return false
		}
		for _, l := range _lints {
			if !l.fn(dir, f, node) {
				switch l.s.l {
				case "e":
					_errors = append(_errors, fmt.Sprintf("%s --> %s", fileset.PositionFor(node.Pos(), true), l.s.d))
				case "w":
					_warns = append(_warns, fmt.Sprintf("%s --> %s", fileset.PositionFor(node.Pos(), true), l.s.d))
				}
			}
		}
		return true
	}
}

func parsePackageFilesByPath(dir string, importPath string) (files []*ast.File, fset *token.FileSet, err error) {
	var pkg *build.Package
	if pkg, err = build.Import(importPath, dir, build.FindOnly); err != nil {
		return
	}
	return parsePackageFiles(pkg.Dir)
}

func parsePackageFiles(absDir string) (files []*ast.File, fset *token.FileSet, err error) {
	var ok bool
	if files, fset, ok = packageCache(absDir); ok {
		return
	}
	defer func() {
		setPackageCache(absDir, fset, files)
	}()

	fset = token.NewFileSet()
	var (
		pkgMap map[string]*ast.Package
	)
	if pkgMap, err = parser.ParseDir(fset, absDir, _defaultFileFilter, parser.ParseComments); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, v := range pkgMap {
		for _, f := range v.Files {
			files = append(files, f)
		}
	}
	return
}
