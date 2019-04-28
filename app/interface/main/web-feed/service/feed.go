package service

import (
	"context"

	"go-common/app/interface/main/web-feed/dao"
	"go-common/app/interface/main/web-feed/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Feed get feed of ups and bangumi.
func (s *Service) Feed(c context.Context, mid int64, pn, ps int) (res []*model.Feed, err error) {
	var (
		feeds, newFeeds []*feedmdl.Feed
		upAids          []int64
		accRes          map[int64]*account.Card
	)
	arg := &feedmdl.ArgFeed{
		Mid:    mid,
		Pn:     pn,
		Ps:     ps,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if feeds, err = s.feedRPC.WebFeed(c, arg); err != nil || len(feeds) == 0 {
		log.Error("s.feedRPC.WebFeed(%v) error(%v)", arg, err)
		if pn == 1 {
			res, err = s.dao.FeedCache(c, mid)
			log.Info("s.dao.FeedCache(%d) len(%d) error(%v)", mid, len(res), err)
		}
		return
	}
	for _, item := range feeds {
		if (item.Type == feedmdl.BangumiType) && (item.Bangumi == nil) {
			dao.PromError("bangumi为空")
			log.Error("s.feedRPC.WebFeed(%v) error(%v)", arg, err)
			continue
		}
		if (item.Type == feedmdl.ArchiveType) && (item.Archive == nil) {
			dao.PromError("archive为空")
			log.Error("s.feedRPC.WebFeed(%v) error(%v)", arg, err)
			continue
		}
		newFeeds = append(newFeeds, item)
	}
	feeds = newFeeds
	for _, item := range feeds {
		if item.Type != feedmdl.ArchiveType {
			continue
		}
		if item.Archive != nil {
			upAids = append(upAids, item.Archive.Author.Mid)
		}
	}
	accArg := &account.ArgMids{Mids: upAids}
	if accRes, err = s.accRPC.Cards3(c, accArg); err != nil {
		dao.PromError("rpc:accRPC.Infos2")
		log.Error("Feed s.accRPC.info(%v) error(%v)", arg, err)
		err = nil
	}
	for _, item := range feeds {
		tmp := model.Feed{Feed: item}
		if tmp.Type == feedmdl.ArchiveType && accRes != nil {
			if ai, ok := accRes[item.Archive.Author.Mid]; ok {
				tmp.OfficialVerify = &ai.Official
			}
		}
		res = append(res, &tmp)
	}
	if pn == 1 {
		s.cache.Save(func() {
			s.dao.SetFeedCache(context.TODO(), mid, res)
		})
	}
	return
}

// UnreadCount get unread count of feed
func (s *Service) UnreadCount(c context.Context, mid int64) (count int, err error) {
	arg := &feedmdl.ArgMid{
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if count, err = s.feedRPC.WebUnreadCount(c, arg); err != nil {
		dao.PromError("rpc:feedRPC.WebUnreadCount")
		log.Error("s.feedRPC.UnreadCount(%v) error(%v)", arg, err)
	}
	return
}

// ArticleUnreadCount get unread count of feed
func (s *Service) ArticleUnreadCount(c context.Context, mid int64) (count int, err error) {
	arg := &feedmdl.ArgMid{
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if count, err = s.feedRPC.ArticleUnreadCount(c, arg); err != nil {
		dao.PromError("feed:ArticleUnreadCount")
		log.Error("s.feedRPC.ArticleUnreadCount(%v) error(%v)", arg, err)
	}
	return
}

// ArticleFeed get feed of ups and bangumi.
func (s *Service) ArticleFeed(c context.Context, mid int64, pn, ps int) (res []*artmdl.Meta, err error) {
	arg := &feedmdl.ArgFeed{
		Mid:    mid,
		Pn:     pn,
		Ps:     ps,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.feedRPC.ArticleFeed(c, arg); err != nil {
		log.Error("s.feedRPC.ArticleFeed(%v) error(%v)", arg, err)
		return
	}
	var mids []int64
	var accRes map[int64]*account.Card
	for _, meta := range res {
		if meta.Author != nil {
			mids = append(mids, meta.Author.Mid)
		}
	}
	accArg := &account.ArgMids{
		Mids: mids,
	}
	if accRes, err = s.accRPC.Cards3(c, accArg); err != nil {
		dao.PromError("rpc:accRPC.Infos2")
		log.Error("Feed s.accRPC.info(%v) error(%v)", arg, err)
		err = nil
		return
	}
	for _, item := range res {
		if (item.Author == nil) || (accRes[item.Author.Mid] == nil) {
			continue
		}
		info := accRes[item.Author.Mid]
		item.Author = &artmdl.Author{
			Mid:  item.Author.Mid,
			Name: info.Name,
			Face: info.Face,
			Pendant: artmdl.Pendant{
				Pid:    int32(info.Pendant.Pid),
				Name:   info.Pendant.Name,
				Image:  info.Pendant.Image,
				Expire: int32(info.Pendant.Expire),
			},
			Nameplate: artmdl.Nameplate{
				Nid:        info.Nameplate.Nid,
				Name:       info.Nameplate.Name,
				Image:      info.Nameplate.Image,
				ImageSmall: info.Nameplate.ImageSmall,
				Level:      info.Nameplate.Level,
				Condition:  info.Nameplate.Condition,
			},
		}
		if info.Official.Role == 0 {
			item.Author.OfficialVerify.Type = -1
		} else {
			if info.Official.Role <= 2 {
				item.Author.OfficialVerify.Type = 0
			} else {
				item.Author.OfficialVerify.Type = 1
			}
			item.Author.OfficialVerify.Desc = info.Official.Title
		}
	}
	return
}
