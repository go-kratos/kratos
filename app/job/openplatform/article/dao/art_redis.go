package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

func sortedKey(categoryID int64, field int) string {
	return fmt.Sprintf("art_sort_%d_%d", categoryID, field)
}

// ExpireSortCache expire sort cache
func (d *Dao) ExpireSortCache(c context.Context, categoryID int64, field int) (ok bool, err error) {
	key := sortedKey(categoryID, field)
	conn := d.artRedis.Get(c)
	defer conn.Close()
	var ttl int64
	if ttl, err = redis.Int64(conn.Do("TTL", key)); err != nil {
		PromError("redis:排序缓存ttl")
		log.Error("conn.Do(TTL, %s) error(%+v)", key, err)
	}
	if ttl > (d.redisSortTTL - d.redisSortExpire) {
		ok = true
	}
	return
}

// AddSortCaches add sort articles cache
func (d *Dao) AddSortCaches(c context.Context, categoryID int64, field int, arts [][2]int64, maxLength int64) (err error) {
	var (
		id, score int64
		key       = sortedKey(categoryID, field)
		conn      = d.artRedis.Get(c)
		count     int
	)
	defer conn.Close()
	if len(arts) == 0 {
		return
	}
	if err = conn.Send("DEL", key); err != nil {
		PromError("redis:删除排序缓存")
		log.Error("conn.Send(DEL, %s) error(%+v)", key, err)
		return
	}
	count++
	for _, art := range arts {
		id = art[0]
		score = art[1]
		if err = conn.Send("ZADD", key, "CH", score, id); err != nil {
			PromError("redis:增加排序缓存")
			log.Error("conn.Send(ZADD, %s, %d, %v) error(%+v)", key, score, id, err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.redisSortTTL); err != nil {
		PromError("redis:排序缓存设定过期")
		log.Error("conn.Send(EXPIRE, %s, %d) error(%+v)", key, d.redisSortTTL, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:增加排序缓存flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:增加排序缓存receive")
			log.Error("conn.Receive error(%+v)", err)
			return
		}
	}
	return
}
