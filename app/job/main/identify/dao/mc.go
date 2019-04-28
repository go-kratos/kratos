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

// CookieCache get cookie info from cache
func (d *Dao) CookieCache(c context.Context, session string) (res *model.Cookie, err error) {
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
	res = new(model.Cookie)
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%v) error(%v)", string(item.Value), err)
	}
	return
}

// DelCookieCache del cache.
func (d *Dao) DelCookieCache(c context.Context, session string) (err error) {
	conn := d.authMC.Get(c)
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

// TokenCache get token from cache
func (d *Dao) TokenCache(c context.Context, sd string) (res *model.Token, err error) {
	key := akKey(sd)
	conn := d.authMC.Get(c)
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
	conn := d.authMC.Get(c)
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
