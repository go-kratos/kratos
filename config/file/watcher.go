package file

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	f  *file
	fw *fsnotify.Watcher
}

func newWatcher(f *file) (config.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw.Add(f.path)
	return &watcher{f: f, fw: fw}, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case event := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			if _, err := os.Stat(event.Name); err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}
		fi, err := os.Stat(w.f.path)
		if err != nil {
			return nil, err
		}
		path := w.f.path
		if fi.IsDir() {
			path = filepath.Join(w.f.path, event.Name)
		}
		kv, err := w.f.loadFile(path)
		if err != nil {
			return nil, err
		}
		return []*config.KeyValue{kv}, nil
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
