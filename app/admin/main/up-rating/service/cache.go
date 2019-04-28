package service

import (
	"time"
)

// Cache simple cache
type Cache struct {
	d    int64                  // duration seconds
	mc   map[string]interface{} // map cache
	snap time.Time
}

// NewCache new cache
func NewCache(d int64) *Cache {
	c := &Cache{
		d:    d,
		mc:   make(map[string]interface{}),
		snap: time.Now(),
	}
	return c
}

// Get ...
func (c *Cache) Get(key string) (val interface{}) {
	c.check()
	return c.mc[key]
}

// Set ...
func (c *Cache) Set(key string, val interface{}) {
	c.check()
	c.mc[key] = val
}

func (c *Cache) check() {
	if int64(time.Since(c.snap).Seconds()) > c.d {
		c.mc = make(map[string]interface{})
		c.snap = time.Now()
	}
}
