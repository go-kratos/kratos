package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/identify/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

var (
	loginCacheValue = []byte("1")
)

// SetAccessCache .
func (d *Dao) SetAccessCache(c context.Context, key string, res *model.IdentifyInfo) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key = cacheKey(key)
	item := &memcache.Item{Key: key, Object: res, Flags: memcache.FlagProtobuf, Expiration: res.Expires}
	if err := conn.Set(item); err != nil {
		log.Error("identify set error(%s,%d,%v)", key, res.Expires, err)
	}
}

// AccessCache .
func (d *Dao) AccessCache(c context.Context, key string) (res *model.IdentifyInfo, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key = cacheKey(key)
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			missedCount.Incr("access_cache")
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	res = &model.IdentifyInfo{}
	if err = conn.Scan(r, res); err != nil {
		PromError("mc:json解析失败")
		log.Error("conn.Scan(%v) error(%v)", string(r.Value), err)
		return
	}
	cachedCount.Incr("access_cache")
	return
}

// DelCache delete access cache.
func (d *Dao) DelCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key = cacheKey(key)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("dao.DelCache(%s) error(%v)", key, err)
	}
	return
}

func cacheKey(key string) string {
	return fmt.Sprintf("i_%s", key)
}

func loginCacheKey(mid int64, ip string) string {
	return fmt.Sprintf("l%d%s", mid, ip)
}

// SetLoginCache set login cache
func (d *Dao) SetLoginCache(c context.Context, mid int64, ip string, expires int32) (err error) {
	key := loginCacheKey(mid, ip)
	conn := d.mcLogin.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Value: loginCacheValue, Flags: memcache.FlagRAW, Expiration: expires}
	// use Add instead of Set
	if err = conn.Set(item); err != nil {
		log.Error("loginCache set error(%s,%v)", key, err)
	}
	return
}

// ExistMIDAndIP check is exist mid
func (d *Dao) ExistMIDAndIP(c context.Context, mid int64, ip string) (ok bool, err error) {
	key := loginCacheKey(mid, ip)
	conn := d.mcLogin.Get(c)
	defer conn.Close()
	_, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			missedCount.Incr("isExistMID")
			err = nil
			return false, nil
		}
		log.Error("loginCache conn.Get(%s) error(%v)", key, err)
		return
	}
	return true, nil
}
