package input

import (
	"fmt"
	"path"
	"path/filepath"

	"go-common/app/tool/gorpc/model"
)

// Returns all the Golang files for the given path. Ignores hidden files.
func Files(srcPath string) ([]model.Path, error) {
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs: %v\n", err)
	}
	if filepath.Ext(srcPath) == "" {
		return dirFiles(srcPath)
	}
	return file(srcPath)
}

func dirFiles(srcPath string) ([]model.Path, error) {
	ps, err := filepath.Glob(path.Join(srcPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob: %v\n", err)
	}
	var srcPaths []model.Path
	for _, p := range ps {
		src := model.Path(p)
		if isHiddenFile(p) || src.IsTestPath() {
			continue
		}
		srcPaths = append(srcPaths, src)
	}
	return srcPaths, nil
}

func file(srcPath string) ([]model.Path, error) {
	src := model.Path(srcPath)
	if filepath.Ext(srcPath) != ".go" || isHiddenFile(srcPath) || src.IsTestPath() {
		return nil, fmt.Errorf("no Go source files found at %v", srcPath)
	}
	return []model.Path{src}, nil
}

func isHiddenFile(path string) bool {
	return []rune(filepath.Base(path))[0] == '.'
}
