package cms

import (
	"context"
	"fmt"

	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_mcEPKey     = "ep_%d"
	_mcSeasonKey = "sn_%d"
)

// SeaCacheKey .
func (d *Dao) SeaCacheKey(sid int64) string {
	return fmt.Sprintf(_mcSeasonKey, sid)
}

// EPCacheKey .
func (d *Dao) EPCacheKey(epid int64) string {
	return fmt.Sprintf(_mcEPKey, epid)
}

// GetSeasonCache get SnAuth cache.
func (d *Dao) GetSeasonCache(c context.Context, sid int64) (s *model.SnAuth, err error) {
	var (
		key  = d.SeaCacheKey(sid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &s); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	return
}

// GetEPCache get EpAuth cache.
func (d *Dao) GetEPCache(c context.Context, epid int64) (ep *model.EpAuth, err error) {
	var (
		key  = d.EPCacheKey(epid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &ep); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	return
}

// AddSnAuthCache save model.SnAuth to memcache
func (d *Dao) AddSnAuthCache(c context.Context, s *model.SnAuth) (err error) {
	var (
		key  = d.SeaCacheKey(s.ID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: s, Flags: memcache.FlagJSON, Expiration: d.expireCMS}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// AddEpAuthCache save model.EpAuth to memcache
func (d *Dao) AddEpAuthCache(c context.Context, ep *model.EpAuth) (err error) {
	var (
		key  = d.EPCacheKey(ep.EPID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: ep, Flags: memcache.FlagJSON, Expiration: d.expireCMS}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// snAuthCache season auth cache
func (d *Dao) snAuthCache(c context.Context, ids []int64) (cached map[int64]*model.SnAuth, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.SnAuth, len(ids))
	idmap, allKeys := keysTreat(ids, d.SeaCacheKey)
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		PromError("mc:获取Season信息缓存")
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		return
	}
	for key, item := range replys {
		art := &model.SnAuth{}
		if err = conn.Scan(item, art); err != nil {
			PromError("mc:获取Season信息缓存json解析")
			log.Error("item.Scan(%s) error(%v)", item.Value, err)
			err = nil
			continue
		}
		cached[idmap[key]] = art
		delete(idmap, key)
	}
	missed = missedTreat(idmap, len(cached))
	return
}

// epAuthCache ep auth cache
func (d *Dao) epAuthCache(c context.Context, ids []int64) (cached map[int64]*model.EpAuth, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.EpAuth, len(ids))
	idmap, allKeys := keysTreat(ids, d.EPCacheKey)
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		PromError("mc:获取EP信息缓存")
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		return
	}
	for key, item := range replys {
		art := &model.EpAuth{}
		if err = conn.Scan(item, art); err != nil {
			PromError("mc:获取EP信息缓存json解析")
			log.Error("item.Scan(%s) error(%v)", item.Value, err)
			err = nil
			continue
		}
		cached[idmap[key]] = art
		delete(idmap, key)
	}
	missed = missedTreat(idmap, len(cached))
	return
}

// SnAuth get's season auth info from Cache & DB
func (d *Dao) SnAuth(c context.Context, sid int64) (sn *model.SnAuth, err error) {
	if sn, err = d.GetSeasonCache(c, sid); err != nil { // mc Error
		return
	} else if sn == nil { // mc not found, go DB
		if sn, err = d.SnAuthDB(c, int64(sid)); err != nil { // DB error
			log.Error("SnAuthDB (%d) ERROR (%v)", sid, err)
			return
		}
		if sn == nil { // DB not found, build a fake item in MC to avoid checking DB next time
			log.Error("SnAuthDB (%d) not found(%v) in DB", sid)
			sn = &model.SnAuth{ID: int64(sid), Check: 0, IsDeleted: 1}
		}
		if err = d.AddSnAuthCache(c, sn); err != nil { // set item in MC ( not found - fake, or true )
			log.Error("AddSnAuthCache fail(%v)", err)
		}
	}
	return
}

// EpAuth get's ep auth info from Cache & DB
func (d *Dao) EpAuth(c context.Context, epid int64) (ep *model.EpAuth, err error) {
	if ep, err = d.GetEPCache(c, epid); err != nil { // MC error
		log.Error("GetEPCache(%d) Error(%v)", epid, err)
		return
	} else if ep == nil { // MC not found, go DB
		if ep, err = d.EpAuthDB(c, epid); err != nil { // DB error
			log.Error("DBSimpleEP(%d) ERROR (%v)", epid, err)
			return
		}
		if ep == nil { // DB not found, build a fake item in MC to avoid checking DB next time
			log.Error("EPID(%d) not found", epid)
			ep = &model.EpAuth{EPID: epid, State: 4, IsDeleted: 1}
		}
		if err = d.AddEpAuthCache(c, ep); err != nil {
			log.Error("AddEpAuthCache fail(%v)", err)
		}
	}
	return
}
