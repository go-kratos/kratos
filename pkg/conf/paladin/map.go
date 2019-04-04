package paladin

import (
	"strings"
	"sync/atomic"
)

// keyNamed key naming to lower case.
func keyNamed(key string) string {
	return strings.ToLower(key)
}

// Map is config map, key(filename) -> value(file).
type Map struct {
	values atomic.Value
}

// Store sets the value of the Value to values map.
func (m *Map) Store(values map[string]*Value) {
	dst := make(map[string]*Value, len(values))
	for k, v := range values {
		dst[keyNamed(k)] = v
	}
	m.values.Store(dst)
}

// Load returns the value set by the most recent Store.
func (m *Map) Load() map[string]*Value {
	return m.values.Load().(map[string]*Value)
}

// Exist check if values map exist a key.
func (m *Map) Exist(key string) bool {
	_, ok := m.Load()[keyNamed(key)]
	return ok
}

// Get return get value by key.
func (m *Map) Get(key string) *Value {
	v, ok := m.Load()[keyNamed(key)]
	if ok {
		return v
	}
	return &Value{}
}

// Keys return map keys.
func (m *Map) Keys() []string {
	values := m.Load()
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
