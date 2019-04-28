package dao

import (
	"context"
	"fmt"
	"go-common/app/interface/main/passport-login/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

func ckKey(session string) string {
	return fmt.Sprintf("ck_%s", session)
}

func akKey(token string) string {
	return fmt.Sprintf("ak_%s", token)
}

func rkKey(refresh string) string {
	return fmt.Sprintf("rk_%s", refresh)
}

// CookieCache get cookie info from cache
func (d *Dao) CookieCache(c context.Context, session string) (res *model.CookieProto, err error) {
	key := ckKey(session)
	conn := d.authMC.Get(c)
	defer conn.Close()
	var item *memcache.Item
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	res = new(model.CookieProto)
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%v) error(%v)", string(item.Value), err)
	}
	return
}

// TokenCache get token info from cache
func (d *Dao) TokenCache(c context.Context, token string) (res *model.TokenProto, err error) {
	key := akKey(token)
	conn := d.authMC.Get(c)
	defer conn.Close()
	var item *memcache.Item
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	res = new(model.TokenProto)
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%v) error(%v)", string(item.Value), err)
	}
	return
}

// SetCookieCache set cookie info to cache
func (d *Dao) SetCookieCache(c context.Context, res *model.CookieProto) (err error) {
	key := ckKey(res.Session)
	conn := d.authMC.Get(c)
	defer conn.Close()
	if res.Expires < 0 {
		log.Error("auth expire error(expires:%d)", res.Expires)
		return
	}
	item := &memcache.Item{Key: key, Object: res, Flags: memcache.FlagProtobuf, Expiration: int32(d.authMCExpire)}
	if err = conn.Set(item); err != nil {
		log.Error("auth set error(%s,%d,%v)", key, res.Expires, err)
	}
	return
}

// DelCookieCache delete cookie cache
func (d *Dao) DelCookieCache(c context.Context, session string) (err error) {
	key := ckKey(session)
	conn := d.authMC.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", key, err)
		return
	}
	return
}

// SetTokenCache set token to cache
func (d *Dao) SetTokenCache(c context.Context, res *model.TokenProto) (err error) {
	key := akKey(res.Token)
	conn := d.authMC.Get(c)
	defer conn.Close()
	if res.Expires < 0 {
		log.Error("auth expire error(expires:%d)", res.Expires)
		return
	}
	item := &memcache.Item{Key: key, Object: res, Flags: memcache.FlagProtobuf, Expiration: int32(d.authMCExpire)}
	if err = conn.Set(item); err != nil {
		log.Error("set token cache error(%s,%d,%v)", key, res.Expires, err)
	}
	return
}

// SetRefreshCache set refresh token to cache .
func (d *Dao) SetRefreshCache(c context.Context, refresh *model.RefreshProto) (err error) {
	key := rkKey(refresh.Refresh)
	conn := d.authMC.Get(c)
	defer conn.Close()
	if refresh.Expires < 0 {
		log.Error("auth expire error(expires:%d)", refresh.Expires)
		return
	}
	item := &memcache.Item{Key: key, Object: refresh, Flags: memcache.FlagProtobuf, Expiration: int32(d.authMCExpire)}
	if err := conn.Set(item); err != nil {
		log.Error("auth set error(%s,%d,%v)", key, refresh.Expires, err)
	}
	return
}

// DelTokenCache delete token cache
func (d *Dao) DelTokenCache(c context.Context, token string) (err error) {
	key := akKey(token)
	conn := d.authMC.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", key, err)
		return
	}
	return
}
