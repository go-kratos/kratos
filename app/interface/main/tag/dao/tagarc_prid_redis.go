package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_tagPridKey = "tp_%d_%d"
)

func tagPridKey(tid, prid int64) string {
	return fmt.Sprintf(_tagPridKey, tid, prid)
}

// ZrangeTagPridArc .
func (d *Dao) ZrangeTagPridArc(c context.Context, tid, prid int64, start, end int) (aids []int64, count int, err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	key := tagPridKey(tid, prid)
	if err = conn.Send("ZREVRANGE", key, start, end); err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("ZCARD", key); err != nil {
		log.Error("conn.Do(ZCARD) err(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() err(%v)", err)
		return
	}
	if aids, err = redis.Int64s(conn.Receive()); err != nil {
		log.Error("redis.Int64s()err(%v)", err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		log.Error("redis.INT64 err(%v)", err)
	}
	return
}

// AddTagPridArcCache .
func (d *Dao) AddTagPridArcCache(c context.Context, tids []int64, prid int64, as ...*api.Arc) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		var ok bool
		if ok, err = d.expireTagArcCache(c, tid, prid); err != nil {
			return
		}
		if !ok {
			return
		}
		key := tagPridKey(tid, prid)
		// del -1 cache
		d.RemoveTagPridArcCache(c, []int64{tid}, prid, -1)
		for _, a := range as {
			if err = conn.Send("ZADD", key, a.PubDate, a.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, a.Aid, err)
				return
			}
		}
		if err = conn.Send("EXPIRE", key, d.expireNewArc); err != nil {
			log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		}
		if err = conn.Flush(); err != nil {
			log.Error("conn.Flush error(%v)", err)
			return
		}
		for i := 0; i < len(as)+1; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("conn.Receive() error(%v)", err)
				return
			}
		}
	}
	return
}

// AddTagPridCache .
func (d *Dao) AddTagPridCache(c context.Context, tids []int64, prid int64, as ...*api.Arc) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		key := tagPridKey(tid, prid)
		for _, a := range as {
			if err = conn.Send("ZADD", key, a.PubDate, a.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, a.Aid, err)
				return
			}
		}
		if err = conn.Send("EXPIRE", key, d.expireNewArc); err != nil {
			log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
		}
		if err = conn.Flush(); err != nil {
			log.Error("conn.Flush error(%v)", err)
			return
		}
		for i := 0; i < len(as)+1; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("conn.Receive() error(%v)", err)
				return
			}
		}
	}
	return
}

// RemoveTagPridArcCache .
func (d *Dao) RemoveTagPridArcCache(c context.Context, tids []int64, prid, aid int64) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	var count int
	for _, tid := range tids {
		key := tagPridKey(tid, prid)
		if err = conn.Send("ZREM", key, aid); err != nil {
			log.Error("conn.Do(ZREM,%s,%d)", key, aid)
			return
		}
		count++
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

func (d *Dao) expireTagArcCache(c context.Context, tid, prid int64) (ok bool, err error) {
	key := tagPridKey(tid, prid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE,%s), error(%v)", key, err)
	}
	return
}
