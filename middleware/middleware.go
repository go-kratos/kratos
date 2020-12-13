package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
)

// Handler defines the handler invoked by Middleware.
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(ctx context.Context, req interface{}, info *transport.ServerInfo, handler Handler) (resp interface{}, err error)
