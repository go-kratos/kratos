// Package group provides a sample lazy load container.
// The group only creating a new object not until the object is needed by user.
// And it will cache all the objects to reduce the creation of object.
package group

import "sync"

// Group is a lazy load container.
type Group struct {
	new  func() interface{}
	objs sync.Map
	sync.RWMutex
}

// NewGroup news a group container.
func NewGroup(new func() interface{}) *Group {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group{
		new: new,
	}
}

// Get gets the object by the given key.
func (g *Group) Get(key string) interface{} {
	g.RLock()
	new := g.new
	g.RUnlock()
	obj, ok := g.objs.Load(key)
	if !ok {
		obj = new()
		g.objs.Store(key, obj)
	}
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
	g.objs.Range(func(key, value interface{}) bool {
		g.objs.Delete(key)
		return true
	})
}
