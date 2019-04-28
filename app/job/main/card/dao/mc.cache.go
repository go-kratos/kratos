package dao

import (
	"context"
	"fmt"
	"strconv"

	mc "go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_prequip = "e_"
)

func equipKey(mid int64) string {
	return _prequip + strconv.FormatInt(mid, 10)
}

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	err = conn.Set(&mc.Item{
		Key:   "ping",
		Value: []byte("pong"),
	})
	return
}

// DelCacheEquip delete data from mc
func (d *Dao) DelCacheEquip(c context.Context, id int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := equipKey(id)
	if err = conn.Delete(key); err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelCacheEquip")
		log.Errorv(c, log.KV("DelCacheEquip", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
