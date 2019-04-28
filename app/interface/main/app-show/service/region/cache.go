package region

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/bangumi"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/app/interface/main/app-show/model/tag"
	"go-common/app/service/main/archive/api"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/log"
)

var (
	// 番剧 动画，音乐，舞蹈，游戏，科技，娱乐，鬼畜，电影，时尚, 生活，国漫
	_tids = []int{13, 1, 3, 129, 4, 36, 5, 119, 23, 155, 160, 11, 167}
)

// loadBanner1 load banner cache.
func (s *Service) loadBanner() {
	var (
		resbs = map[int8]map[int][]*resource.Banner{}
	)
	for plat, resIDStr := range _bannersPlat {
		mobiApp := model.MobiApp(plat)
		res, err := s.res.ResBanner(context.TODO(), plat, 515007, 0, resIDStr, "master", "", "", "", mobiApp, "", "", false)
		if err != nil || len(res) == 0 {
			log.Error("s.res.ResBanner is null or err(%v)", err)
			return
		}
		resbs[plat] = res
	}
	if len(resbs) > 0 {
		s.bannerCache = resbs
	}
	log.Info("loadBannerCahce success")
}

func (s *Service) loadbgmBanner() {
	bgmBanners := map[int8]map[int][]*bangumi.Banner{}
	for plat, b := range _bannersPGC {
		ridBanners := map[int][]*bangumi.Banner{}
		for rid, pgcID := range b {
			bBanner, err := s.bgm.Banners(context.TODO(), pgcID)
			if err != nil {
				log.Error("s.bgmdao.Banners is null or err(%v)", err)
				return
			}
			ridBanners[rid] = bBanner
		}
		bgmBanners[plat] = ridBanners
	}
	s.bannerBmgCache = bgmBanners
}

// loadRegionShow
func (s *Service) loadShow() {
	// default use android regions TODO
	regionkey := fmt.Sprintf(_initRegionKey, model.PlatAndroid, _initlanguage)
	res := s.cachelist[regionkey]
	var (
		// region tmp
		showTmp  = map[int][]*region.ShowItem{} // tmp hotCache
		showNTmp = map[int][]*region.ShowItem{} // tmp newCache
		showD    = map[int][]*region.ShowItem{} // tmp dynamicCache
		sTmp     = map[int]*region.Show{}
		// region tmp overseas
		showTmpOsea  = map[int][]*region.ShowItem{} // tmp overseas hotCache
		showNTmpOsea = map[int][]*region.ShowItem{} // tmp overseas newCache
		showDOsea    = map[int][]*region.ShowItem{} // tmp overseas dynamicCache
		sTmpOsea     = map[int]*region.Show{}
		// aids
		showDynamicAids = map[int][]int64{}
		// new region feed
		regionFeedTmp     = map[int]*region.Show{}
		regionFeedTmpOsea = map[int]*region.Show{}
	)
	for _, v := range res {
		if v.Reid == 0 {
			if v.Rid == 65537 || v.Rid == 65539 || v.Rid == 65541 || v.Rid == 65543 {
				continue
			}
			var (
				tmp, tmpOsea []*region.ShowItem
				aidsTmp      []int64
			)
			tmp, tmpOsea = s.loadShowHot(v.Rid)
			showTmp[v.Rid], showTmpOsea[v.Rid] = s.upCache(tmp, tmpOsea, s.hotCache[v.Rid], s.hotOseaCache[v.Rid])
			tmp, tmpOsea = s.loadShowNewRpc(v.Rid)
			showNTmp[v.Rid], showNTmpOsea[v.Rid] = s.upCache(tmp, tmpOsea, s.newCache[v.Rid], s.newOseaCache[v.Rid])
			tmp, tmpOsea, aidsTmp = s.loadShowDynamic(v.Rid)
			showDynamicAids[v.Rid] = s.upAidsCache(aidsTmp, s.showDynamicAidsCache[v.Rid])
			showD[v.Rid], showDOsea[v.Rid] = s.upCache(tmp, tmpOsea, s.dynamicCache[v.Rid], s.dynamicOseaCache[v.Rid])
			sTmp[v.Rid] = s.mergeShow(showTmp[v.Rid], showNTmp[v.Rid], showD[v.Rid])
			// overseas
			sTmpOsea[v.Rid] = s.mergeShow(showTmpOsea[v.Rid], showNTmpOsea[v.Rid], showDOsea[v.Rid])
			// new region feed
			regionFeedTmp[v.Rid], _ = s.mergeChildShow(showTmp[v.Rid], showNTmp[v.Rid])
			regionFeedTmpOsea[v.Rid], _ = s.mergeChildShow(showTmpOsea[v.Rid], showNTmpOsea[v.Rid])
		}
	}
	s.showCache = sTmp
	s.hotCache = showTmp
	s.newCache = showNTmp
	s.dynamicCache = showD
	// overseas
	s.showOseaCache = sTmpOsea
	s.hotOseaCache = showTmpOsea
	s.newOseaCache = showNTmpOsea
	s.dynamicOseaCache = showDOsea
	// new region feed
	s.regionFeedCache = regionFeedTmp
	s.regionFeedOseaCache = regionFeedTmpOsea
	// aids
	s.showDynamicAidsCache = showDynamicAids
}

// loadShowChild
func (s *Service) loadShowChild() {
	// default use android regions TODO
	regionkey := fmt.Sprintf(_initRegionKey, model.PlatAndroid, _initlanguage)
	res := s.cachelist[regionkey]
	var (
		scTmp  = map[int]*region.Show{}
		showC  = map[int][]*region.ShowItem{} // tmp childHotCache
		showCN = map[int][]*region.ShowItem{} // tmp childNewCache
		// region tmp overseas
		scTmpOsea  = map[int]*region.Show{}
		showCOsea  = map[int][]*region.ShowItem{} // tmp overseas childHotCache
		showCNOsea = map[int][]*region.ShowItem{} // tmp overseas childNewCache
		// aids
		showChildAids = map[int][]int64{}
		showNewAids   = map[int][]int64{}
	)
	for _, v := range res {
		if v.Reid != 0 {
			var (
				tmp, tmpOsea []*region.ShowItem
				aidsTmp      []int64
			)
			tmp, tmpOsea, aidsTmp = s.loadShowChileHot(v.Rid)
			showChildAids[v.Rid] = s.upAidsCache(aidsTmp, s.childHotAidsCache[v.Rid])
			showC[v.Rid], showCOsea[v.Rid] = s.upCache(tmp, tmpOsea, s.childHotCache[v.Rid], s.childHotOseaCache[v.Rid])
			tmp, tmpOsea, aidsTmp = s.loadShowChildNew(v.Rid)
			showNewAids[v.Rid] = s.upAidsCache(aidsTmp, s.childNewAidsCache[v.Rid])
			showCN[v.Rid], showCNOsea[v.Rid] = s.upCache(tmp, tmpOsea, s.childNewCache[v.Rid], s.childNewOseaCache[v.Rid])
			scTmp[v.Rid], showCN[v.Rid] = s.mergeChildShow(showC[v.Rid], showCN[v.Rid])
			// overseas
			scTmpOsea[v.Rid], showCNOsea[v.Rid] = s.mergeChildShow(showCOsea[v.Rid], showCNOsea[v.Rid])
		}
	}
	s.childShowCache = scTmp
	s.childHotCache = showC
	s.childNewCache = showCN
	// overseas
	s.childShowOseaCache = scTmpOsea
	s.childHotOseaCache = showCOsea
	s.childNewOseaCache = showCNOsea
	// region child aids
	s.childHotAidsCache = showChildAids
	s.childNewAidsCache = showNewAids
}

func (s *Service) loadShowChildTagsInfo() {
	// default use android regions TODO
	regionkey := fmt.Sprintf(_initRegionKey, model.PlatAndroid, _initlanguage)
	res := s.cachelist[regionkey]
	reslist := s.regionListCache[regionkey]
	var (
		// tag tmp
		tagsRegionTmp = map[int][]*region.SimilarTag{} // region tags
		tagsTmp       = map[string]string{}            // tagid cache
	)
	for _, v := range res {
		if v.Reid != 0 {
			//tag
			var rTmp *region.Region
			if r, ok := reslist[v.Reid]; ok {
				rTmp = r
			}
			if tids := s.loadShowChildTagIDs(v.Rid); len(tids) > 0 {
				for _, tag := range tids {
					tagInfo := &region.SimilarTag{
						TagId:   int(tag.Tid),
						TagName: tag.Name,
						Rid:     v.Rid,
						Rname:   v.Name,
					}
					if rTmp != nil {
						tagInfo.Reid = rTmp.Rid
						tagInfo.Rename = rTmp.Name
					}
					//tags info
					tagsRegionTmp[v.Rid] = append(tagsRegionTmp[v.Rid], tagInfo)
					key := fmt.Sprintf(_initRegionTagKey, v.Rid, tag.Tid)
					tagsTmp[key] = tag.Name
				}
			}
		}
	}
	// region child aids
	s.regionTagCache = tagsRegionTmp
	s.tagsCache = tagsTmp
}

// loadShowHot
func (s *Service) loadShowHot(rid int) (resData, resOseaData []*region.ShowItem) {
	res, err := s.rcmmnd.RegionHots(context.TODO(), rid)
	if err != nil {
		log.Error("s.rcmmnd.RegionHot(%d) error(%v)", rid, err)
		return
	}
	if len(res) > 8 {
		res = res[:8]
	}
	resData, resOseaData = s.fromAids(context.TODO(), res, false, 0)
	log.Info("loadShowHot(%d) success", rid)
	return
}

// loadShowNewRpc
func (s *Service) loadShowNewRpc(rid int) (resData, resOseaData []*region.ShowItem) {
	arcs, err := s.arc.RankTopArcs(context.TODO(), rid, 1, 20)
	if err != nil {
		log.Error("s.arc.RankTopArcs(%d) error(%v)", rid, err)
		return
	}
	if len(arcs) > 20 {
		arcs = arcs[:20]
	}
	resData, resOseaData = s.fromArchivesPB(arcs)
	log.Info("loadShowNewRpc(%d) success", rid)
	return
}

// loadShowDynamic
func (s *Service) loadShowDynamic(rid int) (resData, resOseaData []*region.ShowItem, arcAids []int64) {
	var (
		err  error
		arcs []*api.Arc
	)
	arcs, arcAids, err = s.dyn.RegionDynamic(context.TODO(), rid, 1, 100)
	if err != nil || len(arcs) < 20 {
		log.Error("s.rcmmnd.RegionDynamic(%d) error(%v)", rid, err)
		return
	}
	resData, resOseaData = s.fromArchivesPB(arcs)
	log.Info("loadShowRPCDynamic(%d) success", rid)
	return
}

// loadShowChileHot
func (s *Service) loadShowChileHot(rid int) (resData, resOseaData []*region.ShowItem, arcAids []int64) {
	var err error
	arcAids, err = s.rcmmnd.RegionChildHots(context.TODO(), rid)
	if err != nil || len(arcAids) < 4 {
		log.Error("s.rcmmnd.RegionChildHots(%d) error(%v)", rid, err)
		return
	}
	resData, resOseaData = s.fromAids(context.TODO(), arcAids, false, 0)
	log.Info("loadShowChileHot(%d) success", rid)
	return
}

// loadShowChildNew
func (s *Service) loadShowChildNew(rid int) (resData, resOseaData []*region.ShowItem, arcAids []int64) {
	var (
		err  error
		arcs []*api.Arc
	)
	arcs, arcAids, err = s.arc.RanksArcs(context.TODO(), rid, 1, 300)
	if err != nil || len(arcAids) < 20 {
		log.Error("s.arc.RanksArc(%d) error(%v)", rid, err)
		return
	}
	resData, resOseaData = s.fromArchivesPB(arcs)
	log.Info("loadShowChildNew(%d) success", rid)
	return
}

// loadShowChildTagIDs
func (s *Service) loadShowChildTagIDs(rid int) (tags []*tag.Tag) {
	tags, err := s.tag.TagHotsId(context.TODO(), rid, time.Now())
	if err != nil || len(tags) == 0 {
		log.Error("s.tag.loadShowChildTagIDs(%d) error(%v)", rid, err)
		return
	}
	return
}

// loadRankRegionCache
func (s *Service) loadRankRegionCache() {
	var (
		tmp     = map[int][]*region.ShowItem{}
		tmpOsea = map[int][]*region.ShowItem{}
	)
	for _, rid := range _tids {
		aids, others, scores, err := s.rcmmnd.RankAppRegion(context.TODO(), rid)
		if err != nil {
			log.Error("s.rcmmnd.RankAppRegion rid (%v) error(%v)", rid, err)
			return
		}
		tRank, tOseaRank := s.fromRankAids(context.TODO(), aids, others, scores)
		tmp[rid] = tRank
		tmpOsea[rid] = tOseaRank
	}
	if len(tmp) > 0 {
		s.rankCache = tmp
	}
	if len(tmpOsea) > 0 {
		s.rankOseaCache = tmpOsea
	}
}

// loadRegionListCache
func (s *Service) loadRegionListCache() {
	res, err := s.dao.RegionPlat(context.TODO())
	if err != nil {
		log.Error("s.dao.RegionPlat error(%v)", err)
		return
	}
	tmpRegion := map[int]*region.Region{}
	tmp := map[int]*region.Region{}
	for _, v := range res {
		// region list map
		tmpRegion[v.Rid] = v
	}
	for _, r := range res {
		if r.Reid != 0 {
			if rerg, ok := tmpRegion[r.Reid]; ok {
				tmp[r.Rid] = rerg
			}
		}
	}
	s.reRegionCache = tmp
}

func (s *Service) loadColumnListCache(now time.Time) {
	var (
		tmpChild = map[int]*card.ColumnList{}
	)
	columns, err := s.cdao.ColumnList(context.TODO(), now)
	if err != nil {
		log.Error("s.cdao.ColumnList error(%v)", err)
		return
	}
	for _, column := range columns {
		tmpChild[column.Cid] = column
	}
	s.columnListCache = tmpChild
}

// loadCardCache load all card cache
func (s *Service) loadCardCache(now time.Time) {
	hdm, err := s.cdao.PosRecs(context.TODO(), now)
	if err != nil {
		log.Error("s.cdao.PosRecs error(%v)", err)
		return
	}
	itm, aids, err := s.cdao.RecContents(context.TODO(), now)
	if err != nil {
		log.Error("s.cdao.RecContents error(%v)", err)
		return
	}
	tmpItem := map[int]map[int64]*region.ShowItem{}
	for recid, aid := range aids {
		tmpItem[recid] = s.fromCardAids(context.TODO(), aid)
	}
	tmp := s.mergeCard(context.TODO(), hdm, itm, tmpItem, now)
	s.cardCache = tmp
}

func (s *Service) mergeCard(c context.Context, hdm map[int8]map[int][]*card.Card, itm map[int][]*card.Content, tmpItems map[int]map[int64]*region.ShowItem, now time.Time) (res map[string][]*region.Head) {
	// default use android regions TODO
	var (
		_topic     = 1
		_activity  = 0
		regionkey  = fmt.Sprintf(_initRegionKey, model.PlatAndroid, _initlanguage)
		regionList = s.cachelist[regionkey]
	)
	res = map[string][]*region.Head{}
	for _, v := range regionList {
		if v.Reid != 0 {
			continue
		}
		for plat, phds := range hdm {
			hds, ok := phds[v.Rid]
			if !ok {
				continue
			}
			for _, hd := range hds {
				key := fmt.Sprintf(_initCardKey, plat, v.Rid)
				var (
					sis []*region.ShowItem
				)
				its, ok := itm[hd.ID]
				if !ok {
					its = []*card.Content{}
				}
				tmpItem, ok := tmpItems[hd.ID]
				if !ok {
					tmpItem = map[int64]*region.ShowItem{}
				}
				// 1 daily 2 topic 3 activity 4 rank 5 polymeric_card
				switch hd.Type {
				case 1:
					for _, ci := range its {
						si := s.fillCardItem(ci, tmpItem)
						if si.Title != "" {
							sis = append(sis, si)
						}
					}
				case 2:
					if topicID, err := strconv.ParseInt(hd.Rvalue, 10, 64); err == nil {
						if actm, err := s.act.Activitys(c, []int64{topicID}, _topic, ""); err != nil {
							log.Error("s.act.Activitys topicID error (%v)", topicID, err)
						} else {
							if act, ok := actm[topicID]; ok && act.H5Cover != "" && act.H5URL != "" {
								si := &region.ShowItem{}
								si.FromTopic(act)
								sis = []*region.ShowItem{si}
							}
						}
					}
				case 3:
					if topicID, err := strconv.ParseInt(hd.Rvalue, 10, 64); err == nil {
						if actm, err := s.act.Activitys(c, []int64{topicID}, _activity, ""); err != nil {
							log.Error("s.act.Activitys topicID error (%v)", topicID, err)
						} else {
							if act, ok := actm[topicID]; ok && act.H5Cover != "" && act.H5URL != "" {
								si := &region.ShowItem{}
								si.FromActivity(act, now)
								sis = []*region.ShowItem{si}
							}
						}
					}
				case 4:
					if tmpRank, ok := s.rankCache[v.Rid]; ok {
						if len(tmpRank) > 3 {
							sis = tmpRank[:3]
						} else {
							sis = tmpRank
						}
					}
				case 5, 6, 8:
					for _, ci := range its {
						si := s.fillCardItem(ci, tmpItem)
						if si.Title != "" {
							sis = append(sis, si)
						}
					}
				case 7:
					si := &region.ShowItem{
						Title: hd.Title,
						Cover: hd.Cover,
						Desc:  hd.Desc,
						Goto:  hd.Goto,
						Param: hd.Param,
					}
					if hd.Goto == model.GotoColumnStage {
						paramInt, _ := strconv.Atoi(hd.Param)
						if c, ok := s.columnListCache[paramInt]; ok {
							cidStr := strconv.Itoa(c.Ceid)
							si.URI = model.FillURICategory(hd.Goto, cidStr, hd.Param)
						}
					} else {
						si.URI = hd.URi
					}
					sis = append(sis, si)
				default:
					continue
				}
				if len(sis) == 0 {
					continue
				}
				sw := &region.Head{
					CardID:    hd.ID,
					Title:     hd.Title,
					Type:      hd.TypeStr,
					Build:     hd.Build,
					Condition: hd.Condition,
					Plat:      hd.Plat,
				}
				if hd.Cover != "" {
					sw.Cover = hd.Cover
				}
				switch sw.Type {
				case model.GotoDaily:
					sw.Date = now.Unix()
					sw.Param = hd.Rvalue
					sw.URI = hd.URi
					sw.Goto = hd.Goto
				case model.GotoCard:
					sw.URI = hd.URi
					sw.Goto = hd.Goto
					sw.Param = hd.Param
				case model.GotoRank:
					sw.Param = strconv.Itoa(v.Rid)
				case model.GotoTopic, model.GotoActivity:
					if sw.Title == "" {
						if len(sis) > 0 {
							sw.Title = sis[0].Title
						}
					}
				case model.GotoVeidoCard:
					sw.Param = hd.Param
					if hd.Goto == model.GotoColumnStage {
						paramInt, _ := strconv.Atoi(hd.Param)
						if c, ok := s.columnListCache[paramInt]; ok {
							cidStr := strconv.Itoa(c.Ceid)
							sw.URI = model.FillURICategory(hd.Goto, cidStr, hd.Param)
						}
						sw.Goto = model.GotoColumn
					} else {
						sw.Goto = hd.Goto
						sw.URI = hd.URi
					}
					if sisLen := len(sis); sisLen > 1 {
						if sisLen%2 != 0 {
							sis = sis[:sisLen-1]
						}
					} else {
						continue
					}
				case model.GotoSpecialCard:
					sw.Cover = ""
				case model.GotoTagCard:
					if hd.TagID > 0 {
						var tagIDInt int64
						sw.Title, tagIDInt = s.fromTagIDByName(c, hd.TagID, now)
						sw.Param = strconv.FormatInt(tagIDInt, 10)
						sw.Goto = model.GotoTagID
					}
				}
				sw.Body = sis
				res[key] = append(res[key], sw)
			}
		}
	}
	return
}

// fillCardItem
func (s *Service) fillCardItem(csi *card.Content, tsi map[int64]*region.ShowItem) (si *region.ShowItem) {
	si = &region.ShowItem{}
	switch csi.Type {
	case model.CardGotoAv:
		si.Goto = model.GotoAv
		si.Param = csi.Value
	}
	si.URI = model.FillURI(si.Goto, si.Param, nil)
	if si.Goto == model.GotoAv {
		aid, err := strconv.ParseInt(si.Param, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", si.Param, err)
		} else {
			if it, ok := tsi[aid]; ok {
				si = it
				if csi.Title != "" {
					si.Title = csi.Title
				}
			} else {
				si = &region.ShowItem{}
			}
		}
	}
	return
}

// fromTagIDByName from tag_id by tag_name
func (s *Service) fromTagIDByName(ctx context.Context, tagID int, now time.Time) (tagName string, tagIDInt int64) {
	tag, err := s.tag.TagInfo(ctx, 0, tagID, now)
	if err != nil {
		log.Error("s.tag.TagInfo(%d) error(%v)", tagID, err)
		return
	}
	tagName = tag.Name
	tagIDInt = tag.Tid
	return
}

// upCahce update cache
func (s *Service) upCache(new, newOsea, old, oldOsea []*region.ShowItem) (res, resOsea []*region.ShowItem) {
	if len(new) > 0 {
		res = new
	} else {
		res = old
	}
	if len(newOsea) > 0 {
		resOsea = newOsea
	} else {
		resOsea = oldOsea
	}
	return
}

// upAidsCache update aids  cache
func (s *Service) upAidsCache(new, old []int64) (aids []int64) {
	if len(new) > 0 {
		aids = new
	} else {
		aids = old
	}
	return
}
