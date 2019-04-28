// Package server generate by warden_gen
package server

import (
	"context"

	pb "go-common/app/service/main/point/api"
	"go-common/app/service/main/point/model"
	service "go-common/app/service/main/point/service"
	"go-common/library/net/rpc/warden"
)

// New Point warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterPointServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ pb.PointServer = &server{}

// Config get point config.
func (s *server) Config(ctx context.Context, req *pb.ConfigReq) (*pb.ConfigReply, error) {
	p, err := s.svr.Config(ctx, int(req.ChangeType), req.Mid, req.Bp)
	if err != nil {
		return nil, err
	}
	return &pb.ConfigReply{Point: p}, err
}

// AllConfig all point config
func (s *server) AllConfig(ctx context.Context, req *pb.AllConfigReq) (*pb.AllConfigReply, error) {
	ac := s.svr.AllConfig(ctx)
	return &pb.AllConfigReply{Data_0: ac}, nil
}

// Ping Service
func (s *server) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingReply, error) {
	return &pb.PingReply{}, nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *pb.CloseReq) (*pb.CloseReply, error) {
	return &pb.CloseReply{}, nil
}

// PointInfo .
func (s *server) PointInfo(ctx context.Context, req *pb.PointInfoReq) (*pb.PointInfoReply, error) {
	p, err := s.svr.PointInfo(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.PointInfoReply{Pi: &pb.ModelPointInfo{
		Mid:          p.Mid,
		PointBalance: p.PointBalance,
		Ver:          p.Ver,
	}}, err
}

// PointHistory .
func (s *server) PointHistory(ctx context.Context, req *pb.PointHistoryReq) (*pb.PointHistoryReply, error) {
	phs, t, nc, err := s.svr.PointHistory(ctx, req.Mid, int(req.Cursor), int(req.Ps))
	if err != nil {
		return nil, err
	}
	var mph []*pb.ModelPointHistory
	for _, v := range phs {
		p := &pb.ModelPointHistory{
			Id:           v.ID,
			Mid:          v.Mid,
			Point:        v.Point,
			OrderId:      v.OrderID,
			ChangeType:   int32(v.ChangeType),
			ChangeTime:   v.ChangeTime.Time().Unix(),
			RelationId:   v.RelationID,
			PointBalance: v.PointBalance,
			Remark:       v.Remark,
			Operator:     v.Operator,
		}
		mph = append(mph, p)
	}
	return &pb.PointHistoryReply{
		Phs:     mph,
		Total:   int32(t),
		Ncursor: int32(nc),
	}, err
}

// OldPointHistory old point history .
func (s *server) OldPointHistory(ctx context.Context, req *pb.OldPointHistoryReq) (*pb.OldPointHistoryReply, error) {
	phs, t, err := s.svr.OldPointHistory(ctx, req.Mid, int(req.Pn), int(req.Ps))
	if err != nil {
		return nil, err
	}
	var mop []*pb.ModelOldPointHistory
	for _, v := range phs {
		p := &pb.ModelOldPointHistory{
			Id:           v.ID,
			Mid:          v.Mid,
			Point:        v.Point,
			OrderId:      v.OrderID,
			ChangeType:   int32(v.ChangeType),
			ChangeTime:   v.ChangeTime,
			RelationId:   v.RelationID,
			PointBalance: v.PointBalance,
			Remark:       v.Remark,
			Operator:     v.Operator,
		}
		mop = append(mop, p)
	}
	return &pb.OldPointHistoryReply{Phs: mop, Total: int32(t)}, err
}

// PointAddByBp by bp.
func (s *server) PointAddByBp(ctx context.Context, req *pb.PointAddByBpReq) (*pb.PointAddByBpReply, error) {
	arg := &model.ArgPointAdd{
		Mid:        req.Pa.Mid,
		ChangeType: int(req.Pa.ChangeType),
		RelationID: req.Pa.RelationId,
		Bcoin:      req.Pa.Bcoin,
		Remark:     req.Pa.Remark,
		OrderID:    req.Pa.OrderId,
	}
	p, err := s.svr.PointAddByBp(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &pb.PointAddByBpReply{P: p}, err
}

// ConsumePoint .
func (s *server) ConsumePoint(ctx context.Context, req *pb.ConsumePointReq) (*pb.ConsumePointReply, error) {
	arg := &model.ArgPointConsume{
		Mid:        req.Pc.Mid,
		ChangeType: req.Pc.ChangeType,
		RelationID: req.Pc.RelationId,
		Point:      req.Pc.Point,
		Remark:     req.Pc.Remark,
	}
	status, err := s.svr.ConsumePoint(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &pb.ConsumePointReply{Status: int32(status)}, err
}

// AddPoint .
func (s *server) AddPoint(ctx context.Context, req *pb.AddPointReq) (*pb.AddPointReply, error) {
	pc := &model.ArgPoint{
		Mid:        req.Pc.Mid,
		ChangeType: req.Pc.ChangeType,
		Point:      req.Pc.Point,
		Remark:     req.Pc.Remark,
		Operator:   req.Pc.Operator,
	}
	status, err := s.svr.AddPoint(ctx, pc)
	if err != nil {
		return nil, err
	}
	return &pb.AddPointReply{Status: int32(status)}, err
}
