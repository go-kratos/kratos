package cms

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

// ArcMetaCache get arc cms cache.
func (d *Dao) ArcMetaCache(c context.Context, aid int64) (s *model.ArcCMS, err error) {
	var (
		key  = d.ArcCMSCacheKey(aid)
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

// SetArcMetaCache save model.ArcCMS to memcache
func (d *Dao) SetArcMetaCache(c context.Context, s *model.ArcCMS) (err error) {
	var (
		key  = d.ArcCMSCacheKey(s.AID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.expireCMS}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// AddArcMetaCache add view relates
func (d *Dao) AddArcMetaCache(arc *model.ArcCMS) {
	d.addCache(func() {
		d.SetArcMetaCache(context.TODO(), arc)
	})
}

// VideoMetaCache get video cms cache.
func (d *Dao) VideoMetaCache(c context.Context, cid int64) (s *model.VideoCMS, err error) {
	var (
		key  = d.VideoCMSCacheKey(cid)
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

// SetVideoMetaCache save model.VideoCMS to memcache
func (d *Dao) SetVideoMetaCache(c context.Context, s *model.VideoCMS) (err error) {
	var (
		key  = d.VideoCMSCacheKey(s.CID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.expireCMS}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// AddVideoMetaCache add view relates
func (d *Dao) AddVideoMetaCache(video *model.VideoCMS) {
	d.addCache(func() {
		d.SetVideoMetaCache(context.TODO(), video)
	})
}
