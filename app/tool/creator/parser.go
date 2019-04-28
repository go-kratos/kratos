package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

type param struct{ K, V, P string }

type parse struct {
	Path    string
	Package string
	Imports []string
	Funcs   []*struct {
		Name                   string
		Method, Params, Result []*param
	}
}

func parseArgs(args []string, res *[]string, index int) (err error) {
	if len(args) <= index {
		return
	}
	if strings.HasPrefix(args[index], "-") {
		index += 2
		parseArgs(args, res, index)
		return
	}
	var f os.FileInfo
	if f, err = os.Stat(args[index]); err != nil {
		return
	}
	if f.IsDir() {
		if !strings.HasSuffix(args[index], "/") {
			args[index] += "/"
		}
		var fs []os.FileInfo
		if fs, err = ioutil.ReadDir(args[index]); err != nil {
			return
		}
		for _, f = range fs {
			args = append(args, args[index]+f.Name())
		}
	} else {
		if strings.HasSuffix(args[index], ".go") &&
			!strings.HasSuffix(args[index], "_test.go") {
			*res = append(*res, args[index])
		}
	}
	index++
	return parseArgs(args, res, index)
}

func parseFile(files ...string) (parses []*parse, err error) {
	for _, file := range files {
		var (
			astFile *ast.File
			fSet    = token.NewFileSet()
			parse   = &parse{}
		)
		if astFile, err = parser.ParseFile(fSet, file, nil, 0); err != nil {
			return
		}
		if astFile.Name != nil {
			parse.Path = file
			parse.Package = astFile.Name.Name
		}
		for _, decl := range astFile.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.(*ast.GenDecl).Specs {
					switch spec.(type) {
					case *ast.ImportSpec:
						parse.Imports = append(parse.Imports, spec.(*ast.ImportSpec).Path.Value)
					}
				}
			case *ast.FuncDecl:
				var (
					dec       = decl.(*ast.FuncDecl)
					parseFunc = &struct {
						Name                   string
						Method, Params, Result []*param
					}{Name: dec.Name.Name}
				)
				if dec.Recv != nil {
					parseFunc.Method = parserParams(dec.Recv.List)
				}
				if dec.Type.Params != nil {
					parseFunc.Params = parserParams(dec.Type.Params.List)
				}
				if dec.Type.Results != nil {
					parseFunc.Result = parserParams(dec.Type.Results.List)
				}
				parse.Funcs = append(parse.Funcs, parseFunc)
			}
		}
		parses = append(parses, parse)
	}
	return
}

func parserParams(fields []*ast.Field) (params []*param) {
	for _, field := range fields {
		p := &param{}
		//TODO:调用parseType解析类型
		p.V = parseType(field.Type)
		if field.Names == nil {
			params = append(params, p)
		}
		for _, name := range field.Names {
			sp := &param{}
			sp.K = name.Name
			if sp.K == "t" {
				sp.K = "no"
			}
			sp.V = p.V
			sp.P = p.P
			params = append(params, sp)
		}
	}
	return
}

func parseType(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.Ident:
		return expr.(*ast.Ident).Name
	case *ast.StarExpr:
		return "*" + parseType(expr.(*ast.StarExpr).X)
	case *ast.ArrayType:
		return "[" + parseType(expr.(*ast.ArrayType).Len) + "]" + parseType(expr.(*ast.ArrayType).Elt)
	case *ast.SelectorExpr:
		return parseType(expr.(*ast.SelectorExpr).X) + "." + expr.(*ast.SelectorExpr).Sel.Name
	case *ast.MapType:
		return "map[" + parseType(expr.(*ast.MapType).Key) + "]" + parseType(expr.(*ast.MapType).Value)
	case *ast.StructType:
		return "struct{}"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		var (
			pTemp string
			rTemp string
		)
		pTemp = parseFuncType(pTemp, expr.(*ast.FuncType).Params)
		if expr.(*ast.FuncType).Results != nil {
			rTemp = parseFuncType(rTemp, expr.(*ast.FuncType).Results)
			return fmt.Sprintf("func(%s) (%s)", pTemp, rTemp)
		}
		return fmt.Sprintf("func(%s)", pTemp)
	case *ast.ChanType:
		return fmt.Sprintf("make(chan %s)", parseType(expr.(*ast.ChanType).Value))
	case *ast.Ellipsis:
		return parseType(expr.(*ast.Ellipsis).Elt)
	}
	return ""
}

func parseFuncType(temp string, data *ast.FieldList) string {
	var params = parserParams(data.List)
	for i, param := range params {
		if i == 0 {
			temp = param.K + " " + param.V
			continue
		}
		t := param.K + " " + param.V
		temp = fmt.Sprintf("%s, %s", temp, t)
	}
	return temp
}
