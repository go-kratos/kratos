package testdata

import (
	"context"
	"testing"
)

func TestArticle(t *testing.T) {
	d := New()
	c := context.TODO()
	art := &Article{ID: 1, Title: "title"}
	err := d.AddCacheArticle(c, art.ID, art)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	art1, err := d.CacheArticle(c, art.ID)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	if (art1.ID != art.ID) || (art.Title != art1.Title) {
		t.Error("art not equal")
		t.FailNow()
	}
	err = d.DelCacheArticle(c, art.ID)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	art1, err = d.CacheArticle(c, art.ID)
	if (art1 != nil) || (err != nil) {
		t.Errorf("art %v, err: %v", art1, err)
		t.FailNow()
	}
}

func TestNone(t *testing.T) {
	d := New()
	c := context.TODO()
	art := &Article{ID: 1, Title: "title"}
	err := d.AddCacheNone(c, art)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	art1, err := d.CacheNone(c)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	if (art1.ID != art.ID) || (art.Title != art1.Title) {
		t.Error("art not equal")
		t.FailNow()
	}
	err = d.DelCacheNone(c)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	art1, err = d.CacheNone(c)
	if (art1 != nil) || (err != nil) {
		t.Errorf("art %v, err: %v", art1, err)
		t.FailNow()
	}
}

func TestArticles(t *testing.T) {
	d := New()
	c := context.TODO()
	art1 := &Article{ID: 1, Title: "title"}
	art2 := &Article{ID: 2, Title: "title"}
	err := d.AddCacheArticles(c, map[int64]*Article{1: art1, 2: art2})
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	arts, err := d.CacheArticles(c, []int64{art1.ID, art2.ID})
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	if (arts[1].Title != art1.Title) || (arts[2].Title != art2.Title) {
		t.Error("art not equal")
		t.FailNow()
	}
	err = d.DelCacheArticles(c, []int64{art1.ID, art2.ID})
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	arts, err = d.CacheArticles(c, []int64{art1.ID, art2.ID})
	if (arts != nil) || (err != nil) {
		t.Errorf("art %v, err: %v", art1, err)
		t.FailNow()
	}
}

func TestString(t *testing.T) {
	d := New()
	c := context.TODO()
	err := d.AddCacheString(c, 1, "abc")
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	res, err := d.CacheString(c, 1)
	if err != nil {
		t.Errorf("err should be nil, get: %v", err)
		t.FailNow()
	}
	if res != "abc" {
		t.Error("res wrong")
		t.FailNow()
	}
}
