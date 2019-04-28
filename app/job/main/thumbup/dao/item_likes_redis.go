package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/thumbup/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

func itemLikesKey(businessID, messageID int64, state int8) string {
	return fmt.Sprintf("i2_m_%d_b_%d_%d", messageID, businessID, state)
}

// ExpireItemLikesCache .
func (d *Dao) ExpireItemLikesCache(c context.Context, messageID, businessID int64, state int8) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisItemLikesExpire)); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(EXPIRE, %s) error(%v)", key, err)))
	}
	return
}

// AddItemLikesCache .
func (d *Dao) AddItemLikesCache(c context.Context, businessID, messageID int64, typ int8, limit int, items []*model.UserLikeRecord) (err error) {
	if len(items) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, typ)
	if err = conn.Send("DEL", key); err != nil {
		log.Errorv(c, log.KV("AddCacheItemLikeList", fmt.Sprintf("AddCacheItemLikeList conn.Send(DEL, %s) error(%+v)", key, err)))
		return
	}
	args := redis.Args{}.Add(key).Add("CH")
	for _, item := range items {
		args = args.Add(int64(item.Time)).Add(item.Mid)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("zadd key(%s) args(%v) error(%v)", key, args, err)
		return
	}
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(limit + 1)); err != nil {
		log.Error("zremrangebyrank error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisItemLikesExpire); err != nil {
		log.Error("expire key(%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("redis flush error(%v)", err)
		return
	}
	for i := 0; i < 4; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("redis receive error(%v)", err)
			return
		}
	}
	return
}

// AppendCacheItemLikeList .
func (d *Dao) AppendCacheItemLikeList(c context.Context, messageID int64, item *model.UserLikeRecord, businessID int64, state int8, limit int) (err error) {
	if item == nil {
		return
	}
	var count int
	conn := d.redis.Get(c)
	defer conn.Close()
	key := itemLikesKey(businessID, messageID, state)
	id := item.Mid
	score := int64(item.Time)
	if err = conn.Send("ZADD", key, "CH", score, id); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZADD, %s, %d, %v) error(%v)", key, score, id, err)))
		return
	}
	count++
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(limit + 1)); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZREMRANGEBYRANK, %s, 0, %d) error(%v)", key, -(limit+1), err)))
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.redisItemLikesExpire); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(EXPIRE, %s, %d) error(%v)", key, d.redisItemLikesExpire, err)))
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Flush error(%v)", err)))
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
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
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZREM, %s, %v) error(%v)", key, mid, err)))
	}
	return
}
