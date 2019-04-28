package databus

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/main/videoup/model/message"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixMsgInfo = "videoup_service_msg"
	_preLock       = "lock_"
)

func lockKey(key string) string {
	return fmt.Sprintf("%s%s", _preLock, key)
}

// PopMsgCache get databus message from redis
func (d *Dao) PopMsgCache(c context.Context) (msg *message.Videoup, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _prefixMsgInfo)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(LPOP, %s) error(%v)", _prefixMsgInfo, err)
		}
		return
	}
	msg = &message.Videoup{}
	if err = json.Unmarshal(bs, msg); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
	}
	return
}

// PushMsgCache add message into redis.
func (d *Dao) PushMsgCache(c context.Context, msg *message.Videoup) (err error) {
	var (
		bs   []byte
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(msg); err != nil {
		log.Error("json.Marshal(%s) error(%v)", bs, err)
		return
	}
	if _, err = conn.Do("RPUSH", _prefixMsgInfo, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s) error(%v)", bs, err)
	}
	return
}

//Lock .
func (d *Dao) Lock(ctx context.Context, key string, ttl int) (gotLock bool, err error) {
	var lockValue = "1"
	conn := d.redis.Get(ctx)
	defer conn.Close()
	realKey := lockKey(key)
	var res interface{}
	//ttl 毫秒(PX)  NX 其实就是 SetNX功能
	res, err = conn.Do("SET", realKey, lockValue, "PX", ttl, "NX")
	if err != nil {
		log.Error("redis_lock failed:%s:%s", realKey, err.Error())
		return
	}
	if res != nil {
		gotLock = true
	}
	return
}

//UnLock .
func (d *Dao) UnLock(ctx context.Context, key string) (err error) {
	realKey := lockKey(key)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", realKey)
	return
}
