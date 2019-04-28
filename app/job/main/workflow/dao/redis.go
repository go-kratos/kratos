package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyUperWeightGroup        = "platform_uper_weight_param_2" // 工作台任务权重参数(用户维度)
	_prefixKeyMissionSortedSet = "platform_missions_"           // 工作台任务池 有序集合(权重排序)
	_prefixKeySingleExpire     = "platform_single_expire_"      // 单条工单的开始处理时间 有序集合
	_relatedMissions           = "platfrom_missions_%d_%d"      // 当前认领任务
)

// SetList .
func (d *Dao) SetList(c context.Context, key string, ids []int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, id := range ids {
		now := time.Now().Format("2006-01-02 15:04:05")
		log.Info("enter queue id is %d time is %s", id, now)
		if err = conn.Send("LPUSH", key, id); err != nil {
			log.Error("d.LPUSH error(%v)", err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
	}
	return
}

// ExistKey .
func (d *Dao) ExistKey(c context.Context, key string) (exist bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	exist, err = redis.Bool(conn.Do("EXISTS", key))
	return
}

// SetString .
func (d *Dao) SetString(c context.Context, key, val string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("SET", key, val)

	return
}

// SetCrash .
func (d *Dao) SetCrash(c context.Context) (err error) {
	key, val := "dead", "1"
	err = d.SetString(c, key, val)

	return
}

// IsCrash .
func (d *Dao) IsCrash(c context.Context) (exist bool, err error) {
	key := "dead"
	exist, err = d.ExistKey(c, key)
	return
}

// UperInfoCache 读取用户维度的申诉weight计算参数
func (d *Dao) UperInfoCache(c context.Context, apIDs []int64) (params []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}
	args = args.Add(_keyUperWeightGroup)
	for _, apid := range apIDs {
		args = args.Add(apid)
	}
	if params, err = redis.Int64s(conn.Do("HMGET", args...)); err != nil {
		log.Error("HMGET %v error(%v)", args, err)
	}
	return
}

// DelUperInfo .
func (d *Dao) DelUperInfo(c context.Context, mids []int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}
	args = args.Add(_keyUperWeightGroup)
	for _, mid := range mids {
		args = args.Add(mid)
	}
	if _, err = conn.Do("HDEL", args...); err != nil {
		log.Error("HDEL %v error(%v)", args, err)
	}
	return
}

// SetWeightSortedSet 覆盖sorted set
func (d *Dao) SetWeightSortedSet(c context.Context, bid int, newWeight map[int64]int64) (err error) {
	key := _prefixKeyMissionSortedSet + strconv.Itoa(bid)
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}
	args = args.Add(key)
	for id, weight := range newWeight {
		args = args.Add(weight, id)
	}
	if _, err = conn.Do("ZADD", args...); err != nil { // ZADD key score member [[score member] [score member] ...]
		log.Error("ZADD %v error(%v)", args, err)
	}
	return
}

// SingleExpire 获取所有的 single expire 信息
func (d *Dao) SingleExpire(c context.Context, bid int) (delIDs []int64, err error) {
	key := _prefixKeySingleExpire + strconv.Itoa(bid)
	conn := d.redis.Get(c)
	defer conn.Close()
	floorTime := time.Now().Add(-480 * time.Second).Unix()
	// fixme if too hash field too many
	delIDs, err = redis.Int64s(conn.Do("ZRANGEBYSCORE", key, "0", floorTime, "LIMIT", "0", "50"))
	return
}

// DelSingleExpire .
func (d *Dao) DelSingleExpire(c context.Context, bid int, ids []int64) (err error) {
	key := _prefixKeySingleExpire + strconv.Itoa(bid)
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}
	args = args.Add(key)
	for _, apID := range ids {
		args = args.Add(apID)
	}
	if _, err = conn.Do("ZREM", args...); err != nil {
		log.Error("ZREM %v error(%v)", args, err)
	}
	return
}

// DelRelatedMissions .
func (d *Dao) DelRelatedMissions(ctx context.Context, bid, transAdmin int, ids []int64) (err error) {
	if len(ids) == 0 {
		return
	}
	key := fmt.Sprintf(_relatedMissions, transAdmin, bid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	args := redis.Args{}.Add(key)
	for _, id := range ids {
		args = args.Add(id)
	}
	_, err = conn.Do("ZREM", args...)
	return
}
