package rpcx

import (
	"context"
	"crypto/tls"
	"github.com/go-kratos/kratos/v2/internal/endpoint"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/transport"

	rpcx "github.com/smallnest/rpcx/server"
)

var _ transport.Server = (*Server)(nil)
var _ transport.Endpointer = (*Server)(nil)

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// Middleware with server middleware.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middleware = m
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

// Options with RPCx options.
func Options(opts ...rpcx.OptionFn) ServerOption {
	return func(s *Server) {
		s.rpcxOpts = opts
	}
}

// Plugin with RPCx Plugin.
func Plugin(plugin ...rpcx.Plugin) ServerOption {
	return func(s *Server) {
		s.rpcxPlugins = plugin
	}
}

type Server struct {
	*rpcx.Server
	ctx         context.Context
	tlsConf     *tls.Config
	lis         net.Listener
	once        sync.Once
	err         error
	network     string
	address     string
	endpoint    *url.URL
	timeout     time.Duration
	log         *log.Helper
	rpcxOpts    []rpcx.OptionFn
	middleware  []middleware.Middleware
	rpcxPlugins []rpcx.Plugin
}

// NewServer creates a RPCx server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		log:     log.NewHelper(log.DefaultLogger),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.Server = rpcx.NewServer(srv.rpcxOpts...)
	if len(srv.rpcxPlugins) > 0 {
		for _, plugin := range srv.rpcxPlugins {
			srv.Server.Plugins.Add(plugin)
		}
	}
	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   rpcx://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, lis)
		_ = lis.Close()
		if err != nil {
			s.err = err
			return
		}
		s.endpoint = endpoint.NewEndpoint("rpcx", addr, s.tlsConf != nil)
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start the RPCx server.
func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	s.log.Infof("[RPCx] server listening on: %s", s.address)
	return s.Serve(s.network, s.address)
}

// Stop the RPCx server.
func (s *Server) Stop(ctx context.Context) error {
	_ = s.Close()
	s.log.Info("[RPCx] server stopping")
	return nil
}
