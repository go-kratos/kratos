package dao

import (
	"context"
	"math/rand"
	"strconv"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_keyPrefixInfoPB              = "pa2_"
	_keyPrefixTokenPB             = "pt_"
	_keyPrefixOriginMissMatchFlag = "m_"
	_keyPrefixOriginToken         = "pot_"

	_missMatchFlagExpireSeconds = 30 // 30 seconds
)

func keyInfoPB(mid int64) string {
	return _keyPrefixInfoPB + strconv.FormatInt(mid, 10)
}

func keyTokenPB(accessToken string) string {
	return _keyPrefixTokenPB + accessToken
}

func keyOriginToken(accessToken string) string {
	return _keyPrefixOriginToken + accessToken
}

func keyOriginMissMatchFlag(identify string) string {
	return _keyPrefixOriginMissMatchFlag + identify
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

// SetInfoCache set info into cache.
func (d *Dao) SetInfoCache(c context.Context, info *model.Info) (err error) {
	item := &memcache.Item{
		Key:        keyInfoPB(info.Mid),
		Object:     info,
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

// InfoCache get info cache from cache.
func (d *Dao) InfoCache(c context.Context, mid int64) (info *model.Info, err error) {
	key := keyInfoPB(mid)
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
	info = new(model.Info)
	if err = conn.Scan(item, info); err != nil {
		log.Error("conn.Scan(%s, %s) error(%v)", key, item.Value, err)
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

// TokenCache get token from cache.
func (d *Dao) TokenCache(c context.Context, accessToken string) (token *model.Perm, err error) {
	key := keyTokenPB(accessToken)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	token = new(model.Perm)
	if err = conn.Scan(item, token); err != nil {
		log.Error("conn.Scan(%s, %s) error(%v)", item.Key, item.Value, err)
	}
	return
}

// DelTokenCache delete token from cache.
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

// SetOriginMissMatchFlagCache set origin miss match flag cache.
func (d *Dao) SetOriginMissMatchFlagCache(c context.Context, identify string, flag []byte) (err error) {
	item := &memcache.Item{
		Key:        keyOriginMissMatchFlag(identify),
		Value:      flag,
		Expiration: _missMatchFlagExpireSeconds,
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", item.Key, err)
	}
	return
}

// OriginMissMatchFlagCache get origin miss match flag.
func (d *Dao) OriginMissMatchFlagCache(c context.Context, identify string) (res []byte, err error) {
	key := keyOriginMissMatchFlag(identify)
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
	res = item.Value
	return
}

// DelOriginMissMatchFlagCache delete origin miss match flag.
func (d *Dao) DelOriginMissMatchFlagCache(c context.Context, identify string) (err error) {
	key := keyOriginMissMatchFlag(identify)
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

// SetOriginTokenCache set passport token into cache.
func (d *Dao) SetOriginTokenCache(c context.Context, token *model.Token) (err error) {
	item := &memcache.Item{
		Key:        keyOriginToken(token.AccessToken),
		Object:     token,
		Flags:      memcache.FlagJSON,
		Expiration: 60 + rand.Int31n(600),
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", item.Key, err)
		return
	}
	return
}

// OriginTokenCache set passport token into cache.
func (d *Dao) OriginTokenCache(c context.Context, accessToken string) (token *model.Token, err error) {
	key := keyOriginToken(accessToken)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	token = new(model.Token)
	if err = conn.Scan(item, token); err != nil {
		log.Error("conn.Scan(%s, %s) error(%v)", item.Key, item.Value, err)
		return
	}
	return
}

// DelOriginTokenCache delete passport token from cache.
func (d *Dao) DelOriginTokenCache(c context.Context, accessToken string) (err error) {
	key := keyOriginToken(accessToken)
	conn := d.mc.Get(c)
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
