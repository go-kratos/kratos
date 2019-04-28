package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	_frozenDelayQueue = "fdq"
	_accLogin         = "al:%d"
)

func accLoginKey(mid int64) string {
	return fmt.Sprintf(_accLogin, mid)
}

// Enqueue put frozen user to redis sortset queue
func (d *Dao) Enqueue(c context.Context, mid, score int64) (err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("ZADD", _frozenDelayQueue, score, mid); err != nil {
		err = errors.Wrap(err, "redis send zadd err(%+v)")
		return
	}
	if err = conn.Send("EXPIRE", _frozenDelayQueue, d.redisExpire); err != nil {
		err = errors.Wrap(err, "redis send expire err(%+v)")
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

// Dequeue get a frozen user from redis sortset queue
func (d *Dao) Dequeue(c context.Context) (mid []int64, err error) {
	var (
		conn = d.redis.Get(c)
		from = time.Now().Add(-1 * time.Minute).Unix()
		to   = time.Now().Unix()
	)
	defer conn.Close()
	if mid, err = redis.Int64s(conn.Do("ZRANGEBYSCORE", _frozenDelayQueue, from, to)); err != nil {
		err = errors.Wrap(err, "redis do zrevrangebyscore err")
		return
	}
	return
}

// RemQueue del a frozen user from redis sortset queue
func (d *Dao) RemQueue(c context.Context, mid int64) (err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREM", _frozenDelayQueue, mid); err != nil {
		err = errors.Wrap(err, "redis do zrem err")
		return
	}
	return
}

// AddLogginIP save user loggin info to sortset
func (d *Dao) AddLogginIP(c context.Context, mid int64, ip uint32) (err error) {
	var (
		conn = d.redis.Get(c)
		key  = accLoginKey(mid)
	)
	defer conn.Close()
	if err = conn.Send("SADD", key, ip); err != nil {
		err = errors.Wrap(err, "redis do zadd err")
		return
	}
	if err = conn.Send("EXPIRE", key, d.frozenExpire); err != nil {
		err = errors.Wrap(err, "redis send expire err(%+v)")
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			return
		}
	}
	return
}

//DelCache del cache.
func (d *Dao) DelCache(c context.Context, mid int64) (err error) {
	var (
		conn = d.redis.Get(c)
		key  = accLoginKey(mid)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// LoginCount get user recent loggin count
func (d *Dao) LoginCount(c context.Context, mid int64) (count int64, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int64(conn.Do("SCARD", accLoginKey(mid))); err != nil {
		err = errors.Wrap(err, "redis send zcard err(%+v)")
	}
	return
}
