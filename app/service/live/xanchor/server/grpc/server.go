package grpc

import (
	v1pb "go-common/app/service/live/xanchor/api/grpc/v1"
	"go-common/app/service/live/xanchor/service"
	"go-common/library/net/rpc/warden"
)

// New new grpc server
func New(svc *service.Service) (wsvr *warden.Server, err error) {
	wsvr = warden.NewServer(nil)
	v1pb.RegisterXAnchorServer(wsvr.Server(), svc.V1Svc())
	if wsvr, err = wsvr.Start(); err != nil {
		return
	}
	return
}
