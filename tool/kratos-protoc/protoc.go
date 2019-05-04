package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

var (
	withGRPC bool
	withBM   bool
)

func protocAction(ctx *cli.Context) (err error) {
	if err = checkProtoc(); err != nil {
		return err
	}
	if !withGRPC && !withBM {
		return errors.New("must be options: [--grpc] or [--bm]")
	}
	if withGRPC {
		if err = installGRPCGen(); err != nil {
			return err
		}
		if err = genGRPC(ctx); err != nil {
			return
		}
	}
	if withBM {
		if err = installBMGen(); err != nil {
			return
		}
		if err = genBM(ctx); err != nil {
			return
		}
	}
	log.Printf("generate %v success.\n", ctx.Args())
	return nil
}

func checkProtoc() error {
	if _, err := exec.LookPath("protoc"); err != nil {
		switch runtime.GOOS {
		case "darwin":
			fmt.Println("brew install protobuf")
			cmd := exec.Command("brew", "install", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		case "linux":
			fmt.Println("snap install --classic protobuf")
			cmd := exec.Command("snap", "install", "--classic", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		default:
			return errors.New("您还没安装protobuf，请进行手动安装：https://github.com/protocolbuffers/protobuf/releases")
		}
	}
	return nil
}
