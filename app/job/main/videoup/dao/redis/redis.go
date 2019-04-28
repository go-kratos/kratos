package redis

import (
	"context"
	"encoding/json"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_failList      = "f_list"
	_clickSort     = "c_sort"
	_prefixMonitor = "va_sc_set"       // key of monitor on videoup job second round hit cache set
	_expire        = 30 * 24 * 60 * 60 // 30 days
	_clickExpire   = 10 * 24 * 60 * 60 // 10 days
)

var (
	_twepoch = time.Date(2017, time.Month(9), 1, 0, 0, 0, 0, time.Local).Unix()
)

// PushFail rpush fail item to redis
func (d *Dao) PushFail(c context.Context, a interface{}) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(a); err != nil {
		log.Error("json.Marshal(%v) error(%v)", a, err)
		return
	}
	if _, err = conn.Do("RPUSH", _failList, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s, %s) error(%v)")
	}
	return
}

// PopFail lpop fail item from redis
func (d *Dao) PopFail(c context.Context) (bs []byte, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", _failList)); err != nil && err != redis.ErrNil {
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%v)", _failList, err)
		return
	}
	return
}

// PushQueue rpush fail item to redis
func (d *Dao) PushQueue(c context.Context, a interface{}, queue string) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(a); err != nil {
		log.Error("json.Marshal(%v) error(%v)", a, err)
		return
	}
	if _, err = conn.Do("RPUSH", queue, bs); err != nil {
		log.Error("conn.Do(RPUSH, %s, %s) error(%v)")
	}
	return
}

// PopQueue lpop fail item from redis
func (d *Dao) PopQueue(c context.Context, queue string) (bs []byte, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("LPOP", queue)); err != nil && err != redis.ErrNil {
		log.Error("redis.Bytes(conn.Do(LPOP, %s)) error(%v)", queue, err)
		return
	}
	return
}

// AddFilename set filename expire time
func (d *Dao) AddFilename(c context.Context, filename string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SETEX", filename, _expire, time.Now().Unix()); err != nil {
		log.Error("conn.Do(SETEX, %s, %d, %d) error(%v)", filename, _expire, time.Now().Unix(), err)
	}
	return
}

// DelFilename set filename expire time
func (d *Dao) DelFilename(c context.Context, filename string) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", filename); err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", filename, err)
	}
	return
}

// AddArcClick add archive click into redis
func (d *Dao) AddArcClick(c context.Context, aid int64, click int) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	now := time.Now().Unix() - _twepoch
	if _, err = conn.Do("ZADD", _clickSort, now<<26|int64(click), aid); err != nil {
		log.Error("conn.Do(ZADD,%s,%d)", _clickSort, aid)
		return
	}
	conn.Do("ZREMRANGEBYSCORE", _clickSort, "-inf", (now-_clickExpire)<<26)
	return
}

// ArcClick find archive click from redis
func (d *Dao) ArcClick(c context.Context, aid int64) (click int, err error) {
	var (
		conn  = d.redis.Get(c)
		score int64
	)
	defer conn.Close()
	if score, err = redis.Int64(conn.Do("ZSCORE", _clickSort, aid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(ZSCORE, %d) error(%v)", aid, err)
		}
		return
	}
	click = int(score & 0x3ffffff) // 0x3ffffff = 26bit
	return
}

func monitorKey() string {
	return _prefixMonitor
}

// SetMonitorCache set monitor cache
func (d *Dao) SetMonitorCache(c context.Context, aid int64) (had bool, err error) {
	var (
		key      = monitorKey()
		conn     = d.redis.Get(c)
		firstAdd bool
	)
	defer conn.Close()
	if firstAdd, err = redis.Bool(conn.Do("SADD", key, aid)); err != nil {
		log.Error("SADD conn.Do error(%v)", err)
		had = false
		return
	}
	had = !firstAdd
	return
}

// DelMonitorCache del monitor cache
func (d *Dao) DelMonitorCache(c context.Context, aid int64) (err error) {
	var (
		key  = monitorKey()
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = redis.Int64(conn.Do("SREM", key, aid)); err != nil {
		log.Error("SREM conn.Do(%s,%d) err(%v)", key, aid, err)
	}
	return
}
