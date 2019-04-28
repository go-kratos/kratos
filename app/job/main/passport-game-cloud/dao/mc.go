package dao

import (
	"context"
	"strconv"

	"go-common/app/job/main/passport-game-cloud/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_keyPrefixInfoPB  = "pa2_"
	_keyPrefixTokenPB = "pt_"
)

func keyInfoPB(mid int64) string {
	return _keyPrefixInfoPB + strconv.FormatInt(mid, 10)
}

func keyTokenPB(accessToken string) string {
	return _keyPrefixTokenPB + accessToken
}

// pingMC check connection success.
func (d *Dao) pingMC(c context.Context) (err error) {
	item := &memcache.Item{
		Key:        "ping",
		Value:      []byte{1},
		Expiration: d.mcExpire,
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DelInfoCache delete info cache.
func (d *Dao) DelInfoCache(c context.Context, mid int64) (err error) {
	key := keyInfoPB(mid)
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

// SetTokenCache set token into cache.
func (d *Dao) SetTokenCache(c context.Context, token *model.Perm) (err error) {
	item := &memcache.Item{
		Key:        keyTokenPB(token.AccessToken),
		Object:     token,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.mcExpire,
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", item.Key, err)
	}
	return
}

// DelTokenCache delete token cache.
func (d *Dao) DelTokenCache(c context.Context, accessToken string) (err error) {
	key := keyTokenPB(accessToken)
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
