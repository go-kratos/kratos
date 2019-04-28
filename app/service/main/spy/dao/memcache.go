package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/spy/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_userKey = "u_%d"
)

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mcUser.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		log.Error("conn.Store(set,ping,1) error(%v)", err)
	}
	return
}

func userInfoCacheKey(mid int64) string {
	return fmt.Sprintf(_userKey, mid)
}

// UserInfoCache get user info to cache.
func (d *Dao) UserInfoCache(c context.Context, mid int64) (ui *model.UserInfo, err error) {
	var (
		key  = userInfoCacheKey(mid)
		conn = d.mcUser.Get(c)
	)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(get, %s) error(%v)", key, err)
		return
	}
	ui = &model.UserInfo{}
	if err = conn.Scan(reply, ui); err != nil {
		log.Error("reply.Scan(%s) error(%v)", string(reply.Value), err)
	}
	return
}

// AddUserInfoCache add info to cache , return err if key is exist already.
func (d *Dao) AddUserInfoCache(c context.Context, ui *model.UserInfo) (err error) {
	if ui == nil {
		return fmt.Errorf("AddUserInfoCache got nil *model.UserInfo")
	}
	var (
		key  = userInfoCacheKey(ui.Mid)
		conn = d.mcUser.Get(c)
	)
	defer conn.Close()

	if err = conn.Add(&memcache.Item{Key: key, Object: ui, Expiration: d.mcUserExpire, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Add(%s, %v) error(%v)", key, ui, err)
		return
	}
	return
}

// SetUserInfoCache set info cache.
func (d *Dao) SetUserInfoCache(c context.Context, ui *model.UserInfo) (err error) {
	if ui == nil {
		return fmt.Errorf("SetUserInfoCache got nil *model.UserInfo")
	}
	var (
		key  = userInfoCacheKey(ui.Mid)
		conn = d.mcUser.Get(c)
	)
	defer conn.Close()

	if err = conn.Set(&memcache.Item{Key: key, Object: ui, Expiration: d.mcUserExpire, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, ui, err)
	}
	return
}

// DelInfoCache delete info cache.
func (d *Dao) DelInfoCache(c context.Context, mid int64) (err error) {
	var (
		key  = userInfoCacheKey(mid)
		conn = d.mcUser.Get(c)
	)
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
