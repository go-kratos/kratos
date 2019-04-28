package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/push-archive/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixUpperLimit    = "pau_%d"
	_prefixFanLimit      = "paf_%d"
	_statisticsKey       = "statistics_push_archive"
	_prefixPerUpperLimit = "perup_%d_%d"
)

func (d *Dao) do(c context.Context, command string, key string, args ...interface{}) (reply interface{}, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	values := []interface{}{key}
	if len(args) > 0 {
		values = append(values, args...)
	}
	reply, err = conn.Do(command, values...)
	return
}

func upperLimitKey(mid int64) string {
	return fmt.Sprintf(_prefixUpperLimit, mid)
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	if _, err = d.do(c, "SET", "PING", "PONG"); err != nil {
		PromError("redis: ping remote")
		log.Error("remote redis: conn.Do(SET,PING,PONG) error(%v)", err)
	}
	return
}

// ExistUpperLimitCache judge that whether upper push limit cache exists.
func (d *Dao) ExistUpperLimitCache(c context.Context, upper int64) (exist bool, err error) {
	key := upperLimitKey(upper)
	if exist, err = redis.Bool(d.do(c, "EXISTS", key)); err != nil {
		PromError("redis:读取upper推送限制")
		log.Error("ExistUpperLimitCache do(EXISTS, %s) error(%v)", key, err)
	}

	return
}

// AddUpperLimitCache sets upper push limit cache.
func (d *Dao) AddUpperLimitCache(c context.Context, upper int64) (err error) {
	key := upperLimitKey(upper)
	if _, err = d.do(c, "SETEX", key, d.UpperLimitExpire, ""); err != nil {
		PromError("redis:添加upper推送限制")
		log.Error("AddUpperLimitCache do(SETEX, %s) error(%v)", key, err)
	}
	return
}

//fanLimitKey 粉丝推送总次数限制key
func fanLimitKey(fan int64, relationType int) string {
	key := fmt.Sprintf(_prefixFanLimit, fan)
	if relationType != model.RelationSpecial {
		key = fmt.Sprintf("%s_%d", key, relationType)
	}
	return key
}

//GetFanLimitCache 读取粉丝限制的当前值
func (d *Dao) GetFanLimitCache(c context.Context, fan int64, relationType int) (limit int, err error) {
	key := fanLimitKey(fan, relationType)
	if limit, err = redis.Int(d.do(c, "GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("GetFanLimitCache do(GET) error(%v)", err)
		}
	}

	return
}

//AddFanLimitCache 添加粉丝限制的缓存
func (d *Dao) AddFanLimitCache(c context.Context, fan int64, relationType int, value int, expire int32) (err error) {
	key := fanLimitKey(fan, relationType)
	if _, err = d.do(c, "SETEX", key, expire, value); err != nil {
		log.Error("AddFanLimitCache do(SETEX) error(%v)", err)
		PromError("redis:添加fan推送限制")
	}
	return
}

//AddStatisticsCache 添加统计数据到redis
func (d *Dao) AddStatisticsCache(c context.Context, ps *model.PushStatistic) (err error) {
	psByte, err := json.Marshal(*ps)
	if err != nil {
		log.Error("AddStatisticsCache json.Marshal error(%v), pushstatistic(%v)", err, ps)
		return
	}

	key := _statisticsKey
	if _, err = d.do(c, "LPUSH", key, string(psByte)); err != nil {
		log.Error("AddStatisticsCache do(LPUSH, %s) error(%v) pushstatistic(%v)", key, err, ps)
		PromError("redis:添加统计数据")
	}

	return
}

//GetStatisticsCache 读取一条统计数据
func (d *Dao) GetStatisticsCache(c context.Context) (ps *model.PushStatistic, err error) {
	key := _statisticsKey
	psStr, err := redis.String(d.do(c, "RPOP", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("GetStatisticsCache do(RPOP, %s) error(%v)", key, err)
		}
		return
	}

	if err = json.Unmarshal([]byte(psStr), &ps); err != nil {
		log.Error("GetStatisticsCache json.Unmarshal error(%v), ps(%s)", err, psStr)
		return
	}
	return
}

//perUpperLimitKey  粉丝每个upper主的推送次数限制key
func perUpperLimitKey(fan int64, upper int64) string {
	return fmt.Sprintf(_prefixPerUpperLimit, fan, upper)
}

//GetPerUpperLimitCache 粉丝每个upper主的已推送次数
func (d *Dao) GetPerUpperLimitCache(c context.Context, fan int64, upper int64) (limit int, err error) {
	key := perUpperLimitKey(fan, upper)
	if limit, err = redis.Int(d.do(c, "GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("GetPerUpperLimitCache do(GET, %s) error(%v)", key, err)
		}
	}

	return
}

//AddPerUpperLimitCache 添加粉丝每个up主的推送次数
func (d *Dao) AddPerUpperLimitCache(c context.Context, fan int64, upper int64, value int, expire int32) (err error) {
	key := perUpperLimitKey(fan, upper)
	if _, err = d.do(c, "SETEX", key, expire, value); err != nil {
		log.Error("AddPerUpperLimitCache do(SETEX, %s, %d, %d) error(%v)", key, expire, value, err)
		PromError("redis:添加perupper推送限制")
	}
	return
}
