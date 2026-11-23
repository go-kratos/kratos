package oauth2

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
)

// Client is a client oauth2 middleware. It fetches an access token from the
// token source and puts it into the Authorization request header. It works
// with both HTTP and gRPC transports.
//
// The token source is wrapped by oauth2.ReuseTokenSource, so the token is
// cached and only refreshed when it expires.
func Client(src oauth2.TokenSource) middleware.Middleware {
	if src != nil {
		src = oauth2.ReuseTokenSource(nil, src)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if src == nil {
				return nil, ErrMissingTokenSource
			}
			tr, ok := transport.FromClientContext(ctx)
			if !ok {
				return nil, ErrWrongContext
			}
			token, err := src.Token()
			if err != nil {
				return nil, ErrTokenSource.WithCause(err)
			}
			tr.RequestHeader().Set(authorizationKey, token.Type()+" "+token.AccessToken)
			return handler(ctx, req)
		}
	}
}
