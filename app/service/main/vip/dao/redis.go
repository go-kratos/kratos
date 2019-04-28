package dao

import (
	"context"
	"fmt"

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

//DelRedisCache del redis cache.
func (d *Dao) DelRedisCache(c context.Context, mid int64) (err error) {
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

//FrozenTime get user frozen time
func (d *Dao) FrozenTime(c context.Context, mid int64) (frozenTime int64, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if frozenTime, err = redis.Int64(conn.Do("ZSCORE", _frozenDelayQueue, mid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}
