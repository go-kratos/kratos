package newcomer

import (
	"context"
	"fmt"
	"go-common/library/log"
)

const (
	_preLock = "creative_task_"
)

func lockKey(key string) string {
	return fmt.Sprintf("%s%s", _preLock, key)
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
		log.Error("receive_lock failed:%s:%s", realKey, err.Error())
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
