package region

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-show/model"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"

	"go-common/app/interface/main/app-show/model/region"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

const (
	_initCardKey      = "card_key_%d_%d"
	_bangumiSeasonID  = 1
	_bangumiEpisodeID = 2
)

var (
	_emptyShow      = &region.Show{}
	_emptyShowItems = []*region.ShowItem{}
	_bangumiRids    = map[int]struct{}{
		33:  struct{}{},
		32:  struct{}{},
		153: struct{}{},
		51:  struct{}{},
		152: struct{}{},
		168: struct{}{},
		169: struct{}{},
		170: struct{}{},
	}
	_bangumiReids = map[int]struct{}{
		13:  struct{}{},
		167: struct{}{},
	}
)

// Show region show
func (s *Service) Show(c context.Context, plat int8, rid, build int, mid int64, channel, buvid, network, mobiApp, device, adExtra string) (res *region.Show) {
	ip := metadata.String(c, metadata.RemoteIP)
	if model.IsOverseas(plat) {
		res = s.showOseaCache[rid]
	} else {
		res = s.showCache[rid]
	}
	if res == nil {
		res = _emptyShow
		return
	}
	if !model.IsIPad(plat) && len(res.Recommend) >= 4 {
		res = &region.Show{
			Recommend: res.Recommend[:4],
			New:       res.New,
			Dynamic:   res.Dynamic,
		}
	} else if model.IsIPad(plat) && len(res.Recommend) < 8 {
		var (
			max    = 8
			hotlen = len(res.Recommend)
		)
		if last := max - hotlen; last < len(res.New) {
			res.Recommend = append(res.Recommend, res.New[:last]...)
		} else if len(res.New) > 0 {
			res.Recommend = append(res.Recommend, res.New...)
		} else {
			if last < len(res.Dynamic) {
				res.Recommend = append(res.Recommend, res.Dynamic[:last]...)
			} else if len(res.Dynamic) > 0 {
				res.Recommend = append(res.Recommend, res.Dynamic...)
			}
		}
	}
	res.Banner = s.getBanners(c, plat, build, rid, mid, channel, ip, buvid, network, mobiApp, device, adExtra)
	return
}

// ShowDynamic show dynamic page
func (s *Service) ShowDynamic(c context.Context, plat int8, build, rid, pn, ps int) (res []*region.ShowItem) {
	var (
		isOsea      = model.IsOverseas(plat) //is overseas
		bangumiType = 0
	)
	if _, isBangumi := _bangumiReids[rid]; isBangumi {
		if (plat == model.PlatIPhone && build > 6090) || (plat == model.PlatAndroid && build > 514000) {
			bangumiType = _bangumiEpisodeID
		}
	}
	start := (pn - 1) * ps
	end := start + ps
	if bangumiType != 0 {
		if end < len(s.showDynamicAidsCache[rid]) {
			aids := s.showDynamicAidsCache[rid][start:end]
			res = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			if len(res) > 0 {
				return
			}
		}
	}
	if !isOsea {
		if end < len(s.dynamicCache[rid]) {
			res = s.dynamicCache[rid][start:end]
			return
		}
	} else {
		if end < len(s.dynamicOseaCache[rid]) {
			res = s.dynamicOseaCache[rid][start:end]
			return
		}
	}
	as, aids, err := s.dyn.RegionDynamic(c, rid, pn, ps)
	if err != nil {
		log.Error("s.rcmmnd.RegionDynamic(%d, %d, %d) error(%v)", rid, pn, ps, err)
		return
	}
	if bangumiType != 0 {
		res = s.fromArchivesPBBangumiOsea(c, as, aids, isOsea, bangumiType)
	} else {
		res = s.fromArchivesPBOsea(as, isOsea)
	}
	return
}

// ChildShow region child show
func (s *Service) ChildShow(c context.Context, plat int8, mid int64, rid, tid, build int, channel, mobiApp string, now time.Time) (res *region.Show) {
	var (
		isOsea      = model.IsOverseas(plat) //is overseas
		bangumiType = 0
		max         = 20
	)
	if _, isBangumi := _bangumiRids[rid]; isBangumi {
		if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
			bangumiType = _bangumiEpisodeID
		} else {
			bangumiType = _bangumiSeasonID
		}
	}
	if (mobiApp == "iphone" && build <= 4280) || (mobiApp == "ipad" && build <= 10400) || (mobiApp == "white" && build <= 101320) ||
		(mobiApp == "android" && build <= 501020) || (mobiApp == "android_tv" && build <= 1310) || mobiApp == "android_G" || mobiApp == "win" {
		bangumiType = 0
	}
	if tid == 0 {
		if bangumiType != 0 {
			var (
				hotTmp, newTmp []*region.ShowItem
				aids           []int64
				hotOk, newOk   bool
			)
			if aids, hotOk = s.childHotAidsCache[rid]; hotOk {
				if len(aids) > max {
					aids = aids[:max]
				}
				hotTmp = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			}
			if aids, newOk = s.childNewAidsCache[rid]; newOk {
				if len(aids) > max {
					aids = aids[:max]
				}
				newTmp = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			}
			if hotOk && newOk && len(hotTmp) > 0 && len(newTmp) > 0 {
				res, _ = s.mergeChildShow(hotTmp, newTmp)
				return
			}
		}
		if !isOsea {
			if res = s.childShowCache[rid]; res == nil {
				res = _emptyShow
				return
			}
		} else {
			if res = s.childShowOseaCache[rid]; res == nil {
				res = _emptyShow
				return
			}
		}
		return
	}
	res = s.mergeTagShow(c, rid, tid, bangumiType, isOsea, now)
	if mid > 0 {
		var err error
		if res.Tag, err = s.tag.TagInfo(c, mid, tid, now); err != nil {
			log.Error("s.tag.TagInfo(%d, %d) error(%v)", mid, tid, err)
		}
	}
	return
}

// ChildListShow region childList show
func (s *Service) ChildListShow(c context.Context, plat int8, rid, tid, pn, ps, build int, mid int64, order, platform, mobiApp, device string) (res []*region.ShowItem) {
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		isOsea      = model.IsOverseas(plat) //is overseas
		bangumiType = 0
	)
	start := (pn - 1) * ps
	end := start + ps
	key := fmt.Sprintf(_initRegionTagKey, rid, tid)
	if _, isBangumi := _bangumiRids[rid]; isBangumi {
		if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
			bangumiType = _bangumiEpisodeID
		} else {
			bangumiType = _bangumiSeasonID
		}
	}
	if (mobiApp == "iphone" && build <= 4280) || (mobiApp == "ipad" && build <= 10400) || (mobiApp == "white" && build <= 101320) ||
		(mobiApp == "android" && build <= 501020) || (mobiApp == "android_tv" && build <= 1310) || mobiApp == "android_G" || mobiApp == "win" {
		bangumiType = 0
	}
	if bangumiType != 0 {
		if tid == 0 && (order == "" || order == "new") && end < len(s.childNewAidsCache[rid]) {
			aids := s.childNewAidsCache[rid][start:end]
			res = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			if len(res) > 0 {
				return
			}
		}
		if tid > 0 && (order == "" || order == "new") && end < len(s.tagNewAidsCache[key]) {
			aids := s.tagNewAidsCache[key][start:end]
			res = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			if len(res) > 0 {
				return
			}
		}
	}
	if !isOsea {
		if tid == 0 && (order == "" || order == "new") && end < len(s.childNewCache[rid]) {
			res = s.childNewCache[rid][start:end]
			return
		}
		if tid > 0 && (order == "" || order == "new") && end < len(s.tagNewCache[key]) {
			res = s.tagNewCache[key][start:end]
			return
		}
	} else {
		if tid == 0 && (order == "" || order == "new") && end < len(s.childNewOseaCache[rid]) {
			res = s.childNewOseaCache[rid][start:end]
			return
		}
		if tid > 0 && (order == "" || order == "new") && end < len(s.tagNewOseaCache[key]) {
			res = s.tagNewOseaCache[key][start:end]
			return
		}
	}
	if (order == "" || order == "new") && tid == 0 {
		arcs, aids, err := s.arc.RanksArcs(c, rid, pn, ps)
		if err != nil {
			log.Error("s.rcmmnd.RegionArcList(%d, %d, %d, %d) error(%v)", rid, pn, ps, err)
		}
		if bangumiType != 0 {
			res = s.fromArchivesPBBangumiOsea(c, arcs, aids, isOsea, bangumiType)
		} else {
			res = s.fromArchivesPBOsea(arcs, isOsea)
		}
		return
	} else if (order == "" || order == "new") && tid > 0 {
		as, err := s.tag.NewArcs(c, rid, tid, pn, ps, time.Now())
		if err != nil {
			log.Error("s.tag.NewArcs(%d, %d) error(%v)", rid, tid, err)
			return
		}
		res = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
		return
	}
	var (
		tname string
		ok    bool
	)
	if tid > 0 {
		if tname, ok = s.tagsCache[key]; !ok {
			return
		}
	}
	as, err := s.search.SearchList(c, rid, build, pn, ps, mid, time.Now(), ip, order, tname, platform, mobiApp, device)
	if err != nil {
		log.Error("s.search.SearchList(%d, %d, %v, %d, %d) error(%v)", rid, tid, tname, pn, ps, err)
	}
	res = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
	return
}

// Dynamic region dynamic
func (s *Service) Dynamic(c context.Context, plat int8, rid, build int, mid int64, channel, buvid, network, mobiApp, device, adExtra string, now time.Time) (res *region.Show) {
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		isOsea   = model.IsOverseas(plat) //is overseas
		resCache *region.Show
	)
	s.prmobi.Incr("region_dynamic_plat_" + mobiApp)
	if isOsea {
		if resCache = s.regionFeedOseaCache[rid]; resCache == nil {
			resCache = _emptyShow
		}
	} else {
		if resCache = s.regionFeedCache[rid]; resCache == nil {
			resCache = _emptyShow
		}
	}
	if dyn, err := s.feedRegionDynamic(c, plat, rid, 0, 0, false, true, 0, mid, resCache.Recommend, now); err == nil && dyn != nil {
		res = dyn
	} else {
		res = resCache
	}
	if res != nil {
		res.Banner = s.getBanners(c, plat, build, rid, mid, channel, ip, buvid, network, mobiApp, device, adExtra)
		res.Card = s.regionCardDisplay(plat, build, rid)
	}
	return
}

// regionCardDisplay
func (s *Service) regionCardDisplay(plat int8, build, rid int) (res []*region.Head) {
	var ss []*region.Head
	key := fmt.Sprintf(_initCardKey, plat, rid)
	ss = s.cardCache[key]
	if len(ss) == 0 {
		res = []*region.Head{}
		return
	}
	res = []*region.Head{}
	for _, sw := range ss {
		if model.InvalidBuild(build, sw.Build, sw.Condition) {
			continue
		}
		tmp := &region.Head{}
		*tmp = *sw
		tmp.FillBuildURI(plat, build)
		res = append(res, tmp)
	}
	return
}

// DynamicList show dynamic list
func (s *Service) DynamicList(c context.Context, plat int8, rid int, pull bool, ctime, mid int64, now time.Time) (res *region.Show) {
	var (
		err error
	)
	if res, err = s.feedRegionDynamic(c, plat, rid, 0, 0, pull, false, ctime, mid, nil, now); err != nil || res == nil {
		res = _emptyShow
		return
	}
	return
}

// DynamicChild region show dynamic
func (s *Service) DynamicChild(c context.Context, plat int8, rid, tid, build int, mid int64, mobiApp string, now time.Time) (res *region.Show) {
	var (
		isOsea      = model.IsOverseas(plat) //is overseas
		bangumiType = 0
		resCache    *region.Show
		max         = 20
	)
	s.prmobi.Incr("region_dynamic_child_plat_" + mobiApp)
	if _, isBangumi := _bangumiRids[rid]; isBangumi {
		if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
			bangumiType = _bangumiEpisodeID
		} else {
			bangumiType = _bangumiSeasonID
		}
	}
	if tid == 0 {
		if bangumiType != 0 {
			var (
				hotTmp, newTmp []*region.ShowItem
				aids           []int64
				hotOk, newOk   bool
			)
			if aids, hotOk = s.childHotAidsCache[rid]; hotOk {
				if len(aids) > max {
					aids = aids[:max]
				}
				hotTmp = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			}
			if aids, newOk = s.childNewAidsCache[rid]; newOk {
				if len(aids) > max {
					aids = aids[:max]
				}
				newTmp = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			}
			if len(hotTmp) > 0 && len(newTmp) > 0 {
				resCache, _ = s.mergeChildShow(hotTmp, newTmp)
			}
		}
		if resCache == nil {
			if !isOsea {
				if resCache = s.childShowCache[rid]; resCache == nil {
					resCache = _emptyShow
				}
			} else {
				if resCache = s.childShowOseaCache[rid]; resCache == nil {
					resCache = _emptyShow
				}
			}
		}
	} else {
		resCache = s.mergeTagShow(c, rid, tid, bangumiType, isOsea, now)
	}
	if resCache == nil {
		resCache = _emptyShow
	}
	if dyn, err := s.feedRegionDynamic(c, plat, rid, tid, bangumiType, false, true, 0, mid, nil, now); err == nil && dyn != nil {
		res = dyn
		s.pHit.Incr("feed_region_dynamic")
	} else {
		res = resCache
		s.pMiss.Incr("feed_region_dynamic")
	}
	if res != nil {
		if tid == 0 {
			if tags, ok := s.regionTagCache[rid]; ok {
				res.TopTag = tags
			}
		} else if mid > 0 {
			var err error
			if res.Tag, err = s.tag.TagInfo(c, mid, tid, now); err != nil {
				log.Error("s.tag.TagInfo(%d, %d) error(%v)", mid, tid, err)
			}
		}
	}
	return
}

// DynamicListChild dynamic childList show
func (s *Service) DynamicListChild(c context.Context, plat int8, rid, tid, build int, pull bool, ctime, mid int64, now time.Time) (res *region.Show) {
	var (
		err         error
		bangumiType = 0
	)
	if _, isBangumi := _bangumiRids[rid]; isBangumi {
		if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
			bangumiType = _bangumiEpisodeID
		} else {
			bangumiType = _bangumiSeasonID
		}
	}
	if res, err = s.feedRegionDynamic(c, plat, rid, tid, bangumiType, pull, false, ctime, mid, nil, now); err != nil || res == nil {
		res = _emptyShow
		return
	}
	return
}

// feedRegionDynamic
func (s *Service) feedRegionDynamic(c context.Context, plat int8, rid, tid, bangumiType int, pull, isRecommend bool, ctime, mid int64, items []*region.ShowItem, now time.Time) (res *region.Show, err error) {
	var (
		smTagNum = 10
		smTagPos = 4
	)
	if res, err = s.feedDynamic(c, plat, rid, tid, bangumiType, pull, isRecommend, ctime, mid, items, now); err != nil || res == nil {
		log.Error("s.feedDynamic is null rid:%v tid:%v", rid, tid)
	} else {
		if isRecommend && tid > 0 {
			if st := s.similarTags(c, rid, tid, now); st != nil {
				if len(st) > smTagNum {
					res.TopTag = st[:smTagNum]
					res.NewTag = &region.NewTag{
						Position: smTagPos,
						Tag:      st[smTagNum:],
					}
				} else {
					res.TopTag = st
				}
			}
		}
	}
	return
}

// feedRegionDynamic
func (s *Service) feedDynamic(c context.Context, plat int8, rid, tid, bangumiType int, pull, isRecommend bool, ctime, mid int64, items []*region.ShowItem, now time.Time) (res *region.Show, err error) {
	var (
		isOsea        = model.IsOverseas(plat) //is overseas
		newAids       = []int64{}
		hotAids       = []int64{}
		ctop, cbottom xtime.Time
	)
	if hotAids, newAids, ctop, cbottom, err = s.rcmmnd.FeedDynamic(c, pull, rid, tid, ctime, mid, now); err != nil || len(newAids) == 0 {
		log.Error("s.rcmmnd.FeedDynamic(%v) error(%v)", rid, err)
		return
	}
	newItems := s.fromAidsOsea(c, newAids, false, isOsea, bangumiType)
	hotItems := s.fromAidsOsea(c, hotAids, false, isOsea, bangumiType)
	if isRecommend && (len(newItems) > 0 || len(hotItems) > 4) && len(newItems) > 0 {
		if len(hotItems) >= 4 {
			res = s.mergeDynamicShow(hotItems, newItems, ctop, cbottom)
		} else {
			log.Error("feedDynamic_hot is null rid:%v tid:%v", rid, tid)
			res = s.mergeDynamicShow(items, newItems, ctop, cbottom)
		}
	} else if !isRecommend && len(newItems) > 0 {
		res = s.mergeDynamicShow(_emptyShowItems, newItems, ctop, cbottom)
	} else {
		log.Error("feedDynamic_newItems is null rid:%v tid:%v", rid, tid)
	}
	s.infoc(mid, hotAids, newAids, rid, tid, pull, now)
	return
}

// fromArchives return region show items from archive archives.
func (s *Service) fromArchivesPB(as []*api.Arc) (is, isOsea []*region.ShowItem) {
	var asLen = len(as)
	if asLen == 0 {
		is = _emptyShowItems
		return
	}
	is = make([]*region.ShowItem, 0, asLen)
	for _, a := range as {
		i := &region.ShowItem{}
		i.FromArchivePB(a)
		if a.AttrVal(archive.AttrBitOverseaLock) == 0 {
			isOsea = append(isOsea, i)
		}
		is = append(is, i)
	}
	return
}

// fromArchivesPBBangumi aid to sid
func (s *Service) fromArchivesPBBangumi(c context.Context, as []*api.Arc, aids []int64, bangumiType int) (is, isOsea []*region.ShowItem) {
	var (
		asLen = len(as)
		err   error
		// bangumi
		// sids map[int64]int64
		sids map[int32]*seasongrpc.CardInfoProto
	)
	if asLen == 0 {
		is = _emptyShowItems
		return
	}
	if sids, err = s.fromSeasonID(c, aids); err != nil {
		log.Error("s.fromSeasonID error(%v)", err)
		return
	}
	is = make([]*region.ShowItem, 0, asLen)
	for _, a := range as {
		if sid, ok := sids[int32(a.Aid)]; ok && sid.SeasonId != 0 {
			i := &region.ShowItem{}
			i.FromBangumiArchivePB(a, sid, bangumiType)
			if a.AttrVal(archive.AttrBitOverseaLock) == 0 {
				isOsea = append(isOsea, i)
			}
			is = append(is, i)
		}
	}
	return
}

// mergeShow merge show
func (s *Service) mergeShow(hotTmp, newTmp, dynTmp []*region.ShowItem) (rs *region.Show) {
	rs = &region.Show{}
	if len(hotTmp) >= 4 {
		rs.Recommend = hotTmp
	} else {
		rs.Recommend = _emptyShowItems
	}
	if len(newTmp) >= 4 {
		rs.New = newTmp[:4]
	} else {
		rs.New = _emptyShowItems
	}
	if len(dynTmp) > 20 {
		rs.Dynamic = dynTmp[:20]
	} else if len(dynTmp) > 0 {
		rs.Dynamic = dynTmp
	} else {
		rs.Dynamic = _emptyShowItems
	}
	return
}

// mergeChildShow merge child show
func (s *Service) mergeChildShow(recTmp, newTmp []*region.ShowItem) (rs *region.Show, new []*region.ShowItem) {
	var (
		last int
	)
	rs = &region.Show{}
	if len(recTmp) >= 4 {
		rs.Recommend = recTmp[:4]
	} else if len(newTmp) >= 4 {
		var (
			max    = 4
			hotlen = len(recTmp)
		)
		rs.Recommend = recTmp
		if last = max - hotlen; last < len(newTmp) {
			rs.Recommend = append(rs.Recommend, newTmp[:last]...)
		} else {
			rs.Recommend = append(rs.Recommend, newTmp...)
		}
	} else {
		rs.Recommend = _emptyShowItems
	}
	new = newTmp
	if last > 0 {
		if len(new) > last {
			new = new[last:]
		}
	}
	if len(new) >= 20 {
		rs.New = new[:20]
	} else if len(new) > 0 {
		rs.New = new
	} else {
		rs.New = _emptyShowItems
	}
	return
}

// mergeDynamicShow
func (s *Service) mergeDynamicShow(recTmp, dynTmp []*region.ShowItem, ctop, cbottom xtime.Time) (rs *region.Show) {
	if ctop == 0 || cbottom == 0 {
		rs = &region.Show{}
	} else {
		rs = &region.Show{
			Ctop:    ctop,
			Cbottom: cbottom,
		}
	}
	if len(recTmp) >= 4 {
		rs.Recommend = recTmp[:4]
	} else {
		rs.Recommend = _emptyShowItems
	}
	if len(dynTmp) > 0 {
		rs.New = dynTmp
	} else {
		rs.New = _emptyShowItems
	}
	return
}

// mergeTagShow merge tag show
func (s *Service) mergeTagShow(c context.Context, rid, tid, bangumiType int, isOsea bool, now time.Time) (rs *region.Show) {
	const (
		strtNum  = 1
		hotNum   = 4
		newNum   = 20
		maxNum   = 50
		smTagNum = 10
		smTagPos = 4
	)
	rs = &region.Show{}
	// hots
	key := fmt.Sprintf(_initRegionTagKey, rid, tid)
	var (
		is []*region.ShowItem
		ok = false
	)
	if bangumiType != 0 {
		var aids []int64
		if aids, ok = s.tagHotAidsCache[key]; ok {
			is = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
		}
	} else if bangumiType == 0 || !ok {
		is, ok = s.tagHotCache[key]
		if isOsea {
			is, ok = s.tagHotOseaCache[key]
		}
	}
	if !ok {
		arcAids, err := s.tag.Hots(c, rid, tid, strtNum, hotNum, now)
		if err != nil {
			log.Error("s.tag.Hots(%d, %d) error(%v)", rid, tid, err)
		} else if is = s.fromAidsOsea(c, arcAids, true, isOsea, bangumiType); len(is) > 0 {
			if len(is) > hotNum {
				is = is[:hotNum]
			}
		}
	}
	rs.Recommend = is
	// news
	if bangumiType != 0 {
		var aids []int64
		if aids, ok = s.tagNewAidsCache[key]; ok {
			is = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
		}
	} else if bangumiType == 0 || !ok {
		is, ok = s.tagNewCache[key]
		if isOsea {
			is, ok = s.tagNewOseaCache[key]
		}
	}
	if !ok {
		as, err := s.tag.NewArcs(c, rid, tid, strtNum, maxNum, now)
		if err != nil {
			log.Error("s.tag.NewArcs(%d, %d) error(%v)", rid, tid, err)
			return
		}
		is = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
	}
	var (
		hotlen = len(rs.Recommend)
		last   int
	)
	if hotlen < hotNum {
		if last = hotNum - hotlen; last < len(is) {
			rs.Recommend = append(rs.Recommend, is[:last]...)
		} else {
			rs.Recommend = append(rs.Recommend, is...)
		}
	}
	if last > 0 {
		if len(is) > last {
			is = is[last:]
		}
	}
	if len(is) >= newNum {
		is = is[:newNum]
	}
	rs.New = is
	// similar tags
	if st := s.similarTags(c, rid, tid, now); st != nil {
		if len(st) > smTagNum {
			rs.TopTag = st[:smTagNum]
			rs.NewTag = &region.NewTag{
				Position: smTagPos,
				Tag:      st[smTagNum+1:],
			}
		} else {
			rs.TopTag = st
		}
	}
	return
}

// fromAids get Aids. bangumiType 1 seasonid  2 epid
func (s *Service) fromAids(c context.Context, arcAids []int64, isCheck bool, bangumiType int) (data, dataOsea []*region.ShowItem) {
	if len(arcAids) == 0 {
		return
	}
	var (
		ok  bool
		aid int64
		arc *api.Arc
		as  map[int64]*api.Arc
		err error
		// bangumi
		sids map[int32]*seasongrpc.CardInfoProto
	)
	if as, err = s.arc.ArchivesPB(c, arcAids); err != nil {
		log.Error("s.arc.ArchivesPB aids(%v) error(%v)", arcAids, err)
		return
	}
	data = make([]*region.ShowItem, 0, len(arcAids))
	if bangumiType != 0 {
		if sids, err = s.fromSeasonID(c, arcAids); err != nil {
			log.Error("s.fromSeasonID error(%v)", err)
			return
		}
	}
	for _, aid = range arcAids {
		if arc, ok = as[aid]; ok {
			if isCheck && !arc.IsNormal() {
				continue
			}
			i := &region.ShowItem{}
			if sid, ok := sids[int32(aid)]; ok && bangumiType != 0 && sid.SeasonId != 0 {
				i.FromBangumiArchivePB(arc, sid, bangumiType)
			} else {
				i.FromArchivePB(arc)
			}
			data = append(data, i)
			if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
				dataOsea = append(dataOsea, i)
			}
		}
	}
	return
}

// fromSeasonID
func (s *Service) fromSeasonID(c context.Context, arcAids []int64) (seasonID map[int32]*seasongrpc.CardInfoProto, err error) {
	if seasonID, err = s.bgm.CardsByAids(c, arcAids); err != nil {
		log.Error("s.bmg.CardsByAids error %v", err)
	}
	return
}

// similarTags similar tags
func (s *Service) similarTags(c context.Context, rid, tid int, now time.Time) (res []*region.SimilarTag) {
	key := fmt.Sprintf(_initRegionTagKey, rid, tid)
	res, ok := s.similarTagCache[key]
	if !ok {
		sts, err := s.tag.SimilarTag(c, rid, tid, now)
		if err != nil {
			log.Error("s.tag.SimilarTag error(%v)", err)
			return
		}
		res = sts
	}
	return
}

// isOverseas
func (s *Service) fromAidsOsea(ctx context.Context, aids []int64, isCheck, isOsea bool, bangumiType int) (data []*region.ShowItem) {
	tmp, tmpOsea := s.fromAids(ctx, aids, isCheck, bangumiType)
	if isOsea {
		data = tmpOsea
	} else {
		data = tmp
	}
	if data == nil {
		data = _emptyShowItems
	}
	return
}

// fromArchivesOsea  isOverseas
func (s *Service) fromArchivesPBOsea(as []*api.Arc, isOsea bool) (data []*region.ShowItem) {
	tmp, tmpOsea := s.fromArchivesPB(as)
	if isOsea {
		data = tmpOsea
	} else {
		data = tmp
	}
	if data == nil {
		data = _emptyShowItems
	}
	return
}

// fromArchivesOsea  isOverseas
func (s *Service) fromArchivesPBBangumiOsea(c context.Context, as []*api.Arc, aids []int64, isOsea bool, bangumiType int) (data []*region.ShowItem) {
	tmp, tmpOsea := s.fromArchivesPBBangumi(c, as, aids, bangumiType)
	if isOsea {
		data = tmpOsea
	} else {
		data = tmp
	}
	if data == nil {
		data = _emptyShowItems
	}
	return
}

// fromRankAids
func (s *Service) fromRankAids(ctx context.Context, aids []int64, others, scores map[int64]int64) (sis, sisOsea []*region.ShowItem) {
	var (
		aid  int64
		as   map[int64]*api.Arc
		arc  *api.Arc
		ok   bool
		err  error
		paid int64
	)
	if as, err = s.arc.ArchivesPB(ctx, aids); err != nil {
		log.Error("s.arc.ArchivesPB aids(%v) error(%v)", aids, err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB aids(%v) length is 0", aids)
		return
	}
	child := map[int64][]*region.ShowItem{}
	childOsea := map[int64][]*region.ShowItem{}
	for _, aid = range aids {
		if arc, ok = as[aid]; ok {
			if paid, ok = others[arc.Aid]; ok {
				i := &region.ShowItem{}
				i.FromArchivePBRank(arc, scores)
				child[paid] = append(child[paid], i)
				if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
					childOsea[paid] = append(childOsea[paid], i)
				}
			}
		}
	}
	for _, aid = range aids {
		if arc, ok = as[aid]; ok {
			if _, ok = others[arc.Aid]; !ok {
				i := &region.ShowItem{}
				i.FromArchivePBRank(arc, scores)
				if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
					if tmpchild, ok := childOsea[arc.Aid]; ok {
						i.Children = tmpchild
					}
					sisOsea = append(sisOsea, i)
				}
				if tmpchild, ok := child[arc.Aid]; ok {
					i.Children = tmpchild
				}
				sis = append(sis, i)
			}
		}
	}
	return
}

// fromCardAids get Aids.
func (s *Service) fromCardAids(ctx context.Context, aids []int64) (data map[int64]*region.ShowItem) {
	var (
		arc *api.Arc
		ok  bool
	)
	as, err := s.arc.ArchivesPB(ctx, aids)
	if err != nil {
		log.Error("s.arc.ArchivesPB error(%v)", err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", aids)
		return
	}
	data = map[int64]*region.ShowItem{}
	for _, aid := range aids {
		if arc, ok = as[aid]; ok {
			if !arc.IsNormal() {
				continue
			}
			i := &region.ShowItem{}
			i.FromArchivePB(arc)
			if region, ok := s.reRegionCache[int(arc.TypeID)]; ok {
				i.Reid = region.Rid
			}
			i.Desc = i.Rname
			data[aid] = i
		}
	}
	return
}
