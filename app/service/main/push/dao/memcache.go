package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go-common/app/service/main/push/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_prefixReport = "r_%d"
	_prefixToken  = "t_%s"

	_bulkSize = 20
)

func reportKey(mid int64) string {
	return fmt.Sprintf(_prefixReport, mid)
}

func tokenKey(token string) string {
	return fmt.Sprintf(_prefixToken, token)
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcReportExpire}
	err = conn.Set(&item)
	return
}

// TokensCache get reports cache by tokens
func (d *Dao) TokensCache(ctx context.Context, tokens []string) (res map[string]*model.Report, missed []string, err error) {
	res = make(map[string]*model.Report)
	if len(tokens) == 0 {
		return
	}
	var (
		mutex    sync.Mutex
		allKeys  []string
		tokenMap = make(map[string]string, len(tokens))
	)
	for _, t := range tokens {
		k := tokenKey(t)
		allKeys = append(allKeys, k)
		tokenMap[k] = t
	}
	keysLen := len(allKeys)
	group, errCtx := errgroup.WithContext(ctx)
	for i := 0; i < keysLen; i += _bulkSize {
		var keys []string
		if (i + _bulkSize) > keysLen {
			keys = allKeys[i:]
		} else {
			keys = allKeys[i : i+_bulkSize]
		}
		group.Go(func() (err error) {
			conn := d.mc.Get(errCtx)
			defer conn.Close()
			replys, err := conn.GetMulti(keys)
			if err != nil {
				PromError("mc:TokensCache GetMulti")
				log.Error("conn.Gets(%v) error(%+v)", keys, err)
				err = nil
				return
			}
			for key, item := range replys {
				r := &model.Report{}
				if err = conn.Scan(item, &r); err != nil {
					PromError("mc:TokensCache Scan")
					log.Error("item.Scan(%s) error(%+v)", item.Value, err)
					err = nil
					continue
				}
				mutex.Lock()
				res[tokenMap[key]] = r
				delete(tokenMap, key)
				mutex.Unlock()
			}
			return
		})
	}
	group.Wait()
	for _, t := range tokenMap {
		missed = append(missed, t)
	}
	return
}

// ReportsCacheByMid gets reports cache by mid.
func (d *Dao) ReportsCacheByMid(c context.Context, mid int64) (res []*model.Report, err error) {
	key := reportKey(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			res = nil
			err = nil
			missedCount.Add("report", 1)
			return
		}
		PromError("mc:获取上报")
		log.Error("conn.Get(%v) error(%v)", key, err)
		return
	}
	rs := make(map[int64]map[string]*model.Report)
	if err = conn.Scan(reply, &rs); err != nil {
		PromError("mc:解析上报")
		log.Error("item.Scan(%s) error(%v)", reply.Value, err)
		return
	}
	for _, v := range rs {
		for _, r := range v {
			res = append(res, r)
		}
	}
	cachedCount.Add("report", 1)
	return
}

// ReportsCacheByMids get report cache by mids.
func (d *Dao) ReportsCacheByMids(c context.Context, mids []int64) (res map[int64][]*model.Report, missed []int64, err error) {
	res = make(map[int64][]*model.Report, len(mids))
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
				rm := make(map[int64]map[string]*model.Report)
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
	missedCount.Add("report", int64(len(missed)))
	cachedCount.Add("report", int64(len(res)))
	return
}

// AddReportsCacheByMids add report cache by mids.
func (d *Dao) AddReportsCacheByMids(c context.Context, mrs map[int64][]*model.Report) (err error) {
	var (
		bs   []byte
		urs  map[int64]map[string]*model.Report
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	for mid, rs := range mrs {
		if urs, err = formatReport(rs); err != nil {
			log.Error("d.AddReportsCacheByMids() formatReport() data(%v) error(%v)", rs, err)
			continue
		}
		if bs, err = json.Marshal(urs); err != nil {
			PromError("增加上报缓存json解析")
			log.Error("json.Marshal(%v) error(%v)", mrs, err)
			continue
		}
		k := reportKey(mid)
		item := &memcache.Item{Key: k, Value: bs, Expiration: d.mcReportExpire}
		if err = conn.Set(item); err != nil {
			PromError("mc:批量添加上报")
			log.Error("conn.Store(%s) error(%v)", k, err)
			return
		}
	}
	return
}

func formatReport(rs []*model.Report) (mrs map[int64]map[string]*model.Report, err error) {
	mrs = make(map[int64]map[string]*model.Report)
	for _, r := range rs {
		if _, ok := mrs[r.APPID]; !ok {
			mrs[r.APPID] = make(map[string]*model.Report)
		}
		mrs[r.APPID][r.DeviceToken] = r
	}
	return
}

// // AddReportCache add report cache.
// func (d *Dao) AddReportCache(c context.Context, r *model.Report) (err error) {
// 	conn := d.mc.Get(c)
// 	defer conn.Close()
// 	k := reportKey(r.Mid)
// 	reply, err := conn.Get(k)
// 	if err != nil {
// 		return
// 	}
// 	rs := make(map[int64]map[string]*model.Report)
// 	if err = conn.Scan(reply, &rs); err != nil {
// 		PromError("mc:解析上报")
// 		log.Error("reply.Scan(%s) error(%v)", reply.Value, err)
// 		return
// 	}
// 	if _, ok := rs[r.APPID]; !ok {
// 		rs[r.APPID] = make(map[string]*model.Report)
// 	}
// 	rs[r.APPID][r.DeviceToken] = r
// 	var bs []byte
// 	if bs, err = json.Marshal(rs); err != nil {
// 		PromError("增加上报缓存json解析")
// 		log.Error("json.Marshal(%v) error(%v)", rs, err)
// 		return
// 	}
// 	item := &memcache.Item{Key: k, Value: bs, Expiration: d.mcReportExpire}
// 	if err = conn.Set(item); err != nil {
// 		PromError("mc:添加上报")
// 		log.Error("conn.Store(%s) error(%v)", k, err)
// 		return
// 	}
// 	PromInfo("mc:新增上报缓存")
// 	return
// }

// DelReportCache delete report cache.
func (d *Dao) DelReportCache(c context.Context, mid int64, appID int64, deviceToken string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	k := reportKey(mid)
	reply, err := conn.Get(k)
	if err != nil {
		if err == memcache.ErrNotFound {
			missedCount.Incr("report")
			err = nil
			return
		}
		PromError("mc:删除上报")
		log.Error("conn.Get(%v) error(%v)", k, err)
		return
	}
	rs := make(map[int64]map[string]*model.Report)
	if err = conn.Scan(reply, &rs); err != nil {
		PromError("mc:解析上报")
		log.Error("reply.Scan(%s) error(%v)", reply.Value, err)
		return
	}
	if _, ok := rs[appID]; !ok {
		return
	}
	if rs[appID][deviceToken] == nil {
		return
	}
	delete(rs[appID], deviceToken)
	var bs []byte
	if bs, err = json.Marshal(rs); err != nil {
		PromError("删除上报缓存 marshal")
		log.Error("json.Marshal(%v) error(%v)", rs, err)
		return
	}
	item := &memcache.Item{Key: k, Value: bs, Expiration: d.mcReportExpire}
	if err = conn.Set(item); err != nil {
		PromError("mc:删除上报缓存")
		log.Error("conn.Store(%s) error(%v)", k, err)
		return
	}
	PromInfo("mc:删除上报缓存")
	return
}
