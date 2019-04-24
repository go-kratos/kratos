package main

import (
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

func buildAction(c *cli.Context) error {
	args := append([]string{"build"}, c.Args()...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return nil
}
