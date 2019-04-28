package dao

import (
	"context"
	"fmt"

	gmc "go-common/library/cache/memcache"
)

const (
	_prefixRedPointFlag = "r_p_f_%d"
)

func redPointFlagKey(mid int64) string {
	return fmt.Sprintf(_prefixRedPointFlag, mid)
}

// PendantPointCache get new pendant info red point cache.
func (d *Dao) PendantPointCache(c context.Context, mid int64) (pid int64, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(redPointFlagKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	err = conn.Scan(item, &pid)
	return
}

// SetPendantPointCache set red point cache.
func (d *Dao) SetPendantPointCache(c context.Context, mid, pid int64) (err error) {
	var (
		item = &gmc.Item{Key: redPointFlagKey(mid), Object: pid, Expiration: d.pointExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// DelPendantPointCache delete new pendant info red point cache.
func (d *Dao) DelPendantPointCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(redPointFlagKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
	}
	return
}
