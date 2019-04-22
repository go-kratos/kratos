package main

import (
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"strings"
)

// runCmd runs the cmd & print output (both stdout & stderr)
func runCmd(cmd string) (err error) {
	out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	if flagVerbose || err != nil {
		logFunc := logf
		if err != nil {
			logFunc = errorf
		}
		logFunc("CMD: %s", cmd)
		logFunc(string(out))
	}
	return
}

func runCmdRet(cmd string) (out string, err error) {
	if flagVerbose {
		logf("CMD: %s", cmd)
	}
	outBytes, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	out = strings.Trim(string(outBytes), "\n\r\t ")
	return
}

func infof(format string, args ...interface{}) {
	color.Green(format, args...)
}

func logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func errorf(format string, args ...interface{}) {
	color.Red(format, args...)
}
