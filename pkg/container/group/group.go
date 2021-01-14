// Package group provides a sample lazy load container.
// The group only creating a new object not until the object is needed by user.
// And it will cache all the objects to reduce the creation of object.
package group

import "sync"

// Group is a lazy load container.
type Group struct {
	new  func() interface{}
	objs map[string]interface{}
	sync.RWMutex
}

// NewGroup news a group container.
func NewGroup(new func() interface{}) *Group {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group{
		new:  new,
		objs: make(map[string]interface{}),
	}
}

// Get gets the object by the given key.
func (g *Group) Get(key string) interface{} {
	g.RLock()
	obj, ok := g.objs[key]
	if ok {
		g.RUnlock()
		return obj
	}
	g.RUnlock()

	// double check
	g.Lock()
	defer g.Unlock()
	obj, ok = g.objs[key]
	if ok {
		return obj
	}
	obj = g.new()
	g.objs[key] = obj
	return obj
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
	g.objs = make(map[string]interface{})
	g.Unlock()
}
