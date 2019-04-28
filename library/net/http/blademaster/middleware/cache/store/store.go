package store

import (
	"context"
)

// Store is the interface of a cache backend
type Store interface {
	// Get retrieves an item from the cache. Returns the item or nil, and a bool indicating
	// whether the key was found.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set sets an item to the cache, replacing any existing item.
	Set(ctx context.Context, key string, value []byte, expire int32) error
}
