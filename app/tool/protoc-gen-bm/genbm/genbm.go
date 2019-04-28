package genbm

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"

	"go-common/app/tool/protoc-gen-bm/generator"
)

// New blademaster server code generator
func New(jsonpb bool) generator.Generator {
	return &genbm{jsonpb: jsonpb}
}

type genbm struct {
	jsonpb bool
}

func (g *genbm) Generate(req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {
	var resp []*plugin.CodeGeneratorResponse_File
	files := req.GetProtoFile()
	for _, file := range files {
		respFile, ok, err := g.generateFile(file)
		if err != nil {
			return resp, err
		}
		if ok {
			resp = append(resp, respFile)
		}
	}
	return resp, nil
}

func (g *genbm) generateFile(file *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, bool, error) {
	glog.V(1).Infof("process proto file %s", file.GetName())
	services := file.GetService()
	if len(services) == 0 {
		glog.V(5).Infof("proto file %s not included service descriptor", file.GetName())
		return nil, false, nil
	}

	var descs []*BMServerDescriptor
	for _, service := range services {
		server, err := ParseBMServer(service)
		if err != nil {
			return nil, false, err
		}
		descs = append(descs, server)
	}

	buf := new(bytes.Buffer)
	goPackageName := GetGoPackageName(file)
	gen := NewBMGenerate(goPackageName, descs, g.jsonpb)
	if err := gen.Generate(buf); err != nil {
		return nil, false, err
	}

	// format code
	data, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, false, err
	}

	content := string(data)
	// no content
	if len(content) == 0 {
		return nil, false, nil
	}
	target := TargetFilePath(file)
	glog.V(1).Infof("generate code to %s", target)
	return &plugin.CodeGeneratorResponse_File{
		Content: &content,
		Name:    &target,
	}, true, nil
}

// TargetFilePath find target file path
func TargetFilePath(file *descriptor.FileDescriptorProto) string {
	fpath := file.GetName()
	protoDir := filepath.Dir(fpath)
	noExt := filepath.Base(fpath)
	for i := len(noExt) - 1; i >= 0 && !os.IsPathSeparator(noExt[i]); i-- {
		if noExt[i] == '.' {
			noExt = noExt[:i]
		}
	}
	target := noExt + ".pb.bm.go"
	options := file.GetOptions()
	if options != nil {
		goPackage := options.GetGoPackage()
		if goPackage != "" {
			goPackage = strings.Split(goPackage, ";")[0]
			if strings.Contains(goPackage, "/") {
				return filepath.Join(goPackage, target)
			}
		}
	}
	return filepath.Join(protoDir, target)
}

// GetGoPackageName last element from proto package name or go_package option
func GetGoPackageName(file *descriptor.FileDescriptorProto) string {
	var goPackageName string
	protoPackage := file.GetPackage()
	goPackageName = splitLastElem(protoPackage, ".")

	options := file.GetOptions()
	if options == nil {
		return goPackageName
	}
	if goPackage := options.GetGoPackage(); goPackage != "" {
		if strings.Contains(goPackage, ";") {
			goPackageName = splitLastElem(goPackage, ";")
		} else {
			goPackageName = splitLastElem(goPackage, "/")
		}
	}
	return goPackageName
}

func splitLastElem(s string, seq string) string {
	seqs := strings.Split(s, seq)
	return seqs[len(seqs)-1]
}
