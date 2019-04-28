package dao

import (
	"context"
	"strconv"

	"go-common/app/service/main/vip/model"
	mc "go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_prefixBindByMid        = "b:m:"
	_prefixOpenInfoByOpenID = "o:i:"
)

func keyBindByMid(mid int64, appID int64) string {
	return _prefixBindByMid + strconv.FormatInt(mid, 10) + ":" + strconv.FormatInt(appID, 10)
}

func keyOpenInfoByOpenID(openID string, appID int64) string {
	return _prefixOpenInfoByOpenID + openID + ":" + strconv.FormatInt(appID, 10)
}

// CacheBindInfoByMid get vip bind by mid cache.
func (d *Dao) CacheBindInfoByMid(c context.Context, mid int64, appID int64) (v *model.OpenBindInfo, err error) {
	key := keyBindByMid(mid, appID)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao cache bind by mid")
		return
	}
	v = new(model.OpenBindInfo)
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrapf(err, "dao cache scan bind by mid")
	}
	return
}

// CacheOpenInfoByOpenID get vip open info by open id cache.
func (d *Dao) CacheOpenInfoByOpenID(c context.Context, openID string, appID int64) (v *model.OpenInfo, err error) {
	key := keyOpenInfoByOpenID(openID, appID)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao cache open by openid")
		return
	}
	v = new(model.OpenInfo)
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrapf(err, "dao cache scan open by openid")
	}
	return
}

// AddCacheBindInfoByMid add bind info cache.
func (d *Dao) AddCacheBindInfoByMid(c context.Context, mid int64, v *model.OpenBindInfo, appID int64) (err error) {
	item := &mc.Item{
		Key:        keyBindByMid(mid, appID),
		Object:     v,
		Expiration: d.mcExpire,
		Flags:      mc.FlagProtobuf,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		err = errors.Wrapf(err, "dao add cache bind by mid")
	}
	return
}

// AddCacheOpenInfoByOpenID add open info cache.
func (d *Dao) AddCacheOpenInfoByOpenID(c context.Context, openID string, v *model.OpenInfo, appID int64) (err error) {
	item := &mc.Item{
		Key:        keyOpenInfoByOpenID(openID, appID),
		Object:     v,
		Expiration: d.mcExpire,
		Flags:      mc.FlagProtobuf,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		err = errors.Wrapf(err, "dao add cache open by openid")
	}
	return
}

// DelBindInfoCache del bind info cache.
func (d *Dao) DelBindInfoCache(c context.Context, mid int64, appID int64) (err error) {
	return d.delCache(c, keyBindByMid(mid, appID))
}

// DelOpenInfoCache del open info cache.
func (d *Dao) DelOpenInfoCache(c context.Context, openID string, appID int64) (err error) {
	return d.delCache(c, keyOpenInfoByOpenID(openID, appID))
}

func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == mc.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "del cache")
		}
	}
	return
}
