package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/job/main/dm2/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// dm xml list v1
	_prefixDM = "dm_v1_%d_%d" // dm_v1_tpe_oid
	divide    = 34359738368   // 2^35
)

func keyDM(tp int32, oid int64) (key string) {
	return fmt.Sprintf(_prefixDM, tp, oid)
}

// 弹幕在redis sortset 中的score
// 通过score保证弹幕在缓存中的排序为:普通弹幕、普通弹幕中的保护弹幕、字幕弹幕、脚本弹幕
func score(dm *model.DM) (score float64) {
	// NOTE redis score最多17位表示，这里采用整数十位+小数部分十位
	v := dm.ID / divide                                // 2^63 / 2^35 = 2^28-1 整数部分最大值：268435455
	k := dm.ID % divide                                // 精度8位，最后5位可忽略
	r := int64(dm.Pool)<<29 | int64(dm.Attr)&1<<28 | v // NOTE v should less than 2^28
	score, _ = strconv.ParseFloat(fmt.Sprintf("%d.%d", r, k), 64)
	return
}

// AddDMCache add dm to redis.
func (d *Dao) AddDMCache(c context.Context, dm *model.DM) (err error) {
	var (
		conn  = d.dmRds.Get(c)
		value []byte
		key   = keyDM(dm.Type, dm.Oid)
	)
	defer conn.Close()
	if value, err = dm.Marshal(); err != nil {
		log.Error("dm.Marshal(%v) error(%v)", dm, err)
		return
	}
	if err = conn.Send("ZADD", key, score(dm), value); err != nil {
		log.Error("conn.Send(ZADD %v) error(%v)", dm, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.dmRdsExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// SetDMCache flush dm list to redis.
func (d *Dao) SetDMCache(c context.Context, tp int32, oid int64, dms []*model.DM) (err error) {
	var (
		value []byte
		conn  = d.dmRds.Get(c)
		key   = keyDM(tp, oid)
	)
	defer conn.Close()
	for _, dm := range dms {
		if value, err = dm.Marshal(); err != nil {
			log.Error("dm.Marshal(%v) error(%v)", dm, err)
			return
		}
		if err = conn.Send("ZADD", key, score(dm), value); err != nil {
			log.Error("conn.Send(ZADD %v) error(%v)", dm, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, d.dmRdsExpire); err != nil {
		log.Error("conn.Send(EXPIRE %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < len(dms)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelDMCache delete redis cache of oid.
func (d *Dao) DelDMCache(c context.Context, tp int32, oid int64) (err error) {
	var (
		key  = keyDM(tp, oid)
		conn = d.dmRds.Get(c)
	)
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// ExpireDMCache expire dm.
func (d *Dao) ExpireDMCache(c context.Context, tp int32, oid int64) (ok bool, err error) {
	key := keyDM(tp, oid)
	conn := d.dmRds.Get(c)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.dmRdsExpire)); err != nil {
		log.Error("conn.Do(EXPIRE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// DMCache 获取redis列表中的弹幕.
func (d *Dao) DMCache(c context.Context, tp int32, oid int64) (res [][]byte, err error) {
	conn := d.dmRds.Get(c)
	key := keyDM(tp, oid)
	if res, err = redis.ByteSlices(conn.Do("ZRANGE", key, 0, -1)); err != nil {
		log.Error("conn.Do(ZRANGE %s) error(%v)", key, err)
	}
	conn.Close()
	return
}

// TrimDMCache 从redis列表中pop掉count条弹幕.
func (d *Dao) TrimDMCache(c context.Context, tp int32, oid, count int64) (err error) {
	conn := d.dmRds.Get(c)
	key := keyDM(tp, oid)
	if _, err = conn.Do("ZREMRANGEBYRANK", key, 0, count-1); err != nil {
		log.Error("conn.Do(ZREMRANGEBYRANK %s) error(%v)", key, err)
	}
	conn.Close()
	return
}
