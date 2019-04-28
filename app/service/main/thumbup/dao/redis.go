package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/xstr"
)

func hashStatsKey(businessID, originID int64) string {
	return fmt.Sprintf("stats_o_%d_b_%d", originID, businessID)
}

// ExpireHashStatsCache .
func (d *Dao) ExpireHashStatsCache(c context.Context, businessID, originID int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := hashStatsKey(businessID, originID)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisStatsExpire)); err != nil {
		PromError("redis:计数缓存设定过期")
		log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, d.redisStatsExpire, err)
	}
	return
}

// DelHashStatsCache del hash cache
func (d *Dao) DelHashStatsCache(c context.Context, businessID, originID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := hashStatsKey(businessID, originID)
	if _, err = conn.Do("del", key); err != nil {
		PromError("redis:计数缓存删除")
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	return
}

// HashStatsCache .
func (d *Dao) HashStatsCache(c context.Context, businessID, originID int64, messageIDs []int64) (res map[int64]*model.Stats, err error) {
	if len(messageIDs) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	key := hashStatsKey(businessID, originID)
	var ss []string
	var commonds []interface{}
	commonds = append(commonds, key)
	for _, m := range messageIDs {
		commonds = append(commonds, m)
	}
	if ss, err = redis.Strings(conn.Do("HMGET", commonds...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(HMGET, %s, %v) error(%v)", key, messageIDs, err)
			PromError("redis:获取统计信息")
		}
		return
	}
	res = make(map[int64]*model.Stats)
	for i, id := range messageIDs {
		if ss[i] == "" {
			continue
		}
		stat := &model.Stats{ID: id, OriginID: originID}
		num, _ := xstr.SplitInts(ss[i])
		if len(num) > 1 {
			stat.Likes = num[0]
			stat.Dislikes = num[1]
		}
		res[id] = stat
	}
	return
}

// AddHashStatsCache .
func (d *Dao) AddHashStatsCache(c context.Context, businessID, originID int64, stats ...*model.Stats) (err error) {
	if len(stats) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	key := hashStatsKey(businessID, originID)
	var commonds = []interface{}{key}
	for _, stat := range stats {
		commonds = append(commonds, stat.ID, xstr.JoinInts([]int64{stat.Likes, stat.Dislikes}))
	}
	if _, err = conn.Do("HMSET", commonds...); err != nil {
		PromError("redis:增加统计信息")
		log.Error("conn.DO(HMSET, %s, %v) error(%v)", key, commonds, err)
	}
	return
}

// AddHashStatsCacheMap .
func (d *Dao) AddHashStatsCacheMap(c context.Context, businessID, originID int64, stats map[int64]*model.Stats) (err error) {
	var s []*model.Stats
	for _, v := range stats {
		s = append(s, v)
	}
	return d.AddHashStatsCache(c, businessID, originID, s...)
}
