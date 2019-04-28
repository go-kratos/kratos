package service

import (
	"context"
	"sort"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	_blankList     = &model.List{ID: -1}
	_blankArtList  = int64(-1)
	_blankListArts = []*model.ListArtMeta{&model.ListArtMeta{ID: -1}}
)

func (s *Service) rawListArticles(c context.Context, listID int64) (res []*model.ListArtMeta, err error) {
	if listID <= 0 {
		return
	}
	arts, err := s.dao.CreativeListArticles(c, listID)
	if err != nil {
		return
	}
	var ids []int64
	for _, art := range arts {
		ids = append(ids, art.ID)
	}
	metas, err := s.dao.CreativeArticles(c, ids)
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

// List .
func (s *Service) List(c context.Context, id int64) (res *model.List, err error) {
	res, err = s.dao.List(c, id)
	res.FillDefaultImage(s.c.Article.ListDefaultImage)
	return
}

// ListArticles get list passed articles
func (s *Service) ListArticles(c context.Context, id, mid int64) (res *model.ListArticles, err error) {
	if id <= 0 {
		return
	}
	res = &model.ListArticles{}
	res.List, err = s.List(c, id)
	if err != nil {
		return
	}
	if res.List == nil {
		err = ecode.NothingFound
		return
	}
	var lastID int64
	group := &errgroup.Group{}
	group.Go(func() (err error) {
		arts, err := s.dao.ListArts(c, id)
		if err != nil {
			return
		}
		for _, a := range arts {
			if a.IsNormal() {
				res.Articles = append(res.Articles, a)
			}
		}
		return
	})
	group.Go(func() (err error) {
		res.Author, err = s.author(c, res.List.Mid)
		return
	})
	group.Go(func() (err error) {
		res.List.Read, err = s.dao.CacheListReadCount(c, id)
		return
	})
	group.Go(func() (err error) {
		if mid > 0 {
			res.Attention, _ = s.isAttention(c, mid, res.List.Mid)
		}
		return
	})
	group.Go(func() (err error) {
		if mid == 0 {
			return
		}
		lastID, _ = s.historyPosition(c, mid, res.List.ID)
		return
	})
	err = group.Wait()
	if res.Articles == nil {
		res.Articles = []*model.ListArtMeta{}
	}
	res.List.ArticlesCount = int64(len(res.Articles))
	res.Last.Strong()
	if lastID == 0 {
		return
	}
	for _, a := range res.Articles {
		if a.ID == lastID {
			res.Last = *a
			break
		}
	}
	return
}

// WebListArticles .
func (s *Service) WebListArticles(c context.Context, id, mid int64) (res *model.WebListArticles, err error) {
	arts, err := s.ListArticles(c, id, mid)
	res = &model.WebListArticles{Attention: arts.Attention, Author: arts.Author, List: arts.List, Last: arts.Last}
	var ids []int64
	for _, a := range arts.Articles {
		ids = append(ids, a.ID)
	}
	var stats map[int64]*model.Stats
	var states map[int64]int8
	group, ctx := errgroup.WithContext(c)
	group.Go(func() (err error) {
		if len(ids) == 0 {
			return
		}
		stats, _ = s.stats(ctx, ids)
		return
	})
	group.Go(func() (err error) {
		if mid == 0 || len(ids) == 0 {
			return
		}
		states, _ = s.HadLikesByMid(ctx, mid, ids)
		return
	})
	group.Wait()
	for _, a := range arts.Articles {
		art := &model.FullListArtMeta{ListArtMeta: a}
		if s.categoriesMap[art.Category.ID] != nil {
			art.Category = s.categoriesMap[art.Category.ID]
			art.Categories = s.categoryParents[art.Category.ID]
		}
		s := stats[art.ID]
		if s != nil {
			art.Stats = *s
		}
		if states != nil {
			art.LikeState = states[art.ID]
		}
		res.Articles = append(res.Articles, art)
	}
	return
}

// RebuildListCache rebuild list cache
func (s *Service) rebuildArticleListCache(c context.Context, aid int64) (listID int64, err error) {
	s.deleteArtsListCache(c, aid)
	lists, _ := s.dao.RawArtsListID(c, []int64{aid})
	listID = lists[aid]
	if listID > 0 {
		err = s.dao.SetArticleListCache(c, listID, []*model.ListArtMeta{&model.ListArtMeta{ID: aid}})
	}
	return
}

// RebuildListCache rebuild list cache
func (s *Service) RebuildListCache(c context.Context, id int64) (err error) {
	err = s.updateListCache(c, id)
	if err != nil {
		return
	}
	arts, err := s.dao.RawListArts(c, id)
	if err != nil {
		return
	}
	if err = s.updateListArtsCache(c, id, arts); err != nil {
		return
	}
	if err = s.updateArtListCache(c, id, arts); err != nil {
		return
	}
	list, err := s.dao.RawList(c, id)
	if err != nil {
		return
	}
	if list == nil {
		err = ecode.NothingFound
		return
	}
	if err = s.dao.RebuildUpListsCache(c, list.Mid); err != nil {
		return
	}
	if err = s.dao.RebuildListReadCountCache(c, id); err != nil {
		return
	}
	return
}

func (s *Service) updateListCache(c context.Context, id int64) (err error) {
	list, err := s.dao.RawList(c, id)
	if err != nil {
		return
	}
	err = s.dao.AddCacheList(c, id, list)
	return
}

func (s *Service) updateListArtsCache(c context.Context, id int64, arts []*model.ListArtMeta) (err error) {
	if arts == nil {
		arts, err = s.dao.RawListArts(c, id)
		if err != nil {
			return
		}
	}
	if len(arts) == 0 {
		err = s.deleteListArtsCache(c, id)
		return
	}
	err = s.dao.AddCacheListArts(c, id, arts)
	return
}

func (s *Service) updateArtListCache(c context.Context, id int64, arts []*model.ListArtMeta) (err error) {
	if arts == nil {
		arts, err = s.dao.RawListArts(c, id)
		if err != nil {
			return
		}
	}
	err = s.dao.SetArticleListCache(c, id, arts)
	return
}

func (s *Service) deleteListArtsCache(c context.Context, listID int64) (err error) {
	err = s.dao.AddCacheListArts(c, listID, _blankListArts)
	return
}

func (s *Service) deleteListCache(c context.Context, listID int64) (err error) {
	l := map[int64]*model.List{listID: _blankList}
	err = s.dao.AddCacheLists(c, l)
	return
}

func (s *Service) deleteArtsListCache(c context.Context, ids ...int64) (err error) {
	if len(ids) == 0 {
		return
	}
	m := make(map[int64]int64)
	for _, id := range ids {
		m[id] = _blankArtList
	}
	if err = s.dao.AddCacheArtsListID(c, m); err != nil {
		log.Errorv(c, log.KV("log", "deleteArtsListCache"), log.KV("err", err), log.KV("arg", ids))
	}
	return
}

// ListInfo list info
func (s *Service) ListInfo(c context.Context, aid int64) (res *model.ListInfo, err error) {
	lists, err := s.dao.ArtsList(c, []int64{aid})
	if err != nil {
		return
	}
	list := lists[aid]
	if list == nil {
		err = ecode.NothingFound
		return
	}
	res = &model.ListInfo{List: list}
	var articles []*model.ListArtMeta
	arts, err := s.dao.ListArts(c, list.ID)
	if err != nil {
		return
	}
	for _, a := range arts {
		if a.IsNormal() {
			articles = append(articles, a)
		}
	}
	res.Total = len(articles)
	for i, l := range articles {
		if l.ID == aid {
			res.Now = i
			if i-1 >= 0 {
				res.Last = articles[i-1]
			}
			if (i + 1) < len(articles) {
				res.Next = articles[i+1]
			}
			break
		}
	}
	res.List.FillDefaultImage(s.c.Article.ListDefaultImage)
	res.List.Read, err = s.dao.CacheListReadCount(c, list.ID)
	return
}

// Lists .
func (s *Service) Lists(c context.Context, keys []int64) (res map[int64]*model.List, err error) {
	res, err = s.dao.Lists(c, keys)
	for _, l := range res {
		l.FillDefaultImage(s.c.Article.ListDefaultImage)
	}
	return
}

// UpLists .
func (s *Service) UpLists(c context.Context, mid int64, sortType int8) (res model.UpLists, err error) {
	lists, err := s.dao.UpLists(c, mid)
	if err != nil {
		return
	}
	lmap, err := s.Lists(c, lists)
	if err != nil {
		return
	}
	arts, err := s.dao.ListsArts(c, lists)
	if err != nil {
		return
	}
	counts, err := s.dao.CacheListsReadCount(c, lists)
	for _, l := range lists {
		if lmap[l] != nil {
			list := lmap[l]
			if arts[l] != nil {
				list.ArticlesCount = int64(len(arts[l]))
			}
			list.Read = counts[l]
			res.Lists = append(res.Lists, list)
		}
	}
	if sortType == model.ListSortView {
		sort.Slice(res.Lists, func(i, j int) bool { return res.Lists[i].Read > res.Lists[j].Read })
	} else {
		sort.Slice(res.Lists, func(i, j int) bool { return res.Lists[i].PublishTime > res.Lists[j].PublishTime })
	}
	res.Total = len(res.Lists)
	return
}

// rebuildAllListReadCount  .
func (s *Service) rebuildAllListReadCount(c context.Context) error {
	i := 0
	for {
		lists, err := s.dao.RawAllListsEx(c, i, 1000)
		if err != nil {
			log.Errorv(c, log.KV("RebuildAllReadCount err", err))
			return err
		}
		if lists == nil {
			return nil
		}
		for _, list := range lists {
			err = s.dao.RebuildListReadCountCache(c, list.ID)
			if err != nil {
				log.Errorv(c, log.KV("RebuildAllReadCount err", err), log.KV("id", list.ID))
			} else {
				log.Infov(c, log.KV("RebuildAllReadCount", list.ID))
			}
		}
		i += 1000
	}
}

// RebuildAllListReadCount  .
func (s *Service) RebuildAllListReadCount(c context.Context) {
	go s.rebuildAllListReadCount(context.TODO())
}

// UpArtMetasAndLists .
func (s *Service) UpArtMetasAndLists(c context.Context, mid int64, pn int, ps int, sortType int) (res model.UpArtMetasLists, err error) {
	metas, err := s.UpArticleMetas(c, mid, pn, ps, sortType)
	if err != nil {
		return
	}
	res.UpArtMetas = metas
	res.UpLists, err = s.UpLists(c, mid, model.ListSortPtime)
	if res.UpArtMetas == nil {
		res.UpArtMetas = &model.UpArtMetas{}
	}
	if res.UpArtMetas.Articles == nil {
		res.UpArtMetas.Articles = []*model.Meta{}
	}
	if res.UpLists.Lists == nil {
		res.UpLists.Lists = []*model.List{}
	}
	return
}
