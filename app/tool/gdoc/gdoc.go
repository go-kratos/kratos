package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// gloabl var.
var (
	ErrParams = errors.New("err params")
	_gopath   = filepath.SplitList(os.Getenv("GOPATH"))
)

var (
	dir         string
	pkgs        = make(map[string]*ast.Package)
	rlpkgs      = make(map[string]*ast.Package)
	definitions = make(map[string]*Schema)
	swagger     = Swagger{
		Definitions:    make(map[string]*Schema),
		Paths:          make(map[string]*Item),
		SwaggerVersion: "2.0",
		Infos: Information{
			Title:       "go-common api",
			Description: "api",
			Version:     "1.0",
			Contact: Contact{
				EMail: "lintanghui@bilibili.com",
			},
			License: &License{
				Name: "Apache 2.0",
				URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
	}
	stdlibObject = map[string]string{
		"&{time Time}": "time.Time",
	}
)

// refer to builtin.go
var basicTypes = map[string]string{
	"bool":       "boolean:",
	"uint":       "integer:int32",
	"uint8":      "integer:int32",
	"uint16":     "integer:int32",
	"uint32":     "integer:int32",
	"uint64":     "integer:int64",
	"int":        "integer:int64",
	"int8":       "integer:int32",
	"int16":      "integer:int32",
	"int32":      "integer:int32",
	"int64":      "integer:int64",
	"uintptr":    "integer:int64",
	"float32":    "number:float",
	"float64":    "number:double",
	"string":     "string:",
	"complex64":  "number:float",
	"complex128": "number:double",
	"byte":       "string:byte",
	"rune":       "string:byte",
	// builtin golang objects
	"time.Time": "string:string",
}

func main() {
	flag.StringVar(&dir, "d", "./", "specific project dir")
	flag.Parse()
	err := ParseFromDir(dir)
	if err != nil {
		panic(err)
	}
	parseModel(pkgs)
	parseModel(rlpkgs)
	parseRouter()
	fd, err := os.Create(path.Join(dir, "swagger.json"))
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(swagger, "", "  ")
	fd.Write(b)
}

// ParseFromDir parse ast pkg from dir.
func ParseFromDir(dir string) (err error) {
	filepath.Walk(dir, func(fpath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fileInfo.IsDir() {
			return nil
		}
		err = parseFromDir(fpath)
		return err
	})
	return
}

func parseFromDir(dir string) (err error) {
	fset := token.NewFileSet()
	pkgFolder, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return
	}
	for k, p := range pkgFolder {
		pkgs[k] = p
	}
	return
}
func parseImport(dir string) (err error) {
	fset := token.NewFileSet()
	pkgFolder, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return
	}
	for k, p := range pkgFolder {
		rlpkgs[k] = p
	}
	return
}
func parseModel(pkgs map[string]*ast.Package) {
	for _, p := range pkgs {
		for _, f := range p.Files {
			for _, im := range f.Imports {
				if !isSystemPackage(im.Path.Value) {
					for _, gp := range _gopath {
						path := gp + "/src/" + strings.Trim(im.Path.Value, "\"")
						if isExist(path) {
							parseImport(path)
						}
					}
				}
			}
			scom := parseStructComment(f)
			for _, obj := range f.Scope.Objects {
				if obj.Kind == ast.Typ {
					objName := obj.Name
					schema := &Schema{
						Title: objName,
						Type:  "object",
					}
					ts, ok := obj.Decl.(*ast.TypeSpec)
					if !ok {
						fmt.Printf("obj type error %v ", obj.Kind)
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}
					properites := make(map[string]*Propertie)
					for _, fd := range st.Fields.List {
						if len(fd.Names) == 0 {
							continue
						}
						name, required, omit, desc := parseFieldTag(fd)
						if omit {
							continue
						}
						isSlice, realType, sType := typeAnalyser(fd)
						if (isSlice && isBasicType(realType)) || sType == "object" {
							if len(strings.Split(realType, " ")) > 1 {
								realType = strings.Replace(realType, " ", ".", -1)
								realType = strings.Replace(realType, "&", "", -1)
								realType = strings.Replace(realType, "{", "", -1)
								realType = strings.Replace(realType, "}", "", -1)
							}
						}
						mp := &Propertie{}
						if isSlice {
							mp.Type = "array"
							if isBasicType(strings.Replace(realType, "[]", "", -1)) {
								typeFormat := strings.Split(sType, ":")
								mp.Items = &Propertie{
									Type:   typeFormat[0],
									Format: typeFormat[1],
								}
							} else {
								ss := strings.Split(realType, ".")
								mp.RefImport = ss[len(ss)-1]
								mp.Type = "array"
								mp.Items = &Propertie{
									Ref:  "#/definitions/" + mp.RefImport,
									Type: sType,
								}
							}
						} else {
							if sType == "object" {
								ss := strings.Split(realType, ".")
								mp.RefImport = ss[len(ss)-1]
								mp.Type = sType
								mp.Ref = "#/definitions/" + mp.RefImport
							} else if isBasicType(realType) {
								typeFormat := strings.Split(sType, ":")
								mp.Type = typeFormat[0]
								mp.Format = typeFormat[1]
							} else if realType == "map" {
								typeFormat := strings.Split(sType, ":")
								mp.AdditionalProperties = &Propertie{
									Type:   typeFormat[0],
									Format: typeFormat[1],
								}
							}
						}
						if name == "" {
							name = fd.Names[0].Name
						}
						if required {
							schema.Required = append(schema.Required, name)
						}
						mp.Description = desc
						if scm, ok := scom[obj.Name]; ok {
							if cm, ok := scm.field[fd.Names[0].Name]; ok {
								mp.Description = cm + desc
							}
						}
						properites[name] = mp
					}
					if scm, ok := scom[obj.Name]; ok {
						schema.Description = scm.comment
					}
					schema.Properties = properites
					definitions[schema.Title] = schema
				}
			}
		}
	}
}
func parseFieldTag(field *ast.Field) (name string, required, omit bool, tagDes string) {
	if field.Tag == nil {
		return
	}
	tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
	param := tag.Get("form")
	if param != "" {
		params := strings.Split(param, ",")
		if len(params) > 0 {
			name = params[0]
		}
		if len(params) == 2 && params[1] == "split" {
			tagDes = "数组,按逗号分隔"
		}
	}
	if def := tag.Get("default"); def != "" {
		tagDes = fmt.Sprintf("%s 默认值 %s", tagDes, def)
	}
	validate := tag.Get("validate")
	if validate != "" {
		params := strings.Split(validate, ",")
		for _, param := range params {
			switch {
			case param == "required":
				required = true
			case strings.HasPrefix(param, "min"):
				tagDes = fmt.Sprintf("%s 最小值 %s", tagDes, strings.Split(param, "=")[1])
			case strings.HasPrefix(param, "max"):
				tagDes = fmt.Sprintf("%s 最大值 %s", tagDes, strings.Split(param, "=")[1])
			}
		}
	}
	// parse json response.
	json := tag.Get("json")
	if json != "" {
		jsons := strings.Split(json, ",")
		if len(jsons) > 0 {
			if jsons[0] == "-" {
				omit = true
				return
			}
		}
	}
	return
}

func parseRouter() {
	for _, p := range pkgs {
		if p.Name != "http" {
			continue
		}
		fmt.Printf("开始解析生成swagger文档\n")
		for _, f := range p.Files {
			for _, decl := range f.Decls {
				if fdecl, ok := decl.(*ast.FuncDecl); ok {
					if fdecl.Doc != nil {
						path, req, resp, item, err := parseFuncDoc(fdecl.Doc)
						if err != nil {
							fmt.Printf("解析失败 注解错误 %v\n", err)
							continue
						}
						if path != "" && err == nil {
							fmt.Printf("解析 %s 完成 请求参数为 %s 返回结构为 %s\n", path, req, resp)
							swagger.Paths[path] = item
						}
					}
				}
			}
		}
	}
}
func parseFuncDoc(f *ast.CommentGroup) (path, reqObj, respObj string, item *Item, err error) {
	item = new(Item)
	op := new(Operation)
	params := make([]*Parameter, 0)
	response := make(map[string]*Response)
	for _, d := range f.List {
		t := strings.TrimSpace(strings.TrimPrefix(d.Text, "//"))
		content := strings.Split(t, " ")
		switch content[0] {
		case "@params":
			if len(content) < 2 {
				err = fmt.Errorf("err params %s", content)
				return
			}
			reqObj = content[1]
			if model, ok := definitions[content[1]]; ok {
				for n, p := range model.Properties {
					param := &Parameter{
						In:          "query",
						Name:        n,
						Description: p.Description,
						Type:        p.Type,
						Format:      p.Format,
					}
					for _, p := range model.Required {
						if p == n {
							param.Required = true
						}
					}
					params = append(params, param)
				}
			} else {
				err = ErrParams
				return
			}
		case "@router":
			if len(content) != 3 {
				err = ErrParams
				return
			}
			switch content[1] {
			case "get":
				item.Get = op
			case "post":
				item.Post = op
			}
			path = content[2]
			op.OperationID = path
		case "@response":
			if len(content) < 2 {
				err = fmt.Errorf("err response %s", content)
				return
			}
			var (
				isarray bool
				ismap   bool
			)
			if strings.HasPrefix(content[1], "[]") {
				isarray = true
				respObj = content[1][2:]
			} else if strings.HasPrefix(content[1], "map[]") {
				ismap = true
				respObj = content[1][5:]
			} else {
				respObj = content[1]
			}
			defini, ok := definitions[respObj]
			if !ok {
				err = ErrParams
				return
			}
			var resp *Propertie
			if isarray {
				resp = &Propertie{
					Type: "array",
					Items: &Propertie{
						Type: "object",
						Ref:  "#/definitions/" + respObj,
					},
				}
			} else if ismap {
				resp = &Propertie{
					Type: "object",
					AdditionalProperties: &Propertie{
						Ref: "#/definitions/" + respObj,
					},
				}
			} else {
				resp = &Propertie{
					Type: "object",
					Ref:  "#/definitions/" + respObj,
				}
			}

			response["200"] = &Response{
				Schema: &Schema{
					Type: "object",
					Properties: map[string]*Propertie{
						"code": &Propertie{
							Type:        "integer",
							Description: "错误码描述",
						},
						"data": resp,
						"message": &Propertie{
							Type:        "string",
							Description: "错误码文本描述",
						},
						"ttl": &Propertie{
							Type:        "integer",
							Format:      "int64",
							Description: "客户端限速时间",
						},
					},
				},
				Description: "服务成功响应内容",
			}
			op.Responses = response
			for _, rl := range defini.Properties {
				if rl.RefImport != "" {
					swagger.Definitions[rl.RefImport] = definitions[rl.RefImport]
				}
			}
			swagger.Definitions[respObj] = defini
		case "@description":
			op.Description = content[1]
		}
	}
	op.Parameters = params
	return
}

type structComment struct {
	comment string
	field   map[string]string
}

func parseStructComment(f *ast.File) (scom map[string]structComment) {
	scom = make(map[string]structComment)
	for _, d := range f.Decls {
		switch specDecl := d.(type) {
		case *ast.GenDecl:
			if specDecl.Tok == token.TYPE {
				for _, s := range specDecl.Specs {
					switch tp := s.(*ast.TypeSpec).Type.(type) {
					case *ast.StructType:
						fcom := make(map[string]string)
						for _, fd := range tp.Fields.List {
							if len(fd.Names) == 0 {
								continue
							}
							if len(fd.Comment.Text()) > 0 {
								fcom[fd.Names[0].Name] = strings.TrimSuffix(fd.Comment.Text(), "\n")
							}
						}
						sspec := s.(*ast.TypeSpec)
						scom[sspec.Name.String()] = structComment{comment: strings.TrimSuffix(specDecl.Doc.Text(), "\n"), field: fcom}
					}
				}
			}
		}
	}
	return
}
func isBasicType(Type string) bool {
	if _, ok := basicTypes[Type]; ok {
		return true
	}
	return false
}

func typeAnalyser(f *ast.Field) (isSlice bool, realType, swaggerType string) {
	if arr, ok := f.Type.(*ast.ArrayType); ok {
		if isBasicType(fmt.Sprint(arr.Elt)) {
			return true, fmt.Sprintf("[]%v", arr.Elt), basicTypes[fmt.Sprint(arr.Elt)]
		}
		if mp, ok := arr.Elt.(*ast.MapType); ok {
			return false, fmt.Sprintf("map[%v][%v]", mp.Key, mp.Value), "object"
		}
		if star, ok := arr.Elt.(*ast.StarExpr); ok {
			return true, fmt.Sprint(star.X), "object"
		}
		basicType := fmt.Sprint(arr.Elt)
		if object, isStdLibObject := stdlibObject[basicType]; isStdLibObject {
			basicType = object

		}
		if k, ok := basicTypes[basicType]; ok {
			return true, basicType, k
		}
		return true, fmt.Sprint(arr.Elt), "object"
	}
	switch t := f.Type.(type) {
	case *ast.StarExpr:
		basicType := fmt.Sprint(t.X)
		if k, ok := basicTypes[basicType]; ok {
			return false, basicType, k
		}
		return false, basicType, "object"
	case *ast.MapType:
		val := fmt.Sprintf("%v", t.Value)
		if isBasicType(val) {
			return false, "map", basicTypes[val]
		}
		return false, val, "object"
	}
	basicType := fmt.Sprint(f.Type)
	if object, isStdLibObject := stdlibObject[basicType]; isStdLibObject {
		basicType = object
	}
	if k, ok := basicTypes[basicType]; ok {
		return false, basicType, k
	}
	return false, basicType, "object"
}

func isSystemPackage(pkgpath string) bool {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		goroot = runtime.GOROOT()
	}
	wg, _ := filepath.EvalSymlinks(filepath.Join(goroot, "src", "pkg", pkgpath))
	if isExist(wg) {
		return true
	}
	wg, _ = filepath.EvalSymlinks(filepath.Join(goroot, "src", pkgpath))
	return isExist(wg)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
