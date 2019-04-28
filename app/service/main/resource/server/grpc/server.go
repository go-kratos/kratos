// Package server generate by warden_gen
package server

import (
	pb "go-common/app/service/main/resource/api/v1"
	"go-common/app/service/main/resource/service"
	"go-common/library/net/rpc/warden"
)

// New Coin warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterResourceServer(ws.Server(), svr)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
