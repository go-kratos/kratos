package mcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kratos/kratos/v2/transport"

	"github.com/mark3labs/mcp-go/server"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
)

// MiddlewareFunc is a function that takes an http.Handler and returns an http.Handler.
type MiddlewareFunc func(http.Handler) http.Handler

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Endpoint with server address.
func Endpoint(endpoint *url.URL) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
	}
}

// Middleware with server middleware.
func Middleware(m MiddlewareFunc) ServerOption {
	return func(s *Server) {
		s.middleware = m
	}
}

// SrvOptions with server options.
func SrvOptions(opts ...server.ServerOption) ServerOption {
	return func(s *Server) {
		s.srvOpts = append(s.srvOpts, opts...)
	}
}

// SSEOptions with server SSE options.
func SSEOptions(opts ...server.SSEOption) ServerOption {
	return func(s *Server) {
		s.sseOpts = append(s.sseOpts, opts...)
	}
}

// Server is a MCP server.
type Server struct {
	*server.MCPServer
	srv        *http.Server
	sse        *server.SSEServer
	middleware MiddlewareFunc
	address    string
	endpoint   *url.URL
	srvOpts    []server.ServerOption
	sseOpts    []server.SSEOption
}

// NewServer creates a new MCP server.
func NewServer(name, version string, opts ...ServerOption) *Server {
	srv := &Server{
		address:    ":8000",
		middleware: func(next http.Handler) http.Handler { return next },
	}
	for _, o := range opts {
		o(srv)
	}
	srv.MCPServer = server.NewMCPServer(name, version, srv.srvOpts...)
	srv.srv = &http.Server{Addr: srv.address, Handler: srv.middleware(srv)}
	srv.sse = server.NewSSEServer(srv.MCPServer, append(srv.sseOpts, server.WithHTTPServer(srv.srv))...)
	return srv
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.sse.ServeHTTP(res, req)
}

// Endpoint return a real address to registry endpoint.
// examples:
// - http://127.0.0.1:8000
func (s *Server) Endpoint() (*url.URL, error) {
	if s.endpoint != nil {
		return s.endpoint, nil
	}
	return url.Parse(fmt.Sprintf("http://%s", s.address))
}

// Start start the MCP server.
func (s *Server) Start(_ context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

// Stop stop the MCP server.
func (s *Server) Stop(ctx context.Context) error {
	return s.sse.Shutdown(ctx)
}
