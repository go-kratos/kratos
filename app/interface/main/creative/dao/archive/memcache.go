package archive

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strconv"

	arcmdl "go-common/app/interface/main/creative/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix               = "porder_"
	_cmPrefix             = "arccm_"
	_addMidAndTitlePrefix = "add_midtitle_"
)

func limitMidSameTitle(mid int64, title string) string {
	ms := md5.Sum([]byte(title))
	return _addMidAndTitlePrefix + strconv.FormatInt(mid, 10) + "_" + hex.EncodeToString(ms[:])
}

func keyPorder(aid int64) string {
	return _prefix + strconv.FormatInt(aid, 10)
}

func keyArcCM(aid int64) string {
	return _cmPrefix + strconv.FormatInt(aid, 10)
}

// POrderCache get stat cache.
func (d *Dao) POrderCache(c context.Context, aid int64) (st *arcmdl.Porder, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyPorder(aid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get2(%d) error(%v)", aid, err)
		}
		return
	}
	if err = conn.Scan(r, &st); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		st = nil
	}
	return
}

// AddPOrderCache add stat cache.
func (d *Dao) AddPOrderCache(c context.Context, aid int64, st *arcmdl.Porder) (err error) {
	var (
		key = keyPorder(aid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// ArcCMCache get stat cache.
func (d *Dao) ArcCMCache(c context.Context, aid int64) (st *arcmdl.Commercial, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyArcCM(aid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get2(%d) error(%v)", aid, err)
		}
		return
	}
	if err = conn.Scan(r, &st); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		st = nil
	}
	return
}

// AddArcCMCache add stat cache.
func (d *Dao) AddArcCMCache(c context.Context, aid int64, st *arcmdl.Commercial) (err error) {
	var (
		key = keyArcCM(aid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// DelSubmitCache func
func (d *Dao) DelSubmitCache(c context.Context, mid int64, title string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(limitMidSameTitle(mid, title)); err == memcache.ErrNotFound {
		err = nil
	}
	return
}
