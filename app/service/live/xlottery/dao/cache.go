package dao

import (
	"context"
	"time"

	stderr "errors"

	"github.com/pkg/errors"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// Incr incr
func (d *Dao) Incr(context context.Context, name string) int {
	conn := d.redis.Get(context)
	defer conn.Close()
	reply, _ := redis.Int(conn.Do("INCR", name))
	return reply
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
		return
	}
	if res != lockValue {
		err = ErrUnLockGet
		return
	}

	_, err = conn.Do("DEL", realKey)

	return
}

// Get get
func (d *Dao) Get(ctx context.Context, key string) (string, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.String(conn.Do("GET", key))

}

// GetInt64 GetInt64
func (d *Dao) GetInt64(ctx context.Context, key string) (int64, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.Int64(conn.Do("GET", key))
}

// Set Set
func (d *Dao) Set(ctx context.Context, key string, value interface{}) (bool, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err := redis.String(conn.Do("SET", key, value)); err != nil {
		return false, errors.WithStack(err)
	}
	return true, nil

}

// SetEx SetEx
func (d *Dao) SetEx(ctx context.Context, key string, value interface{}, ex int) (bool, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err := redis.String(conn.Do("SET", key, value, "EX", ex)); err != nil {
		return false, errors.WithStack(err)
	}
	return true, nil

}

//HMSet HMSet
func (d *Dao) HMSet(ctx context.Context, key string, value interface{}) (bool, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	arg := redis.Args{}.Add(key).AddFlat(value)
	if _, err := redis.String(conn.Do("HMSET", arg...)); err != nil {
		return false, errors.WithStack(err)
	}
	return true, nil
}

//Del Del
func (d *Dao) Del(ctx context.Context, key string) (int64, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.Int64(conn.Do("DEL", key))

}

//SIsMember SIsMember
func (d *Dao) SIsMember(ctx context.Context, key, value string) (bool, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.Bool(conn.Do("SISMEMBER", key, value))
}

//SAdd SAdd
func (d *Dao) SAdd(ctx context.Context, key, value string) (int, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.Int(conn.Do("SADD", key, value))
}

// Expire Expire
func (d *Dao) Expire(ctx context.Context, key string, ttl int) (bool, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.Bool(conn.Do("EXPIRE", key, ttl))
}

//ErrEmptyMap ErrEmptyMap
var ErrEmptyMap = stderr.New("empty map")

// HGetAll HGetAll
func (d *Dao) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	m, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, ErrEmptyMap
	}
	return m, nil

}

// SetWithNxEx SetWithNxEx
func (d *Dao) SetWithNxEx(ctx context.Context, key, value string, ex int64) (string, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.String(conn.Do("SET", key, value, "EX", ex, "NX"))
}
