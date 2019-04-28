package archive

import (
	"context"
	"strconv"

	"go-common/app/service/main/up/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_prefixUpCnt = "uc_"
	_prefixUpPas = "up_"
)

func upCntKey(mid int64) string {
	return _prefixUpCnt + strconv.FormatInt(mid, 10)
}

func upPasKey(mid int64) string {
	return _prefixUpPas + strconv.FormatInt(mid, 10)
}

// AddUpperCountCache the count of up's archives
func (d *Dao) AddUpperCountCache(c context.Context, mid int64, count int64) (err error) {
	var (
		key        = upCntKey(mid)
		conn       = d.upRds.Get(c)
		expireTime = d.upExpire
	)
	defer conn.Close()
	if count == 0 {
		expireTime = 600
	}
	if _, err = conn.Do("SETEX", key, expireTime, count); err != nil {
		log.Error("conn.Do(SETEX, %s, %d, %d)", key, expireTime, count)
		return
	}
	return
}

// UpperCountCache get up count from cache.
func (d *Dao) UpperCountCache(c context.Context, mid int64) (count int64, err error) {
	var (
		key  = upCntKey(mid)
		conn = d.upRds.Get(c)
	)
	defer conn.Close()
	if count, err = redis.Int64(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			count = -1
			err = nil
		} else {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
	}
	return
}

// UppersCountCache return uppers count cache
func (d *Dao) UppersCountCache(c context.Context, mids []int64) (cached map[int64]int64, missed []int64, err error) {
	conn := d.upRds.Get(c)
	defer conn.Close()
	cached = make(map[int64]int64)
	for _, mid := range mids {
		key := upCntKey(mid)
		if err = conn.Send("GET", key); err != nil {
			missed = mids
			continue
		}
	}
	if err = conn.Flush(); err != nil {
		missed = mids
		return
	}
	for _, mid := range mids {
		var cnt int64
		if cnt, err = redis.Int64(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				missed = append(missed, mid)
				err = nil
				continue
			}
		}
		cached[mid] = int64(cnt)
	}
	return
}

// UpperPassedCache get upper passed archives from cache.
func (d *Dao) UpperPassedCache(c context.Context, mid int64, start, end int) (aids []int64, err error) {
	var (
		key  = upPasKey(mid)
		conn = d.upRds.Get(c)
	)
	defer conn.Close()
	if aids, err = redis.Int64s(conn.Do("ZREVRANGE", key, start, end)); err != nil {
		d.errProm.Incr("upper_redis")
		log.Error("conn.Do(ZRANGE, %s, 0, -1) error(%v)", key, err)
	}
	return
}

// UppersPassedCacheWithScore get uppers passed archive from cache with score
func (d *Dao) UppersPassedCacheWithScore(c context.Context, mids []int64, start, end int) (aidm map[int64][]*model.AidPubTime, err error) {
	conn := d.upRds.Get(c)
	defer conn.Close()
	aidm = make(map[int64][]*model.AidPubTime, len(mids))
	for _, mid := range mids {
		key := upPasKey(mid)
		if err = conn.Send("ZREVRANGE", key, start, end, "WITHSCORES"); err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Send(ZREVRANGE, %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for _, mid := range mids {
		aidScores, err := redis.Int64s(conn.Receive())
		if err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Do(GET, %d) error(%v)", mid, err)
			continue
		}
		for i := 0; i < len(aidScores); i += 2 {
			var (
				score     int64
				ptime     int64
				copyright int8
			)
			score = aidScores[i+1]
			if score > 1000000000 {
				ptime = score >> 2
				copyright = int8(score & 3)
				aidm[mid] = append(aidm[mid], &model.AidPubTime{Aid: aidScores[i], PubDate: time.Time(ptime), Copyright: copyright})
			} else {
				aidm[mid] = append(aidm[mid], &model.AidPubTime{Aid: aidScores[i], PubDate: time.Time(score)})
			}
		}
	}
	return
}

// UppersPassedCache get uppers passed archives from cache.
func (d *Dao) UppersPassedCache(c context.Context, mids []int64, start, end int) (aidm map[int64][]int64, err error) {
	conn := d.upRds.Get(c)
	defer conn.Close()
	aidm = make(map[int64][]int64, len(mids))
	for _, mid := range mids {
		key := upPasKey(mid)
		if err = conn.Send("ZREVRANGE", key, start, end); err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Send(ZREVRANGE, %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for _, mid := range mids {
		aids, err := redis.Int64s(conn.Receive())
		if err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Do(GET, %d) error(%v)", mid, err)
			continue
		}
		aidm[mid] = aids
	}
	return
}

// ExpireUpperPassedCache expire up passed cache.
func (d *Dao) ExpireUpperPassedCache(c context.Context, mid int64) (ok bool, err error) {
	var (
		key  = upPasKey(mid)
		conn = d.upRds.Get(c)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.upExpire)); err != nil {
		d.errProm.Incr("upper_redis")
		log.Error("conn.Do(EXPIRE, %s, %d) error(%v)", key, d.upExpire, err)
	}
	return
}

// ExpireUppersCountCache expire ups count cache
func (d *Dao) ExpireUppersCountCache(c context.Context, mids []int64) (cachedUp, missed []int64, err error) {
	var conn = d.upRds.Get(c)
	defer conn.Close()
	defer func() {
		if err != nil {
			d.errProm.Incr("upper_redis")
		}
	}()
	for _, mid := range mids {
		var key = upCntKey(mid)
		if err = conn.Send("GET", key); err != nil {
			log.Error("conn.Send(GET, %s) error(%v)", key, err)
			return
		}
	}
	for _, mid := range mids {
		var key = upCntKey(mid)
		if err = conn.Send("EXPIRE", key, d.upExpire); err != nil {
			log.Error("conn.Send(GET, %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	cachedUp = make([]int64, 0)
	missed = make([]int64, 0)
	for _, mid := range mids {
		var cnt int
		if cnt, err = redis.Int(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
				missed = append(missed, mid)
			} else {
				log.Error("conn.Receive error(%v)", err)
				return
			}
		} else if cnt > 0 {
			cachedUp = append(cachedUp, mid)
		}
	}
	for _, mid := range mids {
		if _, err = redis.Bool(conn.Receive()); err != nil {
			log.Error("conn.Receive mid(%d) error(%v)", mid, err)
			return
		}
	}
	return
}

// ExpireUppersPassedCache expire uppers passed cache.
func (d *Dao) ExpireUppersPassedCache(c context.Context, mids []int64) (res map[int64]bool, err error) {
	conn := d.upRds.Get(c)
	defer conn.Close()
	res = make(map[int64]bool, len(mids))
	for _, mid := range mids {
		key := upPasKey(mid)
		if err = conn.Send("EXPIRE", key, d.upExpire); err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Send(%s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	var ok bool
	for _, mid := range mids {
		if ok, err = redis.Bool(conn.Receive()); err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Receive() error(%v)", err)
			return
		}
		res[mid] = ok
	}
	return
}

// AddUpperPassedCache add up paassed cache.
func (d *Dao) AddUpperPassedCache(c context.Context, mid int64, aids []int64, ptimes []time.Time, copyrights []int8) (err error) {
	var (
		key  = upPasKey(mid)
		conn = d.upRds.Get(c)
	)
	defer conn.Close()
	for k, aid := range aids {
		score := int64(ptimes[k]<<2) | int64(copyrights[k])
		if err = conn.Send("ZADD", key, score, aid); err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)", key, aid, ptimes[k], err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		d.errProm.Incr("upper_redis")
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(aids); i++ {
		if _, err = conn.Receive(); err != nil {
			d.errProm.Incr("upper_redis")
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// DelUpperPassedCache delete up passed cache.
func (d *Dao) DelUpperPassedCache(c context.Context, mid, aid int64) (err error) {
	var (
		key  = upPasKey(mid)
		conn = d.upRds.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, aid); err != nil {
		d.errProm.Incr("upper_redis")
		log.Error("conn.Do(ZERM, %s, %d) error(%v)", key, aid, err)
	}
	return
}
