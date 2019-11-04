package main

import (
	"os/exec"
)

const (
	_getSwaggerGen = "go get -u github.com/bilibili/kratos/tool/protobuf/protoc-gen-bswagger"
	_swaggerProtoc = "protoc --proto_path=%s --proto_path=%s --proto_path=%s --bswagger_out=:."
)

func installSwaggerGen() error {
	if _, err := exec.LookPath("protoc-gen-bswagger"); err != nil {
		if err := goget(_getSwaggerGen); err != nil {
			return err
		}
	}
	return nil
}

func genSwagger(files []string) error {
	return generate(_swaggerProtoc, files)
}
