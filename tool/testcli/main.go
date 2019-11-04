package main

import (
	"flag"
	"os"
	"os/exec"
	"strings"

	"github.com/bilibili/kratos/pkg/testing/lich"
)

func parseArgs() (flags map[string]string) {
	flags = make(map[string]string)
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		if arg == "down" {
			flags["down"] = ""
			return
		}
		if cmds := os.Args[idx+1:]; arg == "run" {
			flags["run"] = strings.Join(cmds, " ")
			return
		}
	}
	return
}

func main() {
	flag.Parse()
	flags := parseArgs()
	if _, ok := flags["down"]; ok {
		lich.Teardown()
		return
	}
	if cmd, ok := flags["run"]; !ok || cmd == "" {
		panic("Your need 'run' flag assign to be run commands.")
	}
	if err := lich.Setup(); err != nil {
		panic(err)
	}
	defer lich.Teardown()
	cmds := strings.Split(flags["run"], " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
