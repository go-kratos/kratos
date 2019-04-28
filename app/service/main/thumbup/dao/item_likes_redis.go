package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func itemLikesKey(businessID, messageID int64, state int8) string {
	return fmt.Sprintf("i2_m_%d_b_%d_%d", messageID, businessID, state)
}

// CacheItemLikeList .
func (d *Dao) CacheItemLikeList(c context.Context, messageID, businessID int64, state int8, start, end int) (res []*model.UserLikeRecord, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	items, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "withscores"))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		PromError("redis:CacheItemLikeList")
		log.Errorv(c, log.KV("CacheItemLikeList", fmt.Sprintf("%+v", err)))
		return
	}
	for len(items) > 0 {
		var id, t int64
		if items, err = redis.Scan(items, &id, &t); err != nil {
			PromError("redis:CacheItemLikeList")
			log.Errorv(c, log.KV("CacheItemLikeList", fmt.Sprintf("%+v", err)))
			return
		}
		res = append(res, &model.UserLikeRecord{Mid: id, Time: xtime.Time(t)})
	}
	return
}

// AddCacheItemLikeList .
func (d *Dao) AddCacheItemLikeList(c context.Context, messageID int64, miss []*model.UserLikeRecord, businessID int64, state int8) (err error) {
	if len(miss) == 0 {
		return
	}
	limit := d.BusinessIDMap[businessID].MessageLikesLimit
	var count int
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	if err = conn.Send("DEL", key); err != nil {
		PromError("redis:项目点赞列表")
		log.Errorv(c, log.KV("AddCacheItemLikeList", fmt.Sprintf("AddCacheItemLikeList conn.Send(DEL, %s) error(%+v)", key, err)))
		return
	}
	count++
	for _, item := range miss {
		id := item.Mid
		score := int64(item.Time)
		if err = conn.Send("ZADD", key, "CH", score, id); err != nil {
			PromError("redis:项目点赞列表")
			log.Errorv(c, log.KV("log", fmt.Sprintf("AddCacheItemLikeList conn.Send(ZADD, %s, %d, %v) error(%v)", key, score, id, err)))
			return
		}
		count++
	}
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(limit + 1)); err != nil {
		PromError("redis:项目点赞列表rm")
		log.Errorv(c, log.KV("log", fmt.Sprintf("AddCacheItemLikeList conn.Send(ZREMRANGEBYRANK, %s, 0, %d) error(%v)", key, -(limit+1), err)))
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.redisItemLikesExpire); err != nil {
		PromError("redis:项目点赞列表过期")
		log.Errorv(c, log.KV("log", fmt.Sprintf("AddCacheItemLikeList conn.Send(EXPIRE, %s, %d) error(%v)", key, d.redisItemLikesExpire, err)))
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:项目点赞列表flush")
		log.Errorv(c, log.KV("log", fmt.Sprintf("AddCacheItemLikeList conn.Flush error(%v)", err)))
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:项目点赞列表receive")
			log.Errorv(c, log.KV("log", fmt.Sprintf("AddCacheItemLikeList conn.Receive error(%v)", err)), log.KV("miss", miss))
			return
		}
	}
	return
}

// ExpireItemLikesCache .
func (d *Dao) ExpireItemLikesCache(c context.Context, messageID, businessID int64, state int8) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisItemLikesExpire)); err != nil {
		PromError("redis:expire项目点赞列表")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(EXPIRE, %s) error(%v)", key, err)))
	}
	return
}

// ItemLikeExists .
func (d *Dao) ItemLikeExists(c context.Context, messageID, businessID int64, mids []int64, state int8) (res map[int64]int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	res = make(map[int64]int64)
	for _, name := range mids {
		conn.Send("ZSCORE", key, name)
	}
	if err = conn.Flush(); err != nil {
		PromError("redis:ItemLikeExists")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Flush() error(%v)", err)))
		return
	}
	for _, name := range mids {
		var s int64
		if s, err = redis.Int64(conn.Receive()); err == nil {
			res[name] = s
		} else if err == redis.ErrNil {
			err = nil
		} else {
			PromError("redis:ItemLikeExists")
			log.Errorv(c, log.KV("log", fmt.Sprintf("ItemLikeExists conn.Receive() error(%v)", err)))
			return
		}
	}
	return
}

// AppendCacheItemLikeList .
func (d *Dao) AppendCacheItemLikeList(c context.Context, messageID int64, item *model.UserLikeRecord, businessID int64, state int8) (err error) {
	if item == nil {
		return
	}
	limit := d.BusinessIDMap[businessID].MessageLikesLimit
	var count int
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	id := item.Mid
	score := int64(item.Time)
	if err = conn.Send("ZADD", key, "CH", score, id); err != nil {
		PromError("redis:项目点赞列表")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZADD, %s, %d, %v) error(%v)", key, score, id, err)))
		return
	}
	count++
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(limit + 1)); err != nil {
		PromError("redis:项目点赞列表rm")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZREMRANGEBYRANK, %s, 0, %d) error(%v)", key, -(limit+1), err)))
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.redisItemLikesExpire); err != nil {
		PromError("redis:项目点赞列表过期")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(EXPIRE, %s, %d) error(%v)", key, d.redisItemLikesExpire, err)))
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:项目点赞列表flush")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Flush error(%v)", err)))
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			PromError("redis:项目点赞列表receive")
			log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Receive error(%v)", err)))
			return
		}
	}
	return
}

// DelItemLikeCache .
func (d *Dao) DelItemLikeCache(c context.Context, messageID, businessID int64, mid int64, state int8) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	if _, err = conn.Do("ZREM", key, mid); err != nil {
		PromError("redis:zrem项目点赞列表")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZREM, %s, %v) error(%v)", key, mid, err)))
	}
	return
}

// ItemLikesCountCache .
func (d *Dao) ItemLikesCountCache(c context.Context, businessID, messageID int64) (res int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, model.StateLike)
	res, err = redis.Int(conn.Do("ZCOUNT", key, "(0", "+inf"))
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("dao.ItemLikesCountCache(%d, %d) err:%v", businessID, messageID, err)))
		PromError("redis:项目点赞总数")
	}
	return
}
