package grpc

import (
	"context"

	pb "go-common/app/admin/main/manager/api"
	"go-common/app/admin/main/manager/service"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// server .
type server struct {
	srv *service.Service
}

// New return warden server.
func New(cfg *warden.ServerConfig, srv *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	pb.RegisterPermitServer(w.Server(), &server{srv: srv})
	var err error
	if w, err = w.Start(); err != nil {
		panic(err)
	}
	return w
}

// Login whether login .
func (s *server) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	sid, uname, err := s.srv.Login(ctx, req.Mngsid, req.Dsbsid)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{Sid: sid, Username: uname}, nil
}

// Permissions .
func (s *server) Permissions(ctx context.Context, req *pb.PermissionReq) (*pb.PermissionReply, error) {
	tmp, err := s.srv.Permissions(ctx, req.Username)
	if err != nil {
		log.Error("s.Permissions error(%v)", err)
		return nil, err
	}
	return &pb.PermissionReply{Uid: tmp.UID, Perms: tmp.Perms}, nil
}
