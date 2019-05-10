package testdata

import (
	"context"
	"errors"
	"testing"
)

func TestNoneCache(t *testing.T) {
	d := New()
	meta := &Article{ID: 1}
	getFromCache := func(c context.Context) (*Article, error) { return meta, nil }
	notGetFromCache := func(c context.Context) (*Article, error) { return nil, errors.New("err") }
	getFromSource := func(c context.Context) (*Article, error) { return meta, nil }
	notGetFromSource := func(c context.Context) (*Article, error) { return meta, errors.New("err") }
	addToCache := func(c context.Context, values *Article) error { return nil }
	// get from cache
	_noneCacheFunc = getFromCache
	_noneRawFunc = notGetFromSource
	_noneAddCacheFunc = addToCache
	res, err := d.None(context.TODO())
	if err != nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res.ID != 1 {
		t.Fatalf("id should be 1")
	}
	// get from source
	_noneCacheFunc = notGetFromCache
	_noneRawFunc = getFromSource
	res, err = d.None(context.TODO())
	if err != nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res.ID != 1 {
		t.Fatalf("id should be 1")
	}
	// with null cache
	nullCache := &Article{ID: -1}
	getNullFromCache := func(c context.Context) (*Article, error) { return nullCache, nil }
	_noneCacheFunc = getNullFromCache
	_noneRawFunc = notGetFromSource
	res, err = d.None(context.TODO())
	if err != nil {
		t.Fatalf("err should be nil, get: %v", err)
	}
	if res != nil {
		t.Fatalf("res should be nil")
	}
}
