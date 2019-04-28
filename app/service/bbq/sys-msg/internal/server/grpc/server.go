package grpc

import (
	pb "go-common/app/service/bbq/sys-msg/api/v1"
	"go-common/app/service/bbq/sys-msg/internal/service"
	"go-common/library/net/rpc/warden"
)

// New new warden rpc server
func New(c *warden.ServerConfig, svc *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterSysMsgServiceServer(ws.Server(), svc)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
