package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_rptCnt                 = "rc_"
	_rptBrig                = "rb_"
	_prfxPaUsrCntKey        = "pa_user_count_"
	_prfxUdAccCntKey        = "usr_dm_acc_cnt_"
	_prfxRecallCntKey       = "recall_cnt_"
	_crontabRedisLock       = "crontab_redis_lock"
	_crontabRedisLockExpire = 2 * 3600
)

// tomorrowUnixTime 明天unix time
func tomorrowUnixTime() int64 {
	tomorrow := time.Now().Add(24 * time.Hour)
	year, month, day := tomorrow.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
}

// paUsrCntKey 用户申请保护弹幕计数 redis key
func keyProtectApplyCmt() string {
	return _prfxPaUsrCntKey + time.Now().Format("20060102")
}

func keyRptCnt() string {
	return _rptCnt + time.Now().Format("04")
}

func keyRptBrig(mid int64) string {
	return _rptBrig + strconv.FormatInt(mid, 10)
}

func keyRecallCnt() string {
	return _prfxRecallCntKey + time.Now().Format("20060102")
}

// UptUsrPaCnt 设置申请保护弹幕数
func (d *Dao) UptUsrPaCnt(c context.Context, uid int64, count int64) (err error) {
	var (
		key, expire = keyProtectApplyCmt(), tomorrowUnixTime()
		conn        = d.redisDM.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("ZINCRBY", key, count, uid); err != nil {
		log.Error("conn.Send(ZADD,%s,%d,%d) error(%v)", key, count, uid, err)
		return
	}
	if err = conn.Send("EXPIREAT", key, expire); err != nil {
		log.Error("conn.Send(EXPIREAT key(%s) expire(%d)) error(%v)", key, expire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// PaUsrCnt 保护弹幕申请计数
func (d *Dao) PaUsrCnt(c context.Context, uid int64) (cnt int, err error) {
	var (
		key  = keyProtectApplyCmt()
		conn = d.redisDM.Get(c)
	)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("ZSCORE", key, uid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(ZSCORE, %s,%d) error(%v)", key, uid, err)
		}
	}
	return
}

// UsrDMAccCnt 我的弹幕访问计数
func (d *Dao) UsrDMAccCnt(c context.Context, uid int64, t int64) (cnt int64, err error) {
	var (
		key  = _prfxUdAccCntKey + strconv.FormatInt(uid, 10)
		conn = d.redisDM.Get(c)
	)
	defer conn.Close()
	if cnt, err = redis.Int64(conn.Do("INCRBY", key, t)); err != nil {
		log.Error("conn.Do(INCRBY, %s,%d) error(%v)", key, t, err)
		return
	}
	if cnt != t {
		return
	}
	_, err = redis.Int(conn.Do("EXPIREAT", key, t+5))
	if err != nil {
		log.Error("conn.Do(EXPIREAT, %s,%d) error(%v)", key, t+5, err)
	}
	return
}

// RecallCnt 撤回弹幕计数
func (d *Dao) RecallCnt(c context.Context, uid int64) (cnt int, err error) {
	var (
		key  = keyRecallCnt()
		conn = d.redisDM.Get(c)
	)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("ZSCORE", key, uid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("redis.Int(conn.Do(GET, %s,%d)) error(%v)", key, uid, err)
		}
	}
	return
}

// UptRecallCnt 更新撤回弹幕计数
func (d *Dao) UptRecallCnt(c context.Context, uid int64) (err error) {
	var (
		key  = keyRecallCnt()
		conn = d.redisDM.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("ZINCRBY", key, 1, uid); err != nil {
		log.Error("conn.Send(ZINCRBY, %s,%d) error(%v)", key, 1, err)
		return
	}
	if err = conn.Send("EXPIRE", key, 3600*24); err != nil {
		log.Error("conn.Send(EXPIRE, %s,%d) error(%v)", key, 3600*24, err)
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

// PaLock redis定时任务锁，保证只有一个地方在执行定时任务
func (d *Dao) PaLock(c context.Context, key string) (incr int, err error) {
	var (
		conn = d.redisDM.Get(c)
		r    interface{}
	)
	key = _crontabRedisLock + key
	defer conn.Close()
	if err = conn.Send("INCRBY", key, 1); err != nil {
		log.Error("conn.Send(INCRBY,%s) error(%v)", key, err)
		return
	}
	if err = conn.Send("EXPIRE", key, _crontabRedisLockExpire); err != nil {
		log.Error("conn.Send(EXPIRE key(%s) expire(%d)) error(%v)", key, _crontabRedisLockExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if r, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
		return
	}
	if incr, err = redis.Int(r, err); err != nil {
		log.Error("redis.Int error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
	}
	return
}

// RptCnt 弹幕举报一分钟计数
func (d *Dao) RptCnt(c context.Context, uid int64) (n int, err error) {
	var (
		k    = keyRptCnt()
		conn = d.redisDM.Get(c)
	)
	if n, err = redis.Int(conn.Do("ZSCORE", k, uid)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(ZSCORE, %s,%d) error(%v)", k, uid, err)
		}
	}
	conn.Close()
	return
}

// UptRptCnt 更新一分钟的弹幕举报数
func (d *Dao) UptRptCnt(c context.Context, uid int64) (n int64, err error) {
	var (
		k    = keyRptCnt()
		conn = d.redisDM.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("ZINCRBY", k, 1, uid); err != nil {
		log.Error("conn.Send(ZINCRBY, %s,%d) error(%v)", k, 1, err)
		return
	}
	if err = conn.Send("EXPIRE", k, 60); err != nil {
		log.Error("conn.Send(EXPIRE, %s,%d) error(%v)", k, 60, err)
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

// AddRptBrig 将用户禁闭，30分钟不让举报
func (d *Dao) AddRptBrig(c context.Context, mid int64) (err error) {
	var (
		k    = keyRptBrig(mid)
		conn = d.redisDM.Get(c)
	)
	if _, err = conn.Do("SETEX", k, 1800, time.Now().Unix()); err != nil {
		log.Error("conn.Do(SETEX, %s,%d) error(%v)", k, 1, err)
	}
	conn.Close()
	return
}

// RptBrigTime 获得禁闭的开始时间
func (d *Dao) RptBrigTime(c context.Context, mid int64) (t int64, err error) {
	var (
		k    = keyRptBrig(mid)
		conn = d.redisDM.Get(c)
	)
	if t, err = redis.Int64(conn.Do("GET", k)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s,%d) error(%v)", k, 1, err)
		}
	}
	conn.Close()
	return
}
