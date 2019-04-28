package dao

import (
	"context"
	"fmt"

	pb "go-common/app/service/main/thumbup/api"
	"go-common/app/service/main/thumbup/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/stat/prom"
	xtime "go-common/library/time"

	pkgerr "github.com/pkg/errors"
)

func userLikesKey(businessID, mid int64, state int8) string {
	return fmt.Sprintf("u_m_%d_b_%d_%d", mid, businessID, state)
}

// CacheUserLikeList .
func (d *Dao) CacheUserLikeList(c context.Context, mid, businessID int64, state int8, start, end int) (res []*model.ItemLikeRecord, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := userLikesKey(businessID, mid, state)
	items, err := redis.Values(conn.Do("ZREVRANGE", key, start, end, "withscores"))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		err = pkgerr.Wrap(err, "get user like cache")
		PromError("redis:CacheUserLikeList")
		log.Errorv(c, log.KV("CacheUserLikeList", fmt.Sprintf("%+v", err)))
		return
	}
	for len(items) > 0 {
		var id, t int64
		if items, err = redis.Scan(items, &id, &t); err != nil {
			err = pkgerr.Wrap(err, "get user like cache")
			PromError("redis:CacheUserLikeList")
			log.Errorv(c, log.KV("CacheUserLikeList", fmt.Sprintf("%+v", err)))
			return
		}
		res = append(res, &model.ItemLikeRecord{MessageID: id, Time: xtime.Time(t)})
	}
	return
}

// ExpireUserLikesCache .
func (d *Dao) ExpireUserLikesCache(c context.Context, mid, businessID int64, state int8) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := userLikesKey(businessID, mid, state)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisUserLikesExpire)); err != nil {
		err = pkgerr.Wrap(err, "")
		PromError("redis:expire用户点赞列表")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(EXPIRE, %s) error(%v)", key, err)))
	}
	return
}

// UserLikeExists .
func (d *Dao) UserLikeExists(c context.Context, mid, businessID int64, messageIDs []int64, state int8) (res map[int64]*pb.UserLikeState, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := userLikesKey(businessID, mid, state)
	res = make(map[int64]*pb.UserLikeState)
	prom.CacheHit.Incr("userLikeList")
	for _, name := range messageIDs {
		conn.Send("ZSCORE", key, name)
	}
	if err = conn.Flush(); err != nil {
		err = pkgerr.Wrap(err, "")
		PromError("redis:UserLikeExists")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Flush() error(%v)", err)))
		return
	}
	for _, name := range messageIDs {
		var ts int64
		if ts, err = redis.Int64(conn.Receive()); err == nil {
			res[name] = &pb.UserLikeState{
				Mid:   mid,
				Time:  xtime.Time(ts),
				State: pb.State(state),
			}
		} else if err == redis.ErrNil {
			err = nil
		} else {
			err = pkgerr.Wrap(err, "")
			PromError("redis:UserLikeExists")
			log.Errorv(c, log.KV("log", fmt.Sprintf("UserLikeExists conn.Receive() error(%v)", err)))
			return
		}
	}
	return
}

// AppendCacheUserLikeList .
func (d *Dao) AppendCacheUserLikeList(c context.Context, mid int64, item *model.ItemLikeRecord, businessID int64, state int8) (err error) {
	if item == nil {
		return
	}
	limit := d.BusinessIDMap[businessID].UserLikesLimit
	var count int
	conn := d.redis.Get(c)
	defer conn.Close()
	key := userLikesKey(businessID, mid, state)
	id := item.MessageID
	score := int64(item.Time)
	if err = conn.Send("ZADD", key, "CH", score, id); err != nil {
		err = pkgerr.Wrap(err, "")
		PromError("redis:用户点赞列表")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZADD, %s, %d, %v) error(%v)", key, score, id, err)))
		return
	}
	count++
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(limit + 1)); err != nil {
		err = pkgerr.Wrap(err, "")
		PromError("redis:用户点赞列表rm")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZREMRANGEBYRANK, %s, 0, %d) error(%v)", key, -(limit+1), err)))
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.redisUserLikesExpire); err != nil {
		err = pkgerr.Wrap(err, "")
		PromError("redis:用户点赞列表过期")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(EXPIRE, %s, %d) error(%v)", key, d.redisUserLikesExpire, err)))
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		PromError("redis:用户点赞列表flush")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Flush error(%v)", err)))
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			err = pkgerr.Wrap(err, "")
			PromError("redis:用户点赞列表receive")
			log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Receive error(%v)", err)))
			return
		}
	}
	return
}

// DelUserLikeCache .
func (d *Dao) DelUserLikeCache(c context.Context, mid, businessID int64, messageID int64, state int8) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := userLikesKey(businessID, mid, state)
	if _, err = conn.Do("ZREM", key, messageID); err != nil {
		err = pkgerr.Wrap(err, "")
		PromError("redis:zrem用户点赞列表")
		log.Errorv(c, log.KV("log", fmt.Sprintf("conn.Send(ZREM, %s, %v) error(%v)", key, messageID, err)))
	}
	return
}

// UserLikesCountCache .
func (d *Dao) UserLikesCountCache(c context.Context, businessID, mid int64) (res int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := userLikesKey(businessID, mid, model.StateLike)
	res, err = redis.Int(conn.Do("ZCOUNT", key, "(0", "+inf"))
	if err != nil {
		err = pkgerr.Wrap(err, "")
		log.Errorv(c, log.KV("log", fmt.Sprintf("dao.UserLikesCountCache(%d, %d) err:%v", businessID, mid, err)))
		PromError("redis:用户点赞总数")
	}
	return
}
