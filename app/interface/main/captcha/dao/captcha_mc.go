package dao

import (
	"context"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

var (
	_defaultCode = "@@@@"
)

// AddTokenCache add token redis cache.
func (d *Dao) AddTokenCache(c context.Context, key string, ttl int32) (err error) {
	conn := d.memcache.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Value: []byte(_defaultCode), Expiration: ttl, Flags: memcache.FlagRAW}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", key, _defaultCode, err)
	}
	return
}

// UpdateTokenCache update token cache.
func (d *Dao) UpdateTokenCache(c context.Context, token, code string, ttl int32) (err error) {
	conn := d.memcache.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: token, Value: []byte(code), Expiration: ttl, Flags: memcache.FlagRAW}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s,%v) error(%v)", token, code, err)
	}
	return
}

// CaptchaCache get captcha cache.
func (d *Dao) CaptchaCache(c context.Context, token string) (code string, isInit bool, err error) {
	conn := d.memcache.Get(c)
	defer conn.Close()
	item, err := conn.Get(token)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(GET %s) error(%v)", token, err)
		}
		return
	}
	if err = conn.Scan(item, &code); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
		return
	}
	isInit = code == _defaultCode
	return
}

// DelCaptchaCache delete captcha cache.
func (d *Dao) DelCaptchaCache(c context.Context, token string) (err error) {
	conn := d.memcache.Get(c)
	defer conn.Close()
	if err = conn.Delete(token); err != nil {
		log.Error("conn.Delete(%s) error(%v)", token, err)
	}
	return
}
