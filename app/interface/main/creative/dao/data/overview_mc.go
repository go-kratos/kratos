package data

import (
	"context"
	"strconv"

	"go-common/app/interface/main/creative/model/data"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_basePrefix         = "bse_"
	_areaPrefix         = "are_"
	_trendPrefix        = "tre_"
	_rfdPrefix          = "rfd_"
	_rfmPrefix          = "rfm_"
	_actPrefix          = "act_"
	_incrPrefix         = "incr_"
	_thirtyDayArcPrefix = "30arc_"
	_thirtyDayArtPrefix = "30art_"
)

func keyBase(mid int64, date string) string {
	return _basePrefix + date + "_" + strconv.FormatInt(mid, 10)
}

func keyArea(mid int64, date string) string {
	return _areaPrefix + date + "_" + strconv.FormatInt(mid, 10)
}

func keyTrend(mid int64, date string) string {
	return _trendPrefix + date + "_" + strconv.FormatInt(mid, 10)
}

func keyRfd(mid int64, date string) string {
	return _rfdPrefix + date + "_" + strconv.FormatInt(mid, 10)
}

func keyRfm(mid int64, date string) string {
	return _rfmPrefix + date + "_" + strconv.FormatInt(mid, 10)
}

func keyAct(mid int64, date string) string {
	return _actPrefix + date + "_" + strconv.FormatInt(mid, 10)
}

func keyViewIncr(mid int64, ty, date string) string {
	return _incrPrefix + ty + "_" + date + "_" + strconv.FormatInt(mid, 10)
}

func keyThirtyDayArchive(mid int64, ty string) string {
	return _thirtyDayArcPrefix + ty + "_" + strconv.FormatInt(mid, 10)
}

func keyThirtyDayArticle(mid int64) string {
	return _thirtyDayArtPrefix + "_" + strconv.FormatInt(mid, 10)
}

// ViewerBaseCache add ViewerBaseCache cache.
func (d *Dao) ViewerBaseCache(c context.Context, mid int64, dt string) (res map[string]*data.ViewerBase, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyBase(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddViewerBaseCache add ViewerBaseCache cache update data by week.
func (d *Dao) AddViewerBaseCache(c context.Context, mid int64, dt string, res map[string]*data.ViewerBase) (err error) {
	key := keyBase(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// ViewerAreaCache add ViewerArea cache.
func (d *Dao) ViewerAreaCache(c context.Context, mid int64, dt string) (res map[string]map[string]int64, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyArea(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddViewerAreaCache add ViewerArea cache update data by week.
func (d *Dao) AddViewerAreaCache(c context.Context, mid int64, dt string, res map[string]map[string]int64) (err error) {
	key := keyArea(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// TrendCache add trend cache.
func (d *Dao) TrendCache(c context.Context, mid int64, dt string) (res map[string]*data.ViewerTrend, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyTrend(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddTrendCache add trend cache update data by week.
func (d *Dao) AddTrendCache(c context.Context, mid int64, dt string, res map[string]*data.ViewerTrend) (err error) {
	key := keyTrend(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// RelationFansDayCache add relation day cache.
func (d *Dao) RelationFansDayCache(c context.Context, mid int64, dt string) (res map[string]map[string]int, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyRfd(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddRelationFansDayCache add relation day cache update data by day.
func (d *Dao) AddRelationFansDayCache(c context.Context, mid int64, dt string, res map[string]map[string]int) (err error) {
	key := keyRfd(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// RelationFansMonthCache add relation month cache.
func (d *Dao) RelationFansMonthCache(c context.Context, mid int64, dt string) (res map[string]map[string]int, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyRfm(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddRelationFansMonthCache add relation month cache update data by day.
func (d *Dao) AddRelationFansMonthCache(c context.Context, mid int64, dt string, res map[string]map[string]int) (err error) {
	key := keyRfm(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// ViewerActionHourCache add ActionHour cache.
func (d *Dao) ViewerActionHourCache(c context.Context, mid int64, dt string) (res map[string]*data.ViewerActionHour, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyAct(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddViewerActionHourCache add ActionHour cache update data by week.
func (d *Dao) AddViewerActionHourCache(c context.Context, mid int64, dt string, res map[string]*data.ViewerActionHour) (err error) {
	key := keyAct(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// ViewerIncrCache get ViewerIncr cache.
func (d *Dao) ViewerIncrCache(c context.Context, mid int64, ty, dt string) (res *data.ViewerIncr, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyViewIncr(mid, ty, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddViewerIncrCache add ViewerIncr cache update data by day.
func (d *Dao) AddViewerIncrCache(c context.Context, mid int64, ty, dt string, res *data.ViewerIncr) (err error) {
	key := keyViewIncr(mid, ty, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// ThirtyDayArchiveCache get archive 30 days cache.
func (d *Dao) ThirtyDayArchiveCache(c context.Context, mid int64, ty string) (res []*data.ThirtyDay, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyThirtyDayArchive(mid, ty))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddThirtyDayArchiveCache add archive 30 days cache update data by day.
func (d *Dao) AddThirtyDayArchiveCache(c context.Context, mid int64, ty string, res []*data.ThirtyDay) (err error) {
	key := keyThirtyDayArchive(mid, ty)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// ThirtyDayArticleCache get article 30 days cache.
func (d *Dao) ThirtyDayArticleCache(c context.Context, mid int64) (res []*artmdl.ThirtyDayArticle, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	// get cache
	r, err = conn.Get(keyThirtyDayArticle(mid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		res = nil
	}
	return
}

// AddThirtyDayArticleCache add article 30 days cache update data by day.
func (d *Dao) AddThirtyDayArticleCache(c context.Context, mid int64, res []*artmdl.ThirtyDayArticle) (err error) {
	key := keyThirtyDayArticle(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: res, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}
