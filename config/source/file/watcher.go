package file

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kratos/kratos/v2/config/source"
)

type watcher struct {
	f  *file
	fw *fsnotify.Watcher
}

func newWatcher(f *file) (source.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw.Add(f.path)
	return &watcher{f: f, fw: fw}, nil
}

func (w *watcher) Next() (*source.KeyValue, error) {
	select {
	case event := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}
		return w.f.loadFile(filepath.Join(w.f.path, event.Name))
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w *watcher) Close() error {
	return w.fw.Close()
}
