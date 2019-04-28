package dao

import (
	"context"
	"encoding/json"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// SetIncomeCache set income cache
func (d *Dao) SetIncomeCache(c context.Context, key string, value map[string]interface{}) (err error) {
	v, err := json.Marshal(value)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	return d.setCacheKV(c, key, v, d.redisExpire)
}

// GetIncomeCache get income cache
func (d *Dao) GetIncomeCache(c context.Context, key string) (data map[string]interface{}, err error) {
	res, err := d.getCacheVal(c, key)
	if err != nil {
		log.Error("d.getCacheVal error(%v)", err)
		return
	}
	if res == nil {
		return
	}
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", res, err)
	}
	return
}

// DelCacheKey del redis key
func (d *Dao) DelCacheKey(c context.Context, key string) (err error) {
	return d.delCacheKey(c, key)
}

func (d *Dao) setCacheKV(c context.Context, key string, value []byte, expire int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if err = conn.Send("SET", key, value); err != nil {
		log.Error("conn.Send(SET, %s, %s) error(%v)", key, value, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.redisExpire, err)
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

func (d *Dao) getCacheVal(c context.Context, key string) (res []byte, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if res, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			res, err = nil, nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	return
}

func (d *Dao) delCacheKey(c context.Context, key string) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}
