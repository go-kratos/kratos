package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "v2.0.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	omitempty := flag.Bool("omitempty", true, "omit if google.api is empty")

	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-http %v\n", version)
		return
	}

	//var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, *omitempty)
		}
		return nil
	})
}
