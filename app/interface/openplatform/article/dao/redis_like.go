package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

func maxLikeKey(aid int64) string {
	return fmt.Sprintf(_prefixMaxLike, aid)
}

// MaxLikeCache max like cache
func (d *Dao) MaxLikeCache(c context.Context, aid int64) (res int64, err error) {
	var (
		conn = d.redis.Get(c)
		key  = maxLikeKey(aid)
	)
	defer conn.Close()
	if res, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			PromError("redis:获取点赞最大数")
			log.Error("MaxLikeMCache GET(%s) error(%+v)", key, err)
		}
		return
	}
	return
}

// ExpireMaxLikeCache expire max like cache
func (d *Dao) ExpireMaxLikeCache(c context.Context, aid int64) (res bool, err error) {
	var (
		conn = d.redis.Get(c)
		key  = maxLikeKey(aid)
	)
	defer conn.Close()
	if res, err = redis.Bool(conn.Do("EXPIRE", key, d.redisMaxLikeExpire)); err != nil {
		PromError("redis:Expire点赞最大数")
		log.Error("MaxLikeCache EXPIRE(%s) error(%+v)", key, err)
	}
	return
}

// SetMaxLikeCache set max like cache
func (d *Dao) SetMaxLikeCache(c context.Context, aid int64, value int64) (err error) {
	var (
		conn  = d.redis.Get(c)
		key   = maxLikeKey(aid)
		count int
	)
	defer conn.Close()
	if err = conn.Send("SET", key, value); err != nil {
		PromError("redis:设定点赞最大数")
		log.Error("conn.Send(SET, %s, %s) error(%+v)", key, value, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.redisMaxLikeExpire); err != nil {
		PromError("redis:Expire点赞最大数")
		log.Error("MaxLikeCache EXPIRE(%s) error(%+v)", key, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:设定点赞最大数flush")
		log.Error("conn.Flush error(%+v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:设定点赞最大数receive")
			log.Error("conn.Receive error(%+v)", err)
			return
		}
	}
	return
}
