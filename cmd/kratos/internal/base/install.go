//go:build go1.17
// +build go1.17

package base

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// GoInstall go get path.
func GoInstall(path ...string) error {
	reg := regexp.MustCompile(`.*@v\d+[\.\d+]+$`)
	for _, p := range path {
		if !reg.MatchString(p) {
			p += "@latest"
		}
		fmt.Printf("go install %s\n", p)
		cmd := exec.Command("go", "install", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
