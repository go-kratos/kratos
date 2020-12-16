package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	"google.golang.org/grpc"
)

// Server is a gRPC server wrapper.
type Server struct {
	opts        serverOptions
	middlewares map[interface{}]middleware.Middleware
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	options := serverOptions{}
	for _, o := range opts {
		o(&options)
	}
	srv := &Server{
		opts:        options,
		middlewares: make(map[interface{}]middleware.Middleware),
	}
	return srv
}

// ServeGRPC returns a unary server interceptor.
func (s *Server) ServeGRPC() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if m, ok := s.middlewares[info.Server]; ok {
			h = m(h)
		}
		if s.opts.middleware != nil {
			h = s.opts.middleware(h)
		}
		return h(ctx, req)
	}
}

// Use use a middleware to the transport.
func (s *Server) Use(srv interface{}, m ...middleware.Middleware) {
	s.middlewares[srv] = middleware.Chain(m[0], m[1:]...)
}
