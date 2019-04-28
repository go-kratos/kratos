package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/notice-service/internal/conf"
	"go-common/app/service/bbq/notice-service/internal/dao"
	"go-common/library/log"
	"go-common/library/log/infoc"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service struct
type Service struct {
	c     *conf.Config
	dao   *dao.Dao
	Infoc *infoc.Infoc
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		dao:   dao.New(c),
		Infoc: infoc.New(c.Infoc),
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

// ListNotices 获取通知列表
func (s *Service) ListNotices(c context.Context, req *v1.ListNoticesReq) (res *v1.ListNoticesReply, err error) {
	res = new(v1.ListNoticesReply)
	res.Mid = req.Mid
	res.List, err = s.dao.ListNotices(c, req.Mid, req.CursorId, req.NoticeType)
	if err != nil {
		log.Errorv(c, log.KV("log", "get list notice fail"))
		return
	}
	// 清理未读数
	s.dao.ClearUnread(c, req.Mid, req.NoticeType)

	return
}

// CreateNotice 创建消息
func (s *Service) CreateNotice(c context.Context, req *v1.NoticeBase) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	id, err := s.dao.CreateNotice(c, req)
	if err != nil {
		log.Warnv(c, log.KV("log", "create notice fail"))
		return
	}
	// 未读数+1
	err = s.dao.IncreaseUnread(c, req.Mid, req.NoticeType, 1)
	if err != nil {
		log.Warnv(c, log.KV("log", "increase notice unread fail: req="+req.String()))
		err = nil
	}
	// TODO:推送
	s.pushNotification(c, id, req)
	return
}

// GetUnreadInfo 获取未读情况
func (s *Service) GetUnreadInfo(ctx context.Context, req *v1.GetUnreadInfoRequest) (res *v1.UnreadInfo, err error) {
	res = new(v1.UnreadInfo)
	res.List, err = s.dao.GetUnreadInfo(ctx, req.Mid)
	if err != nil {
		log.Warnv(ctx, log.KV("log", "get notice unread info fail"), log.KV("req", req.String()))
		return
	}
	log.V(1).Infov(ctx, log.KV("log", "get unread info: res="+res.String()))
	return
}

// PushCallback 推送回调
func (s *Service) PushCallback(ctx context.Context, req *v1.PushCallbackRequest) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	// 埋点
	s.Infoc.Info(req.Tid, req.Mid, req.Buvid, req.Nid, "", "", "", time.Now().Unix(), "")
	return
}

// PushLogout 推送回调
func (s *Service) PushLogout(ctx context.Context, req *v1.UserPushDev) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	var rowsAffected int64
	if rowsAffected, err = s.dao.DeleteUserPushDev(ctx, req); err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("delete user push dev fail: req=%s", req.String())))
		return
	}

	if rowsAffected == 0 {
		log.Warnv(ctx, log.KV("log", fmt.Sprintf("delete user push dev fail due to affected rows is zero: req=%s", req.String())))
	}

	return
}

// PushLogin login
func (s *Service) PushLogin(ctx context.Context, req *v1.UserPushDev) (res *empty.Empty, err error) {
	res = new(empty.Empty)

	// 获取数据库是否存在当前的mid & buvid
	dev, err := s.dao.FetchUserPushDev(ctx, req.Mid, req.Buvid)
	if err != nil {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("fetch user push dev fail: req=%s, err=%s", req.String(), err.Error())))
		return
	}

	// 插入user_push_device
	if dev.Id == 0 {
		if _, err = s.dao.InsertUserPushDev(ctx, req); err != nil {
			log.Errorv(ctx, log.KV("log", fmt.Sprintf("insert user push dev fail: req=%s", req.String())))
			return
		}
	} else {
		// 更新user_push_device
		if _, err = s.dao.UpdateUserPushDev(ctx, req); err != nil {
			log.Errorv(ctx, log.KV("log", fmt.Sprintf("insert user push dev fail: req=%s", req.String())))
			return
		}
	}

	return
}
