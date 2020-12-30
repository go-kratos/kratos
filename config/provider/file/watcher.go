package file

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kratos/kratos/v2/config/provider"
)

type watcher struct {
	f  *file
	fw *fsnotify.Watcher
}

func newWatcher(f *file) (provider.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fw.Add(f.path)
	return &watcher{f: f, fw: fw}, nil
}

func (w *watcher) Next() ([]*provider.KeyValue, error) {
	select {
	case event := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				w.fw.Add(event.Name)
			}
		}
		return w.f.Load()
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w *watcher) Close() error {
	return w.fw.Close()
}
