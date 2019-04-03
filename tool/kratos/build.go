package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"github.com/urfave/cli"
)

func buildAction(c *cli.Context) error {
	args := append([]string{"build"}, c.Args()...)
	// fmt.Println(args)
	cmd := exec.Command("go", args...)
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
	return nil
}
