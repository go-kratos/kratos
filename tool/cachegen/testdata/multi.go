package testdata

import (
	"context"
)

// mock test
var (
	_multiCacheFunc    func(c context.Context, keys []int64) (map[int64]*Article, error)
	_multiRawFunc      func(c context.Context, keys []int64) (map[int64]*Article, error)
	_multiAddCacheFunc func(c context.Context, values map[int64]*Article) error
)

// CacheArticles .
func (d *Dao) CacheArticles(c context.Context, keys []int64) (map[int64]*Article, error) {
	// get data from cache
	return _multiCacheFunc(c, keys)
}

// RawArticles .
func (d *Dao) RawArticles(c context.Context, keys []int64) (map[int64]*Article, error) {
	// get data from db
	return _multiRawFunc(c, keys)
}

// AddCacheArticles .
func (d *Dao) AddCacheArticles(c context.Context, values map[int64]*Article) error {
	// add to cache
	return _multiAddCacheFunc(c, values)
}
