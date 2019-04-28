package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_actionkey   = "ua_%d_%d" // hash key:mid_type field:oid_tid value:action
	_actionField = "%d_%d"
)

func keyAction(mid int64, typ int32) string {
	return fmt.Sprintf(_actionkey, mid, typ)
}

func actionField(oid, tid int64) string {
	return fmt.Sprintf(_actionField, oid, tid)
}

// ExpireAction .
func (d *Dao) ExpireAction(c context.Context, mid int64, typ int32) (ok bool, err error) {
	key := keyAction(mid, typ)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.actionExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s), error(%v)", key, err)
	}
	conn.Close()
	return
}

// AddActionCache .
func (d *Dao) AddActionCache(c context.Context, mid, oid, tid int64, typ, action int32) (err error) {
	key := keyAction(mid, typ)
	field := actionField(oid, tid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HSET", key, field, action); err != nil {
		log.Error("conn.Send(HSET %s,%s) error(%v)", key, field, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.actionExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AddActionsCache .
func (d *Dao) AddActionsCache(c context.Context, mid int64, typ int32, actions []*model.ResourceAction) (err error) {
	key := keyAction(mid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, m := range actions {
		if err = conn.Send("HSET", key, actionField(m.Oid, m.Tid), m.Action); err != nil {
			log.Error("conn.Send(HSET %s,%s) error(%v)", key, actionField(m.Oid, m.Tid), err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.actionExpire); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(actions)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// ActionCache .
func (d *Dao) ActionCache(c context.Context, mid, oid, tid int64, typ int32) (action int32, err error) {
	key := keyAction(mid, typ)
	field := actionField(oid, tid)
	conn := d.redis.Get(c)
	defer conn.Close()
	num, err := redis.Int64(conn.Do("HGET", key, field))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("conn.Do(HGET %s,%s) error(%v)", key, field, err)
		return
	}
	action = int32(num)
	return
}

// ActionsCache .
func (d *Dao) ActionsCache(c context.Context, mid, oid int64, typ int32, tids []int64) (res map[int64]int32, err error) {
	var (
		key  = keyAction(mid, typ)
		args = []interface{}{key}
	)
	for _, tid := range tids {
		args = append(args, actionField(oid, tid))
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Int64s(conn.Do("HMGET", args...))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("conn.Do(HMGET %v) error(%v)", args, err)
		return
	}
	res = make(map[int64]int32, len(tids))
	for i, tid := range tids {
		res[tid] = int32(values[i])
	}
	return
}

// DelActionCache .
func (d *Dao) DelActionCache(c context.Context, mid int64, typ int32) (err error) {
	key := keyAction(mid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}
