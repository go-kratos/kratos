package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/job/main/dm/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyIdxContent = "c_%d_%d"    // dm content hash(c_type_oid, dmid, xml)
	_keyTrimQueue  = "tq_n_%d_%d" // trim queue if dm_count > dm_maxlimit

	divide = int64(34359738368) // 2^35
)

// keyIdxContent return key of different dm.
func keyIdxContent(typ int32, oid int64) string {
	return fmt.Sprintf(_keyIdxContent, typ, oid)
}

// keyIdxQueue return trim queue key.
func keyTrimQueue(tp int32, oid int64) string {
	return fmt.Sprintf(_keyTrimQueue, tp, oid)
}

// ExpireTrimQueue set expire time of index.
func (d *Dao) ExpireTrimQueue(c context.Context, tp int32, oid int64) (ok bool, err error) {
	key := keyTrimQueue(tp, oid)
	conn := d.redis.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.redisExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

func score(attr int32, id int64) (score float64) {
	v := id / divide           // 2^63 / 2^35 = 2^28-1 整数部分最大值：268435455
	k := id % divide           // 精度8位，最后5位可忽略
	r := int64(attr)&1<<28 | v // NOTE v should less than 2^28
	score, _ = strconv.ParseFloat(fmt.Sprintf("%d.%d", r, k), 64)
	return
}

// AddTrimQueueCache add dm index into trim queue.
func (d *Dao) AddTrimQueueCache(c context.Context, tp int32, oid int64, trims []*model.Trim) (count int64, err error) {
	var (
		key  = keyTrimQueue(tp, oid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, trim := range trims {
		if err = conn.Send("ZADD", key, score(trim.Attr, trim.ID), trim.ID); err != nil {
			log.Error("conn.Send(ZADD %s %v) error(%v)", key, trim, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("ZCARD", key); err != nil {
		log.Error("conn.Send(ZCARD %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(trims)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	if count, err = redis.Int64(conn.Receive()); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// FlushTrimCache flush trim queue cache.
func (d *Dao) FlushTrimCache(c context.Context, tp int32, oid int64, trims []*model.Trim) (err error) {
	var (
		key  = keyTrimQueue(tp, oid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, trim := range trims {
		if err = conn.Send("ZADD", key, score(trim.Attr, trim.ID), trim.ID); err != nil {
			log.Error("conn.Send(ZADD %s %v) error(%v)", key, trim, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.redisExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush(%s) error(%v)", key, err)
		return
	}
	for i := 0; i < len(trims)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// ZRemTrimCache ZREM trim from trim queue.
func (d *Dao) ZRemTrimCache(c context.Context, tp int32, oid int64, dmID int64) (err error) {
	var (
		key  = keyTrimQueue(tp, oid)
		conn = d.redis.Get(c)
	)
	if _, err = conn.Do("ZREM", key, dmID); err != nil {
		log.Error("conn.Do(ZREM %s %v) error(%v)", key, dmID, err)
	}
	conn.Close()
	return
}

// TrimCache trim trim queue and return trimed dmid from trim queue.
func (d *Dao) TrimCache(c context.Context, tp int32, oid, count int64) (dmids []int64, err error) {
	var (
		key    = keyTrimQueue(tp, oid)
		conn   = d.redis.Get(c)
		replys [][]byte
		dmID   int64
	)
	defer conn.Close()
	if replys, err = redis.ByteSlices(conn.Do("ZRANGE", key, 0, count-1)); err != nil {
		log.Error("conn.Do(ZRANGE %s) error(%v)", key, err)
		return
	}
	for _, reply := range replys {
		if err = json.Unmarshal(reply, &dmID); err != nil {
			log.Error("json.Unmarshal(%s) error(v)", string(reply), err)
			return
		}
		dmids = append(dmids, dmID)
	}
	if len(dmids) > 0 {
		if _, err = conn.Do("ZREMRANGEBYRANK", key, 0, len(dmids)-1); err != nil {
			log.Error("conn.Do(ZREMRANGEBYRANK %s) error(%v)", key, err)
		}
	}
	return
}

// DelIdxContentCaches del index content cache.
func (d *Dao) DelIdxContentCaches(c context.Context, typ int32, oid int64, dmids ...int64) (err error) {
	key := keyIdxContent(typ, oid)
	conn := d.redis.Get(c)
	args := []interface{}{key}
	for _, dmid := range dmids {
		args = append(args, dmid)
	}
	if _, err = conn.Do("HDEL", args...); err != nil {
		log.Error("conn.Do(HDEL %s) error(%v)", key, err)
	}
	conn.Close()
	return
}
