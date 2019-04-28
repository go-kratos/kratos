package dao

import (
	"context"
	"math/rand"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

var _mainCategory = int64(0)

// AddSortCache add sort articles cache
func (d *Dao) AddSortCache(c context.Context, categoryID int64, field int, aid, score int64) (err error) {
	var (
		key  = sortedKey(categoryID, field)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("ZADD", key, "CH", score, aid); err != nil {
		PromError("redis:增加排序缓存")
		log.Error("conn.Do(ZADD, %s, %d, %v) error(%+v)", key, score, aid, err)
	}
	return
}

// 避免同时回源
func (d *Dao) randomSortTTL() int64 {
	random := rand.Int63() % (d.redisSortTTL / 20)
	if rand.Int()%2 == 0 {
		return d.redisSortTTL - random
	}
	return d.redisSortTTL + random
}

// SortCache get sort cache
func (d *Dao) SortCache(c context.Context, categoryID int64, field int, start, end int) (res []int64, err error) {
	key := sortedKey(categoryID, field)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREVRANGE", key, start, end); err != nil {
		PromError("redis:获取排序列表")
		log.Error("conn.Send(%s) error(%+v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:获取排序列表flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	res, err = redis.Int64s(conn.Receive())
	if err != nil {
		PromError("redis:获取排序列表receive")
		log.Error("conn.Send(ZREVRANGE, %s) error(%+v)", key, err)
	}
	return
}

// SortCacheByValue get new articles cache by aid
func (d *Dao) SortCacheByValue(c context.Context, categoryID int64, field int, value, score int64, ps int) (res []int64, err error) {
	var (
		index  int
		tmpRes []int64
		conn   = d.redis.Get(c)
		key    = sortedKey(categoryID, field)
	)
	defer conn.Close()
	if tmpRes, err = redis.Int64s(conn.Do("ZREVRANGEBYSCORE", key, score, "-inf", "LIMIT", 0, ps+1)); err != nil {
		PromError("redis:获取最新投稿列表")
		log.Error("redis(ZREVRANGEBYSCORE %s,%d,%d) error(%+v)", key, score, ps, err)
		return
	}
	for i, v := range tmpRes {
		if v == value {
			index = i + 1
			break
		}
	}
	res = tmpRes[index:]
	return
}

// ExpireSortCache expire sort cache
func (d *Dao) ExpireSortCache(c context.Context, categoryID int64, field int) (ok bool, err error) {
	key := sortedKey(categoryID, field)
	conn := d.redis.Get(c)
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

// DelSortCache delete sort cache
func (d *Dao) DelSortCache(c context.Context, categoryID int64, field int, aid int64) (err error) {
	key := sortedKey(categoryID, field)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, aid); err != nil {
		PromError("redis:删除排序")
		log.Error("conn.Do(ZERM, %s, %d) error(%+v)", key, aid, err)
	}
	return
}

// NewArticleCount get new article count
func (d *Dao) NewArticleCount(c context.Context, ptime int64) (res int64, err error) {
	var (
		key  = sortedKey(_mainCategory, model.FieldNew)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	begin := "(" + strconv.FormatInt(ptime, 10)
	if res, err = redis.Int64(conn.Do("ZCOUNT", key, begin, "+inf")); err != nil {
		PromError("redis:排序缓存计数")
		log.Error("conn.Do(ZCOUNT, %s, %s, +inf) error(%+v)", key, begin, err)
	}
	return
}
