package main

import (
	"os/exec"
)

const (
	_getGRPCGen = "go get -u github.com/gogo/protobuf/protoc-gen-gofast"
	_grpcProtoc = "protoc --proto_path=%s --proto_path=%s --proto_path=%s --gofast_out=plugins=grpc," +
		"Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types," +
		"Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:."
)

func installGRPCGen() error {
	if _, err := exec.LookPath("protoc-gen-gofast"); err != nil {
		if err := goget(_getGRPCGen); err != nil {
			return err
		}
	}
	return nil
}

func genGRPC(files []string) error {
	return generate(_grpcProtoc, files)
}
