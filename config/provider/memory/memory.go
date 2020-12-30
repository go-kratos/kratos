package memory

import "github.com/go-kratos/kratos/v2/config/provider"

var _ provider.Provider = (*memory)(nil)

type memory struct {
	kvs []*provider.KeyValue
}

// New new a memory provider.
func New(kvs ...*provider.KeyValue) provider.Provider {
	return &memory{kvs: kvs}
}

func (m *memory) Load() ([]*provider.KeyValue, error) {
	return m.kvs, nil
}

func (m *memory) Watch() (provider.Watcher, error) {
	return nil, nil
}
