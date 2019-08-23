package naming

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bilibili/kratos/tool/protobuf/pkg/utils"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/pkg/errors"
	"github.com/siddontang/go/ioutil2"
)

// GetVersionPrefix 根据go包名获取api版本前缀
// @param pkg 从proto获取到的对应的go报名
// @return 如果是v*开始的 返回v*
// 否则返回空
func GetVersionPrefix(pkg string) string {
	if pkg == "" {
		return ""
	}
	if pkg[:1] == "v" {
		return pkg
	}
	return ""
}

// GenFileName returns the output name for the generated Go file.
func GenFileName(f *descriptor.FileDescriptorProto, suffix string) string {
	name := *f.Name
	if ext := path.Ext(name); ext == ".pb" || ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}
	name += suffix
	return name
}

func ServiceName(service *descriptor.ServiceDescriptorProto) string {
	return utils.CamelCase(service.GetName())
}

// MethodName ...
func MethodName(method *descriptor.MethodDescriptorProto) string {
	return utils.CamelCase(method.GetName())
}

// GetGoImportPathForPb 得到 proto 文件对应的 go import路径
// protoFilename is the proto file name
// 可能根本无法得到proto文件的具体路径, 只能假设 proto 的filename 是相对当前目录的
// 假设 protoAbsolutePath = wd/protoFilename
func GetGoImportPathForPb(protoFilename string, moduleImportPath string, moduleDirName string) (importPath string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		panic("cannot get working directory")
	}
	absPath := wd + "/" + protoFilename
	if !ioutil2.FileExists(absPath) {
		err = errors.New("Cannot find proto file path of " + protoFilename)
		return "", err
	}
	index := strings.Index(absPath, moduleDirName)
	if index == -1 {
		return "", errors.Errorf("proto file %s is not inside project %s", protoFilename, moduleDirName)
	}
	relativePath := absPath[index:]
	importPath = filepath.Dir(relativePath)
	return importPath, nil
}

// GoPackageNameForProtoFile returns the Go package name to use in the generated Go file.
// The result explicitly reports whether the name came from an option go_package
// statement. If explicit is false, the name was derived from the protocol
// buffer's package statement or the input file name.
func GoPackageName(f *descriptor.FileDescriptorProto) (name string, explicit bool) {
	// Does the file have a "go_package" option?
	if _, pkg, ok := goPackageOption(f); ok {
		return pkg, true
	}

	// Does the file have a package clause?
	if pkg := f.GetPackage(); pkg != "" {
		return pkg, false
	}
	// Use the file base name.
	return utils.BaseName(f.GetName()), false
}

// goPackageOption interprets the file's go_package option.
// If there is no go_package, it returns ("", "", false).
// If there's a simple name, it returns ("", pkg, true).
// If the option implies an import path, it returns (impPath, pkg, true).
func goPackageOption(f *descriptor.FileDescriptorProto) (impPath, pkg string, ok bool) {
	pkg = f.GetOptions().GetGoPackage()
	if pkg == "" {
		return
	}
	ok = true
	// The presence of a slash implies there's an import path.
	slash := strings.LastIndex(pkg, "/")
	if slash < 0 {
		return
	}
	impPath, pkg = pkg, pkg[slash+1:]
	// A semicolon-delimited suffix overrides the package name.
	sc := strings.IndexByte(impPath, ';')
	if sc < 0 {
		return
	}
	impPath, pkg = impPath[:sc], impPath[sc+1:]
	return
}
