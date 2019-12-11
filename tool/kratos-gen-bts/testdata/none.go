package testdata

import (
	"context"
)

// mock test
var (
	_noneCacheFunc    func(c context.Context) (*Demo, error)
	_noneRawFunc      func(c context.Context) (*Demo, error)
	_noneAddCacheFunc func(c context.Context, value *Demo) error
)

// CacheNone .
func (d *dao) CacheNone(c context.Context) (*Demo, error) {
	// get data from cache
	return _noneCacheFunc(c)
}

// RawNone .
func (d *dao) RawNone(c context.Context) (*Demo, error) {
	// get data from db
	return _noneRawFunc(c)
}

// AddCacheNone .
func (d *dao) AddCacheNone(c context.Context, value *Demo) error {
	// add to cache
	return _noneAddCacheFunc(c, value)
}
