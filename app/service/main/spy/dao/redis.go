package dao

import (
	"context"
	"fmt"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyWaitBlock = "wb_" // b_batch_no wait block
	_preLock      = "lk_"
)

// keyWaitBlock return block cache key.
func keyWaitBlock(batchNo int64) string {
	return _keyWaitBlock + fmt.Sprintf("%d", batchNo)
}

func lockKey(key int64) string {
	return fmt.Sprintf("%s%d", _preLock, key)
}

// AddBlockCache add block cache.
func (d *Dao) AddBlockCache(c context.Context, mid int64, score int8, blockNo int64) (err error) {
	var (
		key = keyWaitBlock(blockNo)
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, score, mid); err != nil {
		log.Error("conn.Send(ZADD %s,%d,%d) error(%v)", key, score, mid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
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

// BlockMidCache get wait block mids.
func (d *Dao) BlockMidCache(c context.Context, batchNo int64, num int64) (res []int64, err error) {
	var (
		conn = d.redis.Get(c)
		key  = keyWaitBlock(batchNo)
	)
	defer conn.Close()
	if res, err = redis.Int64s(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "LIMIT", 0, num)); err != nil {
		log.Error("redis(ZREVRANGEBYSCORE %s,%d) error(%v)", key, num, err)
		return
	}
	return
}

//SetNXLockCache redis lock.
func (d *Dao) SetNXLockCache(c context.Context, k int64) (res bool, err error) {
	var (
		key  = lockKey(k)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", key, "1")); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(SETNX(%d)) error(%v)", key, err)
			return
		}
	}
	if res {
		if _, err = redis.Bool(conn.Do("EXPIRE", key, d.verifyExpire)); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, d.verifyExpire, err)
			return
		}
	}
	return
}

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
