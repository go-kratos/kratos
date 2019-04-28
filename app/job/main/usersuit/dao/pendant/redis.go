package pendant

import (
	"context"
	"strconv"

	"go-common/library/log"
)

const (
	_pendantPKG   = "pkg_" // key of
	_pendantEquip = "pe_"
)

// DelPKGCache del package cache
func (d *Dao) DelPKGCache(c context.Context, mid int64) (err error) {
	key := _pendantPKG + strconv.FormatInt(mid, 10)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	return
}

// DelEquipCache set pendant info cache
func (d *Dao) DelEquipCache(c context.Context, mid int64) (err error) {
	key := _pendantEquip + strconv.FormatInt(mid, 10)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	return
}
