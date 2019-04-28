package dao

import (
	"context"
	"fmt"
	"sync"
	"time"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_prefixReport = "r_%d"

	_bulkSize = 10
)

func reportKey(mid int64) string {
	return fmt.Sprintf(_prefixReport, mid)
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: int32(time.Now().Unix())}
	err = conn.Set(&item)
	return
}

// ReportsCacheByMids get report cache by mids.
func (d *Dao) ReportsCacheByMids(c context.Context, mids []int64) (res map[int64][]*pushmdl.Report, missed []int64, err error) {
	res = make(map[int64][]*pushmdl.Report, len(mids))
	if len(mids) == 0 {
		return
	}
	allKeys := make([]string, 0, len(mids))
	midmap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		k := reportKey(mid)
		allKeys = append(allKeys, k)
		midmap[k] = mid
	}
	group := errgroup.Group{}
	mutex := sync.Mutex{}
	keysLen := len(allKeys)
	for i := 0; i < keysLen; i += _bulkSize {
		var keys []string
		if (i + _bulkSize) > keysLen {
			keys = allKeys[i:]
		} else {
			keys = allKeys[i : i+_bulkSize]
		}
		group.Go(func() error {
			conn := d.mc.Get(context.TODO())
			defer conn.Close()
			replys, err := conn.GetMulti(keys)
			if err != nil {
				PromError("mc:获取上报")
				log.Error("conn.Gets(%v) error(%v)", keys, err)
				return nil
			}
			for k, item := range replys {
				rm := make(map[int64]map[string]*pushmdl.Report)
				if err = conn.Scan(item, &rm); err != nil {
					PromError("mc:解析上报")
					log.Error("item.Scan(%s) error(%v)", item.Value, err)
					continue
				}
				mutex.Lock()
				mid := midmap[k]
				for _, v := range rm {
					for _, r := range v {
						res[mid] = append(res[mid], r)
					}
				}
				delete(midmap, k)
				mutex.Unlock()
			}
			return nil
		})
	}
	group.Wait()
	missed = make([]int64, 0, len(midmap))
	for _, mid := range midmap {
		missed = append(missed, mid)
	}
	return
}
