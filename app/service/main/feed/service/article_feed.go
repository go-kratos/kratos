package service

import (
	"context"
	"sort"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/feed/dao"
	"go-common/app/service/main/feed/model"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/log"
)

// ArticleFeed get feed of ups.
func (s *Service) ArticleFeed(c context.Context, mid int64, pn, ps int, ip string) (res []*artmdl.Meta, err error) {
	var (
		fp     = pn == 1
		cached bool
		ft     = model.TypeArt
		from   = "1"
	)
	if fp {
		// if first page and in 5 minites, return cached feed.
		fp = s.checkLast(c, ft, mid)
	}
	// refresh expire feed cache if user access feed.
	if cached, err = s.dao.ExpireFeedCache(c, ft, mid); err != nil {
		dao.PromError("expire feed cache", "s.dao.ExpireFeedCache(%d) error(%v)", mid, err)
		err = nil
	}
	defer func() {
		for _, meta := range res {
			if meta != nil && meta.Author != nil && meta.Author.Name == "" {
				dao.PromError("bug:noauthor"+from, "bugfix: %+v, author: %+v from: %v", meta, meta.Author, from)
			}
		}
	}()
	if cached && !fp {
		// if cache err, will return directly.
		if res, err = s.articleFeedCache(c, mid, pn, ps, ip); err == nil {
			return
		}
		dao.PromError("获取文章feed cache", "s.articleFeedCache(%v, %v, %v, %v, len: %v) error(%v)", ft, mid, pn, ps, len(res), err)
	}
	res, err = s.articleFeed(c, mid, pn, ps, ip)
	from = "2"
	if fp {
		s.addCache(func() {
			s.dao.AddUnreadCountCache(context.Background(), ft, mid, 0)
		})
	}
	return
}

func (s *Service) articleFeed(c context.Context, mid int64, pn, ps int, ip string) (res []*artmdl.Meta, err error) {
	dao.MissedCount.Incr("Article-feed")
	var (
		start  = (pn - 1) * ps
		end    = start + ps // from slice, end no -1
		tmpRes []*artmdl.Meta
	)
	tmpRes = s.genArticleFeed(c, mid, s.c.Feed.ArticleFeedLength, ip)
	if len(tmpRes) == 0 || len(tmpRes) < start {
		// 当用户取关所有up主时清除缓存
		s.addCache(func() {
			s.dao.AddArticleFeedCache(context.Background(), mid, tmpRes)
		})
		return
	}
	sort.Sort(feedmdl.ArticleFeeds(tmpRes))
	if end < len(tmpRes) {
		res = tmpRes[start:end]
	} else {
		res = tmpRes[start:]
	}
	s.addCache(func() {
		s.dao.AddArticleFeedCache(context.Background(), mid, tmpRes)
	})
	return
}

// articleFeedCache get Article feed by cache.
func (s *Service) articleFeedCache(c context.Context, mid int64, pn, ps int, ip string) (res []*artmdl.Meta, err error) {
	dao.CachedCount.Incr("Article-feed")
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1 // from cache, end-1
		aids  []int64
		am    map[int64]*artmdl.Meta
	)
	if aids, err = s.dao.ArticleFeedCache(c, mid, start, end); err != nil || len(aids) == 0 {
		return
	}
	if am, err = s.articles(c, ip, aids...); err != nil {
		return
	}
	for _, aid := range aids {
		if _, ok := am[aid]; ok {
			res = append(res, am[aid])
		}
	}
	return
}

func (s *Service) genArticleFeed(c context.Context, mid int64, minTotalCount int, ip string) (res []*artmdl.Meta) {
	var (
		marts map[int64][]*artmdl.Meta
		err   error
	)
	if marts, err = s.attenUpArticles(c, minTotalCount, mid, ip); err != nil {
		log.Error("s.attenUpArticles(mid: %v) err: %v", mid, err)
		return
	}
	for _, as := range marts {
		for _, a := range as {
			res = append(res, a)
		}
	}
	return
}
