package grpc

import (
	"context"
	"net"

	"github.com/go-kratos/kratos/v2/transport"

	"google.golang.org/grpc"
)

var _ transport.Server = new(Server)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server

	opts serverOptions
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	options := serverOptions{}
	for _, o := range opts {
		o(&options)
	}
	srv := grpc.NewServer()
	return &Server{Server: srv, opts: options}
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.opts.Address)
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
