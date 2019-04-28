package regexp

import (
	// "fmt"
	"sync"
	// "syscall"
	// "go-common/log"
)

var (
	cache *Cache
)

func init() {
	cache = new(Cache)
	cache.Init()
}

func regexpKey(expr string) string {
	// tid := syscall.Gettid()
	// log.Info("tid", tid)
	// return fmt.Sprintf("%d_%s", tid, expr)
	return expr
}

type Cache struct {
	cache map[string]*Regexp
	sync.RWMutex
}

func (c *Cache) Init() {
	c.cache = make(map[string]*Regexp)
}

func (c *Cache) Get(expr string) *Regexp {
	c.RLock()
	defer c.RUnlock()

	return c.cache[regexpKey(expr)]
}

func (c *Cache) Set(expr string, regexp *Regexp) {
	c.Lock()
	defer c.Unlock()

	c.cache[regexpKey(expr)] = regexp
}
