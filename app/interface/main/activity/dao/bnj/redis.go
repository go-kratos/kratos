package bnj

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

func resetKey(mid int64) string {
	return fmt.Sprintf("bnj_%d", mid)
}

func rewardKey(mid, subID int64, step int) string {
	return fmt.Sprintf("bnj_rwd_%d_%d_%d", mid, subID, step)
}

// CacheResetCD .
func (d *Dao) CacheResetCD(c context.Context, mid int64, cd int32) (bool, error) {
	resetCD := d.resetExpire
	if cd > 0 {
		resetCD = cd
	}
	return d.setNXLockCache(c, resetKey(mid), resetCD)
}

// TTLResetCD get reset cd ttl
func (d *Dao) TTLResetCD(c context.Context, mid int64) (ttl int64, err error) {
	key := resetKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ttl, err = redis.Int64(conn.Do("TTL", key)); err != nil {
		log.Error("TTLResetCD conn.Do(TTL, %s), error(%v)", key, err)
	}
	return
}

// CacheHasReward .
func (d *Dao) CacheHasReward(c context.Context, mid, subID int64, step int) (bool, error) {
	return d.setNXLockCache(c, rewardKey(mid, subID, step), d.rewardExpire)
}

// DelCacheHasReward .
func (d *Dao) DelCacheHasReward(c context.Context, mid, subID int64, step int) error {
	return d.delNXLockCache(c, rewardKey(mid, subID, step))
}

// HasReward .
func (d *Dao) HasReward(c context.Context, mid, subID int64, step int) (res bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := rewardKey(mid, subID, step)
	if res, err = redis.Bool(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("HasReward conn.Do(GET(%s)) error(%v)", key, err)
	}
	return
}

func (d *Dao) setNXLockCache(c context.Context, key string, times int32) (res bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("SETNX", key, "1")); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(SETNX(%s)) error(%v)", key, err)
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

func (d *Dao) delNXLockCache(c context.Context, key string) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return
}
