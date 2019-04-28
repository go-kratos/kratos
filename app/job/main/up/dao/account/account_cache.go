package account

import (
	"context"
	"sync"
	"time"

	account "go-common/app/service/main/account/model"
	"go-common/library/log"
)

var (
	infoCache          = make(map[int64]*account.Info)
	nextClearCacheTime time.Time
	lock               = sync.Mutex{}
)

const (
	cacheClearInterval = 60 * time.Minute
)

//GetCachedInfos get cache info
func (d *Dao) GetCachedInfos(c context.Context, mids []int64, ip string) (infos map[int64]*account.Info, err error) {
	d.checkClearCache()
	var needFindMids []int64
	infos = make(map[int64]*account.Info)
	lock.Lock()
	for _, v := range mids {
		var info, ok = infoCache[v]
		if !ok {
			needFindMids = append(needFindMids, v)
			continue
		}
		log.Info("hit cache! mid=%d", v)
		infos[v] = info
	}
	lock.Unlock()

	if len(needFindMids) == 0 {
		return
	}

	var findinfo, e = d.Infos(c, needFindMids, ip)
	err = e
	if e != nil {
		log.Error("try get uid info fail, err=%v", e)
		return
	}

	lock.Lock()
	for k, v := range findinfo {
		infos[k] = v
		infoCache[k] = v
	}
	lock.Unlock()

	return
}

func (d *Dao) checkClearCache() {
	var now = time.Now()
	if now.Before(nextClearCacheTime) {
		return
	}
	d.clearCache()
}

func (d *Dao) clearCache() {
	nextClearCacheTime = time.Now().Add(cacheClearInterval)
	if len(infoCache) == 0 {
		return
	}
	lock.Lock()
	infoCache = make(map[int64]*account.Info)
	lock.Unlock()
}
