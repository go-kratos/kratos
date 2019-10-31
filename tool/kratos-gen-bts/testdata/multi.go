package testdata

import (
	"context"
)

// mock test
var (
	_multiCacheFunc    func(c context.Context, keys []int64) (map[int64]*Demo, error)
	_multiRawFunc      func(c context.Context, keys []int64) (map[int64]*Demo, error)
	_multiAddCacheFunc func(c context.Context, values map[int64]*Demo) error
)

// CacheDemos .
func (d *dao) CacheDemos(c context.Context, keys []int64) (map[int64]*Demo, error) {
	// get data from cache
	return _multiCacheFunc(c, keys)
}

// RawDemos .
func (d *dao) RawDemos(c context.Context, keys []int64) (map[int64]*Demo, error) {
	// get data from db
	return _multiRawFunc(c, keys)
}

// AddCacheDemos .
func (d *dao) AddCacheDemos(c context.Context, values map[int64]*Demo) error {
	// add to cache
	return _multiAddCacheFunc(c, values)
}

// CacheDemos1 .
func (d *dao) CacheDemos1(c context.Context, keys []int64) (map[int64]*Demo, error) {
    // get data from cache
    return _multiCacheFunc(c, keys)
}

// RawDemos .
func (d *dao) RawDemos1(c context.Context, keys []int64) (map[int64]*Demo, error) {
    // get data from db
    return _multiRawFunc(c, keys)
}

// AddCacheDemos .
func (d *dao) AddCacheDemos1(c context.Context, values map[int64]*Demo) error {
    // add to cache
    return _multiAddCacheFunc(c, values)
}
