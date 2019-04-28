package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// keyZone gets the key of the zone in Redis
func keyZone(category int) string {
	return fmt.Sprintf("zone_idx_%d", category)
}

// ZrevrangeList picks up the page of ids .
func (d *Dao) ZrevrangeList(c context.Context, category int, start, end int) (sids []int64, count int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyZone(category)
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
	if sids, err = redis.Int64s(conn.Receive()); err != nil {
		log.Error("redis.Int64s()err(%v)", err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		log.Error("redis.INT64 err(%v)", err)
	}
	return
}
