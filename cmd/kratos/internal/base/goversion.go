package base

import (
	"os/exec"
	"regexp"
)

func GetGoVersion() (string, error) {
	output, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", nil
	}

	r, err := regexp.Compile(`\d+\.\d+\.\d+`)
	if err != nil {
		return "", err
	}
	return r.FindString(string(output)), nil
}
