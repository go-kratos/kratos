package base

import (
	"os"
	"os/exec"
)

// GoGet go get path.
func Clone(layout string,name string) error {
	cmd := exec.Command("git", "clone", layout, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}