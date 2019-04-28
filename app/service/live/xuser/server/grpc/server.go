package grpc

import (
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/dao"
	"go-common/app/service/live/xuser/service"
	"go-common/library/net/rpc/warden"
)

// New new grpc server
func New(svc *service.Service) (wsvr *warden.Server, err error) {
	wsvr = warden.NewServer(nil)
	dao.InitAPI()

	v1pb.RegisterVipServer(wsvr.Server(), svc.VipV1Svc())
	v1pb.RegisterGuardServer(wsvr.Server(), svc.GuardV1Svc())
	v1pb.RegisterUserExpServer(wsvr.Server(), svc.ExpV1Svc())
	v1pb.RegisterRoomAdminServer(wsvr.Server(), svc.RoomAdminV1Svc())

	if wsvr, err = wsvr.Start(); err != nil {
		return
	}
	return
}
