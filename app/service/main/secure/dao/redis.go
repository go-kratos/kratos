package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/service/main/secure/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixMsg         = "m_"
	_prefixUnNotify    = "d_%d%d_%d"
	_prefixCount       = "c_%d%d_%d"
	_prefixChangePWD   = "cpwd_%d"
	_prefixDoublecheck = "dc_%d"

	_expire    = 24 * 3600
	_expirePWD = 30 * 24 * 3600
)

func doubleCheckKey(mid int64) string {
	return fmt.Sprintf(_prefixDoublecheck, mid)
}

func changePWDKey(mid int64) string {
	return fmt.Sprintf(_prefixChangePWD, mid)
}
func msgKey(mid int64) string {
	return _prefixMsg + strconv.FormatInt(mid, 10)
}

func unnotifyKey(mid int64) string {
	t := time.Now()
	return fmt.Sprintf(_prefixUnNotify, t.Month(), t.Day(), mid)
}

func countKey(mid int64) string {
	t := time.Now()
	return fmt.Sprintf(_prefixCount, t.Month(), t.Day(), mid)
}

// AddExpectionMsg add user login expection msg.
func (d *Dao) AddExpectionMsg(c context.Context, l *model.Log) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
		key  = msgKey(l.Mid)
	)
	defer conn.Close()
	if bs, err = json.Marshal(l); err != nil {
		log.Error("json.Marshal(%v) err(%v)", l, err)
		return
	}
	if _, err = conn.Do("SETEX", key, d.expire, bs); err != nil {
		log.Error("conn.Set msg:%v err(%v)", l, err)
	}
	return
}

// ExpectionMsg get user expection msg.
func (d *Dao) ExpectionMsg(c context.Context, mid int64) (msg *model.Log, err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", msgKey(mid))); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("conn.GET(mid %d) ,err(%v)", mid, err)
		return
	}
	msg = &model.Log{}
	if err = json.Unmarshal(bs, msg); err != nil {
		log.Error("json.Unmarshal err(%v)", err)
	}
	return
}

// AddUnNotify user unnotiry uuid.
func (d *Dao) AddUnNotify(c context.Context, mid int64, uuid string) (err error) {
	var (
		conn = d.redis.Get(c)
		key  = unnotifyKey(mid)
	)
	defer conn.Close()
	if err = conn.Send("SADD", key, uuid); err != nil {
		log.Error("conn.SADD mid:%d err(%v)", mid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, _expire); err != nil {
		log.Error("EXPIRE key :%d err %d", key, err)
		return
	}
	conn.Flush()
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Recive err %v", err)
			return
		}
	}
	return
}

// DelUnNotify del user unnotify record.
func (d *Dao) DelUnNotify(c context.Context, mid int64) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("DEL", unnotifyKey(mid)); err != nil {
		log.Error("conn.DEL mid:%d err:%v", mid, err)
	}
	conn.Close()
	return
}

// UnNotify check if not send notify to user of uuid deveice.
func (d *Dao) UnNotify(c context.Context, mid int64, uuid string) (b bool, err error) {
	conn := d.redis.Get(c)
	if b, err = redis.Bool(conn.Do("SISMEMBER", unnotifyKey(mid), uuid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		log.Error("conn.SISMEMBER (mid:%d) ERR(%v)", mid, err)
	}
	conn.Close()
	return
}

// Count get user close notify count.
func (d *Dao) Count(c context.Context, mid int64, uuid string) (count int64, err error) {
	conn := d.redis.Get(c)
	if count, err = redis.Int64(conn.Do("HGET", countKey(mid), uuid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		log.Error("conn.GET mid:%d err(%v)", mid, err)
	}
	conn.Close()
	return
}

// AddCount add user unnotify count daily.
func (d *Dao) AddCount(c context.Context, mid int64, uuid string) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HINCRBY", countKey(mid), uuid, 1); err != nil {
		log.Error("conn.INCR mid:%d err:%v", mid, err)
		return
	}
	if err = conn.Send("EXPIRE", countKey(mid), _expire); err != nil {
		log.Error("conn.EXPIRE mid:%d err:%v", mid, err)
		return
	}
	conn.Flush()
	for i := 0; i < 2; i++ {
		if _, err1 := conn.Receive(); err1 != nil {
			log.Error("conn.Receive err(%v)", err1)
			return
		}
	}
	return
}

// AddChangePWDRecord set user change passwd record to cache.
func (d *Dao) AddChangePWDRecord(c context.Context, mid int64) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SETEX", changePWDKey(mid), _expirePWD, 1); err != nil {
		log.Error("d.ChangePWDRecord(mid %d) err(%v)", mid, err)
	}
	conn.Close()
	return
}

// ChangePWDRecord check if user had change pwd recently one month.
func (d *Dao) ChangePWDRecord(c context.Context, mid int64) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if b, err = redis.Bool(conn.Do("GET", changePWDKey(mid))); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.ChangePWDRecord err(%v)", err)
	}
	return
}

// DelCount  del count
// for testing clear data.
func (d *Dao) DelCount(c context.Context, mid int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", countKey(mid))
	return
}

// AddDCheckCache add double check cache.
func (d *Dao) AddDCheckCache(c context.Context, mid int64) (err error) {
	conn := d.redis.Get(c)
	if _, err = conn.Do("SETEX", doubleCheckKey(mid), d.doubleCheckExpire, 1); err != nil {
		log.Error("d.AddDCheckCache(mid %d) err(%v)", mid, err)
	}
	conn.Close()
	return
}

// DCheckCache check if user had notify by double check.
func (d *Dao) DCheckCache(c context.Context, mid int64) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if b, err = redis.Bool(conn.Do("GET", doubleCheckKey(mid))); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.DCheckCache err(%v)", err)
	}
	return
}
