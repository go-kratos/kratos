package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/tag/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixSpamDaily = "sd_%d_%d" //  hash  key:mid_now  field:action  value:number
)

func spamKey(mid int64, now time.Time) string {
	return fmt.Sprintf(_prefixSpamDaily, mid, now.Day())
}

// IncrSpamCache increse count of spam cache.
func (d *Dao) IncrSpamCache(c context.Context, mid int64, action int32) (err error) {
	key := spamKey(mid, time.Now())
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HINCRBY", key, action, model.SpamIncrValue); err != nil {
		log.Error("d.IncrSpamCache().HINCRBY(%s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.rdsExpOp); err != nil {
		log.Error("d.IncrSpamCache().EXPIRE(%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("d.IncrSpamCache().Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("d.IncrSpamCache().Receive() error(%v)", err)
			return
		}
	}
	return
}

// SpamCache return spam count cache.
func (d *Dao) SpamCache(c context.Context, mid int64, action int32) (count int, err error) {
	key := spamKey(mid, time.Now())
	conn := d.redis.Get(c)
	if count, err = redis.Int(conn.Do("HGET", key, action)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Int(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}
