package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/tag/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_tagPridKey = "tp_%d_%d"
)

func tagPridKey(tid, prid int64) string {
	return fmt.Sprintf(_tagPridKey, tid, prid)
}

// AddTagPridArcCache .
func (d *Dao) AddTagPridArcCache(c context.Context, arc *model.Archive, prid int64, tids []int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var pubTime int64
	if pubTime, err = d.getPubTime(arc.Aid, arc.PubTime); err != nil {
		return
	}
	for _, tid := range tids {
		var ok bool
		if ok, err = d.expireCache(c, tid, prid); err != nil || !ok {
			return
		}
		key := tagPridKey(tid, prid)
		// del -1 cache
		d.RemTagPridArcCache(c, -1, prid, []int64{tid})
		if err = conn.Send("ZADD", key, pubTime, arc.Aid); err != nil {
			log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)", key, pubTime, arc.Aid, err)
			return
		}
		if err = conn.Send("EXPIRE", key, d.expNewArc); err != nil {
			log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
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
	}
	return
}

// RemTagPridArcCache .
func (d *Dao) RemTagPridArcCache(c context.Context, aid, prid int64, tids []int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
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

// expire
func (d *Dao) expireCache(c context.Context, tid, prid int64) (ok bool, err error) {
	key := tagPridKey(tid, prid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE,%s), error(%v)", key, err)
	}
	return
}
