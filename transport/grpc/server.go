package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/internal/endpoint"
	"github.com/go-kratos/kratos/v2/internal/matcher"

	apimd "github.com/go-kratos/kratos/v2/api/metadata"

	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc/reflection"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

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
// Deprecated: use global logger instead.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {}
}

// Middleware with server middleware.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middleware.Use(m...)
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func StreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}

// Options with grpc options.
func Options(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	baseCtx    context.Context
	tlsConf    *tls.Config
	lis        net.Listener
	err        error
	network    string
	address    string
	endpoint   *url.URL
	timeout    time.Duration
	middleware matcher.Matcher
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
	health     *health.Server
	metadata   *apimd.Server
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx:    context.Background(),
		network:    "tcp",
		address:    ":0",
		timeout:    1 * time.Second,
		health:     health.NewServer(),
		middleware: matcher.New(),
	}
	for _, o := range opts {
		o(srv)
	}
	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	streamInts := []grpc.StreamServerInterceptor{
		srv.streamServerInterceptor(),
	}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	srv.metadata = apimd.NewServer(srv.Server)
	// internal register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)
	reflection.Register(srv.Server)
	return srv
}

// Use uses a service middleware with selector.
// selector:
//   - '/*'
//   - '/helloworld.v1.Greeter/*'
//   - '/helloworld.v1.Greeter/SayHello'
func (s *Server) Use(selector string, m ...middleware.Middleware) {
	s.middleware.Add(selector, m...)
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.err
	}
	s.baseCtx = ctx
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.health.Shutdown()
	s.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}
	return s.err
}
