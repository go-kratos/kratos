package goblin

import (
	"context"
	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"
)

const _hotwordKey = "_tv_search"

// Hotword get hotword cache.
func (d *Dao) Hotword(c context.Context) (s []*model.Hotword, err error) {
	var (
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(_hotwordKey); err != nil {
		return
	}
	err = conn.Scan(item, &s)
	return
}
