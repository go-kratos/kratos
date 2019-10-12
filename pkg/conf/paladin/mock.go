package paladin

import (
	"context"
)

var _ Client = &Mock{}

// Mock is Mock config client.
type Mock struct {
	C chan Event
	*Map
}

// NewMock new a config mock client.
func NewMock(vs map[string]string) Client {
	values := make(map[string]*Value, len(vs))
	for k, v := range vs {
		values[k] = &Value{val: v, raw: v}
	}
	m := new(Map)
	m.Store(values)
	return &Mock{Map: m, C: make(chan Event)}
}

// GetAll return value map.
func (m *Mock) GetAll() *Map {
	return m.Map
}

// WatchEvent watch multi key.
func (m *Mock) WatchEvent(ctx context.Context, key ...string) <-chan Event {
	return m.C
}

// Close close watcher.
func (m *Mock) Close() error {
	close(m.C)
	return nil
}
