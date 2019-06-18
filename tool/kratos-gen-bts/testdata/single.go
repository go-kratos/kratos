package testdata

import (
	"context"
)

// mock test
var (
	_singleCacheFunc    func(c context.Context, key int64) (*Demo, error)
	_singleRawFunc      func(c context.Context, key int64) (*Demo, error)
	_singleAddCacheFunc func(c context.Context, key int64, value *Demo) error
)

// CacheDemo .
func (d *Dao) CacheDemo(c context.Context, key int64) (*Demo, error) {
	// get data from cache
	return _singleCacheFunc(c, key)
}

// RawDemo .
func (d *Dao) RawDemo(c context.Context, key int64) (*Demo, error) {
	// get data from db
	return _singleRawFunc(c, key)
}

// AddCacheDemo .
func (d *Dao) AddCacheDemo(c context.Context, key int64, value *Demo) error {
	// add to cache
	return _singleAddCacheFunc(c, key, value)
}

// CacheDemo1 .
func (d *Dao) CacheDemo1(c context.Context, key int64, pn, ps int) (*Demo, error) {
	// get data from cache
	return nil, nil
}

// RawDemo1 .
func (d *Dao) RawDemo1(c context.Context, key int64, pn, ps int) (*Demo, *Demo, error) {
	// get data from db
	return nil, nil, nil
}

// AddCacheDemo1 .
func (d *Dao) AddCacheDemo1(c context.Context, key int64, value *Demo, pn, ps int) error {
	// add to cache
	return nil
}
