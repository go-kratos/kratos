package memory

import "github.com/go-kratos/kratos/v2/config/provider"

var _ provider.Provider = (*memory)(nil)

type memory struct {
	ch  chan provider.KeyValue
	kvs []provider.KeyValue
}

// New new a memory provider.
func New(ch chan provider.KeyValue, kvs ...provider.KeyValue) provider.Provider {
	return &memory{kvs: kvs}
}

func (m *memory) Load() ([]provider.KeyValue, error) {
	return m.kvs, nil
}

func (m *memory) Watch() <-chan provider.KeyValue {
	return m.ch
}
