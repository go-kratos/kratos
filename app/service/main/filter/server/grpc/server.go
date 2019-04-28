package grpc

import (
	"context"
	"fmt"

	"go-common/app/service/main/filter/api/grpc/v1"
	"go-common/app/service/main/filter/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// New Identify warden rpc server
func New(cfg *warden.ServerConfig, s *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	w.Use(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if resp, err = handler(ctx, req); err == nil {
			log.Infov(ctx,
				log.KV("path", info.FullMethod),
				log.KV("caller", metadata.String(ctx, metadata.Caller)),
				log.KV("remote_ip", metadata.String(ctx, metadata.RemoteIP)),
				log.KV("args", fmt.Sprintf("%s", req)),
				log.KV("retVal", fmt.Sprintf("%s", resp)))
		}
		return
	})
	v1.RegisterFilterServer(w.Server(), &server{s})
	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.FilterServer = &server{}

// Filter filter msg
func (s *server) Filter(ctx context.Context, req *v1.FilterReq) (*v1.FilterReply, error) {
	fmsg, level, typeIds, hitRules, limit, ai, err := s.svr.Filter(ctx, req.Area, req.Message, req.TypeId, req.Id, req.Oid, req.Mid, req.Keys, int8(req.ReplyType))
	if err != nil {
		return nil, err
	}
	if typeIds == nil {
		typeIds = []int64{}
	}
	if hitRules == nil {
		hitRules = []string{}
	}
	modelAi := new(v1.ModelAiScore)
	if ai != nil {
		modelAi.Scores = ai.Scores
		modelAi.Threshold = ai.Threshold
		modelAi.Note = ai.Note
	}
	return &v1.FilterReply{
		Result:   fmsg,
		Level:    int32(level),
		Limit:    int64(limit),
		TypeIds:  typeIds,
		HitRules: hitRules,
		Ai:       modelAi,
	}, nil
}

// MFilter  filter batch msgs
func (s *server) MFilter(ctx context.Context, req *v1.MFilterReq) (*v1.MFilterReply, error) {
	res, err := s.svr.RPCMultiFilter(ctx, req.Area, req.MsgMap, req.TypeId)
	if err != nil {
		return nil, err
	}
	replys := make(map[string]*v1.FilterReply)
	for k, v := range res {
		r := &v1.FilterReply{
			Result: v.Result,
			Level:  int32(v.Level),
			Limit:  int64(v.Limit),
		}
		replys[k] = r
	}
	return &v1.MFilterReply{
		RMap: replys,
	}, nil
}

// Hit return hit words
func (s *server) Hit(ctx context.Context, req *v1.HitReq) (*v1.HitReply, error) {
	res, err := s.svr.Hit(ctx, req.Area, req.Msg, req.TypeId)
	if err != nil {
		return nil, err
	}
	return &v1.HitReply{
		Hits: res,
	}, nil
}

// MHit multi Hit api (max multi <= 20 , max Bytes per msg <= 300)
func (s *server) MHit(ctx context.Context, req *v1.MHitReq) (*v1.MHitReply, error) {
	// 1. validate req
	if len(req.MsgMap) > 20 {
		return nil, ecode.RequestErr
	}
	for _, msg := range req.MsgMap {
		if len(msg) > 300 {
			return nil, ecode.RequestErr
		}
	}
	// 2. do
	res, err := s.svr.MultiHit(ctx, req.Area, req.MsgMap, req.TypeId)
	if err != nil {
		return nil, err
	}
	replys := make(map[string]*v1.HitReply)
	for k, v := range res {
		r := &v1.HitReply{
			Hits: v,
		}
		replys[k] = r
	}
	return &v1.MHitReply{
		RMap: replys,
	}, nil
}
