package dao

import (
	"context"
	"strconv"

	"go-common/app/job/main/dm2/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixRecent = "dm_rct_"
)

func keyRecent(mid int64) string {
	return _prefixRecent + strconv.FormatInt(mid, 10)
}

// AddRecentDM add recent dm of up to redis.
func (d *Dao) AddRecentDM(c context.Context, mid int64, dm *model.DM) (count int64, err error) {
	var (
		conn  = d.dmRctRds.Get(c)
		key   = keyRecent(mid)
		value []byte
	)
	defer conn.Close()
	if value, err = dm.Marshal(); err != nil {
		log.Error("dm.Marshal(%v) error(%v)", dm, err)
		return
	}
	if err = conn.Send("ZREMRANGEBYSCORE", key, dm.ID, dm.ID); err != nil {
		log.Error("conn.Do(ZREMRANGEBYSCORE %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("ZADD", key, dm.ID, value); err != nil {
		log.Error("conn.Send(ZADD %v) error(%v)", dm, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.dmRctExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("ZCARD", key); err != nil {
		log.Error("conn.Send(ZCARD %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 3; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	if count, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	return
}

// ZRemRecentDM remove recent dm of up.
func (d *Dao) ZRemRecentDM(c context.Context, mid, dmid int64) (err error) {
	var (
		conn = d.dmRctRds.Get(c)
		key  = keyRecent(mid)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREMRANGEBYSCORE", key, dmid, dmid); err != nil {
		log.Error("conn.Do(ZREMRANGEBYSCORE %s) error(%v)", key, dmid)
	}
	return
}

// TrimRecentDM zrange remove recent dm of up.
func (d *Dao) TrimRecentDM(c context.Context, mid, count int64) (err error) {
	var (
		conn = d.dmRctRds.Get(c)
		key  = keyRecent(mid)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREMRANGEBYRANK", key, 0, count-1); err != nil {
		log.Error("conn.Do(ZREMRANGEBYRANK %s) error(%v)", key, err)
	}
	return
}
