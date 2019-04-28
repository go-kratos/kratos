package codegenerator

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// ParseRequest parses a code generator request from a proto Message.
func ParseRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read code generator request: %v", err)
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(input, req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code generator request: %v", err)
	}
	return req, nil
}

// WriteResponse write a code generator response
func WriteResponse(w io.Writer, files []*plugin.CodeGeneratorResponse_File, inErr error) error {
	var perrMsg *string
	if inErr != nil {
		errMsg := inErr.Error()
		perrMsg = &errMsg
	}
	resp := &plugin.CodeGeneratorResponse{Error: perrMsg, File: files}
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}
