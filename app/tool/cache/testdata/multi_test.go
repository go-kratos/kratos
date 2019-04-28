package testdata

import (
	"context"
	"errors"
	"testing"
)

func TestMultiCache(t *testing.T) {
	id := int64(1)
	d := New()
	meta := map[int64]*Article{id: {ID: id}}
	getsFromCache := func(c context.Context, keys []int64) (map[int64]*Article, error) { return meta, nil }
	notGetsFromCache := func(c context.Context, keys []int64) (map[int64]*Article, error) { return nil, errors.New("err") }
	// 缓存返回了部分数据
	partFromCache := func(c context.Context, keys []int64) (map[int64]*Article, error) { return meta, errors.New("err") }
	getsFromSource := func(c context.Context, keys []int64) (map[int64]*Article, error) { return meta, nil }
	notGetsFromSource := func(c context.Context, keys []int64) (map[int64]*Article, error) {
		return meta, errors.New("err")
	}
	addToCache := func(c context.Context, values map[int64]*Article) error { return nil }
	// gets from cache
	_multiCacheFunc = getsFromCache
	_multiRawFunc = notGetsFromSource
	_multiAddCacheFunc = addToCache
	res, err := d.Articles(context.TODO(), []int64{id})
	if err != nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res[1].ID != 1 {
		t.Fatalf("id should be 1")
	}
	// get from source
	_multiCacheFunc = notGetsFromCache
	_multiRawFunc = getsFromSource
	res, err = d.Articles(context.TODO(), []int64{id})
	if err != nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res[1].ID != 1 {
		t.Fatalf("id should be 1")
	}
	// 缓存失败 返回部分数据 回源也失败的情况
	_multiCacheFunc = partFromCache
	_multiRawFunc = notGetsFromSource
	res, err = d.Articles(context.TODO(), []int64{id})
	if err == nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res[1].ID != 1 {
		t.Fatalf("id should be 1")
	}
	// with null cache
	nullCache := &Article{ID: -1}
	getNullFromCache := func(c context.Context, keys []int64) (map[int64]*Article, error) {
		return map[int64]*Article{id: nullCache}, nil
	}
	_multiCacheFunc = getNullFromCache
	_multiRawFunc = notGetsFromSource
	res, err = d.Articles(context.TODO(), []int64{id})
	if err != nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res[id] != nil {
		t.Fatalf("res should be nil")
	}
}
