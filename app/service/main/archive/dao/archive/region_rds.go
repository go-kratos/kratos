package archive

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/archive/api"
	commarc "go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	// all type archives
	_allTypeKey = "bl_arc"
	// all archives.
	_suffixAll = "_a"
	// origin archives copyright=1.
	_suffixOrigin = "_o"
	// top type archives
	_suffixTopAll = "_t"
	// top type expire
	topTypeExpire = 3 * 24 * time.Hour
	// all type expire
	allTypeExpire = 7 * 24 * time.Hour
)

func rgAllKey(rid int16) string {
	return fmt.Sprintf("%b", rid) + _suffixAll
}

func rgOriginKey(rid int16) string {
	return fmt.Sprintf("%b", rid) + _suffixOrigin
}

func rgTopKey(rid int16) string {
	return fmt.Sprintf("%b%s", rid, _suffixTopAll)
}

// AddRegionArcCache add archives into region cache.
func (d *Dao) AddRegionArcCache(c context.Context, rid, reid int16, as ...*api.RegionArc) (err error) {
	var (
		key   = rgAllKey(rid)
		okey  = rgOriginKey(rid)
		tpKey = rgTopKey(reid)
		count int
		conn  = d.rgRds.Get(c)
	)
	defer conn.Close()
	defer func() {
		if err != nil {
			d.errProm.Incr("regin_redis")
		}
	}()
	for _, a := range as {
		if !a.AllowShow() {
			continue
		}
		if err = conn.Send("ZADD", tpKey, a.PubDate, a.Aid); err != nil {
			log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)")
			return
		}
		count++
		if err = conn.Send("ZREMRANGEBYSCORE", tpKey, "-inf", time.Now().Add(-time.Duration(topTypeExpire)).Unix()); err != nil {
			log.Error("conn.Send(ZREMRANGEBYSCORE, %s, %d, %d) error(%v)", tpKey)
			return
		}
		count++
		if err = conn.Send("ZADD", key, a.PubDate, a.Aid); err != nil {
			log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, a.Aid, err)
			return
		}
		count++
		if a.Copyright == commarc.CopyrightOriginal {
			if err = conn.Send("ZADD", okey, a.PubDate, a.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d) error(%v)", okey, a.Aid, err)
				return
			}
			count++
		}
		// 连载动画不记录
		if rid != 33 {
			if err = conn.Send("ZADD", _allTypeKey, a.PubDate, a.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)", _allTypeKey, a.PubDate, a.Aid, err)
				return
			}
			count++
			if err = conn.Send("ZREMRANGEBYSCORE", _allTypeKey, "-inf", time.Now().Add(-time.Duration(allTypeExpire)).Unix()); err != nil {
				log.Error("conn.Send()")
				return
			}
			count++
		}
	}
	if count == 0 {
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// RegionTopArcsCache top region archives.
func (d *Dao) RegionTopArcsCache(c context.Context, reid int16, start, end int) (aids []int64, err error) {
	key := rgTopKey(reid)
	if reid == 0 {
		key = _allTypeKey
	}
	if aids, err = d.zrange(c, key, start, end); err != nil {
		d.errProm.Incr("regin_redis")
		log.Error("dao.zrange(%s, %d, %d) error(%v)", key, start, end, err)
	}
	return
}

// RegionArcsCache region archives.
func (d *Dao) RegionArcsCache(c context.Context, rid int16, start, end int) (aids []int64, err error) {
	key := rgAllKey(rid)
	aids, err = d.zrange(c, key, start, end)
	return
}

// RegionOriginArcsCache region origin archives.
func (d *Dao) RegionOriginArcsCache(c context.Context, rid int16, start, end int) (aids []int64, err error) {
	key := rgOriginKey(rid)
	aids, err = d.zrange(c, key, start, end)
	return
}

func (d *Dao) zrange(c context.Context, key string, start, end int) (aids []int64, err error) {
	conn := d.rgRds.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		d.errProm.Incr("regin_redis")
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	err = redis.ScanSlice(values, &aids)
	return
}

// RegionTopCountCache top region count of archives.
func (d *Dao) RegionTopCountCache(c context.Context, reids []int16, min, max int64) (recm map[int16]int, err error) {
	conn := d.rgRds.Get(c)
	defer conn.Close()
	for _, rid := range reids {
		key := rgTopKey(rid)
		if err = conn.Send("ZCOUNT", key, min, max); err != nil {
			d.errProm.Incr("regin_redis")
			log.Error("conn.Do(ZCOUNT, %s) error(%v)", key, err)
			break
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	recm = make(map[int16]int, len(reids))
	for _, rid := range reids {
		if recm[rid], err = redis.Int(conn.Receive()); err != nil {
			d.errProm.Incr("regin_redis")
			log.Error("conn.Receive error(%v)")
			return
		}
	}
	return
}

// RegionAllCountCache count of all type
func (d *Dao) RegionAllCountCache(c context.Context) (count int, err error) {
	conn := d.rgRds.Get(c)
	defer conn.Close()
	count, err = redis.Int(conn.Do("ZCARD", _allTypeKey))
	return
}

// RegionCountCache count of arcs.
func (d *Dao) RegionCountCache(c context.Context, rid int16) (count int, err error) {
	conn := d.rgRds.Get(c)
	defer conn.Close()
	key := rgAllKey(rid)
	count, err = redis.Int(conn.Do("ZCARD", key))
	return
}

// RegionOriginCountCache count of arcs.
func (d *Dao) RegionOriginCountCache(c context.Context, rid int16) (count int, err error) {
	conn := d.rgRds.Get(c)
	defer conn.Close()
	key := rgOriginKey(rid)
	count, err = redis.Int(conn.Do("ZCARD", key))
	return
}

// DelRegionArcCache delete from zset.
func (d *Dao) DelRegionArcCache(c context.Context, rid, reid int16, aid int64) (err error) {
	conn := d.rgRds.Get(c)
	defer conn.Close()
	defer func() {
		if err != nil {
			d.errProm.Incr("regin_redis")
		}
	}()
	var (
		rgAllKey = rgAllKey(rid)
		oKey     = rgOriginKey(rid)
		rgTopKey = rgTopKey(reid)
	)
	if err = conn.Send("ZREM", rgAllKey, aid); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("ZREM", oKey, aid); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("ZREM", rgTopKey, aid); err != nil {
		log.Error("conn.Send(ZREM, %s, %d) error(%v)", rgTopKey, aid, err)
		return
	}
	if err = conn.Send("ZREM", _allTypeKey, aid); err != nil {
		log.Error("conn.Send(ZREM, %s, %d) error(%v)", _allTypeKey, aid, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive error(%v)", err)
	}
	return
}
