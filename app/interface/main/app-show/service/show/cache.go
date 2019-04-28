package show

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"

	clive "go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/card"
	recmod "go-common/app/interface/main/app-show/model/recommend"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/app/interface/main/app-show/model/show"
	creativeAPI "go-common/app/interface/main/creative/api"
	"go-common/app/service/main/archive/api"
	resource "go-common/app/service/main/resource/model"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
)

const (
	_blackUrl = "http://172.18.7.208/privatedata/reco-deny-arcs.json"
)

var (
	// 动画，音乐，舞蹈，游戏，科技，娱乐，鬼畜，电影，时尚, 生活，广告，国漫，影视，纪录片
	_tids       = []string{"1", "3", "129", "4", "36", "5", "119", "23", "155", "160", "11", "165", "167", "181", "177"}
	_emptyItems = []*show.Item{&show.Item{}, &show.Item{}, &show.Item{}, &show.Item{}}
)

// loadRcmmndCache load recommend cahce.
func (s *Service) loadRcmmndCache(now time.Time) {
	aids, err := s.rcmmnd.Hots(context.TODO())
	if err != nil {
		log.Error("s.rcmmnd.Hots(%v) error(%v)", now, err)
		return
	}
	if len(aids) > 200 {
		aids = aids[:200]
	}
	if len(aids) < 60 {
		return
	}
	tmp, tmpOsea := s.fromAids(context.TODO(), aids)
	if len(tmp) > 0 {
		s.rcmmndCache = tmp
	}
	if len(tmpOsea) > 0 {
		s.rcmmndOseaCache = tmpOsea
	}
	log.Info("loadRcmmndCache success")
}

// loadRegionCache load region cahce.
func (s *Service) loadRegionCache(now time.Time) {
	var (
		tmp         = map[string][]*show.Item{}
		tmpOsea     = map[string][]*show.Item{}
		tmpBg       = map[string][]*show.Item{}
		tmpBgOsea   = map[string][]*show.Item{}
		tmpBgEp     = map[string][]*show.Item{}
		tmpBgEpOsea = map[string][]*show.Item{}
		sids        = map[int32]*seasongrpc.CardInfoProto{}
	)
	for _, tid := range _tids {
		rs, err := s.rcmmnd.Region(context.TODO(), tid)
		if len(rs) > 8 {
			rs = rs[:8]
		}
		tidInt, _ := strconv.Atoi(tid)
		if err != nil || len(rs) < 4 {
			log.Error("s.rcmmnd.Region(%v) error(%v)", now, err)
			res, aids, err := s.dyn.RegionDynamic(context.TODO(), tidInt, 1, 8)
			if err != nil || len(res) < 4 {
				log.Error("s.dyn.RegionDynamic(%v) error(%v)", now, err)
				continue
			}
			if len(res) > 8 {
				res = res[:8]
			}
			if _, isBangumi := _bangumiReids[tidInt]; isBangumi {
				sids, _ = s.fromSeasonID(context.TODO(), aids)
			}
			tmp[tid], tmpOsea[tid] = s.fromArchivesPB(res)
			tmpBg[tid], tmpBgOsea[tid] = s.fromArchivesBangumi(context.TODO(), res, aids, sids, _bangumiSeasonID)
			tmpBgEp[tid], tmpBgEpOsea[tid] = s.fromArchivesBangumi(context.TODO(), res, aids, sids, _bangumiEpisodeID)
		} else {
			if _, isBangumi := _bangumiReids[tidInt]; isBangumi {
				sids, _ = s.fromSeasonID(context.TODO(), rs)
			}
			tmp[tid], tmpOsea[tid], tmpBg[tid], tmpBgOsea[tid], tmpBgEp[tid], tmpBgEpOsea[tid] =
				s.fromBgAids(context.TODO(), rs, sids, _bangumiSeasonID)
		}
		log.Info("load show index region(%s) success", tid)
	}
	if len(tmp) > 0 {
		s.regionCache = tmp
	}
	if len(tmpOsea) > 0 {
		s.regionOseaCache = tmpOsea
	}
	if len(tmpBg) > 0 {
		s.regionBgCache = tmpBg
	}
	if len(tmpBgOsea) > 0 {
		s.regionBgOseaCache = tmpBgOsea
	}
	if len(tmpBgEp) > 0 {
		s.regionBgEpCache = tmpBgEp
	}
	if len(tmpBgEpOsea) > 0 {
		s.regionBgEpOseaCache = tmpBgEpOsea
	}
	log.Info("loadRegionCache success")
}

// loadBannerCahce load banner cache.
func (s *Service) loadBannerCahce() {
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

// loadBgmCache load bangumi cache
func (s *Service) loadBgmCache(now time.Time) {
	bgms, err := s.bgm.Recommend(now)
	if err != nil {
		log.Error("s.bgm.Recommend(%v) error(%v)", now, err)
		return
	}
	if len(bgms) < 8 {
		return
	}
	var (
		tmp = map[int8][]*show.Item{}
		si  *show.Item
	)
	for _, bgm := range bgms {
		si = &show.Item{}
		si.FromBangumi(bgm)
		tmp[model.PlatAndroid] = append(tmp[model.PlatAndroid], si)
		tmp[model.PlatIPhone] = append(tmp[model.PlatIPhone], si)
		tmp[model.PlatAndroidG] = append(tmp[model.PlatAndroidG], si)
		tmp[model.PlatAndroidI] = append(tmp[model.PlatAndroidI], si)
		tmp[model.PlatIPhoneI] = append(tmp[model.PlatIPhoneI], si)
		si = &show.Item{}
		si.FromBangumi(bgm)
		si.Cover = bgm.Cover // pad hd get bangumi cover
		tmp[model.PlatIPad] = append(tmp[model.PlatIPad], si)
		tmp[model.PlatIPadI] = append(tmp[model.PlatIPadI], si)
		tmp[model.PlatAndroidTV] = append(tmp[model.PlatAndroidTV], si)
	}
	s.bgmCache = tmp
	log.Info("loadBgmCache success")
}

// loadLiveCache load live cache .
func (s *Service) loadLiveCache(now time.Time) {
	lrs, err := s.lv.Recommend(now)
	if err != nil {
		log.Error("s.live.Recommend(%v) error(%v)", now, err)
		return
	}
	if lrs == nil {
		return
	}
	s.liveCount = lrs.Count
	if subLen := len(lrs.Lives.Subject); subLen > 0 {
		tmp := make([]*show.Item, 0, subLen)
		for _, l := range lrs.Lives.Subject {
			si := &show.Item{}
			si.FromLive(l)
			tmp = append(tmp, si)
		}
		s.liveMoeCache = tmp
	}
	if hotLen := len(lrs.Lives.Hot); hotLen >= 4 {
		tmp := make([]*show.Item, 0, hotLen)
		for _, l := range lrs.Lives.Hot {
			si := &show.Item{}
			si.FromLive(l)
			tmp = append(tmp, si)
		}
		s.liveHotCache = tmp
	}
	log.Info("loadLiveCache success")
}

// loadShowCache load all show cache
func (s *Service) loadShowCache() {
	hdm, err := s.dao.Heads(context.TODO())
	if err != nil {
		log.Error("s.dao.Heads error(%v)", err)
		return
	}
	itm, err := s.dao.Items(context.TODO())
	if err != nil {
		log.Error("s.dao.Items error(%v)", err)
		return
	}
	tmp, tmpbg, tmpbgep := s.mergeShow(hdm, itm)
	if len(tmp) > 0 {
		s.cache = tmp
	}
	if len(tmpbg) > 0 {
		s.cacheBg = tmpbg
	}
	if len(tmpbgep) > 0 {
		s.cacheBgEp = tmpbgep
	}
	log.Info("loadShowCache success")
}

// loadShowTempCache load all show temp cache
func (s *Service) loadShowTempCache() {
	hdm, err := s.dao.TempHeads(context.TODO())
	if err != nil {
		log.Error("s.dao.TempHeads error(%v)", err)
		return
	}
	itm, err := s.dao.TempItems(context.TODO())
	if err != nil {
		log.Error("s.dao.TempItems error(%v)", err)
		return
	}
	s.tempCache, _, _ = s.mergeShow(hdm, itm)
	log.Info("loadShowTempCache success")
}

// loadRegionListCache
func (s *Service) loadRegionListCache() {
	res, err := s.rg.RegionPlat(context.TODO())
	if err != nil {
		log.Error("s.rg.RegionPlat error(%v)", err)
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

// loadRankAllCache
func (s *Service) loadRankAllCache() {
	var (
		rankAids   []int64
		err        error
		as         map[int64]*api.Arc
		c          = context.TODO()
		trankCache = []*rank.Rank{}
	)
	aids, others, scores, err := s.rcmmnd.RankAppAll(c)
	if err != nil {
		log.Error("s.rcmmnd.RankAppAll error(%v)", err)
		return
	}
	for _, aid := range aids {
		if _, ok := others[aid]; !ok {
			rankAids = append(rankAids, aid)
		}
	}
	if len(rankAids) == 0 {
		return
	}
	s.rankAidsCache = rankAids
	s.rankScoreCache = scores
	if as, err = s.arc.ArchivesPB(c, rankAids); err != nil {
		log.Error("s.arc.ArchivesPB aids(%v) error(%v)", aids, err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", rankAids)
		return
	}
	s.rankArchivesCache = as
	tmp, tmpOsea := s.fromRankAids(c, rankAids, scores, as)
	if len(tmp) > 0 {
		s.rankCache = tmp
	}
	if len(tmpOsea) > 0 {
		s.rankOseaCache = tmpOsea
	}
	log.Info("loadRankAllCache success")
	// new cache
	for _, aid := range rankAids {
		r := &rank.Rank{
			Aid:   aid,
			Score: int32(scores[aid]),
		}
		trankCache = append(trankCache, r)
	}
	s.rankCache2 = trankCache

}

// loadColumnListCache
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
	tmpItem := map[int]map[int64]*show.Item{}
	for recid, aid := range aids {
		tmpItem[recid] = s.fromCardAids(context.TODO(), aid)
	}
	tmp := s.mergeCard(context.TODO(), hdm, itm, tmpItem, now)
	s.cardCache = tmp
}

func (s *Service) mergeCard(c context.Context, hdm map[int8]map[int][]*card.Card, itm map[int][]*card.Content, tmpItems map[int]map[int64]*show.Item, now time.Time) (res map[string][]*show.Show) {
	var (
		_topic    = 1
		_activity = 0
	)
	res = map[string][]*show.Show{}
	for plat, phds := range hdm {
		hds, ok := phds[0]
		if !ok {
			continue
		}
		for _, hd := range hds {
			key := fmt.Sprintf(_initCardKey, plat)
			var (
				sis []*show.Item
			)
			its, ok := itm[hd.ID]
			if !ok {
				its = []*card.Content{}
			}
			tmpItem, ok := tmpItems[hd.ID]
			if !ok {
				tmpItem = map[int64]*show.Item{}
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
					} else {
						if act, ok := actm[topicID]; ok && act.H5Cover != "" && act.H5URL != "" {
							si := &show.Item{}
							si.FromTopic(act)
							sis = []*show.Item{si}
						}
					}
				}
			case 3:
				if topicID, err := strconv.ParseInt(hd.Rvalue, 10, 64); err == nil {
					if actm, err := s.act.Activitys(c, []int64{topicID}, _activity, ""); err != nil {
					} else {
						if act, ok := actm[topicID]; ok && act.H5Cover != "" && act.H5URL != "" {
							si := &show.Item{}
							si.FromActivity(act, now)
							sis = []*show.Item{si}
						}
					}
				}
			case 4:
				if len(s.rankCache) > 3 {
					sis = s.rankCache[:3]
				} else {
					sis = s.rankCache
				}
			case 5, 6, 8:
				for _, ci := range its {
					si := s.fillCardItem(ci, tmpItem)
					if si.Title != "" {
						sis = append(sis, si)
					}
				}
			case 7:
				si := &show.Item{
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
			sw := &show.Show{}
			sw.Head = &show.Head{
				CardID:    hd.ID,
				Title:     hd.Title,
				Type:      hd.TypeStr,
				Build:     hd.Build,
				Condition: hd.Condition,
				Plat:      hd.Plat,
				Style:     "small",
			}
			if hd.Cover != "" {
				sw.Head.Cover = hd.Cover
			}
			switch sw.Head.Type {
			case model.GotoDaily:
				sw.Head.Date = now.Unix()
				sw.Head.Param = hd.Rvalue
			case model.GotoCard:
				sw.Head.URI = hd.URi
				sw.Head.Goto = hd.Goto
				sw.Head.Param = hd.Param
			case model.GotoRank:
				sw.Head.Param = "all"
			case model.GotoTopic, model.GotoActivity:
				if sw.Head.Title == "" {
					if len(sis) > 0 {
						sw.Head.Title = sis[0].Title
					}
				}
			case model.GotoVeidoCard:
				sw.Head.Param = hd.Param
				if hd.Goto == model.GotoColumnStage {
					paramInt, _ := strconv.Atoi(hd.Param)
					if c, ok := s.columnListCache[paramInt]; ok {
						cidStr := strconv.Itoa(c.Ceid)
						sw.Head.URI = model.FillURICategory(hd.Goto, cidStr, hd.Param)
					}
					sw.Head.Goto = model.GotoColumn
				} else {
					sw.Head.Goto = hd.Goto
					sw.Head.URI = hd.URi
				}
				if sisLen := len(sis); sisLen > 1 {
					if sisLen%2 != 0 {
						sis = sis[:sisLen-1]
					}
				} else {
					continue
				}
			case model.GotoSpecialCard:
				sw.Head.Cover = ""
			case model.GotoTagCard:
				if hd.TagID > 0 {
					var tagIDInt int64
					sw.Head.Title, tagIDInt = s.fromTagIDByName(c, hd.TagID, now)
					sw.Head.Param = strconv.FormatInt(tagIDInt, 10)
					sw.Head.Goto = model.GotoTagID
				}
			}
			if len(sis) == 0 {
				sis = _emptyShowItems
			}
			sw.Body = sis
			res[key] = append(res[key], sw)
		}
	}
	return
}

func (s *Service) fillCardItem(csi *card.Content, tsi map[int64]*show.Item) (si *show.Item) {
	si = &show.Item{}
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
				si = &show.Item{}
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

// mergeShow merge heads and items
func (s *Service) mergeShow(hdm map[int8][]*show.Head, itm map[int][]*show.Item) (res, resbg, resbgep map[string][]*show.Show) {
	res = map[string][]*show.Show{}
	resbg = map[string][]*show.Show{}
	resbgep = map[string][]*show.Show{}
	for plat, hds := range hdm {
		for _, hd := range hds {
			key := fmt.Sprintf(_initShowKey, plat, hd.Language)
			its, ok := itm[hd.ID]
			if !ok {
				continue
			}
			var (
				sis     []*show.Item
				sisbg   []*show.Item
				sisbgep []*show.Item
				ext     *show.Ext
			)
			switch hd.Type {
			case "recommend":
				for _, si := range its {
					if si.IsRandom() {
						continue
					}
					s.fillItem(plat, si)
					sis = append(sis, si)
				}
				sisbg = sis
				sisbgep = sis
			case "live":
				if plat == model.PlatAndroidTV {
					continue
				}
				ext = &show.Ext{
					LiveCnt: s.liveCount,
				}
			case "bangumi":
				sis = s.sliceCache(plat, s.bgmCache[plat])
				sisbg = sis
				sisbgep = sis
			case "region":
				if model.IsOverseas(plat) {
					sis = s.sliceCache(plat, s.regionOseaCache[hd.Param])
					sisbg = s.sliceCache(plat, s.regionBgOseaCache[hd.Param])
					sisbgep = s.sliceCache(plat, s.regionBgEpOseaCache[hd.Param])
				} else {
					sis = s.sliceCache(plat, s.regionCache[hd.Param])
					sisbg = s.sliceCache(plat, s.regionBgCache[hd.Param])
					sisbgep = s.sliceCache(plat, s.regionBgEpCache[hd.Param])
				}
			case "sp":
				for _, si := range its {
					spidIdx := strings.Split(si.Param, ",")
					si.Goto = model.GotoSp
					si.Param = spidIdx[0]
					si.URI = model.FillURI(model.GotoSp, spidIdx[0], nil)
					if len(spidIdx) == 2 {
						si.Index = spidIdx[1]
					}
					sis = append(sis, si)
				}
				sisbg = sis
				sisbgep = sis
			case "activity":
				for _, si := range its {
					if si.IsRandom() {
						continue
					}
					si.Goto = model.GotoWeb
					si.URI = model.FillURI(model.GotoWeb, si.Param, nil)
					sis = append(sis, si)
				}
				sisbg = sis
				sisbgep = sis
			case "topic":
				for _, si := range its {
					si.Goto = model.GotoWeb
					si.URI = model.FillURI(model.GotoWeb, si.Param, nil)
					sis = append(sis, si)
				}
				sisbg = sis
				sisbgep = sis
			case "focus":
				for _, si := range its {
					if si.IsRandom() {
						continue
					}
					pp := strings.Split(si.Param, ",")
					si.Param = pp[0]
					if len(pp) == 2 {
						si.Goto = pp[1]
					} else {
						si.Goto = model.GotoAv
					}
					si.URI = model.FillURI(si.Goto, si.Param, nil)
					sisbg = append(sisbg, si)
				}
				sisbgep = sisbg
			default:
				continue
			}
			sw := &show.Show{}
			sw.Head = hd
			sw.Body = sis
			sw.Ext = ext
			swbg := &show.Show{}
			swbg.Head = hd
			swbg.Body = sisbg
			swbg.Ext = ext
			swbgep := &show.Show{}
			swbgep.Head = hd
			swbgep.Body = sisbgep
			swbgep.Ext = ext
			// append show.Show
			res[key] = append(res[key], sw)
			resbg[key] = append(resbg[key], swbg)
			resbgep[key] = append(resbgep[key], swbgep)
		}
	}
	return
}

// fillItem used by loadShowCache
func (s *Service) fillItem(plat int8, si *show.Item) {
	pp := strings.Split(si.Param, ",")
	si.Param = pp[0]
	if len(pp) == 2 {
		si.Goto = pp[1]
	} else {
		si.Goto = model.GotoAv
	}
	si.URI = model.FillURI(si.Goto, si.Param, nil)
	if si.Goto == model.GotoAv {
		aid, err := strconv.ParseInt(si.Param, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", si.Param, err)
		} else {
			a, err := s.arc.Archive(context.TODO(), aid)
			if err != nil || a == nil {
				log.Error("s.arc.Archive(%d) error(%v)", aid, err)
			} else {
				si.Play = int(a.Stat.View)
				si.Danmaku = int(a.Stat.Danmaku)
				if si.Title == "" {
					si.Title = a.Title
				}
				if si.Cover == "" {
					si.Cover = a.Pic
				}
			}
		}
	} else {
		si.Play = rand.Intn(1000)
		si.Danmaku = rand.Intn(1000)
	}
}

// sliceCache used by loadShowCache
func (s *Service) sliceCache(plat int8, chc []*show.Item) []*show.Item {
	if len(chc) == 0 {
		return _emptyItems
	}
	cnt := 4
	if plat == model.PlatIPad {
		cnt = 8
	}
	if len(chc) < cnt {
		cnt = len(chc)
	}
	return chc[:cnt]
}

func (s *Service) loadBlackCache() {
	var res []int64
	if err := s.client.Get(context.TODO(), _blackUrl, "", nil, &res); err != nil {
		log.Error("recommend ranking url(%s) error(%v)", _blackUrl, err)
		return
	}
	if len(res) == 0 {
		return
	}
	tmp := map[int64]struct{}{}
	for _, aid := range res {
		tmp[aid] = struct{}{}
	}
	s.blackCache = tmp
	log.Info("reBlackList success")
}

// rcmmndproc get recommend aids and add into cache.
func (s *Service) rcmmndproc() {
	var ctx = context.TODO()
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				r := <-s.rcmmndCh
				s.dao.AddRcmmndCache(ctx, r.key, r.aids...)
			}
		}()
	}
}

// // loadPopularCard load popular card
// func (s *Service) loadPopularCard(now time.Time) {
// 	var (
// 		c             = context.TODO()
// 		err           error
// 		tmp, tmpcache []*card.PopularCard
// 		tmpPlat       = map[int64]map[int8][]*card.PopularCardPlat{}
// 	)
// 	if tmp, err = s.cdao.Card(c, now); err != nil {
// 		log.Error("popular card s.cd.Card error(%v)", err)
// 		return
// 	}
// 	if tmpPlat, err = s.cdao.CardPlat(c); err != nil {
// 		log.Error("popular card s.cd.CardPlat error(%v)", err)
// 		return
// 	}
// 	for _, t := range tmp {
// 		tc := &card.PopularCard{}
// 		*tc = *t
// 		if pconfig, ok := tmpPlat[t.ID]; ok {
// 			tc.PopularCardPlat = pconfig
// 		}
// 		tmpcache = append(tmpcache, tc)
// 	}
// 	s.hotCache = tmpcache
// 	log.Info("hotCache success")
// }

func (s *Service) loadHotTenTabAids() {
	var tmpList = make(map[int][]*recmod.CardList)
	for i := 0; i < 10; i++ {
		var (
			c          = context.TODO()
			err        error
			hottabAids []*recmod.CardList
			flowResp   *creativeAPI.FlowResponse
			oids       []int64
			forbidAids = make(map[int64]struct{})
		)
		if hottabAids, err = s.rcmmnd.HotHeTongTabCard(c, i); err != nil {
			log.Error("%+v", err)
			continue
		}
		for _, hot := range hottabAids {
			if hot.Goto == model.GotoAv {
				oids = append(oids, hot.ID)
			}
		}
		if flowResp, err = s.creativeClient.FlowJudge(context.Background(), &creativeAPI.FlowRequest{
			Oids:     oids,
			Business: 4,
			Gid:      24,
		}); err != nil {
			log.Error("s.creativeClient.FlowJudge error(%v)", err)
			tmpList[i] = hottabAids
		} else {
			for _, oid := range flowResp.Oids {
				forbidAids[oid] = struct{}{}
			}
			for _, list := range hottabAids {
				if list.Goto == model.GotoAv {
					if _, ok := forbidAids[list.ID]; ok {
						log.Info("aid(%d) is flowJundged", list.ID)
						continue
					}
				}
				tmpList[i] = append(tmpList[i], list)
			}
		}
		log.Info("buildHotSuccess(%d) len(%d)", i, len(tmpList[i]))
	}
	if len(tmpList) == 10 {
		s.hotTenTabCardCache = tmpList
	}
}

func (s *Service) loadHotTopicCache() {
	var (
		c      = context.TODO()
		err    error
		topics []*clive.TopicHot
	)
	if topics, err = s.lv.TopicHots(c); err != nil {
		log.Error("topichots s.lv.TopicHots error(%v)", err)
		return
	}
	if len(topics) > 8 {
		s.hottopicsCache = topics[:8]
	} else {
		s.hottopicsCache = topics
	}
	log.Info("loadHotTopicCache success")
}

func (s *Service) loadHotTenMergeRcmdCache(i int) {
	// mc
	var (
		c        = context.TODO()
		hcards   []*recmod.CardList
		rcmdcard []*card.PopularCard
		ok       bool
		err      error
	)
	if hcards, ok = s.hotTenTabCardCache[i]; ok {
		for _, hcard := range hcards {
			rcmdcard = append(rcmdcard, hcard.CardListChange())
		}
		if err = s.cdao.AddPopularCardTenCache(c, i, rcmdcard); err != nil {
			log.Error("cards mc s.cdao.AddPopularCardCache error(%v)", err)
			return
		}
	}
}

// PopularCardList cards
func (s *Service) PopularCardTenList(c context.Context, i int) (res []*card.PopularCard) {
	var err error
	if res, err = s.cdao.PopularCardTenCache(c, i); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

func (s *Service) loadCardSetCache() {
	var (
		cards map[int64]*operate.CardSet
		err   error
	)
	if cards, err = s.cdao.CardSet(context.TODO()); err != nil {
		log.Error("%+v", err)
		return
	}
	s.cardSetCache = cards
}

func (s *Service) loadDynamicHotCache() {
	var (
		liveList []*clive.DynamicHot
		err      error
	)
	if liveList, err = s.lv.DynamicHot(context.TODO()); err != nil {
		log.Error("s.lv.dynamichot error(%v)", err)
		return
	}
	s.dynamicHotCache = liveList
}

func (s *Service) loadEventTopicCache() {
	var (
		eventtopic map[int64]*operate.EventTopic
		err        error
	)
	if eventtopic, err = s.cdao.EventTopic(context.TODO()); err != nil {
		log.Error("s.cdao.eventtopic error(%v)", err)
		return
	}
	s.eventTopicCache = eventtopic
}
