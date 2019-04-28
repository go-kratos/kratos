package dao

import (
	"context"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"github.com/bluele/gcache"
	"github.com/pkg/errors"
)

func (d *Dao) loadStat(ctx context.Context, mid int64) (*model.Stat, error) {
	stat, err := d.statCache(ctx, mid)
	if err != nil {
		return nil, err
	}
	d.storeStat(mid, stat)
	return stat, nil
}

func (d *Dao) storeStat(mid int64, stat *model.Stat) {
	if stat == nil || stat.Follower < int64(d.c.StatCache.LeastFollower) {
		return
	}
	d.statStore.SetWithExpire(mid, stat, time.Duration(d.c.StatCache.Expire))
}

func (d *Dao) localStat(mid int64) (*model.Stat, error) {
	prom.CacheHit.Incr("local_stat_cache")
	item, err := d.statStore.Get(mid)
	if err != nil {
		prom.CacheMiss.Incr("local_stat_cache")
		return nil, err
	}
	stat, ok := item.(*model.Stat)
	if !ok {
		prom.CacheMiss.Incr("local_stat_cache")
		return nil, errors.New("Not a stat")
	}
	return stat, nil
}

// StatCache get stat cache.
func (d *Dao) StatCache(c context.Context, mid int64) (*model.Stat, error) {
	stat, err := d.localStat(mid)
	if err != nil {
		if err != gcache.KeyNotFoundError {
			log.Error("Failed to get stat from local: mid: %d: %+v", mid, err)
		}
		return d.loadStat(c, mid)
	}
	return stat, nil
}

// StatsCache get multi stat cache.
func (d *Dao) StatsCache(c context.Context, mids []int64) (map[int64]*model.Stat, []int64, error) {
	stats := make(map[int64]*model.Stat, len(mids))
	lcMissed := make([]int64, 0, len(mids))
	for _, mid := range mids {
		stat, err := d.localStat(mid)
		if err != nil {
			if err != gcache.KeyNotFoundError {
				log.Error("Failed to get stat from local: mid: %d: %+v", mid, err)
			}
			lcMissed = append(lcMissed, mid)
			continue
		}
		stats[stat.Mid] = stat
	}
	if len(lcMissed) == 0 {
		return stats, nil, nil
	}

	mcStats, mcMissed, err := d.statsCache(c, lcMissed)
	if err != nil {
		return nil, nil, err
	}
	for _, stat := range mcStats {
		d.storeStat(stat.Mid, stat)
		stats[stat.Mid] = stat
	}
	return stats, mcMissed, nil
}
