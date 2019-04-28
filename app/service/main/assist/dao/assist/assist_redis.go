package assist

import (
	"context"
	"go-common/app/service/main/assist/model/assist"
	"go-common/library/cache/redis"
	"go-common/library/log"
	xtime "go-common/library/time"
	"strconv"
	"time"
)

func (d *Dao) logCountKey(mid, assistMid int64) string {
	datetime := time.Now().Format("20060102")
	return datetime + strconv.FormatInt(mid, 10) + "_" + strconv.FormatInt(assistMid, 10)
}

func (d *Dao) assTotalKey(mid int64) string {
	datetime := time.Now().Format("20060102")
	return "ass_" + datetime + "_" + strconv.FormatInt(mid, 10)
}

func (d *Dao) assSameKey(mid int64) string {
	datetime := time.Now().Format("20060102")
	return "assSame_" + datetime + "_" + strconv.FormatInt(mid, 10)
}

func (d *Dao) assUpKey(assistMid int64) string {
	datetime := time.Now().Format("20060102")
	return "assUps_" + datetime + "_" + strconv.FormatInt(assistMid, 10)
}

// DailyLogCount get daily count by type
func (d *Dao) DailyLogCount(c context.Context, mid, assistMid, tp int64) (count int64, err error) {
	var (
		conn  = d.redis.Get(c)
		key   = d.logCountKey(mid, assistMid)
		field = strconv.FormatInt(tp, 10)
	)
	defer conn.Close()
	if err = conn.Send("HGET", key, field); err != nil {
		log.Error("conn.Send HGET (%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	count, err = redis.Int64(conn.Receive())
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			count = 0
			log.Info("conn.Receive(HGET, %d) error(%v)", mid, err)
		} else {
			log.Error("conn.Receive(HGET, %d) error(%v)", mid, err)
		}
		return
	}
	return
}

// IncrLogCount incr daily count
func (d *Dao) IncrLogCount(c context.Context, mid, assistMid, tp int64) (err error) {
	var (
		conn  = d.redis.Get(c)
		key   = d.logCountKey(mid, assistMid)
		field = strconv.FormatInt(tp, 10)
	)
	defer conn.Close()
	_ = conn.Send("EXPIRE", key, d.redisExpire)
	if err = conn.Send("HINCRBY", key, field, 1); err != nil {
		log.Error("conn.Send HINCRBY (%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	return
}

// TotalAssCnt get daily total add count
func (d *Dao) TotalAssCnt(c context.Context, mid int64) (count int64, err error) {
	var (
		conn  = d.redis.Get(c)
		key   = d.assTotalKey(mid)
		field = strconv.FormatInt(mid, 10)
	)
	defer conn.Close()
	if err = conn.Send("HGET", key, field); err != nil {
		log.Error("conn.Send HGET (%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	count, err = redis.Int64(conn.Receive())
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			count = 0
			log.Info("conn.Receive(HGET, %d) error(%v)", mid, err)
		} else {
			log.Error("conn.Receive(HGET, %d) error(%v)", mid, err)
		}
		return
	}
	return
}

// SameAssCnt get add same user count
func (d *Dao) SameAssCnt(c context.Context, mid, assistMid int64) (count int64, err error) {
	var (
		conn  = d.redis.Get(c)
		key   = d.assSameKey(mid)
		field = strconv.FormatInt(assistMid, 10)
	)
	defer conn.Close()
	if err = conn.Send("HGET", key, field); err != nil {
		log.Error("conn.Send HGET (%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	count, err = redis.Int64(conn.Receive())
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			count = 0
			log.Info("conn.Receive(HGET, %d) error(%v)", mid, err)
		} else {
			log.Error("conn.Receive(HGET, %d) error(%v)", mid, err)
		}
		return
	}
	return
}

// IncrAssCnt called when assis added, incr both same and total count
func (d *Dao) IncrAssCnt(c context.Context, mid, assistMid int64) (err error) {
	var (
		conn      = d.redis.Get(c)
		keyAll    = d.assTotalKey(mid)
		fieldAll  = strconv.FormatInt(mid, 10)
		keySame   = d.assSameKey(mid)
		fieldSame = strconv.FormatInt(assistMid, 10)
	)
	defer conn.Close()
	_ = conn.Send("EXPIRE", keyAll, d.redisExpire)
	if err = conn.Send("HINCRBY", keyAll, fieldAll, 1); err != nil {
		log.Error("conn.Send HINCRBY (%s) error(%v)", keyAll, err)
		return
	}
	_ = conn.Send("EXPIRE", keySame, d.redisExpire)
	if err = conn.Send("HINCRBY", keySame, fieldSame, 1); err != nil {
		log.Error("conn.Send HINCRBY (%s) error(%v)", keySame, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	return
}

// DelAssUpAllCache del AssUpAllCache when delete assist relation
func (d *Dao) DelAssUpAllCache(c context.Context, assistMid int64) (err error) {
	var (
		key  = d.assUpKey(assistMid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(ZERM, %s, %d) error(%v)", key, assistMid, err)
		return
	}
	return
}

// AddAssUpAllCache add AssUpAllCache when add assist relation
func (d *Dao) AddAssUpAllCache(c context.Context, assistMid int64, ups map[int64]*assist.Up) (err error) {
	var (
		key  = d.assUpKey(assistMid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	for _, up := range ups {
		if err = conn.Send("ZADD", key, up.CTime, up.Mid); err != nil {
			log.Error("conn.Send(ZADD, %s, %v, %d) error(%v)", key, up.CTime, up.Mid, err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, -1); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, 0, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(ups); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AssUpCacheWithScore get uppers passed mids from cache with scores
func (d *Dao) AssUpCacheWithScore(c context.Context, assistMid int64, start, end int64) (mids []int64, ups map[int64]*assist.Up, total int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	ups = make(map[int64]*assist.Up, end-start)
	key := d.assUpKey(assistMid)
	if err = conn.Send("ZREVRANGE", key, start, end, "WITHSCORES"); err != nil {
		log.Error("conn.Send(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	upsWithScore, err := redis.Int64s(conn.Receive())
	if err != nil {
		log.Error("conn.Do(GET, %d, %d, %d) error(%v)", assistMid, start, end, err)
	}
	for i := 0; i < len(upsWithScore); i += 2 {
		mids = append(mids, upsWithScore[i])
		var t xtime.Time
		t.Scan(strconv.FormatInt(upsWithScore[i+1], 10))
		up := &assist.Up{
			Mid:   upsWithScore[i],
			CTime: t,
		}
		ups[upsWithScore[i]] = up
	}
	if err = conn.Send("ZLEXCOUNT", key, "-", "+"); err != nil {
		log.Error("conn.Send(ZLEXCOUNT, %s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	total, err = redis.Int64(conn.Receive())
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Info("conn.Receive(ZLEXCOUNT, %d) error(%v)", assistMid, err)
		} else {
			log.Error("conn.Receive(ZLEXCOUNT, %d) error(%v)", assistMid, err)
		}
		return
	}
	return
}
