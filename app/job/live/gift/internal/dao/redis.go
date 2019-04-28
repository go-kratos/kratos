package dao

import (
	"context"
	"fmt"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"time"

	"github.com/satori/go.uuid"
)

func bagIDCache(uid, giftID, expireAt int64) string {
	return fmt.Sprintf("bag_id:%d:%d:%d", uid, giftID, expireAt)
}

// SetBagIDCache SetBagIDCache
func (d *Dao) SetBagIDCache(ctx context.Context, uid, giftID, expireAt, bagID, expire int64) (err error) {
	key := bagIDCache(uid, giftID, expireAt)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key, expire, bagID)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

func bagListKey(uid int64) string {
	return fmt.Sprintf("bag_list:%d", uid)
}

// ClearBagListCache ClearBagListCache
func (d *Dao) ClearBagListCache(ctx context.Context, uid int64) (err error) {
	key := bagListKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	return
}

// GetBagIDCache GetBagIDCache
func (d *Dao) GetBagIDCache(ctx context.Context, uid, giftID, expireAt int64) (bagID int64, err error) {
	key := bagIDCache(uid, giftID, expireAt)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	bagID, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	return
}

func bagNumKey(uid, giftID, expireAt int64) string {
	return fmt.Sprintf("bag_num:%d:%d:%d", uid, giftID, expireAt)
}

// SetBagNumCache SetBagNumCache
func (d *Dao) SetBagNumCache(ctx context.Context, uid, giftID, expireAt, giftNum, expire int64) (err error) {
	key := bagNumKey(uid, giftID, expireAt)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key, expire, giftNum)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

//Lock Lock
func (d *Dao) Lock(ctx context.Context, key string, ttl int, retry int, retryDelay int) (gotLock bool, lockValue string, err error) {

	if retry <= 0 {
		retry = 1
	}
	lockValue = uuid.NewV4().String()
	retryTimes := 0
	conn := d.redis.Get(ctx)
	defer conn.Close()

	realKey := lockKey(key)

	for ; retryTimes < retry; retryTimes++ {
		var res interface{}
		res, err = conn.Do("SET", realKey, lockValue, "PX", ttl, "NX")
		if err != nil {
			log.Error("redis_lock failed:%s:%v", realKey, err)
			break
		}

		if res != nil {
			gotLock = true
			break
		}
		time.Sleep(time.Duration(retryDelay) * time.Millisecond)
	}
	return
}

func lockKey(key string) string {
	return fmt.Sprintf("gift_job_lock:%s", key)
}
