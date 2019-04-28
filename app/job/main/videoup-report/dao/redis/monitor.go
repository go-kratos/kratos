package redis

import (
	"context"
	"go-common/library/log"
	"time"
)

// AddMonitorStats add stay stats
func (d *Dao) AddMonitorStats(c context.Context, key string, oid int64) (err error) {
	var (
		conn = d.secondary.Get(c)
		now  = time.Now().Unix()
		age  = 7 * 24 * 60 * 60
	)
	defer conn.Close()
	if _, err = conn.Do("ZADD", key, now, oid); err != nil {
		log.Error("conn.Do(ZADD, %s, %d, %d) error(%v)", key, now, oid, err)
		return
	}
	if _, err = conn.Do("EXPIRE", key, age); err != nil {
		log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, age, err)
	}
	return
}

// RemMonitorStats remove stay stats
func (d *Dao) RemMonitorStats(c context.Context, key string, oid int64) (err error) {
	var (
		conn = d.secondary.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, oid); err != nil {
		log.Error("conn.Do(ZREM, %s, %d) error(%v)", key, oid, err)
	}
	return
}

// ClearMonitorStats clear expire stats
func (d *Dao) ClearMonitorStats(c context.Context, key string) (err error) {
	var (
		conn = d.secondary.Get(c)
		now  = time.Now().Unix()
		min  int64
		max  = now - 7*24*60*60
	)
	defer conn.Close()
	if _, err = conn.Do("ZREMRANGEBYSCORE", key, min, max); err != nil {
		log.Error("conn.Do(ZREMRANGEBYSCORE, %s, %d, %d) error(%v)", key, min, max, err)
	}
	return
}
