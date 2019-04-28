package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-common/app/admin/main/videoup/model/monitor"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"sort"
	"strconv"
	"time"
)

const (
	FieldKeyFormat = "%d_%d_%d" //监控规则配置的Redis key中的field格式
)

// StatsResult 获取稿件停留统计
func (d *Dao) StatsResult(c context.Context, key string, conf *monitor.RuleConf) (res *monitor.Stats, err error) {
	var (
		conn                = d.redis.Get(c)
		totalC, moniC, maxT int
		now                 = time.Now().Unix()
		tFrom, tTo          int64
		timeCdt             int64
		compCdt             string
		ok                  bool
	)
	defer conn.Close()
	if _, ok = conf.NotifyCdt["time"]; !ok {
		err = errors.New("配置的 NotifyCdt 中不存在 time")
		return
	}
	timeCdt = conf.NotifyCdt["time"].Value
	compCdt = conf.NotifyCdt["time"].Comp
	switch compCdt {
	case monitor.CompGT:
		tFrom = 0
		tTo = now - timeCdt
	case monitor.CompLT:
		tFrom = now - timeCdt
		tTo = now
	default:
		err = errors.New("配置的 NotifyCdt 中 comparison 不合法: " + compCdt)
		return
	}
	if totalC, err = redis.Int(conn.Do("ZCOUNT", key, 0, now)); err != nil {
		log.Error("conn.Do(ZCOUNT,%s,0,%d) error(%v)", key, now, err)
		return
	}
	if moniC, err = redis.Int(conn.Do("ZCOUNT", key, tFrom, tTo)); err != nil {
		log.Error("conn.Do(ZCOUNT,%s,%d,%d) error(%v)", key, tFrom, tTo, err)
		return
	}
	var oldest map[string]string //进入列表最久的项
	oldest, err = redis.StringMap(conn.Do("ZRANGE", key, 0, 0, "WITHSCORES"))
	for _, t := range oldest {
		var i int
		if i, err = strconv.Atoi(t); err != nil {
			return
		}
		maxT = int(now) - i
	}
	res = &monitor.Stats{
		TotalCount: totalC,
		MoniCount:  moniC,
		MaxTime:    maxT,
	}
	return
}

// GetAllRules 获取所有规则
func (d *Dao) GetAllRules(c context.Context, all bool) (rules []*monitor.Rule, err error) {
	var (
		conn = d.redis.Get(c)
		res  = make(map[string]string)
	)
	defer conn.Close()
	if res, err = redis.StringMap(conn.Do("HGETALL", monitor.RulesKey)); err != nil {
		if err != redis.ErrNil {
			log.Error("conn.Do(HGETALL, %s) error(%v)", monitor.RulesKey, err)
			return
		}
	}
	for _, v := range res {
		rule := &monitor.Rule{}
		if err = json.Unmarshal([]byte(v), rule); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", v, err)
			break
		}
		if !all && rule.State != 1 {
			continue
		}
		rules = append(rules, rule)
	}
	return
}

// GetRules 获取业务下的规则
func (d *Dao) GetRules(c context.Context, tp, bid int8, all bool) (rules []*monitor.Rule, err error) {
	if rules, err = d.GetAllRules(c, all); err != nil {
		return
	}
	for k := 0; k < len(rules); k++ {
		v := rules[k]
		if v.Type != tp || v.Business != bid { //去掉非当前业务开头的配置
			rules = append(rules[:k], rules[k+1:]...)
			k--
			continue
		}
	}
	return
}

// SetRule 修改/添加监控规则
func (d *Dao) SetRule(c context.Context, rule *monitor.Rule) (err error) {
	if rule.ID == 0 {
		if rule.ID, err = d.RuleIDIncKey(c); err != nil {
			return
		}
	}
	var (
		conn  = d.redis.Get(c)
		field = fmt.Sprintf(FieldKeyFormat, rule.Type, rule.Business, rule.ID)
		bs    []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(rule); err != nil {
		log.Error("json.Marshal(%v) error(%v)", rule, err)
		return
	}
	if _, err = conn.Do("HSET", monitor.RulesKey, field, bs); err != nil {
		log.Error("conn.Do(HSET,%s,%s,%s) error(%v)", monitor.RulesKey, field, bs, err)
		return
	}
	return
}

// GetRule 获取某条监控规则
func (d *Dao) GetRule(c context.Context, tp, bid int8, id int64) (rule *monitor.Rule, err error) {
	var (
		conn  = d.redis.Get(c)
		field = fmt.Sprintf(FieldKeyFormat, tp, bid, id)
		bs    []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("HGET", monitor.RulesKey, field)); err != nil {
		log.Error("conn.Do(HGET,%s,%s) error(%v)", monitor.RulesKey, field, err)
		return
	}
	rule = &monitor.Rule{}
	if err = json.Unmarshal(bs, rule); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", bs, err)
		return
	}
	return
}

// SetRuleState 修改监控规则的状态
func (d *Dao) SetRuleState(c context.Context, tp, bid int8, id int64, state int8) (err error) {
	var (
		rule *monitor.Rule
	)
	if rule, err = d.GetRule(c, tp, bid, id); err != nil {
		return
	}
	rule.State = state
	if err = d.SetRule(c, rule); err != nil {
		return
	}
	return
}

// RuleIDIncKey 自增配置id
func (d *Dao) RuleIDIncKey(c context.Context) (id int64, err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if id, err = redis.Int64(conn.Do("INCR", monitor.RuleIDIncKey)); err != nil {
		log.Error("conn.Do(INCR,%s) error(%v)", monitor.RuleIDIncKey, err)
	}
	return
}

// BusStatsKeys 获取某业务统计的所有keys
func (d *Dao) BusStatsKeys(c context.Context, bid int8) (prefix string, keys []string, err error) {
	var (
		conf *monitor.KeyConf
		ok   bool
	)
	if conf, ok = monitor.RedisKeyConf[bid]; !ok {
		err = errors.New("业务redis key配置不存在")
		log.Error("d.BusStatsKeys(%d) error(%v)", bid, err)
		return
	}
	prefix = fmt.Sprintf(monitor.BusPrefix, bid)
	//TODO 递归实现
	if bid == monitor.BusVideo {
		for _, v := range conf.KFields["state"] {
			key := prefix + fmt.Sprintf(monitor.SuffixVideo, v)
			keys = append(keys, key)
		}
	} else if bid == monitor.BusArc {
		for _, round := range conf.KFields["round"] {
			for _, state := range conf.KFields["state"] {
				key := prefix + fmt.Sprintf(monitor.SuffixArc, round, state)
				keys = append(keys, key)
			}
		}
	}
	return
}

// StayOids 获取多个key 中的滞留oid
func (d *Dao) StayOids(c context.Context, rule *monitor.Rule, keys []string) (oidMap map[int64]int, total int, err error) {
	var (
		conn     = d.redis.Get(c)
		intMap   map[string]int
		min, max int64
		now      = time.Now().Unix()
	)
	defer conn.Close()
	oidMap = make(map[int64]int)
	intMap = make(map[string]int)
	if _, ok := rule.RuleConf.NotifyCdt["time"]; !ok {
		log.Error("StayOids(%+v) Rule配置中NotifyCdt 没有time", *rule)
		err = errors.New(fmt.Sprintf("Rule(%d) NotifyCdt Error: no time", rule.ID))
		return
	}
	timeConf := rule.RuleConf.NotifyCdt["time"]
	switch timeConf.Comp {
	case monitor.CompGT:
		min = 0
		max = now - timeConf.Value
	case monitor.CompLT:
		min = now - timeConf.Value
		max = now
	default:
		log.Error("StayOids(%+v) Rule配置NotifyCdt中time的表达式错误", *rule)
		err = errors.New(fmt.Sprintf("Rule(%d) NotifyCdt Error: unknown time comp", rule.ID))
		return
	}
	//key排序
	sort.Strings(keys)
	//计算count 翻页
	for _, key := range keys {
		count := 0
		if count, err = redis.Int(conn.Do("ZCOUNT", key, min, max)); err != nil {
			log.Error("redis.Int(conn.Do(\"ZCOUNT\", %s, %d, %d)) error(%v)", key, min, max, err)
			return
		}
		total += count

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
	}
	return
}

// RemMonitorStats remove stay stats
func (d *Dao) RemMonitorStats(c context.Context, key string, oid int64) (err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, oid); err != nil {
		log.Error("conn.Do(ZADD, %s, %d) error(%v)", key, oid, err)
	}
	return
}
