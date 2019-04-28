package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/point/model"
	gmc "go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_pointInfo = "pti:%d"
)

func pointKey(mid int64) string {
	return fmt.Sprintf(_pointInfo, mid)
}

//DelPointInfoCache .
func (d *Dao) DelPointInfoCache(c context.Context, mid int64) (err error) {
	return d.delCache(c, pointKey(mid))
}

// PointInfoCache .
func (d *Dao) PointInfoCache(c context.Context, mid int64) (pi *model.PointInfo, err error) {
	var (
		item *gmc.Item
	)
	key := pointKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			pi = nil
			return
		}
		err = errors.Wrapf(err, "d.PointInfoCache(%d)", mid)
		d.errProm.Incr("get_mc")
		return
	}
	pi = new(model.PointInfo)
	if err = conn.Scan(item, pi); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%d)", pi.Mid)
		d.errProm.Incr("scan_mc")
	}
	return
}

// SetPointInfoCache .
func (d *Dao) SetPointInfoCache(c context.Context, pi *model.PointInfo) (err error) {
	key := pointKey(pi.Mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{
		Key:        key,
		Object:     pi,
		Expiration: d.mcExpire,
		Flags:      gmc.FlagJSON,
	}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "d.SetPointInfoCache(%d)", pi.Mid)
		d.errProm.Incr("set_mc")
		return
	}
	return
}

// DelCache del cache.
func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "conn.Delete(%s)", key)
			d.errProm.Incr("del_mc")
		}
	}
	return
}
