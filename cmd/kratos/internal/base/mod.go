package base

import (
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

// ModulePath returns go module path.
func ModulePath(filename string) (string, error) {
	modBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}
