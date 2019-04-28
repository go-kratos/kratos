package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

// AddcategoriesAuthors .
func (d *Dao) AddcategoriesAuthors(c context.Context, data map[int64][]int64) (err error) {
	if len(data) == 0 {
		return
	}
	conn := d.artRedis.Get(c)
	defer conn.Close()
	for k, x := range data {
		key := fmt.Sprintf("recommends:authors:%d", k)
		args := redis.Args{}.Add(key)
		for _, v := range x {
			args = args.Add(fmt.Sprintf("%d", v))
		}
		if err = conn.Send("DEL", key); err != nil {
			log.Error("conn.Send(SADD, %s, %+v) error(%v)", key, args, err)
			PromError("redis:分区作者数据")
			return
		}
		if err = conn.Send("SADD", args...); err != nil {
			log.Error("conn.Send(SADD, %s, %+v) error(%v)", key, args, err)
			PromError("redis:分区作者数据")
			return
		}
		if err = conn.Send("EXPIRE", key, d.c.Job.RecommendExpire); err != nil {
			log.Error("conn.Send(EXPIRE, %s) error(%v)", d.c.Job.RecommendExpire, err)
			PromError("redis:分区作者数据")
			return
		}
		if err = conn.Flush(); err != nil {
			log.Error("conn.Flush error(%v)", err)
			PromError("redis:分区作者数据")
			return
		}
		for i := 0; i < 3; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("conn.Receive error(%v)", err)
				PromError("redis:分区作者数据")
				return
			}
		}
	}
	return
}

// AddAuthorMostCategories .
func (d *Dao) AddAuthorMostCategories(c context.Context, mid int64, categories []int64) (err error) {
	if len(categories) == 0 {
		return
	}
	conn := d.artRedis.Get(c)
	defer conn.Close()
	key := fmt.Sprintf("author:categories:%d", mid)
	args := redis.Args{}.Add(key)
	for _, v := range categories {
		args = args.Add(fmt.Sprintf("%d", v))
	}
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(SADD, %s, %+v) error(%v)", key, args, err)
		PromError("redis:作者分区数据")
		return
	}
	if err = conn.Send("SADD", args...); err != nil {
		log.Error("conn.Send(SADD, %s, %+v) error(%v)", key, args, err)
		PromError("redis:作者分区数据")
		return
	}
	if err = conn.Send("EXPIRE", key, d.c.Job.RecommendExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", d.c.Job.RecommendExpire, err)
		PromError("redis:作者分区数据")
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		PromError("redis:作者分区数据")
		return
	}
	for i := 0; i < 3; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			PromError("redis:作者分区数据")
			return
		}
	}
	return
}
