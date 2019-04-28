package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"go-common/app/tool/warden/generator"
	"go-common/app/tool/warden/goparser"
	"go-common/app/tool/warden/types"
)

const (
	// GoCommon .
	GoCommon = "go-common"
)

var (
	name            string
	dir             string
	recvName        string
	workDir         string
	protoOut        string
	csCode          string
	goPackage       string
	protoPackage    string
	ignoreTypeError bool
	noprotoc        bool
	importPaths     string
)

func init() {
	flag.StringVar(&name, "name", "", "service name")
	flag.StringVar(&dir, "dir", "service", "service go code dir")
	flag.StringVar(&recvName, "recv", "Service", "receiver name")
	flag.StringVar(&workDir, "workdir", ".", "workdir")
	flag.StringVar(&csCode, "cs-code", "server/grpc", "server code directory")
	flag.StringVar(&protoOut, "proto-out", "api/api.proto", "proto file save path")
	flag.StringVar(&goPackage, "go-package", "", "go-package")
	flag.StringVar(&protoPackage, "proto-package", "", "proto-package")
	flag.BoolVar(&ignoreTypeError, "ignore-type-error", true, "ignore type error")
	flag.BoolVar(&noprotoc, "noprotoc", false, "don't run protoc")
	flag.StringVar(&importPaths, "proto-path", defaultImportPath(), "specify the directory in which to search for imports.")
}

func defaultImportPath() string {
	for _, goPath := range strings.Split(os.Getenv("GOPATH"), ":") {
		fixPath := path.Join(goPath, "src", GoCommon)
		if _, err := os.Stat(fixPath); err == nil {
			return fixPath
		}
	}
	return ""
}

func main() {
	var err error
	if !flag.Parsed() {
		flag.Parse()
	}
	if name == "" {
		log.Fatal("service name required")
	}
	var servicePackage string
	servicePackage, err = goparser.GoPackage(dir)
	if err != nil {
		log.Fatalf("auto detect gopackage error %s", err)
	}
	if goPackage == "" {
		// auto set go package
		goPackage = path.Join(path.Dir(servicePackage), csCode)
	}
	if protoPackage == "" {
		log.Fatal("proto package name required")
	}
	var spec *types.ServiceSpec
	spec, err = goparser.Parse(name, dir, recvName, workDir)
	if err != nil {
		log.Fatal(err)
	}

	var paths []string
	if importPaths != "" {
		paths = strings.Split(importPaths, ",")
	}
	options := &generator.ServiceProtoOptions{
		GoPackage:    goPackage,
		ProtoPackage: protoPackage,
		IgnoreType:   ignoreTypeError,
		ImportPaths:  paths,
	}

	protoFile := path.Join(workDir, protoOut)
	if err = os.MkdirAll(path.Dir(protoFile), 0755); err != nil {
		log.Print(err)
	}
	var protoFp *os.File
	protoFp, err = os.OpenFile(protoFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer protoFp.Close()
	if err := generator.GenServiceProto(protoFp, spec, options); err != nil {
		log.Fatal(err)
	}

	if !noprotoc {
		if err := generator.Protoc(protoFile, "", "", paths); err != nil {
			log.Fatal(err)
		}
	}

	csOptions := &generator.GenCSCodeOptions{
		PbPackage:   path.Join(path.Dir(servicePackage), path.Dir(protoOut)),
		RecvName:    recvName,
		RecvPackage: servicePackage,
	}
	if err := generator.GenCSCode(csCode, spec, csOptions); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(`
🍺  (゜-゜)つロ 干杯~ !
proto file: %s
server: %s
`, protoFile, path.Join(csCode, "server.go"))
}
