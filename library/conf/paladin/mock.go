package paladin

import (
	"context"
)

var _ Client = &mock{}

// mock is mock config client.
type mock struct {
	ch     chan Event
	values *Map
}

// NewMock new a config mock client.
func NewMock(vs map[string]string) Client {
	values := make(map[string]*Value, len(vs))
	for k, v := range vs {
		values[k] = &Value{val: v, raw: v}
	}
	m := new(Map)
	m.Store(values)
	return &mock{values: m, ch: make(chan Event)}
}

// Get return value by key.
func (m *mock) Get(key string) *Value {
	return m.values.Get(key)
}

// GetAll return value map.
func (m *mock) GetAll() *Map {
	return m.values
}

// WatchEvent watch multi key.
func (m *mock) WatchEvent(ctx context.Context, key ...string) <-chan Event {
	return m.ch
}

// Close close watcher.
func (m *mock) Close() error {
	close(m.ch)
	return nil
}
