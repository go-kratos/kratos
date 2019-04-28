package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/passport-sns/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

func snsKey(platform string, mid int64) string {
	return fmt.Sprintf("sns_%s_%d", platform, mid)
}

func oauth2Key(platform string, openID string) string {
	return fmt.Sprintf("sns_oauth2_%s_%s", platform, openID)
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

// SetOauth2Cache set oauth2 info to cache
func (d *Dao) SetOauth2Cache(c context.Context, openID, platform string, sns *model.Oauth2Proto) (err error) {
	key := oauth2Key(platform, openID)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: sns, Flags: memcache.FlagProtobuf, Expiration: 300}
	if err = conn.Set(item); err != nil {
		log.Error("fail to set oauth2 info to mc, key(%s) expire(%d) error(%+v)", key, 300, err)
	}
	return
}

// SnsCache sns cache
func (d *Dao) SnsCache(c context.Context, mid int64, platform string) (res *model.SnsProto, err error) {
	key := snsKey(platform, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return nil, err
	}
	res = new(model.SnsProto)
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
		return nil, err
	}
	return res, nil
}

// Oauth2Cache oauth2 info cache
func (d *Dao) Oauth2Cache(c context.Context, openID, platform string) (res *model.Oauth2Proto, err error) {
	key := oauth2Key(platform, openID)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return nil, err
	}
	res = new(model.Oauth2Proto)
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
		return nil, err
	}
	return res, nil
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

// DelOauth2Cache del oauth2 cache
func (d *Dao) DelOauth2Cache(c context.Context, openID, platform string) (err error) {
	key := oauth2Key(platform, openID)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("fail to del oauth2 cache, key(%s) error(%+v)", key, err)
	}
	return
}
