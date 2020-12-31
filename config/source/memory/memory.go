package memory

import (
	"github.com/go-kratos/kratos/v2/config/source"
)

var _ source.Source = (*memory)(nil)

type memory struct {
	kvs []*source.KeyValue
}

// New new a memory provider.
func New(kvs ...*source.KeyValue) source.Source {
	return &memory{kvs: kvs}
}

func (m *memory) Load() ([]*source.KeyValue, error) {
	return m.kvs, nil
}

func (m *memory) Watch() (source.Watcher, error) {
	return newWatcher(), nil
}
