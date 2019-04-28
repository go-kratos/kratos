package service

import (
	"context"

	"go-common/app/interface/main/web/conf"
	"go-common/app/interface/main/web/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const _allRank = 0

// NewList get new list by region id.
func (s *Service) NewList(c context.Context, rid int32, tp int8, pn, ps int) (arcs []*api.Arc, count int, err error) {
	var (
		addCache          = true
		first             bool
		allArcs, moreArcs []*api.Arc
		start, end, max   int
		ip                = metadata.String(c, metadata.RemoteIP)
	)
	// check first or second region id.
	if _, ok := s.rids[rid]; ok {
		first = true
	}
	if first {
		max = s.c.Rule.MaxFirstCacheSize
		end = max
	} else {
		max = s.c.Rule.MaxSecondCacheSize
		start = (pn - 1) * ps
		end = start + ps - 1
		if start < 0 || end < 0 {
			err = ecode.RequestErr
			return
		}
	}
	// get from local cache.
	if arcs, count, err = s.dao.NewListCache(c, rid, tp, start, end); err != nil {
		err = nil
		addCache = false
	} else if len(arcs) > 0 {
		return
	}
	// get from archive-service
	if first {
		if allArcs, count, err = s.rankTopArcs(c, rid, tp, 1, max, ip); err != nil {
			err = nil
		} else if tid, ok := model.NewListRid[rid]; ok && len(allArcs) < conf.Conf.Rule.MinNewListCnt {
			if moreArcs, _, err = s.rankArcs(c, tid, tp, _samplePn, conf.Conf.Rule.MinNewListCnt, ip); err != nil {
				err = nil
			} else {
				allAidMap := make(map[int64]int64, len(allArcs))
				for _, arc := range allArcs {
					allAidMap[arc.Aid] = arc.Aid
				}
				for _, v := range moreArcs {
					if _, ok := allAidMap[v.Aid]; ok {
						continue
					}
					allArcs = append(allArcs, v)
				}
				log.Info("NewList more arcs rid(%d) tid(%d) len allArcs(%d) len moreArcs(%d)", tid, rid, len(allArcs), len(moreArcs))
			}
		}
	} else {
		if start >= max {
			if arcs, count, err = s.rankArcs(c, rid, tp, pn, ps, ip); err != nil {
				return
			}
			fmtArcs3(arcs)
			return
		}
		if allArcs, count, err = s.rankArcs(c, rid, tp, 1, max, ip); err != nil {
			err = nil
		}
	}
	if len(allArcs) > 0 {
		fmtArcs3(allArcs)
		length := len(allArcs)
		if length < start {
			arcs = []*api.Arc{}
			return
		}
		if length > end {
			arcs = allArcs[start : end+1]
		} else {
			arcs = allArcs[start:]
		}
		if addCache {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetNewListCache(c, rid, tp, allArcs, count)
			})
		}
		return
	}
	// get from remote cache.
	arcs, count, err = s.dao.NewListBakCache(c, rid, tp, start, end)
	if len(arcs) == 0 {
		arcs = []*api.Arc{}
	}
	return
}

func (s *Service) rankArcs(c context.Context, rid int32, tp int8, pn, ps int, ip string) (arcs []*api.Arc, count int, err error) {
	var res = &archive.RankArchives3{}
	if rid == _allRank {
		if res, err = s.arc.RankAllArcs3(c, &archive.ArgRankAll2{Pn: pn, Ps: ps}); err != nil {
			log.Error("arcrpc.RankAllArcs2(%d,%d) error(%v)", pn, ps, err)
			return
		}
	} else {
		arg := &archive.ArgRank2{Rid: int16(rid), Type: tp, Pn: pn, Ps: ps, RealIP: ip}
		if res, err = s.arc.RankArcs3(c, arg); err != nil {
			log.Error("arcrpc.RankArcs2(%d,%d,%d,%s) error(%v)", rid, pn, ps, ip, err)
			return
		}
	}
	arcs = res.Archives
	count = res.Count
	return
}

func (s *Service) rankTopArcs(c context.Context, rid int32, tp int8, pn, ps int, ip string) (arcs []*api.Arc, count int, err error) {
	arg := &archive.ArgRankTop2{ReID: int16(rid), Pn: pn, Ps: ps}
	if arcs, err = s.arc.RankTopArcs3(c, arg); err != nil {
		log.Error("arcrpc.RankTopArcs3(%d,%d,%d,%s) error(%v)", rid, pn, ps, ip, err)
	}
	return
}
