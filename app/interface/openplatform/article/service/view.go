package service

import (
	"context"
	"strconv"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/library/sync/errgroup"
)

// ViewInfo get view info
func (s *Service) ViewInfo(c context.Context, mid, id int64, ip string, cheat *model.CheatInfo, plat int8, from string) (res *model.ViewInfo, err error) {
	var art *model.Meta
	group := &errgroup.Group{}
	if art, err = s.ArticleMeta(c, id); (err != nil) || (art == nil) || (!art.IsNormal()) {
		err = ecode.NothingFound
		return
	}
	res = &model.ViewInfo{
		Title:           art.Title,
		BannerURL:       art.BannerURL,
		Mid:             art.Author.Mid,
		ImageURLs:       art.ImageURLs,
		OriginImageURLs: art.OriginImageURLs,
		ShowLaterWatch:  s.setting.ShowLaterWatch,
		ShowSmallWindow: s.setting.ShowSmallWindow,
	}
	res.Shareable = !art.AttrVal(model.AttrBitNoDistribute)
	if mid > 0 {
		res.IsAuthor, _, _ = s.IsAuthor(c, mid)
		group.Go(func() error {
			res.Like, _ = s.isLike(c, mid, id)
			return nil
		})
		group.Go(func() error {
			if art.Author != nil {
				res.Attention, _ = s.isAttention(c, mid, art.Author.Mid)
			}
			return nil
		})
		group.Go(func() error {
			res.Favorite, _ = s.IsFav(c, mid, id)
			return nil
		})
		group.Go(func() error {
			res.Coin, _ = s.Coin(c, mid, id, ip)
			return nil
		})
	}
	group.Go(func() error {
		if stat, e := s.stat(c, id); (e == nil) && (stat != nil) {
			res.Stats = *stat
			res.Stats.Dynamic, _ = s.dao.DynamicCount(c, id)
		}
		return nil
	})
	group.Go(func() error {
		var lid int64
		lists, _ := s.dao.ArtsList(c, []int64{id})
		if lists[id] != nil {
			res.InList = true
			lid = lists[id].ID
		}
		if mid > 0 {
			s.AddHistory(c, mid, id, lid, ip, plat, from)
		}
		return nil
	})
	group.Go(func() error {
		info, _ := s.authorDetail(c, art.Author.Mid)
		if info != nil {
			res.AuthorName = info.Name
		}
		return nil
	})
	cache.Save(func() {
		if mid == 0 || from == "articleSlide" {
			return
		}
		if info, _ := s.accountInfo(context.TODO(), mid); info != nil {
			cheat.Lv = strconv.FormatInt(int64(info.Level), 10)
		}
		s.dao.PubView(context.TODO(), mid, id, ip, cheat)
	})
	group.Wait()
	return
}

func (s *Service) isAttention(c context.Context, mid, up int64) (ok bool, err error) {
	arg := account.ArgRelation{Mid: mid, Owner: up}
	relation, err := s.accountRPC.Relation3(c, &arg)
	if err != nil {
		dao.PromError("view:获取关注列表")
		log.Error("s.accountRPC.Relation2(%+v) err: %+v", arg, err)
		return
	}
	ok = relation.Following
	return
}

func (s *Service) isAttentions(c context.Context, mid int64, ups []int64) (res map[int64]bool, err error) {
	arg := account.ArgRelations{Mid: mid, Owners: ups}
	relations, err := s.accountRPC.Relations3(c, &arg)
	if err != nil {
		dao.PromError("view:批量获取关注列表")
		log.Error("s.accountRPC.Relations3(%+v) err: %+v", arg, err)
		return
	}
	res = make(map[int64]bool)
	for id, r := range relations {
		res[id] = r.Following
	}
	return
}

func (s *Service) isBlacks(c context.Context, mid int64, ups []int64) (res map[int64]struct{}, err error) {
	arg := account.ArgMid{Mid: mid}
	res, err = s.accountRPC.Blacks3(c, &arg)
	if err != nil {
		dao.PromError("view:获取黑名单列表")
		log.Error("s.accountRPC.Blacks3(%+v) err: %+v", arg, err)
		return
	}
	return
}

func (s *Service) checkArticle(c context.Context, id int64) (err error) {
	var art *model.Meta
	if art, err = s.ArticleMeta(c, id); (err != nil) || (art == nil) {
		err = ecode.NothingFound
		return
	}
	return
}
