package service

import (
	"context"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"

	"go-common/library/sync/errgroup"
)

// SetStat sets article's stat.
func (s *Service) SetStat(c context.Context, aid int64, stat *artmdl.Stats) (err error) {
	group, ctx := errgroup.WithContext(c)
	group.Go(func() error {
		return s.dao.AddArticleStatsCache(ctx, aid, stat)
	})
	group.Go(func() error {
		s.SendMessage(ctx, aid, stat)
		return nil
	})
	err = group.Wait()
	return
}

func (s *Service) stat(c context.Context, aid int64) (res *artmdl.Stats, err error) {
	var addCache = true
	if res, err = s.dao.ArticleStatsCache(c, aid); err != nil {
		err = nil
		addCache = false
	} else if res != nil {
		return
	}
	res, err = s.dao.ArticleStats(c, aid)
	if res != nil && addCache {
		cache.Save(func() {
			s.dao.AddArticleStatsCache(context.TODO(), aid, res)
		})
	}
	return
}

func (s *Service) stats(c context.Context, aids []int64) (res map[int64]*artmdl.Stats, err error) {
	var (
		cachedArtStats  map[int64]*artmdl.Stats
		missed          []int64
		missedArtsStats map[int64]*artmdl.Stats
		addCache        = true
	)
	if cachedArtStats, missed, err = s.dao.ArticlesStatsCache(c, aids); err != nil {
		addCache = false
		missed = aids
		err = nil
	}
	if len(missed) > 0 {
		if missedArtsStats, err = s.dao.ArticlesStats(c, missed); err != nil {
			addCache = false
			dao.PromError("stat:文章计数")
			err = nil
		}
	}
	res = make(map[int64]*artmdl.Stats)
	for aid, art := range cachedArtStats {
		res[aid] = art
	}
	if missedArtsStats == nil {
		missedArtsStats = make(map[int64]*artmdl.Stats)
	} else {
		for aid, art := range missedArtsStats {
			res[aid] = art
		}
	}
	if addCache {
		for _, aid := range aids {
			if _, ok := res[aid]; !ok {
				missedArtsStats[aid] = new(artmdl.Stats)
			}
		}
		cache.Save(func() {
			for id, stats := range missedArtsStats {
				s.dao.AddArticleStatsCache(context.TODO(), id, stats)
			}
		})
	}
	return
}
