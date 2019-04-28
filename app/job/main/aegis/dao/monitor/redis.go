package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/job/main/aegis/model/monitor"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"strconv"
	"time"
)

const (
	// _maxAge Sorted
	_maxAge = 604800 //7天
)

// AddToSet add monitor stats
func (d *Dao) AddToSet(c context.Context, keys []string, oid int64) (logs []string, err error) {
	if len(keys) == 0 {
		return
	}
	var (
		conn = d.redis.Get(c)
		now  = time.Now().Unix()
	)
	defer conn.Close()
	for _, key := range keys {
		//先判断key是否存在，存在则忽略
		if v, _ := redis.Int(conn.Do("ZSCORE", key, oid)); v != 0 {
			logs = append(logs, fmt.Sprintf("AddToSet() conn.Do(ZSCORE, %s, %d) member exists success", key, oid))
			continue
		}
		if _, err = conn.Do("ZADD", key, now, oid); err != nil {
			log.Error("conn.Do(ZADD, %s, %d, %d) error(%v)", key, now, oid, err)
			logs = append(logs, fmt.Sprintf("AddToSet() conn.Do(ZADD, %s, %d, %d) error(%v)", key, now, oid, err))
		} else {
			logs = append(logs, fmt.Sprintf("AddToSet() conn.Do(ZADD, %s, %d, %d) success", key, now, oid))
		}
		if _, err = conn.Do("EXPIRE", key, _maxAge); err != nil {
			log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, _maxAge, err)
			logs = append(logs, fmt.Sprintf("AddToSet() conn.Do(EXPIRE, %s, %d) error(%v)", key, _maxAge, err))
		} else {
			logs = append(logs, fmt.Sprintf("AddToSet() conn.Do(EXPIRE, %s, %d) success", key, _maxAge))
		}
	}
	return
}

// RemFromSet remove monitor stats
func (d *Dao) RemFromSet(c context.Context, keys []string, oid int64) (logs []string, err error) {
	if len(keys) == 0 {
		return
	}
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, key := range keys {
		if _, er := conn.Do("ZREM", key, oid); er != nil {
			err = er
			log.Error("conn.Do(ZREM, %s, %d) error(%v)", key, oid, err)
			logs = append(logs, fmt.Sprintf("RemFromSet() conn.Do(ZREM, %s, %d) error(%v)", key, oid, err))
			continue
		}
		logs = append(logs, fmt.Sprintf("RemFromSet() conn.Do(ZREM, %s, %d) success", key, oid))
	}
	return
}

// ClearExpireSet clear expire stats
func (d *Dao) ClearExpireSet(c context.Context, keys []string) (logs []string, err error) {
	if len(keys) == 0 {
		return
	}
	var (
		conn = d.redis.Get(c)
		now  = time.Now().Unix()
		min  int64
		max  = now - _maxAge
	)
	defer conn.Close()
	for _, key := range keys {
		if _, er := conn.Do("ZREMRANGEBYSCORE", key, min, max); er != nil {
			err = er
			log.Error("conn.Do(ZREMRANGEBYSCORE, %s, %d, %d) error(%v)", key, min, max, err)
			logs = append(logs, fmt.Sprintf("ClearExpireSet() key: %s min:%d max:%d error:%v", key, min, max, err))
			continue
		}
		logs = append(logs, fmt.Sprintf("ClearExpireSet() key: %s min:%d max:%d success", key, min, max))
	}
	return
}

// AddToDelArc 添加稿件信息到
func (d *Dao) AddToDelArc(c context.Context, a *monitor.BinlogArchive) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	info := &monitor.DelArcInfo{
		AID:   a.ID,
		MID:   a.MID,
		Time:  a.MTime,
		Title: a.Title,
	}
	if bs, err = json.Marshal(info); err != nil {
		log.Error("json.Marshal(%+v) error:%v", info, err)
		return
	}
	if _, err = conn.Do("HSET", monitor.RedisDelArcInfo, a.ID, string(bs)); err != nil {
		log.Error("conn.Send(HSET,%s,%d,%s) error(%v)", monitor.RedisDelArcInfo, a.ID, bs, err)
		return
	}
	return
}

// ArcDelInfos 获取被删除稿件的信息
func (d *Dao) ArcDelInfos(c context.Context, aids []int64) (infos map[int64]*monitor.DelArcInfo, err error) {
	var (
		conn = d.redis.Get(c)
		strs []string
	)
	defer conn.Close()
	infos = make(map[int64]*monitor.DelArcInfo)
	if len(aids) == 0 {
		return
	}
	args := redis.Args{}
	args = args.Add(monitor.RedisDelArcInfo)
	for _, id := range aids {
		args = args.Add(id)
	}
	log.Info("s.monitorNotify() ArcDelInfos. aids(%v) args(%+v)", aids, args)
	if strs, err = redis.Strings(conn.Do("HMGET", args...)); err != nil {
		log.Error("conn.Send(HMGET,%v) error(%v)", args, err)
		return
	}
	log.Info("s.monitorNotify() ArcDelInfos. aids(%v) strs(%v)", aids, strs)
	for _, v := range strs {
		info := &monitor.DelArcInfo{}
		if err = json.Unmarshal([]byte(v), info); err != nil {
			log.Error("json.Unmarshal(%s) error:%v", v, err)
			continue
		}
		infos[info.AID] = info
	}
	return
}

// MoniRuleStats 获取监控统计
func (d *Dao) MoniRuleStats(c context.Context, id int64, min, max int64) (stats *monitor.Stats, err error) {
	var (
		conn = d.redis.Get(c)
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

// MoniRuleOids 获取监控的id
func (d *Dao) MoniRuleOids(c context.Context, id int64, min, max int64) (oidMap map[int64]int, err error) {
	var (
		conn   = d.redis.Get(c)
		key    = fmt.Sprintf(monitor.RedisPrefix, id)
		intMap map[string]int
	)
	oidMap = make(map[int64]int)
	intMap = make(map[string]int)
	defer conn.Close()
	if intMap, err = redis.IntMap(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES")); err != nil {
		log.Error("redis.IntMap(conn.Do(\"ZRANGEBYSCORE\", %s, %d, %d, \"WITHSCORES\")) error(%v)", key, min, max, err)
		return
	}
	for k, v := range intMap {
		oid := 0
		if oid, err = strconv.Atoi(k); err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", k, err)
		}
		oidMap[int64(oid)] = v
	}
	return
}
