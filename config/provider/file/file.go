package file

import (
	"io/ioutil"
	"path"

	"github.com/go-kratos/kratos/v2/config/provider"
)

// Option is config file option.
type Option func(*file)

type file struct {
	path string
}

// New new a file provider.
func New() provider.Provider {
	return &file{}
}

func (f *file) Load() (kvs []provider.KeyValue, err error) {
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(path.Join(f.path, file.Name()))
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, provider.KeyValue{
			Key:       file.Name(),
			Value:     data,
			Format:    format(file.Name()),
			Timestamp: file.ModTime(),
		})
	}
	return nil, nil
}

func (f *file) Watch() (provider.Watcher, error) {
	return nil, nil
}
