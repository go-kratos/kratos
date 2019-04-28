package cache

import (
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/cache/store"
)

// Cache is the abstract struct for any cache impl
type Cache struct {
	store store.Store
}

// Filter is used to check is cache required for every request
type Filter func(*bm.Context) bool

// Policy is used to abstract different cache policy
type Policy interface {
	Key(*bm.Context) string
	Handler(store.Store) bm.HandlerFunc
}

// New will create a new Cache struct
func New(store store.Store) *Cache {
	c := &Cache{
		store: store,
	}
	return c
}

// Cache is used to mark path as customized cache policy
func (c *Cache) Cache(policy Policy, filter Filter) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		if filter != nil && !filter(ctx) {
			return
		}
		policy.Handler(c.store)(ctx)
	}
}
