package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/service/main/archive/api"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	// new list
	_keyNl = "n_%d_%d"
)

func keyNl(rid int32, tp int8) string {
	return fmt.Sprintf(_keyNl, rid, tp)
}

func keyNlBak(rid int32, tp int8) string {
	return _keyBakPrefix + keyNl(rid, tp)
}

// NewListCache get region rank list from cache.
func (d *Dao) NewListCache(c context.Context, rid int32, tp int8, start, end int) (arcs []*api.Arc, count int, err error) {
	key := keyNl(rid, tp)
	conn := d.redis.Get(c)
	defer conn.Close()
	arcs, count, err = d.nlCache(c, conn, key, start, end)
	return
}

// NewListBakCache get region rank list from bak cache.
func (d *Dao) NewListBakCache(c context.Context, rid int32, tp int8, start, end int) (arcs []*api.Arc, count int, err error) {
	d.cacheProm.Incr("newlist_remote_cache")
	key := keyNlBak(rid, tp)
	conn := d.redisBak.Get(c)
	defer conn.Close()
	arcs, count, err = d.nlCache(c, conn, key, start, end)
	if len(arcs) == 0 {
		log.Error("NewlistBakCache(%s,%d,%d) is nil", key, start, end)
	}
	return
}

func (d *Dao) nlCache(c context.Context, conn redis.Conn, key string, start, end int) (arcs []*api.Arc, count int, err error) {
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		arc := &api.Arc{}
		if err = json.Unmarshal(bs, arc); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		arcs = append(arcs, arc)
	}
	count = from(num)
	return
}

// SetNewListCache set region cache.
func (d *Dao) SetNewListCache(c context.Context, rid int32, tp int8, arcs []*api.Arc, count int) (err error) {
	key := keyNl(rid, tp)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = d.setNlCache(c, conn, key, d.redisNlExpire, arcs, count); err != nil {
		return
	}
	key = keyNlBak(rid, tp)
	connBak := d.redisBak.Get(c)
	err = d.setNlCache(c, connBak, key, d.redisNlBakExpire, arcs, count)
	connBak.Close()
	return
}

func (d *Dao) setNlCache(c context.Context, conn redis.Conn, key string, expire int32, arcs []*api.Arc, num int) (err error) {
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	for _, arc := range arcs {
		bs, _ := json.Marshal(arc)
		if err = conn.Send("ZADD", key, combine(arc.PubDate, num), bs); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, string(bs), err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, expire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, expire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

func from(i int64) int {
	return int(i & 0xffffff)
}

func combine(pubdate time.Time, count int) int64 {
	return pubdate.Time().Unix()<<24 | int64(count)
}
