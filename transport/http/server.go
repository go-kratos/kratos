package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gorilla/mux"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

var _ transport.Server = new(Server)

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server

	router      *mux.Router
	opts        serverOptions
	middlewares map[interface{}]middleware.Middleware
}

// NewServer creates a HTTP server by options.
func NewServer(addr string, opts ...ServerOption) *Server {
	options := serverOptions{
		errorHandler: DefaultErrorHandler,
	}
	for _, o := range opts {
		o(&options)
	}
	router := mux.NewRouter()
	return &Server{
		opts:   options,
		router: router,
		Server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		middlewares: make(map[interface{}]middleware.Middleware),
	}
}

func (s *Server) Use(srv interface{}, m middleware.Middleware) {
	s.middlewares[srv] = m
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
	return s.ListenAndServe()
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
