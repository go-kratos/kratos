package anchorReward

import (
	"context"
	"go-common/library/cache/redis"

	"go-common/library/log"
)

const (
	_preLock        = "lk_"
	_preExpireCount = "ec_"
)

func lockKey(key string) string {
	return _preLock + key
}

func expireCountKey(key string) string {
	return _preExpireCount + key
}

//DelLockCache del lock cache.
func (d *Dao) DelLockCache(c context.Context, k string) (err error) {
	var (
		key  = lockKey(k)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("DelLockCache.conn.Do(del,%v) err(%v)", key, err)
	} else {
		log.Info("DelLockCache.conn.Do(del,%v)", key)
	}
	return
}

// GetExpireCountCache .
func (d *Dao) GetExpireCountCache(c context.Context, k string) (count int64, err error) {
	var (
		key  = expireCountKey(k)
		conn = d.redis.Get(c)
	)

	//spew.Dump(k)
	defer conn.Close()
	item, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			count = 0
			err = nil
		} else {
			log.Error("conn.Do(GET %s) error(%v)", key, err)
			return 0, err
		}
	}
	count = int64(item)

	return
}

// AddExpireCountCache .
func (d *Dao) AddExpireCountCache(c context.Context, k string, times int64) (err error) {
	var (
		key  = expireCountKey(k)
		conn = d.redis.Get(c)
	)
	//spew.Dump(k)
	defer conn.Close()
	if _, err = conn.Do("INCR", key); err != nil {
		log.Error("conn.Do(incr,%v) err(%v)", key, err)
	}

	if _, err = redis.Bool(conn.Do("EXPIRE", key, times)); err != nil {
		log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, times, err)
		return
	}
	return

}

// ClearExpireCountCache .
func (d *Dao) ClearExpireCountCache(c context.Context, k string) (err error) {
	var (
		key  = expireCountKey(k)
		conn = d.redis.Get(c)
	)
	//spew.Dump(k)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(del,%v) err(%v)", key, err)
		return
	}

	return
}

//SetNxLock redis lock.
func (d *Dao) SetNxLock(c context.Context, k string, times int64) (res bool, err error) {
	var (
		key  = lockKey(k)
		conn = d.redis.Get(c)
	)

	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", key, "1")); err != nil {
		log.Error("conn.Do(SETNX(%d)) error(%v)", key, err)
		return
	}
	//spew.Dump(res, err )
	if res {
		if _, err = redis.Bool(conn.Do("EXPIRE", key, times)); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, times, err)
			return
		}
		log.Info("conn.Do(EXPIRE, %s, %d) ", key, times)
	}
	return
}
