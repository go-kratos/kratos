package http

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/server"
)

var _ server.Server = (*Server)(nil)

// Option is HTTP server option.
type Option func(o *options)

// options is HTTP server options.
type options struct {
	handler      http.Handler
	tlsConfig    *tls.Config
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

// Handler with server handler.
func Handler(h http.Handler) Option {
	return func(o *options) {
		o.handler = h
	}
}

// TLSConfig with server tls config.
func TLSConfig(c *tls.Config) Option {
	return func(o *options) {
		o.tlsConfig = c
	}
}

// ReadTimeout with read timeout.
func ReadTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.readTimeout = timeout
	}
}

// WriteTimeout with write timeout.
func WriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.writeTimeout = timeout
	}
}

// IdleTimeout with read timeout.
func IdleTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.idleTimeout = timeout
	}
}

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server

	network string
	addr    string
	opts    options
}

// NewServer creates a HTTP server by options.
func NewServer(network, addr string, opts ...Option) *Server {
	options := options{
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
