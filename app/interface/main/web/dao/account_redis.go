package dao

import (
	"context"
	"fmt"

	"encoding/json"
	"go-common/app/interface/main/web/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyCardFmt = "ac_%d"
)

func keyCard(mid int64) string {
	return fmt.Sprintf(_keyCardFmt, mid)
}

// SetCardBakCache set card data to cache.
func (d *Dao) SetCardBakCache(c context.Context, mid int64, rs *model.Card) (err error) {
	var bs []byte
	key := keyCard(mid)
	if bs, err = json.Marshal(rs); err != nil {
		log.Error("json.Marshal(%v) error(%v)", rs, err)
		return
	}
	err = d.commonSetBakCache(c, key, bs)
	return
}

//CardBakCache get card data from cache.
func (d *Dao) CardBakCache(c context.Context, mid int64) (rs *model.Card, err error) {
	key := keyCard(mid)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var values []byte
	if values, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Warn("CardBakCache (%s) return nil", key)
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(values, &rs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", values, err)
	}
	return
}

func (d *Dao) commonSetBakCache(c context.Context, key string, bs []byte) (err error) {
	conn := d.redisBak.Get(c)
	defer conn.Close()
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisCardBakExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s,%d) error(%v)", key, d.redisCardBakExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive(%d) error(%v)", i, err)
		}
	}
	return
}
