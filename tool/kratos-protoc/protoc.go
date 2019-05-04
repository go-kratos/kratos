package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

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

func generate(ctx *cli.Context, protoc string) error {
	pwd, _ := os.Getwd()
	gosrc := path.Join(os.Getenv("GOPATH"), "src")
	ext, err := latestKratos()
	if err != nil {
		return err
	}
	line := fmt.Sprintf(protoc, gosrc, ext, pwd)
	log.Println(line, strings.Join(ctx.Args(), " "))
	args := strings.Split(line, " ")
	args = append(args, ctx.Args()...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = pwd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func goget(url string) error {
	args := strings.Split(url, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(url)
	return cmd.Run()
}

func latestKratos() (string, error) {
	gopath := os.Getenv("GOPATH")
	ext := path.Join(gopath, "src/github.com/bilibili/kratos/tool/protobuf/extensions")
	if _, err := os.Stat(ext); !os.IsNotExist(err) {
		return ext, nil
	}
	ext = path.Join(gopath, "pkg/mod/github.com/bilibili")
	files, err := ioutil.ReadDir(ext)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("not found kratos package")
	}
	return path.Join(ext, files[len(files)-1].Name(), "tool/protobuf/extensions"), nil
}
