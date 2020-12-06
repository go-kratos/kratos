package base

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func kratosHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := path.Join(dir, ".kratos")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func kratosHomeWithDir(dir string) string {
	home := path.Join(kratosHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func copyDir(src, dst string, ignores []string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}

		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp, ignores); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}
