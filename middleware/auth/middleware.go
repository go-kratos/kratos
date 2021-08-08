package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

const (
	//HeaderKey holds the key used to store the JWT Token in the request header.
	HeaderKey string = "Authorization"

	//InfoKey holds the key used to store the auth info in the context
	InfoKey string = "AuthInfo"
)

var (
	ErrWrongContext      = errors.Unauthorized("Something wrong", "Wrong context for middelware")
	ErrNeedTokenProvider = errors.Unauthorized("Missing info", "Token provider is missing")
)

//TokenParser parse auth token from header
type TokenParser interface {
	ParseToken(token string) (interface{}, error)
}

//TokenProvider provider jwt token
type TokenProvider interface {
	GetToken() string
}

//Server is a server auth middleware
func Server(parser TokenParser) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				tokenInfo, err := parser.ParseToken(header.RequestHeader().Get(HeaderKey))
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, InfoKey, tokenInfo)
				return handler(ctx, req)
			} else {
				return nil, ErrWrongContext
			}
		}
	}
}

//Client is a client jwt middleware
func Client(provider TokenProvider) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if provider == nil {
				return nil, ErrNeedTokenProvider
			}
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(HeaderKey, provider.GetToken())
				return handler(ctx, req)
			} else {
				return nil, ErrWrongContext
			}
		}
	}
}
