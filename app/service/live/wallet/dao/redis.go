package dao

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"time"
)

var (
	UnLockGetWrong   = "UnLockGetWrong"
	LockFailed       = "LockFailed"
	UserWalletPrefix = "user_wallet_lock_uid_"
	ErrUnLockGet     = errors.New(UnLockGetWrong)
	ErrLockFailed    = errors.New(LockFailed)
)

func (d *Dao) IsLockFailedError(err error) bool {
	return err == ErrLockFailed
}

func lockKey(k string) string {
	return "wallet_lock_key:" + k
}

/*
ttl ms
retry 重试次数
retryDelay us
*/
func (d *Dao) Lock(ctx context.Context, key string, ttl int, retry int, retryDelay int) (err error, gotLock bool, lockValue string) {

	if retry <= 0 {
		retry = 1
	}
	lockValue = "locked:" + randomString(5)
	retryTimes := 0
	conn := d.redis.Get(ctx)
	defer conn.Close()

	realKey := lockKey(key)

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

func (d *Dao) UnLock(ctx context.Context, key string, lockValue string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	realKey := lockKey(key)
	res, err := redis.String(conn.Do("GET", realKey))
	if err != nil {
		return
	}
	if res != lockValue {
		err = ErrUnLockGet
		return
	}

	_, err = conn.Do("DEL", realKey)

	return
}

func (d *Dao) ForceUnLock(ctx context.Context, key string) (err error) {
	realKey := lockKey(key)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", realKey)

	return

}

func (d *Dao) LockTransactionId(ctx context.Context, tid string) (err error) {
	err, gotLock, _ := d.Lock(ctx, tid, 300*1000, 0, 200000)
	if err != nil {
		return
	}
	if !gotLock {
		err = ErrLockFailed
	}
	return
}

func (d *Dao) LockUser(ctx context.Context, uid int64) (err error, gotLock bool, lockValue string) {
	lockTime := 600
	retry := 1
	retryDelay := 10

	return d.Lock(ctx, getUserLockKey(uid), lockTime*1000, retry, retryDelay*1000)
}

func (d *Dao) UnLockUser(ctx context.Context, uid int64, lockValue string) error {
	return d.UnLock(ctx, getUserLockKey(uid), lockValue)
}

func getUserLockKey(uid int64) string {
	return fmt.Sprintf("%s%v", UserWalletPrefix, uid)
}
