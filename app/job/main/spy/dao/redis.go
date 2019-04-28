package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyWaitBlock   = "wb_" // b_batch_no wait block
	_keyBlock       = "bl_" // b_batch_no  block
	_preLock        = "lk_"
	_keyUniqueCheck = "uc:"
	times           = 3
)

// keyWaitBlock return block cache key.
func keyWaitBlock(batchNo int64) string {
	return _keyWaitBlock + fmt.Sprintf("%d", batchNo)
}

// keyBlock return block cache key.
func keyBlock() string {
	return _keyBlock
}

func lockKey(key string) string {
	return _preLock + key
}

func uniqueCheckKey(uuid string) string {
	return _keyUniqueCheck + uuid
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

// DelBlockCache delete the wait block redis.
func (d *Dao) DelBlockCache(c context.Context, batchNo int64, mid int64) (err error) {
	var (
		key  = keyWaitBlock(batchNo)
		args = []interface{}{key, mid}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", args...); err != nil {
		log.Error("conn.Send(ZREM %s,%v) error(%v)", key, mid, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}

//SetNXLockCache redis lock.
func (d *Dao) SetNXLockCache(c context.Context, k string, times int64) (res bool, err error) {
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
		if _, err = redis.Bool(conn.Do("EXPIRE", key, times)); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, times, err)
			return
		}
	}
	return
}

//DelLockCache del lock cache.
func (d *Dao) DelLockCache(c context.Context, k string) (err error) {
	var (
		key  = lockKey(k)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(del,%v) err(%v)", key, err)
	}
	return
}

//AddBlockCache add block cache.
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

// SetBlockCache block.
func (d *Dao) SetBlockCache(c context.Context, mids []int64) (err error) {
	var (
		key  = keyBlock()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, mid := range mids {
		if err = conn.Send("SADD", key, mid); err != nil {
			log.Error("SADD conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.expire); err != nil {
		log.Error("EXPIRE conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(mids); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("SetBlockCache Receive error(%v)", err)
			return
		}
	}
	return
}

//SPOPBlockCache pop mid.
func (d *Dao) SPOPBlockCache(c context.Context) (mid int64, err error) {
	var (
		key  = keyBlock()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if mid, err = redis.Int64(conn.Do("SPOP", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("SPOP conn.Do(%s,%v) err(%v)", key, err)
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

// PfaddCache SetNX.
func (d *Dao) PfaddCache(c context.Context, uuid string) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := uniqueCheckKey(uuid)
	if err = conn.Send("SETNX", key, 1); err != nil {
		log.Error("SETNX conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.msgUUIDExpire); err != nil {
		log.Error("conn.Send(EXPIRE) error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("DelLock conn.Flush() error(%v)", err)
		return
	}
	if ok, err = redis.Bool(conn.Receive()); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// TTL get redis cache ttl.
func (d *Dao) TTL(c context.Context, key string) (ttl int64, err error) {
	conn := d.redis.Get(c)
	ttl, err = redis.Int64(conn.Do("TTL", key))
	conn.Close()
	return
}
