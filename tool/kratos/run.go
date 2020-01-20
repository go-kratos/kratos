package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func runAction(c *cli.Context) error {
	base, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := buildDir(base, "cmd", 5)
	conf := path.Join(filepath.Dir(dir), "configs")
	args := append([]string{"run", "main.go", "-conf", conf}, c.Args().Slice()...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return nil
}
