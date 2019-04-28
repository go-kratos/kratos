package share

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const _keySha = "sm_%d"

func keyShare(mid int64) string {
	return fmt.Sprintf(_keySha, mid)
}

// SharesCache get shares from  cache.
func (d *Dao) SharesCache(c context.Context, mid int64) (res map[string]int64, err error) {
	var (
		key  = keyShare(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Int64Map(conn.Do("HGETALL", key)); err != nil {
		log.Error("HGETALL %v failed error(%v)", key, err)
	}
	return
}

// SetSharesCache back cache from db.
func (d *Dao) SetSharesCache(c context.Context, expire int, mid int64, shares map[string]int64) (err error) {
	var (
		key  = keyShare(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for k, v := range shares {
		args = args.Add(k).Add(v)
	}
	if err = conn.Send("HMSET", args...); err != nil {
		log.Error("conn.Send(HMSET, %s, %d) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
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
