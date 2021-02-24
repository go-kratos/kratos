package memory

import (
	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	Id      string
	Updates chan *config.KeyValue
	Source  *memory
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	kv := <-w.Updates
	return []*config.KeyValue{kv}, nil
}

func (w *watcher) Close() error {
	w.Source.Lock()
	delete(w.Source.Watchers, w.Id)
	w.Source.Unlock()
	return nil
}
