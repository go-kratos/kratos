package ads

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/dgryski/go-farm"
)

const (
	_prefixBuvid = "buvid:%d"
)

func (d *Dao) keyBuvid(buvid string) (key string) {
	num := int64(farm.Hash32([]byte(buvid)))
	key = fmt.Sprintf(_prefixBuvid, num%d.c.HashNum)
	return
}

// ExistsAuth if existes buvid in redis.
func (d *Dao) ExistsAuth(c context.Context, key string) (ok bool, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		log.Error("EXISTS key(%s), error(%v)", key, err)
	}
	return
}

// BuvidCount get buvid count info from redis.
func (d *Dao) BuvidCount(c context.Context, faid int64, buvid string) (res int64, err error) {
	var (
		key   = d.keyBuvid(buvid)
		conn  = d.redis.Get(c)
		field = faid
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("HGET", key, field)); err != nil {
		if err != redis.ErrNil {
			log.Error("BuvidCount conn.Send HGET(%v, %v) error(%v)", key, field, err)
			return
		}
		err = nil
	}
	return
}

// AddBuvidCount add buvid count info into redis.
func (d *Dao) AddBuvidCount(c context.Context, buvidCounts map[string]map[int64]int64) (err error) {
	var (
		key       string
		count     int
		conn      = d.redis.Get(c)
		faid      int64
		playCount int64
		ok        bool
	)
	defer conn.Close()
	for buvid, buvidCount := range buvidCounts {
		key = d.keyBuvid(buvid)
		if ok, err = d.ExistsAuth(c, key); err != nil {
			log.Error("EXISTS key(%s) error(%v)", key, err)
			return
		}
		for faid, playCount = range buvidCount {
			if err = conn.Send("HSET", key, faid, playCount); err != nil {
				log.Error("HSET key(%s) field(%d) playCount(%d) error(%v)", key, faid, playCount, err)
				return
			}
			count++
		}
		if !ok {
			if err = conn.Send("EXPIRE", key, d.expire); err != nil {
				log.Error("EXPIRE key(%s) expire(%v) error(%v)", key, d.expire, err)
				return
			}
			count++
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("BuvidCount conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("BuvidCount conn.Receive error(%v)", err)
			return
		}
	}
	return
}
