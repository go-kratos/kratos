package dao

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/web/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const _onlineListKey = "olk"

// OnlineListBakCache get online list bak cache.
func (d *Dao) OnlineListBakCache(c context.Context) (rs []*model.OnlineArc, err error) {
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var values []byte
	if values, err = redis.Bytes(conn.Do("GET", _onlineListKey)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Warn("OnlineListBakCache redis (%s) return nil ", _onlineListKey)
		} else {
			log.Error("OnlineListBakCache conn.Do(GET,%s) error(%v)", _onlineListKey, err)
		}
		return
	}
	if err = json.Unmarshal(values, &rs); err != nil {
		log.Error("OnlineListBakCache json.Unmarshal(%v) error(%v)", values, err)
	}
	return
}

// SetOnlineListBakCache set online list bak cache.
func (d *Dao) SetOnlineListBakCache(c context.Context, data []*model.OnlineArc) (err error) {
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = json.Marshal(data); err != nil {
		log.Error("SetOnlineListBakCache json.Marshal(%v) error(%v)", data, err)
		return
	}
	if err = conn.Send("SET", _onlineListKey, bs); err != nil {
		log.Error("SetOnlineListBakCache conn.Send(SET,%s,%s) error(%v)", _onlineListKey, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", _onlineListKey, d.redisOlListBakExpire); err != nil {
		log.Error("SetOnlineListBakCache conn.Send(EXPIRE,%s,%d) error(%v)", _onlineListKey, d.redisOlListBakExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("SetOnlineListBakCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("SetOnlineListBakCache conn.Recevie(%d) error(%v0", i, err)
		}
	}
	return
}
