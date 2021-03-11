package base

import (
	"fmt"
	"os"
	"os/exec"
)

// GoGet go get path.
func GoGet(path ...string) error {
	for _, p := range path {
		fmt.Printf("go get -u %s\n", p)
		cmd := exec.Command("go", "get", "-u", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
