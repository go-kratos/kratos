package main

import (
	"os/exec"

	"github.com/urfave/cli"
)

const (
	_getGRPCGen = "go get github.com/gogo/protobuf/protoc-gen-gofast"
	_grpcProtoc = `protoc --proto_path=%s --proto_path=%s --proto_path=%s --gofast_out=plugins=grpc,Mgithub.com/gogo/protobuf/gogoproto/gogo.proto=github.com/bilibili/kratos/third_party/github.com/gogo/protobuf/gogoproto:.`
)

func installGRPCGen() error {
	if _, err := exec.LookPath("protoc-gen-gofast"); err != nil {
		if err := goget(_getGRPCGen); err != nil {
			return err
		}
	}
	return nil
}

func genGRPC(ctx *cli.Context) error {
	return generate(ctx, _grpcProtoc)
}
