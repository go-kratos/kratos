package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/passport-sns/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

func snsKey(platform string, mid int64) string {
	return fmt.Sprintf("sns_%s_%d", platform, mid)
}

// SetSnsCache set sns to cache
func (d *Dao) SetSnsCache(c context.Context, mid int64, platform string, sns *model.SnsProto) (err error) {
	key := snsKey(platform, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: sns, Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set sns to mc, key(%s) expire(%d) error(%+v)", key, d.mcExpire, err)
	}
	return
}

// DelSnsCache del sns cache
func (d *Dao) DelSnsCache(c context.Context, mid int64, platform string) (err error) {
	key := snsKey(platform, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del sns cache, key(%s) error(%+v)", key, err)
	}
	return
}

// GetUnionIDCache .
func (d *Dao) GetUnionIDCache(c context.Context, key string) (v string, err error) {
	conn := d.mc.Get(c)
	r, err := conn.Get(key)
	conn.Close()
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("GetUnionIDCache, key(%s) error(%+v)", key, err)
		return
	}
	v = string(r.Value)
	return
}
