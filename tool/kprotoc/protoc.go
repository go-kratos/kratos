package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

const (
	_installProtocGen = "go get github.com/gogo/protobuf/protoc-gen-gogofast"
	_protocGRPC       = "protoc --proto_path=%s/src --proto_path=%s --gogofast_out=plugins=grpc:."
)

func protocAction(c *cli.Context) error {
	if _, err := exec.LookPath("protoc"); err != nil {
		switch runtime.GOOS {
		case "darwin":
			// MacOS: install protobuf
			fmt.Println("brew install protobuf..")
			cmd := exec.Command("brew", "install", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		default:
			return errors.New("您还没安装protobuf，请进行手动安装：https://github.com/protocolbuffers/protobuf/releases")
		}
	}
	if err := installProtocGen(); err != nil {
		return err
	}
	args := strings.Split(protocGRPC(), " ")
	args = append(args, c.Args()...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir, _ = os.Getwd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Printf("generate %v success.\n", c.Args())
	return nil
}

func installProtocGen() error {
	if _, err := exec.LookPath("protoc-gen-gogofast"); err != nil {
		log.Println(_installProtocGen)
		args := strings.Split(_installProtocGen, " ")
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

func protocGRPC() string {
	source, _ := os.Getwd()
	cmd := fmt.Sprintf(_protocGRPC, os.Getenv("GOPATH"), source)
	log.Println(cmd)
	return cmd
}
