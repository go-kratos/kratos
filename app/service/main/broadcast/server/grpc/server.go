// Package server generate by warden_gen
package server

import (
	"context"
	"net"

	pb "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/app/service/main/broadcast/service"
	"go-common/library/conf/paladin"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
	// use gzip decoder
	_ "google.golang.org/grpc/encoding/gzip"
)

// New Zerg warden rpc server
func New(svr *service.Service) (*warden.Server, string) {
	var rc struct {
		Server *warden.ServerConfig
	}
	if err := paladin.Get("grpc.toml").UnmarshalTOML(&rc); err != nil {
		panic(err)
	}
	_, port, _ := net.SplitHostPort(rc.Server.Addr)
	ws := warden.NewServer(rc.Server, grpc.MaxRecvMsgSize(32*1024*1024), grpc.MaxSendMsgSize(32*1024*1024))
	pb.RegisterZergServer(ws.Server(), &server{svr})
	return ws, port
}

type server struct {
	srv *service.Service
}

var _ pb.ZergServer = &server{}

// Ping Service
func (s *server) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingReply, error) {
	return &pb.PingReply{}, nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *pb.CloseReq) (*pb.CloseReply, error) {
	return &pb.CloseReply{}, nil
}

// Connect connect a conn.
func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	mid, key, room, platform, accepts, err := s.srv.Connect(ctx, req.Server, req.ServerKey, req.Cookie, req.Token)
	if err != nil {
		return &pb.ConnectReply{}, err
	}
	return &pb.ConnectReply{Mid: mid, Key: key, RoomID: room, Accepts: accepts, Platform: platform}, nil
}

// Disconnect disconnect a conn.
func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
	has, err := s.srv.Disconnect(ctx, req.Mid, req.Key, req.Server)
	if err != nil {
		return &pb.DisconnectReply{}, err
	}
	return &pb.DisconnectReply{Has: has}, nil
}

// Heartbeat beartbeat a conn.
func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
	if err := s.srv.Heartbeat(ctx, req.Mid, req.Key, req.Server); err != nil {
		return &pb.HeartbeatReply{}, err
	}
	return &pb.HeartbeatReply{}, nil
}

// RenewOnline renew server online.
func (s *server) RenewOnline(ctx context.Context, req *pb.OnlineReq) (*pb.OnlineReply, error) {
	roomCount, err := s.srv.RenewOnline(ctx, req.Server, req.Sharding, req.RoomCount)
	if err != nil {
		return &pb.OnlineReply{}, err
	}
	return &pb.OnlineReply{RoomCount: roomCount}, nil
}

// Receive receive a message.
func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq) (*pb.ReceiveReply, error) {
	if err := s.srv.Receive(ctx, req.Mid, req.Proto); err != nil {
		return &pb.ReceiveReply{}, err
	}
	return &pb.ReceiveReply{}, nil
}

// ServerList return server list.
func (s *server) ServerList(ctx context.Context, req *pb.ServerListReq) (*pb.ServerListReply, error) {
	return s.srv.ServerList(ctx, req.Platform), nil
}
