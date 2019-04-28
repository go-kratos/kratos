package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/vip/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_vipInfo   = "vo:%d"
	_open      = "open:%d"
	_pointTip  = "pointTip:%d:%d"
	_sign      = "sign:%d"
	_vipfrozen = "vipfrozen:%d"
)

func pointTip(mid, id int64) string {
	return fmt.Sprintf(_pointTip, mid, id)
}
func openCode(mid int64) string {
	return fmt.Sprintf(_open, mid)
}
func vipInfoKey(mid int64) string {
	return fmt.Sprintf(_vipInfo, mid)
}
func signVip(mid int64) string {
	return fmt.Sprintf(_sign, mid)
}
func vipfrozen(mid int64) string {
	return fmt.Sprintf(_vipfrozen, mid)
}

//GetVipFrozen set vip frozen.
func (d *Dao) GetVipFrozen(c context.Context, mid int64) (val int, err error) {
	var (
		key = vipfrozen(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("vipfrozen_mc")
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//DelVipFrozen del vip frozen.
func (d *Dao) DelVipFrozen(c context.Context, mid int64) (err error) {
	return d.DelCache(c, vipfrozen(mid))
}

// DelVipInfoCache delete vipinfo cache.
func (d *Dao) DelVipInfoCache(c context.Context, mid int64) (err error) {
	return d.DelCache(c, vipInfoKey(mid))
}

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&gmc.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}); err != nil {
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
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
	item := &gmc.Item{Key: key, Object: v, Expiration: d.mcExpire, Flags: gmc.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, " conn.Set(%s)", key)
		d.errProm.Incr("vipinfo_mc")
	}
	return
}

// VipInfoCache get vip info.
func (d *Dao) VipInfoCache(c context.Context, mid int64) (v *model.VipInfo, err error) {
	var (
		key = vipInfoKey(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("vipinfo_mc")
		return
	}
	v = new(model.VipInfo)
	if err = conn.Scan(item, v); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%s)", key)
	}
	return
}

//DelCache del cache.
func (d *Dao) DelCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
			d.errProm.Incr("del_mc")
		}
	}
	return
}

// GetOpenCodeCount get open code count.
func (d *Dao) GetOpenCodeCount(c context.Context, mid int64) (val int, err error) {
	var (
		key = openCode(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("opencode_mc")
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//SetOpenCode set open code.
func (d *Dao) SetOpenCode(c context.Context, mid int64, count int) (err error) {
	var (
		key = openCode(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: count, Expiration: d.mcExpire, Flags: gmc.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//GetPointTip get pointTip.
func (d *Dao) GetPointTip(c context.Context, mid, id int64) (val int, err error) {
	var (
		key = pointTip(mid, id)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("opencode_mc")
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		return
	}
	return
}

//SetPointTip set point tip.
func (d *Dao) SetPointTip(c context.Context, mid, id int64, val int, expired int32) (err error) {
	var (
		key = pointTip(mid, id)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: val, Expiration: expired, Flags: gmc.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//SetSignVip .
func (d *Dao) SetSignVip(c context.Context, mid int64, t int) (err error) {
	var (
		key = signVip(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: 1, Expiration: int32(t), Flags: gmc.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//GetSignVip .
func (d *Dao) GetSignVip(c context.Context, mid int64) (val int, err error) {
	var (
		key = signVip(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		d.errProm.Incr("signvip_mc")
		return
	}
	if err = conn.Scan(item, &val); err != nil {

		return
	}
	return
}

// AddTransferLock add lock.
func (d *Dao) AddTransferLock(c context.Context, key string) (succeed bool) {
	var (
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Value:      []byte("0"),
		Expiration: 3600,
	}
	if err := conn.Add(item); err != nil {
		if err != gmc.ErrNotStored {
			log.Error("conn.Add(%s) error(%v)", key, err)
		}
	} else {
		succeed = true
	}
	return
}
