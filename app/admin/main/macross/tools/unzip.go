package tools

import (
	"archive/zip"
	"go-common/library/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip unzip a file to the target directory.
func Unzip(filePath string, targetDir string) (err error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		log.Error("ozip.OpenReader(%s) error(%v)", filePath, err)
		return
	}
	defer reader.Close()

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		log.Error("os.MkdirAll(%s) error(%v)", targetDir, err)
		return
	}

	for _, file := range reader.File {
		path := filepath.Join(targetDir, file.Name)
		if strings.Contains(path, "__MACOSX") {
			continue
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		} else {
			os.MkdirAll(filepath.Dir(path), 0755)
		}

		fileReader, err := file.Open()
		if err != nil {
			log.Error("file.Open() error(%v)", err)
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			log.Error("os.OpenFile() error(%v)", err)
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			log.Error("io.Copy() error(%v)", err)
			return err
		}
	}
	return
}
