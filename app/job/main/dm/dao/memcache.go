package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/job/main/dm/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixSub   = "s_"
	_keyDuration = "d_" // video duration
)

func keySubject(tp int32, oid int64) string {
	return _prefixSub + fmt.Sprintf("%d_%d", tp, oid)
}

// keyDuration return video duration key.
func keyDuration(oid int64) string {
	return _keyDuration + strconv.FormatInt(oid, 10)
}

// SubjectCache get subject from memcache.
func (d *Dao) SubjectCache(c context.Context, tp int32, oid int64) (sub *model.Subject, err error) {
	var (
		conn = d.mc.Get(c)
		key  = keySubject(tp, oid)
		rp   *memcache.Item
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	sub = &model.Subject{}
	if err = conn.Scan(rp, &sub); err != nil {
		log.Error("conn.Scan() error(%v)", err)
	}
	return
}

// SetSubjectCache add subject cache.
func (d *Dao) SetSubjectCache(c context.Context, sub *model.Subject) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = keySubject(sub.Type, sub.Oid)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     sub,
		Flags:      memcache.FlagJSON,
		Expiration: d.mcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DurationCache return duration of video.
func (d *Dao) DurationCache(c context.Context, oid int64) (duration int64, err error) {
	var (
		key  = keyDuration(oid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			duration = model.NotFound
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if duration, err = strconv.ParseInt(string(item.Value), 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", item.Value, err)
	}
	return
}

// SetDurationCache set video duration to redis.
func (d *Dao) SetDurationCache(c context.Context, oid, duration int64) (err error) {
	key := keyDuration(oid)
	conn := d.mc.Get(c)
	item := memcache.Item{
		Key:        key,
		Value:      []byte(fmt.Sprint(duration)),
		Expiration: d.mcExpire,
		Flags:      memcache.FlagRAW,
	}
	if err = conn.Set(&item); err != nil {
		log.Error("mc.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}
