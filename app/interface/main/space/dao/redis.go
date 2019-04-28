package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/space/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

func keyUpArt(mid int64) string {
	return fmt.Sprintf("%s_%d", "uat", mid)
}

func keyUpArc(mid int64) string {
	return fmt.Sprintf("%s_%d", "uar", mid)
}

// UpArtCache get up article cache.
func (d *Dao) UpArtCache(c context.Context, mid int64) (data *model.UpArtStat, err error) {
	var (
		value []byte
		key   = keyUpArt(mid)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	data = new(model.UpArtStat)
	if err = json.Unmarshal(value, &data); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetUpArtCache set up article cache.
func (d *Dao) SetUpArtCache(c context.Context, mid int64, data *model.UpArtStat) (err error) {
	var (
		bs   []byte
		key  = keyUpArt(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%v) error (%v)", data, err)
		return
	}
	err = setKvCache(conn, key, bs, d.upArtExpire)
	return
}

// UpArcCache get up archive cache.
func (d *Dao) UpArcCache(c context.Context, mid int64) (data *model.UpArcStat, err error) {
	var (
		value []byte
		key   = keyUpArc(mid)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	data = new(model.UpArcStat)
	if err = json.Unmarshal(value, &data); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", value, err)
	}
	return
}

// SetUpArcCache set up archive cache.
func (d *Dao) SetUpArcCache(c context.Context, mid int64, data *model.UpArcStat) (err error) {
	var (
		bs   []byte
		key  = keyUpArc(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%v) error (%v)", data, err)
		return
	}
	err = setKvCache(conn, key, bs, d.upArcExpire)
	return
}

func setKvCache(conn redis.Conn, key string, value []byte, expire int32) (err error) {
	if err = conn.Send("SET", key, value); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, string(value), err)
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
