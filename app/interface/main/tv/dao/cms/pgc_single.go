package cms

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

// GetSnCMSCache get SeasonCMS cache.
func (d *Dao) GetSnCMSCache(c context.Context, sid int64) (s *model.SeasonCMS, err error) {
	var (
		key  = snCMSCacheKey(sid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			missedCount.Add("tv-meta", 1)
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &s); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	cachedCount.Add("tv-meta", 1)
	return
}

// SetSnCMSCache save model.SeasonCMS to memcache
func (d *Dao) SetSnCMSCache(c context.Context, s *model.SeasonCMS) (err error) {
	var (
		key  = snCMSCacheKey(s.SeasonID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.expireCMS}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// GetEpCMSCache get EpCMS cache.
func (d *Dao) GetEpCMSCache(c context.Context, epid int64) (s *model.EpCMS, err error) {
	var (
		key  = epCMSCacheKey(epid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			missedCount.Add("tv-meta", 1)
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &s); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	cachedCount.Add("tv-meta", 1)
	return
}

// SetEpCMSCache save model.EpCMS to memcache
func (d *Dao) SetEpCMSCache(c context.Context, s *model.EpCMS) (err error) {
	var (
		key  = epCMSCacheKey(s.EPID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.expireCMS}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}
