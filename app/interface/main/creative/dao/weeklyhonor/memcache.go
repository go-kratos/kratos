package weeklyhonor

import (
	"context"
	"fmt"

	model "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_statKey  = "hs_%d_%s"
	_honorKey = "ho_%d_%s"
	_clickKey = "hc_%d"
)

func statKey(mid int64, date string) string {
	return fmt.Sprintf(_statKey, mid, date)
}

func honorKey(mid int64, date string) string {
	return fmt.Sprintf(_honorKey, mid, date)
}

func honorClickKey(mid int64) string {
	return fmt.Sprint(_clickKey, mid)
}

// StatMC get stat cache.
func (d *Dao) StatMC(c context.Context, mid int64, date string) (hs *model.HonorStat, err error) {
	var (
		key  = statKey(mid, date)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(r, &hs); err != nil {
		log.Error("conn.Scan(%s) error(%v)", r.Value, err)
		hs = nil
	}
	return
}

// HonorMC get honor cache.
func (d *Dao) HonorMC(c context.Context, mid int64, date string) (res *model.HonorLog, err error) {
	var (
		key  = honorKey(mid, date)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("conn.Scan(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// SetStatMC add stat cache.
func (d *Dao) SetStatMC(c context.Context, mid int64, date string, hs *model.HonorStat) (err error) {
	var (
		key  = statKey(mid, date)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: hs, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%s) error(%v)", key, err)
	}
	return
}

// SetHonorMC add honor cache.
func (d *Dao) SetHonorMC(c context.Context, mid int64, date string, hs *model.HonorLog) (err error) {
	var (
		key  = honorKey(mid, date)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: hs, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%s) error(%v)", key, err)
	}
	return
}

// SetClickMC add click cache
func (d *Dao) SetClickMC(c context.Context, mid int64) (err error) {
	key := honorClickKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Value: []byte{1}, Expiration: d.mcClickExpire}); err != nil {
		log.Error("memcache.Set(%s) error(%v)", key, err)
	}
	return
}

// ClickMC get click cache
func (d *Dao) ClickMC(c context.Context, mid int64) (err error) {
	key := honorClickKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	_, err = conn.Get(key)
	return
}
