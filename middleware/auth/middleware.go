package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type authkey string

const (
	//HeaderKey holds the key used to store the JWT Token in the request header.
	HeaderKey authkey = "Authorization"

	//InfoKey holds the key used to store the auth info in the context
	InfoKey authkey = "AuthInfo"
)

var (
	ErrWrongContext      = errors.Unauthorized("Something wrong", "Wrong context for middelware")
	ErrNeedTokenProvider = errors.Unauthorized("Missing info", "Token provider is missing")
)

type Option func(options *options)

//options an option setting
type options struct {
	authHeaderKey string
}

//WithAuthHeaderKey set key that hold auth token in header
func WithAuthHeaderKey(headerKey string) Option {
	return func(options *options) {
		options.authHeaderKey = headerKey
	}
}

//TokenParser parse auth token from header
type TokenParser interface {
	ParseToken(token string) (interface{}, error)
}

//TokenProvider provider jwt token
type TokenProvider interface {
	GetToken() string
}

//Server is a server auth middleware
func Server(parser TokenParser, opts ...Option) middleware.Middleware {
	o := initOptions(opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				tokenInfo, err := parser.ParseToken(header.RequestHeader().Get(o.authHeaderKey))
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, InfoKey, tokenInfo)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

//Client is a client jwt middleware
func Client(provider TokenProvider, opts ...Option) middleware.Middleware {
	o := initOptions(opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if provider == nil {
				return nil, ErrNeedTokenProvider
			}
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(o.authHeaderKey, provider.GetToken())
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

//initOptions init the option
func initOptions(opts ...Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.authHeaderKey == "" {
		o.authHeaderKey = "Authorization"
	}
	return o
}
