package service

import (
	"context"
	"sort"

	"go-common/app/service/main/feed/dao"
	feedmdl "go-common/app/service/main/feed/model"
)

// ArchiveFeed get feed of ups.
func (s *Service) ArchiveFeed(c context.Context, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	if res, err = s.archiveFeedCache(c, mid, pn, ps, ip); (err == nil) && (len(res) > 0) {
		return
	}
	res, err = s.archiveFeed(c, mid, pn, ps, ip)
	return
}

func (s *Service) archiveFeed(c context.Context, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	dao.MissedCount.Incr("archive-feed")
	var (
		start  = (pn - 1) * ps
		end    = start + ps // from slice, end no -1
		tmpRes []*feedmdl.Feed
	)
	tmpRes = s.genArchiveFeed(c, true, mid, s.c.Feed.ArchiveFeedLength, ip)
	sort.Sort(feedmdl.Feeds(tmpRes))
	res = s.sliceFeeds(tmpRes, start, end)
	res, err = s.fillArchiveFeeds(c, res, ip)
	s.addCache(func() {
		s.dao.AddArchiveFeedCache(context.Background(), mid, tmpRes)
	})
	return
}

// archiveFeedCache get archive feed by cache.
func (s *Service) archiveFeedCache(c context.Context, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	dao.CachedCount.Incr("archive-feed")
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1 // from cache, end-1
	)
	if res, err = s.dao.ArchiveFeedCache(c, mid, start, end); err != nil {
		return
	}
	res, err = s.fillArchiveFeeds(c, res, ip)
	return
}
