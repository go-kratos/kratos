package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_keyDuration = "d_" // video duration
	_keySegMC    = "sg_%d_%d_%d_%d"
)

func keyXMLSeg(tp int32, oid, cnt, num int64) string {
	return fmt.Sprintf("%d_%d_%d_%d", tp, oid, cnt, num)
}

func keySegMC(tp int32, oid, total, num int64) string {
	return fmt.Sprintf(_keySegMC, tp, oid, total, num)
}

// keyDuration return video duration key.
func keyDuration(oid int64) string {
	return _keyDuration + strconv.FormatInt(oid, 10)
}

// XMLSegCache get dm segment xml content from memcache.
func (d *Dao) XMLSegCache(c context.Context, tp int32, oid, cnt, num int64) (res []byte, err error) {
	var (
		conn = d.dmMC.Get(c)
		key  = keyXMLSeg(tp, oid, cnt, num)
		rp   *memcache.Item
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			PromCacheMiss("dm_xml_seg", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("dm_xml_seg", 1)
	if err = conn.Scan(rp, &res); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// SetXMLSegCache set dm xml content into memcache.
func (d *Dao) SetXMLSegCache(c context.Context, tp int32, oid, cnt, num int64, value []byte) (err error) {
	key := keyXMLSeg(tp, oid, cnt, num)
	conn := d.dmMC.Get(c)
	item := memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: d.dmExpire,
		Flags:      memcache.FlagRAW,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("mc.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}

// DurationCache return duration of video.
func (d *Dao) DurationCache(c context.Context, oid int64) (duration int64, err error) {
	var (
		key  = keyDuration(oid)
		conn = d.dmMC.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			duration = model.NotFound
			err = nil
			PromCacheMiss("video_duration", 1)
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("video_duration", 1)
	if duration, err = strconv.ParseInt(string(item.Value), 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", item.Value, err)
	}
	return
}

// SetDurationCache set video duration to redis.
func (d *Dao) SetDurationCache(c context.Context, oid, duration int64) (err error) {
	key := keyDuration(oid)
	conn := d.dmMC.Get(c)
	item := memcache.Item{
		Key:        key,
		Value:      []byte(fmt.Sprint(duration)),
		Expiration: d.dmExpire,
		Flags:      memcache.FlagRAW,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("mc.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}

// SetDMSegCache set segment dm to cache.
func (d *Dao) SetDMSegCache(c context.Context, tp int32, oid, total, num int64, dmSeg *model.DMSeg) (err error) {
	key := keySegMC(tp, oid, total, num)
	conn := d.dmSegMC.Get(c)
	item := memcache.Item{
		Key:        key,
		Object:     dmSeg,
		Expiration: d.dmSegMCExpire,
		Flags:      memcache.FlagProtobuf | memcache.FlagGzip,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}

// DMSegCache dm segment pb cache.
func (d *Dao) DMSegCache(c context.Context, tp int32, oid, total, num int64) (dmSeg *model.DMSeg, err error) {
	var (
		key  = keySegMC(tp, oid, total, num)
		conn = d.dmSegMC.Get(c)
		item *memcache.Item
	)
	dmSeg = new(model.DMSeg)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			dmSeg = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, dmSeg); err != nil {
		log.Error("conn.Scan() error(%v)", err)
	}
	return
}
