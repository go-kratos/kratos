package dao

import (
	"context"

	"go-common/app/interface/main/shorturl/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

func cacheKey(short string) string {
	return _prefix + short
}

// Cache get short url cache.
func (d *Dao) Cache(c context.Context, short string) (su *model.ShortUrl, err error) {
	var (
		key  = cacheKey(short)
		conn = d.memchDB.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &su); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	return
}

// SetEmptyCache set empty cache for a few time
func (d *Dao) SetEmptyCache(c context.Context, short string) (err error) {
	var (
		key  = cacheKey(short)
		conn = d.memchDB.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: &model.ShortUrl{}, Flags: memcache.FlagJSON, Expiration: 300}); err != nil {
		log.Error("conn.Set error(%v)", err)
	}
	return
}

// SetCache save model.ShortUrl to memcache
func (d *Dao) SetCache(c context.Context, su *model.ShortUrl) (err error) {
	var (
		key  = cacheKey(su.Short)
		conn = d.memchDB.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: su, Flags: memcache.FlagJSON, Expiration: 0}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}
