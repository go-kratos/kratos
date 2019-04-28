package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixLimitDay     = "ld_%d_%d_%s"
	_prefixLimitBiz     = "lb_%d_%d_%s_%d"
	_prefixLimitNotLive = "lnl_%d_%s"
)

func limitDayKey(day string, app, mid int64) string {
	return fmt.Sprintf(_prefixLimitDay, app, mid, day)
}

func limitBizKey(day string, app, mid, biz int64) string {
	return fmt.Sprintf(_prefixLimitBiz, app, mid, day, biz)
}

func limitNotLiveKey(day string, mid int64) string {
	return fmt.Sprintf(_prefixLimitNotLive, mid, day)
}

// pingRedis ping redis.
func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		PromError("redis: ping remote")
		log.Error("remote redis: conn.Do(SET,PING,PONG) error(%v)", err)
	}
	return
}

// LimitDayCache gets limit cache by day & mid.
// 测试用，业务用不着
func (d *Dao) LimitDayCache(ctx context.Context, day string, app, mid int64) (count int, err error) {
	var (
		key  = limitDayKey(day, app, mid)
		conn = d.redis.Get(ctx)
	)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("GET", key)); err != nil {
		PromError("redis:LimitDayCache")
		log.Error("LimitDayCache(%s,%d,%d) error(%v)", day, app, mid, err)
	}
	return
}

// IncrLimitDayCache increases and gets limit cache by day & mid.
func (d *Dao) IncrLimitDayCache(ctx context.Context, day string, app, mid int64) (count int, err error) {
	var (
		key  = limitDayKey(day, app, mid)
		conn = d.redis.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Send("INCR", key); err != nil {
		PromError("redis:IncrLimitDayCache")
		log.Error("IncrLimitDayCache(%s,%d,%d) error(%v)", day, app, mid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisLimitDayExpire); err != nil {
		PromError("redis:IncrLimitDayCache:expire")
		log.Error("IncrLimitDayCache(%s,%d,%d) expire error(%v)", day, app, mid, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:IncrLimitDayCache:flush")
		log.Error("IncrLimitDayCache(%s,%d,%d) flush error(%v)", day, app, mid, err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		PromError("redis:IncrLimitDayCache:receive:incr")
		log.Error("IncrLimitDayCache(%s,%d,%d) receive incr error(%+v)", day, app, mid, err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		PromError("redis:IncrLimitDayCache:receive:expire")
		log.Error("IncrLimitDayCache(%s,%d,%d) receive expire error(%+v)", day, app, mid, err)
	}
	return
}

// IncrLimitBizCache increases and gets limit cache by day & mid & bisiness.
func (d *Dao) IncrLimitBizCache(ctx context.Context, day string, app, mid, biz int64) (count int, err error) {
	var (
		key  = limitBizKey(day, app, mid, biz)
		conn = d.redis.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Send("INCR", key); err != nil {
		PromError("redis:IncrLimitBizCache")
		log.Error("IncrLimitBizCache(%s,%d,%d,%d) error(%v)", day, app, mid, biz, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisLimitDayExpire); err != nil {
		PromError("redis:IncrLimitBizCache:expire")
		log.Error("IncrLimitBizCache(%s,%d,%d,%d) expire error(%v)", day, app, mid, biz, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:IncrLimitBizCache:flush")
		log.Error("IncrLimitBizCache(%s,%d,%d,%d) flush error(%v)", day, app, mid, biz, err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		PromError("redis:IncrLimitBizCache:receive:incr")
		log.Error("IncrLimitBizCache(%s,%d,%d,%d) receive incr error(%+v)", day, app, mid, biz, err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		PromError("redis:IncrLimitBizCache:receive:expire")
		log.Error("IncrLimitBizCache(%s,%d,%d,%d) receive expire error(%+v)", day, app, mid, biz, err)
	}
	return
}

// IncrLimitNotLiveCache increases and gets not live limit cache by day & mid.
func (d *Dao) IncrLimitNotLiveCache(ctx context.Context, day string, mid int64) (count int, err error) {
	var (
		key  = limitNotLiveKey(day, mid)
		conn = d.redis.Get(ctx)
	)
	defer conn.Close()
	if err = conn.Send("INCR", key); err != nil {
		PromError("redis:IncrLimitNotLiveCache")
		log.Error("IncrLimitNotLiveCache(%s,%d) error(%v)", day, mid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisLimitDayExpire); err != nil {
		PromError("redis:IncrLimitNotLiveCache:expire")
		log.Error("IncrLimitNotLiveCache(%s,%d) expire error(%v)", day, mid, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:IncrLimitNotLiveCache:flush")
		log.Error("IncrLimitNotLiveCache(%s,%d) flush error(%v)", day, mid, err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		PromError("redis:IncrLimitNotLiveCache:receive:incr")
		log.Error("IncrLimitNotLiveCache(%s,%d) receive incr error(%+v)", day, mid, err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		PromError("redis:IncrLimitNotLiveCache:receive:expire")
		log.Error("IncrLimitNotLiveCache(%s,%d) receive expire error(%+v)", day, mid, err)
	}
	return
}
