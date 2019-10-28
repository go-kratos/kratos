package main

import (
	"os/exec"
)

const (
	_getGRPCGen = "go get -u github.com/gogo/protobuf/protoc-gen-gofast"
	_grpcProtoc = `protoc --proto_path=%s --proto_path=%s --proto_path=%s --gofast_out=plugins=grpc:.`
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
