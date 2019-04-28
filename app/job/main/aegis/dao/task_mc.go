package dao

import (
	"context"
	"fmt"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

// IsConsumerOn .
func (d *Dao) IsConsumerOn(c context.Context, bizid, flowid int, uid int64) (isOn bool, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := mcKey(bizid, flowid, uid)
	if _, err = conn.Get(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("IsConsumerOn error(%v)", err)
		}
		return
	}
	isOn = true
	return
}

func mcKey(bizid, flowid int, uid int64) string {
	return fmt.Sprintf("aegis%d_%d_%d", bizid, flowid, uid)
}
