package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bilibili/kratos/tool/pkg"
	"github.com/urfave/cli"
)

func runNew(ctx *cli.Context) (err error) {
	if p.onlyGRPC && p.onlyHTTP {
		p.onlyGRPC = false
		p.onlyHTTP = false
	}
	if p.path != "" {
		if p.path, err = filepath.Abs(p.path); err != nil {
			return
		}
		p.path = filepath.Join(p.path, p.Name)
	} else {
		pwd, _ := os.Getwd()
		p.path = filepath.Join(pwd, p.Name)
	}
	p.ModPrefix = modPath(p.path)
	// creata a project
	if err := create(); err != nil {
		return err
	}
	fmt.Printf("Project: %s\n", p.Name)
	fmt.Printf("OnlyGRPC: %t\n", p.onlyGRPC)
	fmt.Printf("OnlyHTTP: %t\n", p.onlyHTTP)
	fmt.Printf("Directory: %s\n\n", p.path)
	fmt.Println("项目创建成功.")
	return nil
}

func modPath(p string) string {
	dir := filepath.Dir(p)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			content, _ := ioutil.ReadFile(filepath.Join(dir, "go.mod"))
			mod := pkg.RegexpReplace(`module\s+(?P<name>[\S]+)`, string(content), "$name")
			return fmt.Sprintf("%s/%s/", mod, strings.TrimPrefix(filepath.Dir(p), dir+string(os.PathSeparator)))
		}
		parent := filepath.Dir(dir)
		if dir == parent {
			return ""
		}
		dir = parent
	}
}
