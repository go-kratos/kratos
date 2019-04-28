package v1

import (
	"context"
	v1pb "go-common/app/interface/live/app-room/api/http/v1"
	"go-common/app/interface/live/app-room/conf"
	"go-common/app/interface/live/app-room/dao"
	"go-common/app/service/live/xuserex/api/grpc/v1"
	"go-common/library/ecode"
)

// RoomNoticeService struct
type RoomNoticeService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewRoomNoticeService init
func NewRoomNoticeService(c *conf.Config) (s *RoomNoticeService) {
	s = &RoomNoticeService{
		conf: c,
		dao:  dao.New(c),
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
		err = ecode.ParamInvalid
		return
	}

	ret, err := s.dao.XuserexAPI.BuyGuard(ctx, &v1.RoomNoticeBuyGuardReq{
		Uid:      UID,
		TargetId: targetID,
	})
	if err != nil {
		return
	}
	if ret == nil {
		return
	}
	resp = &v1pb.RoomNoticeBuyGuardResp{
		ShouldNotice: ret.GetShouldNotice(),
		Begin:        ret.GetBegin(),
		End:          ret.GetEnd(),
		Now:          ret.GetNow(),
		Title:        ret.GetTitle(),
		Content:      ret.GetContent(),
		Button:       ret.GetButton(),
	}
	return
}
