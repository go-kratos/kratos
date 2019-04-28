package dao

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

//TryLock ...
func (d *Dao) TryLock(c context.Context, key string, value string, timeout int) (ok bool, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	_, err = redis.String(conn.Do("SET", key, value, "EX", timeout, "NX"))
	if err == redis.ErrNil {
		log.Info("TryLock redis key(%s) is ErrNil!", key)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnLock ...
func (d *Dao) UnLock(c context.Context, key string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return
}
