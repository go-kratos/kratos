// Package server generate by warden_gen
package server

import (
	"context"
	"time"

	pb "go-common/app/interface/main/broadcast/api/grpc/v1"
	"go-common/app/interface/main/broadcast/conf"
	service "go-common/app/interface/main/broadcast/server"
	"go-common/library/net/rpc/warden"
)

// New boradcast grpc server.
func New(c *conf.Config, svr *service.Server) (ws *warden.Server) {
	var (
		err error
	)
	ws = warden.NewServer(c.WardenServer)
	pb.RegisterZergServer(ws.Server(), &server{svr})
	if ws, err = ws.Start(); err != nil {
		panic(err)
	}
	return
}

type server struct {
	srv *service.Server
}

var _ pb.ZergServer = &server{}

// Ping Service
func (s *server) Ping(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	// TODO: some graceful close
	return &pb.Empty{}, nil
}

// PushMsg push a message to specified sub keys.
func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq) (reply *pb.PushMsgReply, err error) {
	if len(req.Keys) == 0 || req.Proto == nil {
		return nil, service.ErrPushMsgArg
	}
	for _, key := range req.Keys {
		if channel := s.srv.Bucket(key).Channel(key); channel != nil {
			if !channel.NeedPush(req.ProtoOp, "") {
				continue
			}
			if err = channel.Push(req.Proto); err != nil {
				return
			}
		}
	}
	return &pb.PushMsgReply{}, nil
}

// Broadcast broadcast msg to all user.
func (s *server) Broadcast(ctx context.Context, req *pb.BroadcastReq) (*pb.BroadcastReply, error) {
	if req.Proto == nil {
		return nil, service.ErrBroadCastArg
	}
	go func() {
		for _, bucket := range s.srv.Buckets() {
			bucket.Broadcast(req.GetProto(), req.ProtoOp, req.Platform)
			if req.Speed > 0 {
				t := bucket.ChannelCount() / int(req.Speed)
				time.Sleep(time.Duration(t) * time.Second)
			}
		}
	}()

	return &pb.BroadcastReply{}, nil
}

// BroadcastRoom broadcast msg to specified room.
func (s *server) BroadcastRoom(ctx context.Context, req *pb.BroadcastRoomReq) (*pb.BroadcastRoomReply, error) {
	if req.Proto == nil || req.RoomID == "" {
		return nil, service.ErrBroadCastRoomArg
	}
	for _, bucket := range s.srv.Buckets() {
		bucket.BroadcastRoom(service.ProtoRoom{
			RoomID: req.RoomID,
			Proto:  req.Proto,
		})
	}
	return &pb.BroadcastRoomReply{}, nil
}

// Rooms gets all the room ids for the server.
func (s *server) Rooms(ctx context.Context, req *pb.RoomsReq) (*pb.RoomsReply, error) {
	var (
		roomIds = make(map[string]bool)
	)
	for _, bucket := range s.srv.Buckets() {
		for roomID := range bucket.Rooms() {
			roomIds[roomID] = true
		}
	}
	return &pb.RoomsReply{Rooms: roomIds}, nil
}
