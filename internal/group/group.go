// Package group provides a sample lazy load container.
// The group only creating a new object not until the object is needed by user.
// And it will cache all the objects to reduce the creation of object.
package group

import "sync"

// Factory is a function that creates an object of type T.
type Factory[T any] func() T

// Group is a lazy load container.
type Group[T any] struct {
	factory func() T
	vals    map[string]T
	sync.RWMutex
}

// NewGroup news a group container.
func NewGroup[T any](factory Factory[T]) *Group[T] {
	if factory == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group[T]{
		factory: factory,
		vals:    make(map[string]T),
	}
}

// Get gets the object by the given key.
func (g *Group[T]) Get(key string) T {
	g.RLock()
	v, ok := g.vals[key]
	if ok {
		g.RUnlock()
		return v
	}
	g.RUnlock()

	// slow path for group don`t have specified key value
	g.Lock()
	defer g.Unlock()
	v, ok = g.vals[key]
	if ok {
		return v
	}
	v = g.factory()
	g.vals[key] = v
	return v
}

// Reset resets the new function and deletes all existing objects.
func (g *Group[T]) Reset(factory Factory[T]) {
	if factory == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	g.Lock()
	g.factory = factory
	g.Unlock()
	g.Clear()
}

// Clear deletes all objects.
func (g *Group[T]) Clear() {
	g.Lock()
	g.vals = make(map[string]T)
	g.Unlock()
}
