package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/api/metadata"
	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
		s.middleware = middleware.Chain(m...)
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
	lis        net.Listener
	network    string
	address    string
	timeout    time.Duration
	log        *log.Helper
	middleware middleware.Middleware
	grpcOpts   []grpc.ServerOption
	health     *health.Server
	metadata   *metadata.Server
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: time.Second,
		middleware: middleware.Chain(
			recovery.Recovery(),
		),
		health: health.NewServer(),
		log:    log.NewHelper(log.DefaultLogger),
	}
	for _, o := range opts {
		o(srv)
	}
	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			unaryServerInterceptor(srv.middleware, srv.timeout),
		),
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	srv.metadata = metadata.NewServer(srv.Server)
	// internal register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	metadata.RegisterMetadataServer(srv.Server, srv.metadata)
	reflection.Register(srv.Server)
	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	if s.lis == nil && strings.HasSuffix(s.address, ":0") {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return "", err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
}

// Start start the gRPC server.
func (s *Server) Start() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	s.log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop() error {
	s.GracefulStop()
	s.health.Shutdown()
	s.log.Info("[gRPC] server stopping")
	return nil
}

func unaryServerInterceptor(m middleware.Middleware, timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindGRPC})
		ctx = NewServerContext(ctx, ServerInfo{Server: info.Server, FullMethod: info.FullMethod})
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if m != nil {
			h = m(h)
		}
		return h(ctx, req)
	}
}
