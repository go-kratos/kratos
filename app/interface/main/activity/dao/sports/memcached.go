package sports

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_preQq = "q_"
)

func keyQq(tp int64) string {
	return fmt.Sprintf("%s%d", _preQq, tp)
}

// QqCache get qq from cache
func (d *Dao) QqCache(c context.Context, tp int64) (rs *json.RawMessage, err error) {
	var (
		mckey = keyQq(tp)
		conn  = d.mc.Get(c)
		item  *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(mckey); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get error(%v)", err)
		}
		return
	}
	if err = conn.Scan(item, &rs); err != nil {
		log.Error("item.Scan error(%v)", err)
	}
	return
}

// SetQqCache set qq to cache
func (d *Dao) SetQqCache(c context.Context, v *json.RawMessage, tp int64) (err error) {
	var (
		conn  = d.mc.Get(c)
		mckey = keyQq(tp)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: mckey, Object: v, Flags: memcache.FlagJSON, Expiration: d.mcQqExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}
