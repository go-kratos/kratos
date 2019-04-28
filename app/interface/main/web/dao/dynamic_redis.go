package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/service/main/archive/api"
	dymdl "go-common/app/service/main/dynamic/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyRegion       = "dr_"
	_keyRegionTag    = "drt_"
	_keyRegions      = "drs_"
	_keyRegionFmt    = "%d_%d_%d"
	_keyRegionTagFmt = "%d_%d_%d_%d"
)

func keyRegion(rid int32, pn, ps int) string {
	return _keyRegion + fmt.Sprintf(_keyRegionFmt, rid, pn, ps)
}

func keyRegionTag(tagID int64, rid int32, pn, ps int) string {
	return _keyRegionTag + fmt.Sprintf(_keyRegionTagFmt, tagID, rid, pn, ps)
}

// SetRegionBakCache set dynamic region data to cache.
func (d *Dao) SetRegionBakCache(c context.Context, rid int32, pn, ps int, rs *dymdl.DynamicArcs3) (err error) {
	key := keyRegion(rid, pn, ps)
	err = d.setBakCache(c, key, rs)
	return
}

// RegionBakCache get dynamic region from cache.
func (d *Dao) RegionBakCache(c context.Context, rid int32, pn, ps int) (rs *dymdl.DynamicArcs3, err error) {
	d.cacheProm.Incr("dynamic_region_remote_cache")
	key := keyRegion(rid, pn, ps)
	rs, err = d.bakCache(c, key)
	if rs == nil {
		log.Error("RegionBakCache d.bakCache(%d,%d,%d,%s) is nill", rid, pn, ps, key)
	}
	return
}

// SetRegionTagBakCache set dynamic region tag data to cache.
func (d *Dao) SetRegionTagBakCache(c context.Context, tagID int64, rid int32, pn, ps int, rs *dymdl.DynamicArcs3) (err error) {
	key := keyRegionTag(tagID, rid, pn, ps)
	err = d.setBakCache(c, key, rs)
	return
}

// RegionTagBakCache get dynamic region tag from cache.
func (d *Dao) RegionTagBakCache(c context.Context, tagID int64, rid int32, pn, ps int) (rs *dymdl.DynamicArcs3, err error) {
	d.cacheProm.Incr("dynamic_tag_remote_cache")
	key := keyRegionTag(tagID, rid, pn, ps)
	rs, err = d.bakCache(c, key)
	if rs == nil {
		log.Error("RegionTagBakCache d.bakCache(%d,%d,%d,%s) is nill", rid, pn, ps, key)
	}
	return
}

// SetRegionsBakCache set regions data to cache.
func (d *Dao) SetRegionsBakCache(c context.Context, rs map[int32][]*api.Arc) (err error) {
	key := _keyRegions
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = json.Marshal(rs); err != nil {
		log.Error("json.Marshal(%v) error(%v)", rs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisDynamicBakExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s,%d) error(%v)", key, d.redisDynamicBakExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Recevie(%d) error(%v0", i, err)
		}
	}
	return
}

// RegionsBakCache get dynamic region data from cache.
func (d *Dao) RegionsBakCache(c context.Context) (rs map[int32][]*api.Arc, err error) {
	d.cacheProm.Incr("dynamic_regions_remote_cache")
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var values []byte
	if values, err = redis.Bytes(conn.Do("GET", _keyRegions)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Error("RegionsBakCache (%s) return nil ", _keyRegions)
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", _keyRegions, err)
		}
		return
	}
	rs = make(map[int32][]*api.Arc)
	if err = json.Unmarshal(values, &rs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", values, err)
	}
	return
}

func (d *Dao) setBakCache(c context.Context, key string, rs *dymdl.DynamicArcs3) (err error) {
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = json.Marshal(rs); err != nil {
		log.Error("json.Marshal(%v) error(%v)", rs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.redisDynamicBakExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s,%d) error(%v)", key, d.redisDynamicBakExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive(%d) error(%v)", i, err)
		}
	}
	return
}

func (d *Dao) bakCache(c context.Context, key string) (rs *dymdl.DynamicArcs3, err error) {
	conn := d.redisBak.Get(c)
	defer conn.Close()
	var values []byte
	if values, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	rs = &dymdl.DynamicArcs3{}
	if err = json.Unmarshal(values, rs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", values, err)
	}
	return
}
