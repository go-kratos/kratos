package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

func mailGroupKey(name string, t *model.Target, code string) string {
	return fmt.Sprintf("%s:%s", name, targetKey(t, code))
}

func targetKey(t *model.Target, code string) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s", t.Source, t.Product, t.Event, t.SubEvent, code)
}

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		log.Error("remote redis: conn.Do(SET,PING,PONG) error(%+v)", err)
	}
	conn.Close()
	return
}

// GetMailLock .
func (d *Dao) GetMailLock(c context.Context, name string, interval int, t *model.Target, code string) (ok bool, err error) {
	if name == "" || interval == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	key := mailGroupKey(name, t, code)
	if ok, err = redis.Bool(conn.Do("SETNX", key, "1")); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("d.redis.conn.Do(SETNX(%s)) error(%v)", key, err)
			return
		}
	}
	if ok {
		conn.Do("EXPIRE", key, interval)
	}
	return
}

// TargetIncr get current target error amount.
func (d *Dao) TargetIncr(c context.Context, t *model.Target, code string) (res int) {
	var (
		conn = d.redis.Get(c)
		err  error
		key  = targetKey(t, code)
	)
	defer conn.Close()
	if res, err = redis.Int(conn.Do("INCR", key)); err != nil {
		log.Error("d.redis.intr error(%+v)", err)
		return
	}
	if res == 1 {
		conn.Do("EXPIRE", key, t.Duration)
	}
	return
}
