package testdata

import (
	"context"
)

// mock test
var (
	_singleCacheFunc    func(c context.Context, key int64) (*Article, error)
	_singleRawFunc      func(c context.Context, key int64) (*Article, error)
	_singleAddCacheFunc func(c context.Context, key int64, value *Article) error
)

// CacheArticle .
func (d *Dao) CacheArticle(c context.Context, key int64) (*Article, error) {
	// get data from cache
	return _singleCacheFunc(c, key)
}

// RawArticle .
func (d *Dao) RawArticle(c context.Context, key int64) (*Article, error) {
	// get data from db
	return _singleRawFunc(c, key)
}

// AddCacheArticle .
func (d *Dao) AddCacheArticle(c context.Context, key int64, value *Article) error {
	// add to cache
	return _singleAddCacheFunc(c, key, value)
}

// CacheArticle1 .
func (d *Dao) CacheArticle1(c context.Context, key int64, pn, ps int) (*Article, error) {
	// get data from cache
	return nil, nil
}

// RawArticle1 .
func (d *Dao) RawArticle1(c context.Context, key int64, pn, ps int) (*Article, *Article, error) {
	// get data from db
	return nil, nil, nil
}

// AddCacheArticle1 .
func (d *Dao) AddCacheArticle1(c context.Context, key int64, value *Article, pn, ps int) error {
	// add to cache
	return nil
}
