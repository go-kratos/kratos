package dao

import (
	"context"
	"strconv"

	"go-common/app/service/main/secure/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixLocs = "locs_"
)

func locsKey(mid int64) string {
	return _prefixLocs + strconv.FormatInt(mid, 10)
}
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Set(&gmc.Item{Key: "ping", Value: []byte{1}, Expiration: d.locsExpire}); err != nil {
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
	}
	conn.Close()
	return
}

// AddLocsCache add login locs count to cache.
func (d *Dao) AddLocsCache(c context.Context, mid int64, locs *model.Locs) (err error) {
	item := &gmc.Item{Key: locsKey(mid), Object: locs, Expiration: d.locsExpire, Flags: gmc.FlagProtobuf}
	conn := d.mc.Get(c)
	if err = conn.Set(item); err != nil {
		log.Error("AddLocs err(%v)", err)
	}
	conn.Close()
	return
}

// LocsCache get login locs count.
func (d *Dao) LocsCache(c context.Context, mid int64) (locs map[int64]int64, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(locsKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	loc := &model.Locs{}
	if err = conn.Scan(item, loc); err != nil {
		log.Error("Locs err(%v)", err)
	}
	locs = loc.LocsCount
	return
}
