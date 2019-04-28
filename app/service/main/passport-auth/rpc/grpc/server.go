// Package grpc server generate by warden_gen
package grpc

import (
	"context"

	"go-common/app/service/main/passport-auth/api/grpc/v1"
	"go-common/app/service/main/passport-auth/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// New Auth warden rpc server
func New(cfg *warden.ServerConfig, s *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	v1.RegisterAuthServer(w.Server(), &server{s})
	w.Use(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if resp, err = handler(ctx, req); err == nil {
			log.Infov(ctx,
				log.KV("path", info.FullMethod),
				log.KV("caller", metadata.String(ctx, metadata.Caller)),
				log.KV("args", req), log.KV("retVal", resp))
		}
		return
	})
	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.AuthServer = &server{}
var (
	emptyCookieReply = &v1.GetCookieInfoReply{
		IsLogin: false,
	}

	emptyTokenReply = &v1.GetTokenInfoReply{
		IsLogin: false,
	}
)

// CookieInfo verify user info by cookie.
func (s *server) GetCookieInfo(c context.Context, req *v1.GetCookieInfoReq) (*v1.GetCookieInfoReply, error) {
	res, err := s.svr.CookieInfo(c, req.GetCookie())
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
		Csrf:    res.CSRF,
	}, nil
}

// TokenInfo verify user info by accesskey.
func (s *server) GetTokenInfo(c context.Context, req *v1.GetTokenInfoReq) (*v1.GetTokenInfoReply, error) {
	res, err := s.svr.TokenInfo(c, req.GetToken())
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
		Csrf:    res.CSRF,
	}, nil
}
