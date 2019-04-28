package tag

import (
	"fmt"
	"strconv"
	"time"

	"context"
	"go-common/app/interface/main/app-tag/model"
	"go-common/app/interface/main/app-tag/model/bangumi"
	"go-common/app/interface/main/app-tag/model/feed"
	"go-common/app/interface/main/app-tag/model/region"
	"go-common/app/interface/main/app-tag/model/tag"
	bustag "go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	_emptySimilar   = []*tag.SimilarTag{}
	_emptyShowItems = []*region.ShowItem{}
	_emptyShow      = &region.Show{}
	_bangumiReids   = map[int]struct{}{
		13:  struct{}{},
		167: struct{}{},
	}
	_tagBlack = -4
)

func (s *Service) TagDetail(c context.Context, plat int8, mid, tid int64, now time.Time) (t *tag.Tag, similar []*tag.SimilarTag, err error) {
	var (
		info *bustag.Tag
	)
	if info, err = s.tg.InfoByID(c, mid, tid); err != nil {
		log.Error("s.tg.InfoByID(%d,%d) error(%v)", mid, tid, err)
		return
	}
	if info != nil {
		t = &tag.Tag{
			ID:      info.ID,
			Name:    info.Name,
			IsAtten: info.IsAtten,
			Count:   &tag.Count{Atten: info.Count.Atten},
		}
	}
	if similar, err = s.similarChangeTag(c, tid, now); err != nil {
		log.Error("s.tg.Similar(%d) error(%v)", tid, err)
		err = nil
	}
	if similar == nil {
		similar = _emptySimilar
	}
	return
}

func (s *Service) TagDefault(c context.Context, plat int8, tid, mid, ctime int64, pull bool, build int, buvid string, now time.Time) (is []*feed.Item, ctop, cbottom xtime.Time, err error) {
	var (
		isOsea = model.IsOverseas(plat)
		aids   []int64
		as     map[int64]*api.Arc
	)
	if aids, ctop, cbottom, err = s.rcmd.TagRecommand(c, plat, 0, build, tid, 10, mid, buvid); err != nil {
		if err == ecode.Int(_tagBlack) {
			log.Error("s.rcmd.TagRecommand(%) error(%v)", err)
			return
		} else {
			err = nil
			if _, aids, ctop, cbottom, err = s.rcmd.FeedDynamic(c, pull, 0, tid, 10, ctime, mid, now); err != nil {
				log.Error("s.FeedDynamic.Top(%d) error(%v)", tid, err)
				return
			}
		}
	}

	if as, err = s.arc.Archives(c, aids); err != nil {
		log.Error("s.arc.Arcivies(%v) error(%v)", aids, err)
		return
	}
	is = make([]*feed.Item, 0, len(as))
	for _, aid := range aids {
		if a, ok := as[aid]; ok {
			if !a.IsNormal() || (isOsea && a.AttrVal(archive.AttrBitOverseaLock) == 0) {
				continue
			}
			i := &feed.Item{}
			i.FromArc(a)
			is = append(is, i)
		}
	}
	return
}

func (s *Service) TagNew(c context.Context, plat int8, tid int64, pn, ps int, now time.Time) (is []*feed.Item, err error) {
	var (
		isOsea = model.IsOverseas(plat)
		as     map[int64]*api.Arc
		arc    *api.Arc
		ok     bool
		aids   []int64
	)
	if aids, err = s.tg.Detail(c, tid, pn, ps, now); err != nil {
		log.Error("s.tag.Detail(%d) error(%v)", tid, err)
		return
	}
	if as, err = s.arc.Archives(c, aids); err != nil {
		log.Error("s.arc.Archives() error(%v)", err)
		return
	}
	for _, aid := range aids {
		if arc, ok = as[aid]; ok {
			if !arc.IsNormal() || (isOsea && arc.AttrVal(archive.AttrBitOverseaLock) == 0) {
				continue
			}
			i := &feed.Item{}
			i.FromArc(arc)
			is = append(is, i)
		}
	}
	return
}

// similarChangeTag
func (s *Service) similarChangeTag(c context.Context, tid int64, now time.Time) (res []*tag.SimilarTag, err error) {
	var ok bool
	if res, ok = s.similarTagCache[tid]; !ok {
		if res, err = s.tg.SimilarTagChange(context.TODO(), tid, time.Now()); err != nil || len(res) == 0 {
			log.Error("s.tag.similarChangeTag(%d) error(%v)", tid, err)
			return
		}
	}
	return
}

// TagDynamic Tag region dynamic
func (s *Service) TagDynamic(c context.Context, plat int8, build, rid, reid int, tid, mid int64, mobiApp, buvid, device string, now time.Time, tagtab bool) (res *region.Show) {
	var (
		isOsea      = model.IsOverseas(plat) //is overseas
		resCache    *region.Show
		bangumiType = 0
		rn          = 10
		isRec       int
	)
	switch plat {
	case model.PlatIPad, model.PlatIPadI:
		rn = 20
	}
	s.prmobi.Incr("tag_dynamic_plat_" + mobiApp)
	if _, isBangumi := _bangumiReids[reid]; isBangumi {
		if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
			bangumiType = _bangumiEpisodeID
		} else {
			bangumiType = _bangumiSeasonID
		}
	}
	if (plat == model.PlatIPhone && build <= 5730) || (plat == model.PlatIPad) {
		bangumiType = 0
	}
	resCache = s.mergeTagRankShow(c, rid, reid, tid, bangumiType, isOsea, tagtab, now)
	if dyn, err := s.feedTagDynamic(c, plat, rid, reid, build, tid, 0, rn, bangumiType, true, tagtab, false, mid, buvid, now); err == nil && dyn != nil {
		res = dyn
		s.pHit.Incr("feed_tag_dynamic")
		isRec = 1
	} else {
		res = resCache
		s.pMiss.Incr("feed_tag_dynamic")
		isRec = 0
	}
	var (
		err  error
		info *bustag.Tag
	)
	if res != nil {
		if info, err = s.tg.InfoByID(c, mid, tid); err != nil {
			log.Error("s.tag.InfoByID(%d, %d) error(%v)", mid, tid, err)
			return
		}
		if tagtab {
			res.Tag = &tag.Tag{
				ID:      info.ID,
				Name:    info.Name,
				IsAtten: info.IsAtten,
				Count:   &tag.Count{Atten: info.Count.Atten},
			}
		}
	} else {
		res = &region.Show{New: _emptyShowItems}
	}
	infoc := &feedInfoc{
		mobiApp:    mobiApp,
		device:     device,
		build:      strconv.Itoa(build),
		now:        now.Format("2006-01-02 15:04:05"),
		pull:       "false",
		loginEvent: "0",
		tagID:      strconv.FormatInt(tid, 10),
		tagName:    info.Name,
		mid:        strconv.FormatInt(mid, 10),
		buvid:      buvid,
		displayID:  "0",
		feed:       res,
		isRec:      strconv.Itoa(isRec),
		topChannel: "0",
	}
	s.infoc(infoc)
	return
}

// TagDynamicList dynamic childList show
func (s *Service) TagDynamicList(c context.Context, plat int8, build, rid, reid int, tid int64, pull bool, ctime, mid int64,
	mobiApp, buvid, device, tagName string, now time.Time) (res *region.Show) {
	var (
		err         error
		bangumiType = 0
		rn          = 10
	)
	switch plat {
	case model.PlatIPad, model.PlatIPadI:
		rn = 20
	}
	if reid > 0 {
		if _, isBangumi := _bangumiReids[reid]; isBangumi {
			if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
				bangumiType = _bangumiEpisodeID
			} else {
				bangumiType = _bangumiSeasonID
			}
		}
		if (plat == model.PlatIPhone && build <= 5730) || (plat == model.PlatIPad) {
			bangumiType = 0
		}
	}
	if res, err = s.feedTagDynamic(c, plat, rid, reid, build, tid, ctime, rn, bangumiType, false, false, pull, mid, buvid, now); err != nil || res == nil {
		res = _emptyShow
		return
	}
	infoc := &feedInfoc{
		mobiApp:    mobiApp,
		device:     device,
		build:      strconv.Itoa(build),
		now:        now.Format("2006-01-02 15:04:05"),
		pull:       strconv.FormatBool(pull),
		loginEvent: "0",
		tagID:      strconv.FormatInt(tid, 10),
		tagName:    tagName,
		mid:        strconv.FormatInt(mid, 10),
		buvid:      buvid,
		displayID:  "0",
		feed:       res,
		isRec:      "1",
		topChannel: "0",
	}
	s.infoc(infoc)
	return
}

// TagRankList tag list rank
func (s *Service) TagRankList(c context.Context, plat int8, build, reid int, tid int64, pn, ps int, order, mobiApp string) (res []*region.ShowItem) {
	var (
		isOsea      = model.IsOverseas(plat) //is overseas
		bangumiType = 0
		key         string
	)
	s.prmobi.Incr("tag_rank_list_plat_" + mobiApp)
	start := (pn - 1) * ps
	end := start + ps
	if reid > 0 {
		if _, isBangumi := _bangumiReids[reid]; isBangumi {
			if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
				bangumiType = _bangumiEpisodeID
			} else {
				bangumiType = _bangumiSeasonID
			}
		}
		if (plat == model.PlatIPhone && build <= 5730) || (plat == model.PlatIPad) {
			bangumiType = 0
		}
		key = fmt.Sprintf(_initRegionTagKey, reid, tid)
	}
	if bangumiType != 0 {
		if reid > 0 && (order == "" || order == "new") && end < len(s.tagsDetailRankingAidsCache[key]) {
			aids := s.tagsDetailRankingAidsCache[key][start:end]
			res = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			if len(res) > 0 {
				return
			}
		}
		if reid == 0 && (order == "" || order == "new") && end < len(s.tagsDetailAidsCache[tid]) {
			aids := s.tagsDetailAidsCache[tid][start:end]
			res = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			if len(res) > 0 {
				return
			}
		}
	}
	if !isOsea {
		if reid > 0 && (order == "" || order == "new") && end < len(s.tagsDetailRankingCache[key]) {
			res = s.tagsDetailRankingCache[key][start:end]
			return
		}
		if reid == 0 && (order == "" || order == "new") && end < len(s.tagsDetailCache[tid]) {
			res = s.tagsDetailCache[tid][start:end]
			return
		}
	} else {
		if reid > 0 && (order == "" || order == "new") && end < len(s.tagsDetailRankingOseaCache[key]) {
			res = s.tagsDetailRankingOseaCache[key][start:end]
			return
		}
		if reid == 0 && (order == "" || order == "new") && end < len(s.tagsDetailOseaCache[tid]) {
			res = s.tagsDetailOseaCache[tid][start:end]
			return
		}
	}
	if order == "" || order == "new" {
		if reid > 0 {
			as, err := s.tg.DetailRanking(c, reid, tid, pn, ps, time.Now())
			if err != nil {
				log.Error("s.tag.DetailRanking(%d, %d) error(%v)", reid, tid, err)
				return
			}
			res = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
			return
		}
		if reid == 0 {
			as, err := s.tg.Detail(c, tid, pn, ps, time.Now())
			if err != nil {
				log.Error("s.tag.Detail(%d) error(%v)", tid, err)
				return
			}
			res = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
			return
		}
	}
	if res == nil {
		res = _emptyShowItems
	}
	return
}

// TagTab
func (s *Service) TagTab(c context.Context, rid int, tid, mid int64, now time.Time) (res *region.TagTab) {
	var (
		err  error
		info *bustag.Tag
	)
	res = &region.TagTab{}
	if tags := s.similarTagChange(c, rid, tid, now); len(tags) > 0 {
		res.TopTag = tags
	}
	if info, err = s.tg.InfoByID(c, mid, tid); err != nil {
		log.Error("s.tag.InfoByID(%d, %d) error(%v)", mid, tid, err)
		return
	}
	res.Tag = &tag.Tag{
		ID:      info.ID,
		Name:    info.Name,
		IsAtten: info.IsAtten,
		Count:   &tag.Count{Atten: info.Count.Atten},
	}
	return
}

// TagIDByName
func (s *Service) TagIDByName(c context.Context, tname string) (tid int64, err error) {
	var (
		ok   bool
		info *bustag.Tag
	)
	if tid, ok = s.tagsNameCache[tname]; !ok {
		if info, err = s.tg.InfoByName(c, tname); err != nil {
			log.Error("s.tag.InfoByName(%v) error(%v)", tname, err)
			return
		}
		tid = info.ID
	}
	return
}

// TagNameByID
func (s *Service) TagNameByID(c context.Context, tid int64) (tagName string, err error) {
	var (
		info *bustag.Tag
	)
	if info, err = s.tg.InfoByID(c, 0, tid); err != nil {
		log.Error("s.tag.InfoByID(%v) error(%v)", tid, err)
		return
	}
	tagName = info.Name
	return
}

// mergeTagShow merge tag rank show
func (s *Service) mergeTagRankShow(c context.Context, rid, reid int, tid int64, bangumiType int, isOsea, tagtab bool, now time.Time) (rs *region.Show) {
	const (
		strtNum = 1
		newNum  = 20
		maxNum  = 50
	)
	rs = &region.Show{}
	var (
		is []*region.ShowItem
		ok = false
	)
	if reid > 0 {
		// hots
		key := fmt.Sprintf(_initRegionTagKey, reid, tid)
		// news
		if bangumiType != 0 {
			var aids []int64
			if aids, ok = s.tagsDetailRankingAidsCache[key]; ok {
				is = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			}
		} else if bangumiType == 0 || !ok {
			is, ok = s.tagsDetailRankingCache[key]
			if isOsea {
				is, ok = s.tagsDetailRankingOseaCache[key]
			}
		}
		if !ok {
			as, err := s.tg.DetailRanking(c, reid, tid, strtNum, maxNum, now)
			if err != nil {
				log.Error("s.tag.DetailRanking(%d, %d) error(%v)", reid, tid, err)
				return
			}
			is = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
		}
	} else {
		if bangumiType != 0 {
			var aids []int64
			if aids, ok = s.tagsDetailAidsCache[tid]; ok {
				is = s.fromAidsOsea(c, aids, false, isOsea, bangumiType)
			}
		} else if bangumiType == 0 || !ok {
			is, ok = s.tagsDetailCache[tid]
			if isOsea {
				is, ok = s.tagsDetailOseaCache[tid]
			}
		}
		if !ok {
			as, err := s.tg.Detail(c, tid, strtNum, maxNum, now)
			if err != nil {
				log.Error("s.tag.Detail(%d) error(%v)", tid, err)
				return
			}
			is = s.fromAidsOsea(c, as, false, isOsea, bangumiType)
		}
	}
	if len(is) >= newNum {
		is = is[:newNum]
	}
	rs.New = is
	// similar tags
	if tagtab {
		if st := s.similarTagChange(c, rid, tid, now); st != nil {
			rs.TopTag = st
		}
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

// fromAids get Aids. bangumiType 1 seasonid  2 epid
func (s *Service) fromAids(c context.Context, arcAids []int64, isCheck bool, bangumiType int) (data, dataOsea []*region.ShowItem) {
	if len(arcAids) == 0 {
		return
	}
	var (
		ok  bool
		aid int64
		as  map[int64]*api.Arc
		arc *api.Arc
		err error
		// bangumi
		sids map[int64]*bangumi.SeasonInfo
	)
	if as, err = s.arc.Archives(c, arcAids); err != nil {
		log.Error("s.arc.Archives() error(%v)", err)
		return
	}
	data = make([]*region.ShowItem, 0, len(arcAids))
	if bangumiType != 0 {
		if sids, err = s.fromSeasonID(c, arcAids, time.Now()); err != nil {
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
			if sid, ok := sids[aid]; ok && bangumiType != 0 && sid.SeasonID != 0 {
				i.FromBangumi(arc, sid, bangumiType)
			} else {
				i.FromArchive(arc)
			}
			data = append(data, i)
			if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
				dataOsea = append(dataOsea, i)
			}
		}
	}
	return
}

func (s *Service) feedTagDynamic(c context.Context, plat int8, rid, reid, build int, tid, ctime int64, rn, bangumiType int, isRecommend, tagtab, pull bool, mid int64,
	buvid string, now time.Time) (res *region.Show, err error) {
	var regionID int
	if rid != 0 {
		regionID = rid
	} else {
		regionID = reid
	}
	if res, err = s.feedDynamic(c, plat, regionID, build, tid, ctime, rn, bangumiType, mid, buvid, isRecommend, pull, now); err != nil || res == nil {
		log.Error("s.feedDynamic is null rid:%v tid:%v", rid, tid)
	} else if isRecommend && tagtab {
		if tags := s.similarTagChange(c, rid, tid, now); len(tags) > 0 {
			res.TopTag = tags
		}
	}
	return
}

// feedRegionDynamic
func (s *Service) feedDynamic(c context.Context, plat int8, rid, build int, tid, ctime int64, rn, bangumiType int, mid int64, buvid string, isRecommend, pull bool,
	now time.Time) (res *region.Show, err error) {
	var (
		isOsea        = model.IsOverseas(plat) //is overseas
		newAids       = []int64{}
		hotAids       = []int64{}
		ctop, cbottom xtime.Time
	)
	if newAids, ctop, cbottom, err = s.rcmd.TagRecommand(c, plat, rid, build, tid, rn, mid, buvid); err != nil {
		log.Error("s.rcmd.TagRecommand(%) error(%v)", err)
		return
	}
	if len(newAids) == 0 {
		return
	}
	newItems := s.fromAidsOsea(c, newAids, false, isOsea, bangumiType)
	hotItems := s.fromAidsOsea(c, hotAids, false, isOsea, bangumiType)
	if isRecommend && len(hotItems) >= 4 {
		res = s.mergeDynamicShow(hotItems, newItems, ctop, cbottom)
	} else if len(newItems) > 0 {
		res = s.mergeDynamicShow(nil, newItems, ctop, cbottom)
	} else {
		log.Error("feedDynamic_newItems is null rid:%v tid:%v", rid, tid)
	}
	return
}

// similarTag
func (s *Service) similarTagChange(c context.Context, rid int, tid int64, now time.Time) (tag []*tag.SimilarTag) {
	if rid != 0 {
		if tags, ok := s.regionTagCache[rid]; ok {
			tag = tags
		}
	} else if tid != 0 {
		if tags, ok := s.similarTagCache[tid]; ok {
			tag = tags
		} else if st, err := s.tg.SimilarTagChange(c, tid, now); err == nil && st != nil {
			tag = st
		}
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

// fromSeasonID
func (s *Service) fromSeasonID(c context.Context, arcAids []int64, now time.Time) (seasonID map[int64]*bangumi.SeasonInfo, err error) {
	if seasonID, err = s.bgm.Seasonid(arcAids, now); err != nil {
		log.Error("s.bmg.Seasonid error %v", err)
	}
	return
}
