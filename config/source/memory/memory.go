package memory

import (
	"github.com/go-kratos/kratos/v2/config/source"
)

var _ source.Source = (*memory)(nil)

type memory struct {
	ch  chan *source.KeyValue
	kvs []*source.KeyValue
}

// New new a memory provider.
func New(ch chan *source.KeyValue, kvs ...*source.KeyValue) source.Source {
	return &memory{ch: ch, kvs: kvs}
}

func (m *memory) Load() ([]*source.KeyValue, error) {
	return m.kvs, nil
}

func (m *memory) Watch() (source.Watcher, error) {
	return newWatcher(m.ch), nil
}
