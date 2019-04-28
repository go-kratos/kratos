package service

import (
	"context"
	"strconv"
	"strings"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
)

var (
	_slideCount = 20
	_left       = 0
	_right      = 1
)

// var _slideRecommends = 2

// ViewList .
func (s *Service) ViewList(c context.Context, aid int64, buvid string, from string, ip string, build int, plat int8, mid int64) (pre int64, next int64) {
	var (
		cid  int64
		res  *model.Article
		err  error
		arts = &model.ArticleViewList{From: from, Buvid: buvid, Plat: plat, Build: build, Mid: mid}
	)
	if res, err = s.Article(c, aid); err != nil {
		res = nil
		return
	}
	if res == nil {
		res = nil
		return
	}
	if strings.Contains(from, "_") {
		fs := strings.Split(from, "_")
		if len(fs) < 2 {
			return
		}
		cid, _ = strconv.ParseInt(fs[1], 10, 64)
		switch fs[0] {
		case "category":
			//分区 天马
			if cid == 0 {
				arts.From = "recommend"
				pre, next = s.listFromRecommends(c, aid, arts)
				// arts.From = "skyhorse"
				// pre, next = s.listFromSkyHorseEx(c, aid, cid, arts, int64(res.PublishTime), res.Author.Mid)
			} else {
				arts.From = "default"
				pre, next = s.listByAuthor(c, aid, res.Author.Mid, int64(res.PublishTime), arts)
			}
		case "rank":
			//排行榜
			pre, next = s.listFromRank(c, aid, cid, arts, ip)
		case "readlist":
			//文集
			pre, next = s.listFromReadList(c, aid, arts)
		default:
			//default
			arts.From = "default"
			pre, next = s.listByAuthor(c, aid, res.Author.Mid, int64(res.PublishTime), arts)
		}
	} else {
		switch from {
		// case "mainCard":
		// 	//天马
		// 	arts.From = "skyhorse"
		// 	pre, next = s.listFromSkyHorseEx(c, aid, 0, arts, int64(res.PublishTime), res.Author.Mid)
		case "favorite":
			// TODO
			pre, next = s.validFavsList(c, aid, arts, true)
		case "records":
			// TODO
			// pre, next = 0, 0
			// pre, next = s.historyCursor(c, aid, arts, true)
		case "readlist":
			//文集
			pre, next = s.listFromReadList(c, aid, arts)
		case "articleSlide":
			//滑动
			pre, next, arts = s.listFromCache(c, aid, arts.Buvid)
		default:
			//作者其他的文章
			arts.From = "default"
			pre, next = s.listByAuthor(c, aid, res.Author.Mid, int64(res.PublishTime), arts)
		}
	}

	s.dao.AddCacheListArtsId(c, buvid, arts)
	return
}

// func (s *Service) listFromSkyHorseEx(c context.Context, aid int64, cid int64, arts *model.ArticleViewList, pt int64, authorMid int64) (pre int64, next int64) {
// 	pre, next = s.listFromSkyHorse(c, aid, cid, arts)
// 	if pre == 0 && next == 0 {
// 		arts.From = "default"
// 		pre, next = s.listByAuthor(c, aid, authorMid, pt, arts)
// 	}
// 	return
// }

func (s *Service) listFromCache(c context.Context, aid int64, buvid string) (pre int64, next int64, arts *model.ArticleViewList) {
	var (
		res      *model.Article
		err      error
		position = -1
	)
	if arts, err = s.dao.CacheListArtsId(c, buvid); err != nil {
		log.Error("s.dao.CacheListArtsId, error(%+v)", err)
		return
	}
	if arts == nil || arts.Aids == nil {
		return
	}
	for i, id := range arts.Aids {
		if id == aid {
			position = i
			break
		}
	}
	if position == -1 {
		return
	}
	switch arts.From {
	case "recommend":
		if position == 0 {
			if len(arts.Aids) > 1 {
				next = arts.Aids[1]
			}
			pre = s.newRecommend(c, arts, _left)
			// pre = s.newSkyHorse(c, arts, _left)
			arts.Position = position
			return
		}
		if position == len(arts.Aids)-1 {
			if len(arts.Aids) > 1 {
				pre = arts.Aids[len(arts.Aids)-2]
			}
			next = s.newRecommend(c, arts, _right)
			// next = s.newSkyHorse(c, arts, _right)
			arts.Position = position
			return
		}
	case "default":
		if position == 0 || position == len(arts.Aids)-1 {
			if res, err = s.Article(c, aid); err != nil || res == nil {
				return
			}
			pre, next = s.listByAuthor(c, aid, res.Author.Mid, int64(res.PublishTime), arts)
			arts.Position = position
			return
		}
		next = arts.Aids[position-1]
		pre = arts.Aids[position+1]
		return
	case "records":
		return
		// if position == len(arts.Aids)-1 {
		// 	if position > 0 {
		// 		pre = arts.Aids[position-1]
		// 	}
		// 	_, next = s.historyCursor(c, aid, arts, false)
		// 	return
		// }
	case "favorite":
		if position == len(arts.Aids)-1 {
			if position > 0 {
				pre = arts.Aids[position-1]
			}
			_, next = s.validFavsList(c, aid, arts, false)
			return
		}
	}
	if position > 0 {
		pre = arts.Aids[position-1]
	}
	if position < len(arts.Aids)-1 {
		next = arts.Aids[position+1]
	}
	return
}

func (s *Service) listFromReadList(c context.Context, aid int64, arts *model.ArticleViewList) (pre int64, next int64) {
	var (
		listID    int64
		list      *model.List
		artsMetas []*model.ListArtMeta
		ok        bool
	)
	lists, err := s.dao.ArtsList(c, []int64{aid})
	if err != nil {
		log.Error("s.dao.ArtsList, error(%+v)", err)
		return
	}
	if list, ok = lists[aid]; !ok {
		return
	}
	listID = list.ID

	if artsMetas, err = s.dao.ListArts(c, listID); err != nil {
		log.Error("s.dao.ListArts, error(%+v)", err)
		return
	}
	if artsMetas == nil {
		return
	}
	for i, art := range artsMetas {
		if art.ID == aid {
			arts.Position = i
		}
		arts.Aids = append(arts.Aids, art.ID)
	}
	if arts.Position > 0 {
		pre = artsMetas[arts.Position-1].ID
	}
	if arts.Position < len(artsMetas)-1 {
		next = artsMetas[arts.Position+1].ID
	}
	return
}

func (s *Service) listByAuthor(c context.Context, aid int64, mid int64, pt int64, arts *model.ArticleViewList) (pre int64, next int64) {
	var (
		beforeAids, afterAids, tmpAids []int64
		aidTimes                       [][2]int64
		exist                          bool
		j                              int
		err                            error
		metas                          map[int64]*model.Meta
		aids                           []int64
		addCache                       = true
		position                       = -1
	)
	if exist, err = s.dao.ExpireUpperCache(c, mid); err != nil {
		addCache = false
		err = nil
	} else if exist {
		if beforeAids, afterAids, err = s.dao.MoreArtsCaches(c, mid, int64(pt), _slideCount); err != nil {
			addCache = false
			exist = false
		}
		if len(beforeAids)+len(afterAids) == 0 {
			return
		}
		for i := len(beforeAids) - 1; i >= 0; i-- {
			tmpAids = append(tmpAids, beforeAids[i])
		}
		tmpAids = append(tmpAids, aid)
		tmpAids = append(tmpAids, afterAids...)
	} else {
		if aidTimes, err = s.dao.UpperPassed(c, mid); err != nil {
			log.Error("s.dao.UpperPassed, error(%+v)", err)
			return
		}
		if addCache {
			cache.Save(func() {
				s.dao.AddUpperCaches(context.TODO(), map[int64][][2]int64{mid: aidTimes})
			})
		}
		for i := len(aidTimes) - 1; i >= 0; i-- {
			aidTime := aidTimes[i]
			tmpAids = append(tmpAids, aidTime[0])
		}
	}
	if metas, err = s.ArticleMetas(c, tmpAids); err != nil {
		log.Error("s.ArticleMetas, error(%+v)", err)
		return
	}
	//过滤禁止显示的稿件
	filterNoDistributeArtsMap(metas)
	if len(metas) == 0 {
		return
	}
	for _, id := range tmpAids {
		if _, ok := metas[id]; !ok {
			continue
		}
		if id == aid {
			position = j
		}
		j++
		aids = append(aids, id)
	}
	if position == -1 {
		return
	}
	if position > 0 {
		next = aids[position-1]
	}
	if position < len(aids)-1 {
		pre = aids[position+1]
	}
	arts.Position = position
	arts.Aids = aids
	return
}

func (s *Service) listFromRank(c context.Context, aid int64, cid int64, arts *model.ArticleViewList, ip string) (pre int64, next int64) {
	var (
		exist    bool
		err      error
		aids     []int64
		rank     model.RankResp
		addCache = true
		position = -1
	)
	s.dao.DelCacheListArtsId(c, arts.Buvid)
	if !s.ranksMap[cid] {
		return
	}
	if exist, err = s.dao.ExpireRankCache(c, cid); err != nil {
		addCache = false
		err = nil
	}
	if exist {
		if rank, err = s.dao.RankCache(c, cid); err != nil {
			exist = false
			err = nil
			addCache = false
		}
	} else {
		if rank, err = s.dao.Rank(c, cid, ip); err != nil {
			if rank, err = s.dao.RankCache(c, cid); err != nil {
				log.Error("s.dao.RankCache, error(%+v)", err)
				return
			}
		} else {
			if addCache && len(rank.List) > 0 {
				cache.Save(func() {
					s.dao.AddRankCache(context.TODO(), cid, rank)
				})
			}
		}
	}
	if len(rank.List) == 0 {
		return
	}
	for i, a := range rank.List {
		aids = append(aids, a.Aid)
		if a.Aid == aid {
			position = i
		}
	}
	if position == -1 {
		return
	}
	arts.Position = position
	arts.Aids = aids
	if position > 0 {
		pre = aids[position-1]
	}
	if position < len(aids)-1 {
		next = aids[position+1]
	}
	return
}

// func (s *Service) listFromSkyHorse(c context.Context, aid int64, cid int64, arts *model.ArticleViewList) (pre int64, next int64) {
// 	var (
// 		beforeAids, afterAids, tmpAids, aids []int64
// 		err                                  error
// 		position                             = -1
// 		half                                 = _slideRecommends / 2
// 	)

// 	if tmpAids, err = s.filterAidsFromSkyHorse(c, aid, arts, _slideRecommends); err != nil || len(tmpAids) == 0 {
// 		return
// 	}

// 	if len(tmpAids) < half {
// 		half = len(tmpAids) / 2
// 	}
// 	beforeAids = make([]int64, half)
// 	afterAids = make([]int64, len(tmpAids)-half)
// 	copy(beforeAids, tmpAids[:half])
// 	copy(afterAids, tmpAids[half:])
// 	position = half
// 	aids = append([]int64{}, beforeAids...)
// 	aids = append(aids, aid)
// 	aids = append(aids, afterAids...)
// 	arts.Position = position
// 	arts.Aids = aids
// 	if len(beforeAids) > 0 {
// 		pre = beforeAids[len(beforeAids)-1]
// 	}
// 	if len(afterAids) > 0 {
// 		next = afterAids[0]
// 	}
// 	return
// }

// func (s *Service) newSkyHorse(c context.Context, arts *model.ArticleViewList, side int) (aid int64) {
// 	var (
// 		aids []int64
// 		err  error
// 	)
// 	if aids, err = s.filterAidsFromSkyHorse(c, arts.Aids[arts.Position], arts, _slideRecommends/2); err != nil || len(aids) == 0 {
// 		return
// 	}
// 	aids = uniqIDs(aids, arts.Aids)
// 	if len(aids) == 0 {
// 		return
// 	}
// 	if side == _left {
// 		aid = aids[len(aids)-1]
// 		arts.Aids = append(aids, arts.Aids...)
// 		arts.Position += len(aids)
// 	} else {
// 		aid = aids[0]
// 		arts.Aids = append(arts.Aids, aids...)
// 	}
// 	return
// }

// func (s *Service) filterAidsFromSkyHorse(c context.Context, aid int64, arts *model.ArticleViewList, size int) (aids []int64, err error) {
// 	var (
// 		tmpIds []int64
// 		sky    *model.SkyHorseResp
// 		metas  map[int64]*model.Meta
// 	)
// 	if sky, err = s.dao.SkyHorse(c, arts.Mid, arts.Build, arts.Buvid, arts.Plat, size); err != nil {
// 		return
// 	}
// 	if len(sky.Data) == 0 {
// 		return
// 	}
// 	for _, item := range sky.Data {
// 		tmpIds = append(tmpIds, item.ID)
// 	}
// 	if metas, err = s.ArticleMetas(c, tmpIds); err != nil {
// 		log.Error("s.ArticleMetas, error(%+v)", err)
// 		return
// 	}
// 	//过滤禁止显示的稿件
// 	filterNoDistributeArtsMap(metas)
// 	filterNoRegionArts(metas)
// 	if len(metas) == 0 {
// 		return
// 	}
// 	for _, meta := range metas {
// 		if meta.ID == aid {
// 			continue
// 		}
// 		aids = append(aids, meta.ID)
// 	}
// 	return
// }

func (s *Service) validFavsList(c context.Context, aid int64, arts *model.ArticleViewList, ok bool) (pre int64, next int64) {
	arts.Position++
	var (
		favs []*model.Favorite
		err  error
		page = &model.Page{Total: arts.Position*_slideCount + 1}
	)
	for ok && page.Total > arts.Position*_slideCount {
		if favs, page, err = s.Favs(c, arts.Mid, 0, arts.Position, _slideCount, ""); err != nil {
			return
		}
		for _, fav := range favs {
			if !fav.Valid {
				continue
			}
			arts.Aids = append(arts.Aids, fav.ID)
		}
		for i, id := range arts.Aids {
			if id != aid {
				continue
			}
			if i > 0 {
				pre = arts.Aids[i-1]
			}
			if i < len(arts.Aids)-1 {
				next = arts.Aids[i+1]
				ok = false
			}
		}
		arts.Position++
	}
	if next > 0 {
		return
	}
	ok = true
	for ok && page.Total > arts.Position*_slideCount {
		if favs, page, err = s.Favs(c, arts.Mid, 0, arts.Position, _slideCount, ""); err != nil {
			return
		}
		for _, fav := range favs {
			if !fav.Valid {
				continue
			}
			if next == 0 {
				next = fav.ID
			}
			ok = false
			arts.Aids = append(arts.Aids, fav.ID)
		}
		arts.Position++
	}
	return
}

// func (s *Service) historyCursor(c context.Context, aid int64, arts *model.ArticleViewList, ok bool) (pre int64, next int64) {
// 	var (
// 		res    []*history.Resource
// 		err    error
// 		viewAt = int64(arts.Position)
// 		aids   []int64
// 	)
// 	for ok {
// 		arg := &history.ArgCursor{Mid: arts.Mid, Businesses: []string{"article", "article-list"}, Ps: _slideCount, ViewAt: viewAt}
// 		if res, err = s.hisRPC.HistoryCursor(c, arg); err != nil || len(res) < 2 {
// 			return
// 		}
// 		for _, r := range res {
// 			viewAt = r.Unix
// 			id := r.Oid
// 			if r.TP == history.TypeCorpus {
// 				id = r.Cid
// 			}
// 			aids = append(aids, id)
// 		}
// 		for i, id := range aids {
// 			if id != aid {
// 				continue
// 			}
// 			ok = false
// 			if i > 0 {
// 				pre = aids[i-1]
// 			}
// 			if i < len(aids)-1 {
// 				next = aids[i+1]
// 			}
// 		}
// 		arts.Aids = append(arts.Aids, aids...)
// 		arts.Position = int(viewAt)
// 	}
// 	if next > 0 {
// 		return
// 	}
// 	arg := &history.ArgCursor{Mid: arts.Mid, Businesses: []string{"article", "article-list"}, Ps: _slideCount, ViewAt: viewAt}
// 	if res, err = s.hisRPC.HistoryCursor(c, arg); err != nil || len(res) == 0 {
// 		return
// 	}
// 	for _, r := range res {
// 		viewAt = r.Unix
// 		id := r.Oid
// 		if r.TP == history.TypeCorpus {
// 			id = r.Cid
// 		}
// 		if next == 0 {
// 			next = id
// 		}
// 		arts.Aids = append(arts.Aids, id)
// 	}
// 	arts.Position = int(viewAt)
// 	return
// }

func (s *Service) listFromRecommends(c context.Context, aid int64, arts *model.ArticleViewList) (pre int64, next int64) {
	recommends := s.genRecommendArtFromPool(s.RecommendsMap[_recommendCategory], _slideCount)
	tmpAids, _ := s.dealRecommends(recommends)
	var aids []int64
	for _, id := range tmpAids {
		if id != aid {
			aids = append(aids, id)
		}
	}
	if len(aids) == 0 {
		return
	}
	if len(aids) == 1 {
		next = aids[0]
		arts.Aids = append([]int64{aid}, aids[0])
		return
	}
	half := len(aids) / 2
	beforeAids := make([]int64, half)
	afterAids := make([]int64, len(aids)-half)
	copy(beforeAids, aids[:half])
	copy(afterAids, aids[half:])
	aids = append([]int64{}, beforeAids...)
	aids = append(aids, aid)
	aids = append(aids, afterAids...)
	arts.Position = half
	arts.Aids = aids
	pre = arts.Aids[half-1]
	next = arts.Aids[half+1]
	return
}

func (s *Service) newRecommend(c context.Context, arts *model.ArticleViewList, side int) (res int64) {
	var (
		m    = make(map[int64]bool)
		nids []int64
	)
	recommends := s.genRecommendArtFromPool(s.RecommendsMap[_recommendCategory], _slideCount)
	aids, _ := s.dealRecommends(recommends)
	if len(aids) == 0 {
		return
	}
	for _, aid := range arts.Aids {
		m[aid] = true
	}
	for _, aid := range aids {
		if !m[aid] {
			nids = append(nids, aid)
		}
	}
	if len(nids) == 0 {
		return
	}
	if side == _left {
		res = nids[len(nids)-1]
		arts.Aids = append(nids, arts.Aids...)
	}
	if side == _right {
		res = nids[0]
		arts.Aids = append(arts.Aids, nids...)
	}
	return
}
