package v1

import (
	"context"
	v1pb "go-common/app/service/live/xuserex/api/grpc/v1"
	"go-common/app/service/live/xuserex/conf"
	"go-common/app/service/live/xuserex/dao/notice"
	"go-common/library/ecode"
	"go-common/library/log"
)

// RoomNoticeService struct
type RoomNoticeService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *notice.Dao
}

//NewRoomNoticeService init
func NewRoomNoticeService(c *conf.Config) (s *RoomNoticeService) {
	s = &RoomNoticeService{
		conf: c,
		dao:  notice.New(c),
	}
	return s
}

// 房间提示 相关服务

// BuyGuard implementation
// 是否弹出大航海购买提示
func (s *RoomNoticeService) BuyGuard(ctx context.Context, req *v1pb.RoomNoticeBuyGuardReq) (resp *v1pb.RoomNoticeBuyGuardResp, err error) {
	resp = &v1pb.RoomNoticeBuyGuardResp{}
	UID := req.GetUid()
	targetID := req.GetTargetId()

	if UID <= 0 || targetID <= 0 {
		err = ecode.NoLogin
		return resp, err
	}
	resp, err = s.dao.IsNotice(ctx, UID, targetID)

	return
}

// IsTaskFinish implementation
// habse 任务是否结束
func (s *RoomNoticeService) IsTaskFinish(ctx context.Context, req *v1pb.RoomNoticeIsTaskFinishReq) (resp *v1pb.RoomNoticeIsTaskFinishResp, err error) {
	resp = &v1pb.RoomNoticeIsTaskFinishResp{}
	ret, err := s.dao.GetTaskFinish(ctx, s.dao.GetTermBegin())
	if err != nil {
		log.Error("s.dao.GetTaskFinish(%v) error(%v)", s.dao.GetTermBegin(), err)
		return
	}
	if ret {
		resp.Result = 1
	}
	return
}

// SetTaskFinish implementation
// 手动设置base 任务结束
func (s *RoomNoticeService) SetTaskFinish(ctx context.Context, req *v1pb.RoomNoticeSetTaskFinishReq) (resp *v1pb.RoomNoticeSetTaskFinishResp, err error) {
	resp = &v1pb.RoomNoticeSetTaskFinishResp{}
	isFinish := req.GetResult()
	log.Info("SetTaskFinish finish (%v) term (%v)", isFinish, s.dao.GetTermBegin())
	err = s.dao.SetTaskFinish(ctx, s.dao.GetTermBegin(), isFinish)

	if err != nil {
		log.Error("SetTaskFinish(%v) term (%v) error(%v)", isFinish, s.dao.GetTermBegin(), err)
		return
	}

	resp.Result = 1
	return
}
