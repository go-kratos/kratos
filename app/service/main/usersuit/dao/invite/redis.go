package dao

import (
	"context"
	"strconv"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixBuyFlag   = "b_"
	_prefixApplyFlag = "a_"
)

func keyBuyFlag(mid int64) string {
	return _prefixBuyFlag + strconv.FormatInt(mid, 10)
}

func keyApplyFlag(code string) string {
	return _prefixApplyFlag + code
}

// pingRedis ping redis.
func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = redis.String(conn.Do("PING")); err != nil {
		log.Error("redis.String(conn.Do(PING)) error(%v)", err)
	}
	return
}

// SetBuyFlagCache set buy flag cache
func (d *Dao) SetBuyFlagCache(c context.Context, mid int64, f string) (ok bool, err error) {
	key := keyBuyFlag(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	var res interface{}
	if res, err = conn.Do("SET", key, f, "EX", d.inviteExpire, "NX"); err != nil {
		log.Error("conn.Do(SET, %s, %s, EX, %d, NX)", key, f)
		return
	}
	if res != nil {
		ok = true
	}
	return
}

// DelBuyFlagCache del buy flag cache
func (d *Dao) DelBuyFlagCache(c context.Context, mid int64) (err error) {
	key := keyBuyFlag(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL, %s)", key)
	}
	return
}

// SetApplyFlagCache set apply flag cache
func (d *Dao) SetApplyFlagCache(c context.Context, code, f string) (ok bool, err error) {
	key := keyApplyFlag(code)
	conn := d.redis.Get(c)
	defer conn.Close()
	var res interface{}
	if res, err = conn.Do("SET", key, f, "EX", d.inviteExpire, "NX"); err != nil {
		log.Error("conn.Do(SET, %s, %s, EX, %d, NX)", key, f, d.inviteExpire)
		return
	}
	if res != nil {
		ok = true
	}
	return
}

// DelApplyFlagCache del apply Flag Cache
func (d *Dao) DelApplyFlagCache(c context.Context, code string) (err error) {
	key := keyApplyFlag(code)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL, %s)", key)
	}
	return
}
