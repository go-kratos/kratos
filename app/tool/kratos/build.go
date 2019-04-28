package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kr/pty"
	"github.com/urfave/cli"
)

func runbazel(param ...string) {
	command := append([]string{"build", "--watchfs"}, param...)
	fmt.Println(command)
	cmd := exec.Command("bazel", command...)
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
}

func bazelAction(c *cli.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	index := strings.Index(pwd, "go-common")
	if index == -1 {
		fmt.Println("not in go-common")
		os.Exit(1)
	}
	result := strings.Split(pwd[index:], "/")
	runPath := strings.Join(result[1:], "/")
	if c.NArg() > 0 {
		param := []string{}
		for index := 0; index < c.NArg(); index++ {
			name := path.Join(runPath, path.Clean(c.Args().Get(index)))
			if name == "." {
				continue
			}
			if strings.HasSuffix(name, "/...") {
				param = append(param, "//"+name)
			} else {
				param = append(param, "//"+name+"/...")
			}

		}
		runbazel(param...)
	} else {
		if len(runPath) == 0 {
			runbazel("//app/...", "//library/...", "//vendor/...")
		} else {
			runbazel("//" + strings.Join(result[1:], "/") + "/...")
		}
	}
	return nil
}
