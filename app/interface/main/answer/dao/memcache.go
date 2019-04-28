package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/answer/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_answerTimeKey         = "v3_at_%d" // key of user's answer limit time
	_answerHistoryKey      = "v3_ah_%d" // key of user's answer history
	_answerQidListKey      = "v3_aqbi_%d"
	_answerExtraQidListKey = "v3_aql_%d_%d"
	_answerBlockKey        = "v3_ablk_%d" // key of user's answer block flag

	_answerHistory = "hid_%d" // ah_hid
)

var (
	_blockFlag = []byte("1")
)

func answerQidListKey(mid int64, t int8) (key string) {
	switch t {
	case model.BaseExtraPassQ:
		key = fmt.Sprintf(_answerExtraQidListKey, t, mid)
	case model.BaseExtraNoPassQ:
		key = fmt.Sprintf(_answerExtraQidListKey, t, mid)
	default:
		key = fmt.Sprintf(_answerQidListKey, mid)
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

// ExpireCache get user's answer stime and base answer error times cache.
func (d *Dao) ExpireCache(c context.Context, mid int64) (at *model.AnswerTime, err error) {
	key := fmt.Sprintf(_answerTimeKey, mid)
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
	at = &model.AnswerTime{}
	if err = conn.Scan(item, at); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
	}
	return
}

// SetExpireCache set user's answer stime and base answer error times cache.
func (d *Dao) SetExpireCache(c context.Context, mid int64, at *model.AnswerTime) (err error) {
	key := fmt.Sprintf(_answerTimeKey, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     at,
		Flags:      memcache.FlagJSON,
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, at, err)
	}
	return
}

// DelExpireCache delete user's answer stime and base answer error times cache.
func (d *Dao) DelExpireCache(c context.Context, mid int64) (err error) {
	key := fmt.Sprintf(_answerTimeKey, mid)
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

// HistoryCache get user's answer history cache
func (d *Dao) HistoryCache(c context.Context, mid int64) (ah *model.AnswerHistory, err error) {
	key := fmt.Sprintf(_answerHistoryKey, mid)
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
	ah = &model.AnswerHistory{}
	if err = conn.Scan(item, ah); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
	}
	return
}

// SetHistoryCache set user's answer history cache
func (d *Dao) SetHistoryCache(c context.Context, mid int64, ah *model.AnswerHistory) (err error) {
	key := fmt.Sprintf(_answerHistoryKey, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     ah,
		Flags:      memcache.FlagJSON,
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, ah, err)
	}
	return
}

// DelHistoryCache delete user's answer history cache
func (d *Dao) DelHistoryCache(c context.Context, mid int64) (err error) {
	key := fmt.Sprintf(_answerHistoryKey, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DelHistoryCache(%d),err:%+v", mid, err)
	}
	return
}

// IdsCache get user's base question ids
func (d *Dao) IdsCache(c context.Context, mid int64, t int8) (ids []int64, err error) {
	key := answerQidListKey(mid, t)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("d.IdsCache(%d,%d) error(%v)", mid, t, err)
		return
	}
	if err = conn.Scan(item, &ids); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
	}
	return
}

// SetIdsCache set user's base question ids
func (d *Dao) SetIdsCache(c context.Context, mid int64, ids []int64, t int8) (err error) {
	key := answerQidListKey(mid, t)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     ids,
		Flags:      memcache.FlagJSON,
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, ids, err)
	}
	return
}

// DelIdsCache delete user's base question ids
func (d *Dao) DelIdsCache(c context.Context, mid int64, t int8) (err error) {
	key := answerQidListKey(mid, t)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("DelIdsCache(%d,%d),err:%+v", mid, t, err)
	}
	return
}

// SetBlockCache set user's block.
func (d *Dao) SetBlockCache(c context.Context, mid int64) (err error) {
	key := fmt.Sprintf(_answerBlockKey, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Value:      _blockFlag,
		Expiration: d.answerBlockExpire,
	}); err != nil {
		log.Error("conn.Store(%s, %v) error(%v)", key, string(_blockFlag), err)
	}
	return
}

// CheckBlockCache check user's block.
func (d *Dao) CheckBlockCache(c context.Context, mid int64) (exist bool, err error) {
	key := fmt.Sprintf(_answerBlockKey, mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	_, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", key, err)
		return
	}
	exist = true
	return
}

// HidCache get user's answer history cache
func (d *Dao) HidCache(c context.Context, hid int64) (ah *model.AnswerHistory, err error) {
	key := fmt.Sprintf(_answerHistory, hid)
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
	ah = &model.AnswerHistory{}
	if err = conn.Scan(item, ah); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(item.Value), err)
	}
	return
}

// SetHidCache set user's answer history cache
func (d *Dao) SetHidCache(c context.Context, ah *model.AnswerHistory) (err error) {
	key := fmt.Sprintf(_answerHistory, ah.Hid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     ah,
		Flags:      memcache.FlagJSON,
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, ah, err)
	}
	return
}
