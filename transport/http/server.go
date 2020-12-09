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

	handlers []http.Handler

	opts serverOptions
}

// NewServer creates a HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	options := serverOptions{}
	for _, o := range opts {
		o(&options)
	}
	s := &Server{opts: options}
	s.Server = &http.Server{Handler: s}
	return s
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

// AddHandler add a HTTP handler.
func (s *Server) AddHandler(h http.Handler) {
	s.handlers = append(s.handlers, h)
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, h := range s.handlers {
		h.ServeHTTP(res, req)
	}
}
