package main

import (
	"os/exec"

	"github.com/urfave/cli"
)

const (
	_getBMGen = "go get -u github.com/bilibili/kratos/tool/protobuf/protoc-gen-bm"
	_bmProtoc = "protoc --proto_path=%s --proto_path=%s --proto_path=%s --bm_out=:."
)

func installBMGen() error {
	if _, err := exec.LookPath("protoc-gen-bm"); err != nil {
		if err := goget(_getBMGen); err != nil {
			return err
		}
	}
	return nil
}

func genBM(ctx *cli.Context) error {
	return generate(ctx, _bmProtoc)
}
