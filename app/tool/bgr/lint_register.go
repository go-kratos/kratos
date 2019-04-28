package main

import (
	"fmt"
	"go/ast"
	"strings"
)

var (
	_parsers = make(map[string]func(curDir string, f *ast.File, n ast.Node) (v string, hit bool))
)

func init() {
	// 函数定义注释
	_parsers["decl.func.comment"] = func(curDir string, f *ast.File, n ast.Node) (v string, hit bool) {
		var ele *ast.FuncDecl
		if ele, hit = n.(*ast.FuncDecl); !hit {
			return
		}
		v = ele.Doc.Text()
		return
	}

	// 标签语句
	_parsers["stmt.label"] = func(curDir string, f *ast.File, n ast.Node) (v string, hit bool) {
		var ele *ast.LabeledStmt
		if ele, hit = n.(*ast.LabeledStmt); !hit {
			return
		}
		v = ele.Label.Name
		return
	}

	// 调用函数的定义注释
	_parsers["expr.call.decl.comment"] = func(curDir string, f *ast.File, n ast.Node) (v string, hit bool) {
		var ele *ast.CallExpr
		if ele, hit = n.(*ast.CallExpr); !hit {
			return
		}
		var (
			files    []*ast.File
			declName string
			err      error
			doc      *ast.CommentGroup
		)
		switch ele2 := ele.Fun.(type) {
		case *ast.SelectorExpr: // like a.b()
			for _, impt := range f.Imports {
				if fmt.Sprintf("%s", ele2.X) == importName(impt) {
					var imptPath = strings.Trim(impt.Path.Value, `"`)
					if files, _, err = parsePackageFilesByPath(curDir, imptPath); err != nil {
						_log.Debugf("parsePackageFiles err: %+v", err)
						continue
					}
					declName = ele2.Sel.Name
				}
			}
		case *ast.Ident: // like a()
			if files, _, err = parsePackageFiles(curDir); err != nil {
				_log.Debugf("parsePackageFiles err: %+v", err)
				return
			}
			declName = ele2.Name
		}
		doc = declDocFromFiles(files, declName)
		if doc != nil {
			v = doc.Text()
		}
		return
	}

	// struct定义方法列表
	// _parsers["decl.gen.decl.method"] = func(curDir string, f *ast.File, n ast.Node) (v string, hit bool) {
	// 	var ele *ast.GenDecl
	// 	if ele, hit = n.(*ast.GenDecl); !hit {
	// 		return
	// 	}
	// 	switch ele.Tok {
	// 	case token.DEFINE, token.VAR:
	// 		for _, spec := range ele.Specs {
	// 			ele2, ok := spec.(*ast.ValueSpec)
	// 			if !ok {
	// 				_log.Debugf("spec: %T %+v can't convert to *ast.ValueSpec")
	// 				continue
	// 			}
	// 			switch ele3 := ele2.Type.(type) {
	// 			case *ast.CallExpr:
	// 				ele3.Fun.()
	// 			}
	// 		}
	// 	default:
	// 		hit = false
	// 		return
	// 	}

	// }
}

func importName(impt *ast.ImportSpec) string {
	if impt.Name != nil {
		return impt.Name.Name
	}
	sep := strings.Split(strings.Trim(impt.Path.Value, `"`), "/")
	return sep[len(sep)-1]
}

func declDocFromFiles(files []*ast.File, declName string) (doc *ast.CommentGroup) {
	for _, f := range files {
		for _, fd := range f.Decls {
			switch decl := fd.(type) {
			case *ast.FuncDecl: // like func a(){}
				if decl.Name.Name == declName {
					return decl.Doc
				}
			case *ast.GenDecl: // like var a int = 0
				for _, s := range decl.Specs {
					switch spec := s.(type) {
					case *ast.ValueSpec:
						for _, specName := range spec.Names {
							if specName.Name == declName {
								return spec.Doc
							}
						}
					}
				}
			default:
				_log.Debugf("decl(%+v,%s) unknown fDecl.(type): %T %+v", files, declName, decl, decl)
			}
		}
	}
	return
}
