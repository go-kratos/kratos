package dao

import (
	"context"
	"encoding/json"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	_failList = "player_job_list"
)

func keyRetry() string {
	return _failList
}

// PushList rpush item to redis
func (d *Dao) PushList(c context.Context, a interface{}) (err error) {
	var bs []byte
	conn := d.redis.Get(c)
	defer conn.Close()
	if bs, err = json.Marshal(a); err != nil {
		err = errors.Wrapf(err, "%v", a)
		return
	}
	if _, err = conn.Do("RPUSH", keyRetry(), bs); err != nil {
		err = errors.Wrapf(err, "conn.Do(RPUSH,%s,%s)", keyRetry(), bs)
	}
	return
}

// PopList lpop item from redis
func (d *Dao) PopList(c context.Context) (bs []byte, err error) {
	conn := d.redis.Get(c)
	if bs, err = redis.Bytes(conn.Do("LPOP", keyRetry())); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrapf(err, "redis.Bytes(conn.Do(LPOP, %s))", keyRetry())
		}
	}
	conn.Close()
	return
}

// PingRedis is
func (d *Dao) PingRedis(c context.Context) (err error) {
	var conn = d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
