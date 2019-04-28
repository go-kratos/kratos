package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go-common/library/log"

	"go-common/app/service/bbq/sys-msg/api/v1"
	"go-common/app/service/bbq/sys-msg/internal/conf"
	"go-common/app/service/bbq/sys-msg/internal/dao"
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

// ListSysMsg 获取系统通知
func (s *Service) ListSysMsg(ctx context.Context, req *v1.ListSysMsgReq) (res *v1.ListSysMsgReply, err error) {
	res = new(v1.ListSysMsgReply)
	if len(req.Ids) == 0 {
		return
	}
	msgMap, err := s.dao.SysMsg(ctx, req.Ids)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "list sys msg fail: req="+req.String()))
	}
	for _, id := range req.Ids {
		val, exists := msgMap[id]
		if !exists {
			log.Errorv(ctx)
			continue
		}
		// 只返回状态可见的msg，当id=0时表示穿透，也就是无效id
		if val.Id != 0 && val.State == 0 {
			res.List = append(res.List, val)
		} else {
			log.V(1).Infov(ctx, log.KV("log", "no show msg"), log.KV("id", id), log.KV("msg", val.String()))
		}
	}
	log.V(1).Infov(ctx, log.KV("req_size", len(req.Ids)), log.KV("rsp_size", len(res.List)))
	return
}

// CreateSysMsg 创建消息
func (s *Service) CreateSysMsg(ctx context.Context, req *v1.SysMsg) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	err = s.dao.CreateSysMsg(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "create sys msg fail: req="+req.String()))
	}
	return
}

// UpdateSysMsg 更新，一般是状态更新
func (s *Service) UpdateSysMsg(ctx context.Context, req *v1.SysMsg) (res *empty.Empty, err error) {
	res = new(empty.Empty)

	return
}
