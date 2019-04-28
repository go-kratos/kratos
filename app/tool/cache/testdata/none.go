package testdata

import (
	"context"
)

// mock test
var (
	_noneCacheFunc    func(c context.Context) (*Article, error)
	_noneRawFunc      func(c context.Context) (*Article, error)
	_noneAddCacheFunc func(c context.Context, value *Article) error
)

// CacheNone .
func (d *Dao) CacheNone(c context.Context) (*Article, error) {
	// get data from cache
	return _noneCacheFunc(c)
}

// RawNone .
func (d *Dao) RawNone(c context.Context) (*Article, error) {
	// get data from db
	return _noneRawFunc(c)
}

// AddCacheNone .
func (d *Dao) AddCacheNone(c context.Context, value *Article) error {
	// add to cache
	return _noneAddCacheFunc(c, value)
}
