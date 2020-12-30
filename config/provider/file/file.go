package file

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/config/provider"
)

type file struct {
	path string
}

// New new a file provider.
func New(path string) provider.Provider {
	return &file{path: path}
}

func (f *file) Load() (kvs []*provider.KeyValue, err error) {
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		data, err := ioutil.ReadFile(path.Join(f.path, file.Name()))
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, &provider.KeyValue{
			Key:       file.Name(),
			Value:     data,
			Format:    format(file.Name()),
			Timestamp: file.ModTime(),
		})
	}
	return nil, nil
}

func (f *file) Watch() (provider.Watcher, error) {
	return newWatcher(f)
}
