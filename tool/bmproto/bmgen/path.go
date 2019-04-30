package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func fileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func initGopath() string {
	root, err := goPath()
	if err != nil || root == "" {
		logf("can not read GOPATH, use ~/go as default GOPATH")
		root = path.Join(os.Getenv("HOME"), "go")
	}
	return root
}

func goPath() (string, error) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	if len(gopaths) == 1 {
		return gopaths[0], nil
	}
	for _, gp := range gopaths {
		absgp, err := filepath.Abs(gp)
		if err != nil {
			return "", err
		}
		if fileExist(absgp + "/src/github.com/bilibili/kratos/") || fileExist(absgp + "/src/kratos/") {
			return absgp, nil
		}
	}
	return "", fmt.Errorf("can't found current gopath")
}

// files all files with suffix
func filesWithSuffix(dir string, suffixes ...string) (result []string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			for _, suffix := range suffixes {
				if strings.HasSuffix(info.Name(), suffix) {
					result = append(result, path)
					break
				}
			}
		}
		return nil
	})
	return
}
