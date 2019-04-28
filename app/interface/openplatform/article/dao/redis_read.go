package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// GetsetReadPing 设置并获取上次阅读心跳时间，不存在则返回0
func (d *Dao) GetsetReadPing(c context.Context, buvid string, aid int64, cur int64) (last int64, err error) {
	var (
		key  = readPingKey(buvid, aid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if last, err = redis.Int64(conn.Do("GETSET", key, cur)); err != nil && err != redis.ErrNil {
		log.Error("conn.Do(GETSET, %s, %d) error(%+v)", key, cur, err)
		return
	}
	if _, err = conn.Do("EXPIRE", key, d.redisReadPingExpire); err != nil {
		log.Error("conn.Do(EXPIRE, %s, %d) error(%+v)", key, cur, err)
		return
	}
	return
}

// AddReadPingSet 添加新的阅读记录
func (d *Dao) AddReadPingSet(c context.Context, buvid string, aid int64, mid int64, ip string, cur int64, source string) (err error) {
	var (
		key   = readPingSetKey()
		value = fmt.Sprintf("%s|%d|%d|%s|%d|%s", buvid, aid, mid, ip, cur, source)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("SADD", key, value); err != nil {
		log.Error("conn.Do(SADD, %s, %s) error(%+v)", key, value, err)
		return
	}
	if _, err = conn.Do("EXPIRE", key, d.redisReadSetExpire); err != nil {
		log.Error("conn.Do(EXPIRE, %s, %d) error(%+v)", key, cur, err)
		return
	}
	return
}
