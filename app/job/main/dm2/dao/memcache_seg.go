package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_keySegMC = "sg_%d_%d_%d_%d"
)

func keySegMC(tp int32, oid, total, num int64) string {
	return fmt.Sprintf(_keySegMC, tp, oid, total, num)
}

func keyXMLSeg(tp int32, oid, cnt, num int64) string {
	return fmt.Sprintf("%d_%d_%d_%d", tp, oid, cnt, num)
}

// DelXMLSegCache delete segment xml content.
func (d *Dao) DelXMLSegCache(c context.Context, tp int32, oid, cnt, num int64) (err error) {
	conn := d.mc.Get(c)
	key := keyXMLSeg(tp, oid, cnt, num)
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

// SetXMLSegCache set dm xml content into memcache.
func (d *Dao) SetXMLSegCache(c context.Context, tp int32, oid, cnt, num int64, value []byte) (err error) {
	key := keyXMLSeg(tp, oid, cnt, num)
	conn := d.mc.Get(c)
	item := memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: d.mcExpire,
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
		Expiration: d.mcExpire,
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
