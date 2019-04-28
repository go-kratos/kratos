package dao

import (
	"context"
	"strconv"
	"sync"

	"go-common/app/service/main/archive/api"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/cache/memcache"

	"go-common/library/sync/errgroup"
)

const (
	_prefixArc     = "ap_"
	_prefixBangumi = "bp_"
	_bulkSize      = 100
)

func arcKey(aid int64) string {
	return _prefixArc + strconv.FormatInt(aid, 10)
}

func bangumiKey(bid int64) string {
	return _prefixBangumi + strconv.FormatInt(bid, 10)
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// AddArchivesCache batch set archives cache.
func (d *Dao) AddArchivesCache(c context.Context, vs ...*api.Arc) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		item := memcache.Item{Key: arcKey(v.Aid), Object: v, Flags: memcache.FlagProtobuf, Expiration: d.mcExpire}
		if err = conn.Set(&item); err != nil {
			PromError("mc:增加稿件缓存", "conn.Store(%s) error(%v)", arcKey(v.Aid), err)
			return
		}
	}
	return
}

// AddArchivesCacheMap batch set archives cache.
func (d *Dao) AddArchivesCacheMap(c context.Context, arcm map[int64]*api.Arc) (err error) {
	var arcs []*api.Arc
	for _, arc := range arcm {
		arcs = append(arcs, arc)
	}
	return d.AddArchivesCache(c, arcs...)
}

// ArchivesCache batch get archive from cache.
func (d *Dao) ArchivesCache(c context.Context, aids []int64) (cached map[int64]*api.Arc, missed []int64, err error) {
	if len(aids) == 0 {
		return
	}
	cached = make(map[int64]*api.Arc, len(aids))
	allKeys := make([]string, 0, len(aids))
	aidmap := make(map[string]int64, len(aids))
	for _, aid := range aids {
		k := arcKey(aid)
		allKeys = append(allKeys, k)
		aidmap[k] = aid
	}

	group, errCtx := errgroup.WithContext(c)
	mutex := sync.Mutex{}
	keysLen := len(allKeys)
	for i := 0; i < keysLen; i += _bulkSize {
		var keys []string
		if (i + _bulkSize) > keysLen {
			keys = allKeys[i:]
		} else {
			keys = allKeys[i : i+_bulkSize]
		}

		group.Go(func() (err error) {
			conn := d.mc.Get(errCtx)
			replys, err := conn.GetMulti(keys)
			defer conn.Close()
			if err != nil {
				PromError("mc:获取稿件缓存", "conn.Gets(%v) error(%v)", keys, err)
				err = nil
				return
			}
			for _, reply := range replys {
				arc := &api.Arc{}
				if err = conn.Scan(reply, arc); err != nil {
					PromError("获取稿件缓存json解析", "json.Unmarshal(%v) error(%v)", reply.Value, err)
					err = nil
					continue
				}
				mutex.Lock()
				cached[aidmap[reply.Key]] = arc
				delete(aidmap, reply.Key)
				mutex.Unlock()
			}
			return
		})
	}
	group.Wait()
	missed = make([]int64, 0, len(aidmap))
	for _, aid := range aidmap {
		missed = append(missed, aid)
	}
	MissedCount.Add("archive", int64(len(missed)))
	CachedCount.Add("archive", int64(len(cached)))
	return
}

// DelArchiveCache delete archive cache.
func (d *Dao) DelArchiveCache(c context.Context, aid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(arcKey(aid)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			PromError("mc:删除稿件缓存", "conn.Delete(%s) error(%v)", arcKey(aid), err)
			return
		}
	}
	return
}

// AddBangumisCacheMap batch set bangumis cache.
func (d *Dao) AddBangumisCacheMap(c context.Context, bm map[int64]*feedmdl.Bangumi) (err error) {
	var bs []*feedmdl.Bangumi
	for _, b := range bm {
		bs = append(bs, b)
	}
	return d.AddBangumisCache(c, bs...)
}

// AddBangumisCache add batch set bangumi cache.
func (d *Dao) AddBangumisCache(c context.Context, bs ...*feedmdl.Bangumi) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, b := range bs {
		if b == nil {
			continue
		}
		item := memcache.Item{Key: bangumiKey(b.SeasonID), Object: b, Flags: memcache.FlagProtobuf, Expiration: d.bangumiExpire}
		if err = conn.Set(&item); err != nil {
			PromError("mc:增加番剧缓存", "conn.Store(%s) error(%v)", bangumiKey(b.SeasonID), err)
			return
		}
	}
	return
}

// BangumisCache batch get archive from cache.
func (d *Dao) BangumisCache(c context.Context, bids []int64) (cached map[int64]*feedmdl.Bangumi, missed []int64, err error) {
	cached = make(map[int64]*feedmdl.Bangumi, len(bids))
	if len(bids) == 0 {
		return
	}
	keys := make([]string, 0, len(bids))
	bidmap := make(map[string]int64, len(bids))
	for _, bid := range bids {
		k := bangumiKey(bid)
		keys = append(keys, k)
		bidmap[k] = bid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(keys)
	if err != nil {
		PromError("mc:获取番剧", "conn.Gets(%v) error(%v)", keys, err)
		return
	}
	for _, reply := range replys {
		b := &feedmdl.Bangumi{}
		if err = conn.Scan(reply, b); err != nil {
			PromError("获取番剧json解析", "json.Unmarshal(%v) error(%v)", reply.Value, err)
			return
		}
		cached[bidmap[reply.Key]] = b
		delete(bidmap, reply.Key)
	}
	missed = make([]int64, 0, len(bidmap))
	for _, bid := range bidmap {
		missed = append(missed, bid)
	}
	MissedCount.Add("bangumi", int64(len(missed)))
	CachedCount.Add("bangumi", int64(len(cached)))
	return
}
