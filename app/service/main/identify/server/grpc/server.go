package grpc

import (
	"context"
	"fmt"

	"go-common/app/service/main/identify/api/grpc"
	"go-common/app/service/main/identify/service"
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
		if resp, err = handler(ctx, req); err != nil {
			log.Infov(ctx,
				log.KV("path", info.FullMethod),
				log.KV("caller", metadata.String(ctx, metadata.Caller)),
				log.KV("args", fmt.Sprintf("%v", req)),
				log.KV("args", fmt.Sprintf("%+v", err)))
		}
		return
	})
	v1.RegisterIdentifyServer(w.Server(), &server{s})
	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.IdentifyServer = &server{}
var (
	emptyCookieReply = &v1.GetCookieInfoReply{
		IsLogin: false,
	}

	emptyTokenReply = &v1.GetTokenInfoReply{
		IsLogin: false,
	}
)

// CookieInfo verify user info by cookie.
func (s *server) GetCookieInfo(ctx context.Context, req *v1.GetCookieInfoReq) (*v1.GetCookieInfoReply, error) {
	res, err := s.svr.GetCookieInfo(ctx, req.GetCookie())
	if err != nil {
		if err == ecode.NoLogin {
			return emptyCookieReply, nil
		}
		return nil, err
	}

	return &v1.GetCookieInfoReply{
		IsLogin: true,
		Mid:     res.Mid,
		Expires: res.Expires,
		Csrf:    res.Csrf,
	}, nil
}

// TokenInfo verify user info by token.
func (s *server) GetTokenInfo(ctx context.Context, req *v1.GetTokenInfoReq) (*v1.GetTokenInfoReply, error) {
	token := &v1.GetTokenInfoReq{
		Buvid: req.Buvid,
		Token: req.Token,
	}
	res, err := s.svr.GetTokenInfo(ctx, token)
	if err != nil {
		if err == ecode.NoLogin {
			return emptyTokenReply, nil
		}
		return nil, err
	}
	return &v1.GetTokenInfoReply{
		IsLogin: true,
		Mid:     res.Mid,
		Expires: res.Expires,
		Csrf:    res.Csrf,
	}, nil
}
