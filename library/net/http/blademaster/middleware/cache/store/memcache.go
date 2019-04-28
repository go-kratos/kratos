package store

import (
	"context"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

// Memcache represents the cache with memcached persistence
type Memcache struct {
	pool *memcache.Pool
}

// NewMemcache new a memcache store.
func NewMemcache(c *memcache.Config) *Memcache {
	if c == nil {
		panic("cache config is nil")
	}
	return &Memcache{
		pool: memcache.NewPool(c),
	}
}

// Set save the result to memcache store.
func (ms *Memcache) Set(ctx context.Context, key string, value []byte, expire int32) (err error) {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expire,
	}
	conn := ms.pool.Get(ctx)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", key, err)
	}
	return
}

// Get get result from mc by locaiton+params.
func (ms *Memcache) Get(ctx context.Context, key string) ([]byte, error) {
	conn := ms.pool.Get(ctx)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			//ignore not found error
			return nil, nil
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return nil, err
	}
	return r.Value, nil
}
