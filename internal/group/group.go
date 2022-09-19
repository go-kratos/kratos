// Package group provides a sample lazy load container.
// The group only creating a new object not until the object is needed by user.
// And it will cache all the objects to reduce the creation of object.
package group

import "sync"

// Group is a lazy load container.
type Group struct {
	new  func() interface{}
	vals map[string]interface{}
	sync.RWMutex
}

// NewGroup news a group container.
func NewGroup(new func() interface{}) *Group {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group{
		new:  new,
		vals: make(map[string]interface{}),
	}
}

// Get gets the object by the given key.
func (g *Group) Get(key string) interface{} {
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
	v = g.new()
	g.vals[key] = v
	return v
}

// Reset resets the new function and deletes all existing objects.
func (g *Group) Reset(new func() interface{}) {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	g.Lock()
	g.new = new
	g.Unlock()
	g.Clear()
}

// Clear deletes all objects.
func (g *Group) Clear() {
	g.Lock()
	g.vals = make(map[string]interface{})
	g.Unlock()
}
