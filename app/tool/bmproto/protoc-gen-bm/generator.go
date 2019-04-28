// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen"
	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen/stringutils"
	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen/typemap"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
	"github.com/siddontang/go/ioutil2"
)

var legacyPathMapping = map[string]string{
	"live.webucenter":       "live.web-ucenter",
	"live.webroom":          "live.web-room",
	"live.appucenter":       "live.app-ucenter",
	"live.appblink":         "live.app-blink",
	"live.approom":          "live.app-room",
	"live.appinterface":     "live.app-interface",
	"live.liveadmin":        "live.live-admin",
	"live.resource":         "live.resource",
	"live.livedemo":         "live.live-demo",
	"live.lotteryinterface": "live.lottery-interface",
}

type bm struct {
	filesHandled int

	reg *typemap.Registry

	// Map to record whether we've built each package
	pkgs          map[string]string
	pkgNamesInUse map[string]bool

	importPrefix string            // String to prefix to imported package file names.
	importMap    map[string]string // Mapping from .proto file name to import path.
	tpl          bool              // only generate service template file, no docs, no .bm.go, default false

	// Package naming:
	genPkgName          string // Name of the package that we're generating
	fileToGoPackageName map[*descriptor.FileDescriptorProto]string

	// List of files that were inputs to the generator. We need to hold this in
	// the struct so we can write a header for the file that lists its inputs.
	genFiles []*descriptor.FileDescriptorProto

	// Output buffer that holds the bytes we want to write out for a single file.
	// Gets reset after working on a file.
	output *bytes.Buffer

	deps map[string]string
}

// if current dir is a go-common project
// or is the internal directory of a go-common project
// this present a project info
type projectInfo struct {
	absolutePath string
	// relative to go-common
	importPath string
	name       string
	department string
	// interface, service, admin ...
	typ            string
	hasInternalPkg bool
	// 从工作目录到project目录的相对路径 比如./ ../
	pathRefToProj string
}

// projectInfo for current directory
var projInfo *projectInfo

func bmGenerator() *bm {
	t := &bm{
		pkgs:                make(map[string]string),
		pkgNamesInUse:       make(map[string]bool),
		importMap:           make(map[string]string),
		fileToGoPackageName: make(map[*descriptor.FileDescriptorProto]string),
		output:              bytes.NewBuffer(nil),
	}

	return t
}

func (t *bm) Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	params, err := parseCommandLineParams(in.GetParameter())
	if err != nil {
		gen.Fail("could not parse parameters passed to --bm_out", err.Error())
	}
	t.importPrefix = params.importPrefix
	t.importMap = params.importMap
	t.tpl = params.tpl

	t.genFiles = gen.FilesToGenerate(in)

	// Collect information on types.
	t.reg = typemap.New(in.ProtoFile)

	t.registerPackageName("context")
	t.registerPackageName("ioutil")
	t.registerPackageName("proto")

	// Time to figure out package names of objects defined in protobuf. First,
	// we'll figure out the name for the package we're generating.
	genPkgName, err := deduceGenPkgName(t.genFiles)
	if err != nil {
		gen.Fail(err.Error())
	}
	t.genPkgName = genPkgName

	// Next, we need to pick names for all the files that are dependencies.
	if len(in.ProtoFile) > 0 {
		t.initProjInfo(in.ProtoFile[0])
	}
	for _, f := range in.ProtoFile {
		if fileDescSliceContains(t.genFiles, f) {
			// This is a file we are generating. It gets the shared package name.
			t.fileToGoPackageName[f] = t.genPkgName
		} else {
			// This is a dependency. Use its package name.
			name := f.GetPackage()
			if name == "" {
				name = stringutils.BaseName(f.GetName())
			}
			name = stringutils.CleanIdentifier(name)
			alias := t.registerPackageName(name)
			t.fileToGoPackageName[f] = alias
		}
	}

	// Showtime! Generate the response.
	resp := new(plugin.CodeGeneratorResponse)
	for _, f := range t.genFiles {
		respFile := t.generate(f)
		if respFile != nil {
			resp.File = append(resp.File, respFile)
		}
		for _, s := range f.Service {
			docResp := t.generateDoc(f, s)
			if docResp != nil {
				resp.File = append(resp.File, docResp)
			}
		}

		if t.tpl {
			if projInfo != nil {
				for _, s := range f.Service {
					serviceResp := t.generateServiceImpl(f, s)
					if serviceResp != nil {
						resp.File = append(resp.File, serviceResp)
					}
				}
			}
		}
	}
	return resp
}

// lookupProjPath get project path by proto absolute path
// assume that proto is in the project's model directory
func lookupProjPath(protoAbs string) (result string) {
	lastIndex := len(protoAbs)
	curPath := protoAbs

	for lastIndex > 0 {
		if ioutil2.FileExists(curPath+"/cmd") && ioutil2.FileExists(curPath+"/api") {
			result = curPath
			return
		}
		lastIndex = strings.LastIndex(curPath, string(os.PathSeparator))
		curPath = protoAbs[:lastIndex]
	}
	result = ""
	return
}

func (t *bm) initProjInfo(file *descriptor.FileDescriptorProto) {
	var err error
	projInfo = &projectInfo{}
	defer func() {
		if err != nil {
			projInfo = nil
		}
	}()
	wd, err := os.Getwd()
	if err != nil {
		panic("cannot get working directory")
	}
	protoAbs := wd + "/" + file.GetName()
	appIndex := strings.Index(wd, "go-common/app/")
	if appIndex == -1 {
		err = errors.New("not in go-common/app/")
		return
	}

	projPath := lookupProjPath(protoAbs)
	if projPath == "" {
		err = errors.New("not in project")
		return
	}
	if strings.Contains(wd, projPath) {
		rest := strings.Replace(wd, projPath, "", 1)
		projInfo.pathRefToProj = "./"
		if rest != "" {
			split := strings.Split(rest, "/")
			ref := ""
			for i := 0; i < len(split)-1; i++ {
				ref = ref + "../"
			}
			projInfo.pathRefToProj = ref
		}
	}
	projInfo.absolutePath = projPath
	if ioutil2.FileExists(projPath + "/internal") {
		projInfo.hasInternalPkg = true
	}

	relativePath := projInfo.absolutePath[appIndex+len("go-common/app/"):]
	projInfo.importPath = "go-common/app/" + relativePath
	split := strings.Split(relativePath, "/")
	projInfo.typ = split[0]
	projInfo.department = split[1]
	projInfo.name = split[2]
}

// find tag between backtick, start & end is the position of backtick
func getLineTag(line string) (tag reflect.StructTag, start int, end int) {
	start = strings.Index(line, "`")
	end = strings.LastIndex(line, "`")
	if end <= start {
		return
	}
	tag = reflect.StructTag(line[start+1 : end])
	return
}

func getCommentWithoutTag(comment string) []string {
	var lines []string
	if comment == "" {
		return lines
	}
	split := strings.Split(strings.TrimRight(comment, "\n\r"), "\n")
	for _, line := range split {
		tag, _, _ := getLineTag(line)
		if tag == "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func getTagsInComment(comment string) []reflect.StructTag {
	split := strings.Split(comment, "\n")
	var tagsInComment []reflect.StructTag
	for _, line := range split {
		tag, _, _ := getLineTag(line)
		if tag != "" {
			tagsInComment = append(tagsInComment, tag)
		}
	}
	return tagsInComment
}

func getTagValue(key string, tags []reflect.StructTag) string {
	for _, t := range tags {
		val := t.Get(key)
		if val != "" {
			return val
		}
	}
	return ""
}

// Is this field repeated?
func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

func (t *bm) isMap(field *descriptor.FieldDescriptorProto) bool {
	if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return false
	}
	md := t.reg.MessageDefinition(field.GetTypeName())
	if md == nil || !md.Descriptor.GetOptions().GetMapEntry() {
		return false
	}
	return true
}

func (t *bm) generateToc(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	for _, method := range service.Method {
		comment, _ := t.reg.MethodComments(file, service, method)
		tags := getTagsInComment(comment.Leading)
		_, path, newPath := t.getHttpInfo(file, service, method, tags)

		cleanComments := getCommentWithoutTag(comment.Leading)
		var title string
		if len(cleanComments) > 0 {
			title = cleanComments[0]
		}

		// 如果有老的路径，只显示老的路径文档
		if path != "" {
			anchor := strings.Replace(path, "/", "", -1)
			t.P(fmt.Sprintf("- [%s](#%s) %s", path, anchor, title))
		} else {
			anchor := strings.Replace(newPath, "/", "", -1)
			t.P(fmt.Sprintf("- [%s](#%s) %s", newPath, anchor, title))
		}
	}
}

func (t *bm) generateDoc(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) *plugin.CodeGeneratorResponse_File {
	resp := new(plugin.CodeGeneratorResponse_File)
	var name = goFileName(file, "."+lcFirst(service.GetName())+".md")
	resp.Name = &name
	t.P("<!-- package=" + file.GetPackage() + " -->")
	t.generateToc(file, service)
	for _, method := range service.Method {
		comment, err := t.reg.MethodComments(file, service, method)
		tags := getTagsInComment(comment.Leading)
		cleanComments := getCommentWithoutTag(comment.Leading)
		midwaresStr := getTagValue("midware", tags)
		needAuth := false
		if midwaresStr != "" {
			split := strings.Split(midwaresStr, ",")
			for _, m := range split {
				if m == "auth" {
					needAuth = true
					break
				}
			}
		}
		t.P()
		httpMethod, legacyPath, path := t.getHttpInfo(file, service, method, tags)
		if legacyPath != "" {
			path = legacyPath
		}
		t.P("## " + path)

		if err == nil {
			if len(cleanComments) == 0 {
				t.P(`### 无标题`)
			} else {
				t.P(`###`, strings.Join(cleanComments, "\n"))
			}
		}
		t.P()

		if needAuth {
			t.P(`> `, "需要登录")
			t.P()
		}

		t.P("#### 方法：" + httpMethod)
		t.P()

		t.genRequestParam(file, service, method)
		t.P("#### 响应")
		t.P()

		t.P("```javascript")
		t.P(`{`)
		t.P(`    "code": 0,`)
		t.P(`    "message": "ok",`)
		t.P(t.getExampleJson(file, service, method))
		t.P(`}`)
		t.P("```")
		t.P()
	}
	resp.Content = proto.String(t.output.String())
	t.output.Reset()
	return resp
}

func (t *bm) genRequestParam(
	file *descriptor.FileDescriptorProto,
	svc *descriptor.ServiceDescriptorProto,
	method *descriptor.MethodDescriptorProto) {
	md := t.reg.MessageDefinition(method.GetInputType())
	t.P(`#### 请求参数`)
	t.P()

	var outputs []string
	for i, f := range md.Descriptor.Field {
		if f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			// 如果有message 只能以json的方式显示参数了
			var buf = &[]string{}
			t.exampleJsonForMsg(md, file, buf, "", 0, "")
			j := strings.Join(*buf, "\n")
			t.P("```javascript")
			t.P(j)
			t.P("```")
			t.P()
			return
		}
		if i == 0 {
			outputs = append(outputs, `|参数名|必选|类型|描述|`)
			outputs = append(outputs, `|:---|:---|:---|:---|`)
		}
		fComment, _ := t.reg.FieldComments(file, md, f)
		var tags []reflect.StructTag
		{
			//get required info from gogoproto.moretags
			moretags := getMoreTags(f)
			if moretags != nil {
				tags = []reflect.StructTag{reflect.StructTag(*moretags)}
			}
		}
		if len(tags) == 0 {
			tags = getTagsInComment(fComment.Leading)
		}
		validateTag := getTagValue("validate", tags)
		var validateRules []string
		if validateTag != "" {
			validateRules = strings.Split(validateTag, ",")
		}
		required := false
		for _, rule := range validateRules {
			if rule == "required" {
				required = true
			}
		}
		requiredDesc := "是"
		if !required {
			requiredDesc = "否"
		}
		_, typeName := t.mockValueForField(f, tags)
		split := strings.Split(fComment.Leading, "\n")
		desc := ""
		for _, line := range split {
			if line != "" {
				tag, _, _ := getLineTag(line)
				if tag == "" {
					desc += line
				}
			}
		}
		outputs = append(outputs, fmt.Sprintf(`|%s|%s|%s|%s|`, getJsonTag(f), requiredDesc, typeName, desc))
	}
	for _, s := range outputs {
		t.P(s)
	}
	t.P()
}

func (t *bm) getExampleJson(file *descriptor.FileDescriptorProto,
	svc *descriptor.ServiceDescriptorProto,
	method *descriptor.MethodDescriptorProto) string {
	md := t.reg.MessageDefinition(method.GetOutputType())
	var buf = &[]string{}
	t.exampleJsonForMsg(md, file, buf, "data", 4, "")
	return strings.Join(*buf, "\n")
}

func makeIndentStr(i int) string {
	return strings.Repeat(" ", i)
}

func (t *bm) mockValueForField(field *descriptor.FieldDescriptorProto,
	tags []reflect.StructTag) (mockVal string, typeName string) {
	tagMock := getTagValue("mock", tags)
	mockVal = "\"unknown\""
	typeName = "unknown"
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if tagMock == "true" || tagMock == "false" {
			mockVal = tagMock
		} else {
			mockVal = "true"
		}
		typeName = "bool"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		mockVal = "0.1"
		if tagMock != "" {
			if _, err := strconv.ParseFloat(tagMock, 64); err == nil {
				mockVal = tagMock
			}
		}
		typeName = "float"
	case
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		mockVal = "0"
		if tagMock != "" {
			if _, err := strconv.Atoi(tagMock); err == nil {
				mockVal = tagMock
			}
		}
		typeName = "integer"
	case
		descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		mockVal = `""`
		if tagMock != "" {
			mockVal = strconv.Quote(tagMock)
		}
		typeName = "string"
	}
	if isRepeated(field) {
		typeName = "多个" + typeName
	}
	return
}

func (t *bm) exampleJsonForMsg(
	msg *typemap.MessageDefinition,
	file *descriptor.FileDescriptorProto,
	buf *[]string, fieldName string, indent int, outEndComma string) {
	if fieldName == "" {
		*buf = append(*buf, makeIndentStr(indent)+"{")
	} else {
		*buf = append(*buf, makeIndentStr(indent)+fmt.Sprintf(`"%s": {`, fieldName))
	}
	num := len(msg.Descriptor.Field)
	for i, f := range msg.Descriptor.Field {
		isScalar := isScalar(f)
		fComment, _ := t.reg.FieldComments(file, msg, f)
		cleanComment := getCommentWithoutTag(fComment.Leading)
		for _, line := range cleanComment {
			if strings.Trim(line, " \t\n\r") != "" {
				*buf = append(*buf, makeIndentStr(indent+4)+"// "+line)
			}
		}

		endComma := ""
		if i < (num - 1) {
			endComma = ","
		}
		repeated := isRepeated(f)
		tags := getTagsInComment(fComment.Leading)
		if isScalar {
			mockVal, _ := t.mockValueForField(f, tags)

			if repeated {
				// "key" : [
				// 	value
				// ]
				*buf = append(*buf, makeIndentStr(indent+4)+`"`+getJsonTag(f)+`": [`)
				*buf = append(*buf, makeIndentStr(indent+8)+mockVal)
				*buf = append(*buf, makeIndentStr(indent+4)+`]`+endComma)
			} else {
				// "key" : value
				*buf = append(*buf, makeIndentStr(indent+4)+`"`+getJsonTag(f)+`": `+mockVal+endComma)
			}
		} else {
			isMap := t.isMap(f)
			if repeated {
				if isMap {
					*buf = append(*buf, makeIndentStr(indent+4)+`"`+getJsonTag(f)+`": {`)
				} else {
					*buf = append(*buf, makeIndentStr(indent+4)+`"`+getJsonTag(f)+`": [`)
				}
			}
			subMsg := t.reg.MessageDefinition(f.GetTypeName())
			if subMsg == nil {
				panic(fmt.Sprintf("%v%v", f.TypeName, f.Type))
			}
			nextIndent := indent + 4
			nextFname := getJsonTag(f)
			if repeated {
				nextIndent = indent + 8
				nextFname = ""
			}
			if isMap {
				mapKeyField := subMsg.Descriptor.Field[0]
				mapValueField := subMsg.Descriptor.Field[1]
				keyDesc := "mapKey"
				if mapKeyField.GetType() != descriptor.FieldDescriptorProto_TYPE_STRING {
					keyDesc = "1"
				}

				if mapValueField.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					// "mapKey" : {
					// ...
					// }
					mapValueMsg := t.reg.MessageDefinition(mapValueField.GetTypeName())
					t.exampleJsonForMsg(mapValueMsg, file, buf, keyDesc, nextIndent, "")
				} else {
					// "mapKey" : "map value"
					val, _ := t.mockValueForField(mapValueField, tags)

					*buf = append(*buf, makeIndentStr(indent+8)+`"`+keyDesc+`": `+val)
				}
				*buf = append(*buf, makeIndentStr(indent+4)+`}`+endComma)
			} else {
				if repeated {
					t.exampleJsonForMsg(subMsg, file, buf, nextFname, nextIndent, "")
					*buf = append(*buf, makeIndentStr(indent+4)+`]`+endComma)
				} else {
					t.exampleJsonForMsg(subMsg, file, buf, nextFname, nextIndent, endComma)
				}
			}
		}
	}
	*buf = append(*buf, makeIndentStr(indent)+"}"+outEndComma)
}

// Is this field a scalar numeric type?
func isScalar(field *descriptor.FieldDescriptorProto) bool {
	if field.Type == nil {
		return false
	}
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_BYTES,
		descriptor.FieldDescriptorProto_TYPE_STRING:
		return true
	default:
		return false
	}
}

func (t *bm) registerPackageName(name string) (alias string) {
	alias = name
	i := 1
	for t.pkgNamesInUse[alias] {
		alias = name + strconv.Itoa(i)
		i++
	}
	t.pkgNamesInUse[alias] = true
	t.pkgs[name] = alias
	return alias
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

// generateServiceImpl returns service implementation file service/{prefix}/service.go
// if file not exists
// else it returns nil
func (t *bm) generateServiceImpl(file *descriptor.FileDescriptorProto, svc *descriptor.ServiceDescriptorProto) *plugin.CodeGeneratorResponse_File {
	resp := new(plugin.CodeGeneratorResponse_File)
	prefix := t.getVersionPrefix()
	importPath := t.getPbImportPath(file.GetName())
	var alias = t.getPkgAlias()
	confPath := projInfo.importPath + "/conf"
	if projInfo.hasInternalPkg {
		confPath = projInfo.importPath + "/internal/conf"
	}
	name := "service/" + prefix + "/" + lcFirst(svc.GetName()) + ".go"
	if projInfo.hasInternalPkg {
		name = "internal/" + name
	}
	name = projInfo.pathRefToProj + name
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

		t.output.Reset()
		buf, err := ioutil.ReadFile(name)
		if err != nil {
			panic("cannot read file:" + name)
		}
		t.P(string(buf))
		t.generateBmImpl(file, svc, v.funcMap)
		resp.Content = proto.String(t.formattedOutput())
		t.output.Reset()
		return resp
	}

	tplPkg := "service"
	if t.genPkgName[:1] == "v" {
		tplPkg = t.genPkgName
	}
	t.P(`package `, tplPkg)
	t.P()
	t.P(`import (`)
	t.P(`	`, alias, ` "`, importPath, `"`)
	t.P(`	"`, confPath, `"`)
	t.P(`	"context"`)
	t.P(`)`)
	for pkg, importPath := range t.deps {
		t.P(`import `, pkg, ` `, importPath)
	}
	svcStructName := serviceName(svc) + "Service"
	t.P(`// `, svcStructName, ` struct`)
	t.P(`type `, svcStructName, ` struct {`)
	t.P(`	conf   *conf.Config`)
	t.P(`	// optionally add other properties here, such as dao`)
	t.P(`	// dao *dao.Dao`)
	t.P(`}`)
	t.P()
	t.P(`//New`, svcStructName, ` init`)
	t.P(`func New`, svcStructName, `(c *conf.Config) (s *`, svcStructName, `) {`)
	t.P(`	s = &`, svcStructName, `{`)
	t.P(`		conf:   c,`)
	t.P(`	}`)
	t.P(`	return s`)
	t.P(`}`)

	comments, err := t.reg.ServiceComments(file, svc)
	if err == nil {
		t.printComments(comments)
	}
	t.P()
	t.generateBmImpl(file, svc, map[string]bool{})
	resp.Content = proto.String(t.formattedOutput())
	t.output.Reset()
	return resp
}

func (t *bm) generate(file *descriptor.FileDescriptorProto) *plugin.CodeGeneratorResponse_File {
	resp := new(plugin.CodeGeneratorResponse_File)
	if len(file.Service) == 0 {
		return nil
	}

	t.generateFileHeader(file, t.genPkgName)

	t.generateImports(file)

	t.generateMiddlewareInfo(file)
	for i, service := range file.Service {
		t.generateService(file, service, i)
		t.generateSingleRoute(file, service, i)
	}

	t.generateFileDescriptor(file)

	resp.Name = proto.String(goFileName(file, ".bm.go"))
	resp.Content = proto.String(t.formattedOutput())
	t.output.Reset()

	t.filesHandled++
	return resp
}

func (t *bm) generateMiddlewareInfo(file *descriptor.FileDescriptorProto) {
	t.P()
	for _, service := range file.Service {
		name := serviceName(service)
		for _, method := range service.Method {
			_, _, path := t.getHttpInfo(file, service, method, nil)
			t.P(`var Path`, name, methodName(method), ` = "`, path, `"`)
		}
		t.P()
	}
}
func (t *bm) generateFileHeader(file *descriptor.FileDescriptorProto, pkgName string) {
	t.P("// Code generated by protoc-gen-bm ", gen.Version, ", DO NOT EDIT.")
	t.P("// source: ", file.GetName())
	t.P()
	if t.filesHandled == 0 {
		t.P("/*")
		t.P("Package ", t.genPkgName, " is a generated blademaster stub package.")
		t.P("This code was generated with go-common/app/tool/bmgen/protoc-gen-bm ", gen.Version, ".")
		t.P()
		comment, err := t.reg.FileComments(file)
		if err == nil && comment.Leading != "" {
			for _, line := range strings.Split(comment.Leading, "\n") {
				line = strings.TrimPrefix(line, " ")
				// ensure we don't escape from the block comment
				line = strings.Replace(line, "*/", "* /", -1)
				t.P(line)
			}
			t.P()
		}
		t.P("It is generated from these files:")
		for _, f := range t.genFiles {
			t.P("\t", f.GetName())
		}
		t.P("*/")
	}
	t.P(`package `, pkgName)
	t.P()
}

func (t *bm) generateImports(file *descriptor.FileDescriptorProto) {
	if len(file.Service) == 0 {
		return
	}
	t.P(`import (`)
	//t.P(`	`,t.pkgs["context"], ` "context"`)
	t.P(`	"context"`)
	t.P()
	t.P(`	bm "go-common/library/net/http/blademaster"`)
	t.P(`	"go-common/library/net/http/blademaster/binding"`)

	t.P(`)`)
	// It's legal to import a message and use it as an input or output for a
	// method. Make sure to import the package of any such message. First, dedupe
	// them.
	deps := make(map[string]string) // Map of package name to quoted import path.
	ourImportPath := path.Dir(goFileName(file, ""))
	for _, s := range file.Service {
		for _, m := range s.Method {
			defs := []*typemap.MessageDefinition{
				t.reg.MethodInputDefinition(m),
				t.reg.MethodOutputDefinition(m),
			}
			for _, def := range defs {
				// By default, import path is the dirname of the Go filename.
				importPath := path.Dir(goFileName(def.File, ""))
				if importPath == ourImportPath {
					continue
				}
				if substitution, ok := t.importMap[def.File.GetName()]; ok {
					importPath = substitution
				}
				importPath = t.importPrefix + importPath
				pkg := t.goPackageName(def.File)
				deps[pkg] = strconv.Quote(importPath)
			}
		}
	}
	t.deps = deps
	for pkg, importPath := range deps {
		t.P(`import `, pkg, ` `, importPath)
	}
	if len(deps) > 0 {
		t.P()
	}
	t.P()
	t.P(`// to suppressed 'imported but not used warning'`)
	t.P(`var _ *bm.Context`)
	t.P(`var _ context.Context`)
	t.P(`var _ binding.StructValidator`)

}

// P forwards to g.gen.P, which prints output.
func (t *bm) P(args ...string) {
	for _, v := range args {
		t.output.WriteString(v)
	}
	t.output.WriteByte('\n')
}

// Big header comments to makes it easier to visually parse a generated file.
func (t *bm) sectionComment(sectionTitle string) {
	t.P()
	t.P(`// `, strings.Repeat("=", len(sectionTitle)))
	t.P(`// `, sectionTitle)
	t.P(`// `, strings.Repeat("=", len(sectionTitle)))
	t.P()
}

func (t *bm) generateService(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto, index int) {
	servName := serviceName(service)

	t.sectionComment(servName + ` Interface`)
	t.generateBMInterface(file, service)

}

// import project/api的路径
func (t *bm) getPbImportPath(filename string) (importPath string) {
	wd, err := os.Getwd()
	if err != nil {
		panic("cannot get working directory")
	}
	index := strings.Index(wd, "go-common")
	if index == -1 {
		gen.Fail("must use inside go-common")
	}
	dir := filepath.Dir(filename)
	if dir != "." {
		importPath = wd + "/" + dir
	} else {
		importPath = wd
	}
	importPath = importPath[index:]
	return
}

// getProjPath return project path relative to GOPATH
func (t *bm) getProjPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic("cannot get working directory")
	}
	index := strings.Index(wd, "go-common")
	if index == -1 {
		gen.Fail("must use inside go-common")
	}
	projPkgPath := wd[index:]
	return projPkgPath
}

func lcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// TODO rename
func (t *bm) getLegacyPathPrefix(
	svc *descriptor.ServiceDescriptorProto, pathParts []string, isInternal bool) (uriPrefix string) {
	var parts []string
	parts = append(parts, pathParts[0])
	if isInternal {
		parts = append(parts, "internal")
	}
	parts = append(parts, pathParts[1:]...)
	uriPrefix = fmt.Sprintf("/x%s/%s", strings.Join(parts, "/"), lcFirst(svc.GetName()))
	return
}

func (t *bm) getHttpInfo(
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
	method *descriptor.MethodDescriptorProto,
	tags []reflect.StructTag,
) (httpMethod string, oldPath string, newPath string) {
	googleOptionInfo, err := ParseBMMethod(method)
	if err == nil {
		httpMethod = strings.ToUpper(googleOptionInfo.Method)
		p := googleOptionInfo.PathPattern
		if p != "" {
			oldPath = p
			newPath = p
			return
		}
	}

	if httpMethod == "" {
		// resolve http method
		httpMethod = getTagValue("method", tags)
		if httpMethod == "" {
			httpMethod = "GET"
		} else {
			httpMethod = strings.ToUpper(httpMethod)
		}
	}

	isLegacy, parts := t.convertLegacyPackage(file.GetPackage())
	if isLegacy {
		apiInternal := getTagValue("internal", tags) == "true"
		pathPrefix := t.getLegacyPathPrefix(service, parts, apiInternal)
		oldPath = pathPrefix + `/` + method.GetName()
	}
	newPath = "/" + file.GetPackage() + "." + service.GetName() + "/" + method.GetName()
	return
}

// 返回空，则不用考虑历史package
// 如果非空，则表示按照返回的pathParts做url规则
func (t *bm) convertLegacyPackage(pkgName string) (isLegacy bool, pathParts []string) {
	var splits = strings.Split(pkgName, ".")
	var remain []string
	if len(splits) >= 2 {
		splits = splits[0:2]
		remain = splits[2:]
	}
	var pkgPrefix = strings.Join(splits, ".")
	legacyPkg, isLegacy := legacyPathMapping[pkgPrefix]
	if isLegacy {
		legacyPkg = strings.Replace(pkgName, pkgPrefix, legacyPkg, 1)
		pathParts = append(pathParts, strings.Split(legacyPkg, ".")...)
		pathParts = append(pathParts, remain...)
	}
	return
}

func (t *bm) generateSingleRoute(
	file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
	index int) {
	// old mode is generate xx.route.go in the http pkg
	// new mode is generate route code in the same .bm.go
	// route rule /x{department}/{project-name}/{path_prefix}/method_name
	// generate each route method
	servName := serviceName(service)
	versionPrefix := t.getVersionPrefix()
	svcName := lcFirst(stringutils.CamelCase(versionPrefix)) + servName + "Svc"
	t.P(`var `, svcName, ` `, servName, `BMServer`)

	type methodInfo struct {
		httpMethod    string
		midwares      []string
		routeFuncName string
		path          string
		legacyPath    string
		methodName    string
	}
	var methList []methodInfo
	var allMidwareMap = make(map[string]bool)
	var isLegacyPkg = false
	for _, method := range service.Method {
		var httpMethod string
		var midwares []string
		comments, _ := t.reg.MethodComments(file, service, method)
		tags := getTagsInComment(comments.Leading)
		if getTagValue("dynamic", tags) == "true" {
			continue
		}
		httpMethod, legacyPath, path := t.getHttpInfo(file, service, method, tags)
		if legacyPath != "" {
			isLegacyPkg = true
		}

		midStr := getTagValue("midware", tags)
		if midStr != "" {
			midwares = strings.Split(midStr, ",")
			for _, m := range midwares {
				allMidwareMap[m] = true
			}
		}

		methName := methodName(method)
		inputType := t.goTypeName(method.GetInputType())

		routeName := lcFirst(stringutils.CamelCase(servName) +
			stringutils.CamelCase(methName))

		methList = append(methList, methodInfo{
			httpMethod:    httpMethod,
			midwares:      midwares,
			routeFuncName: routeName,
			path:          path,
			legacyPath:    legacyPath,
			methodName:    method.GetName(),
		})

		t.P(fmt.Sprintf("func %s (c *bm.Context) {", routeName))
		t.P(`	p := new(`, inputType, `)`)
		t.P(`	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {`)
		t.P(`		return`)
		t.P(`	}`)
		t.P(`	resp, err := `, svcName, `.`, methName, `(c, p)`)
		t.P(`	c.JSON(resp, err)`)
		t.P(`}`)
		t.P(``)
	}

	// generate route group
	var midList []string
	for m := range allMidwareMap {
		midList = append(midList, m+" bm.HandlerFunc")
	}

	sort.Strings(midList)

	// 注册老的路由的方法
	if isLegacyPkg {
		funcName := `Register` + stringutils.CamelCase(versionPrefix) + servName + `Service`
		t.P(`// `, funcName, ` Register the blademaster route with middleware map`)
		t.P(`// midMap is the middleware map, the key is defined in proto`)
		t.P(`func `, funcName, `(e *bm.Engine, svc `, servName, "BMServer, midMap map[string]bm.HandlerFunc)", ` {`)
		var keys []string
		for m := range allMidwareMap {
			keys = append(keys, m)
		}
		// to keep generated code consistent
		sort.Strings(keys)
		for _, m := range keys {
			t.P(m, ` := midMap["`, m, `"]`)
		}

		t.P(svcName, ` = svc`)
		for _, methInfo := range methList {
			var midArgStr string
			if len(methInfo.midwares) == 0 {
				midArgStr = ""
			} else {
				midArgStr = strings.Join(methInfo.midwares, ", ") + ", "
			}
			t.P(`e.`, methInfo.httpMethod, `("`, methInfo.legacyPath, `", `, midArgStr, methInfo.routeFuncName, `)`)
		}
		t.P(`	}`)
	}
	// 新的注册路由的方法
	var bmFuncName = fmt.Sprintf("Register%sBMServer", servName)
	t.P(`// `, bmFuncName, ` Register the blademaster route`)
	t.P(`func `, bmFuncName, `(e *bm.Engine, server `, servName, `BMServer) {`)
	t.P(svcName, ` = server`)
	for _, methInfo := range methList {
		t.P(`e.`, methInfo.httpMethod, `("`, methInfo.path, `",`, methInfo.routeFuncName, ` )`)
	}
	t.P(`	}`)
}

func (t *bm) generateBMInterface(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	servName := serviceName(service)

	comments, err := t.reg.ServiceComments(file, service)
	if err == nil {
		t.printComments(comments)
	}
	t.P(`type `, servName, `BMServer interface {`)
	for _, method := range service.Method {
		t.generateSignature(file, service, method, comments)
		t.P()
	}
	t.P(`}`)
}

// pb包的别名
// 用户生成service实现模板时，对pb文件的引用
// 如果是v*的package 则为v*pb
// 其他为pb
func (t *bm) getPkgAlias() string {
	if t.genPkgName == "" {
		return "pb"
	}
	if t.genPkgName[:1] == "v" {
		return t.genPkgName + "pb"
	}
	return "pb"
}

// 如果是v*开始的 返回v*
// 否则返回空
func (t *bm) getVersionPrefix() string {
	if t.genPkgName == "" {
		return ""
	}
	if t.genPkgName[:1] == "v" {
		return t.genPkgName
	}
	return ""
}

func (t *bm) generateBmImpl(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto,
	existMap map[string]bool) {
	var pkgName = t.getPkgAlias()
	svcName := serviceName(service) + "Service"

	for _, method := range service.Method {
		methName := methodName(method)
		if existMap[methName] {
			continue
		}
		comments, err := t.reg.MethodComments(file, service, method)
		tags := getTagsInComment(comments.Leading)
		respDynamic := getTagValue("dynamic_resp", tags) == "true"
		genImp := func(dynamicRet bool) {
			t.P(`// `, methName, " implementation")
			if err == nil {
				t.printComments(comments)
			}
			outputType := t.goTypeName(method.GetOutputType())
			inputType := t.goTypeName(method.GetInputType())
			var body string
			var ownPkg = t.isOwnPackage(method.GetOutputType())
			var respType string
			if ownPkg {
				respType = pkgName + "." + outputType
			} else {
				respType = outputType
			}
			if dynamicRet {
				body = fmt.Sprintf(`func (s *%s) %s(ctx context.Context, req *%s.%s) (resp interface{}, err error) {`,
					svcName, methName, pkgName, inputType)
			} else {
				body = fmt.Sprintf(`func (s *%s) %s(ctx context.Context, req *%s.%s) (resp *%s, err error) {`,
					svcName, methName, pkgName, inputType, respType)
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

func (t *bm) generateSignature(file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
	method *descriptor.MethodDescriptorProto,
	comments typemap.DefinitionComments) {
	comments, err := t.reg.MethodComments(file, service, method)

	methName := methodName(method)
	outputType := t.goTypeName(method.GetOutputType())
	inputType := t.goTypeName(method.GetInputType())
	tags := getTagsInComment(comments.Leading)
	if getTagValue("dynamic", tags) == "true" {
		return
	}

	if err == nil {
		t.printComments(comments)
	}

	respDynamic := getTagValue("dynamic_resp", tags) == "true"
	if respDynamic {
		t.P(fmt.Sprintf(`	%s(ctx context.Context, req *%s) (resp interface{}, err error)`,
			methName, inputType))
	} else {
		t.P(fmt.Sprintf(`	%s(ctx context.Context, req *%s) (resp *%s, err error)`,
			methName, inputType, outputType))
	}
}

func (t *bm) generateFileDescriptor(file *descriptor.FileDescriptorProto) {
	// Copied straight of of protoc-gen-go, which trims out comments.
	pb := proto.Clone(file).(*descriptor.FileDescriptorProto)
	pb.SourceCodeInfo = nil

	b, err := proto.Marshal(pb)
	if err != nil {
		gen.Fail(err.Error())
	}

	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	w.Write(b)
	w.Close()
	buf.Bytes()
}

func (t *bm) printComments(comments typemap.DefinitionComments) bool {
	text := strings.TrimSuffix(comments.Leading, "\n")
	if len(strings.TrimSpace(text)) == 0 {
		return false
	}
	split := strings.Split(text, "\n")
	for _, line := range split {
		t.P("// ", strings.TrimPrefix(line, " "))
	}
	return len(split) > 0
}

// Given a protobuf name for a Message, return the Go name we will use for that
// type, including its package prefix.
func (t *bm) goTypeName(protoName string) string {
	def := t.reg.MessageDefinition(protoName)
	if def == nil {
		gen.Fail("could not find message for", protoName)
	}

	var prefix string
	if pkg := t.goPackageName(def.File); pkg != t.genPkgName {
		prefix = pkg + "."
	}

	var name string
	for _, parent := range def.Lineage() {
		name += parent.Descriptor.GetName() + "_"
	}
	name += def.Descriptor.GetName()
	return prefix + name
}

func (t *bm) isOwnPackage(protoName string) bool {
	def := t.reg.MessageDefinition(protoName)
	if def == nil {
		gen.Fail("could not find message for", protoName)
	}
	pkg := t.goPackageName(def.File)
	return pkg == t.genPkgName
}

func (t *bm) goPackageName(file *descriptor.FileDescriptorProto) string {
	return t.fileToGoPackageName[file]
}

func (t *bm) formattedOutput() string {
	// Reformat generated code.
	fset := token.NewFileSet()
	raw := t.output.Bytes()
	ast, err := parser.ParseFile(fset, "", raw, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code,
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(raw))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		gen.Fail("bad Go source code was generated:", err.Error(), "\n"+src.String())
	}

	out := bytes.NewBuffer(nil)
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(out, fset, ast)
	if err != nil {
		gen.Fail("generated Go source code could not be reformatted:", err.Error())
	}

	return out.String()
}

func serviceName(service *descriptor.ServiceDescriptorProto) string {
	return stringutils.CamelCase(service.GetName())
}

func methodName(method *descriptor.MethodDescriptorProto) string {
	return stringutils.CamelCase(method.GetName())
}

func fileDescSliceContains(slice []*descriptor.FileDescriptorProto, f *descriptor.FileDescriptorProto) bool {
	for _, sf := range slice {
		if f == sf {
			return true
		}
	}
	return false
}

// deduceGenPkgName figures out the go package name to use for generated code.
// Will try to use the explicit go_package setting in a file (if set, must be
// consistent in all files). If no files have go_package set, then use the
// protobuf package name (must be consistent in all files)
func deduceGenPkgName(genFiles []*descriptor.FileDescriptorProto) (string, error) {
	var genPkgName string
	for _, f := range genFiles {
		name, explicit := goPackageName(f)
		if explicit {
			name = stringutils.CleanIdentifier(name)
			if genPkgName != "" && genPkgName != name {
				// Make sure they're all set consistently.
				return "", errors.Errorf("files have conflicting go_package settings, must be the same: %q and %q", genPkgName, name)
			}
			genPkgName = name
		}
	}
	if genPkgName != "" {
		return genPkgName, nil
	}

	// If there is no explicit setting, then check the implicit package name
	// (derived from the protobuf package name) of the files and make sure it's
	// consistent.
	for _, f := range genFiles {
		name, _ := goPackageName(f)
		name = stringutils.CleanIdentifier(name)
		if genPkgName != "" && genPkgName != name {
			return "", errors.Errorf("files have conflicting package names, must be the same or overridden with go_package: %q and %q", genPkgName, name)
		}
		genPkgName = name
	}

	// All the files have the same name, so we're good.
	return genPkgName, nil
}
