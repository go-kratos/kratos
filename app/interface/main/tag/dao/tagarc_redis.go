package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	//  newest arc of tag of all region
	_newArcKey = "ta_%d"
	// newest arcs of tag of region
	_regionArcKey = "ra_%d_%d"
	// origin newest arcs of tag of region
	_regionOriArcKey = "ro_%d_%d"
)

// arcKey -
func arcKey(tid int64) string {
	return fmt.Sprintf(_newArcKey, tid)
}

// regionArcKey
func regionArcKey(rid int32, tid int64) string {
	return fmt.Sprintf(_regionArcKey, rid, tid)
}

// regionOriArcKey
func regionOriArcKey(rid int32, tid int64) string {
	return fmt.Sprintf(_regionOriArcKey, rid, tid)
}

func (d *Dao) zrange(c context.Context, key string, start, end int) (aids []int64, count int, err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREVRANGE", key, start, end); err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if err = conn.Send("ZCARD", key); err != nil {
		log.Error("conn.Do(ZCARD) err(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() err(%v)", err)
		return
	}
	if aids, err = redis.Int64s(conn.Receive()); err != nil {
		log.Error("redis.Int64s()err(%v)", err)
		return
	}
	if count, err = redis.Int(conn.Receive()); err != nil {
		log.Error("redis.INT64 err(%v)", err)
	}
	return
}

// NewArcsCache get archives of tag by tid of all regions.
func (d *Dao) NewArcsCache(c context.Context, tid int64, start, end int) (aids []int64, count int, err error) {
	key := arcKey(tid)
	aids, count, err = d.zrange(c, key, start, end)
	return
}

// RemoveNewArcsCache remove arc of tag by tids
func (d *Dao) RemoveNewArcsCache(c context.Context, aid int64, tids ...int64) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		key := arcKey(tid)
		if err = conn.Send("ZREM", key, aid); err != nil {
			log.Error("")
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(tids); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelNewArcsCache del tid from cache
func (d *Dao) DelNewArcsCache(c context.Context, tid int64) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	key := arcKey(tid)
	_, err = conn.Do("DEL", key)
	return
}

// RemTidArcCache .
func (d *Dao) RemTidArcCache(c context.Context, aid int64, tids []int64) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		key := arcKey(tid)
		if err = conn.Send("ZREM", key, aid); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(tids); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DeleteNewArcCache delete arcs by tid from cache.
func (d *Dao) DeleteNewArcCache(c context.Context, tid int64, aids string) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	key := arcKey(tid)
	if err = conn.Send("ZREM", key, aids); err != nil {
		log.Error("conn.Receive() error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// AddNewArcsCache add new arcs of tid by tid-arcs.
func (d *Dao) AddNewArcsCache(c context.Context, tid int64, as ...*api.Arc) (err error) {
	key := arcKey(tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, a := range as {
		if err = conn.Send("ZADD", key, a.PubDate, a.Aid); err != nil {
			log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, a.Aid, err)
			return
		}
	}
	if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(d.tagArcMaxNum + 1)); err != nil {
		log.Error("conn.Send(ZREMRANGEBYRANK, %s) error(%v)", key, err)
	}
	if err = conn.Send("EXPIRE", key, d.expireNewArc); err != nil {
		log.Error("conn.Send(EXPIRE, %s) error(%v)", key, err)
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(as)+2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddNewArcCache add new arc into tids by arc-tids.
func (d *Dao) AddNewArcCache(c context.Context, a *api.Arc, tids ...int64) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	// del 0 cache
	d.RemoveNewArcsCache(c, 0, tids...)
	for _, tid := range tids {
		var ok bool
		key := arcKey(tid)
		if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireNewArc)); err != nil {
			log.Error("Conn.Do(EXPIRE) err(%v)", err)
		} else if ok {
			if err = conn.Send("ZADD", key, a.PubDate, a.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, a.Aid, err)
				continue
			}
			if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(d.tagArcMaxNum + 1)); err != nil {
				log.Error("conn.Send(ZREMRANGEBYRANK, %s) error(%v)", key, err)
				continue
			}
			if err = conn.Flush(); err != nil {
				log.Error("conn.Flush error(%v)", err)
				continue
			}
			for i := 0; i < 2; i++ {
				if _, err = conn.Receive(); err != nil {
					log.Error("conn.Receive() error(%v)", err)
					continue
				}
			}
		}
	}
	return
}

// ========RegionNewArcs========

// OriginRegionNewArcsCache get origin newest arc of tag by rid and tid
func (d *Dao) OriginRegionNewArcsCache(c context.Context, rid int32, tid int64, start, end int) (aids []int64, count int, err error) {
	key := regionOriArcKey(rid, tid)
	aids, count, err = d.zrange(c, key, start, end)
	return
}

// ExpireRegionNewArcsCache expire region new arcs.
func (d *Dao) ExpireRegionNewArcsCache(c context.Context, rid int32, tid int64) (ok bool, err error) {
	key := regionArcKey(rid, tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
	}
	return
}

// ExpireOriginalNewestArcCache expire region new arcs.
func (d *Dao) ExpireOriginalNewestArcCache(c context.Context, rid int32, tid int64) (ok bool, err error) {
	key := regionOriArcKey(rid, tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
	}
	return
}

// RegionNewArcsCache get newest arcs of tag by rid and tid.
func (d *Dao) RegionNewArcsCache(c context.Context, rid int32, tid int64, start, end int) (aids []int64, count int, err error) {
	key := regionArcKey(rid, tid)
	aids, count, err = d.zrange(c, key, start, end)
	return
}

// AddRegionNewArcCache 增加分区下热门tag的最新视频.
func (d *Dao) AddRegionNewArcCache(c context.Context, rid int32, arc *api.Arc, tids ...int64) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	var ok bool
	for _, tid := range tids {
		// del 0 cache
		d.RemoveRegionNewArcCache(c, rid, &api.Arc{Aid: 0, Copyright: int32(archive.CopyrightOriginal)}, tid)
		key := regionArcKey(rid, tid)
		if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireNewArc)); err != nil {
			log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
			return
		}
		if ok {
			if _, err = conn.Do("ZADD", key, arc.PubDate, arc.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, arc.Aid, err)
				return
			}
		}
		if arc.Copyright == int32(archive.CopyrightOriginal) {
			key := regionOriArcKey(rid, tid)
			if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expireNewArc)); err != nil {
				log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
				return
			}
			if ok {
				if _, err = conn.Do("ZADD", key, arc.PubDate, arc.Aid); err != nil {
					log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, arc.Aid, err)
					return
				}
			}
		}
	}
	return
}

// RemoveRegionNewArcCache del new arc from rid-tids by arc.
func (d *Dao) RemoveRegionNewArcCache(c context.Context, rid int32, arc *api.Arc, tids ...int64) (err error) {
	var count int
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		key := regionArcKey(rid, tid)
		if err = conn.Send("ZREM", key, arc.Aid); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
		count++
		if arc.Copyright == int32(archive.CopyrightOriginal) {
			key := regionOriArcKey(rid, tid)
			if err = conn.Send("ZREM", key, arc.Aid); err != nil {
				log.Error("conn.Receive() error(%v)", err)
				return
			}
			count++
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// RegionNewestArcCache add region newest arc.
func (d *Dao) RegionNewestArcCache(c context.Context, rid int32, tid int64, aids []int64) (exist, none []int64, err error) {
	key := regionArcKey(rid, tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	exist = make([]int64, 0)
	none = make([]int64, 0)
	for _, aid := range aids {
		if err = conn.Send("ZSCORE", key, aid); err != nil {
			log.Error("conn.Send(ZSCORE, %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for _, aid := range aids {
		var pubData int64
		if pubData, err = redis.Int64(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Error("redis.Int64 err(%v)", err)
				return
			}
		}
		if pubData == 0 {
			none = append(none, aid)
		} else {
			exist = append(exist, aid)
		}
	}
	return
}

// DeleteRegionNewArcsCache delete region NewArcs Cache.
func (d *Dao) DeleteRegionNewArcsCache(c context.Context, tid int64, rid int32, arcs []*api.Arc) (err error) {
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	var count int
	for _, arc := range arcs {
		aids := strconv.FormatInt(arc.Aid, 10)
		key := regionArcKey(arc.TypeID, tid)
		if err = conn.Send("ZREM", key, aids); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			continue
		}
		count++
		if rid != 0 && rid != arc.TypeID {
			key := regionArcKey(rid, tid)
			if err = conn.Send("ZREM", key, aids); err != nil {
				log.Error("conn.Receive() error(%v)", err)
				continue
			}
			count++
		}
		if arc.Copyright == int32(archive.CopyrightOriginal) {
			key := regionOriArcKey(arc.TypeID, tid)
			if err = conn.Send("ZREM", key, aids); err != nil {
				log.Error("conn.Receive() error(%v)", err)
				continue
			}
			count++
			if rid != 0 && rid != arc.TypeID {
				key := regionOriArcKey(rid, tid)
				if err = conn.Send("ZREM", key, aids); err != nil {
					log.Error("conn.Receive() error(%v)", err)
					continue
				}
				count++
			}
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
		}
	}
	return
}

// AddRegionNewestArcCache add region newest arc.
func (d *Dao) AddRegionNewestArcCache(c context.Context, rid int32, tid int64, as []*api.Arc) (err error) {
	key := regionArcKey(rid, tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, a := range as {
		if err = conn.Send("ZADD", key, a.PubDate, a.Aid); err != nil {
			log.Error("conn.Send(ZADD, %s, %d) error(%v)", key, a.Aid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(as); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// OriginalNewestArcCache add region newest arc.
func (d *Dao) OriginalNewestArcCache(c context.Context, rid int32, tid int64, aids []int64) (exist, none []int64, err error) {
	exist = make([]int64, 0)
	none = make([]int64, 0)
	key := regionOriArcKey(rid, tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	for _, aid := range aids {
		if err = conn.Send("ZSCORE", key, aid); err != nil {
			log.Error("conn.Send(ZSCORE, %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for _, aid := range aids {
		var pubData int64
		if pubData, err = redis.Int64(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Error("redis.Int64 err(%v)", err)
				return
			}
		}
		if pubData == 0 {
			none = append(none, aid)
		} else {
			exist = append(exist, aid)
		}
	}
	return
}

// AddOriginalNewestArcCache region original newest archive .
func (d *Dao) AddOriginalNewestArcCache(c context.Context, rid int32, tid int64, as []*api.Arc) (err error) {
	okey := regionOriArcKey(rid, tid)
	conn := d.rankRedis.Get(c)
	defer conn.Close()
	var count int
	for _, a := range as {
		if a.Copyright == int32(archive.CopyrightOriginal) {
			if err = conn.Send("ZADD", okey, a.PubDate, a.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d) error(%v)", okey, a.Aid, err)
				return
			}
			count++
		}
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
