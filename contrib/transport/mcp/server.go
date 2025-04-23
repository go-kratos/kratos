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

// Server is a MCP server.
type Server struct {
	*server.MCPServer
	sse      *server.SSEServer
	address  string
	endpoint *url.URL
	mcpOpts  []server.ServerOption
	sseOpts  []server.SSEOption
}

// NewServer creates a new MCP server.
func NewServer(name, version string, opts ...ServerOption) *Server {
	srv := &Server{
		address: ":7000",
	}
	for _, o := range opts {
		o(srv)
	}
	srv.MCPServer = server.NewMCPServer(name, version, srv.mcpOpts...)
	srv.sse = server.NewSSEServer(srv.MCPServer, srv.sseOpts...)
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
	if err := s.sse.Start(s.address); err != nil {
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
