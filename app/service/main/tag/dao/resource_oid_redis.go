package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_oidKey = "ot_%d_%d" // sortedset  key:ot_oid_type value:tid score:ctime
)

func oidKey(oid int64, typ int32) string {
	return fmt.Sprintf(_oidKey, oid, typ)
}

// ExpireOidCache .
func (d *Dao) ExpireOidCache(c context.Context, oid int64, typ int32) (ok bool, err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("conn.Do(Expire %s) error(%v)", key, err)
	}
	return
}

// OidCache .
func (d *Dao) OidCache(c context.Context, oid int64, typ int32) (tids []int64, err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	if tids, err = redis.Int64s(conn.Do("ZRANGE", key, 0, -1)); err != nil {
		log.Error("redis.Int64s()err(%v)", err)
	}
	return
}

// AddOidCache .
func (d *Dao) AddOidCache(c context.Context, oid int64, typ int32, r *model.Resource) (err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, r.CTime, r.Tid); err != nil {
		log.Error("conn.Send(ZADD %v) error(%v)", r, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddOidsCache .
func (d *Dao) AddOidsCache(c context.Context, oid int64, typ int32, rs []*model.Resource) (err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, r := range rs {
		if err = conn.Send("ZADD", key, r.CTime, r.Tid); err != nil {
			log.Error("conn.Send(ZADD %v) error(%v)", r, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(rs)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddOidMapCache .
func (d *Dao) AddOidMapCache(c context.Context, oid int64, typ int32, rsMap map[int64]*model.Resource) (err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, r := range rsMap {
		if err = conn.Send("ZADD", key, r.CTime, r.Tid); err != nil {
			log.Error("conn.Send(ZADD %v) error(%v)", r, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(rsMap)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ZremOidCache .
func (d *Dao) ZremOidCache(c context.Context, oid, tid int64, typ int32) (err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, tid); err != nil {
		log.Error("conn.Do(ZREM %s,%d) error(%v)", key, tid, err)
	}
	return
}

// DelOidCache .
func (d *Dao) DelOidCache(c context.Context, oid int64, typ int32) (err error) {
	key := oidKey(oid, typ)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("del", key); err != nil {
		log.Error("conn.Do(ZREM %s) error(%v)", key, err)
	}
	return
}
