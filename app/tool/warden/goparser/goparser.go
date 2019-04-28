package goparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"go-common/app/tool/warden/types"
)

var protoFileRegexp *regexp.Regexp

const (
	optionsPrefix = "+wd:"
)

func init() {
	protoFileRegexp = regexp.MustCompile(`//\s+source:\s+(.*\.proto)`)
}

// GoPackage get go package name from file or directory path
func GoPackage(dpath string) (string, error) {
	if strings.HasSuffix(dpath, ".go") {
		dpath = filepath.Dir(dpath)
	}
	absDir, err := filepath.Abs(dpath)
	if err != nil {
		return "", err
	}
	goPaths := os.Getenv("GOPATH")
	if goPaths == "" {
		return "", fmt.Errorf("GOPATH not set")
	}
	for _, goPath := range strings.Split(goPaths, ":") {
		srcPath := path.Join(goPath, "src")
		if !strings.HasPrefix(absDir, srcPath) {
			continue
		}
		return strings.Trim(absDir[len(srcPath):], "/"), nil
	}
	return "", fmt.Errorf("give package not under $GOPATH")
}

// Parse service spec with gived path and receiver name
func Parse(name, dpath, recvName, workDir string) (*types.ServiceSpec, error) {
	if workDir == "" {
		workDir, _ = os.Getwd()
	}
	ps := &parseState{
		name:     strings.Title(name),
		dpath:    dpath,
		recvName: recvName,
		workDir:  workDir,
	}
	return ps.parse()
}

type parseState struct {
	dpath    string
	recvName string
	name     string
	workDir  string

	typedb      map[string]types.Typer
	importPath  string
	packageName string
	methods     []*types.Method
}

func (p *parseState) parse() (spec *types.ServiceSpec, err error) {
	p.typedb = make(map[string]types.Typer)
	if p.importPath, err = GoPackage(p.dpath); err != nil {
		return
	}
	if err := p.searchMethods(); err != nil {
		return nil, err
	}
	return &types.ServiceSpec{
		ImportPath: p.importPath,
		Name:       p.name,
		Package:    p.packageName,
		Receiver:   p.recvName,
		Methods:    p.methods,
	}, nil
}

func (p *parseState) searchMethods() error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, p.dpath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("no package found on %s", p.dpath)
	}
	if len(pkgs) > 1 {
		return fmt.Errorf("multiple package found on %s", p.dpath)
	}
	for pkgName, pkg := range pkgs {
		//log.Printf("search method in package %s", pkgName)
		p.packageName = pkgName
		for fn, f := range pkg.Files {
			//log.Printf("search method in file %s", fn)
			if err = p.searchMethodsInFile(pkg, f); err != nil {
				log.Printf("search method in %s err %s", fn, err)
			}
		}
	}
	return nil
}

func (p *parseState) searchMethodsInFile(pkg *ast.Package, f *ast.File) error {
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || !funcDecl.Name.IsExported() || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}
		var recvIdent *ast.Ident
		recvField := funcDecl.Recv.List[0]
		switch rt := recvField.Type.(type) {
		case *ast.Ident:
			recvIdent = rt
		case *ast.StarExpr:
			recvIdent = rt.X.(*ast.Ident)
		}
		if recvIdent == nil {
			return fmt.Errorf("unknown recv %v", recvField)
		}
		if recvIdent.Name != p.recvName {
			continue
		}
		log.Printf("find method %s", funcDecl.Name.Name)
		if err := p.parseFuncDecl(pkg, f, funcDecl); err != nil {
			return err
		}
	}
	return nil
}

func (p *parseState) parseFuncDecl(pkg *ast.Package, f *ast.File, funcDecl *ast.FuncDecl) error {
	//log.Printf("parse method %s", funcDecl.Name.Name)
	comments, options := parseComments(funcDecl)

	for _, option := range options {
		if option == "ignore" {
			log.Printf("ignore method %s", funcDecl.Name.Name)
			return nil
		}
	}

	ps := typeState{
		File:       f,
		ImportPath: p.importPath,
		Pkg:        pkg,
		WorkDir:    p.workDir,
		typedb:     p.typedb,
		PkgDir:     p.dpath,
	}
	parameters, err := ps.parseFieldList(funcDecl.Type.Params, false)
	if err != nil {
		return err
	}
	results, err := ps.parseFieldList(funcDecl.Type.Results, false)
	if err != nil {
		return err
	}
	method := &types.Method{
		Name:       funcDecl.Name.Name,
		Comments:   comments,
		Options:    options,
		Parameters: parameters,
		Results:    results,
	}
	p.methods = append(p.methods, method)
	return nil
}

type typeState struct {
	typedb     map[string]types.Typer
	ImportPath string
	Pkg        *ast.Package
	File       *ast.File
	WorkDir    string
	PkgDir     string
}

func (t *typeState) parseType(expr ast.Expr, ident string) (types.Typer, error) {
	oldFile := t.File
	defer func() {
		t.File = oldFile
	}()
	switch exp := expr.(type) {
	case *ast.Ident:
		if isBuildIn(exp.Name) {
			return &types.BasicType{Name: exp.Name}, nil
		}
		tid := fmt.Sprintf("%s-%s-%s", t.ImportPath, t.Pkg.Name, exp.Name)
		if ty, ok := t.typedb[tid]; ok {
			return ty, nil
		}
		ty, err := t.searchType(exp)
		if err != nil {
			return nil, err
		}
		t.typedb[tid] = ty
		return ty, nil
	case *ast.StarExpr:
		t, err := t.parseType(exp.X, ident)
		if err != nil {
			return nil, err
		}
		return t.SetReference(), nil
	case *ast.SelectorExpr:
		return t.parseSel(exp)
	case *ast.ArrayType:
		et, err := t.parseType(exp.Elt, ident)
		if err != nil {
			return nil, err
		}
		return &types.ArrayType{EltType: et}, nil
	case *ast.MapType:
		kt, err := t.parseType(exp.Key, ident)
		if err != nil {
			return nil, err
		}
		vt, err := t.parseType(exp.Value, ident)
		if err != nil {
			return nil, err
		}
		return &types.MapType{KeyType: kt, ValueType: vt}, nil
	case *ast.InterfaceType:
		return &types.InterfaceType{
			ImportPath: t.ImportPath,
			Package:    t.Pkg.Name,
			IdentName:  ident,
		}, nil
	case *ast.StructType:
		fields, err := t.parseFieldList(exp.Fields, true)
		return &types.StructType{
			IdentName:  ident,
			ImportPath: t.ImportPath,
			Package:    t.Pkg.Name,
			Fields:     fields,
			ProtoFile:  findProtoFile(t.PkgDir, t.File),
		}, err
	}
	return nil, fmt.Errorf("unexpect expr %v", expr)
}

func (t *typeState) searchType(ident *ast.Ident) (types.Typer, error) {
	//log.Printf("search type %s", ident.Name)
	for fn, f := range t.Pkg.Files {
		//log.Printf("search in %s", fn)
		for _, decl := range f.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						return nil, fmt.Errorf("expect typeSpec get %v in file %s", spec, fn)
					}
					if typeSpec.Name.Name == ident.Name {
						//log.Printf("found in %s", fn)
						t.File = f
						return t.parseType(typeSpec.Type, ident.Name)
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("type %s not found in package %s", ident.Name, t.Pkg.Name)
}

func lockType(pkg *ast.Package, ident *ast.Ident) (*ast.File, error) {
	//log.Printf("lock type %s", ident.Name)
	for fn, f := range pkg.Files {
		//log.Printf("search in %s", fn)
		for _, decl := range f.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						return nil, fmt.Errorf("expect typeSpec get %v in file %s fn", spec, fn)
					}
					if typeSpec.Name.Name == ident.Name {
						return f, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("type %s not found in package %s", ident.Name, pkg.Name)
}

func (t *typeState) parseFieldList(fl *ast.FieldList, filterExported bool) ([]*types.Field, error) {
	fields := make([]*types.Field, 0, fl.NumFields())
	if fl == nil {
		return fields, nil
	}
	for _, af := range fl.List {

		ty, err := t.parseType(af.Type, "")
		if err != nil {
			return nil, err
		}
		if af.Names == nil {
			fields = append(fields, &types.Field{Type: ty})
		} else {
			for _, name := range af.Names {
				if filterExported && !name.IsExported() {
					continue
				}
				fields = append(fields, &types.Field{Type: ty, Name: name.Name})
			}
		}
	}
	return fields, nil
}

func (t *typeState) parseSel(sel *ast.SelectorExpr) (types.Typer, error) {
	//log.Printf("parse sel %v.%v", sel.X, sel.Sel)
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("unsupport sel.X type %v", sel.X)
	}
	var pkg *ast.Package
	var pkgPath string
	var err error
	var importPath string
	var found bool
	var pkgs map[string]*ast.Package
	for _, spec := range t.File.Imports {
		importPath = strings.Trim(spec.Path.Value, "\"")

		if spec.Name != nil && spec.Name.Name == x.Name {
			pkgPath, err = importPackage(t.WorkDir, importPath)
			if err != nil {
				return nil, err
			}

			pkgs, err = parser.ParseDir(token.NewFileSet(), pkgPath, nil, parser.ParseComments)
			if err != nil {
				return nil, err
			}

			pkg, err = filterPkgs(pkgs)
			if err != nil {
				return nil, err
			}
			found = true
			break
		}

		pkgPath, err = importPackage(t.WorkDir, importPath)
		if err != nil {
			return nil, err
		}

		pkgs, err = parser.ParseDir(token.NewFileSet(), pkgPath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if pkg, ok = pkgs[x.Name]; ok {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("can't found type %s.%s", x.Name, sel.Sel.Name)
	}

	file, err := lockType(pkg, sel.Sel)
	if err != nil {
		return nil, err
	}
	ts := &typeState{
		File:       file,
		Pkg:        pkg,
		ImportPath: importPath,
		WorkDir:    t.WorkDir,
		typedb:     t.typedb,
		PkgDir:     pkgPath,
	}
	return ts.searchType(sel.Sel)
}

func filterPkgs(pkgs map[string]*ast.Package) (*ast.Package, error) {
	for pname, pkg := range pkgs {
		if strings.HasSuffix(pname, "_test") {
			continue
		}
		return pkg, nil
	}
	return nil, fmt.Errorf("no package found")
}

func importPackage(workDir, importPath string) (string, error) {
	//log.Printf("import package %s", importPath)
	searchPaths := make([]string, 0, 3)
	searchPaths = append(searchPaths, path.Join(runtime.GOROOT(), "src"))
	if vendorDir, ok := searchVendor(workDir); ok {
		searchPaths = append(searchPaths, vendorDir)
	}
	for _, goPath := range strings.Split(os.Getenv("GOPATH"), ":") {
		searchPaths = append(searchPaths, path.Join(goPath, "src"))
	}
	var pkgPath string
	var found bool
	for _, basePath := range searchPaths {
		pkgPath = path.Join(basePath, importPath)
		if stat, err := os.Stat(pkgPath); err == nil && stat.IsDir() {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("can't import package %s", importPath)
	}
	return pkgPath, nil
}

func searchVendor(workDir string) (vendorDir string, ok bool) {
	var err error
	if workDir, err = filepath.Abs(workDir); err != nil {
		return "", false
	}
	goPath := os.Getenv("GOPATH")
	for {
		if !strings.HasPrefix(workDir, goPath) {
			break
		}
		vendorDir := path.Join(workDir, "vendor")
		if stat, err := os.Stat(vendorDir); err == nil && stat.IsDir() {
			return vendorDir, true
		}
		workDir = filepath.Dir(workDir)
	}
	return
}

func parseComments(funcDecl *ast.FuncDecl) (comments []string, options []string) {
	if funcDecl.Doc == nil {
		return
	}
	for _, comment := range funcDecl.Doc.List {
		text := strings.TrimLeft(comment.Text, "/ ")
		if strings.HasPrefix(text, optionsPrefix) {
			options = append(options, text[len(optionsPrefix):])
		} else {
			comments = append(comments, text)
		}
	}
	return
}

func isBuildIn(t string) bool {
	switch t {
	case "bool", "byte", "complex128", "complex64", "error", "float32",
		"float64", "int", "int16", "int32", "int64", "int8",
		"rune", "string", "uint", "uint16", "uint32", "uint64", "uint8", "uintptr":
		return true
	}
	return false
}

func findProtoFile(pkgDir string, f *ast.File) string {
	if f.Comments == nil {
		return ""
	}
	for _, comment := range f.Comments {
		if comment.List == nil {
			continue
		}
		for _, line := range comment.List {
			if protoFile := extractProtoFile(line.Text); protoFile != "" {
				fixPath := path.Join(pkgDir, protoFile)
				if s, err := os.Stat(fixPath); err == nil && !s.IsDir() {
					return fixPath
				}
				return protoFile
			}
		}
	}
	return ""
}

func extractProtoFile(line string) string {
	matchs := protoFileRegexp.FindStringSubmatch(line)
	if len(matchs) > 1 {
		return matchs[1]
	}
	return ""
}
