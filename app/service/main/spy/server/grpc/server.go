// Package server generate by warden_gen
package server

import (
	"context"

	pb "go-common/app/service/main/spy/api"
	"go-common/app/service/main/spy/model"
	service "go-common/app/service/main/spy/service"
	"go-common/library/net/rpc/warden"
)

// New Spy warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterSpyServer(ws.Server(), &server{svr})
	_, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ pb.SpyServer = &server{}

// Ping check dao health.
func (s *server) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingReply, error) {
	return &pb.PingReply{}, nil
}

// StatByID spy stat by id or mid.
func (s *server) StatByID(ctx context.Context, req *pb.StatByIDReq) (*pb.StatByIDReply, error) {
	statistics, err := s.svr.StatByID(ctx, req.Mid, req.Id)
	if err != nil {
		return nil, err
	}
	reply := new(pb.StatByIDReply)
	reply.DeepCopyFromStatistics(statistics)
	return reply, nil
}

// StatByIDGroupEvent spy stat by id or mid.
func (s *server) StatByIDGroupEvent(ctx context.Context, req *pb.StatByIDGroupEventReq) (*pb.StatByIDGroupEventReply, error) {
	statistics, err := s.svr.StatByIDGroupEvent(ctx, req.Mid, req.Id)
	if err != nil {
		return nil, err
	}
	reply := new(pb.StatByIDGroupEventReply)
	reply.DeepCopyFromStatistics(statistics)
	return reply, nil
}

// PurgeUser purge  user info
func (s *server) PurgeUser(ctx context.Context, req *pb.PurgeUserReq) (*pb.PurgeUserReply, error) {
	return &pb.PurgeUserReply{}, s.svr.PurgeUser(ctx, req.Mid, req.Action)
}

// HandleEvent handle spy-event.
func (s *server) HandleEvent(ctx context.Context, req *pb.HandleEventReq) (*pb.HandleEventReply, error) {
	eventMsg := new(model.EventMessage)
	req.DeepCopyAsIntoEventMessage(eventMsg)
	return &pb.HandleEventReply{}, s.svr.HandleEvent(ctx, eventMsg)
}

// UserInfo get UserInfo by mid , from cache or db or generate.
func (s *server) UserInfo(ctx context.Context, req *pb.UserInfoReq) (*pb.UserInfoReply, error) {
	ui, err := s.svr.UserInfo(ctx, req.Mid, req.Ip)
	if err != nil {
		return nil, err
	}
	reply := new(pb.UserInfoReply)
	reply.DeepCopyFromUserInfo(ui)
	return reply, nil
}

// UserInfoAsyn get UserInfo by mid , from cache or db or asyn generate.
func (s *server) UserInfoAsyn(ctx context.Context, req *pb.UserInfoAsynReq) (*pb.UserInfoAsynReply, error) {
	ui, err := s.svr.UserInfoAsyn(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	reply := new(pb.UserInfoAsynReply)
	reply.DeepCopyFromUserInfo(ui)
	return reply, nil
}

// ReBuildPortrait reBuild user info.
func (s *server) ReBuildPortrait(ctx context.Context, req *pb.ReBuildPortraitReq) (*pb.ReBuildPortraitReply, error) {
	return &pb.ReBuildPortraitReply{}, s.svr.ReBuildPortrait(ctx, req.Mid, req.Reason)
}

// UpdateUserScore update user score
func (s *server) UpdateUserScore(ctx context.Context, req *pb.UpdateUserScoreReq) (*pb.UpdateUserScoreReply, error) {
	return &pb.UpdateUserScoreReply{}, s.svr.UpdateUserScore(ctx, req.Mid, req.Ip, req.Effect)
}

// RefreshBaseScore refresh base score.
func (s *server) RefreshBaseScore(ctx context.Context, req *pb.RefreshBaseScoreReq) (*pb.RefreshBaseScoreReply, error) {
	argReset := new(model.ArgReset)
	req.DeepCopyAsIntoArgReset(argReset)
	return &pb.RefreshBaseScoreReply{}, s.svr.RefreshBaseScore(ctx, argReset)
}

// UpdateBaseScore update base score.
func (s *server) UpdateBaseScore(ctx context.Context, req *pb.UpdateBaseScoreReq) (*pb.UpdateBaseScoreReply, error) {
	argReset := new(model.ArgReset)
	req.DeepCopyAsIntoArgReset(argReset)
	return &pb.UpdateBaseScoreReply{}, s.svr.UpdateBaseScore(ctx, argReset)
}

// UpdateEventScore update event score.
func (s *server) UpdateEventScore(ctx context.Context, req *pb.UpdateEventScoreReq) (*pb.UpdateEventScoreReply, error) {
	argReset := new(model.ArgReset)
	req.DeepCopyAsIntoArgReset(argReset)
	return &pb.UpdateEventScoreReply{}, s.svr.UpdateEventScore(ctx, argReset)
}

// ClearReliveTimes clear times.
func (s *server) ClearReliveTimes(ctx context.Context, req *pb.ClearReliveTimesReq) (*pb.ClearReliveTimesReply, error) {
	argReset := new(model.ArgReset)
	req.DeepCopyAsIntoArgReset(argReset)
	return &pb.ClearReliveTimesReply{}, s.svr.ClearReliveTimes(ctx, argReset)
}

// Info get user info by mid.
func (s *server) Info(ctx context.Context, req *pb.InfoReq) (*pb.InfoReply, error) {
	ui, err := s.svr.Info(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	reply := new(pb.InfoReply)
	reply.DeepCopyFromUserInfo(ui)
	return reply, nil
}
