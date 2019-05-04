package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

const (
	_installGRPCGen = "go get github.com/gogo/protobuf/protoc-gen-gogofast"
	_grpcProtoc     = "protoc --proto_path=%s/src --proto_path=%s --gogofast_out=plugins=grpc:."
)

func installGRPCGen() error {
	if _, err := exec.LookPath("protoc-gen-gogofast"); err != nil {
		log.Println(_installGRPCGen)
		args := strings.Split(_installGRPCGen, " ")
		cmd := exec.Command(args[0], args[1:]...)
		// 安装gogofast不能直接用go mod安装，因为需要安装到gopath，protoc将会依赖到gogo.proto
		for _, env := range os.Environ() {
			if strings.Contains(env, "GO111MODULE") {
				continue
			}
			cmd.Env = append(cmd.Env, env)
		}
		cmd.Env = append(cmd.Env, "GO111MODULE=off")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func genGRPC(ctx *cli.Context) error {
	pwd, _ := os.Getwd()
	protoc := fmt.Sprintf(_grpcProtoc, os.Getenv("GOPATH"), pwd)
	log.Println(protoc, strings.Join(ctx.Args(), " "))
	args := strings.Split(protoc, " ")
	args = append(args, ctx.Args()...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir, _ = os.Getwd()
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
