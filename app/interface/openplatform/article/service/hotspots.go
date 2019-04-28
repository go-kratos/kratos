package service

import (
	"context"
	"time"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var _hotspotArtTime = time.Hour * 24 * 30

// UpdateHotspots update all hotspots
func (s *Service) UpdateHotspots(force bool) (err error) {
	var c = context.TODO()
	hotspots, err := s.dao.Hotspots(c)
	if err != nil {
		return
	}
	if len(hotspots) == 0 {
		err = s.dao.DelCacheHotspots(c)
		return
	}
	for _, hot := range hotspots {
		if err = s.genHotspot(c, hot, force); err != nil {
			dao.PromError("hotspots:生成")
			log.Error("hotspots s.genHotspot(%v, %v) err:%v", hot.Tag, force, err)
			return
		}
	}
	err = s.dao.AddCacheHotspots(c, hotspots)
	for _, h := range hotspots {
		s.dao.AddCacheHotspot(c, h.ID, h)
	}
	return
}

func (s *Service) genHotspot(c context.Context, hot *model.Hotspot, force bool) (err error) {
	var ok bool
	for _, typ := range model.HotspotTypes {
		ok, err = s.dao.ExpireHotspotArtsCache(c, typ, hot.ID)
		if err != nil {
			return
		}
		if !ok || force {
			break
		}
	}
	if ok && !force {
		//不重新生成 赋值计数 返回
		if s, _ := s.dao.CacheHotspot(c, hot.ID); s != nil {
			hot.Stats = s.Stats
		}
		return
	}
	ptime := time.Now().Add(-_hotspotArtTime)
	arts, err := s.dao.SearchArts(c, ptime.Unix())
	if err != nil {
		return
	}
	// filter tags && remove top arts
	tops := make(map[int64]bool)
	for _, x := range hot.TopArticles {
		tops[x] = true
	}
	var newArts []*model.SearchArt
	for _, art := range arts {
		for _, t := range art.Tags {
			if t == hot.Tag {
				hot.Stats.Read += art.StatsView
				hot.Stats.Reply += art.StatsReply
				if !tops[art.ID] {
					newArts = append(newArts, art)
				}
				break
			}
		}
	}
	arts = newArts
	if len(arts) == 0 {
		return
	}
	// add cache
	for _, typ := range model.HotspotTypes {
		var as [][2]int64
		for _, art := range arts {
			as = append(as, hotspotValue(typ, art))
		}
		if err = s.dao.AddCacheHotspotArts(c, typ, hot.ID, as, true); err != nil {
			return
		}
	}
	return
}

func hotspotValue(typ int8, art *model.SearchArt) [2]int64 {
	switch typ {
	case model.HotspotTypeView:
		return [2]int64{art.ID, art.StatsView}
	case model.HotspotTypePtime:
		return [2]int64{art.ID, art.PublishTime}
	}
	return [2]int64{0, 0}
}

// AddCacheHotspotArt check article in hotspots list and add cache
func (s *Service) AddCacheHotspotArt(c context.Context, art *model.SearchArt) (err error) {
	if art.PublishTime < time.Now().Add(-_hotspotArtTime).Unix() {
		return
	}
	hots, all, err := s.tagsHotspots(c, art.Tags)
	if err != nil {
		dao.PromError("hotspots:AddCacheHotspotArt")
		return
	}
	if len(hots) == 0 {
		return
	}
	for _, hot := range hots {
		hot.Stats.Read += art.StatsView
		hot.Stats.Reply += art.StatsReply
		if err = s.addCacheHotspotArt(c, hot.ID, art); err != nil {
			dao.PromError("hotspots:addCacheHotspotArt")
			return
		}
	}
	err = s.dao.AddCacheHotspots(c, all)
	return
}

// DelCacheHotspotArt delete art from hotspots
func (s *Service) DelCacheHotspotArt(c context.Context, aid int64) (err error) {
	hots, err := s.dao.CacheHotspots(c)
	if err != nil {
		dao.PromError("hotspots:DelCacheHotspotArt")
		return
	}
	for _, hot := range hots {
		for _, typ := range model.HotspotTypes {
			if err = s.dao.DelHotspotArtsCache(c, typ, hot.ID, aid); err != nil {
				return
			}
		}
	}
	return
}

// tagsHotspots get hotspots form tags
func (s *Service) tagsHotspots(c context.Context, tags []string) (res, all []*model.Hotspot, err error) {
	all, err = s.dao.CacheHotspots(c)
	if err != nil {
		dao.PromError("hotspots:tagsHotspots")
		return
	}
	for _, hot := range all {
		for _, t := range tags {
			if t == hot.Tag {
				res = append(res, hot)
				break
			}
		}
	}
	return
}

func (s *Service) addCacheHotspotArt(c context.Context, hotID int64, art *model.SearchArt) (err error) {
	for _, typ := range model.HotspotTypes {
		var ok bool
		if ok, err = s.dao.ExpireHotspotArtsCache(c, typ, hotID); err != nil {
			return
		}
		if ok {
			if err = s.dao.AddCacheHotspotArts(c, typ, hotID, [][2]int64{hotspotValue(typ, art)}, false); err != nil {
				return
			}
		}
	}
	return
}

func (s *Service) metaToSearch(c context.Context, m *model.Meta) (res *model.SearchArt) {
	if m == nil {
		return
	}
	res = &model.SearchArt{
		ID:          m.ID,
		PublishTime: int64(m.PublishTime),
	}
	for _, t := range m.Tags {
		res.Tags = append(res.Tags, t.Name)
	}
	stats, _ := s.stat(c, m.ID)
	if stats != nil {
		res.StatsView = stats.View
		res.StatsReply = stats.Reply
	}
	return
}

// HotspotArts get hotspot articles
func (s *Service) HotspotArts(c context.Context, id int64, pn, ps int, lastAids []int64, sort int8, mid int64) (hotspot *model.Hotspot, res []*model.MetaWithLike, err error) {
	if pn <= 0 {
		pn = 1
	}
	var (
		start = (pn - 1) * ps
		// 多取一些用于去重
		end           = start + ps - 1 + len(lastAids)
		allIDs        []int64
		metas         map[int64]*model.Meta
		aidsm         map[int64]bool
		withRecommend bool
	)
	if hotspot, err = s.dao.CacheHotspot(c, id); err != nil {
		return
	}
	if hotspot == nil {
		err = ecode.NothingFound
		return
	}
	hotspot.Stats.Count, _ = s.dao.HotspotArtsCacheCount(c, sort, id)
	if sort == model.HotspotTypeView {
		withRecommend = true
	}
	recommendsLen := len(hotspot.TopArticles)
	// 只是最新文章 无推荐
	if (start >= recommendsLen) || !withRecommend {
		var (
			nids        []int64
			newArtStart = start
			newArtEnd   = end
		)
		if withRecommend {
			newArtStart = start - recommendsLen
			newArtEnd = end - recommendsLen
		}
		nids, _ = s.dao.HotspotArtsCache(c, sort, id, newArtStart, newArtEnd)
		if withRecommend {
			allIDs = uniqIDs(nids, hotspot.TopArticles)
		} else {
			allIDs = nids
		}
	} else {
		if end < recommendsLen {
			allIDs = hotspot.TopArticles[start : end+1]
		} else {
			// 混合推荐和最新文章
			var (
				nids []int64
				rs   = hotspot.TopArticles[start:]
			)
			newArtStart := 0
			newArtEnd := (end - start) - len(rs)
			nids, _ = s.dao.HotspotArtsCache(c, sort, id, newArtStart, newArtEnd)
			nids = uniqIDs(nids, hotspot.TopArticles)
			allIDs = append(rs, nids...)
		}
	}
	if len(allIDs) == 0 {
		return
	}
	if metas, err = s.ArticleMetas(c, allIDs); err != nil {
		return
	}
	//过滤禁止显示的稿件
	filterNoDistributeArtsMap(metas)
	filterNoRegionArts(metas)
	//填充数据
	aidsm = make(map[int64]bool, len(lastAids))
	for _, aid := range lastAids {
		aidsm[aid] = true
	}
	for _, id := range allIDs {
		if (metas == nil) || (metas[id] == nil) || aidsm[id] {
			continue
		}
		art := &model.MetaWithLike{Meta: *metas[id]}
		res = append(res, art)
	}
	//截断分页数据
	if ps > len(res) {
		ps = len(res)
	}
	res = res[:ps]
	// fill like state
	aids := []int64{}
	for _, m := range res {
		aids = append(aids, m.ID)
	}
	states, _ := s.HadLikesByMid(c, mid, aids)
	if states == nil {
		return
	}
	for _, m := range res {
		m.LikeState = int(states[m.ID])
	}
	return
}
