package grpc

import (
	pb "go-common/app/service/main/sms/api"
	"go-common/app/service/main/sms/service"
	"go-common/library/net/rpc/warden"
)

// New Sms warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterSmsServer(ws.Server(), svr)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
