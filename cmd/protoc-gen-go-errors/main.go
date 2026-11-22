package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var showVersion = flag.Bool("version", false, "print the version and exit")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-errors %v\n", release)
		return
	}
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		setSupportedFeatures(gen)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

// setSupportedFeatures declares the protoc features this plugin handles.
// Editions files are accepted for the same range as protoc-gen-go.
func setSupportedFeatures(gen *protogen.Plugin) {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
		pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	gen.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO2
	gen.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2024
}
