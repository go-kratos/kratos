package feed

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	"go-common/app/interface/main/app-feed/model/live"
	"go-common/app/interface/main/app-feed/model/tag"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	busfeed "go-common/app/service/main/feed/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_iosBanBangumi      = 4310
	_androidBanBangumi  = 502000
	_androidIBanBangumi = 104000
)

func (s *Service) Upper(c context.Context, mid int64, plat int8, build int, pn, ps int, now time.Time) (is []*feed.Item, lp bool) {
	if (plat == model.PlatIPhone && build > _iosBanBangumi) || (plat == model.PlatAndroid && build > _androidBanBangumi) || (plat == model.PlatAndroidI && build >= _androidIBanBangumi) || plat == model.PlatIPhoneB {
		is, lp = s.UpperFeed(c, mid, plat, build, pn, ps, now)
	} else {
		is, lp = s.UpperArchive(c, mid, plat, build, pn, ps, now)
	}
	return
}

// UpperFeed get the archives and bangumi for feed
// if archives are less then `_minTotalCnt` then will fill with recommended archives
func (s *Service) UpperFeed(c context.Context, mid int64, plat int8, build int, pn, ps int, now time.Time) (is []*feed.Item, lp bool) {
	var (
		err      error
		hc       bool
		uis, fis []*busfeed.Feed
		fp       = pn == 1
	)
	if fp {
		var (
			unread int
			count  = ps
		)
		if unread, err = s.upper.UnreadCountCache(c, mid); err != nil {
			log.Error("%+v", err)
		} else if unread != 0 {
			count = s.c.Feed.FeedCacheCount
		}
		if uis, err = s.upper.Feed(c, mid, pn, count); err != nil {
			log.Error("%+v", err)
			// get cache from redis
			if hc, err = s.upper.ExpireUpItem(c, mid); err != nil {
				log.Error("%+v", err)
			} else if hc {
				if uis, err = s.upperCache(c, mid, plat, build, pn, ps, now); err != nil {
					log.Error("%+v", err)
				}
			}
			if len(uis) == 0 {
				is = _emptyItem
				lp = true
				return
			}
		} else if unread != 0 {
			s.addCache(func() {
				s.upper.AddUpItemCaches(context.Background(), mid, uis...)
			})
		}
	} else {
		if uis, err = s.upper.Feed(c, mid, pn, ps); err != nil {
			log.Error("%+v", err)
			// get cache from redis
			if hc, err = s.upper.ExpireUpItem(c, mid); err != nil {
				log.Error("%+v", err)
			} else if hc {
				if uis, err = s.upperCache(c, mid, plat, build, pn, ps, now); err != nil {
					log.Error("%+v", err)
				}
			}
			if len(uis) == 0 {
				is = _emptyItem
				lp = true
				return
			}
		}
	}
	// handle feed
	if len(uis) > ps {
		fis = uis[:ps]
	} else {
		fis = uis
	}
	is = s.upperItem(c, fis, mid, now)
	if len(is) < ps {
		lp = true
	}
	if len(is) == 0 {
		is = _emptyItem
	}
	return
}

func (s *Service) upperCache(c context.Context, mid int64, plat int8, build, pn, ps int, now time.Time) (uis []*busfeed.Feed, err error) {
	var (
		start     = (pn - 1) * ps
		end       = start + ps // from slice, end no -1
		aids      []int64
		am        map[int64]*api.Arc
		seasonIDs []int64
		psm       map[int64]*busfeed.Bangumi
	)
	if uis, aids, seasonIDs, err = s.upper.UpItemCaches(c, mid, start, end); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(aids) > 0 {
		if am, err = s.arc.Archives(c, aids); err != nil {
			log.Error("%+v", err)
			return
		}
	}
	if len(seasonIDs) > 0 && ((plat == model.PlatIPhone && build > _iosBanBangumi) || (plat == model.PlatAndroid && build > _androidBanBangumi) || plat == model.PlatIPhoneB) {
		if psm, err = s.bgm.PullSeasons(c, seasonIDs, now); err != nil {
			log.Error("%+v", err)
		}
	}
	for _, ui := range uis {
		switch ui.Type {
		case busfeed.ArchiveType:
			if a, ok := am[ui.ID]; ok {
				ui.Archive = a
			}
			for _, r := range ui.Fold {
				if a, ok := am[r.Aid]; ok {
					r = a
				} else {
					r = nil
				}
			}
		case busfeed.BangumiType:
			if s, ok := psm[ui.ID]; ok {
				ui.Bangumi = s
			}
		}
	}
	return
}

func (s *Service) UpperArchive(c context.Context, mid int64, plat int8, build int, pn, ps int, now time.Time) (is []*feed.Item, lp bool) {
	var (
		err error
		uis []*busfeed.Feed
		hc  bool
	)
	if uis, err = s.upper.ArchiveFeed(c, mid, pn, ps); err != nil {
		log.Error("%+v", err)
		// get cache from redis
		if hc, err = s.upper.ExpireUpItem(c, mid); err != nil {
			log.Error("%+v", err)
		} else if hc {
			if uis, err = s.upperCache(c, mid, plat, build, pn, ps, now); err != nil {
				log.Error("%+v", err)
			}
		}
		if len(uis) == 0 {
			is = _emptyItem
			lp = true
			return
		}
	}
	// handle feed
	is = s.upperItem(c, uis, mid, now)
	if len(is) < ps {
		lp = true
	}
	if len(is) == 0 {
		is = _emptyItem
	}
	return
}

func (s *Service) UpperBangumi(c context.Context, mid int64, plat int8, build int, pn, ps int, now time.Time) (is []*feed.Item, lp bool) {
	var (
		err error
		uis []*busfeed.Feed
		hc  bool
	)
	if uis, err = s.upper.BangumiFeed(c, mid, pn, ps); err != nil {
		log.Error("%+v", err)
		// get cache from redis
		if hc, err = s.upper.ExpireUpItem(c, mid); err != nil {
			log.Error("%+v", err)
		} else if hc {
			if uis, err = s.upperCache(c, mid, plat, build, pn, ps, now); err != nil {
				log.Error("%+v", err)
			}
		}
		if len(uis) == 0 {
			is = _emptyItem
			lp = true
			return
		}
	}
	// handle feed
	is = s.upperItem(c, uis, mid, now)
	if len(is) < ps {
		lp = true
	}
	if len(is) == 0 {
		is = _emptyItem
	}
	return
}

func (s *Service) UpperRecent(c context.Context, mid, upperID, aid int64, now time.Time) (is []*feed.Item) {
	var (
		err error
		uis []*busfeed.Feed
	)
	if uis, err = s.upper.Recent(c, upperID, aid); err != nil {
		log.Error("%+v", err)
	}
	// handle feed
	is = s.upperItem(c, uis, mid, now)
	if len(is) == 0 {
		is = _emptyItem
	}
	return
}

func (s *Service) upperItem(c context.Context, uis []*busfeed.Feed, mid int64, now time.Time) (is []*feed.Item) {
	var (
		g            *errgroup.Group
		ctx          context.Context
		owners, aids []int64
		follows      map[int64]bool
		tm           map[string][]*tag.Tag
		err          error
	)
	owners = make([]int64, 0, len(uis))
	for _, ui := range uis {
		if ui != nil {
			if ui.Archive != nil {
				owners = append(owners, ui.Archive.Author.Mid)
				aids = append(aids, ui.Archive.Aid)
			}
			for _, r := range ui.Fold {
				if r != nil {
					aids = append(aids, r.Aid)
				}
			}
		}
	}
	g, ctx = errgroup.WithContext(c)
	if len(owners) != 0 {
		g.Go(func() (err error) {
			follows = s.acc.Relations3(ctx, owners, mid)
			return
		})
	}
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if tm, err = s.tg.Tags(ctx, mid, aids, now); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	is = make([]*feed.Item, 0, len(uis))
	for _, ui := range uis {
		if ui != nil {
			switch ui.Type {
			case busfeed.ArchiveType:
				if ui.Archive != nil && ui.Archive.IsNormal() {
					i := &feed.Item{}
					i.FromAv(&archive.ArchiveWithPlayer{Archive3: archive.BuildArchive3(ui.Archive)})
					i.RecCnt = len(ui.Fold)
					if len(ui.Fold) > 0 {
						ris := make([]*feed.Item, 0, len(ui.Fold))
						for _, r := range ui.Fold {
							if r != nil && r.IsNormal() {
								ri := &feed.Item{}
								ri.FromAv(&archive.ArchiveWithPlayer{Archive3: archive.BuildArchive3(r)})
								if infos, ok := tm[strconv.FormatInt(r.Aid, 10)]; ok {
									if len(infos) != 0 {
										ri.Tag = &feed.Tag{TagID: infos[0].ID, TagName: infos[0].Name, IsAtten: infos[0].IsAtten, Count: &feed.TagCount{Atten: infos[0].Count.Atten}}
									}
								}
								if follows[i.Mid] {
									ri.IsAtten = 1
								}
								ris = append(ris, ri)
							}
						}
						i.Recent = ris
					}
					if infos, ok := tm[strconv.FormatInt(ui.Archive.Aid, 10)]; ok {
						if len(infos) != 0 {
							i.Tag = &feed.Tag{TagID: infos[0].ID, TagName: infos[0].Name, IsAtten: infos[0].IsAtten, Count: &feed.TagCount{Atten: infos[0].Count.Atten}}
						}
					}
					if follows[i.Mid] {
						i.IsAtten = 1
					}
					is = append(is, i)
				}
			case busfeed.BangumiType:
				if ui.Bangumi != nil {
					i := &feed.Item{}
					i.FromUpBangumi(ui.Bangumi)
					is = append(is, i)
				}
			}
		}
	}
	return
}

func (s *Service) UpperLive(c context.Context, mid int64) (is []*feed.Item, count int) {
	var (
		err error
		fs  []*live.Feed
		pn  = 1
		ps  = s.c.Feed.LiveFeedCount
	)
	if fs, count, err = s.lv.FeedList(c, mid, pn, ps); err != nil {
		log.Error("%+v", err)
	}
	for _, f := range fs {
		i := &feed.Item{}
		i.FromUpLive(f)
		is = append(is, i)
	}
	return
}

func (s *Service) UpperArticle(c context.Context, mid int64, plat int8, build int, pn, ps int, now time.Time) (is []*feed.Item, lp bool) {
	var (
		err error
		uis []*article.Meta
	)
	if uis, err = s.upper.ArticleFeed(c, mid, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	// handle feed
	is = s.articleItem(c, uis, mid)
	if len(is) < ps {
		lp = true
	}
	if len(is) == 0 {
		is = _emptyItem
	}
	return
}

func (s *Service) articleItem(c context.Context, uis []*article.Meta, mid int64) (is []*feed.Item) {
	is = make([]*feed.Item, 0, len(uis))
	for _, ui := range uis {
		if ui != nil {
			i := &feed.Item{}
			i.FromUpArticle(ui)
			is = append(is, i)
		}
	}
	return
}

func (s *Service) UnreadCount(c context.Context, mid int64, plat int8, build int, now time.Time) (total, feedCount, articleCount int) {
	var (
		withoutBangumi = true
		err            error
	)
	if (plat == model.PlatIPhone && build > _iosBanBangumi) || (plat == model.PlatAndroid && build > _androidBanBangumi) || plat == model.PlatIPhoneB {
		withoutBangumi = false
	}
	if feedCount, err = s.upper.AppUnreadCount(c, mid, withoutBangumi); err != nil {
		log.Error("%+v", err)
	}
	if true {
		if articleCount, err = s.upper.ArticleUnreadCount(c, mid); err != nil {
			log.Error("%+v", err)
		}
	}
	total = feedCount + articleCount
	// add feed unread count cache
	if feedCount > 0 {
		s.addCache(func() {
			s.upper.AddUnreadCountCache(context.Background(), mid, feedCount)
		})
	}
	return
}
