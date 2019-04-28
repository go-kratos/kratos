package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
	"github.com/urfave/cli"
)

func runUpdate(root, param string) {
	cmd := exec.Command("make", "bazel-update")
	cmd.Dir = param
	cmd.Env = append(cmd.Env, "KRATOS_ROOT="+root, "GOPATH="+root, "PATH="+os.Getenv("PATH"))
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
}

func runProwUpdate(root, param string) {
	cmd := exec.Command("make", "prow-update")
	cmd.Dir = param
	cmd.Env = append(cmd.Env, "KRATOS_ROOT="+root, "GOPATH="+root, "PATH="+os.Getenv("PATH"))
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
}

func updateAction(c *cli.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	index := strings.Index(pwd, "go-common")
	if index == -1 {
		fmt.Println("not in go-common")
		os.Exit(1)
	}
	path := strings.Split(pwd[:index-1], "/")
	result := strings.Split(pwd[index:], "/")
	path = append(path, result[0])
	runPath := strings.Join(path, "/")
	runUpdate(strings.Join(path[:len(path)-2], "/"), runPath)
	return nil
}

func updateProwAction() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	index := strings.Index(pwd, "go-common")
	if index == -1 {
		fmt.Println("not in go-common")
		os.Exit(1)
	}
	path := strings.Split(pwd[:index-1], "/")
	result := strings.Split(pwd[index:], "/")
	path = append(path, result[0])
	runPath := strings.Join(path, "/")
	runProwUpdate(strings.Join(path[:len(path)-2], "/"), runPath)
	return nil
}
