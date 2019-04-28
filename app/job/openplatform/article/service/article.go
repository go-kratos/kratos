package service

import (
	"context"
	"encoding/json"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/conf/env"
	"go-common/library/log"
)

func (s *Service) upArticles(c context.Context, action string, newMsg []byte, oldMsg []byte) {
	log.Info("s.upArticles action(%s) old(%s) new(%s)", action, string(oldMsg), string(newMsg))
	if action != model.ActUpdate && action != model.ActInsert && action != model.ActDelete {
		return
	}
	var (
		err        error
		aid        int64
		newArticle = &model.Article{}
		oldArticle = &model.Article{}
	)
	if err = json.Unmarshal(newMsg, newArticle); err != nil {
		log.Error("json.Unmarshal(%s) error(%+v)", newMsg, err)
		dao.PromError("article:解析过审文章databus新内容")
		return
	}
	if action == model.ActUpdate {
		if err = json.Unmarshal(oldMsg, oldArticle); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", oldMsg, err)
			dao.PromError("article:解析过审文章databus旧内容")
			return
		}
	}
	aid = newArticle.ID
	mid := newArticle.Mid
	show := true
	var comment string
	if artmdl.NoDistributeAttr(newArticle.Attributes) {
		show = false
		comment = "禁止分发"
	}
	switch action {
	case model.ActInsert:
		s.openReply(c, aid, mid)
		s.addArtCache(c, aid)
		s.updateSearchArt(c, newArticle)
		s.flowSync(c, aid, mid)
		s.addURLNode(c, aid)
		if comment == "" {
			comment = "文章过审"
		}
	case model.ActUpdate:
		s.updateArtCache(c, aid, oldArticle.CategoryID)
		s.addURLNode(c, aid)
		if !artmdl.NoDistributeAttr(oldArticle.Attributes) && artmdl.NoDistributeAttr(newArticle.Attributes) {
			s.delSearchArt(c, aid)
		} else {
			s.updateSearchArt(c, newArticle)
		}
		if comment == "" {
			comment = "文章修改"
		}
	case model.ActDelete:
		s.closeReply(c, aid, mid)
		s.deleteArtCache(c, aid, mid)
		s.deleteArtRecommendCache(c, aid, newArticle.CategoryID)
		s.delSearchArt(c, aid)
		// s.delMediaScore(c, aid, newArticle.MediaID, mid)
		show = false
		comment = "文章不可见"
	}
	if e := s.dao.PubDynamic(c, mid, aid, show, comment, int64(newArticle.PublishTime), newArticle.DynamicIntro); e != nil {
		s.dao.PushDynamicCache(c, &model.DynamicCacheRetry{
			Aid:          aid,
			Mid:          mid,
			Show:         show,
			Comment:      comment,
			Ts:           int64(newArticle.PublishTime),
			DynamicIntro: newArticle.DynamicIntro,
		})
	}
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	s.gameSync(c, aid, mid, action)
	urls := model.ReadURLs(aid)
	for _, u := range urls {
		if err = s.dao.PurgeCDN(c, u); err == nil {
			log.Info("s.dao.PurgeCDN(%s) success.", u)
			dao.PromInfo("article:刷新CDN")
		}
	}
}

func (s *Service) gameSync(c context.Context, aid, mid int64, action string) {
	authors, err := s.dao.GameList(c)
	if err != nil {
		dao.PromError("service:获得游戏数据")
		log.Error("s.gameSync(aid: %v, mid: %v) init err: %+v", aid, mid, err)
		return
	}
	var exist bool
	for _, author := range authors {
		if author == mid {
			exist = true
			break
		}
	}
	if !exist {
		return
	}
	if err = s.dao.GameSync(c, action, aid); err != nil {
		log.Error("s.gameSync(%d, %d, %s) error(%+v)", aid, mid, action, err)
		dao.PromError("service:同步游戏数据")
		s.dao.PushGameCache(c, &model.GameCacheRetry{
			Action: action,
			Aid:    aid,
		})
		return
	}
	log.Info("s.gameSync(%d, %d, %s) success", aid, mid, action)
}

func (s *Service) flowSync(c context.Context, aid, mid int64) {
	if err := s.dao.FlowSync(c, mid, aid); err != nil {
		s.dao.PushFlowCache(c, &model.FlowCacheRetry{
			Aid: aid,
			Mid: mid,
		})
		return
	}
	log.Info("s.flowSync(aid: %d, mid: %d) success", aid, mid)
}

func (s *Service) openReply(c context.Context, aid, mid int64) (err error) {
	if err = s.dao.OpenReply(c, aid, mid); err == nil {
		log.Info("OpenReply(%d,%d) success.", aid, mid)
		dao.PromInfo("article:打开评论区")
	}
	return
}

func (s *Service) closeReply(c context.Context, aid, mid int64) (err error) {
	if err = s.dao.CloseReply(c, aid, mid); err == nil {
		log.Info("CloseReply(%d,%d) success.", aid, mid)
		dao.PromInfo("article:关闭评论区")
	}
	return
}

func (s *Service) addArtCache(c context.Context, aid int64) (err error) {
	arg := &artmdl.ArgAid{Aid: aid}
	if err = s.articleRPC.AddArticleCache(c, arg); err != nil {
		log.Error("s.articleRPC.AddArticleCache(%d) error(%+v)", aid, err)
		dao.PromError("article:新增文章缓存")
		s.dao.PushArtCache(c, &dao.CacheRetry{
			Action: dao.RetryAddArtCache,
			Aid:    aid,
		})
	} else {
		log.Info("s.articleRPC.AddArticleCache(%d) success", aid)
		dao.PromInfo("article:新增文章缓存")
	}
	return
}

func (s *Service) updateArtCache(c context.Context, aid, cid int64) (err error) {
	arg := &artmdl.ArgAidCid{Aid: aid, Cid: cid}
	if err = s.articleRPC.UpdateArticleCache(c, arg); err != nil {
		log.Error("s.articleRPC.UpdateArticleCache(%d,%d) error(%+v)", aid, cid, err)
		dao.PromError("article:更新文章缓存")
		s.dao.PushArtCache(c, &dao.CacheRetry{
			Action: dao.RetryUpdateArtCache,
			Aid:    aid,
			Cid:    cid,
		})
	} else {
		log.Info("s.articleRPC.UpdateArticleCache(%d,%d) success", aid, cid)
		dao.PromInfo("article:更新文章缓存")
	}
	return
}

func (s *Service) deleteArtCache(c context.Context, aid, mid int64) (err error) {
	arg := &artmdl.ArgAidMid{Aid: aid, Mid: mid}
	if err = s.articleRPC.DelArticleCache(c, arg); err != nil {
		log.Error("s.articleRPC.DelArticleCache(%d,%d) error(%+v)", aid, mid, err)
		dao.PromError("article:删除文章缓存")
		s.dao.PushArtCache(c, &dao.CacheRetry{
			Action: dao.RetryDeleteArtCache,
			Aid:    aid,
			Mid:    mid,
		})
	} else {
		log.Info("s.articleRPC.DelArticleCache(%d,%d) success", aid, mid)
		dao.PromInfo("article:删除文章缓存")
	}
	return
}

func (s *Service) deleteArtRecommendCache(c context.Context, aid, cid int64) (err error) {
	arg := &artmdl.ArgAidCid{Aid: aid, Cid: cid}
	if err = s.articleRPC.DelRecommendArtCache(c, arg); err != nil {
		log.Error("s.articleRPC.DelRecommendArtCache(%d,%d) error(%+v)", aid, cid, err)
		dao.PromError("article:删除文章推荐缓存")
		s.dao.PushArtCache(c, &dao.CacheRetry{
			Action: dao.RetryDeleteArtRecCache,
			Aid:    aid,
			Cid:    cid,
		})
	} else {
		log.Info("s.articleRPC.DelRecommendArtCache(%d,%d) success", aid, cid)
		dao.PromInfo("article:删除文章推荐缓存")
	}
	return
}

func (s *Service) delMediaScore(c context.Context, aid, mediaID, mid int64) (err error) {
	err = s.dao.DelScore(c, aid, mediaID, mid)
	return
}
