package dao

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/xstr"

	"go-common/library/sync/errgroup"
)

const (
	_bulkSize = 100
)

func statsKey(businessID, messageID int64) string {
	return fmt.Sprintf("m_%d_b_%d", messageID, businessID)
}

func recoverStatsValue(c context.Context, s string) (res *model.Stats) {
	var (
		vs  []int64
		err error
	)
	res = new(model.Stats)
	if s == "" {
		return
	}
	if vs, err = xstr.SplitInts(s); err != nil || len(vs) < 2 {
		PromError("mc:stats解析")
		log.Error("dao.recoverStatsValue(%s) err: %v", s, err)
		return
	}
	res = &model.Stats{Likes: vs[0], Dislikes: vs[1]}
	return
}

// AddStatsCache .
func (d *Dao) AddStatsCache(c context.Context, businessID int64, vs ...*model.Stats) (err error) {
	if len(vs) == 0 {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, v := range vs {
		if v == nil {
			continue
		}
		key := statsKey(businessID, v.ID)
		bs := xstr.JoinInts([]int64{v.Likes, v.Dislikes})
		item := memcache.Item{Key: key, Value: []byte(bs), Expiration: d.mcStatsExpire}
		if err = conn.Set(&item); err != nil {
			PromError("mc:增加计数缓存")
			log.Error("conn.Set(%s) error(%v)", key, err)
			return
		}
	}
	return
}

// DelStatsCache del stats cache
func (d *Dao) DelStatsCache(c context.Context, businessID int64, messageID int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := statsKey(businessID, messageID)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		PromError("mc:DelStatsCache")
		log.Error("d.DelStatsCache(%s) error(%+v)", key, err)
	}
	return
}

// AddStatsCacheMap .
func (d *Dao) AddStatsCacheMap(c context.Context, businessID int64, stats map[int64]*model.Stats) (err error) {
	var s []*model.Stats
	for _, v := range stats {
		s = append(s, v)
	}
	return d.AddStatsCache(c, businessID, s...)
}

// StatsCache .
func (d *Dao) StatsCache(c context.Context, businessID int64, messageIDs []int64) (cached map[int64]*model.Stats, missed []int64, err error) {
	if len(messageIDs) == 0 {
		return
	}
	cached = make(map[int64]*model.Stats, len(messageIDs))
	allKeys := make([]string, 0, len(messageIDs))
	midmap := make(map[string]int64, len(messageIDs))
	for _, id := range messageIDs {
		k := statsKey(businessID, id)
		allKeys = append(allKeys, k)
		midmap[k] = id
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
				PromError("mc:获取计数缓存")
				log.Error("conn.Gets(%v) error(%v)", keys, err)
				err = nil
				return
			}
			for _, reply := range replys {
				var s string
				if err = conn.Scan(reply, &s); err != nil {
					PromError("获取计数缓存json解析")
					log.Error("json.Unmarshal(%v) error(%v)", reply.Value, err)
					err = nil
					continue
				}
				stat := recoverStatsValue(c, s)
				stat.ID = midmap[reply.Key]
				mutex.Lock()
				cached[midmap[reply.Key]] = stat
				delete(midmap, reply.Key)
				mutex.Unlock()
			}
			return
		})
	}
	group.Wait()
	missed = make([]int64, 0, len(midmap))
	for _, aid := range midmap {
		missed = append(missed, aid)
	}
	return
}
