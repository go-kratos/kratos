package dao

import (
	"context"
	"math/rand"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

func randomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Lock .
func (d *Dao) Lock(ctx context.Context, realKey string, ttl int, retry int, retryDelay int) (gotLock bool, lockValue string, err error) {

	if retry <= 0 {
		retry = 1
	}
	lockValue = "locked:" + randomString(5)
	retryTimes := 0
	conn := d.redis.Get(ctx)
	defer conn.Close()

	for ; retryTimes < retry; retryTimes++ {
		var res interface{}
		res, err = conn.Do("SET", realKey, lockValue, "PX", ttl, "NX")
		if err != nil {
			log.Error("redis_lock failed:%s:%s", realKey, err.Error())
			break
		}

		if res != nil {
			gotLock = true
			break
		}
		time.Sleep(time.Duration(retryDelay * 1000))
	}
	return
}

// UnLock .
func (d *Dao) UnLock(ctx context.Context, realKey string, lockValue string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	res, err := redis.String(conn.Do("GET", realKey))
	if err != nil {
		log.Error("redis_unlock get error:%s:%v", realKey, err)
		return
	}
	if res != lockValue {
		err = ErrUnLockGet
		return
	}

	_, err = conn.Do("DEL", realKey)
	if err != nil {
		log.Error("redis_unlock del error:%s:%v", realKey, err)
	}
	return
}
