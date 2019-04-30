package main

import (
	"fmt"
	"github.com/bilibili/kratos/tool/bmproto/pkg/generator"
	"github.com/bilibili/kratos/tool/bmproto/pkg/naming"
	"github.com/bilibili/kratos/tool/bmproto/pkg/project"
	"github.com/bilibili/kratos/tool/bmproto/pkg/tag"
	"github.com/bilibili/kratos/tool/bmproto/pkg/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/bilibili/kratos/tool/bmproto/pkg/gen"
)

type tpl struct {
	generator.Base
	// key serviceName+methodName
	*project.ProjectInfo
}

func main() {
	g := TplGenerator()
	gen.Main(g)
}

// BmGenerator ...
func TplGenerator() *tpl {
	t := &tpl{}
	return t
}

func (t *tpl) Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	t.Setup(in)
	t.ProjectInfo, _ = project.NewProjInfo(t.GenFiles[0].GetName(),
		generator.GoModuleDirName, generator.GoModuleImportPath)
	resp := &plugin.CodeGeneratorResponse{}
	for _, f := range t.GenFiles {
		for _, svc := range f.Service {
			respFile := t.generateServiceImpl(f, svc)
			if respFile != nil {
				resp.File = append(resp.File, respFile)
			}
		}
	}
	return resp
}

// generateServiceImpl returns service implementation file service/{prefix}/service.go
// if file not exists
// else it returns nil
func (t *tpl) generateServiceImpl(file *descriptor.FileDescriptorProto, svc *descriptor.ServiceDescriptorProto) *plugin.CodeGeneratorResponse_File {
	resp := new(plugin.CodeGeneratorResponse_File)
	prefix := naming.GetVersionPrefix(t.GenPkgName)
	importPath, err := naming.GetGoImportPathForPb(file.GetName(),
		generator.GoModuleImportPath, generator.GoModuleDirName)
	// panic(fmt.Sprintf("%v %v %v %v %v", file.GetName(),
	// 	generator.GoModuleImportPath, generator.GoModuleDirName,importPath, err))
	if err != nil {
		importPath = "UNKNOWN IMPORT PATH, PLEASE CHANGE THIS YOURSELF, " + err.Error()
	}
	var alias = t.getPkgAlias()

	name := "service/" + prefix + "/" + utils.LcFirst(svc.GetName()) + ".go"
	if t.ProjectInfo.HasInternalPkg {
		name = "internal/" + name
	}
	name = t.ProjectInfo.PathRefToProj + "/" + name
	resp.Name = &name
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		// Insert methods if file already exists
		fset := token.NewFileSet()
		astTree, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			panic("parse file error: " + name + " err: " + err.Error())
		}
		v := visitor{funcMap: map[string]bool{}}
		ast.Walk(v, astTree)

		t.Output.Reset()
		buf, err := ioutil.ReadFile(name)
		if err != nil {
			panic("cannot read file:" + name)
		}
		t.P(string(buf))
		t.generateImplMethods(file, svc, v.funcMap)
		resp.Content = proto.String(t.FormattedOutput())
		t.Output.Reset()
		return resp
	}

	tplPkg := "service"
	if t.GenPkgName[:1] == "v" {
		tplPkg = t.GenPkgName
	}
	t.P(`package `, tplPkg)
	t.P()
	t.P(`import (`)
	t.P(`	`, alias, ` "`, importPath, `"`)
	t.P(`	"context"`)
	t.P(`)`)
	for pkg, importPath := range t.Deps {
		//t.RegisterPackageName()
		t.P(`import `, pkg, ` `, importPath)
	}
	svcStructName := naming.ServiceName(svc) + "Service"
	t.P(`// `, svcStructName, ` struct`)
	t.P(`type `, svcStructName, ` struct {`)
	t.P(`//	conf   *conf.Config`)
	t.P(`	// optionally add other properties here, such as dao`)
	t.P(`	// 值得注意的是，多个service公用一个Dao时，可以在外面New一个Dao对象传进来`)
	t.P(`	// 不必每个service都New一个Dao，以免造成资源浪费(例如mysql连接等)`)
	t.P(`	// dao *dao.Dao`)
	t.P(`}`)
	t.P()
	t.P(`//New`, svcStructName, ` init`)
	t.P(`func New`, svcStructName, `(`)
	t.P(`// c *conf.Config`)
	t.P(`) (s *`, svcStructName, `) {`)
	t.P(`	s = &`, svcStructName, `{`)
	t.P(`//		conf:   c,`)
	t.P(`	}`)
	t.P(`	return s`)
	t.P(`}`)

	comments, err := t.Reg.ServiceComments(file, svc)
	if err == nil {
		t.PrintComments(comments)
	}
	t.P()
	t.generateImplMethods(file, svc, map[string]bool{})
	resp.Content = proto.String(t.FormattedOutput())
	t.Output.Reset()
	return resp
}

func (t *tpl) generateImplMethods(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto,
	existMap map[string]bool) {
	var pkgName = t.getPkgAlias()
	svcName := naming.ServiceName(service) + "Service"

	for _, method := range service.Method {
		methName := naming.MethodName(method)
		if existMap[methName] {
			continue
		}
		comments, err := t.Reg.MethodComments(file, service, method)
		tags := tag.GetTagsInComment(comments.Leading)
		respDynamic := tag.GetTagValue("dynamic_resp", tags) == "true"
		genImp := func(dynamicRet bool) {
			t.P(`// `, methName, " implementation")
			if err == nil {
				t.PrintComments(comments)
			}
			outputType := t.GoTypeName(method.GetOutputType())
			inputType := t.GoTypeName(method.GetInputType())
			var body string
			var ownPkg = t.IsOwnPackage(method.GetOutputType())
			var reqOwnPkg = t.IsOwnPackage(method.GetInputType())
			var respType string
			var reqType string
			if ownPkg {
				respType = pkgName + "." + outputType
			} else {
				respType = outputType
			}
			if reqOwnPkg {
				reqType = pkgName + "." + inputType
			} else {
				reqType = inputType
			}
			if dynamicRet {
				body = fmt.Sprintf(`func (s *%s) %s(ctx context.Context, req *%s) (resp interface{}, err error) {`,
					svcName, methName, reqType)
			} else {
				body = fmt.Sprintf(`func (s *%s) %s(ctx context.Context, req *%s) (resp *%s, err error) {`,
					svcName, methName, reqType, respType)
			}

			t.P(body)
			t.P(fmt.Sprintf("resp = &%s{}", respType))
			t.P(`	return`)
			t.P(`}`)
			t.P()
		}
		genImp(respDynamic)
	}
}

type visitor struct {
	funcMap map[string]bool
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	switch d := n.(type) {
	case *ast.FuncDecl:
		v.funcMap[d.Name.Name] = true
	}
	return v
}

// pb包的别名
// 用户生成service实现模板时，对pb文件的引用
// 如果是v*的package 则为v*pb
// 其他为pb
func (t *tpl) getPkgAlias() string {
	if t.GenPkgName == "" {
		return "pb"
	}
	if t.GenPkgName[:1] == "v" {
		return t.GenPkgName + "pb"
	}
	return "pb"
}
