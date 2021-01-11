package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"google.golang.org/grpc"
)

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(ctx context.Context, err error) error

// ServerMiddleware with server middleware.
func ServerMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.middleware = middleware.Chain(m[0], m[1:]...)
	}
}

// Server is a gRPC server wrapper.
type Server struct {
	middleware       middleware.Middleware
	serverMiddleware map[interface{}]middleware.Middleware
	errorEncoder     EncodeErrorFunc
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		errorEncoder:     DefaultErrorEncoder,
		serverMiddleware: make(map[interface{}]middleware.Middleware),
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

// Use use a middleware to the transport.
func (s *Server) Use(srv interface{}, m ...middleware.Middleware) {
	s.serverMiddleware[srv] = middleware.Chain(m[0], m[1:]...)
}

// Interceptor returns a unary server interceptor.
func (s *Server) Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = transport.NewContext(ctx, transport.Transport{Kind: "GRPC"})
		ctx = NewContext(ctx, ServerInfo{Server: info.Server, FullMethod: info.FullMethod})
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if m, ok := s.serverMiddleware[info.Server]; ok {
			h = m(h)
		}
		if s.middleware != nil {
			h = s.middleware(h)
		}
		resp, err := h(ctx, req)
		if err != nil {
			return nil, s.errorEncoder(ctx, err)
		}
		return resp, nil
	}
}
