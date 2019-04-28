// Package generator provides an abstract interface to code generators.
package generator

import (
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Generator is an abstraction of code generators.
type Generator interface {
	// Generate generates output files from input .proto files.
	Generate(req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error)
}
