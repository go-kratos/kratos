package service

import (
	"context"
	"strings"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/log"

	"go-common/library/sync/errgroup"

	"github.com/jaytaylor/html2text"
)

func (s *Service) updateSearchArt(c context.Context, art *model.Article) (err error) {
	if art == nil {
		return
	}
	a := artmdl.Meta{Attributes: art.Attributes}
	if a.AttrVal(artmdl.AttrBitNoDistribute) {
		return
	}
	searchArt := &model.SearchArticle{Article: *art}
	var (
		group  *errgroup.Group
		errCtx context.Context
	)

	group, errCtx = errgroup.WithContext(c)
	group.Go(func() (err error) {
		if searchArt.Content, err = s.dao.ArticleContent(errCtx, art.ID); err != nil {
			return
		}
		searchArt.Content, err = extractText(searchArt.Content)
		return
	})
	group.Go(func() (err error) {
		var stat *artmdl.StatMsg
		if stat, err = s.dao.Stat(errCtx, art.ID); (err != nil) || (stat == nil) {
			return
		}
		searchArt.StatsDisLike = *stat.Dislike
		searchArt.StatsFavorite = *stat.Favorite
		searchArt.StatsLikes = *stat.Like
		searchArt.StatsReply = *stat.Reply
		searchArt.StatsShare = *stat.Share
		searchArt.StatsView = *stat.View
		searchArt.StatsCoin = *stat.Coin
		return
	})
	group.Go(func() (err error) {
		searchArt.Tags, err = s.tags(errCtx, art.ID)
		return
	})
	group.Go(func() (err error) {
		searchArt.Keywords, err = s.dao.Keywords(errCtx, art.ID)
		return
	})
	if err = group.Wait(); err == nil {
		err = s.dao.UpdateSearch(c, searchArt)
	}
	if err != nil {
		dao.PromError("search:更新文章")
		log.Error("updateSearchArt(%+v) err: %+v", searchArt, err)
		return
	}
	log.Info("success updateSearchArt(id: %v, title: %v)", searchArt.ID, searchArt.Title)
	return
}

func (s *Service) delSearchArt(c context.Context, aid int64) (err error) {
	err = s.dao.DelSearch(c, aid)
	return
}

func extractText(html string) (res string, err error) {
	ops := html2text.Options{PrettyTables: false}
	res, err = html2text.FromString(html, ops)
	res = strings.Replace(res, "*", "", -1)
	return
}

func (s *Service) updateSearchStats(c context.Context, stat *artmdl.StatMsg) (err error) {
	if err = s.dao.UpdateSearchStats(c, stat); err != nil {
		log.Error("updateSearchStats(%v) err: %+v", stat.String(), err)
		return
	}
	log.Info("updateSearchStats(%v) success", stat.String())
	return
}
