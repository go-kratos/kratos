package service

import (
	"context"
	"encoding/json"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/log"
)

func (s *Service) upAuthors(c context.Context, action string, newMsg []byte, oldMsg []byte) {
	log.Info("s.upAuthors action(%s) old(%s) new(%s)", action, string(oldMsg), string(newMsg))
	var (
		err       error
		newAuthor = &model.Author{}
	)
	if err = json.Unmarshal(newMsg, newAuthor); err != nil {
		log.Error("json.Unmarshal(%s) error(%+v)", newMsg, err)
		dao.PromError("article:解析作者表databus新内容")
		return
	}
	s.updateAuthorCache(c, newAuthor.Mid)
}

// updateAuthorCache update author cache
func (s *Service) updateAuthorCache(c context.Context, mid int64) (err error) {
	arg := &artmdl.ArgAuthor{Mid: mid}
	if err = s.articleRPC.UpdateAuthorCache(c, arg); err != nil {
		log.Error("s.articleRPC.UpdateAuthorCache(%+v) error(%+v)", arg, err)
		dao.PromError("article:更新作者缓存")
		return
	}
	log.Info("s.articleRPC.UpdateAuthorCache(%+v) success", arg)
	dao.PromInfo("article:更新作者缓存")
	return
}
