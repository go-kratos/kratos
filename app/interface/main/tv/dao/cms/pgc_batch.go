package cms

import (
	"context"

	"fmt"
	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_mcSnCMSKey = "sn_cms_%d"
	_mcEPCMSKey = "ep_cms_%d"
)

func snCMSCacheKey(sid int64) string {
	return fmt.Sprintf(_mcSnCMSKey, sid)
}

func epCMSCacheKey(epid int64) string {
	return fmt.Sprintf(_mcEPCMSKey, epid)
}

func keysTreat(ids []int64, keyFunc func(int64) string) (idmap map[string]int64, allKeys []string) {
	idmap = make(map[string]int64, len(ids))
	for _, id := range ids {
		k := keyFunc(id)
		allKeys = append(allKeys, k)
		idmap[k] = id
	}
	return
}

func missedTreat(idmap map[string]int64, lenCached int) (missed []int64) {
	missed = make([]int64, 0, len(idmap))
	for _, id := range idmap {
		missed = append(missed, id)
	}
	missedCount.Add("tv-meta", int64(len(missed)))
	cachedCount.Add("tv-meta", int64(lenCached))
	return
}

// SeasonsMetaCache season cms meta cache
func (d *Dao) SeasonsMetaCache(c context.Context, ids []int64) (cached map[int64]*model.SeasonCMS, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.SeasonCMS, len(ids))
	idmap, allKeys := keysTreat(ids, snCMSCacheKey)
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		PromError("mc:获取Season信息缓存")
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		err = nil
		return
	}
	for key, item := range replys {
		art := &model.SeasonCMS{}
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

// EpMetaCache season cms meta cache
func (d *Dao) EpMetaCache(c context.Context, ids []int64) (cached map[int64]*model.EpCMS, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.EpCMS, len(ids))
	allKeys := make([]string, 0, len(ids))
	idmap := make(map[string]int64, len(ids))
	for _, id := range ids {
		k := epCMSCacheKey(id)
		allKeys = append(allKeys, k)
		idmap[k] = id
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		PromError("mc:获取EP信息缓存")
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		err = nil
		return
	}
	for key, item := range replys {
		art := &model.EpCMS{}
		if err = conn.Scan(item, art); err != nil {
			PromError("mc:获取EP信息缓存json解析")
			log.Error("item.Scan(%s) error(%v)", item.Value, err)
			err = nil
			continue
		}
		cached[idmap[key]] = art
		delete(idmap, key)
	}
	missed = make([]int64, 0, len(idmap))
	for _, id := range idmap {
		missed = append(missed, int64(id))
	}
	missedCount.Add("tv-meta", int64(len(missed)))
	cachedCount.Add("tv-meta", int64(len(cached)))
	return
}

//AddSeasonMetaCache add season meta cache
func (d *Dao) AddSeasonMetaCache(c context.Context, vs ...*model.SeasonCMS) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		item := &memcache.Item{Key: snCMSCacheKey(v.SeasonID), Object: v, Flags: memcache.FlagJSON, Expiration: d.expireCMS}
		if err = conn.Set(item); err != nil {
			PromError("mc:增加Season信息缓存")
			log.Error("conn.Store(%s) error(%v)", snCMSCacheKey(v.SeasonID), err)
			return
		}
	}
	return
}

//AddEpMetaCache add ep meta cache
func (d *Dao) AddEpMetaCache(c context.Context, vs ...*model.EpCMS) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		item := &memcache.Item{Key: epCMSCacheKey(v.EPID), Object: v, Flags: memcache.FlagJSON, Expiration: d.expireCMS}
		if err = conn.Set(item); err != nil {
			PromError("mc:增加EP信息缓存")
			log.Error("conn.Store(%s) error(%v)", epCMSCacheKey(v.EPID), err)
			return
		}
	}
	return
}
