package service

import (
	"context"
	"fmt"

	"go-common/app/interface/bbq/bullet/api"
	"go-common/app/interface/bbq/bullet/internal/conf"
	"go-common/app/interface/bbq/bullet/internal/dao"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// ContentGet .
func (s *Service) ContentGet(ctx context.Context, req *api.ListBulletReq) (res []*api.Bullet, err error) {
	videoBase, err := s.dao.VideoBase(ctx, req.Oid)
	if err != nil {
		log.Warnw(ctx, "log", "get video base fail", "err", err, "req", req)
		return
	}
	if video.IsLimitSet(videoBase.Limits, video.VideoLimitBitBullet) {
		log.Infow(ctx, "log", "video limit is set")
		return []*api.Bullet{}, nil
	}

	if res, err = s.dao.ContentGet(ctx, req); err != nil {
		log.Warnv(ctx, log.KV("log", "content get fail"))
		return
	}

	return
}

// ContentList .
func (s *Service) ContentList(ctx context.Context, req *api.ListBulletReq) (res *api.ListBulletReply, err error) {
	videoBase, err := s.dao.VideoBase(ctx, req.Oid)
	if err != nil {
		log.Warnw(ctx, "log", "get video base fail", "err", err, "req", req)
		return
	}
	if video.IsLimitSet(videoBase.Limits, video.VideoLimitBitBullet) {
		log.Infow(ctx, "log", "video limit is set")
		return &api.ListBulletReply{}, nil
	}

	if res, err = s.dao.ContentList(ctx, req); err != nil {
		log.Warnv(ctx, log.KV("log", "get content list fail"))
		return
	}

	return
}

// ContentPost .
func (s *Service) ContentPost(ctx context.Context, req *api.Bullet) (dmid int64, err error) {
	videoBase, err := s.dao.VideoBase(ctx, req.Oid)
	if err != nil {
		log.Warnw(ctx, "log", "get video base fail", "err", err, "req", req)
		return
	}
	if video.IsLimitSet(videoBase.Limits, video.VideoLimitBitBullet) {
		log.Infow(ctx, "log", "video limit is set")
		err = ecode.DanmuLimitErr
		return
	}

	// 屏蔽词
	level, filterErr := s.dao.Filter(ctx, req.Content, dao.FilterAreaDanmu)
	if filterErr != nil {
		log.Errorv(ctx, log.KV("log", "filter fail"))
	} else if level >= dao.FilterLevel {
		err = ecode.FilterErr
		log.Warnv(ctx, log.KV("log", fmt.Sprintf("content filter fail: content=%s, level=%d", req.Content, level)))
		return
	}

	// 发布弹幕
	dmid, err = s.dao.ContentPost(ctx, req)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "publish danmu fail"))
		return
	}
	return
}

//PhoneCheck ..
func (s *Service) PhoneCheck(c context.Context, mid int64) (err error) {
	telStatus, err := s.dao.PhoneCheck(c, mid)
	if err != nil {
		log.Errorw(c, "log", "call phone check fail", "mid", mid)
		return
	}
	if telStatus == 0 {
		err = ecode.BBQNoBindPhone
		log.Infow(c, "log", "no bind phone", "mid", mid)
		return
	}
	return
}
