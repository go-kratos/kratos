package service

import (
	"context"
	"sort"

	"go-common/app/service/main/feed/dao"
	feedmdl "go-common/app/service/main/feed/model"
)

// BangumiFeed get feed of bangumi.
func (s *Service) BangumiFeed(c context.Context, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	if res, err = s.bangumiFeedCache(c, mid, pn, ps, ip); (err == nil) && (len(res) > 0) {
		return
	}
	res, err = s.bangumiFeed(c, mid, pn, ps, ip)
	return
}

func (s *Service) bangumiFeed(c context.Context, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	dao.MissedCount.Incr("bangumi-feed")
	var (
		start  = (pn - 1) * ps
		end    = start + ps // from slice, end no -1
		tmpRes []*feedmdl.Feed
	)
	tmpRes = s.genBangumiFeed(c, mid, ip)
	sort.Sort(feedmdl.Feeds(tmpRes))
	res = s.sliceFeeds(tmpRes, start, end)
	s.addCache(func() {
		s.dao.AddBangumiFeedCache(c, mid, tmpRes)
	})
	return
}

// bangumiFeedCache get bangumi feed by cache.
func (s *Service) bangumiFeedCache(c context.Context, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	dao.CachedCount.Incr("bangumi-feed")
	var (
		start  = (pn - 1) * ps
		end    = start + ps - 1 // from cache, end-1
		endPos = end
		bids   []int64
	)
	if bids, err = s.dao.BangumiFeedCache(c, mid, start, endPos); err != nil {
		return
	}
	if res, err = s.bangumiFeedFromSeason(c, bids, ip); err != nil {
		dao.PromError("获取番剧feed", "s.bangumiFeed(bids: %v) err: %v", bids, err)
		return
	}
	return
}
