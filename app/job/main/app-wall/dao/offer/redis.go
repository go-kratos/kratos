package offer

import (
	"context"
	"encoding/json"

	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

const (
	_failList = "act_list"
)

func keyRetry() string {
	return _failList
}

// PushFail rpush fail item to redis
func (d *Dao) PushFail(c context.Context, a interface{}) (err error) {
	var bs []byte
	conn := d.redis.Get(c)
	key := keyRetry()
	defer conn.Close()
	if bs, err = json.Marshal(a); err != nil {
		err = errors.Wrapf(err, "%v", a)
		return
	}
	if _, err = conn.Do("RPUSH", key, bs); err != nil {
		err = errors.Wrapf(err, "conn.Do(RPUSH,%s,%s)", key, bs)
	}
	return
}

// PopFail lpop fail item from redis
func (d *Dao) PopFail(c context.Context) (bs []byte, err error) {
	conn := d.redis.Get(c)
	key := keyRetry()
	if bs, err = redis.Bytes(conn.Do("LPOP", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			err = errors.Wrapf(err, "redis.Bytes(conn.Do(LPOP, %s))", key)
		}
	}
	conn.Close()
	return
}

func (d *Dao) PingRedis(c context.Context) (err error) {
	var conn = d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
