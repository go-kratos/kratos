package grpc

import (
	"context"

	v1 "go-common/app/interface/main/creative/api"
	newcMdl "go-common/app/interface/main/creative/model/newcomer"
	"go-common/app/interface/main/creative/service/archive"
	"go-common/app/interface/main/creative/service/newcomer"
	"go-common/library/net/rpc/warden"
)

// New boradcast grpc server.
func New(cfg *warden.ServerConfig, asrv *archive.Service, nsrv *newcomer.Service) (ws *warden.Server) {
	ws = warden.NewServer(cfg)
	v1.RegisterCreativeServer(ws.Server(), &server{asrv, nsrv})

	var err error
	if ws, err = ws.Start(); err != nil {
		panic(err)
	}
	return
}

type server struct {
	asrv *archive.Service
	nsrv *newcomer.Service
}

// Ping Service
func (s *server) Ping(ctx context.Context, req *v1.Empty) (*v1.Empty, error) {
	return &v1.Empty{}, nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *v1.Empty) (*v1.Empty, error) {
	// TODO: some graceful close
	return &v1.Empty{}, nil
}

// AddReport .
func (s *server) CheckTaskState(c context.Context, arg *v1.TaskRequest) (res *v1.TaskReply, err error) {
	req := &newcMdl.CheckTaskStateReq{
		MID:    arg.Mid,
		TaskID: arg.TaskId,
	}
	return &v1.TaskReply{FinishState: s.nsrv.CheckTaskState(c, req)}, nil

}

// FlowJudge .
func (s *server) FlowJudge(c context.Context, arg *v1.FlowRequest) (res *v1.FlowResponse, err error) {
	res = new(v1.FlowResponse)
	a, err := s.asrv.FlowJudge(c, arg.Business, arg.Gid, arg.Oids)
	if err != nil {
		return
	}
	res.Oids = a
	return
}
