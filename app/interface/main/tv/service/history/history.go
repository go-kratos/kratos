package history

import (
	"context"
	hismodel "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/tv/model/history"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_videoPGC = "pgc"
	_videoUGC = "archive"
	_typePGC  = 1
	_typeUGC  = 2
)

// pick history from cursor and cache, then compare to tell whether we could use cache or not
func (s *Service) cacheHis(c context.Context, mid int64) (resp *history.RespCacheHis, err error) {
	var (
		cfg   = s.conf.Cfg.HisCfg
		hismc *history.HisMC
	)
	resp = &history.RespCacheHis{
		UseCache: true,
	}
	if resp.Res, err = s.dao.Cursor(c, mid, 0, cfg.Pagesize, 0, cfg.Businesses); err != nil {
		log.Error("history dao.Cursor Mid %d, Err %v", mid, err)
		return
	}
	if len(resp.Res) == 0 {
		log.Info("Mid %d, No history", mid)
		return
	}
	if hismc, err = s.dao.HisCache(c, mid); err != nil {
		log.Error("history dao.HisCache Mid %d, Err %v", mid, err)
		return
	}
	if hismc != nil { // if the first item in cache and from cursor is the same, return with cache
		if resp.Res[0].Unix == hismc.LastViewAt {
			resp.Filtered = hismc.Res
			return
		}
	}
	resp.UseCache = false
	return
}

func (s *Service) combineHis(c context.Context, req *history.ReqCombineHis) (filtered []*history.HisRes) {
	var (
		durs           = make(map[int64]int64)
		pgcRes, ugcRes []*hismodel.Resource
		pgcMap, ugcMap map[int64]*history.HisRes
	)
	g, _ := errgroup.WithContext(c)
	for _, v := range req.OriRes { // combine pgc & ugc data
		if v.Business == _videoPGC { // combine pgc history data
			if _, ok := req.OkSids[v.Sid]; !ok {
				continue
			}
			pgcRes = append(pgcRes, v)
		} else if v.Business == _videoUGC { // combine ugc history data
			if _, ok := req.OkAids[v.Oid]; !ok {
				continue
			}
			ugcRes = append(ugcRes, v)
		} else {
			continue
		}
	}
	okRes := mergeRes(pgcRes, ugcRes)
	g.Go(func() (err error) { // get pgc info
		pgcMap, err = s.pgcHisRes(context.Background(), pgcRes)
		return
	})
	g.Go(func() (err error) { // get ugc info
		ugcMap, err = s.ugcHisRes(context.Background(), ugcRes)
		return
	})
	g.Go(func() (err error) { // get duration info
		durs = s.getDuration(context.Background(), okRes)
		return nil
	})
	if err := g.Wait(); err != nil { // wait history combine media info
		log.Error("getHistory For Mid %d, Err %v", req.Mid, err)
	}
	for _, v := range okRes {
		var resrc *history.HisRes
		if v.Business == _videoPGC {
			if res, ok := pgcMap[v.Sid]; ok {
				resrc = res
			}
		} else if v.Business == _videoUGC {
			if res, ok := ugcMap[v.Oid]; ok {
				resrc = res
			}
		}
		if resrc == nil {
			log.Error("okRes Business %s, CID %d, %d, Empty", v.Business, v.Sid, v.Oid)
			continue
		}
		if dur, ok := durs[v.Oid]; ok { // duration
			resrc.PageDuration = dur
		}
		filtered = append(filtered, resrc)
	}
	return
}

// GetHistory picks history from rpc and combine the media data from Cache & DB
func (s *Service) GetHistory(c context.Context, mid int64) (filtered []*history.HisRes, err error) {
	var respCache *history.RespCacheHis
	if respCache, err = s.cacheHis(c, mid); err != nil {
		return
	}
	if respCache.UseCache {
		return respCache.Filtered, nil
	}
	okSids, okAids := s.filterIDs(c, mid, respCache.Res)
	filtered = s.combineHis(c, &history.ReqCombineHis{
		Mid:    mid,
		OkAids: okAids,
		OkSids: okSids,
		OriRes: respCache.Res,
	})
	s.dao.SaveHisCache(c, filtered)
	log.Info("Mid %d, OriLen %d, Filtered %d", mid, len(respCache.Res), len(filtered))
	return
}

// filterIDs picks the original history resource, arrange them into pgc and ugc and then filter by DAO
func (s *Service) filterIDs(ctx context.Context, mid int64, res []*hismodel.Resource) (okSids, okAids map[int64]int) {
	var ugcAIDs, pgcSIDs []int64
	for _, v := range res { // we pick only pgc & archive from History rpc
		if v.Business == _videoPGC {
			pgcSIDs = append(pgcSIDs, v.Sid)
		} else if v.Business == _videoUGC {
			ugcAIDs = append(ugcAIDs, v.Oid)
		}
	}
	okSids, okAids = s.cmsDao.MixedFilter(ctx, pgcSIDs, ugcAIDs) // we filter the okSids and okAids
	log.Info("Mid %d, okSids %v, okAids %v", mid, okSids, okAids)
	return
}

// mergeRes merges two slices and return a new slice
func mergeRes(s1 []*hismodel.Resource, s2 []*hismodel.Resource) []*hismodel.Resource {
	slice := make([]*hismodel.Resource, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
