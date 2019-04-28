package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/tag/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	//  tag new arc
	_newArcKey = "ta_%d"
	// newest arcs of tag of region
	_regionArcKey = "ra_%d_%d"
	// origin newest arcs of tag of region
	_regionOriArcKey = "ro_%d_%d"
)

// arcKey
func arcKey(tid int64) string {
	return fmt.Sprintf(_newArcKey, tid)
}

// regionArcKey
func regionArcKey(rid int16, tid int64) string {
	return fmt.Sprintf(_regionArcKey, rid, tid)
}

// regionOriArcKey
func regionOriArcKey(rid int16, tid int64) string {
	return fmt.Sprintf(_regionOriArcKey, rid, tid)
}

// AddTagNewArcCache .
func (d *Dao) AddTagNewArcCache(c context.Context, arc *model.Archive, tids ...int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	// del 0 cache
	d.RemTidArcCache(c, 0, tids...)
	var pubTime int64
	if pubTime, err = d.getPubTime(arc.Aid, arc.PubTime); err != nil {
		return
	}
	for _, tid := range tids {
		var ok bool
		if ok, err = d.expireTagNewArc(c, tid); err != nil {
			log.Error("d.expireTagNewArc(%d) error(%v)", tid, err)
			continue
		}
		if ok {
			key := arcKey(tid)
			if err = conn.Send("ZADD", key, pubTime, arc.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)", key, pubTime, arc.Aid, err)
				continue
			}
			if err = conn.Send("ZREMRANGEBYRANK", key, 0, -(d.maxNum + 1)); err != nil {
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

// RemTidArcCache .
func (d *Dao) RemTidArcCache(c context.Context, aid int64, tids ...int64) (err error) {
	conn := d.redis.Get(c)
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

func (d *Dao) expireTagNewArc(c context.Context, tid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := arcKey(tid)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
	}
	return
}

// ======== region hot tag new arc ========

// AddRegionTagNewArcCache .
func (d *Dao) AddRegionTagNewArcCache(c context.Context, arc *model.Archive, tids ...int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var pubTime int64
	if pubTime, err = d.getPubTime(arc.Aid, arc.PubTime); err != nil {
		return
	}
	for _, tid := range tids {
		// del 0 cache
		d.RemoveRegionNewArcCache(c, &model.Archive{ID: 0, Copyright: archive.CopyrightOriginal}, tid)
		key := regionArcKey(arc.TypeID, tid)
		var ok bool
		if ok, err = d.expireRegionNewArc(c, arc.TypeID, tid); err != nil {
			log.Error("d.expireRegionNewArc(%d,%d) error(%v)", arc.TypeID, tid, err)
		}
		if ok {
			if _, err = conn.Do("ZADD", key, pubTime, arc.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)", key, pubTime, arc.Aid, err)
			}
		}
		if arc.Copyright == archive.CopyrightOriginal {
			if ok, err = d.expireRegionOriArc(c, arc.TypeID, tid); err != nil {
				log.Error("d.expireRegionOriArc(%d,%d) error(%v)", arc.TypeID, tid, err)
			}
			if ok {
				key := regionOriArcKey(arc.TypeID, tid)
				if _, err = conn.Do("ZADD", key, pubTime, arc.Aid); err != nil {
					log.Error("conn.Send(ZADD, %s, %v, %d) error(%v)", key, pubTime, arc.Aid, err)
				}
			}
		}
	}
	return
}

// RemoveRegionNewArcCache .
func (d *Dao) RemoveRegionNewArcCache(c context.Context, arc *model.Archive, tids ...int64) (err error) {
	var count int
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, tid := range tids {
		key := regionArcKey(arc.TypeID, tid)
		if err = conn.Send("ZREM", key, arc.Aid); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
		count++
		if arc.Copyright == archive.CopyrightOriginal {
			key := regionOriArcKey(arc.TypeID, tid)
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

func (d *Dao) expireRegionNewArc(c context.Context, tpID int16, tid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := regionArcKey(tpID, tid)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
	}
	return
}

func (d *Dao) expireRegionOriArc(c context.Context, tpID int16, tid int64) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := regionOriArcKey(tpID, tid)
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, d.expNewArc)); err != nil {
		log.Error("conn.Do(EXPIRE, %s) error(%v)", key, err)
	}
	return
}

func (d *Dao) getPubTime(aid int64, pubDate string) (timeStamp int64, err error) {
	var t time.Time
	loc, _ := time.LoadLocation("Local")
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", pubDate, loc); err != nil {
		log.Error("time.Parse(%s) error(%v)", pubDate, err)
		return
	}
	timeStamp = t.Unix()
	return
}

// UpdateTagNewArcCache .
func (d *Dao) UpdateTagNewArcCache(c context.Context, tid int64, arcs map[int64]*api.Arc) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var count int
	for _, arc := range arcs {
		if arc.Aid == 0 {
			continue
		}
		count++
		pubTime := arc.PubDate.Time().Unix()
		if err = conn.Send("ZADD", arcKey(tid), pubTime, arc.Aid); err != nil {
			log.Error("conn.Send(ZADD, %d, %d, %d) error(%v)", tid, pubTime, arc.Aid, err)
		}
		count++
		key := regionArcKey(int16(arc.TypeID), tid)
		if err = conn.Send("ZADD", key, pubTime, arc.Aid); err != nil {
			log.Error("conn.Send(ZADD, %s, %d, %d) error(%v)", key, pubTime, arc.Aid, err)
		}
		if arc.Copyright == int32(archive.CopyrightOriginal) {
			count++
			key := regionOriArcKey(int16(arc.TypeID), tid)
			if err = conn.Send("ZADD", key, pubTime, arc.Aid); err != nil {
				log.Error("conn.Send(ZADD, %s, %v, %d) error(%v)", key, pubTime, arc.Aid, err)
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
			return
		}
	}
	return
}
