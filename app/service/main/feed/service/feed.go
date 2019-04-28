package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/feed/dao"
	feedmdl "go-common/app/service/main/feed/model"
	xtime "go-common/library/time"

	"go-common/library/sync/errgroup"
)

func archiveAppName(ft int) string {
	if ft == feedmdl.TypeApp {
		return "app"
	} else if ft == feedmdl.TypeWeb {
		return "web"
	}
	return "article"
}

func (s *Service) pullInterval(ft int) xtime.Duration {
	if ft == feedmdl.TypeApp {
		return s.c.Feed.AppPullInterval
	} else if ft == feedmdl.TypeWeb {
		return s.c.Feed.WebPullInterval
	} else {
		return s.c.Feed.ArtPullInterval
	}
}

// Feed get feed of ups and bangumi.
func (s *Service) Feed(c context.Context, app bool, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	var (
		fp     = pn == 1
		cached bool
		ft     = feedmdl.FeedType(app)
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
	// check feed if exist same aid
	defer func() {
		var exists = map[int64]bool{}
		for _, feed := range res {
			if _, ok := exists[feed.ID]; ok {
				dao.PromError("重复动态", "feed same user: %d, id:%d", mid, feed.ID)
			} else {
				exists[feed.ID] = true
			}
		}
	}()
	if cached && !fp {
		// when gen feed err will return
		if res, err = s.feedCache(c, ft, mid, pn, ps, ip); err != nil {
			dao.PromError("获取Feed cache", "s.feedCache(%v, %v, %v, %v, len: %v) error(%v)", app, mid, pn, ps, len(res), err)
		}
		return
	}
	res, err = s.feed(c, ft, mid, pn, ps, ip)
	if fp {
		s.addCache(func() {
			s.dao.AddUnreadCountCache(context.Background(), ft, mid, 0)
		})
	}
	return
}

// checkLast check if last access is in 5 minites.
func (s *Service) checkLast(c context.Context, ft int, mid int64) (fp bool) {
	var (
		t            int64
		err          error
		now          = time.Now()
		pullInterval xtime.Duration
	)
	fp = true
	if t, err = s.dao.LastAccessCache(c, ft, mid); err != nil {
		return
	}
	last := time.Unix(t, 0)
	pullInterval = s.pullInterval(ft)
	if now.Sub(last) < time.Duration(pullInterval) {
		fp = false
	} else {
		s.addCache(func() {
			s.dao.AddLastAccessCache(context.TODO(), ft, mid, now.Unix())
		})
	}
	return
}

func (s *Service) genArchiveFeed(c context.Context, fold bool, mid int64, minTotalCount int, ip string) (res []*feedmdl.Feed) {
	var (
		marcs map[int64][]*arcmdl.AidPubTime
		err   error
	)
	if marcs, err = s.attenUpArcs(c, minTotalCount, mid, ip); err != nil {
		dao.PromError("获取关注up主稿件", "s.attenUpArcs(mid: %v) err: %v", mid, err)
		return
	}
	if fold {
		for _, as := range marcs {
			switch len(as) {
			case 0: // no archives from upper
			case 1: // just have one archive of upper
				res = append(res, arcTimeToFeed(as[0]))
			default: // fold
				for _, appFeed := range s.fold(as) {
					res = append(res, appFeed)
				}
			}
		}
	} else {
		for _, as := range marcs {
			for _, a := range as {
				res = append(res, arcTimeToFeed(a))
			}
		}
	}
	return
}

func (s *Service) genBangumiFeed(c context.Context, mid int64, ip string) (res []*feedmdl.Feed) {
	var (
		seasonIDs []int64
		err       error
	)
	if seasonIDs, err = s.bangumi.BangumiPull(c, mid, ip); err != nil {
		return nil
	}
	res, _ = s.bangumiFeedFromSeason(c, seasonIDs, ip)
	return
}

// feed get feed.
func (s *Service) feed(c context.Context, ft int, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	dao.MissedCount.Incr(archiveAppName(ft) + "-feed")
	var (
		start         = (pn - 1) * ps
		end           = start + ps // from slice, end no -1
		tmpRes        []*feedmdl.Feed
		bangumiFeeds  []*feedmdl.Feed
		group         *errgroup.Group
		errCtx        context.Context
		minTotalCount int
		fold          bool
	)
	group, errCtx = errgroup.WithContext(c)

	if ft == feedmdl.TypeApp {
		minTotalCount = s.c.Feed.AppLength
		fold = true
	} else {
		minTotalCount = s.c.Feed.WebLength
		fold = false
	}
	// fetch archives
	group.Go(func() error {
		tmpRes = s.genArchiveFeed(errCtx, fold, mid, minTotalCount, ip)
		return nil
	})
	// fetch bangumis
	group.Go(func() error {
		bangumiFeeds = s.genBangumiFeed(errCtx, mid, ip)
		return nil
	})
	group.Wait()
	// merge archives and bangumis
	tmpRes = append(tmpRes, bangumiFeeds...)
	if len(tmpRes) == 0 || len(tmpRes) < start {
		// 当用户取关所有up主时清除缓存
		s.addCache(func() {
			s.dao.AddFeedCache(context.Background(), ft, mid, tmpRes)
		})
		return
	}
	// sort by aid desc ,then set cache.
	sort.Sort(feedmdl.Feeds(tmpRes))
	res = s.sliceFeeds(tmpRes, start, end)
	res, err = s.fillArchiveFeeds(c, res, ip)
	s.addCache(func() {
		s.dao.AddFeedCache(context.Background(), ft, mid, tmpRes)
	})
	return
}

func (s *Service) sliceFeeds(fs []*feedmdl.Feed, start, end int) (res []*feedmdl.Feed) {
	if start > len(fs) {
		return
	}
	if end < len(fs) {
		res = fs[start:end]
	} else {
		res = fs[start:]
	}
	return
}

// feedCache get feed by cache.
func (s *Service) feedCache(c context.Context, ft int, mid int64, pn, ps int, ip string) (res []*feedmdl.Feed, err error) {
	dao.CachedCount.Incr(archiveAppName(ft) + "-feed")
	var (
		start              = (pn - 1) * ps
		end                = start + ps - 1 // from cache, end-1
		endPos             = end
		bids               []int64
		bangumiFeeds       []*feedmdl.Feed
		group              *errgroup.Group
		errCtx             context.Context
		arcErr, bangumiErr error
	)
	if res, bids, err = s.dao.FeedCache(c, ft, mid, start, endPos); err != nil {
		err = nil
		return
	}
	group, errCtx = errgroup.WithContext(c)
	group.Go(func() error {
		res, arcErr = s.fillArchiveFeeds(errCtx, res, ip)
		return nil
	})
	group.Go(func() error {
		bangumiFeeds, bangumiErr = s.bangumiFeedFromSeason(errCtx, bids, ip)
		return nil
	})
	group.Wait()
	if (arcErr != nil) && (bangumiErr != nil) {
		dao.PromError("生成feed流", "s.feedCache(mid: %v) arc:%v bangumi: %v", mid, arcErr, bangumiErr)
		err = fmt.Errorf("s.feedCache(mid: %v) arc:%v bangumi: %v", mid, arcErr, bangumiErr)
		return
	}
	s.replaceFeeds(res, bangumiFeeds)
	return
}

func (s *Service) fillArchiveFeeds(c context.Context, fs []*feedmdl.Feed, ip string) (res []*feedmdl.Feed, err error) {
	var (
		allAids []int64
		am      map[int64]*api.Arc
	)
	if len(fs) == 0 {
		return
	}
	for _, fe := range fs {
		if fe.Type != feedmdl.ArchiveType {
			continue
		}
		allAids = append(allAids, fe.ID)
		for _, a := range fe.Fold {
			allAids = append(allAids, a.Aid)
		}
	}
	if am, err = s.archives(c, allAids, ip); err != nil {
		return
	}
	for _, fe := range fs {
		if fe.Type == feedmdl.ArchiveType {
			fe = fmtArc(fe, am)
		}
		if fe != nil {
			res = append(res, fe)
		}
	}
	return
}

func fmtArc(feed *feedmdl.Feed, archives map[int64]*api.Arc) *feedmdl.Feed {
	if feed == nil {
		return nil
	}
	var arcs []*api.Arc
	if v, ok := archives[feed.ID]; ok {
		arcs = append(arcs, v)
	}
	for _, arc := range feed.Fold {
		if v, ok := archives[arc.Aid]; ok {
			arcs = append(arcs, v)
		}
	}
	if len(arcs) == 0 {
		return nil
	}
	feed.Archive = arcs[0]
	feed.PubDate = arcs[0].PubDate
	feed.ID = arcs[0].Aid
	feed.Fold = arcs[1:]
	return feed
}

func (s *Service) replaceFeeds(res []*feedmdl.Feed, fs []*feedmdl.Feed) {
	var (
		f      *feedmdl.Feed
		key    string
		format = "%v-%v"
		m      = make(map[string]*feedmdl.Feed)
	)
	for _, f = range fs {
		key = fmt.Sprintf(format, f.Type, f.ID)
		m[key] = f
	}
	for i, f := range res {
		key = fmt.Sprintf(format, f.Type, f.ID)
		if _, ok := m[key]; ok {
			res[i] = m[key]
		}
	}
}

// fold archives in every 4 hours
func (s *Service) fold(as []*arcmdl.AidPubTime) (res []*feedmdl.Feed) {
	if len(as) == 0 {
		return
	}
	sort.Sort(feedmdl.Arcs(as))
	var (
		fa = arcTimeToFeed(as[0]) // the cover archive of item
		ft = as[0].PubDate.Time() // first archive pubdate
		// at every 4 o'clock we will fold archives with a cover archive and the count of folded archives
		// ps. hour at [0,4),ch is 0;hour at [4,8),ch is 4;hour at [8,12),ch is 8;hour at [12,16),ch is 12;hour at [16,20),ch is 16;hour at [20,24),ch is 20;
		ch = (ft.Hour() / 4) * 4 // check hour
		at time.Time             // archive pubdate
	)
	for k, a := range as[1:] {
		isEnd := k == (len(as[1:]) - 1)
		y1, m1, d1 := ft.Date()
		at = a.PubDate.Time()
		y2, m2, d2 := at.Date()
		// NOTE: original video(copyright == 1) does not fold
		if a.Copyright != 1 && (y1 == y2 && m1 == m2 && d1 == d2 && at.Hour() >= ch) {
			fa.Fold = append(fa.Fold, &api.Arc{Aid: a.Aid})
			if isEnd {
				res = append(res, fa)
			}
		} else {
			res = append(res, fa)
			fa = arcTimeToFeed(a)
			if isEnd {
				res = append(res, fa)
			} else {
				// next item
				ft = a.PubDate.Time()
				ch = (ft.Hour() / 4) * 4
			}
		}
	}
	return
}

//Fold get fold archives.
func (s *Service) Fold(c context.Context, mid int64, aid int64, ip string) (res []*feedmdl.Feed, err error) {
	var (
		arcsm map[int64][]*arcmdl.AidPubTime
		arcs  []*arcmdl.AidPubTime
	)
	if arcsm, err = s.upArcs(c, s.c.Feed.AppLength, ip, mid); err != nil {
		return
	}
	sort.Sort(feedmdl.Arcs(arcsm[mid]))
	for i, arc := range arcsm[mid] {
		if arc.Aid == aid {
			arcs = arcsm[mid][i:]
			break
		}
	}
	switch len(arcs) {
	case 0, 1:
	default:
		for _, appFeed := range s.fold(arcs) {
			if appFeed.ID == aid {
				for _, arc := range appFeed.Fold {
					res = append(res, arcToFeed(arc))
				}
			}
		}
	}
	res, err = s.fillArchiveFeeds(c, res, ip)
	return
}

// PurgeFeedCache purge cache when attention/unattention upper
func (s *Service) PurgeFeedCache(c context.Context, mid int64, ip string) (err error) {
	if err = s.dao.PurgeFeedCache(c, feedmdl.TypeApp, mid); err != nil {
		return
	}
	if err = s.dao.PurgeFeedCache(c, feedmdl.TypeWeb, mid); err != nil {
		return
	}
	err = s.dao.PurgeFeedCache(c, feedmdl.TypeArt, mid)
	return
}

func (s *Service) bangumiFeedFromSeason(c context.Context, seasonIDs []int64, ip string) (feeds []*feedmdl.Feed, err error) {
	var (
		bm       map[int64]*feedmdl.Bangumi
		cached   map[int64]*feedmdl.Bangumi
		missed   []int64
		addCache = true
	)
	if cached, missed, err = s.dao.BangumisCache(c, seasonIDs); err != nil {
		dao.PromError("番剧feed中调用缓存", "s.dao.BangumisCache err: %v", err)
		err = nil
		addCache = false
		missed = seasonIDs
	}
	if len(missed) > 0 {
		if bm, err = s.bangumi.BangumiSeasons(c, missed, ip); err != nil {
			return
		}
		if addCache {
			s.addCache(func() {
				s.dao.AddBangumisCacheMap(context.Background(), bm)
			})
		}
		for bid, b := range bm {
			cached[bid] = b
		}
	}
	for _, bid := range seasonIDs {
		if bangumi, ok := cached[bid]; ok {
			feeds = append(feeds, bangumiToFeed(bangumi))
		}
	}
	return
}

func arcToFeed(arc *api.Arc) *feedmdl.Feed {
	return &feedmdl.Feed{ID: arc.Aid, Type: feedmdl.ArchiveType, PubDate: arc.PubDate}
}
func arcTimeToFeed(arc *arcmdl.AidPubTime) *feedmdl.Feed {
	return &feedmdl.Feed{ID: arc.Aid, Type: feedmdl.ArchiveType, PubDate: arc.PubDate}
}

func bangumiToFeed(b *feedmdl.Bangumi) *feedmdl.Feed {
	return &feedmdl.Feed{ID: b.SeasonID, Bangumi: b, Type: feedmdl.BangumiType, PubDate: xtime.Time(b.Ts)}
}
