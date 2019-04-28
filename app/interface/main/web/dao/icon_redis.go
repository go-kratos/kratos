package dao

import (
	"context"
	"encoding/json"

	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_indexIconKey    = "iik"
	_indexIconBakKey = "b_iik"
)

// SetIndexIconCache set index icon cache and bak cache
func (d *Dao) SetIndexIconCache(c context.Context, data []*resmdl.IndexIcon) (err error) {
	key := _indexIconKey
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setIndexIconCache(conn, key, d.redisRcExpire, data); err != nil {
		return
	}
	key = _indexIconBakKey
	connBak := d.redisBak.Get(c)
	err = d.setIndexIconCache(connBak, key, d.redisRcBakExpire, data)
	connBak.Close()
	return
}

func (d *Dao) setIndexIconCache(conn redis.Conn, key string, expire int32, data []*resmdl.IndexIcon) (err error) {
	var bs []byte
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%v) error (%v)", data, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// IndexIconCache get index icon cache
func (d *Dao) IndexIconCache(c context.Context) (res []*resmdl.IndexIcon, err error) {
	key := _indexIconKey
	conn := d.redis.Get(c)
	defer conn.Close()
	res, err = d.indexIconCache(conn, key)
	return
}

// IndexIconBakCache get index icon bak cache
func (d *Dao) IndexIconBakCache(c context.Context) (res []*resmdl.IndexIcon, err error) {
	d.cacheProm.Incr("indexicon_remote_cache")
	key := _indexIconBakKey
	conn := d.redisBak.Get(c)
	defer conn.Close()
	res, err = d.indexIconCache(conn, key)
	return
}

func (d *Dao) indexIconCache(conn redis.Conn, key string) (res []*resmdl.IndexIcon, err error) {
	var value []byte
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(value, &res); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}
