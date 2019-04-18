package common

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// Source source
type Source struct {
	Fset *token.FileSet
	Src  string
	F    *ast.File
}

// NewSource new source
func NewSource(src string) *Source {
	s := &Source{
		Fset: token.NewFileSet(),
		Src:  src,
	}
	f, err := parser.ParseFile(s.Fset, "", src, 0)
	if err != nil {
		log.Fatal("无法解析源文件")
	}
	s.F = f
	return s
}

// ExprString expr string
func (s *Source) ExprString(typ ast.Expr) string {
	fset := s.Fset
	s1 := fset.Position(typ.Pos()).Offset
	s2 := fset.Position(typ.End()).Offset
	return s.Src[s1:s2]
}

// pkgPath package path
func (s *Source) pkgPath(name string) (res string) {
	for _, im := range s.F.Imports {
		if im.Name != nil && im.Name.Name == name {
			return im.Path.Value
		}
	}
	for _, im := range s.F.Imports {
		if strings.HasSuffix(im.Path.Value, name+"\"") {
			return im.Path.Value
		}
	}
	return
}

// GetDef get define code
func (s *Source) GetDef(name string) string {
	c := s.F.Scope.Lookup(name).Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType)
	s1 := s.Fset.Position(c.Pos()).Offset
	s2 := s.Fset.Position(c.End()).Offset
	line := s.Fset.Position(c.Pos()).Line
	lines := []string{strings.Split(s.Src, "\n")[line-1]}
	for _, l := range strings.Split(s.Src[s1:s2], "\n")[1:] {
		lines = append(lines, "\t"+l)
	}
	return strings.Join(lines, "\n")
}

// RegexpReplace replace regexp
func RegexpReplace(reg, src, temp string) string {
	result := []byte{}
	pattern := regexp.MustCompile(reg)
	for _, submatches := range pattern.FindAllStringSubmatchIndex(src, -1) {
		result = pattern.ExpandString(result, temp, src, submatches)
	}
	return string(result)
}

// formatPackage format package
func formatPackage(name, path string) (res string) {
	if path != "" {
		if strings.HasSuffix(path, name+"\"") {
			res = path
			return
		}
		res = fmt.Sprintf("%s %s", name, path)
	}
	return
}

// SourceText get source file text
func SourceText() string {
	file := os.Getenv("GOFILE")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("can't open file", file)
	}
	return string(data)
}

// FormatCode format code
func FormatCode(source string) string {
	src, err := format.Source([]byte(source))
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: 输出文件不合法: %s", err)
		log.Printf("warning: 详细错误请编译查看")
		return source
	}
	return string(src)
}

// Packages get import packages
func (s *Source) Packages(f *ast.Field) (res []string) {
	fs := f.Type.(*ast.FuncType).Params.List
	fs = append(fs, f.Type.(*ast.FuncType).Results.List...)
	var types []string
	resMap := make(map[string]bool)
	for _, field := range fs {
		if p, ok := field.Type.(*ast.MapType); ok {
			types = append(types, s.ExprString(p.Key))
			types = append(types, s.ExprString(p.Value))
		} else if p, ok := field.Type.(*ast.ArrayType); ok {
			types = append(types, s.ExprString(p.Elt))
		} else {
			types = append(types, s.ExprString(field.Type))
		}
	}

	for _, t := range types {
		name := RegexpReplace(`(?P<pkg>\w+)\.\w+`, t, "$pkg")
		if name == "" {
			continue
		}
		pkg := formatPackage(name, s.pkgPath(name))
		if !resMap[pkg] {
			resMap[pkg] = true
		}
	}
	for pkg := range resMap {
		res = append(res, pkg)
	}
	return
}
