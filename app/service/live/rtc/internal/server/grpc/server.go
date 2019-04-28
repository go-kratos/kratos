package grpc

import (
	v1pb "go-common/app/service/live/rtc/api/v1"
	"go-common/app/service/live/rtc/internal/conf"
	v1srv "go-common/app/service/live/rtc/internal/service/v1"
	"go-common/library/net/rpc/warden"
)

// TODO

func New(c *conf.Config) *warden.Server {
	ws := warden.NewServer(nil)
	v1pb.RegisterRtcServer(ws.Server(), v1srv.NewRtcService(c))
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
