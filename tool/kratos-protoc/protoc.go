package main

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

var (
	withBM      bool
	withGRPC    bool
	withSwagger bool
	withEcode   bool
)

func protocAction(ctx *cli.Context) (err error) {
	if err = checkProtoc(); err != nil {
		return err
	}
	files := []string(ctx.Args())
	if len(files) == 0 {
		files, _ = filepath.Glob("*.proto")
	}
	if !withGRPC && !withBM && !withSwagger && !withEcode {
		withBM = true
		withGRPC = true
		withSwagger = true
		withEcode = true
	}
	if withBM {
		if err = installBMGen(); err != nil {
			return
		}
		if err = genBM(files); err != nil {
			return
		}
	}
	if withGRPC {
		if err = installGRPCGen(); err != nil {
			return err
		}
		if err = genGRPC(files); err != nil {
			return
		}
	}
	if withSwagger {
		if err = installSwaggerGen(); err != nil {
			return
		}
		if err = genSwagger(files); err != nil {
			return
		}
	}
	if withEcode {
		if err = installEcodeGen(); err != nil {
			return
		}
		if err = genEcode(files); err != nil {
			return
		}
	}
	log.Printf("generate %s success.\n", strings.Join(files, " "))
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

func generate(protoc string, files []string) error {
	pwd, _ := os.Getwd()
	gosrc := path.Join(gopath(), "src")
	ext, err := latestKratos()
	if err != nil {
		return err
	}
	line := fmt.Sprintf(protoc, gosrc, ext, pwd)
	log.Println(line, strings.Join(files, " "))
	args := strings.Split(line, " ")
	args = append(args, files...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = pwd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	gopath := gopath()
	ext := path.Join(gopath, "src/github.com/bilibili/kratos/third_party")
	if _, err := os.Stat(ext); !os.IsNotExist(err) {
		return ext, nil
	}
	baseMod := path.Join(gopath, "pkg/mod/github.com/bilibili")
	files, err := ioutil.ReadDir(baseMod)
	if err != nil {
		return "", err
	}
	for i := len(files) - 1; i >= 0; i-- {
		if strings.HasPrefix(files[i].Name(), "kratos@") {
			return path.Join(baseMod, files[i].Name(), "third_party"), nil
		}
	}
	return "", errors.New("not found kratos package")
}

func gopath() (gp string) {
	gopaths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))

	if len(gopaths) == 1 && gopaths[0] != "" {
		return gopaths[0]
	}
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	abspwd, err := filepath.Abs(pwd)
	if err != nil {
		return
	}
	for _, gopath := range gopaths {
		if gopath == "" {
			continue
		}
		absgp, err := filepath.Abs(gopath)
		if err != nil {
			return
		}
		if strings.HasPrefix(abspwd, absgp) {
			return absgp
		}
	}
	return build.Default.GOPATH
}
