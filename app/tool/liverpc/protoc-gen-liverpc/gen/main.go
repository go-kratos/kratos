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

package gen

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Generator ...
type Generator interface {
	Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse
}

// Main ...
func Main(g Generator) {
	req := readGenRequest()
	resp := g.Generate(req)
	writeResponse(os.Stdout, resp)
}

// FilesToGenerate ...
func FilesToGenerate(req *plugin.CodeGeneratorRequest) []*descriptor.FileDescriptorProto {
	genFiles := make([]*descriptor.FileDescriptorProto, 0)
Outer:
	for _, name := range req.FileToGenerate {
		for _, f := range req.ProtoFile {
			if f.GetName() == name {
				genFiles = append(genFiles, f)
				continue Outer
			}
		}
		Fail("could not find file named", name)
	}

	return genFiles
}

func readGenRequest() *plugin.CodeGeneratorRequest {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Error(err, "reading input")
	}

	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		Error(err, "parsing input proto")
	}

	if len(req.FileToGenerate) == 0 {
		Fail("no files to generate")
	}

	return req
}

func writeResponse(w io.Writer, resp *plugin.CodeGeneratorResponse) {
	data, err := proto.Marshal(resp)
	if err != nil {
		Error(err, "marshaling response")
	}
	_, err = w.Write(data)
	if err != nil {
		Error(err, "writing response")
	}
}
