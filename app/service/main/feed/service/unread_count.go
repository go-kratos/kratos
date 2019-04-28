package service

import (
	"context"
	"go-common/app/service/main/feed/dao"
	"go-common/app/service/main/feed/model"
	"sync/atomic"
	"time"

	"go-common/app/service/main/archive/model/archive"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	"go-common/library/sync/errgroup"
)

// UnreadCount get count of unread archives
func (s *Service) UnreadCount(c context.Context, app bool, withoutBangumi bool, mid int64, ip string) (count int, err error) {
	var (
		t             int64
		last          time.Time
		pullInterval  xtime.Duration
		marcs         map[int64][]*archive.AidPubTime
		minTotalCount int
		bangumiFeeds  []*feedmdl.Feed
		unreadCount   int64
		ft            = model.FeedType(app)
	)
	if t, err = s.dao.LastAccessCache(c, ft, mid); err != nil {
		dao.PromError("未读数获取上次访问时间缓存", "s.dao.LastAccessCache(app:%v, mid:%v, err: %v", app, mid, err)
		return 0, nil
	}
	if t == 0 {
		return
	}
	last = time.Unix(t, 0)
	pullInterval = s.pullInterval(ft)
	if time.Now().Sub(last) < time.Duration(pullInterval) {
		count, _ = s.dao.UnreadCountCache(c, ft, mid)
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if app {
			minTotalCount = s.c.Feed.AppLength
		} else {
			minTotalCount = s.c.Feed.WebLength
		}
		if marcs, err = s.attenUpArcs(errCtx, minTotalCount, mid, ip); err != nil {
			dao.PromError("未读数attenUpArcs", "s.attenUpArcs(count:%v, mid:%v, err: %v", minTotalCount, mid, err)
			err = nil
			return nil
		}
		for _, arcs := range marcs {
			for _, arc := range arcs {
				if int64(arc.PubDate) > t {
					atomic.AddInt64(&unreadCount, 1)
				}
			}
		}
		return nil
	})
	group.Go(func() error {
		if withoutBangumi {
			return nil
		}
		bangumiFeeds = s.genBangumiFeed(errCtx, mid, ip)
		for _, f := range bangumiFeeds {
			if int64(f.PubDate) > t {
				atomic.AddInt64(&unreadCount, 1)
			}
		}
		return nil
	})
	group.Wait()
	count = int(unreadCount)
	if count > s.c.Feed.MaxTotalCnt {
		count = s.c.Feed.MaxTotalCnt
	}
	s.addCache(func() {
		s.dao.AddUnreadCountCache(context.Background(), ft, mid, count)
	})
	return
}

// ArticleUnreadCount get count of unread articles
func (s *Service) ArticleUnreadCount(c context.Context, mid int64, ip string) (count int, err error) {
	var (
		t            int64
		last         time.Time
		pullInterval xtime.Duration
		ft           = model.TypeArt
	)
	if t, err = s.dao.LastAccessCache(c, ft, mid); err != nil {
		log.Error("s.dao.LastAccessCache(app:%v, mid:%v), err: %v", ft, mid, err)
		return 0, nil
	}
	if t == 0 {
		// no access, no unread
		return
	}
	last = time.Unix(t, 0)
	pullInterval = s.pullInterval(ft)
	if time.Now().Sub(last) < time.Duration(pullInterval) {
		count, _ = s.dao.UnreadCountCache(c, ft, mid)
		return
	}
	res := s.genArticleFeed(c, mid, s.c.Feed.ArticleFeedLength, ip)
	for _, f := range res {
		if int64(f.PublishTime) > t {
			count++
		}
	}
	if count > s.c.Feed.MaxTotalCnt {
		count = s.c.Feed.MaxTotalCnt
	}
	s.addCache(func() {
		s.dao.AddUnreadCountCache(context.Background(), ft, mid, count)
	})
	return
}
