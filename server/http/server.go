package http

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/server"
)

var _ server.Server = new(Server)

// ServerOption is HTTP server option.
type ServerOption func(o *serverOptions)

// serverOptions is HTTP server options.
type serverOptions struct {
	handler      http.Handler
	tlsConfig    *tls.Config
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// ServerHandler with server handler.
func ServerHandler(h http.Handler) ServerOption {
	return func(o *serverOptions) {
		o.handler = h
	}
}

// ServerTLSConfig with server tls config.
func ServerTLSConfig(c *tls.Config) ServerOption {
	return func(o *serverOptions) {
		o.tlsConfig = c
	}
}

// ServerReadTimeout with read timeout.
func ServerReadTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.readTimeout = timeout
	}
}

// ServerWriteTimeout with write timeout.
func ServerWriteTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.writeTimeout = timeout
	}
}

// ServerIdleTimeout with read timeout.
func ServerIdleTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.idleTimeout = timeout
	}
}

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server

	network string
	addr    string
	opts    serverOptions
}

// NewServer creates a HTTP server by options.
func NewServer(network, addr string, opts ...ServerOption) *Server {
	options := serverOptions{
		readTimeout:  time.Second,
		writeTimeout: time.Second,
		idleTimeout:  time.Minute,
	}
	for _, o := range opts {
		o(&options)
	}
	return &Server{
		network: network,
		addr:    addr,
		opts:    options,
		Server: &http.Server{
			Handler:      options.handler,
			TLSConfig:    options.tlsConfig,
			ReadTimeout:  options.readTimeout,
			WriteTimeout: options.writeTimeout,
			IdleTimeout:  options.idleTimeout,
		},
	}
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
