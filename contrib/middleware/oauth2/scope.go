package oauth2

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
)

// RequireScopes is a server middleware that requires the token to be granted
// all the given scopes. It must be used after Server.
//
// Combine it with selector to require scopes per operation:
//
//	selector.Server(oauth2.RequireScopes("order.write")).
//		Path(v1.OperationOrderCreateOrder).Build()
func RequireScopes(scopes ...string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			info, ok := FromContext(ctx)
			if !ok {
				return nil, ErrMissingToken
			}
			if !info.HasScopes(scopes...) {
				if tr, ok := transport.FromServerContext(ctx); ok {
					setChallenge(tr, "", challengeInsufficientScope, strings.Join(scopes, " "))
				}
				return nil, ErrInsufficientScope
			}
			return handler(ctx, req)
		}
	}
}
