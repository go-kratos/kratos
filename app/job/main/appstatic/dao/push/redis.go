package push

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const _pushKey = "appstatic-admin-topush"

// ZrangeList picks up all the to push resIDs, ctime
func (d *Dao) ZrangeList(c context.Context) (resIDs map[string]int64, err error) {
	var (
		conn = d.redis.Get(c)
		key  = _pushKey
	)
	defer conn.Close()
	// get all the resIDs in one shot
	if err = conn.Send("ZRANGE", _pushKey, 0, -1, "WITHSCORES"); err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() err(%v)", err)
		return
	}
	if resIDs, err = redis.Int64Map(conn.Receive()); err != nil {
		log.Error("redis.Int64s()err(%v)", err)
		return
	}
	return
}

// ZRem ZREM trim from trim queue.
func (d *Dao) ZRem(c context.Context, resID string) (err error) {
	var (
		conn = d.redis.Get(c)
		key  = _pushKey
	)
	if _, err = conn.Do("ZREM", key, resID); err != nil {
		log.Error("conn.Send(ZADD %s - %v) error(%v)", key, resID, err)
	}
	conn.Close()
	return
}
