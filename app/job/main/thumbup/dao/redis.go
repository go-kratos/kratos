package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/thumbup/model"
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
		log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, d.redisStatsExpire, err)
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
		log.Error("conn.DO(HMSET, %s, %v) error(%v)", key, commonds, err)
	}
	return
}
