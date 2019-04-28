package dao

import (
	"context"
	"fmt"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_vipInfo = "big:%d"
	_ui      = "ui:%s"
)

func userinfo(name string) string {
	return fmt.Sprintf(_ui, name)
}
func vipInfoKey(mid int64) string {
	return fmt.Sprintf(_vipInfo, mid)
}

// DelVipInfoCache delete vipinfo cache.
func (d *Dao) DelVipInfoCache(c context.Context, mid int64) (err error) {
	return d.delCache(c, vipInfoKey(mid))
}

func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		err = errors.WithStack(err)
		d.errProm.Incr("conn_del")
		return
	}
	return
}

// DelSelCode .
func (d *Dao) DelSelCode(c context.Context, username string) (err error) {
	return d.delCache(c, userinfo(username))
}

// SetSelCode .
func (d *Dao) SetSelCode(c context.Context, username string, linkmap map[int64]int64) (err error) {
	var (
		key = userinfo(username)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: key, Object: linkmap, Expiration: d.mcExpire, Flags: gmc.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s)", key)
		d.errProm.Incr("set_selcode")
	}
	return
}

// GetSelCode .
func (d *Dao) GetSelCode(c context.Context, username string) (linkmap map[int64]int64, err error) {
	var (
		key = userinfo(username)
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
	linkmap = make(map[int64]int64)
	if err = conn.Scan(item, &linkmap); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%s)", key)
	}
	return
}

// PingMC ping memcache.
func (d *Dao) PingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&gmc.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}); err != nil {
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
	}
	return
}
