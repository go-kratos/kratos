package grpc

import (
	"context"

	"go-common/app/service/main/share/api"
	"go-common/app/service/main/share/model"
	"go-common/app/service/main/share/service"
	"go-common/library/net/rpc/warden"
)

// server .
type server struct {
	srv *service.Service
}

// New share warden rpc server.
func New(cfg *warden.ServerConfig, srv *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	v1.RegisterShareServer(w.Server(), &server{srv: srv})
	var err error
	w, err = w.Start()
	if err != nil {
		panic(err)
	}
	return w
}

// AddShare .
func (s *server) AddShare(ctx context.Context, req *v1.AddShareRequest) (*v1.AddShareReply, error) {
	p := &model.ShareParams{
		OID: req.Oid,
		MID: req.Mid,
		TP:  int(req.Type),
		IP:  req.Ip,
	}
	shares, err := s.srv.Add(ctx, p)
	if err != nil {
		return nil, err
	}
	addShareReply := &v1.AddShareReply{
		Shares: shares,
	}
	return addShareReply, nil
}
