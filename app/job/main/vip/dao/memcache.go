package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/job/main/vip/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_vipInfo   = "vo:%d"
	_vipMadel  = "madel:%d"
	_vipbuy    = "vipbuy:%d"
	_vipfrozen = "vipfrozen:%d"

	madelExpired = 3600 * 24 * 15

	vipbuyExpired = 3600 * 24 * 8

	vipFrozenExpired = 60 * 30

	vipFrozenFlag = 1

	_prefixInfo = "i:"
)

func vipfrozen(mid int64) string {
	return fmt.Sprintf(_vipfrozen, mid)
}

func vipbuy(mid int64) string {
	return fmt.Sprintf(_vipbuy, mid)
}
func vipInfoKey(mid int64) string {
	return fmt.Sprintf(_vipInfo, mid)
}

func vipMadel(mid int64) string {
	return fmt.Sprintf(_vipMadel, mid)
}

func keyInfo(mid int64) string {
	return _prefixInfo + strconv.FormatInt(mid, 10)
}

//SetVipFrozen .
func (d *Dao) SetVipFrozen(c context.Context, mid int64) (err error) {
	var (
		key = vipfrozen(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: vipFrozenFlag, Expiration: vipFrozenExpired, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "d.setVipFrozen")
		d.errProm.Incr("vip frozen_mc")
		return
	}
	return
}

//DelVipFrozen .
func (d *Dao) DelVipFrozen(c context.Context, mid int64) (err error) {
	return d.delCache(c, vipfrozen(mid))
}

//SetVipMadelCache set vip madel cache
func (d *Dao) SetVipMadelCache(c context.Context, mid int64, val int64) (err error) {
	var (
		key = vipMadel(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: val, Expiration: madelExpired, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, "d.SetVipMadelCache")
		d.errProm.Incr("vipmadel_mc")
	}
	return
}

//GetVipBuyCache get vipbuy cache by key
func (d *Dao) GetVipBuyCache(c context.Context, mid int64) (val int64, err error) {
	var (
		key  = vipbuy(mid)
		item *memcache.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("vipinfo_mc")
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("vipMadelCache_mc")
	}
	return
}

//SetVipBuyCache set vipbuy cache
func (d *Dao) SetVipBuyCache(c context.Context, mid int64, val int64) (err error) {
	var (
		key = vipbuy(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: val, Expiration: vipbuyExpired, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, "d.SetVipBuyCache")
		d.errProm.Incr("vipbuy_mc")
	}
	return
}

//GetVipMadelCache get madel info by mid
func (d *Dao) GetVipMadelCache(c context.Context, mid int64) (val int64, err error) {
	var (
		key  = vipMadel(mid)
		item *memcache.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("vipinfo_mc")
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("vipMadelCache_mc")
	}
	return
}

// SetVipInfoCache set vip info cache.
func (d *Dao) SetVipInfoCache(c context.Context, mid int64, v *model.VipInfo) (err error) {
	var (
		key = vipInfoKey(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: key, Object: v, Expiration: d.mcExpire, Flags: memcache.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		err = errors.Wrap(err, "d.SetVipInfo")
		d.errProm.Incr("vipinfo_mc")
	}
	return
}

// DelVipInfoCache del vip info cache.
func (d *Dao) DelVipInfoCache(c context.Context, mid int64) (err error) {
	err = d.delCache(c, vipInfoKey(mid))
	return
}

// DelInfoCache del vip info cache.
func (d *Dao) DelInfoCache(c context.Context, mid int64) (err error) {
	if err = d.delCache(c, keyInfo(mid)); err != nil {
		log.Error("del vipinfo cache(mid:%d) error(%+v)", mid, err)
	}
	return
}

func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Delete(%s)", key)
			d.errProm.Incr("del_mc")
		}
	}
	return

}

// AddTransferLock add lock.
func (d *Dao) AddTransferLock(c context.Context, key string) (succeed bool) {
	var (
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: 3600,
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
