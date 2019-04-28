package history

import (
	"context"
	"fmt"

	hismodel "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/tv/model/history"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"time"

	"github.com/pkg/errors"
)

func keyHis(mid int64) string {
	return fmt.Sprintf("tv_his_%d", mid)
}

// Cursor get history rpc data
func (d *Dao) Cursor(c context.Context, mid, max int64, ps int, tp int8, businesses []string) (res []*hismodel.Resource, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &hismodel.ArgCursor{Mid: mid, Max: max, Ps: ps, RealIP: ip, TP: tp, ViewAt: max, Businesses: businesses}
	if res, err = d.hisRPC.HistoryCursor(c, arg); err != nil {
		err = errors.Wrapf(err, "d.historyRPC.HistoryCursor(%+v)", arg)
	}
	return
}

// HisCache get history cms cache.
func (d *Dao) HisCache(c context.Context, mid int64) (s *history.HisMC, err error) {
	var (
		key  = keyHis(mid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			missedCount.Add("tv-his", 1)
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &s); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	cachedCount.Add("tv-his", 1)
	return
}

// SaveHisCache save the member's history into cache
func (d *Dao) SaveHisCache(ctx context.Context, filtered []*history.HisRes) {
	var hismc = &history.HisMC{}
	if len(filtered) == 0 {
		hismc.LastViewAt = time.Now().Unix()
	} else {
		firstItem := filtered[0]
		hismc.LastViewAt = firstItem.Unix
		hismc.MID = firstItem.Mid
	}
	hismc.Res = filtered
	d.addHisCache(ctx, hismc)
}

// addHisCache adds the history into cache
func (d *Dao) addHisCache(ctx context.Context, his *history.HisMC) {
	d.addCache(func() {
		d.setHisCache(ctx, his)
	})
}

// setHisCache add his cache
func (d *Dao) setHisCache(c context.Context, his *history.HisMC) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{Key: keyHis(his.MID), Object: his, Flags: memcache.FlagJSON, Expiration: d.expireHis}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Store(%s) error(%v)", keyHis(his.MID), err)
	}
	log.Info("set HisMC Mid %d", his.MID)
	return
}
