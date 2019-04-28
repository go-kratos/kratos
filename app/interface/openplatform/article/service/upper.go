package service

import (
	"context"
	"sort"
	"sync"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"

	"go-common/library/sync/errgroup"
)

// UpArticleMetas list up article metas
func (s *Service) UpArticleMetas(c context.Context, mid int64, pn int, ps int, sortType int) (res *artmdl.UpArtMetas, err error) {
	var (
		ups   map[int64][]*artmdl.Meta
		start = (pn - 1) * ps
		end   = start + ps - 1 // from cache, end-1
		metas []*artmdl.Meta
	)
	if sortType == artmdl.FieldDefault {
		ups, err = s.UpsArticleMetas(c, []int64{mid}, start, end)
	} else {
		ups, err = s.UpsArticleMetas(c, []int64{mid}, 0, -1)
	}
	if err != nil {
		return
	}
	metas = ups[mid]
	if sortType != artmdl.FieldDefault {
		end++
		switch sortType {
		case artmdl.FieldFav:
			sort.Slice(metas, func(i, j int) bool { return metas[i].Stats.Favorite > metas[j].Stats.Favorite })
		case artmdl.FieldView:
			sort.Slice(metas, func(i, j int) bool { return metas[i].Stats.View > metas[j].Stats.View })
		}
		if start > len(metas) {
			start = len(metas)
		}
		if end > len(metas) {
			end = len(metas)
		}
		metas = metas[start:end]
	}
	res = new(artmdl.UpArtMetas)
	res.Articles = filterNoDistributeArts(metas)
	res.Pn = pn
	res.Ps = ps
	if res.Count, err = s.UpperArtsCount(c, mid); err != nil {
		dao.PromError("upper:获取作者文章数")
	}
	return
}

// UpsArticleMetas list up article metas
func (s *Service) UpsArticleMetas(c context.Context, mids []int64, start int, end int) (res map[int64][]*artmdl.Meta, err error) {
	var (
		group = &errgroup.Group{}
		mutex = &sync.Mutex{}
	)
	res = make(map[int64][]*artmdl.Meta)
	upArtIDs, _ := s.upArtIDs(c, mids, start, end)
	for mid, ids := range upArtIDs {
		mid := mid
		ids := ids
		group.Go(func() (err error) {
			var (
				artsm map[int64]*artmdl.Meta
				arts  []*artmdl.Meta
			)
			artsm, _ = s.FeedArticleMetas(c, ids)
			for _, art := range artsm {
				arts = append(arts, art)
			}
			mutex.Lock()
			sort.Sort(artmdl.Metas(arts))
			res[mid] = arts
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	return
}

func (s *Service) upArtIDs(c context.Context, mids []int64, start, end int) (res map[int64][]int64, err error) {
	var (
		exists        map[int64]bool
		addCache      = true
		missMids      = make([]int64, 0, len(mids))
		cacheMids     = make([]int64, 0, len(mids))
		group         = &errgroup.Group{}
		cacheUpArtIDs map[int64][]int64
		missUpArts    map[int64][][2]int64
		// missUpArtIDs    map[int64][][2]int64
	)
	res = make(map[int64][]int64)
	if exists, err = s.dao.ExpireUppersCache(c, mids); err != nil {
		addCache = false
		err = nil
	}
	for _, mid := range mids {
		if !exists[mid] {
			missMids = append(missMids, mid)
		} else {
			cacheMids = append(cacheMids, mid)
		}
	}
	// from cache
	group.Go(func() (err error) {
		if cacheUpArtIDs, err = s.dao.UppersCaches(c, cacheMids, start, end); err != nil {
			dao.PromError("upper:获取up主文章列表")
		}
		return
	})
	group.Go(func() (err error) {
		if len(missMids) > 0 {
			missUpArts, err = s.dao.UppersPassed(c, missMids)
		}
		return
	})
	group.Wait()
	for mid, ids := range cacheUpArtIDs {
		res[mid] = ids
	}
	for mid, arts := range missUpArts {
		var ids []int64
		for _, art := range arts {
			ids = append(ids, art[0])
		}
		if (start == 0) && (end == -1) {
			res[mid] = ids
			continue
		}
		if len(ids) <= start {
			res[mid] = []int64{}
			continue
		}
		if len(ids) < end {
			res[mid] = ids[start:]
		} else {
			res[mid] = ids[start:end]
		}
	}
	if addCache && (len(missUpArts) > 0) {
		s.dao.AddUpperCaches(c, missUpArts)
	}
	return
}

// UpperArtsCount count upper article
func (s *Service) UpperArtsCount(c context.Context, mid int64) (res int, err error) {
	var (
		exists map[int64]bool
		arts   map[int64][][2]int64
	)
	if exists, err = s.dao.ExpireUppersCache(c, []int64{mid}); err != nil {
		err = nil
		return
	}
	if exists[mid] {
		if res, err = s.dao.UpperArtsCountCache(c, mid); err == nil {
			return
		}
		log.Error("s.dao.UpperArtsCountCache(%v) err: %+v", mid, err)
	}
	if arts, err = s.dao.UppersPassed(c, []int64{mid}); err != nil {
		dao.PromError("upper:获取作者文章列表")
		return
	}
	res = len(arts[mid])
	cache.Save(func() {
		s.dao.AddUpperCaches(context.TODO(), arts)
	})
	return
}
