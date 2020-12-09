package http

import (
	"context"
	"net"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gorilla/mux"
)

var _ transport.Server = new(Server)

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server

	router *mux.Router

	opts ServerOptions
}

// NewServer creates a HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	options := ServerOptions{
		ErrorHandler: DefaultErrorHandler,
	}
	for _, o := range opts {
		o(&options)
	}
	router := mux.NewRouter()
	return &Server{
		Server: &http.Server{
			Handler: router,
		},
		router: router,
		opts:   options,
	}
}

// Handle registers a new route with a matcher for the URL path.
func (s *Server) Handle(path string, handler http.Handler) {
	s.router.Handle(path, handler)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (s *Server) HandleFunc(path string, h func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(path, h)
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
