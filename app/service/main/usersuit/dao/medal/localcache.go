package medal

import (
	"context"
	"time"

	"go-common/library/log"
	"go-common/library/stat/prom"

	"github.com/bluele/gcache"
	"github.com/pkg/errors"
)

func (d *Dao) loadMedal(c context.Context, mid int64) (int64, bool, error) {
	nid, nofound, err := d.medalActivatedCache(c, mid)
	if err != nil {
		return 0, nofound, err
	}
	d.storeMedal(mid, nid, nofound)
	return nid, nofound, nil
}

func (d *Dao) storeMedal(mid int64, nid int64, nofound bool) {
	if nofound {
		return
	}
	d.medalStore.SetWithExpire(mid, nid, time.Duration(d.c.MedalCache.Expire))
}

func (d *Dao) localMedal(mid int64) (int64, error) {
	prom.CacheHit.Incr("local_medal_cache")
	item, err := d.medalStore.Get(mid)
	if err != nil {
		prom.CacheMiss.Incr("local_medal_cache")
		return 0, err
	}
	nid, ok := item.(int64)
	if !ok {
		prom.CacheMiss.Incr("local_medal_cache")
		return 0, errors.New("Not a medal")
	}
	return nid, nil
}

// MedalActivatedCache get medal cache.
func (d *Dao) MedalActivatedCache(c context.Context, mid int64) (int64, bool, error) {
	nid, err := d.localMedal(mid)
	if err != nil {
		if err != gcache.KeyNotFoundError {
			log.Error("Failed to get medal from local: mid: %d: %+v", mid, err)
		}
		return d.loadMedal(c, mid)
	}
	return nid, false, nil
}

// MedalsActivatedCache get multi medals cache.
func (d *Dao) MedalsActivatedCache(c context.Context, mids []int64) (map[int64]int64, []int64, error) {
	nids := make(map[int64]int64, len(mids))
	lcMissed := make([]int64, 0, len(mids))
	for _, mid := range mids {
		nid, err := d.localMedal(mid)
		if err != nil {
			if err != gcache.KeyNotFoundError {
				log.Error("Failed to get medal from local: mid: %d: %+v", mid, err)
			}
			lcMissed = append(lcMissed, mid)
			continue
		}
		nids[mid] = nid
	}
	if len(lcMissed) == 0 {
		return nids, nil, nil
	}
	mcMedals, mcMissed, err := d.medalsActivatedCache(c, lcMissed)
	if err != nil {
		return nil, nil, err
	}
	for mid, nid := range mcMedals {
		d.storeMedal(mid, nid, false)
		nids[mid] = nid
	}
	return nids, mcMissed, nil
}
