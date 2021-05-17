package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
)

var _ config.Source = (*file)(nil)

type file struct {
	path string
}

// NewSource new a file source.
func NewSource(path string) config.Source {
	return &file{path: path}
}

func (f *file) loadFile(path string) (*config.KeyValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &config.KeyValue{
		Key:    info.Name(),
		Format: format(info.Name()),
		Value:  data,
	}, nil
}

func (f *file) loadDir(path string) (kvs []*config.KeyValue, err error) {
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// ignore hidden files
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		kv, err := f.loadFile(filepath.Join(f.path, file.Name()))
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
	}
	return
}

func (f *file) Load() (kvs []*config.KeyValue, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	return []*config.KeyValue{kv}, nil
}

func (f *file) Watch() (config.Watcher, error) {
	return newWatcher(f)
}
