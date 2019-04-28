package store

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// FileConfig config of File.
type FileConfig struct {
	RootDir string
}

// File is a degrade file service.
type File struct {
	c *FileConfig
}

var _ Store = &File{}

// NewFile new a file degrade service.
func NewFile(fc *FileConfig) *File {
	if fc == nil {
		panic(errors.New("file config is nil"))
	}
	fs := &File{c: fc}
	if err := os.MkdirAll(fs.c.RootDir, 0755); err != nil {
		panic(errors.Wrapf(err, "dir: %s", fs.c.RootDir))
	}
	return fs
}

// Set save the result of location to file.
// expire is not implemented in file storage.
func (fs *File) Set(ctx context.Context, key string, bs []byte, _ int32) (err error) {
	file := path.Join(fs.c.RootDir, key)
	tmp := file + ".tmp"
	if err = ioutil.WriteFile(tmp, bs, 0644); err != nil {
		log.Error("ioutil.WriteFile(%s, bs, 0644): error(%v)", tmp, err)
		err = errors.Wrapf(err, "key: %s", key)
		return
	}
	if err = os.Rename(tmp, file); err != nil {
		log.Error("os.Rename(%s, %s): error(%v)", tmp, file, err)
		err = errors.Wrapf(err, "key: %s", key)
		return
	}
	return
}

// Get get result from file by locaiton+params.
func (fs *File) Get(ctx context.Context, key string) (bs []byte, err error) {
	p := path.Join(fs.c.RootDir, key)
	if bs, err = ioutil.ReadFile(p); err != nil {
		log.Error("ioutil.ReadFile(%s): error(%v)", p, err)
		err = errors.Wrapf(err, "key: %s", key)
		return
	}
	return
}
