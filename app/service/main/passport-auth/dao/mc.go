package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/passport-auth/model"
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

// SetCookieCache set cookie info to cache
func (d *Dao) SetCookieCache(c context.Context, session string, res *model.Cookie) (err error) {
	key := ckKey(session)
	conn := d.mc.Get(c)
	defer conn.Close()
	if res.Expires < 0 {
		log.Error("auth expire error(expires:%d)", res.Expires)
		return
	}
	item := &memcache.Item{Key: key, Object: res, Flags: memcache.FlagProtobuf, Expiration: int32(res.Expires)}
	if err = conn.Set(item); err != nil {
		log.Error("auth set error(%s,%d,%v)", key, res.Expires, err)
	}
	return
}

// CookieCache get cookie info from cache
func (d *Dao) CookieCache(c context.Context, session string) (res *model.Cookie, err error) {
	key := ckKey(session)
	conn := d.mc.Get(c)
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
	res = new(model.Cookie)
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%v) error(%v)", string(item.Value), err)
	}
	return
}

// DelCookieCache del cache.
func (d *Dao) DelCookieCache(c context.Context, session string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(ckKey(session)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%s) error(%v)", ckKey(session), err)
	}
	return
}

// SetTokenCache set token to cache
func (d *Dao) SetTokenCache(c context.Context, k string, res *model.Token) (err error) {
	key := akKey(k)
	conn := d.mc.Get(c)
	defer conn.Close()
	if res.Expires < 0 {
		log.Error("auth expire error(expires:%d)", res.Expires)
		return
	}
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     res,
		Flags:      memcache.FlagProtobuf,
		Expiration: int32(res.Expires),
	}); err != nil {
		log.Error("set token cache error(%s,%d,%v)", key, res.Expires, err)
	}
	return
}

// TokenCache get token from cache
func (d *Dao) TokenCache(c context.Context, sd string) (res *model.Token, err error) {
	key := akKey(sd)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	res = new(model.Token)
	if err = conn.Scan(r, res); err != nil {
		log.Error("conn.Scan(%v) error(%v)", string(r.Value), err)
	}
	return
}

// DelTokenCache del cache.
func (d *Dao) DelTokenCache(c context.Context, token string) (err error) {
	key := akKey(token)
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

// SetRefreshCache set refresh token to cache .
func (d *Dao) SetRefreshCache(c context.Context, refresh *model.Refresh) (err error) {
	key := rkKey(refresh.Refresh)
	conn := d.mc.Get(c)
	defer conn.Close()
	if refresh.Expires < 0 {
		log.Error("auth expire error(expires:%d)", refresh.Expires)
		return
	}
	if err := conn.Set(&memcache.Item{
		Key:        key,
		Object:     refresh,
		Flags:      memcache.FlagProtobuf,
		Expiration: int32(refresh.Expires),
	}); err != nil {
		log.Error("auth set error(%s,%d,%v)", key, refresh.Expires, err)
	}
	return
}

// RefreshCache get refresh token from cache
func (d *Dao) RefreshCache(c context.Context, refresh string) (res *model.Refresh, err error) {
	key := rkKey(refresh)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	res = new(model.Refresh)
	if err = conn.Scan(r, res); err != nil {
		log.Error("conn.Scan(%v) error(%v)", string(r.Value), err)
	}
	return
}

// DelRefreshCache del refresh token from cache
func (d *Dao) DelRefreshCache(c context.Context, refresh string) (err error) {
	key := akKey(refresh)
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

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        "ping",
		Value:      []byte{1},
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("conn.Set(ping, 1) error(%v)", err)
	}
	return
}
