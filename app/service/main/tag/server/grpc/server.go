package grpc

import (
	"context"

	v1 "go-common/app/service/main/tag/api"
	"go-common/app/service/main/tag/model"
	"go-common/app/service/main/tag/service"
	"go-common/library/net/rpc/warden"
)

type grpcServer struct {
	svr *service.Service
}

// New new a grpc server.
func New(cfg *warden.ServerConfig, svr *service.Service) *warden.Server {
	grpc := warden.NewServer(cfg)
	v1.RegisterTagRPCServer(grpc.Server(), &grpcServer{svr: svr})
	grpc, err := grpc.Start()
	if err != nil {
		panic(err)
	}
	return grpc
}

// AddReport .
func (s *grpcServer) AddReport(c context.Context, arg *v1.AddReportReq) (res *v1.AddReportReply, err error) {
	req := &model.AddReportReq{
		Oid:      arg.Oid,
		Type:     arg.Type,
		Tid:      arg.Tid,
		Mid:      arg.Mid,
		PartID:   arg.PartId,
		ReasonID: arg.ReasonId,
		Score:    arg.Score,
	}
	return s.svr.AddReport(c, req)
}
