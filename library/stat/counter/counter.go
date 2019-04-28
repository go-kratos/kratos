package counter

import (
	"sync"
)

// Counter is a counter interface.
type Counter interface {
	Add(int64)
	Reset()
	Value() int64
}

// Group is a counter group.
type Group struct {
	mu   sync.RWMutex
	vecs map[string]Counter

	// New optionally specifies a function to generate a counter.
	// It may not be changed concurrently with calls to other functions.
	New func() Counter
}

// Add add a counter by a specified key, if counter not exists then make a new one and return new value.
func (g *Group) Add(key string, value int64) {
	g.mu.RLock()
	vec, ok := g.vecs[key]
	g.mu.RUnlock()
	if !ok {
		vec = g.New()
		g.mu.Lock()
		if g.vecs == nil {
			g.vecs = make(map[string]Counter)
		}
		if _, ok = g.vecs[key]; !ok {
			g.vecs[key] = vec
		}
		g.mu.Unlock()
	}
	vec.Add(value)
}

// Value get a counter value by key.
func (g *Group) Value(key string) int64 {
	g.mu.RLock()
	vec, ok := g.vecs[key]
	g.mu.RUnlock()
	if ok {
		return vec.Value()
	}
	return 0
}

// Reset reset a counter by key.
func (g *Group) Reset(key string) {
	g.mu.RLock()
	vec, ok := g.vecs[key]
	g.mu.RUnlock()
	if ok {
		vec.Reset()
	}
}
