package service

import (
	"context"
	"strconv"
	"sync"

	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/dynamic/conf"
	"go-common/app/service/main/dynamic/dao"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

var (
	_emptyArcs3 = make([]*api.Arc, 0)
)

const (
	_arcMuch   = 5
	_guoChuang = 168
)

// RegionTotal total  dynamic of regin
func (s *Service) RegionTotal(c context.Context) (res map[string]int) {
	res = map[string]int{}
	for k, v := range s.regionTotal {
		res[strconv.FormatInt(int64(k), 10)] = v
	}
	res["live"] = s.live
	return
}

// regionNeedCache
func regionNeedCache(regionArcs map[int32][]int64) bool {
	if len(regionArcs) <= 0 {
		return false
	}
	for k, v := range regionArcs {
		if len(v) < conf.Conf.Rule.MinRegionCount {
			log.Error("arcNeedCache key(%d) len(%v)<=0", k, v)
			return false
		}
	}
	return true
}

// RegionArcs3 get new arcs of region.
func (s *Service) RegionArcs3(c context.Context, rid int32, pn, ps int) (arcs []*api.Arc, count int, err error) {
	var (
		ok         bool
		start, end int
		aids       []int64
	)
	if aids, ok = s.regionArcs[rid]; !ok {
		arcs = _emptyArcs3
		return
	}
	count = len(aids)
	start = (pn - 1) * ps
	end = start + ps + _arcMuch
	if start > count {
		arcs = _emptyArcs3
		return
	}
	if end > count {
		aids = aids[start:]
	} else {
		aids = aids[start:end]
	}
	if arcs, err = s.normalArcs3(c, aids, ps); err != nil {
		log.Error("archives(%v) error(%v)", aids, err)
	}
	return
}

// normalArcs3 .
func (s *Service) normalArcs3(c context.Context, aids []int64, ps int) (res []*api.Arc, err error) {
	var (
		arcs   []*api.Arc
		tmpRes map[int64]*api.Arc
	)
	archivesLog("normalArcs3", aids)
	if tmpRes, err = s.archives3(c, aids); err != nil {
		return
	}
	for _, aid := range aids {
		if arc, ok := tmpRes[aid]; ok {
			arcs = append(arcs, arc)
		} else {
			log.Error("normalArcs s.archives aid(%d) nil", aid)
		}
	}
	res = filterArc3(arcs, ps)
	return
}

// archives3 .
func (s *Service) archives3(c context.Context, aids []int64) (res map[int64]*api.Arc, err error) {
	arg := &arcmdl.ArgAids2{Aids: aids}
	if res, err = s.arcRPC.Archives3(c, arg); err != nil {
		dao.PromError("稿件rpc接口:Archives2", "archives s.arcRPC.Archives3(%v) error(%v)", aids, err)
	}
	return
}

func filterArc3(arcs []*api.Arc, count int) (res []*api.Arc) {
	tmpPs := 1
	for _, arc := range arcs {
		if tmpPs <= count && isShow3(arc) {
			res = append(res, arc)
			tmpPs = tmpPs + 1
		} else if tmpPs > count {
			break
		}
	}
	return
}

// isShow3
func isShow3(a *api.Arc) bool {
	return a.IsNormal() && a.AttrVal(arcmdl.AttrBitNoDynamic) == arcmdl.AttrNo
}

// RegionTagArcs3 get new arcs of region and hot tag.
func (s *Service) RegionTagArcs3(c context.Context, rid int32, tagID int64, pn, ps int) (arcs []*api.Arc, count int, err error) {
	var (
		ok         bool
		start, end int
		aids       []int64
	)
	key := regionTagKey(rid, tagID)
	if aids, ok = s.regionTagArcs[key]; !ok {
		arcs = _emptyArcs3
		return
	}
	count = len(aids)
	start = (pn - 1) * ps
	end = start + ps + _arcMuch
	if start > count {
		arcs = _emptyArcs3
		return
	}
	if end > count {
		aids = aids[start:]
	} else {
		aids = aids[start:end]
	}
	if arcs, err = s.normalArcs3(c, aids, ps); err != nil {
		log.Error("archives(%v) error(%v)", aids, err)
	}
	return
}

// RegionsArcs3 get batch new arcs of regions.
func (s *Service) RegionsArcs3(c context.Context, rids []int32, count int) (mArcs map[int32][]*api.Arc, err error) {
	var (
		ok      bool
		noRids  []int32
		allAids []int64
		aids    []int64
		mAids   map[int32][]int64
		res     map[int64]*api.Arc
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	mAids = make(map[int32][]int64, len(rids))
	for _, rid := range rids {
		end := count + _arcMuch
		if aids, ok = s.regionArcs[rid]; !ok || len(aids) == 0 {
			continue
		}
		if end > len(aids) {
			end = len(aids)
		}
		allAids = append(allAids, aids[:end]...)
		mAids[rid] = aids[:end]
	}
	archivesLog("RegionsArcs3", allAids)
	if res, err = s.archives3(c, allAids); err != nil {
		log.Error("archives(%v) error(%v)", allAids, err)
		return
	}
	mArcs = make(map[int32][]*api.Arc, len(rids))
	for _, rid := range rids {
		var arcs []*api.Arc
		for _, aid := range mAids[rid] {
			if arc, ok := res[aid]; ok {
				arcs = append(arcs, arc)
			} else {
				log.Error("RegionsArcs s.archives aid(%d) nil", aid)
			}
		}
		mArcs[rid] = filterArc3(arcs, count)
		if len(mArcs[rid]) < count {
			dao.PromError("一级分区数据错误", "RegionsArcs rid(%d) len(mArcs[rid])(%d) count(%d)", rid, len(mArcs[rid]), count)
			noRids = append(noRids, rid)
		}
	}
	//last back up from rankIndexArc.
	if len(noRids) > 0 {
		err = s.rankIndexArc3(c, noRids, 1, count, ip, mArcs)
	}
	return
}

// rankIndexArc3 archives3.
func (s *Service) rankIndexArc3(c context.Context, rids []int32, pn, ps int, ip string, mArcs map[int32][]*api.Arc) (err error) {
	var mutex = sync.Mutex{}
	group, errCtx := errgroup.WithContext(c)
	for _, rid := range rids {
		trid := rid
		group.Go(func() (err error) {
			var (
				topRes, tmpRes []*api.Arc
				regionRes      = &arcmdl.RankArchives3{}
			)
			if trid == _guoChuang {
				arg := &arcmdl.ArgRank2{Rid: int16(trid), Type: 0, Pn: pn, Ps: ps + _arcMuch, RealIP: ip}
				if regionRes, err = s.arcRPC.RankArcs3(errCtx, arg); err != nil {
					dao.PromError("RankArcs2接口错误", "arcRPC.RankArcs2(%d,%d,%d,%s) error(%v)", trid, pn, ps, ip, err)
					return nil
				}
				topRes = regionRes.Archives
			} else {
				arg := &arcmdl.ArgRankTop2{ReID: int16(trid), Pn: pn, Ps: ps + _arcMuch}
				if topRes, err = s.arcRPC.RankTopArcs3(errCtx, arg); err != nil {
					dao.PromError("RankTopArcs2接口错误", "arcRPC.RankTopArcs2(%d,%d,%d,%s) error(%v)", trid, pn, ps, ip, err)
					return nil
				}
			}
			tmpRes = filterArc3(topRes, ps)
			mutex.Lock()
			if len(tmpRes) == 0 {
				mArcs[trid] = _emptyArcs3
			} else {
				mArcs[trid] = tmpRes
			}
			mutex.Unlock()
			return nil
		})
	}
	group.Wait()
	return
}
