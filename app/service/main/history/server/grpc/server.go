package grpc

import (
	pb "go-common/app/service/main/history/api/grpc"
	"go-common/app/service/main/history/service"
	"go-common/library/net/rpc/warden"
)

// New Coin warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterHistoryServer(ws.Server(), svr)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
