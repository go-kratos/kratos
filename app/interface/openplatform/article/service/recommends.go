package service

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"time"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/library/sync/errgroup"
)

var (
	_recommendCategory = int64(0)
)

// Recommends list recommend arts by category id
func (s *Service) Recommends(c context.Context, cid int64, pn, ps int, lastAids []int64, sort int) (res []*model.RecommendArt, err error) {
	var (
		start = (pn - 1) * ps
		// 多取一些用于去重
		end           = start + ps - 1 + len(lastAids)
		rems          = make(map[int64]*model.Recommend)
		allIDs, aids  []int64
		metas         map[int64]*model.Meta
		aidsm         map[int64]bool
		withRecommend bool
	)
	if cid != _recommendCategory {
		if (s.categoriesMap == nil) || (s.categoriesMap[cid] == nil) {
			err = ecode.RequestErr
			return
		}
	}
	if sort == model.FieldDefault {
		withRecommend = true
		sort = model.FieldNew
	}
	allRec := s.recommendAids[cid]
	var recommends [][]*model.Recommend
	if cid == _recommendCategory {
		recommends = s.genRecommendArtFromPool(s.RecommendsMap[cid], s.c.Article.RecommendRegionLen)
	} else {
		recommends = s.RecommendsMap[cid]
	}
	recommendsLen := len(recommends)
	// 只是最新文章 无推荐
	if (start >= recommendsLen) || !withRecommend {
		if (cid == _recommendCategory) && !s.setting.ShowRecommendNewArticles {
			return
		}
		var (
			nids        []int64
			newArtStart = start
			newArtEnd   = end
		)
		if withRecommend {
			newArtStart = start - recommendsLen
			newArtEnd = end - recommendsLen
		}
		nids, _ = s.dao.SortCache(c, cid, sort, newArtStart, newArtEnd)
		if withRecommend {
			allIDs = uniqIDs(nids, allRec)
		} else {
			allIDs = nids
		}
	} else {
		aids, rems = s.dealRecommends(recommends)
		if end < recommendsLen {
			allIDs = aids[start : end+1]
		} else {
			if (cid == _recommendCategory) && !s.setting.ShowRecommendNewArticles {
				allIDs = aids[start:]
			} else {
				// 混合推荐和最新文章
				var (
					nids []int64
					rs   = aids[start:]
				)
				newArtStart := 0
				newArtEnd := (end - start) - len(rs)
				nids, _ = s.dao.SortCache(c, cid, sort, newArtStart, newArtEnd)
				nids = uniqIDs(nids, allRec)
				allIDs = append(rs, nids...)
			}
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
		art := &model.RecommendArt{Meta: *metas[id]}
		if rems[id] != nil {
			art.Recommend = *rems[id]
		}
		res = append(res, art)
	}
	//截断分页数据
	if ps > len(res) {
		ps = len(res)
	}
	res = res[:ps]
	if cid == _recommendCategory {
		sortRecs(res)
	}
	return
}

func (s *Service) dealRecommends(recommends [][]*model.Recommend) (aids []int64, rems map[int64]*model.Recommend) {
	rems = make(map[int64]*model.Recommend)
	for _, recs := range recommends {
		rec := &model.Recommend{}
		*rec = *recs[s.randPosition(len(recs))]
		aids = append(aids, rec.ArticleID)
		// 不在推荐大图时间 去掉大图
		if rec.RecImageURL != "" {
			var now = time.Now().Unix()
			if (now < rec.RecImageStartTime) || (now > rec.RecImageEndTime) {
				rec.RecImageURL = ""
				rec.RecFlag = false
			}
		}
		if rec.RecFlag {
			rec.RecText = "编辑推荐"
		}
		rems[rec.ArticleID] = rec
		// 推荐文章id置空
		rec.ArticleID = 0
	}
	return
}

// 过滤禁止分区投稿
func filterNoRegionArts(as map[int64]*model.Meta) {
	for id, a := range as {
		if (a != nil) && a.AttrVal(model.AttrBitNoRegion) {
			delete(as, id)
		}
	}
}

// 按照发布时间排序
func sortRecs(res []*model.RecommendArt) {
	if len(res) == 0 {
		return
	}
	var len int
	for i, v := range res {
		if v.Rec {
			len = i
		}
	}
	sort.Slice(res[:len+1], func(i, j int) bool { return res[i].PublishTime > res[j].PublishTime })
}

// array a - array b
func uniqIDs(a []int64, b []int64) (res []int64) {
	bm := make(map[int64]bool)
	for _, v := range b {
		bm[v] = true
	}
	for _, v := range a {
		if !bm[v] {
			res = append(res, v)
		}
	}
	return
}

// UpdateRecommends update recommends
func (s *Service) UpdateRecommends(c context.Context) (err error) {
	var (
		recommendsMap = make(map[int64][][]*model.Recommend)
		recommendAids = make(map[int64][]int64)
		mutex         = &sync.Mutex{}
	)
	group, ctx := errgroup.WithContext(c)
	group.Go(func() error {
		recommends, err1 := s.dao.RecommendByCategory(ctx, _recommendCategory)
		if err1 != nil {
			return err1
		}
		// 推荐分区无位置 为推荐池
		rs := [][]*model.Recommend{recommends}
		mutex.Lock()
		recommendsMap[_recommendCategory] = rs
		mutex.Unlock()
		return nil
	})
	for _, category := range s.categoriesMap {
		category := category
		group.Go(func() error {
			recommends, err1 := s.dao.RecommendByCategory(ctx, category.ID)
			if err1 != nil {
				return err1
			}
			rs := calculateRecommends(recommends)
			mutex.Lock()
			recommendsMap[category.ID] = rs
			mutex.Unlock()
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return
	}
	s.RecommendsMap = recommendsMap
	for cid, v := range recommendsMap {
		for _, vv := range v {
			for _, vvv := range vv {
				recommendAids[cid] = append(recommendAids[cid], vvv.ArticleID)
			}
		}
	}
	s.recommendAids = recommendAids
	log.Info("s.UpdateRecommends success! len:(%v)", len(recommendsMap))
	return
}

func calculateRecommends(rs []*model.Recommend) (res [][]*model.Recommend) {
	m := make(map[int][]*model.Recommend)
	// 位置去重+ 随机选择
	for _, r := range rs {
		if r == nil {
			continue
		}
		if len(m[r.Position]) == 0 {
			m[r.Position] = append(m[r.Position], r)
		} else {
			var endTime bool
			for _, x := range m[r.Position] {
				if x.EndTime != 0 {
					endTime = true
					break
				}
			}
			if endTime {
				if r.EndTime == 0 {
					continue
				} else {
					m[r.Position] = append(m[r.Position], r)
				}
			} else {
				if r.EndTime == 0 {
					m[r.Position] = append(m[r.Position], r)
				} else {
					m[r.Position] = []*model.Recommend{r}
				}
			}
		}
	}
	for _, recommends := range m {
		res = append(res, recommends)
	}
	sort.Sort(model.Recommends(res))
	return
}

func (s *Service) randPosition(max int) (res int) {
	res = rand.Intn(max)
	return
}

func (s *Service) genRecommendArtFromPool(rs [][]*model.Recommend, recLen int) (res [][]*model.Recommend) {
	var pool []*model.Recommend
	if len(rs) > 0 {
		pool = rs[0]
	}
	if len(pool) == 0 {
		return
	}
	recs := append([]*model.Recommend{}, pool...)
	for i := range recs {
		j := rand.Intn(i + 1)
		recs[i], recs[j] = recs[j], recs[i]
	}
	if len(recs) < recLen {
		recLen = len(recs)
	}
	for _, r := range recs[:recLen] {
		res = append(res, []*model.Recommend{r})
	}
	return
}

// DelRecommendArtCache delete recommend article cache
func (s *Service) DelRecommendArtCache(c context.Context, aid, cid int64) (err error) {
	s.DelRecommendArt(_recommendCategory, aid)
	if cid, err = s.CategoryToRoot(cid); err != nil {
		dao.PromError("cache:删除文章推荐缓存")
		log.Error("s.DelRecommendArtCache.RootCategory(c, %v, %v) err: %+v", aid, cid, err)
		return
	}
	s.DelRecommendArt(cid, aid)
	return
}

// DelRecommendArt delete recommend article
func (s *Service) DelRecommendArt(categoryID int64, aid int64) {
	select {
	case s.recommendChan <- [2]int64{categoryID, aid}:
	default:
		dao.PromError("recommends:删除推荐文章 chan已满")
		log.Error("s.DelRecommendArt(%v, %v) chan full!", categoryID, aid)
	}
}

func (s *Service) deleteRecommendproc() {
	for {
		info, ok := <-s.recommendChan
		if !ok {
			return
		}
		if s.RecommendsMap == nil {
			continue
		}
		categoryID, aid := info[0], info[1]
		newRecommendsMap := map[int64][][]*model.Recommend{}
		for cid, rss := range s.RecommendsMap {
			if cid != categoryID {
				newRecommendsMap[cid] = rss
				continue
			}
			var newRecommends [][]*model.Recommend
			for _, rs := range rss {
				var newRs []*model.Recommend
				for _, r := range rs {
					if r.ArticleID != aid {
						newRs = append(newRs, r)
					}
				}
				if len(newRs) > 0 {
					newRecommends = append(newRecommends, newRs)
				}
			}
			newRecommendsMap[cid] = newRecommends
		}
		s.RecommendsMap = newRecommendsMap
	}
}

// RecommendHome recommend home
func (s *Service) RecommendHome(c context.Context, plat int8, build int, pn, ps int, aids []int64, mid int64, ip string, t time.Time, buvid string) (res *model.RecommendHome, sky *model.SkyHorseResp, err error) {
	res = &model.RecommendHome{IP: ip, Categories: s.primaryCategories}
	plus, sky, err := s.RecommendPlus(c, _recommendCategory, plat, build, pn, ps, aids, mid, t, model.FieldDefault, buvid)
	if plus != nil {
		res.RecommendPlus = *plus
	}
	return
}

// RecommendPlus recommend plus
func (s *Service) RecommendPlus(c context.Context, cid int64, plat int8, build int, pn, ps int, aids []int64, mid int64, t time.Time, sort int, buvid string) (res *model.RecommendPlus, sky *model.SkyHorseResp, err error) {
	res = &model.RecommendPlus{Banners: []*model.Banner{}, Articles: []*model.RecommendArtWithLike{}, Ranks: []*model.RankMeta{}, Hotspots: []*model.Hotspot{}}
	var group *errgroup.Group
	group, _ = errgroup.WithContext(c)
	group.Go(func() error {
		var arts []*model.RecommendArtWithLike
		if arts, sky, err = s.SkyHorse(c, cid, pn, ps, aids, sort, mid, build, buvid, plat); err == nil {
			res.Articles = arts
		}
		return nil
	})
	group.Go(func() error {
		if bs, e := s.Banners(c, plat, build, t); (e == nil) && (len(bs) > 0) {
			res.Banners = bs
		}
		return nil
	})
	group.Go(func() error {
		if s.setting.ShowHotspot {
			if hs, _ := s.dao.CacheHotspots(c); len(hs) > 0 {
				res.Hotspots = hs
			}
		}
		return nil
	})
	group.Go(func() error {
		if !s.setting.ShowAppHomeRank {
			return nil
		}
		if ranks, _, err := s.Ranks(c, model.RankWeek, mid, ""); (err == nil) && (len(ranks) > 0) {
			if len(ranks) > 3 {
				ranks = ranks[:3]
			}
			res.Ranks = ranks
		}
		return nil
	})
	group.Wait()
	return
}

// AllRecommends all recommends articles
func (s *Service) AllRecommends(c context.Context, pn, ps int) (count int64, res []*model.Meta, err error) {
	if pn < 1 {
		pn = 1
	}
	t := time.Now()
	count, _ = s.dao.AllRecommendCount(c, t)
	res = []*model.Meta{}
	ids, err := s.dao.AllRecommends(c, t, pn, ps)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	metas, err := s.ArticleMetas(c, ids)
	if err != nil {
		return
	}
	for _, id := range ids {
		if metas[id] != nil {
			res = append(res, metas[id])
		}
	}
	return
}

// SkyHorse .
func (s *Service) SkyHorse(c context.Context, cid int64, pn, ps int, lastAids []int64, sort int, mid int64, build int, buvid string, plat int8) (res []*model.RecommendArtWithLike, sky *model.SkyHorseResp, err error) {
	if (cid != _recommendCategory) || !s.skyHorseGray(buvid, mid) {
		res, err = s.RecommendsWithLike(c, cid, pn, ps, lastAids, sort, mid)
		return
	}
	var aids []int64
	var metas map[int64]*model.Meta
	var rems map[int64]*model.Recommend
	if pn == 1 {
		size := ps
		if size > s.c.Article.SkyHorseRecommendRegionLen {
			size = s.c.Article.SkyHorseRecommendRegionLen
		}
		recommends := s.genRecommendArtFromPool(s.RecommendsMap[_recommendCategory], size)
		aids, rems = s.dealRecommends(recommends)
	}
	if len(aids) < ps {
		sky, err = s.dao.SkyHorse(c, mid, build, buvid, plat, ps-len(aids))
		if (err != nil) || (len(sky.Data) == 0) {
			res, err = s.RecommendsWithLike(c, cid, pn, ps, lastAids, sort, mid)
			sky = nil
			return
		}
		for _, item := range sky.Data {
			if rems[item.ID] == nil {
				aids = append(aids, item.ID)
			}
		}
	}
	if metas, err = s.ArticleMetas(c, aids); err != nil {
		return
	}
	//过滤禁止显示的稿件
	filterNoDistributeArtsMap(metas)
	filterNoRegionArts(metas)
	states, _ := s.HadLikesByMid(c, mid, aids)
	for _, aid := range aids {
		if metas[aid] == nil {
			continue
		}
		art := model.RecommendArt{Meta: *metas[aid]}
		r := &model.RecommendArtWithLike{RecommendArt: art}
		if states != nil {
			r.LikeState = int(states[aid])
		}
		if rems[aid] != nil {
			r.Recommend = *rems[aid]
		}
		res = append(res, r)
	}
	return
}

func (s *Service) skyHorseGray(buvid string, mid int64) bool {
	if (mid == 0) && (buvid == "") {
		return false
	}
	for _, id := range s.c.Article.SkyHorseGrayUsers {
		if mid == id {
			return true
		}
	}
	for _, id := range s.c.Article.SkyHorseGray {
		if mid%10 == id {
			return true
		}
	}
	return false
}

func (s *Service) groupRecommend(c context.Context) (err error) {
	var (
		m     = make(map[int64]map[int64]bool)
		mutex = &sync.Mutex{}
	)
	for _, recommends := range s.RecommendsMap {
		var (
			rs   []*model.Recommend
			arts map[int64]*model.Meta
			aids = []int64{}
		)
		if len(recommends) > 0 {
			rs = recommends[0]
		}
		for _, r := range rs {
			aids = append(aids, r.ArticleID)
		}
		if arts, err = s.ArticleMetas(c, aids); err != nil || arts == nil {
			return
		}

		for _, art := range arts {
			if _, ok := m[art.Category.ID]; !ok {
				m[art.Category.ID] = make(map[int64]bool)
			}
			m[art.Category.ID][art.ID] = true
		}
	}
	mutex.Lock()
	s.RecommendsGroups = m
	mutex.Unlock()
	return
}

func (s *Service) getRecommentsGroups(c context.Context, cid int64, aid int64) (res []int64) {
	for i := range s.RecommendsGroups[cid] {
		if i != aid {
			res = append(res, i)
		}
	}
	return
}
