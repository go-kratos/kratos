package redis

import (
	"context"
	"fmt"
	"go-common/app/admin/main/aegis/model/monitor"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"strconv"
	"time"
)

// MoniRuleStats 获取监控统计
func (d *Dao) MoniRuleStats(c context.Context, id, min, max int64) (stats *monitor.Stats, err error) {
	var (
		conn = d.cluster.Get(c)
		key  = fmt.Sprintf(monitor.RedisPrefix, id)
		now  = time.Now().Unix()
	)
	stats = &monitor.Stats{}
	defer conn.Close()

	if stats.TotalCount, err = redis.Int(conn.Do("ZCOUNT", key, 0, now)); err != nil {
		log.Error("conn.Do(ZCOUNT,%s,0,%d) error(%v)", key, now, err)
		return
	}
	if stats.MoniCount, err = redis.Int(conn.Do("ZCOUNT", key, min, max)); err != nil {
		log.Error("conn.Do(ZCOUNT,%s,%d,%d) error(%v)", key, min, max, err)
		return
	}
	var oldest map[string]string //进入列表最久的项
	oldest, err = redis.StringMap(conn.Do("ZRANGE", key, 0, 0, "WITHSCORES"))
	for _, t := range oldest {
		var i int
		if i, err = strconv.Atoi(t); err != nil {
			return
		}
		stats.MaxTime = int(now) - i
	}
	return
}

// MoniRuleOids 获取监控符合报警时长的
func (d *Dao) MoniRuleOids(c context.Context, id, min, max int64) (oidMap map[int64]int, err error) {
	var (
		key    = fmt.Sprintf(monitor.RedisPrefix, id)
		conn   = d.cluster.Get(c)
		strMap map[string]int
	)
	oidMap = make(map[int64]int)
	strMap = make(map[string]int)
	defer conn.Close()
	if strMap, err = redis.IntMap(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES")); err != nil {
		log.Error("redis.IntMap(conn.Do(\"ZRANGEBYSCORE\", %s, %d, %d, \"WITHSCORES\")) error(%v)", key, min, max, err)
		return
	}
	for k, v := range strMap {
		oid := 0
		if oid, err = strconv.Atoi(k); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", k, err)
		}
		oidMap[int64(oid)] = v
	}
	return
}
