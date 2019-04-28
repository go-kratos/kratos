package account

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_addMidAndTitlePrefix = "add_midtitle_"
	_addMidHalfMinPrefix  = "add_midhafmin_"
)

func limitMidHafMin(mid int64) string {
	return _addMidHalfMinPrefix + strconv.FormatInt(mid, 10)
}

func limitMidSameTitle(mid int64, title string) string {
	ms := md5.Sum([]byte(title))
	return _addMidAndTitlePrefix + strconv.FormatInt(mid, 10) + "_" + hex.EncodeToString(ms[:])
}

// HalfMin fn
func (d *Dao) HalfMin(c context.Context, mid int64) (exist bool, ts uint64, err error) {
	var (
		conn = d.mc.Get(c)
		rp   *memcache.Item
	)
	defer conn.Close()
	key := limitMidHafMin(mid)
	rp, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get error(%v) | key(%s) mid(%d)", err, key, mid)
		}
		return
	}
	if err = conn.Scan(rp, &ts); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
		return
	}
	log.Info("HalfMin key(%s) ts(%d)", key, ts)
	if ts != 0 {
		exist = true
	}
	return
}

// AddHalfMin fn
func (d *Dao) AddHalfMin(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := limitMidHafMin(mid)
	ts := time.Now().Unix()
	if err = conn.Set(&memcache.Item{Key: key, Object: ts, Flags: memcache.FlagJSON, Expiration: d.mcLimitAddBasicExp}); err != nil {
		log.Error("memcache.set error(%v) | key(%s) mid(%d)", err, key, mid)
	}
	return
}

// DelHalfMin func
func (d *Dao) DelHalfMin(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(limitMidHafMin(mid)); err == memcache.ErrNotFound {
		err = nil
	}
	return
}

// SubmitCache get user submit cache.
func (d *Dao) SubmitCache(c context.Context, mid int64, title string) (exist int8, err error) {
	var (
		conn = d.mc.Get(c)
		rp   *memcache.Item
	)
	defer conn.Close()
	key := limitMidSameTitle(mid, title)
	rp, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get error(%v) | key(%s) mid(%d) title(%s)", err, key, mid, title)
		}
	}
	if rp != nil {
		exist = int8(binary.BigEndian.Uint64(rp.Value))
	}
	return
}

// AddSubmitCache add submit cache into mc.
func (d *Dao) AddSubmitCache(c context.Context, mid int64, title string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := limitMidSameTitle(mid, title)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, 1)
	if err = conn.Set(&memcache.Item{Key: key, Object: bs, Flags: memcache.FlagJSON, Expiration: d.mcSubExp}); err != nil {
		log.Error("memcache.set error(%v) | key(%s) mid(%d) title(%s)", err, key, mid, title)
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

func (d *Dao) pingMemcache(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("mc.ping.Store error(%v)", err)
		return
	}
	return
}
