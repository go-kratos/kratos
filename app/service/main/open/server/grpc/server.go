// Package grpc generate by warden_gen
package grpc

import (
	"context"

	"go-common/app/service/main/open/api/grpc/v1"
	"go-common/app/service/main/open/service"
	"go-common/library/ecode"
	"go-common/library/net/rpc/warden"
)

// New Open warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	w := warden.NewServer(c)
	v1.RegisterOpenServer(w.Server(), &openService{svr})

	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type openService struct {
	svr *service.Service
}

var _ v1.OpenServer = &openService{}

// Ping check dao health.
func (s *openService) Ping(ctx context.Context, req *v1.PingReq) (*v1.PingReply, error) {
	return nil, ecode.MethodNotAllowed
}

// Close close all dao.
func (s *openService) Close(ctx context.Context, req *v1.CloseReq) (*v1.CloseReply, error) {
	return nil, ecode.MethodNotAllowed
}

// Secret .
func (s *openService) Secret(ctx context.Context, req *v1.SecretReq) (*v1.SecretReply, error) {
	return nil, ecode.MethodNotAllowed
}

// AppID .
func (s *openService) AppID(ctx context.Context, req *v1.AppIDReq) (*v1.AppIDReply, error) {
	appID, err := s.svr.AppID(ctx, req.AppKey)
	if err != nil {
		return nil, err
	}
	appIDReply := &v1.AppIDReply{
		AppId: appID,
	}
	return appIDReply, nil
}
