package middleware

import (
	"context"
)

// Handler defines the handler invoked by Middleware.
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// ServerInfo is the server request infomation.
type ServerInfo interface {
	Server() interface{}
	Path() string
	Method() string
}

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(ctx context.Context, req interface{}, info *ServerInfo, handler Handler) (resp interface{}, err error)
