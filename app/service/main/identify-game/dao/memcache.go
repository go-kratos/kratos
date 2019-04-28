package dao

import (
	"context"

	"go-common/app/service/main/identify-game/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

// SetAccessCache .
func (d *Dao) SetAccessCache(c context.Context, key string, res *model.AccessInfo) (err error) {
	if res.Expires < 0 {
		log.Error("identify-game expire error(expires:%d)", res.Expires)
		return
	}
	item := &memcache.Item{Key: key, Object: res, Flags: memcache.FlagProtobuf, Expiration: int32(res.Expires)}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		log.Error("identify-game set error(%s,%d,%v)", key, res.Expires, err)
	}
	return
}

// AccessCache .
func (d *Dao) AccessCache(c context.Context, key string) (res *model.AccessInfo, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
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
	res = &model.AccessInfo{}
	if err = conn.Scan(r, res); err != nil {
		PromError("mc:json解析失败")
		log.Error("conn.Scan(%v) error(%v)", string(r.Value), err)
		return
	}
	cachedCount.Incr("access_cache")
	return
}

// DelAccessCache del cache.
func (d *Dao) DelAccessCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", key, err)
	}
	return
}
