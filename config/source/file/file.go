package file

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/config/source"
)

type file struct {
	path string
}

// New new a file source.
func New(path string) source.Source {
	return &file{path: path}
}

func (f *file) loadFile(name string) (*source.KeyValue, error) {
	file, err := os.Open(path.Join(f.path, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &source.KeyValue{
		Key:       name,
		Value:     data,
		Format:    format(name),
		Timestamp: info.ModTime(),
	}, nil
}

func (f *file) Load() (kvs []*source.KeyValue, err error) {
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		kv, err := f.loadFile(file.Name())
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
	}
	return nil, nil
}

func (f *file) Watch() (source.Watcher, error) {
	return newWatcher(f)
}
