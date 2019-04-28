package dao

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-common/library/cache/redis"
)

const (
	_expShard       = 10000
	_expAddedPrefix = "ea_%s_%d_%d"
	_expire         = 86400
)

func expAddedKey(tp string, mid int64, day int) string {
	return fmt.Sprintf(_expAddedPrefix, tp, day, mid/_expShard)
}

func leader() (key string, value string) {
	value, _ = os.Hostname()
	key = fmt.Sprintf("member-job:leader:%d", time.Now().Day())
	return key, value
}

// SetExpAdded set user exp added of tp,
func (d *Dao) SetExpAdded(c context.Context, mid int64, day int, tp string) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("SETBIT", expAddedKey(tp, mid, day), mid%_expShard, 1); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", expAddedKey(tp, mid, day), _expire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if b, err = redis.Bool(conn.Receive()); err != nil {
		return
	}
	conn.Receive()
	return
}

// ExpAdded check if user had add exp.
func (d *Dao) ExpAdded(c context.Context, mid int64, day int, tp string) (b bool, err error) {
	conn := d.redis.Get(c)
	b, err = redis.Bool(conn.Do("GETBIT", expAddedKey(tp, mid, day), mid%_expShard))
	conn.Close()
	return
}

// LeaderEleciton eleciton job leader.
func (d *Dao) LeaderEleciton(c context.Context) (elected bool) {
	key, value := leader()
	conn := d.redis.Get(c)
	elected, _ = redis.Bool(conn.Do("SETNX", key, value))
	if elected {
		conn.Do("EXPIRE", key, 600)
	}
	conn.Close()
	return
}
