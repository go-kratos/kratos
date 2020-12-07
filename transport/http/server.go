package http

import (
	"context"
	"net"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
)

var _ transport.Server = new(Server)

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server

	opts serverOptions
}

// NewServer creates a HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	options := serverOptions{}
	for _, o := range opts {
		o(&options)
	}
	srv := &http.Server{}
	return &Server{Server: srv, opts: options}
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

// Stop stop the HTT.
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
