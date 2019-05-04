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
	_installBMGen = "go get github.com/bilibili/kratos/tool/protoc-gen-bm"
	_bmProtoc     = "protoc --proto_path=%s/src --proto_path=%s --bm_out=."
)

func installBMGen() error {
	if _, err := exec.LookPath("protoc-gen-bm"); err != nil {
		log.Println(_installGRPCGen)
		args := strings.Split(_installGRPCGen, " ")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
func genBM(ctx *cli.Context) error {
	pwd, _ := os.Getwd()
	protoc := fmt.Sprintf(_bmProtoc, os.Getenv("GOPATH"), pwd)
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
