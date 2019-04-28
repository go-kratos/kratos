// Package grpc generate by warden_gen
package grpc

import (
	"context"

	pb "go-common/app/interface/main/history/api/grpc"
	"go-common/app/interface/main/history/model"
	service "go-common/app/interface/main/history/service"
	"go-common/library/net/rpc/warden"
)

// New History warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterHistoryServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ pb.HistoryServer = &server{}

// AddHistory add hisotry progress into hbase.
func (s *server) AddHistory(ctx context.Context, req *pb.AddHistoryReq) (*pb.AddHistoryReply, error) {
	tp, err := model.MustCheckBusiness(req.H.Business)
	if err != nil {
		return nil, err
	}
	h := &model.History{
		Mid:      req.H.Mid,
		Aid:      req.H.Aid,
		Sid:      req.H.Sid,
		Epid:     req.H.Epid,
		TP:       tp,
		Business: req.H.Business,
		STP:      int8(req.H.Stp),
		Cid:      req.H.Cid,
		DT:       int8(req.H.Dt),
		Pro:      req.H.Pro,
		Unix:     req.H.Unix,
	}
	return nil, s.svr.AddHistory(ctx, req.Mid, req.Rtime, h)
}

// Progress get view progress from cache/hbase.
func (s *server) Progress(ctx context.Context, req *pb.ProgressReq) (*pb.ProgressReply, error) {
	histories, err := s.svr.Progress(ctx, req.Mid, req.Aids)
	if err != nil {
		return nil, err
	}
	reply := &pb.ProgressReply{Res: make(map[int64]*pb.ModelHistory)}
	for k, v := range histories {
		reply.Res[k] = &pb.ModelHistory{
			Mid:      v.Mid,
			Aid:      v.Aid,
			Sid:      v.Sid,
			Epid:     v.Epid,
			Business: v.Business,
			Stp:      int32(v.STP),
			Cid:      v.Cid,
			Dt:       int32(v.DT),
			Pro:      v.Pro,
			Unix:     v.Unix,
		}
	}
	return reply, nil
}

// Position get view progress from cache/hbase.
func (s *server) Position(ctx context.Context, req *pb.PositionReq) (*pb.PositionReply, error) {
	tp, err := model.MustCheckBusiness(req.Business)
	if err != nil {
		return nil, err
	}
	h, err := s.svr.Position(ctx, req.Mid, req.Aid, tp)
	if err != nil {
		return nil, err
	}
	reply := &pb.PositionReply{
		Res: &pb.ModelHistory{
			Mid:      h.Mid,
			Aid:      h.Aid,
			Sid:      h.Sid,
			Epid:     h.Epid,
			Business: h.Business,
			Stp:      int32(h.STP),
			Cid:      h.Cid,
			Dt:       int32(h.DT),
			Pro:      h.Pro,
			Unix:     h.Unix,
		},
	}
	return reply, err
}

// ClearHistory clear user's historys.
func (s *server) ClearHistory(ctx context.Context, req *pb.ClearHistoryReq) (*pb.ClearHistoryReply, error) {
	var tps []int8
	for _, b := range req.Businesses {
		tp, err := model.MustCheckBusiness(b)
		if err != nil {
			return nil, err
		}
		tps = append(tps, tp)
	}
	return nil, s.svr.ClearHistory(ctx, req.Mid, tps)
}

// Histories return the user all av  history.
func (s *server) Histories(ctx context.Context, req *pb.HistoriesReq) (*pb.HistoriesReply, error) {
	tp, err := model.MustCheckBusiness(req.Business)
	if err != nil {
		return nil, err
	}
	resources, err := s.svr.Histories(ctx, req.Mid, tp, int(req.Pn), int(req.Ps))
	if err != nil {
		return nil, err
	}
	reply := &pb.HistoriesReply{}
	for _, v := range resources {
		reply.Res = append(reply.Res,
			&pb.ModelResource{
				Mid:      v.Mid,
				Oid:      v.Oid,
				Sid:      v.Sid,
				Epid:     v.Epid,
				Business: v.Business,
				Stp:      int32(v.STP),
				Cid:      v.Cid,
				Dt:       int32(v.DT),
				Pro:      v.Pro,
				Unix:     v.Unix,
			})
	}
	return reply, nil
}

// HistoryCursor return the user all av  history.
func (s *server) HistoryCursor(ctx context.Context, req *pb.HistoryCursorReq) (*pb.HistoryCursorReply, error) {
	tp, err := model.MustCheckBusiness(req.Business)
	if err != nil {
		return nil, err
	}
	var tps []int8
	for _, b := range req.Businesses {
		t, er := model.MustCheckBusiness(b)
		if er != nil {
			return nil, er
		}
		tps = append(tps, t)
	}
	resources, err := s.svr.HistoryCursor(ctx, req.Mid, req.Max, req.ViewAt, int(req.Ps), tp, tps, req.Ip)
	if err != nil {
		return nil, err
	}
	reply := &pb.HistoryCursorReply{}
	for _, v := range resources {
		reply.Res = append(reply.Res,
			&pb.ModelResource{
				Mid:      v.Mid,
				Oid:      v.Oid,
				Sid:      v.Sid,
				Epid:     v.Epid,
				Business: v.Business,
				Stp:      int32(v.STP),
				Cid:      v.Cid,
				Dt:       int32(v.DT),
				Pro:      v.Pro,
				Unix:     v.Unix,
			})
	}
	return reply, nil
}

// Delete .
func (s *server) Delete(ctx context.Context, req *pb.DeleteReq) (*pb.DeleteReply, error) {
	var his []*model.History
	for _, b := range req.His {
		tp, err := model.MustCheckBusiness(b.Business)
		if err != nil {
			return nil, err
		}
		h := &model.History{
			Mid:      b.Mid,
			Aid:      b.Aid,
			Sid:      b.Sid,
			Epid:     b.Epid,
			TP:       tp,
			Business: b.Business,
			STP:      int8(b.Stp),
			Cid:      b.Cid,
			DT:       int8(b.Dt),
			Pro:      b.Pro,
			Unix:     b.Unix,
		}
		his = append(his, h)
	}
	return nil, s.svr.Delete(ctx, req.Mid, his)
}

// FlushHistory flush to hbase from cache.
func (s *server) FlushHistory(ctx context.Context, req *pb.FlushHistoryReq) (*pb.FlushHistoryReply, error) {
	return nil, s.svr.FlushHistory(ctx, req.Mids, req.Stime)
}
