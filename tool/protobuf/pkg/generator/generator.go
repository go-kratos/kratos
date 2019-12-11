package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"strconv"
	"strings"

	"github.com/bilibili/kratos/tool/protobuf/pkg/gen"
	"github.com/bilibili/kratos/tool/protobuf/pkg/naming"
	"github.com/bilibili/kratos/tool/protobuf/pkg/typemap"
	"github.com/bilibili/kratos/tool/protobuf/pkg/utils"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

const Version = "v0.1"

var GoModuleImportPath = "github.com/bilibili/kratos"
var GoModuleDirName = "github.com/bilibili/kratos"

type Base struct {
	Reg *typemap.Registry

	// Map to record whether we've built each package
	// pkgName => alias name
	pkgs          map[string]string
	pkgNamesInUse map[string]bool

	ImportPrefix string            // String to prefix to imported package file names.
	importMap    map[string]string // Mapping from .proto file name to import path.

	// Package naming:
	GenPkgName          string // Name of the package that we're generating
	PackageName         string // Name of the proto file package
	fileToGoPackageName map[*descriptor.FileDescriptorProto]string

	// List of files that were inputs to the generator. We need to hold this in
	// the struct so we can write a header for the file that lists its inputs.
	GenFiles []*descriptor.FileDescriptorProto

	// Output buffer that holds the bytes we want to write out for a single file.
	// Gets reset after working on a file.
	Output *bytes.Buffer

	// key: pkgName
	// value: importPath
	Deps map[string]string

	Params *ParamsBase

	httpInfoCache map[string]*HTTPInfo
}

// RegisterPackageName name is the go package name or proto pkg name
// return go pkg alias
func (t *Base) RegisterPackageName(name string) (alias string) {
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

func (t *Base) Setup(in *plugin.CodeGeneratorRequest, paramsOpt ...GeneratorParamsInterface) {
	t.httpInfoCache = make(map[string]*HTTPInfo)
	t.pkgs = make(map[string]string)
	t.pkgNamesInUse = make(map[string]bool)
	t.importMap = make(map[string]string)
	t.Deps = make(map[string]string)
	t.fileToGoPackageName = make(map[*descriptor.FileDescriptorProto]string)
	t.Output = bytes.NewBuffer(nil)

	var params GeneratorParamsInterface
	if len(paramsOpt) > 0 {
		params = paramsOpt[0]
	} else {
		params = &BasicParam{}
	}
	err := ParseGeneratorParams(in.GetParameter(), params)
	if err != nil {
		gen.Fail("could not parse parameters", err.Error())
	}
	t.Params = params.GetBase()
	t.ImportPrefix = params.GetBase().ImportPrefix
	t.importMap = params.GetBase().ImportMap

	t.GenFiles = gen.FilesToGenerate(in)

	// Collect information on types.
	t.Reg = typemap.New(in.ProtoFile)
	t.RegisterPackageName("context")
	t.RegisterPackageName("ioutil")
	t.RegisterPackageName("proto")
	// Time to figure out package names of objects defined in protobuf. First,
	// we'll figure out the name for the package we're generating.
	genPkgName, err := DeduceGenPkgName(t.GenFiles)
	if err != nil {
		gen.Fail(err.Error())
	}
	t.GenPkgName = genPkgName
	// Next, we need to pick names for all the files that are dependencies.
	if len(in.ProtoFile) > 0 {
		t.PackageName = t.GenFiles[0].GetPackage()
	}

	for _, f := range in.ProtoFile {
		if fileDescSliceContains(t.GenFiles, f) {
			// This is a file we are generating. It gets the shared package name.
			t.fileToGoPackageName[f] = t.GenPkgName
		} else {
			// This is a dependency. Use its package name.
			name := f.GetPackage()
			if name == "" {
				name = utils.BaseName(f.GetName())
			}
			name = utils.CleanIdentifier(name)
			alias := t.RegisterPackageName(name)
			t.fileToGoPackageName[f] = alias
		}
	}

	for _, f := range t.GenFiles {
		deps := t.DeduceDeps(f)
		for k, v := range deps {
			t.Deps[k] = v
		}
	}
}

func (t *Base) DeduceDeps(file *descriptor.FileDescriptorProto) map[string]string {
	deps := make(map[string]string) // Map of package name to quoted import path.
	ourImportPath := path.Dir(naming.GoFileName(file, ""))
	for _, s := range file.Service {
		for _, m := range s.Method {
			defs := []*typemap.MessageDefinition{
				t.Reg.MethodInputDefinition(m),
				t.Reg.MethodOutputDefinition(m),
			}
			for _, def := range defs {
				if def.File.GetPackage() == t.PackageName {
					continue
				}
				// By default, import path is the dirname of the Go filename.
				importPath := path.Dir(naming.GoFileName(def.File, ""))
				if importPath == ourImportPath {
					continue
				}
				importPath = t.SubstituteImportPath(importPath, def.File.GetName())
				importPath = t.ImportPrefix + importPath
				pkg := t.GoPackageNameForProtoFile(def.File)
				deps[pkg] = strconv.Quote(importPath)
			}
		}
	}
	return deps
}

// DeduceGenPkgName figures out the go package name to use for generated code.
// Will try to use the explicit go_package setting in a file (if set, must be
// consistent in all files). If no files have go_package set, then use the
// protobuf package name (must be consistent in all files)
func DeduceGenPkgName(genFiles []*descriptor.FileDescriptorProto) (string, error) {
	var genPkgName string
	for _, f := range genFiles {
		name, explicit := naming.GoPackageName(f)
		if explicit {
			name = utils.CleanIdentifier(name)
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
		name, _ := naming.GoPackageName(f)
		name = utils.CleanIdentifier(name)
		if genPkgName != "" && genPkgName != name {
			return "", errors.Errorf("files have conflicting package names, must be the same or overridden with go_package: %q and %q", genPkgName, name)
		}
		genPkgName = name
	}

	// All the files have the same name, so we're good.
	return genPkgName, nil
}

func (t *Base) GoPackageNameForProtoFile(file *descriptor.FileDescriptorProto) string {
	return t.fileToGoPackageName[file]
}

func fileDescSliceContains(slice []*descriptor.FileDescriptorProto, f *descriptor.FileDescriptorProto) bool {
	for _, sf := range slice {
		if f == sf {
			return true
		}
	}
	return false
}

// P forwards to g.gen.P, which prints output.
func (t *Base) P(args ...string) {
	for _, v := range args {
		t.Output.WriteString(v)
	}
	t.Output.WriteByte('\n')
}

func (t *Base) FormattedOutput() string {
	// Reformat generated code.
	fset := token.NewFileSet()
	raw := t.Output.Bytes()
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

func (t *Base) PrintComments(comments typemap.DefinitionComments) bool {
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

// IsOwnPackage ...
// protoName is fully qualified name of a type
func (t *Base) IsOwnPackage(protoName string) bool {
	def := t.Reg.MessageDefinition(protoName)
	if def == nil {
		gen.Fail("could not find message for", protoName)
	}
	return def.File.GetPackage() == t.PackageName
}

// Given a protobuf name for a Message, return the Go name we will use for that
// type, including its package prefix.
func (t *Base) GoTypeName(protoName string) string {
	def := t.Reg.MessageDefinition(protoName)
	if def == nil {
		gen.Fail("could not find message for", protoName)
	}

	var prefix string
	if def.File.GetPackage() != t.PackageName {
		prefix = t.GoPackageNameForProtoFile(def.File) + "."
	}

	var name string
	for _, parent := range def.Lineage() {
		name += parent.Descriptor.GetName() + "_"
	}
	name += def.Descriptor.GetName()
	return prefix + name
}

func streamingMethod(method *descriptor.MethodDescriptorProto) bool {
	return (method.ServerStreaming != nil && *method.ServerStreaming) || (method.ClientStreaming != nil && *method.ClientStreaming)
}

func (t *Base) ShouldGenForMethod(file *descriptor.FileDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
	method *descriptor.MethodDescriptorProto) bool {
	if streamingMethod(method) {
		return false
	}
	if !t.Params.ExplicitHTTP {
		return true
	}
	httpInfo := t.GetHttpInfoCached(file, service, method)
	return httpInfo.HasExplicitHTTPPath
}
func (t *Base) SubstituteImportPath(importPath string, importFile string) string {
	if substitution, ok := t.importMap[importFile]; ok {
		importPath = substitution
	}
	return importPath
}
