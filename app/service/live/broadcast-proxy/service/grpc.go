package service

import (
	"errors"
	v1pb "go-common/app/service/live/broadcast-proxy/api/v1"
	"go-common/app/service/live/broadcast-proxy/server"
	v1srv "go-common/app/service/live/broadcast-proxy/service/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
	"google.golang.org/grpc"
	"time"
)

func NewGrpcService(p *server.BroadcastProxy, d *server.CometDispatcher) (*warden.Server, error) {
	if p == nil || d == nil {
		return nil, errors.New("empty proxy")
	}
	ws := warden.NewServer(&warden.ServerConfig{
		Timeout: xtime.Duration(30 * time.Second),
	}, grpc.MaxRecvMsgSize(1024 * 1024 * 1024), grpc.MaxSendMsgSize(1024 * 1024 * 1024))
	v1pb.RegisterDanmakuServer(ws.Server(), v1srv.NewDanmakuService(p, d))
	ws, err := ws.Start()
	if err != nil {
		return nil, err
	}
	return ws, nil
}
