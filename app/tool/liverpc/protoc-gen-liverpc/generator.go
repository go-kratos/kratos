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
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"strconv"
	"strings"

	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen"
	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen/stringutils"
	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen/typemap"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

type liverpc struct {
	filesHandled int

	reg *typemap.Registry

	// Map to record whether we've built each package
	pkgs          map[string]string
	pkgNamesInUse map[string]bool

	importPrefix string            // String to prefix to imported package file names.
	importMap    map[string]string // Mapping from .proto file name to import path.

	// Package naming:
	genPkgName          string // Name of the package that we're generating
	fileToGoPackageName map[*descriptor.FileDescriptorProto]string

	// List of files that were inputs to the generator. We need to hold this in
	// the struct so we can write a header for the file that lists its inputs.
	genFiles []*descriptor.FileDescriptorProto

	// Output buffer that holds the bytes we want to write out for a single file.
	// Gets reset after working on a file.
	output *bytes.Buffer
}

func liveRPCGenerator() *liverpc {
	t := &liverpc{
		pkgs:                make(map[string]string),
		pkgNamesInUse:       make(map[string]bool),
		importMap:           make(map[string]string),
		fileToGoPackageName: make(map[*descriptor.FileDescriptorProto]string),
		output:              bytes.NewBuffer(nil),
	}

	return t
}

func (t *liverpc) Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	params, err := parseCommandLineParams(in.GetParameter())
	if err != nil {
		gen.Fail("could not parse parameters passed to --liverpc_out", err.Error())
	}
	t.importPrefix = params.importPrefix
	t.importMap = params.importMap

	t.genFiles = gen.FilesToGenerate(in)

	// Collect information on types.
	t.reg = typemap.New(in.ProtoFile)

	t.registerPackageName("context")
	t.registerPackageName("ioutil")
	t.registerPackageName("proto")
	t.registerPackageName("liverpc")

	// Time to figure out package names of objects defined in protobuf. First,
	// we'll figure out the name for the package we're generating.
	genPkgName, err := deduceGenPkgName(t.genFiles)
	if err != nil {
		gen.Fail(err.Error())
	}
	t.genPkgName = genPkgName

	// Next, we need to pick names for all the files that are dependencies.
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
	var servicesNames []string
	for _, f := range t.genFiles {
		respFile := t.generate(f)
		for _, s := range f.Service {
			servicesNames = append(servicesNames, *s.Name)
		}
		if respFile != nil {
			resp.File = append(resp.File, respFile)
		}
	}
	// generate a temp file of service names
	// because a protobuf plugin can only generate for a single package
	// therefore we generate these temp files for other script to combine
	// a single client for all packages
	var filename = "client." + genPkgName + ".txt"
	var respFile = &plugin.CodeGeneratorResponse_File{}
	respFile.Name = &filename
	var content = strings.Join(servicesNames, "\n")
	content += "\n"
	respFile.Content = &content
	resp.File = append(resp.File, respFile)
	return resp
}

func (t *liverpc) registerPackageName(name string) (alias string) {
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

func (t *liverpc) generate(file *descriptor.FileDescriptorProto) *plugin.CodeGeneratorResponse_File {
	resp := new(plugin.CodeGeneratorResponse_File)
	if len(file.Service) == 0 {
		return nil
	}

	t.generateFileHeader(file)

	t.generateImports(file)
	if t.filesHandled == 0 {
		t.generateUtilImports()
	}

	// For each service, generate client stubs and server
	for i, service := range file.Service {
		t.generateService(file, service, i)
	}

	// Util functions only generated once per package
	if t.filesHandled == 0 {
		t.generateUtils()
	}

	t.generateFileDescriptor(file)

	resp.Name = proto.String(goFileName(file))
	resp.Content = proto.String(t.formattedOutput())
	t.output.Reset()

	t.filesHandled++
	return resp
}

func (t *liverpc) generateFileHeader(file *descriptor.FileDescriptorProto) {
	t.P("// Code generated by protoc-gen-liverpc ", gen.Version, ", DO NOT EDIT.")
	t.P("// source: ", file.GetName())
	t.P()
	if t.filesHandled == 0 {
		t.P("/*")
		t.P("Package ", t.genPkgName, " is a generated liverpc stub package.")
		t.P("This code was generated with go-common/app/tool/liverpc/protoc-gen-liverpc ", gen.Version, ".")
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
	t.P(`package `, t.genPkgName)
	t.P()
}

func (t *liverpc) generateImports(file *descriptor.FileDescriptorProto) {
	if len(file.Service) == 0 {
		return
	}
	t.P(`import `, t.pkgs["context"], ` "context"`)
	t.P()
	t.P(`import `, t.pkgs["proto"], ` "github.com/golang/protobuf/proto"`)
	t.P(`import "go-common/library/net/rpc/liverpc"`)
	t.P()

	// It's legal to import a message and use it as an input or output for a
	// method. Make sure to import the package of any such message. First, dedupe
	// them.
	deps := make(map[string]string) // Map of package name to quoted import path.
	ourImportPath := path.Dir(goFileName(file))
	for _, s := range file.Service {
		for _, m := range s.Method {
			defs := []*typemap.MessageDefinition{
				t.reg.MethodInputDefinition(m),
				t.reg.MethodOutputDefinition(m),
			}
			for _, def := range defs {
				// By default, import path is the dirname of the Go filename.
				importPath := path.Dir(goFileName(def.File))
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
	for pkg, importPath := range deps {
		t.P(`import `, pkg, ` `, importPath)
	}
	if len(deps) > 0 {
		t.P()
	}
	t.P(`var _ proto.Message // generate to suppress unused imports`)
}

func (t *liverpc) generateUtilImports() {
	t.P("// Imports only used by utility functions:")
	//t.P(`import `, t.pkgs["io"], ` "io"`)
	//t.P(`import `, t.pkgs["strconv"], ` "strconv"`)
	//t.P(`import `, t.pkgs["json"], ` "encoding/json"`)
	//t.P(`import `, t.pkgs["url"], ` "net/url"`)
}

// Generate utility functions used in LiveRpc code.
// These should be generated just once per package.
func (t *liverpc) generateUtils() {
	t.sectionComment(`Utils`)
	t.P(`func doRPCRequest(ctx `, t.pkgs["context"], `.Context, client *liverpc.Client, version int, method string, in, out `, t.pkgs["proto"], `.Message, opts []liverpc.CallOption) (err error) {`)
	t.P(`    err = client.Call(ctx, version, method, in, out, opts...)`)
	t.P(`    return`)
	t.P(`}`)
	t.P()
}

// P forwards to g.gen.P, which prints output.
func (t *liverpc) P(args ...string) {
	for _, v := range args {
		t.output.WriteString(v)
	}
	t.output.WriteByte('\n')
}

// Big header comments to makes it easier to visually parse a generated file.
func (t *liverpc) sectionComment(sectionTitle string) {
	t.P()
	t.P(`// `, strings.Repeat("=", len(sectionTitle)))
	t.P(`// `, sectionTitle)
	t.P(`// `, strings.Repeat("=", len(sectionTitle)))
	t.P()
}

func (t *liverpc) generateService(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto, index int) {
	servName := serviceName(service)

	t.sectionComment(servName + ` Interface`)
	t.generateLiveRPCInterface(file, service)

	t.sectionComment(servName + ` Live Rpc Client`)
	t.generateClient(file, service)

}

func (t *liverpc) generateLiveRPCInterface(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	comments, err := t.reg.ServiceComments(file, service)
	if err == nil {
		t.printComments(comments)
	}
	t.P(`type `, clientName(service), ` interface {`)
	for _, method := range service.Method {
		comments, err = t.reg.MethodComments(file, service, method)
		if err == nil {
			t.printComments(comments)
		}
		t.P(t.generateSignature(method))
		t.P()
	}
	t.P(`}`)
}

func (t *liverpc) generateSignature(method *descriptor.MethodDescriptorProto) string {
	methName := methodName(method)
	inputBodyType := t.goTypeName(method.GetInputType())
	outputType := t.goTypeName(method.GetOutputType())
	return fmt.Sprintf(`	%s(ctx %s.Context, req *%s, opts ...liverpc.CallOption) (resp *%s, err error)`, methName, t.pkgs["context"], inputBodyType, outputType)
}

// valid names: 'JSON', 'Protobuf'
func (t *liverpc) generateClient(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	clientName := clientName(service)
	structName := unexported(clientName)
	newClientFunc := "New" + clientName

	t.P(`type `, structName, ` struct {`)
	t.P(`  client *liverpc.Client`)
	t.P(`}`)
	t.P()
	t.P(`// `, newClientFunc, ` creates a client that implements the `, clientName, ` interface.`)
	t.P(`func `, newClientFunc, `(client *liverpc.Client) `, clientName, ` {`)
	t.P(`  return &`, structName, `{`)
	t.P(`    client: client,`)
	t.P(`  }`)
	t.P(`}`)
	t.P()

	for _, method := range service.Method {
		methName := methodName(method)
		pkgName := pkgName(file)

		inputType := t.goTypeName(method.GetInputType())
		outputType := t.goTypeName(method.GetOutputType())

		parts := strings.Split(pkgName, ".")
		if len(parts) < 2 {
			panic("package name must contain at least to parts, eg: service.v1, get " + pkgName + "!")
		}
		vStr := parts[len(parts)-1]
		if len(vStr) < 2 {
			panic("package name must contain a valid version, eg: service.v1")
		}
		_, err := strconv.Atoi(vStr[1:])
		if err != nil {
			panic("package name must contain a valid version, eg: service.v1, get " + vStr)
		}

		rpcMethod := method.GetName()
		rpcCtrl := service.GetName()
		rpcCmd := rpcCtrl + "." + rpcMethod

		t.P(`func (c *`, structName, `) `, methName, `(ctx `, t.pkgs["context"], `.Context, in *`, inputType, `, opts ...liverpc.CallOption) (*`, outputType, `, error) {`)
		t.P(`  out := new(`, outputType, `)`)
		t.P(`  err := doRPCRequest(ctx,c.client, `, vStr[1:], `, "`, rpcCmd, `", in, out, opts)`)
		t.P(`  if err != nil {`)
		t.P(`    return nil, err`)
		t.P(`  }`)
		t.P(`  return out, nil`)
		t.P(`}`)
		t.P()
	}
}

func (t *liverpc) generateFileDescriptor(file *descriptor.FileDescriptorProto) {
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

func (t *liverpc) printComments(comments typemap.DefinitionComments) bool {
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
func (t *liverpc) goTypeName(protoName string) string {
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

func (t *liverpc) goPackageName(file *descriptor.FileDescriptorProto) string {
	return t.fileToGoPackageName[file]
}

func (t *liverpc) formattedOutput() string {
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

func unexported(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func pkgName(file *descriptor.FileDescriptorProto) string {
	return file.GetPackage()
}

func serviceName(service *descriptor.ServiceDescriptorProto) string {
	return stringutils.CamelCase(service.GetName())
}

func clientName(service *descriptor.ServiceDescriptorProto) string {
	return serviceName(service) + "RPCClient"
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
