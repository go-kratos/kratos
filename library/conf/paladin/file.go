package paladin

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var _ Client = &file{}

// file is file config client.
type file struct {
	ch     chan Event
	values *Map
}

// NewFile new a config file client.
// conf = /data/conf/app/
// conf = /data/conf/app/xxx.toml
func NewFile(base string) (Client, error) {
	// paltform slash
	base = filepath.FromSlash(base)
	fi, err := os.Stat(base)
	if err != nil {
		panic(err)
	}
	// dirs or file to paths
	var paths []string
	if fi.IsDir() {
		files, err := ioutil.ReadDir(base)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if !file.IsDir() {
				paths = append(paths, path.Join(base, file.Name()))
			}
		}
	} else {
		paths = append(paths, base)
	}
	// laod config file to values
	values := make(map[string]*Value, len(paths))
	for _, file := range paths {
		if file == "" {
			return nil, errors.New("paladin: path is empty")
		}
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		s := string(b)
		values[path.Base(file)] = &Value{val: s, raw: s}
	}
	m := new(Map)
	m.Store(values)
	return &file{values: m, ch: make(chan Event, 10)}, nil
}

// Get return value by key.
func (f *file) Get(key string) *Value {
	return f.values.Get(key)
}

// GetAll return value map.
func (f *file) GetAll() *Map {
	return f.values
}

// WatchEvent watch multi key.
func (f *file) WatchEvent(ctx context.Context, key ...string) <-chan Event {
	return f.ch
}

// Close close watcher.
func (f *file) Close() error {
	close(f.ch)
	return nil
}
