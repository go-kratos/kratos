package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
)

const (
	_keyVersionID = "abtest:versionid:%d"
)

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("PING")
	return
}

//RedisVersionID 获取redis中的分组版本
func (d *Dao) RedisVersionID(c context.Context, group int) (ver int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	ver, err = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, group)))
	return
}

//SetnxRedisVersionID 使用v设置redis中的版本号
func (d *Dao) SetnxRedisVersionID(c context.Context, group int, v int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("SETNX", fmt.Sprintf(_keyVersionID, group), v)
	return
}

//UpdateRedisVersionID 使用v更新redis中的分组版本
func (d *Dao) UpdateRedisVersionID(c context.Context, group int, v int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("SETEX", fmt.Sprintf(_keyVersionID, group), d.verifyExpire, v)
	return
}
