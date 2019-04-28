package cache

import (
	"sync"
	"time"

	"go-common/app/service/bbq/recsys-recall/conf"
)

// Solt .
type Solt struct {
	data       []byte
	ctime      int64
	maxAge     int64
	lastUsed   int64
	keepAlived int64
}

// LocalCache .
type LocalCache struct {
	d      map[string]*Solt
	l1tags map[string]byte
	l2tags map[string]byte
	c      *conf.LocalCacheConfig
	lock   *sync.RWMutex
}

// NewLocalCache .
func NewLocalCache(c *conf.LocalCacheConfig) *LocalCache {
	l1 := make(map[string]byte)
	l2 := make(map[string]byte)
	for _, v := range c.L1Tags {
		l1[v] = byte(1)
	}
	for _, v := range c.L2Tags {
		l2[v] = byte(1)
	}

	return &LocalCache{
		d:      make(map[string]*Solt),
		l1tags: l1,
		l2tags: l2,
		c:      c,
		lock:   &sync.RWMutex{},
	}
}

// Set .
func (lc *LocalCache) Set(key string, val []byte) bool {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	keep := lc.c.Level3
	if _, ok := lc.l1tags[key]; ok {
		keep = lc.c.Level1
	} else if _, ok := lc.l2tags[key]; ok {
		keep = lc.c.Level2
	}

	lc.d[key] = &Solt{
		data:       val,
		ctime:      time.Now().Unix(),
		maxAge:     int64(lc.c.MaxAge),
		lastUsed:   time.Now().Unix(),
		keepAlived: int64(keep),
	}
	return true
}

// Get .
func (lc *LocalCache) Get(key string) []byte {
	lc.lock.RLock()
	defer lc.lock.RUnlock()

	current := time.Now().Unix()
	s := lc.d[key]
	if s == nil {
		return nil
	}

	keepAlived := s.keepAlived / int64(time.Second)
	maxAge := s.maxAge / int64(time.Second)
	if keepAlived < (current-s.lastUsed) || maxAge < (current-s.ctime) {
		return nil
	}

	s.lastUsed = current
	return s.data
}
