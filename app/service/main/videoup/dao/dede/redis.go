package dede

import (
	"context"
	"encoding/json"

	"go-common/app/service/main/videoup/model/dede"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixPadInfo = "padinfo"
)

// PopPadInfoCache get padinfo from redis
func (d *Dao) PopPadInfoCache(c context.Context) (pad *dede.PadInfo, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _prefixPadInfo)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(LPOP, %s) error(%v)", _prefixPadInfo, err)
		}
		return
	}
	pad = &dede.PadInfo{}
	if err = json.Unmarshal(bs, pad); err != nil {
		log.Error("s.padproc json.Unmarshal error(%v)", err)
	}
	return
}

// PushPadCache add padinfo into redis.
func (d *Dao) PushPadCache(c context.Context, pad *dede.PadInfo) (err error) {
	var (
		bs   []byte
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(pad); err != nil {
		log.Error("json.Marshal(%s) error(%v)", bs, err)
		return
	}
	if _, err = conn.Do("RPUSH", _prefixPadInfo, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%v)", bs, err)
	}
	return
}
