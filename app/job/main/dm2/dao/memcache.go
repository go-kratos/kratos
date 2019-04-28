package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/job/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixXML   = "dm_xml_"
	_prefixSub   = "s_"
	_prefixAjax  = "dm_ajax_"
	_keyDuration = "d_" // video duration

)

func keyXML(oid int64) string {
	return _prefixXML + strconv.FormatInt(oid, 10)
}

func keySubject(tp int32, oid int64) string {
	return _prefixSub + fmt.Sprintf("%d_%d", tp, oid)
}

func keyAjax(oid int64) string {
	return _prefixAjax + strconv.FormatInt(oid, 10)
}

// keyDuration return video duration key.
func keyDuration(oid int64) string {
	return _keyDuration + strconv.FormatInt(oid, 10)
}

func keyTransferLock() string {
	return "dm_transfer_lock"
}

// DelXMLCache delete xml content.
func (d *Dao) DelXMLCache(c context.Context, oid int64) (err error) {
	conn := d.mc.Get(c)
	key := keyXML(oid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// AddXMLCache add xml content to memcache.
func (d *Dao) AddXMLCache(c context.Context, oid int64, value []byte) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keyXML(oid),
		Value:      value,
		Expiration: d.mcExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s) error(%v)", keyXML(oid), err)
	}
	return
}

// XMLCache get xml content.
func (d *Dao) XMLCache(c context.Context, oid int64) (data []byte, err error) {
	key := keyXML(oid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	data = item.Value
	return
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
			sub = nil
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	sub = &model.Subject{}
	if err = conn.Scan(rp, &sub); err != nil {
		log.Error("mc.Scan(%d) error(%v)", oid, err)
	}
	return
}

// SubjectsCache multi get subject from memcache.
func (d *Dao) SubjectsCache(c context.Context, tp int32, oids []int64) (cached map[int64]*model.Subject, missed []int64, err error) {
	var (
		conn   = d.mc.Get(c)
		keys   []string
		oidMap = make(map[string]int64, len(oids))
	)
	cached = make(map[int64]*model.Subject, len(oids))
	defer conn.Close()
	for _, oid := range oids {
		k := keySubject(tp, oid)
		if _, ok := oidMap[k]; !ok {
			keys = append(keys, k)
			oidMap[k] = oid
		}
	}
	rs, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	for k, r := range rs {
		sub := &model.Subject{}
		if err = conn.Scan(r, sub); err != nil {
			log.Error("conn.Scan(%s) error(%v)", r.Value, err)
			err = nil
			continue
		}
		cached[oidMap[k]] = sub
		// delete hit key
		delete(oidMap, k)
	}
	// missed key
	missed = make([]int64, 0, len(oidMap))
	for _, oid := range oidMap {
		missed = append(missed, oid)
	}
	return
}

// AddSubjectCache add subject cache.
func (d *Dao) AddSubjectCache(c context.Context, sub *model.Subject) (err error) {
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

// DelSubjectCache delete subject memcache cache.
func (d *Dao) DelSubjectCache(c context.Context, tp int32, oid int64) (err error) {
	conn := d.mc.Get(c)
	key := keySubject(tp, oid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// AddTransferLock 添加弹幕转移并发锁
func (d *Dao) AddTransferLock(c context.Context) (succeed bool) {
	var (
		key  = keyTransferLock()
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: 60,
	}
	if err := conn.Add(item); err != nil {
		if err != memcache.ErrNotStored {
			log.Error("conn.Add(%s) error(%v)", key, err)
		}
	} else {
		succeed = true
	}
	return
}

// DelTransferLock 删除弹幕转移并发锁
func (d *Dao) DelTransferLock(c context.Context) (err error) {
	var (
		key  = keyTransferLock()
		conn = d.mc.Get(c)
	)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// DelAjaxDMCache delete ajax dm from memcache.
func (d *Dao) DelAjaxDMCache(c context.Context, oid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := keyAjax(oid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("DelAjaxDMCache.conn.Delete(%s) error(%v)", key, err)
		}
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
