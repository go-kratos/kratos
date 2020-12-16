package grpc

import (
	"context"
	"net"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"google.golang.org/grpc"
)

var _ transport.Server = new(Server)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server

	addr        string
	opts        serverOptions
	middlewares map[interface{}]middleware.Middleware
}

// NewServer creates a gRPC server by options.
func NewServer(addr string, opts ...ServerOption) *Server {
	options := serverOptions{}
	for _, o := range opts {
		o(&options)
	}
	srv := &Server{
		addr:        addr,
		opts:        options,
		middlewares: make(map[interface{}]middleware.Middleware),
	}
	srv.Server = grpc.NewServer(append(options.grpcOpts, grpc.UnaryInterceptor(srv.interceptor()))...)
	return srv
}

func (s *Server) interceptor() grpc.UnaryServerInterceptor {
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

// Use .
func (s *Server) Use(srv interface{}, m middleware.Middleware) {
	s.middlewares[srv] = m
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	return nil
}
