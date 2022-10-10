package cache

import (
	"context"
	"time"
)

// Cache is a generic key value cache interface.
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Put(ctx context.Context, key string, value interface{}, expired time.Duration) error
	Delete(ctx context.Context, key string) error
}
