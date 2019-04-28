package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/web/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyArchiveFmt = "va_%d"
)

func keyArchive(aid int64) string {
	return fmt.Sprintf(_keyArchiveFmt, aid)
}

// SetViewBakCache  set view archive page data to cache.
func (d *Dao) SetViewBakCache(c context.Context, aid int64, a *model.View) (err error) {
	key := keyArchive(aid)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = json.Marshal(a); err != nil {
		log.Error("SetViewBakCache json.Marshal(%v) error(%v)", a, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("SetViewBakCache conn.Send(SET,%s,%s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisArchiveBakExpire); err != nil {
		log.Error("SetViewBakCache conn.Send(EXPIRE,%s,%d) error(%v)", key, d.redisArchiveBakExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("SetViewBakCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("SetViewBakCache conn.Recevie(%d) error(%v0", i, err)
		}
	}
	return
}

// ViewBakCache get view archive  page data from cache.
func (d *Dao) ViewBakCache(c context.Context, aid int64) (rs *model.View, err error) {
	key := keyArchive(aid)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var values []byte
	if values, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Warn("ViewBakCache redis (%s) return nil ", key)
		} else {
			log.Error("ViewBakCache conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(values, &rs); err != nil {
		log.Error("ViewBakCache json.Unmarshal(%v) error(%v)", values, err)
	}
	return
}
